//go:build windows

package clipboard

import (
	"syscall"
	"unsafe"
)

var (
	user32        = syscall.NewLazyDLL("user32.dll")
	procSendInput = user32.NewProc("SendInput")
)

const (
	inputKeyboard = 1
	keyEventKeyUp = 0x0002
	vkControl     = 0x11
	vkV           = 0x56
)

type keyboardInput struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
}

type input struct {
	dtype   uint32
	ki      keyboardInput
	padding [8]byte
}

// triggerPastePlatform sends Ctrl+V directly via Windows API.
func triggerPastePlatform() error {
	var inputs []input

	// Helper to create input struct
	createKey := func(vk uint16, flags uint32) input {
		return input{
			dtype: inputKeyboard,
			ki: keyboardInput{
				wVk:     vk,
				dwFlags: flags,
			},
		}
	}

	// Sequence: Ctrl Down -> V Down -> V Up -> Ctrl Up
	inputs = append(inputs, createKey(vkControl, 0))
	inputs = append(inputs, createKey(vkV, 0))
	inputs = append(inputs, createKey(vkV, keyEventKeyUp))
	inputs = append(inputs, createKey(vkControl, keyEventKeyUp))

	cbSize := int(unsafe.Sizeof(inputs[0]))
	procSendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		uintptr(cbSize),
	)

	return nil
}
