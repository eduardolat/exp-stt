package record

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"

	"github.com/gen2brain/malgo"
)

var (
	ErrAlreadyRecording = fmt.Errorf("recording is already in progress")
)

type Recorder struct {
	device      *malgo.Device
	ctx         *malgo.AllocatedContext
	isRecording bool
	data        []byte
	mu          sync.Mutex
}

func NewRecorder() (*Recorder, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, err
	}
	return &Recorder{ctx: ctx}, nil
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
	r.isRecording = false
	r.mu.Unlock()

	if r.device != nil {
		r.device.Stop()
		r.device.Uninit()
	}
}

// SaveWAV saves the recorded audio data to a WAV file at the specified path.
func (r *Recorder) SaveWAV(path string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write WAV header manually (44 bytes)
	writeWavHeader(f, len(r.data), 16000, 1)
	_, err = f.Write(r.data)
	return err
}

// writeWavHeader is a helper function to create the standard WAV header
func writeWavHeader(f *os.File, dataSize, sampleRate, channels int) {
	binary.Write(f, binary.LittleEndian, []byte("RIFF"))
	binary.Write(f, binary.LittleEndian, int32(36+dataSize))
	binary.Write(f, binary.LittleEndian, []byte("WAVE"))
	binary.Write(f, binary.LittleEndian, []byte("fmt "))
	binary.Write(f, binary.LittleEndian, int32(16))
	binary.Write(f, binary.LittleEndian, int16(1)) // Audio format (PCM)
	binary.Write(f, binary.LittleEndian, int16(channels))
	binary.Write(f, binary.LittleEndian, int32(sampleRate))
	binary.Write(f, binary.LittleEndian, int32(sampleRate*channels*2))
	binary.Write(f, binary.LittleEndian, int16(channels*2))
	binary.Write(f, binary.LittleEndian, int16(16)) // Bits por sample
	binary.Write(f, binary.LittleEndian, []byte("data"))
	binary.Write(f, binary.LittleEndian, int32(dataSize))
}
