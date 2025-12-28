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
	parakeet *ParakeetModel
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

// CheckModels checks if all required models exist.
// Returns true if all models exist, false otherwise with the list of missing models.
func (i *Instance) CheckModels() (bool, []ModelFile) {
	return i.parakeet.CheckModelsExist()
}

// DownloadModels downloads all missing model files.
func (i *Instance) DownloadModels(progressCallback DownloadProgressCallback) error {
	return i.parakeet.DownloadModels(progressCallback)
}

// LoadModels loads the vocabulary and prepares the model for transcription.
func (i *Instance) LoadModels() error {
	// Check if models exist
	allExist, missing := i.CheckModels()
	if !allExist {
		var missingNames []string
		for _, m := range missing {
			missingNames = append(missingNames, m.Name)
		}
		return fmt.Errorf("missing model files: %v. Call DownloadModels first", missingNames)
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
