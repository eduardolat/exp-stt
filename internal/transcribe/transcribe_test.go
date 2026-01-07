package transcribe

import (
	"os"
	"testing"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper to generate a valid WAV buffer
func generateWAV(t *testing.T, sampleRate int, bitDepth int, numChannels int, data []int) []byte {
	// Create a temp file since wav.Encoder needs a WriteSeeker
	tmpFile, err := os.CreateTemp("", "genwav*.wav")
	require.NoError(t, err)
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	enc := wav.NewEncoder(tmpFile, sampleRate, bitDepth, numChannels, 1)

	// Create int buffer
	audioBuf := &audio.IntBuffer{
		Format: &audio.Format{
			SampleRate:  sampleRate,
			NumChannels: numChannels,
		},
		Data:           data,
		SourceBitDepth: bitDepth,
	}

	err = enc.Write(audioBuf)
	require.NoError(t, err)
	err = enc.Close()
	require.NoError(t, err)

	_ = tmpFile.Close()

	// Read back the file content
	content, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)
	return content
}

func TestProcessWAVBytes(t *testing.T) {
	t.Run("Valid 16-bit Mono WAV", func(t *testing.T) {
		// Create a simple WAV with known values
		// Max positive 16-bit: 32767 -> should be ~1.0
		// Max negative 16-bit: -32768 -> should be -1.0
		// Zero: 0 -> 0.0
		inputData := []int{0, 32767, -32768, 16384}
		wavData := generateWAV(t, 16000, 16, 1, inputData)

		samples, err := processWAVBytes(wavData)
		require.NoError(t, err)
		require.Len(t, samples, 4)

		// Check normalization
		assert.InDelta(t, 0.0, samples[0], 0.0001)
		assert.InDelta(t, 1.0, samples[1], 0.0001)  // 32767 / 32768 ~= 1.0
		assert.InDelta(t, -1.0, samples[2], 0.0001)  // -32768 / 32768
		assert.InDelta(t, 0.5, samples[3], 0.0001)   // 16384 / 32768
	})

	t.Run("Empty WAV Data", func(t *testing.T) {
		samples, err := processWAVBytes([]byte{})
		require.Error(t, err)
		assert.Nil(t, samples)
		assert.Contains(t, err.Error(), "invalid WAV file") // Or specific error from library
	})

	t.Run("Garbage Data", func(t *testing.T) {
		samples, err := processWAVBytes([]byte("not a wav file"))
		require.Error(t, err)
		assert.Nil(t, samples)
	})

	t.Run("Invalid Header", func(t *testing.T) {
		// Valid RIFF header but not WAV
		data := []byte("RIFF\x24\x00\x00\x00WAVEfmt \x10\x00\x00\x00\x01\x00\x01\x00\x80>\x00\x00\x00}\x00\x00\x02\x00\x10\x00data\x00\x00\x00\x00")
		// Corrupt it slightly
		data[0] = 'X'
		samples, err := processWAVBytes(data)
		require.Error(t, err)
		assert.Nil(t, samples)
	})

	t.Run("Stereo WAV (Should process but logic assumes mono)", func(t *testing.T) {
		// The current implementation takes FullPCMBuffer().Data which interleaves stereo samples.
		// It doesn't explicitly check for mono, so it will just normalize all samples.
		// If the requirement is strict 16kHz Mono, this test documents current behavior.
		inputData := []int{100, 200, 300, 400} // L R L R
		wavData := generateWAV(t, 16000, 16, 2, inputData)

		samples, err := processWAVBytes(wavData)
		require.NoError(t, err)
		require.Len(t, samples, 4)
		// It just treats them as a stream of samples
		assert.InDelta(t, 100.0/32768.0, samples[0], 0.0001)
	})

	t.Run("Wrong Sample Rate (44.1kHz)", func(t *testing.T) {
		// Code accepts any sample rate, just reads raw samples.
		// Testing to ensure it doesn't crash.
		inputData := []int{100, 200, 300, 400}
		wavData := generateWAV(t, 44100, 16, 1, inputData)

		samples, err := processWAVBytes(wavData)
		require.NoError(t, err)
		require.Len(t, samples, 4)
	})

	t.Run("8-bit Audio", func(t *testing.T) {
		// 8-bit WAVs are usually unsigned 0-255. Center is 128.
		// go-audio/wav decoder handles conversion to IntBuffer.
		// However, processWAVBytes divides by 32768.0, which assumes 16-bit.
		// If the decoder returns small values (0-255), the output samples will be tiny (~0.003).
		// This test documents that behavior (or correct behavior if decoder scales it up).
		inputData := []int{0, 128, 255}
		wavData := generateWAV(t, 16000, 8, 1, inputData)

		samples, err := processWAVBytes(wavData)
		require.NoError(t, err)
		require.Len(t, samples, 3)
		// If decoder returns raw 8-bit values (0-255), normalized values will be small.
		// Let's check what we get. If this fails, we know how it behaves.
		// Note: wav.Decoder might return 8-bit data in IntBuffer.Data as is.
		assert.InDelta(t, 0.0, samples[0], 0.01) // 0/32768
	})

	t.Run("Truncated Data", func(t *testing.T) {
		// Generate valid WAV then chop off the end
		inputData := []int{1, 2, 3, 4}
		wavData := generateWAV(t, 16000, 16, 1, inputData)
		// Remove nearly all bytes (header corruption)
		truncated := wavData[:10]

		// It should error
		samples, err := processWAVBytes(truncated)
		require.Error(t, err)
		assert.Nil(t, samples)
	})
}

func TestReadWAVFile(t *testing.T) {
	// Setup temp file
	tmpFile, err := os.CreateTemp("", "test*.wav")
	require.NoError(t, err)
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	// Write valid WAV
	wavData := generateWAV(t, 16000, 16, 1, []int{1, 2, 3})
	_, err = tmpFile.Write(wavData)
	require.NoError(t, err)
	_ = tmpFile.Close()

	// Test ReadWAVFile
	data, err := ReadWAVFile(tmpFile.Name())
	require.NoError(t, err)
	assert.Equal(t, wavData, data)

	// Test non-existent file
	_, err = ReadWAVFile("nonexistent.wav")
	require.Error(t, err)
}
