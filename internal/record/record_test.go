package record

import (
	"encoding/binary"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildWavBytes(t *testing.T) {
	t.Run("Valid WAV Header Generation", func(t *testing.T) {
		// 16kHz, 1 channel, 2 bytes per sample (16-bit)
		sampleRate := 16000
		channels := 1
		// Create 4 bytes of data (2 samples)
		data := []byte{0x00, 0x00, 0xFF, 0x7F} // 0, 32767

		wavData := buildWavBytes(data, sampleRate, channels)

		require.NotNil(t, wavData)
		// Header is 44 bytes
		assert.Equal(t, 44+len(data), len(wavData))

		// Verify RIFF header
		assert.Equal(t, "RIFF", string(wavData[0:4]))
		// File size: 36 + dataSize
		expectedFileSize := uint32(36 + len(data))
		assert.Equal(t, expectedFileSize, binary.LittleEndian.Uint32(wavData[4:8]))
		assert.Equal(t, "WAVE", string(wavData[8:12]))

		// Verify fmt chunk
		assert.Equal(t, "fmt ", string(wavData[12:16]))
		assert.Equal(t, uint32(16), binary.LittleEndian.Uint32(wavData[16:20])) // Subchunk1Size
		assert.Equal(t, uint16(1), binary.LittleEndian.Uint16(wavData[20:22]))  // AudioFormat (PCM)
		assert.Equal(t, uint16(channels), binary.LittleEndian.Uint16(wavData[22:24]))
		assert.Equal(t, uint32(sampleRate), binary.LittleEndian.Uint32(wavData[24:28]))

		// ByteRate = SampleRate * NumChannels * BitsPerSample/8
		expectedByteRate := uint32(sampleRate * channels * 16 / 8)
		assert.Equal(t, expectedByteRate, binary.LittleEndian.Uint32(wavData[28:32]))

		// BlockAlign = NumChannels * BitsPerSample/8
		expectedBlockAlign := uint16(channels * 16 / 8)
		assert.Equal(t, expectedBlockAlign, binary.LittleEndian.Uint16(wavData[32:34]))

		assert.Equal(t, uint16(16), binary.LittleEndian.Uint16(wavData[34:36])) // BitsPerSample

		// Verify data chunk
		assert.Equal(t, "data", string(wavData[36:40]))
		assert.Equal(t, uint32(len(data)), binary.LittleEndian.Uint32(wavData[40:44]))

		// Verify payload
		assert.Equal(t, data, wavData[44:])
	})

	t.Run("Empty Data", func(t *testing.T) {
		data := []byte{}
		wavData := buildWavBytes(data, 44100, 2)

		assert.Equal(t, 44, len(wavData))
		assert.Equal(t, uint32(36), binary.LittleEndian.Uint32(wavData[4:8]))  // 36 + 0
		assert.Equal(t, uint32(0), binary.LittleEndian.Uint32(wavData[40:44])) // Data size 0
	})
}

func TestParseDeviceID(t *testing.T) {
	t.Run("Valid Device ID", func(t *testing.T) {
		// Valid hex string (usually 32 chars for 16 bytes? depends on malgo.DeviceID size)
		// Assuming DeviceID is large enough.
		hexStr := "0102030405060708"
		id := parseDeviceID(hexStr)

		require.NotNil(t, id)
		// Check first few bytes
		assert.Equal(t, byte(0x01), id[0])
		assert.Equal(t, byte(0x08), id[7])
	})

	t.Run("Invalid Hex String", func(t *testing.T) {
		id := parseDeviceID("invalid-hex")
		assert.Nil(t, id)
	})

	t.Run("Short Hex String", func(t *testing.T) {
		// Should verify it pads with zero or just fills what it has
		// The implementation copies what it has.
		hexStr := "AB"
		id := parseDeviceID(hexStr)
		require.NotNil(t, id)
		assert.Equal(t, byte(0xAB), id[0])
		assert.Equal(t, byte(0x00), id[1]) // Should remain zero
	})

	t.Run("Long Hex String", func(t *testing.T) {
		// Implementation truncates if longer than ID len.
		// Construct a very long string.
		longHex := hex.EncodeToString(make([]byte, 100))
		id := parseDeviceID(longHex)
		require.NotNil(t, id)
		// We don't check exact length here since we don't know malgo.DeviceID size easily without reflecting,
		// but we know it shouldn't panic.
	})
}
