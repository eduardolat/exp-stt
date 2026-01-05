//go:build linux

package clipboard

import (
	"errors"
	"fmt"
	"os/exec"
	"sync"

	"github.com/varavelio/tribar/internal/state"
)

// pasteDriver represents the available paste drivers on Linux.
type pasteDriver int

const (
	pasteDriverNone pasteDriver = iota
	pasteDriverYdotool
	pasteDriverWtype
	pasteDriverXdotool
)

var (
	detectedDriver     pasteDriver
	detectedDriverErr  error
	detectedDriverOnce sync.Once
)

// detectPasteDriver determines which paste driver is available on the system.
// Priority order: ydotool (universal) > wtype (Wayland) > xdotool (X11/XWayland)
func detectPasteDriver() (pasteDriver, error) {
	detectedDriverOnce.Do(func() {
		if _, err := exec.LookPath("ydotool"); err == nil {
			detectedDriver = pasteDriverYdotool
			return
		}
		if state.RuntimeInfo.DisplayServer == "wayland" {
			if _, err := exec.LookPath("wtype"); err == nil {
				detectedDriver = pasteDriverWtype
				return
			}
		}
		if _, err := exec.LookPath("xdotool"); err == nil {
			detectedDriver = pasteDriverXdotool
			return
		}
		detectedDriver = pasteDriverNone

		if state.RuntimeInfo.DisplayServer == "wayland" {
			detectedDriverErr = errors.New("no paste driver available: install ydotool, wtype or xdotool (with xwayland)")
		} else {
			detectedDriverErr = errors.New("no paste driver available: install ydotool or xdotool")
		}
	})
	return detectedDriver, detectedDriverErr
}

// triggerPastePlatform sends the paste shortcut using the best available driver.
func triggerPastePlatform(shortcut PasteShortcut) error {
	driver, err := detectPasteDriver()
	if err != nil {
		return fmt.Errorf("failed to detect paste driver: %w", err)
	}

	switch driver {
	case pasteDriverYdotool:
		return triggerPasteYdotool(shortcut)
	case pasteDriverWtype:
		return triggerPasteWtype(shortcut)
	case pasteDriverXdotool:
		return triggerPasteXdotool(shortcut)
	default:
		return fmt.Errorf("no paste driver available")
	}
}

// triggerPasteYdotool sends the paste shortcut using ydotool, it uses raw keycodes
// from input-event-codes.h.
//
// Key format: "KEYCODE:1" for press, "KEYCODE:0" for release.
//
//   - https://github.com/ReimuNotMoe/ydotool
//   - https://github.com/torvalds/linux/blob/master/include/uapi/linux/input-event-codes.h
func triggerPasteYdotool(shortcut PasteShortcut) error {
	// Keycode reference (from input-event-codes.h):
	const (
		keyCtrl   = "29"
		keyShift  = "42"
		keyV      = "47"
		keyInsert = "110"
	)

	var args []string
	switch shortcut {
	case PasteShortcutCtrlShiftV:
		args = []string{
			"key",
			keyCtrl + ":1", keyShift + ":1", keyV + ":1", // Press Ctrl, Shift, V
			keyV + ":0", keyShift + ":0", keyCtrl + ":0", // Release V, Shift, Ctrl
		}
	case PasteShortcutShiftInsert:
		args = []string{
			"key",
			keyShift + ":1", keyInsert + ":1", // Press Shift, Insert
			keyInsert + ":0", keyShift + ":0", // Release Insert, Shift
		}
	default: // PasteShortcutCtrlV
		args = []string{
			"key",
			keyCtrl + ":1", keyV + ":1", // Press Ctrl, V
			keyV + ":0", keyCtrl + ":0", // Release V, Ctrl
		}
	}
	return exec.Command("ydotool", args...).Run()
}

// triggerPasteWtype sends the paste shortcut using wtype.
// wtype uses -M/-m for modifier press/release and -P/-p for named key press/release.
//
//   - https://github.com/atx/wtype
func triggerPasteWtype(shortcut PasteShortcut) error {
	var args []string
	switch shortcut {
	case PasteShortcutCtrlShiftV:
		args = []string{"-M", "ctrl", "-M", "shift", "v", "-m", "shift", "-m", "ctrl"}
	case PasteShortcutShiftInsert:
		args = []string{"-M", "shift", "-P", "Insert", "-p", "Insert", "-m", "shift"}
	default: // PasteShortcutCtrlV
		args = []string{"-M", "ctrl", "v", "-m", "ctrl"}
	}
	return exec.Command("wtype", args...).Run()
}

// triggerPasteXdotool sends the paste shortcut using xdotool.
// it works with X11 and XWayland.
//
//   - https://github.com/jordansissel/xdotool
func triggerPasteXdotool(shortcut PasteShortcut) error {
	var keyseq string
	switch shortcut {
	case PasteShortcutCtrlShiftV:
		keyseq = "ctrl+shift+v"
	case PasteShortcutShiftInsert:
		keyseq = "shift+Insert"
	default: // PasteShortcutCtrlV
		keyseq = "ctrl+v"
	}
	return exec.Command("xdotool", "key", "--clearmodifiers", keyseq).Run()
}
