//go:build darwin

package clipboard

import "os/exec"

// triggerPastePlatform sends the paste shortcut using AppleScript.
func triggerPastePlatform(shortcut PasteShortcut) error {
	var script string
	switch shortcut {
	case PasteShortcutCtrlShiftV:
		// Cmd+Shift+V for paste without formatting
		script = `tell application "System Events" to keystroke "v" using {command down, shift down}`
	case PasteShortcutShiftInsert:
		// Shift+Insert (key code 114 is Insert/Help key on macOS)
		script = `tell application "System Events" to key code 114 using {shift down}`
	default: // PasteShortcutCtrlV
		// Cmd+V is the standard paste on macOS
		script = `tell application "System Events" to keystroke "v" using {command down}`
	}
	return exec.Command("osascript", "-e", script).Run()
}
