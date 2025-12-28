//go:build linux

package notify

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/varavelio/tribar/assets/logo"
	"github.com/varavelio/tribar/internal/config"
)

var (
	iconPath string
	once     sync.Once
)

// send dispatches a system notification via notify-send.
func (n *Instance) send(ctx context.Context, title, message string) {
	// Extract notification icon to temporary storage once.
	once.Do(func() {
		tempPath := filepath.Join(os.TempDir(), "tribar-notify-icon.png")
		if err := os.WriteFile(tempPath, logo.LogoBlackWhite.PNG.Size128.Logo, 0644); err == nil {
			iconPath = tempPath
		} else {
			n.logger.Error(ctx, "failed to create temporary notification icon", "err", err)
		}
	})

	args := []string{
		"-a", fmt.Sprintf("%s v%s", config.AppName, config.AppVersion),
		title, message,
	}
	if iconPath != "" {
		args = append(args, "-i", iconPath)
	}

	cmd := exec.CommandContext(ctx, "notify-send", args...)
	if err := cmd.Run(); err != nil {
		n.logger.Error(
			ctx, "failed to send desktop notification",
			"title", title,
			"message", message,
			"err", err,
		)
	}
}
