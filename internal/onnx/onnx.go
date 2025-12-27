// Package onnx provides functionality to extract and manage the ONNX Runtime shared library.
// The library is embedded as a compressed archive (tgz for Unix-like systems, zip for Windows)
// and extracted to the user's data directory on first run.
package onnx

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/eduardolat/exp-stt/internal/config"
	"github.com/eduardolat/exp-stt/internal/logger"
)

// SharedLibraryPath holds the absolute path to the extracted ONNX Runtime shared library.
// This value is set by EnsureSharedLibrary after successful extraction.
var SharedLibraryPath = ""

// EnsureSharedLibrary extracts the ONNX Runtime shared library from the embedded archive
// if it doesn't already exist. It sets SharedLibraryPath to the location of the extracted library.
func EnsureSharedLibrary(logger logger.Logger) error {
	extractDir := filepath.Join(config.DirectoryOnnxRuntime, runtimeVersion, runtimePlatform)
	SharedLibraryPath = filepath.Join(extractDir, "lib", sharedLibName)

	logger.Debug(
		context.Background(), "ensuring ONNX Runtime shared library",
		"shared_library_path", SharedLibraryPath,
	)

	if fileExists(SharedLibraryPath) {
		logger.Debug(context.Background(), "ONNX Runtime shared library already exists, skipping extraction")
		return nil
	}

	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return fmt.Errorf("creating extraction directory %s: %w", extractDir, err)
	}

	if isTgz {
		if err := extractTgz(CompressedLib, extractDir); err != nil {
			return fmt.Errorf("extracting tgz archive: %w", err)
		}
		logger.Debug(context.Background(), "ONNX Runtime shared library extracted from tgz archive")
		return nil
	}

	if isZip {
		if err := extractZip(CompressedLib, extractDir); err != nil {
			return fmt.Errorf("extracting zip archive: %w", err)
		}
		logger.Debug(context.Background(), "ONNX Runtime shared library extracted from zip archive")
		return nil
	}

	return fmt.Errorf("unknown archive format: neither tgz nor zip")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// extractTgz extracts a gzipped tar archive to the destination directory.
// It strips the top-level directory from the archive paths.
func extractTgz(data []byte, destDir string) error {
	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("creating gzip reader: %w", err)
	}
	defer func() { _ = gzReader.Close() }()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading tar header: %w", err)
		}

		relativePath := stripTopLevelDir(header.Name)
		if relativePath == "" {
			continue
		}

		targetPath := filepath.Join(destDir, relativePath)

		if err := extractTarEntry(header, tarReader, targetPath); err != nil {
			return err
		}
	}

	return nil
}

func extractTarEntry(header *tar.Header, reader io.Reader, targetPath string) error {
	switch header.Typeflag {
	case tar.TypeDir:
		if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
			return fmt.Errorf("creating directory %s: %w", targetPath, err)
		}

	case tar.TypeReg:
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("creating parent directory for %s: %w", targetPath, err)
		}

		if err := writeFile(targetPath, os.FileMode(header.Mode), reader); err != nil {
			return err
		}

	case tar.TypeSymlink:
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("creating parent directory for symlink %s: %w", targetPath, err)
		}

		if err := os.Symlink(header.Linkname, targetPath); err != nil {
			return fmt.Errorf("creating symlink %s -> %s: %w", targetPath, header.Linkname, err)
		}
	}

	return nil
}

// extractZip extracts a zip archive to the destination directory.
// It strips the top-level directory from the archive paths.
func extractZip(data []byte, destDir string) error {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("creating zip reader: %w", err)
	}

	for _, file := range zipReader.File {
		relativePath := stripTopLevelDir(file.Name)
		if relativePath == "" {
			continue
		}

		targetPath := filepath.Join(destDir, relativePath)

		if err := extractZipEntry(file, targetPath); err != nil {
			return err
		}
	}

	return nil
}

func extractZipEntry(file *zip.File, targetPath string) error {
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(targetPath, file.Mode()); err != nil {
			return fmt.Errorf("creating directory %s: %w", targetPath, err)
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("creating parent directory for %s: %w", targetPath, err)
	}

	srcFile, err := file.Open()
	if err != nil {
		return fmt.Errorf("opening zip entry %s: %w", file.Name, err)
	}
	defer func() { _ = srcFile.Close() }()

	return writeFile(targetPath, file.Mode(), srcFile)
}

func writeFile(targetPath string, mode os.FileMode, reader io.Reader) error {
	outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", targetPath, err)
	}

	_, copyErr := io.Copy(outFile, reader)
	closeErr := outFile.Close()

	if copyErr != nil {
		return fmt.Errorf("writing file %s: %w", targetPath, copyErr)
	}
	if closeErr != nil {
		return fmt.Errorf("closing file %s: %w", targetPath, closeErr)
	}

	return nil
}

// stripTopLevelDir removes the first directory component from a path.
// For example, "onnxruntime-linux-x64-1.23.2/lib/file.so" becomes "lib/file.so".
func stripTopLevelDir(path string) string {
	cleaned := filepath.Clean(path)
	idx := 0
	for i, c := range cleaned {
		if c == '/' || c == filepath.Separator {
			idx = i + 1
			break
		}
	}
	if idx >= len(cleaned) {
		return ""
	}
	return cleaned[idx:]
}
