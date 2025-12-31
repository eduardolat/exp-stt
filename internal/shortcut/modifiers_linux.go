//go:build linux

package shortcut

import (
	"golang.design/x/hotkey"
)

// modMap maps string modifier names to hotkey.Modifier constants for Linux.
// On X11, Mod1 is typically Alt, and Mod4 is typically Super/Meta/Win.
var modMap = map[string]hotkey.Modifier{
	"ctrl":  hotkey.ModCtrl,
	"alt":   hotkey.Mod1,
	"shift": hotkey.ModShift,
	"meta":  hotkey.Mod4,
}
