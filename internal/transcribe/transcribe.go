package transcribe

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/eduardolat/exp-stt/internal/onnx"
	"github.com/go-audio/wav"
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

// TranscribeWAV transcribes audio from WAV bytes.
// The WAV can be in any format (sample rate, channels, bit depth) - it will be
// automatically converted to the required format (16kHz, mono, float32).
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

// processWAVBytes reads WAV bytes and converts to 16kHz mono float32 samples.
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

	// Convert to float32 normalized
	rawSamples := make([]float32, len(buf.Data))
	for j, val := range buf.Data {
		rawSamples[j] = float32(val) / 32768.0
	}

	// Convert to mono if stereo
	numChannels := buf.Format.NumChannels
	var monoSamples []float32
	if numChannels > 1 {
		monoSamples = convertToMono(rawSamples, numChannels)
	} else {
		monoSamples = rawSamples
	}

	// Resample to 16kHz if needed
	originalSampleRate := buf.Format.SampleRate
	targetSampleRate := 16000

	var samples []float32
	if originalSampleRate != targetSampleRate {
		samples = resample(monoSamples, originalSampleRate, targetSampleRate)
	} else {
		samples = monoSamples
	}

	return samples, nil
}

// convertToMono converts multi-channel audio to mono by averaging channels.
func convertToMono(samples []float32, numChannels int) []float32 {
	numSamples := len(samples) / numChannels
	mono := make([]float32, numSamples)

	for i := range numSamples {
		var sum float32
		for ch := range numChannels {
			sum += samples[i*numChannels+ch]
		}
		mono[i] = sum / float32(numChannels)
	}

	return mono
}

// resample performs linear interpolation resampling.
func resample(input []float32, fromRate, toRate int) []float32 {
	if fromRate == toRate {
		return input
	}

	ratio := float64(fromRate) / float64(toRate)
	targetLength := int(float64(len(input)) / ratio)
	output := make([]float32, targetLength)

	for i := range targetLength {
		pos := float64(i) * ratio
		index := int(pos)
		frac := float32(pos - float64(index))

		low := index
		high := index + 1
		if high >= len(input) {
			high = len(input) - 1
		}

		output[i] = (1-frac)*input[low] + frac*input[high]
	}

	return output
}

// ReadWAVFile is a helper function to read a WAV file into bytes.
func ReadWAVFile(filepath string) ([]byte, error) {
	return os.ReadFile(filepath)
}
