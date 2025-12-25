package systray

import (
	"fyne.io/systray"
	"github.com/eduardolat/exp-stt/assets/logo"
	"github.com/eduardolat/exp-stt/internal/config"
)

type Instance struct {
	systrayStart func()
	systrayEnd   func()
}

func New() *Instance {
	i := &Instance{}

	systray.SetTitle(config.App.AppName)
	systray.SetTooltip(config.App.AppName)
	systray.SetIcon(logo.LogoBlackGray.PNG.Size512.Logo)

	start, end := systray.RunWithExternalLoop(i.onReady, i.onExit)
	i.systrayStart = start
	i.systrayEnd = end

	return i
}

func (s *Instance) onReady() {}

func (s *Instance) onExit() {}

func (s *Instance) Start() {
	s.systrayStart()
}

func (s *Instance) Shutdown() {
	s.systrayEnd()
}
