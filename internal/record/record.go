package record

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/gen2brain/malgo"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/logger"
)

var (
	ErrAlreadyRecording = fmt.Errorf("recording is already in progress")
)

type Recorder struct {
	logger          logger.Logger
	device          *malgo.Device
	ctx             *malgo.AllocatedContext
	settingsManager *config.SettingsManager
	isRecording     bool
	data            []byte
	mu              sync.Mutex
}

func NewRecorder(logger logger.Logger, settingsManager *config.SettingsManager) (*Recorder, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, err
	}
	return &Recorder{
		logger:          logger,
		ctx:             ctx,
		settingsManager: settingsManager,
	}, nil
}

// Shutdown cleans up resources used by the recorder.
func (r *Recorder) Shutdown() {
	r.Stop()
	if r.ctx != nil {
		if err := r.ctx.Uninit(); err != nil {
			r.logger.Error(context.Background(), "failed to uninit context", "err", err)
		}
	}
}

// Start begins the recording process. It cleans the buffer and starts capturing audio data.
func (r *Recorder) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isRecording {
		return ErrAlreadyRecording
	}

	r.data = []byte{} // Clean the buffer before starting
	r.isRecording = true

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = 16000

	// Set specific device if configured
	settings := r.settingsManager.Get()
	if settings.InputDevice != "" && settings.InputDevice != "default" {
		deviceID := parseDeviceID(settings.InputDevice)
		if deviceID != nil {
			deviceConfig.Capture.DeviceID = deviceID.Pointer()
		}
	}

	onData := func(pOutput, pInput []byte, frameCount uint32) {
		r.mu.Lock()
		if r.isRecording {
			r.data = append(r.data, pInput...)
		}
		r.mu.Unlock()
	}

	var err error
	r.device, err = malgo.InitDevice(r.ctx.Context, deviceConfig, malgo.DeviceCallbacks{Data: onData})
	if err != nil {
		return err
	}

	return r.device.Start()
}

// Stop stops the recording process.
func (r *Recorder) Stop() {
	r.mu.Lock()
	wasRecording := r.isRecording
	r.isRecording = false
	r.mu.Unlock()

	// Only stop and uninit device if we were actually recording
	if r.device != nil && wasRecording {
		if err := r.device.Stop(); err != nil {
			r.logger.Error(context.Background(), "failed to stop device", "err", err)
		}
		r.device.Uninit()
		r.device = nil
	}
}

// GetData returns a copy of the raw PCM audio data.
func (r *Recorder) GetData() []byte {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]byte, len(r.data))
	copy(result, r.data)
	return result
}

// BuildWAV returns the recorded audio as a complete WAV file in memory.
func (r *Recorder) BuildWAV() []byte {
	r.mu.Lock()
	defer r.mu.Unlock()

	return buildWavBytes(r.data, 16000, 1)
}

// parseDeviceID parses a hex string device ID back to malgo.DeviceID.
func parseDeviceID(hexStr string) *malgo.DeviceID {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil
	}

	var id malgo.DeviceID
	if len(data) > len(id) {
		data = data[:len(id)]
	}
	copy(id[:], data)
	return &id
}

// buildWavBytes creates a WAV file as bytes from raw PCM data.
func buildWavBytes(data []byte, sampleRate, channels int) []byte {
	dataSize := len(data)
	headerSize := 44
	totalSize := headerSize + dataSize

	buf := make([]byte, totalSize)

	// RIFF header
	copy(buf[0:4], "RIFF")
	binary.LittleEndian.PutUint32(buf[4:8], uint32(36+dataSize))
	copy(buf[8:12], "WAVE")

	// fmt subchunk
	copy(buf[12:16], "fmt ")
	binary.LittleEndian.PutUint32(buf[16:20], 16) // Subchunk1Size
	binary.LittleEndian.PutUint16(buf[20:22], 1)  // AudioFormat (PCM)
	binary.LittleEndian.PutUint16(buf[22:24], uint16(channels))
	binary.LittleEndian.PutUint32(buf[24:28], uint32(sampleRate))
	binary.LittleEndian.PutUint32(buf[28:32], uint32(sampleRate*channels*2)) // ByteRate
	binary.LittleEndian.PutUint16(buf[32:34], uint16(channels*2))            // BlockAlign
	binary.LittleEndian.PutUint16(buf[34:36], 16)                            // BitsPerSample

	// data subchunk
	copy(buf[36:40], "data")
	binary.LittleEndian.PutUint32(buf[40:44], uint32(dataSize))

	// Audio data
	copy(buf[44:], data)

	return buf
}
