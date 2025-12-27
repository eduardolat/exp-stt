package systray

import (
	"runtime"
	"time"

	"fyne.io/systray"
	"github.com/varavelio/tribar/assets/logo"
	"github.com/varavelio/tribar/internal/app"
	"github.com/varavelio/tribar/internal/config"
)

const animationFrameDuration = time.Millisecond * 200

type animationPosition int

const (
	animationPositionUnknown animationPosition = iota
	animationPositionMiddle
	animationPositionRight
	animationPositionLeft
)

type Instance struct {
	app *app.Instance

	systrayStart func()
	systrayEnd   func()

	animationBackward bool // True if animation is moving backward, false if forward
	animationPosCurr  animationPosition
	animationPosPrev  animationPosition
	animationTimer    *time.Timer

	isShuttingDown bool
}

func New(app *app.Instance) *Instance {
	i := &Instance{
		app:              app,
		animationPosCurr: animationPositionMiddle,
		animationTimer:   time.NewTimer(0),
	}

	start, end := systray.RunWithExternalLoop(i.onReady, func() {})
	i.systrayStart = start
	i.systrayEnd = end

	return i
}

func (i *Instance) onReady() {
	i.animate()
}

func (i *Instance) Start() {
	i.systrayStart()
}

func (i *Instance) Shutdown() {
	i.isShuttingDown = true
	i.animationTimer.Stop()
	i.systrayEnd()
}

// setNextAnimationPosition advances the animation position.
//
// For unloaded and loaded statuses, the animation position is always set to middle, for
// other statuses, it cycles through middle, right, and left positions.
func (i *Instance) setNextAnimationPosition() {
	i.animationPosPrev = i.animationPosCurr

	if i.app.StatusCurrent == app.StatusUnloaded || i.app.StatusCurrent == app.StatusLoaded {
		i.animationPosCurr = animationPositionMiddle
		return
	}

	// Ping-Pong animation between left, middle, and right positions
	switch i.animationPosCurr {
	case animationPositionMiddle:
		if i.animationBackward {
			i.animationPosCurr = animationPositionLeft
		} else {
			i.animationPosCurr = animationPositionRight
		}
	case animationPositionRight:
		i.animationPosCurr = animationPositionMiddle
		i.animationBackward = true
	case animationPositionLeft:
		i.animationPosCurr = animationPositionMiddle
		i.animationBackward = false
	}
}

// setTitle updates the systray title and tooltip based on the current status.
func (i *Instance) setTitle() {
	title := config.AppName

	switch i.app.StatusCurrent {
	case app.StatusUnloaded:
		title += " - Model not loaded"
	case app.StatusLoading:
		title += " - Loading model..."
	case app.StatusLoaded:
		title += " - Model loaded"
	case app.StatusListening:
		title += " - Listening..."
	case app.StatusTranscribing:
		title += " - Transcribing..."
	case app.StatusPostProcessing:
		title += " - Post-processing..."
	}

	systray.SetTitle(title)
	systray.SetTooltip(title)
}

// setIcon updates the systray icon based on the current status and animation position.
func (i *Instance) setIcon() {
	pngOrIco := func(logoRes logo.LogoResources) logo.ResourceSet {
		if runtime.GOOS == "windows" {
			return logoRes.ICO
		}
		return logoRes.PNG.Size32
	}

	res := pngOrIco(logo.LogoBlackGray)

	switch i.app.StatusCurrent {
	case app.StatusUnloaded:
		res = pngOrIco(logo.LogoBlackGray)
	case app.StatusLoading:
		res = pngOrIco(logo.LogoBlackAmber)
	case app.StatusLoaded:
		res = pngOrIco(logo.LogoBlackWhite)
	case app.StatusListening:
		res = pngOrIco(logo.LogoBlackPink)
	case app.StatusTranscribing:
		res = pngOrIco(logo.LogoBlackBlue)
	case app.StatusPostProcessing:
		res = pngOrIco(logo.LogoBlackGreen)
	}

	switch i.animationPosCurr {
	case animationPositionLeft:
		systray.SetIcon(res.Left)
	case animationPositionMiddle:
		systray.SetIcon(res.Middle)
	case animationPositionRight:
		systray.SetIcon(res.Right)
	}
}

// animate runs the animation loop, updating the systray icon and title based on the current status
// and animation position at regular intervals defined by animationFrameDuration.
func (i *Instance) animate() {
	for range i.animationTimer.C {
		if i.isShuttingDown {
			return
		}

		if i.app.StatusPrevious != i.app.StatusCurrent {
			i.setTitle()
		}

		if i.app.StatusPrevious != i.app.StatusCurrent || i.animationPosPrev != i.animationPosCurr {
			i.setIcon()
		}

		i.setNextAnimationPosition()
		i.animationTimer.Reset(animationFrameDuration)
	}
}
