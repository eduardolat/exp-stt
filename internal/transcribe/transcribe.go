package transcribe

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/go-audio/wav"
	"github.com/varavelio/tribar/internal/onnx"
	ort "github.com/yalue/onnxruntime_go"
)

// Instance represents a transcription engine instance.
type Instance struct {
	parakeet          *ParakeetModel
	integrityVerified bool
}

// New creates a new transcription instance.
func New() (*Instance, error) {
	ort.SetSharedLibraryPath(onnx.SharedLibraryPath)

	if err := ort.InitializeEnvironment(); err != nil {
		return nil, fmt.Errorf("error initializing onnx runtime: %w", err)
	}

	parakeet, err := NewParakeetModel()
	if err != nil {
		return nil, fmt.Errorf("error creating parakeet model: %w", err)
	}

	return &Instance{
		parakeet: parakeet,
	}, nil
}

// Shutdown cleans up resources used by the transcription instance.
func (i *Instance) Shutdown() error {
	if err := ort.DestroyEnvironment(); err != nil {
		return fmt.Errorf("error destroying onnx runtime environment: %w", err)
	}
	return nil
}

// CheckModels checks if all required models exist with full SHA256 verification.
// After successful verification, marks integrity as verified to skip checksums on reload.
func (i *Instance) CheckModels() (bool, []ModelFile) {
	allExist, missing := i.parakeet.CheckModelsExist()
	if allExist {
		i.integrityVerified = true
	}
	return allExist, missing
}

// CheckModelsQuick checks if all required model files exist (no checksum verification).
// Used after integrity has already been verified to avoid 1-2 second lag on reload.
func (i *Instance) CheckModelsQuick() bool {
	return i.parakeet.CheckModelsExistQuick()
}

// DownloadModels downloads all missing model files.
// After successful download, marks integrity as verified.
func (i *Instance) DownloadModels(progressCallback DownloadProgressCallback) error {
	if err := i.parakeet.DownloadModels(progressCallback); err != nil {
		return err
	}
	i.integrityVerified = true
	return nil
}

// LoadModels loads the vocabulary and prepares the model for transcription.
func (i *Instance) LoadModels() error {
	// Use quick check if integrity was already verified, otherwise do full check
	if i.integrityVerified {
		if !i.CheckModelsQuick() {
			// Files disappeared, need to re-download
			i.integrityVerified = false
			return errors.New("missing model files, please restart the application")
		}
	} else {
		if allExist, _ := i.CheckModels(); !allExist {
			return errors.New("missing model files, please restart the application")
		}
	}

	// Load vocabulary
	if err := i.parakeet.LoadVocabulary(); err != nil {
		return fmt.Errorf("error loading vocabulary: %w", err)
	}

	return nil
}

// UnloadModels clears the loaded vocabulary to free memory.
// The models can be reloaded later by calling LoadModels again.
func (i *Instance) UnloadModels() {
	i.parakeet.UnloadVocabulary()
}

// TranscribeWAV transcribes audio from WAV bytes.
// The WAV must be 16kHz mono PCM format (as produced by the recorder).
func (i *Instance) TranscribeWAV(wavData []byte) (string, error) {
	samples, err := processWAVBytes(wavData)
	if err != nil {
		return "", fmt.Errorf("error processing WAV data: %w", err)
	}

	return i.parakeet.Transcribe(samples)
}

// TranscribeSamples transcribes audio from float32 samples.
// Samples must already be 16kHz mono audio normalized to [-1, 1].
func (i *Instance) TranscribeSamples(samples []float32) (string, error) {
	return i.parakeet.Transcribe(samples)
}

// processWAVBytes reads WAV bytes and converts to float32 samples.
// Expects 16kHz mono PCM format (as produced by the recorder).
func processWAVBytes(wavData []byte) ([]float32, error) {
	reader := bytes.NewReader(wavData)
	decoder := wav.NewDecoder(reader)

	if !decoder.IsValidFile() {
		return nil, errors.New("invalid WAV file")
	}

	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, fmt.Errorf("error decoding WAV: %w", err)
	}

	// Convert to float32 normalized to [-1, 1]
	samples := make([]float32, len(buf.Data))
	for j, val := range buf.Data {
		samples[j] = float32(val) / 32768.0
	}

	return samples, nil
}

// ReadWAVFile is a helper function to read a WAV file into bytes.
func ReadWAVFile(filepath string) ([]byte, error) {
	return os.ReadFile(filepath)
}
