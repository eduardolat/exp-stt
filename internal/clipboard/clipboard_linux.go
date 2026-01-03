//go:build linux

package clipboard

import (
	"os/exec"
)

// triggerPastePlatform sends the paste shortcut using xdotool (requires xwayland on wayland).
func triggerPastePlatform(shortcut PasteShortcut) error {
	// Map the enum to xdotool key names
	var keyseq string
	switch shortcut {
	case PasteShortcutCtrlShiftV:
		keyseq = "ctrl+shift+v"
	case PasteShortcutShiftInsert:
		keyseq = "shift+Insert"
	default: // PasteShortcutCtrlV
		keyseq = "ctrl+v"
	}
	return exec.Command("xdotool", "key", keyseq).Run()
}
