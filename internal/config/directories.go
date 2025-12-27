package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/eduardolat/exp-stt/internal/logger"
)

const dirAppName = "stt"

var (
	DirectoryConfig         = ""
	DirectoryData           = ""
	DirectoryOnnxRuntime    = ""
	DirectoryModels         = ""
	DirectoryModelsParakeet = ""
)

// EnsureDirectories creates all necessary directories if they don't exist.
func EnsureDirectories(logger logger.Logger) error {
	configDir, err := calculateConfigDir()
	if err != nil {
		return fmt.Errorf("could not determine config directory: %w", err)
	}

	dataDir, err := calculateDataDir()
	if err != nil {
		return fmt.Errorf("could not determine data directory: %w", err)
	}

	DirectoryConfig = configDir
	DirectoryData = dataDir
	DirectoryOnnxRuntime = filepath.Join(DirectoryData, "onnxruntime")
	DirectoryModels = filepath.Join(DirectoryData, "models")
	DirectoryModelsParakeet = filepath.Join(DirectoryModels, "parakeet")

	// We only have to create the deepest directories, as os.MkdirAll will create all necessary parents.
	ensureDirs := []string{DirectoryConfig, DirectoryOnnxRuntime, DirectoryModelsParakeet}
	for _, dir := range ensureDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(fmt.Errorf("failed to create application directory %s: %w", dir, err))
		}
	}

	logger.Debug(
		context.Background(), "application directories ensured",
		"directory_config", DirectoryConfig,
		"directory_data", DirectoryData,
		"directory_onnx_runtime", DirectoryOnnxRuntime,
		"directory_models", DirectoryModels,
		"directory_parakeet_models", DirectoryModelsParakeet,
	)

	return nil
}

// calculateConfigDir returns the base config directory for the application.
//
// This follows OS-specific conventions:
//   - Windows: %APPDATA%\{app_name}
//   - macOS: ~/Library/Application Support/{app_name}
//   - Linux: ~/.config/{app_name}
func calculateConfigDir() (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		baseDir = os.Getenv("LOCALAPPDATA")
		if baseDir == "" {
			baseDir = os.Getenv("APPDATA") // Fallback
		}
		if baseDir == "" {
			return "", fmt.Errorf("the LOCALAPPDATA or APPDATA environment variable is not set")
		}
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine user home directory: %w", err)
		}
		baseDir = filepath.Join(homeDir, "Library", "Application Support")
	default: // Linux and other Unix-like systems
		baseDir = os.Getenv("XDG_CONFIG_HOME")
		if baseDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("cannot determine user home directory: %w", err)
			}
			baseDir = filepath.Join(homeDir, ".config")
		}
	}

	return filepath.Join(baseDir, dirAppName), nil
}

// calculateDataDir returns the base data directory for the application.
//
// This follows OS-specific conventions:
//   - Windows: %LOCALAPPDATA%\{app_name}
//   - macOS: ~/Library/Application Support/{app_name}
//   - Linux: ~/.local/share/{app_name}
func calculateDataDir() (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		baseDir = os.Getenv("LOCALAPPDATA")
		if baseDir == "" {
			baseDir = os.Getenv("APPDATA") // Fallback
		}
		if baseDir == "" {
			return "", fmt.Errorf("the LOCALAPPDATA or APPDATA environment variable is not set")
		}
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine user home directory: %w", err)
		}
		baseDir = filepath.Join(homeDir, "Library", "Application Support")
	default: // Linux and other Unix-like systems
		baseDir = os.Getenv("XDG_DATA_HOME")
		if baseDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("cannot determine user home directory: %w", err)
			}
			baseDir = filepath.Join(homeDir, ".local", "share")
		}
	}

	return filepath.Join(baseDir, dirAppName), nil
}
