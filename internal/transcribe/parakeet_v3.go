package transcribe

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/varavelio/tribar/internal/config"
	ort "github.com/yalue/onnxruntime_go"
)

// parakeetFile defines a model file with all required download information.
type parakeetFile struct {
	Name        string // Display name for progress
	FileName    string // Local file name
	CDNURL      string // Primary download URL (CDN)
	FallbackURL string // Fallback download URL (HuggingFace)
	SHA256      string // Expected SHA256 checksum
}

// Parakeet V3 model files configuration
var parakeetFiles = []parakeetFile{
	{
		Name:        "Vocabulary",
		FileName:    "vocab.txt",
		CDNURL:      "https://cdn.varavel.com/tribar/models/parakeet-v3/abd2878/vocab.txt",
		FallbackURL: "https://huggingface.co/istupakov/parakeet-tdt-0.6b-v3-onnx/resolve/abd2878d52a678ce380088ef9d9b1d9664404565/vocab.txt?download=true",
		SHA256:      "d58544679ea4bc6ac563d1f545eb7d474bd6cfa467f0a6e2c1dc1c7d37e3c35d",
	},
	{
		Name:        "Preprocessor",
		FileName:    "nemo128.onnx",
		CDNURL:      "https://cdn.varavel.com/tribar/models/parakeet-v3/abd2878/nemo128.onnx",
		FallbackURL: "https://huggingface.co/istupakov/parakeet-tdt-0.6b-v3-onnx/resolve/abd2878d52a678ce380088ef9d9b1d9664404565/nemo128.onnx?download=true",
		SHA256:      "a9fde1486ebfcc08f328d75ad4610c67835fea58c73ba57e3209a6f6cf019e9f",
	},
	{
		Name:        "Encoder",
		FileName:    "encoder-model.int8.onnx",
		CDNURL:      "https://cdn.varavel.com/tribar/models/parakeet-v3/abd2878/encoder-model.int8.onnx",
		FallbackURL: "https://huggingface.co/istupakov/parakeet-tdt-0.6b-v3-onnx/resolve/abd2878d52a678ce380088ef9d9b1d9664404565/encoder-model.int8.onnx?download=true",
		SHA256:      "6139d2fa7e1b086097b277c7149725edbab89cc7c7ae64b23c741be4055aff09",
	},
	{
		Name:        "Decoder",
		FileName:    "decoder_joint-model.int8.onnx",
		CDNURL:      "https://cdn.varavel.com/tribar/models/parakeet-v3/abd2878/decoder_joint-model.int8.onnx",
		FallbackURL: "https://huggingface.co/istupakov/parakeet-tdt-0.6b-v3-onnx/resolve/abd2878d52a678ce380088ef9d9b1d9664404565/decoder_joint-model.int8.onnx?download=true",
		SHA256:      "eea7483ee3d1a30375daedc8ed83e3960c91b098812127a0d99d1c8977667a70",
	},
}

// Parakeet model constants
const (
	parakeetSubsamplingFactor = 8
	parakeetDecoderHiddenSize = 640
	parakeetEncoderHiddenSize = 1024
	parakeetNumMelBins        = 128
	parakeetHopLength         = 160 // 10ms @ 16kHz
	parakeetNumDurations      = 5   // TDT duration options
)

// ParakeetModel represents the Parakeet TDT model for speech recognition.
type ParakeetModel struct {
	vocab    []string
	blankIdx int32

	vocabPath        string
	preprocessorPath string
	encoderPath      string
	decoderPath      string

	// ONNX sessions (kept in memory for faster inference)
	preprocessorSession *ort.DynamicAdvancedSession
	encoderSession      *ort.DynamicAdvancedSession
	decoderSession      *ort.DynamicAdvancedSession
}

// NewParakeetModel creates a new ParakeetModel instance.
func NewParakeetModel() (*ParakeetModel, error) {
	parakeetDir := config.DirectoryModelsParakeetV3

	return &ParakeetModel{
		vocabPath:        path.Join(parakeetDir, parakeetFiles[0].FileName),
		preprocessorPath: path.Join(parakeetDir, parakeetFiles[1].FileName),
		encoderPath:      path.Join(parakeetDir, parakeetFiles[2].FileName),
		decoderPath:      path.Join(parakeetDir, parakeetFiles[3].FileName),
	}, nil
}

// ModelFile represents a model file with its URLs, path, and checksum.
type ModelFile struct {
	Name        string
	CDNURL      string
	FallbackURL string
	Path        string
	SHA256      string
}

// GetModelFiles returns all model files with their URLs, paths, and checksums.
func (p *ParakeetModel) GetModelFiles() []ModelFile {
	paths := []string{p.vocabPath, p.preprocessorPath, p.encoderPath, p.decoderPath}
	files := make([]ModelFile, len(parakeetFiles))

	for i, pf := range parakeetFiles {
		files[i] = ModelFile{
			Name:        pf.Name,
			CDNURL:      pf.CDNURL,
			FallbackURL: pf.FallbackURL,
			Path:        paths[i],
			SHA256:      pf.SHA256,
		}
	}

	return files
}

// CheckModelsExist checks if all required model files exist with valid checksums.
func (p *ParakeetModel) CheckModelsExist() (bool, []ModelFile) {
	var missing []ModelFile

	for _, file := range p.GetModelFiles() {
		if !isFileValidWithChecksum(file.Path, file.SHA256) {
			missing = append(missing, file)
		}
	}

	return len(missing) == 0, missing
}

// CheckModelsExistQuick checks if all required model files exist (no checksum verification).
// This is faster than CheckModelsExist and should be used when integrity was already verified.
func (p *ParakeetModel) CheckModelsExistQuick() bool {
	for _, file := range p.GetModelFiles() {
		if _, err := os.Stat(file.Path); os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// isFileValidWithChecksum checks if a file exists and has a valid SHA256 checksum.
func isFileValidWithChecksum(filepath, expectedSHA256 string) bool {
	file, err := os.Open(filepath)
	if err != nil {
		return false
	}
	defer func() { _ = file.Close() }()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false
	}

	actualSHA256 := hex.EncodeToString(hash.Sum(nil))
	return actualSHA256 == expectedSHA256
}

// DownloadProgressCallback is called during download with progress information.
type DownloadProgressCallback func(filename string, downloaded, total int64, percent float64)

// DownloadModels downloads all missing or invalid model files.
func (p *ParakeetModel) DownloadModels(progressCallback DownloadProgressCallback) error {
	_, missing := p.CheckModelsExist()
	if len(missing) == 0 {
		return nil
	}

	for _, file := range missing {
		// Delete any existing file (could be corrupt/partial)
		_ = os.Remove(file.Path)

		// Try CDN first
		err := downloadFileWithProgress(file.Path, file.CDNURL, file.Name, progressCallback)
		if err == nil && isFileValidWithChecksum(file.Path, file.SHA256) {
			continue
		}

		// CDN failed or checksum invalid - try fallback
		_ = os.Remove(file.Path)
		err = downloadFileWithProgress(file.Path, file.FallbackURL, file.Name, progressCallback)
		if err != nil {
			_ = os.Remove(file.Path)
			return fmt.Errorf("failed to download %s: %w", file.Name, err)
		}

		// Verify fallback download checksum
		if !isFileValidWithChecksum(file.Path, file.SHA256) {
			_ = os.Remove(file.Path)
			return fmt.Errorf("checksum verification failed for %s", file.Name)
		}
	}

	return nil
}

// downloadFileWithProgress downloads a file from URL to the specified path with progress tracking.
func downloadFileWithProgress(filepath, url, name string, progressCallback DownloadProgressCallback) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", config.UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		_ = os.Remove(filepath)
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		_ = os.Remove(filepath)
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	contentLength := resp.ContentLength
	var written int64
	buf := make([]byte, 32*1024)

	for {
		nr, readErr := resp.Body.Read(buf)
		if nr > 0 {
			nw, writeErr := out.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if writeErr != nil {
				_ = os.Remove(filepath)
				return writeErr
			}
			if nr != nw {
				_ = os.Remove(filepath)
				return io.ErrShortWrite
			}

			if progressCallback != nil && contentLength > 0 {
				percent := float64(written) / float64(contentLength) * 100
				progressCallback(name, written, contentLength, percent)
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			_ = os.Remove(filepath)
			return readErr
		}
	}

	return nil
}

// LoadModel loads vocabulary and ONNX sessions into memory for transcription.
func (p *ParakeetModel) LoadModel() error {
	// Load vocabulary
	file, err := os.Open(p.vocabPath)
	if err != nil {
		return fmt.Errorf("error opening vocab file: %w", err)
	}
	defer func() { _ = file.Close() }()

	var vocab []string
	var blankIdx int32 = -1
	scanner := bufio.NewScanner(file)
	for idx := 0; scanner.Scan(); idx++ {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) > 0 {
			token := parts[0]
			vocab = append(vocab, token)
			if token == "<blk>" {
				blankIdx = int32(idx)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading vocab file: %w", err)
	}

	// If no <blk> token found, assume last token is blank
	if blankIdx == -1 {
		blankIdx = int32(len(vocab) - 1)
	}

	p.vocab = vocab
	p.blankIdx = blankIdx

	// Create ONNX sessions
	p.preprocessorSession, err = ort.NewDynamicAdvancedSession(
		p.preprocessorPath,
		[]string{"waveforms", "waveforms_lens"},
		[]string{"features", "features_lens"},
		nil,
	)
	if err != nil {
		p.UnloadModel()
		return fmt.Errorf("error creating preprocessor session: %w", err)
	}

	p.encoderSession, err = ort.NewDynamicAdvancedSession(
		p.encoderPath,
		[]string{"audio_signal", "length"},
		[]string{"outputs", "encoded_lengths"},
		nil,
	)
	if err != nil {
		p.UnloadModel()
		return fmt.Errorf("error creating encoder session: %w", err)
	}

	p.decoderSession, err = ort.NewDynamicAdvancedSession(
		p.decoderPath,
		[]string{"encoder_outputs", "targets", "target_length", "input_states_1", "input_states_2"},
		[]string{"outputs", "output_states_1", "output_states_2"},
		nil,
	)
	if err != nil {
		p.UnloadModel()
		return fmt.Errorf("error creating decoder session: %w", err)
	}

	return nil
}

// UnloadModel clears vocabulary and destroys ONNX sessions to free memory.
func (p *ParakeetModel) UnloadModel() {
	// Destroy ONNX sessions
	if p.preprocessorSession != nil {
		_ = p.preprocessorSession.Destroy()
		p.preprocessorSession = nil
	}
	if p.encoderSession != nil {
		_ = p.encoderSession.Destroy()
		p.encoderSession = nil
	}
	if p.decoderSession != nil {
		_ = p.decoderSession.Destroy()
		p.decoderSession = nil
	}

	// Clear vocabulary
	p.vocab = nil
	p.blankIdx = 0
}

// Transcribe performs speech-to-text on audio samples.
// samples should be 16kHz mono float32 audio normalized to [-1, 1].
func (p *ParakeetModel) Transcribe(samples []float32) (string, error) {
	if len(p.vocab) == 0 {
		return "", fmt.Errorf("model not loaded, call LoadModel first")
	}

	features, featuresLen, err := p.runPreprocessor(samples)
	if err != nil {
		return "", fmt.Errorf("preprocessor error: %w", err)
	}

	encoderOut, encoderLen, err := p.runEncoder(features, featuresLen)
	if err != nil {
		return "", fmt.Errorf("encoder error: %w", err)
	}

	text, err := p.runDecoder(encoderOut, encoderLen)
	if err != nil {
		return "", fmt.Errorf("decoder error: %w", err)
	}

	return text, nil
}

func (p *ParakeetModel) runPreprocessor(samples []float32) ([]float32, int64, error) {
	samplesLen := int64(len(samples))

	// Input tensors
	waveformsTensor, err := ort.NewTensor(ort.NewShape(1, samplesLen), samples)
	if err != nil {
		return nil, 0, fmt.Errorf("error creating waveforms tensor: %w", err)
	}
	defer func() { _ = waveformsTensor.Destroy() }()

	waveformsLensTensor, err := ort.NewTensor(ort.NewShape(1), []int64{samplesLen})
	if err != nil {
		return nil, 0, fmt.Errorf("error creating waveforms_lens tensor: %w", err)
	}
	defer func() { _ = waveformsLensTensor.Destroy() }()

	// Output tensors - calculate expected size
	expectedTimeSteps := (samplesLen / parakeetHopLength) + 1

	featShape := ort.NewShape(1, parakeetNumMelBins, expectedTimeSteps)
	featTensor, err := ort.NewEmptyTensor[float32](featShape)
	if err != nil {
		return nil, 0, fmt.Errorf("error creating features tensor: %w", err)
	}
	defer func() { _ = featTensor.Destroy() }()

	featLensTensor, err := ort.NewEmptyTensor[int64](ort.NewShape(1))
	if err != nil {
		return nil, 0, fmt.Errorf("error creating features_lens tensor: %w", err)
	}
	defer func() { _ = featLensTensor.Destroy() }()

	// Run using persistent session
	if err := p.preprocessorSession.Run(
		[]ort.ArbitraryTensor{waveformsTensor, waveformsLensTensor},
		[]ort.ArbitraryTensor{featTensor, featLensTensor},
	); err != nil {
		return nil, 0, fmt.Errorf("error running preprocessor: %w", err)
	}

	features := make([]float32, len(featTensor.GetData()))
	copy(features, featTensor.GetData())
	featLen := featLensTensor.GetData()[0]

	return features, featLen, nil
}

func (p *ParakeetModel) runEncoder(features []float32, featuresLen int64) ([]float32, int64, error) {
	timeSteps := int64(len(features)) / parakeetNumMelBins

	// Input tensors
	audioSignalTensor, err := ort.NewTensor(ort.NewShape(1, parakeetNumMelBins, timeSteps), features)
	if err != nil {
		return nil, 0, fmt.Errorf("error creating audio_signal tensor: %w", err)
	}
	defer func() { _ = audioSignalTensor.Destroy() }()

	lengthTensor, err := ort.NewTensor(ort.NewShape(1), []int64{featuresLen})
	if err != nil {
		return nil, 0, fmt.Errorf("error creating length tensor: %w", err)
	}
	defer func() { _ = lengthTensor.Destroy() }()

	// Output tensors
	encoderTimeSteps := (featuresLen + parakeetSubsamplingFactor - 1) / parakeetSubsamplingFactor

	encOutShape := ort.NewShape(1, parakeetEncoderHiddenSize, encoderTimeSteps)
	encOutTensor, err := ort.NewEmptyTensor[float32](encOutShape)
	if err != nil {
		return nil, 0, fmt.Errorf("error creating encoder output tensor: %w", err)
	}
	defer func() { _ = encOutTensor.Destroy() }()

	encLensTensor, err := ort.NewEmptyTensor[int64](ort.NewShape(1))
	if err != nil {
		return nil, 0, fmt.Errorf("error creating encoder lengths tensor: %w", err)
	}
	defer func() { _ = encLensTensor.Destroy() }()

	// Run using persistent session
	if err := p.encoderSession.Run(
		[]ort.ArbitraryTensor{audioSignalTensor, lengthTensor},
		[]ort.ArbitraryTensor{encOutTensor, encLensTensor},
	); err != nil {
		return nil, 0, fmt.Errorf("error running encoder: %w", err)
	}

	encoderOut := make([]float32, len(encOutTensor.GetData()))
	copy(encoderOut, encOutTensor.GetData())
	encoderLen := encLensTensor.GetData()[0]

	return encoderOut, encoderLen, nil
}

func (p *ParakeetModel) runDecoder(encoderOut []float32, encoderLen int64) (string, error) {
	var transcribedTokens []string
	var lastEmittedToken int32 = -1 // Track last emitted for deduplication

	vocabSize := len(p.vocab)
	lastToken := p.blankIdx

	// --- Pre-allocate tensors and buffers ---

	// Constant shapes
	encOutShape := ort.NewShape(1, parakeetEncoderHiddenSize, 1)
	targetsShape := ort.NewShape(1, 1)
	targetLenShape := ort.NewShape(1)
	stateShape := ort.NewShape(2, 1, parakeetDecoderHiddenSize)
	outputSize := int64(len(p.vocab) + parakeetNumDurations)
	logitsShape := ort.NewShape(1, 1, 1, outputSize)
	outStateShape := ort.NewShape(2, 1, parakeetDecoderHiddenSize)

	// Create Tensors (using NewEmptyTensor for ones we fill, NewTensor for constants)
	// Input: encoder_outputs (1, 1024, 1)
	encOutTensor, err := ort.NewEmptyTensor[float32](encOutShape)
	if err != nil {
		return "", fmt.Errorf("error creating encoder_outputs tensor: %w", err)
	}
	defer func() { _ = encOutTensor.Destroy() }()

	// Input: targets (1, 1) - Initialized with blankIdx
	targetsTensor, err := ort.NewTensor(targetsShape, []int32{lastToken})
	if err != nil {
		return "", fmt.Errorf("error creating targets tensor: %w", err)
	}
	defer func() { _ = targetsTensor.Destroy() }()

	// Input: target_length (1) - Always 1
	targetLenTensor, err := ort.NewTensor(targetLenShape, []int32{1})
	if err != nil {
		return "", fmt.Errorf("error creating target_length tensor: %w", err)
	}
	defer func() { _ = targetLenTensor.Destroy() }()

	// Input: input_states_1 (2, 1, 640) - Initialized to 0
	state1Tensor, err := ort.NewEmptyTensor[float32](stateShape)
	if err != nil {
		return "", fmt.Errorf("error creating input_states_1 tensor: %w", err)
	}
	defer func() { _ = state1Tensor.Destroy() }()
	clear(state1Tensor.GetData()) // Zero initialization

	// Input: input_states_2 (2, 1, 640) - Initialized to 0
	state2Tensor, err := ort.NewEmptyTensor[float32](stateShape)
	if err != nil {
		return "", fmt.Errorf("error creating input_states_2 tensor: %w", err)
	}
	defer func() { _ = state2Tensor.Destroy() }()
	clear(state2Tensor.GetData()) // Zero initialization

	// Output: outputs (1, 1, 1, outputSize)
	logitsTensor, err := ort.NewEmptyTensor[float32](logitsShape)
	if err != nil {
		return "", fmt.Errorf("error creating outputs tensor: %w", err)
	}
	defer func() { _ = logitsTensor.Destroy() }()

	// Output: output_states_1 (2, 1, 640)
	outState1Tensor, err := ort.NewEmptyTensor[float32](outStateShape)
	if err != nil {
		return "", fmt.Errorf("error creating output_states_1 tensor: %w", err)
	}
	defer func() { _ = outState1Tensor.Destroy() }()

	// Output: output_states_2 (2, 1, 640)
	outState2Tensor, err := ort.NewEmptyTensor[float32](outStateShape)
	if err != nil {
		return "", fmt.Errorf("error creating output_states_2 tensor: %w", err)
	}
	defer func() { _ = outState2Tensor.Destroy() }()

	// --- Loop ---
	inputs := []ort.ArbitraryTensor{encOutTensor, targetsTensor, targetLenTensor, state1Tensor, state2Tensor}
	outputs := []ort.ArbitraryTensor{logitsTensor, outState1Tensor, outState2Tensor}

	encOutData := encOutTensor.GetData()
	targetsData := targetsTensor.GetData()
	state1Data := state1Tensor.GetData()
	state2Data := state2Tensor.GetData()

	for t := range encoderLen {
		// 1. Update encoder_outputs tensor data (copy from global encoder output)
		// Instead of creating a slice and copying, update directly
		// Optimization: encOutData is a slice mapped to C memory
		for k := range parakeetEncoderHiddenSize {
			idx := int64(k)*encoderLen + t
			if idx < int64(len(encoderOut)) {
				encOutData[k] = encoderOut[idx]
			} else {
				encOutData[k] = 0 // Padding/Safety
			}
		}

		// 2. Targets tensor data is already set (from init or previous loop)

		// 3. Run Session
		if err := p.decoderSession.Run(inputs, outputs); err != nil {
			return "", fmt.Errorf("error running decoder at t=%d: %w", t, err)
		}

		// 4. Get logits and find best token
		logits := logitsTensor.GetData()
		vocabLogits := logits[:vocabSize] // Slice without duration logits
		bestToken := argmax(vocabLogits)

		if bestToken != p.blankIdx && bestToken != lastEmittedToken {
			// Emit non-blank token
			transcribedTokens = append(transcribedTokens, p.vocab[bestToken])

			// Update state for next step
			lastToken = bestToken
			lastEmittedToken = bestToken

			// Update tensors for next iteration:
			// Update targetsTensor
			targetsData[0] = lastToken

			// Update state tensors (copy output states to input states)
			// We can copy directly between tensor data slices
			copy(state1Data, outState1Tensor.GetData())
			copy(state2Data, outState2Tensor.GetData())

		} else if bestToken == p.blankIdx {
			// Reset deduplication on blank
			lastEmittedToken = -1
			// State remains the same (from previous accepted token), so we don't copy outState -> state
		}
	}

	// Post-process result
	result := strings.Join(transcribedTokens, "")
	result = strings.ReplaceAll(result, "▁", " ")
	result = strings.ReplaceAll(result, "\u2581", " ")
	return strings.TrimSpace(result), nil
}

func argmax(slice []float32) int32 {
	if len(slice) == 0 {
		return 0
	}
	var maxIdx int32
	maxVal := slice[0]
	for i, val := range slice {
		if val > maxVal {
			maxVal = val
			maxIdx = int32(i)
		}
	}
	return maxIdx
}
