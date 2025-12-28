//go:build darwin

package notify

import (
	"context"
	"fmt"
	"sync"

	"github.com/gen2brain/beeep"
	"github.com/varavelio/tribar/assets/logo"
	"github.com/varavelio/tribar/internal/config"
)

var once sync.Once

// send dispatches a notification to the desktop.
func (n *Instance) send(ctx context.Context, title, message string) {
	once.Do(func() {
		beeep.AppName = fmt.Sprintf("%s v%s", config.AppName, config.AppVersion)
	})

	if err := beeep.Notify(title, message, logo.LogoBlackWhite.PNG.Size128.Logo); err != nil {
		n.logger.Error(
			ctx, "failed to send desktop notification",
			"title", title,
			"message", message,
			"err", err,
		)
	}
}
