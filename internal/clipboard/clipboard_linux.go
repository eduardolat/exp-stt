//go:build linux

package clipboard

import (
	"os/exec"
)

// triggerPastePlatform sends Ctrl+V using xdotool (requires xwayland on wayland).
func triggerPastePlatform() error {
	return exec.Command("xdotool", "key", "ctrl+v").Run()
}
