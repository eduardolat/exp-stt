//go:build darwin

package shortcut

import (
	"golang.design/x/hotkey"
)

// modMap maps string modifier names to hotkey.Modifier constants for macOS.
// On macOS, ModOption is Alt, and ModCmd is Command (⌘).
var modMap = map[string]hotkey.Modifier{
	"ctrl":  hotkey.ModCtrl,
	"alt":   hotkey.ModOption,
	"shift": hotkey.ModShift,
	"meta":  hotkey.ModCmd,
}
