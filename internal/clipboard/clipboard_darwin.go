//go:build darwin

package clipboard

import "os/exec"

// triggerPastePlatform sends Cmd+V using AppleScript.
func triggerPastePlatform() error {
	script := `tell application "System Events" to keystroke "v" using {command down}`
	return exec.Command("osascript", "-e", script).Run()
}
