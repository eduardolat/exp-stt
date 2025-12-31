//go:build windows

package shortcut

import (
	"golang.design/x/hotkey"
)

// modMap maps string modifier names to hotkey.Modifier constants for Windows.
// On Windows, meta is the Windows key, the other modifiers are the same as the constants.
var modMap = map[string]hotkey.Modifier{
	"ctrl":  hotkey.ModCtrl,
	"alt":   hotkey.ModAlt,
	"shift": hotkey.ModShift,
	"meta":  hotkey.ModWin,
}
