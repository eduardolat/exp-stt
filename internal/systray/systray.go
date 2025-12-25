package systray

import (
	"time"

	"fyne.io/systray"
	"github.com/eduardolat/exp-stt/assets/logo"
	"github.com/eduardolat/exp-stt/internal/config"
)

const animationFrameDuration = time.Millisecond * 200

type animationPosition int

const (
	animationPositionUnknown animationPosition = iota
	animationPositionMiddle
	animationPositionRight
	animationPositionLeft
)

type status int

const (
	statusUnknown status = iota
	statusUnloaded
	statusLoading
	statusLoaded
	statusListening
	statusTranscribing
	statusPostProcessing
)

type Instance struct {
	systrayStart func()
	systrayEnd   func()

	statusCurr status
	statusPrev status

	animationForward bool // True if animation is moving forward, false if backward
	animationPosCurr animationPosition
	animationPosPrev animationPosition
	animationTimer   *time.Timer

	isShuttingDown bool
}

func New() *Instance {
	i := &Instance{
		statusCurr:       statusUnloaded,
		animationPosCurr: animationPositionMiddle,
		animationTimer:   time.NewTimer(0),
	}

	start, end := systray.RunWithExternalLoop(i.onReady, func() {})
	i.systrayStart = start
	i.systrayEnd = end

	return i
}

func (i *Instance) onReady() {
	i.SetStatus(statusUnloaded)
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

	if i.statusCurr == statusUnloaded || i.statusCurr == statusLoaded {
		i.animationPosCurr = animationPositionMiddle
		return
	}

	// Ping-Pong animation between left, middle, and right positions
	switch i.animationPosCurr {
	case animationPositionMiddle:
		if i.animationForward {
			i.animationPosCurr = animationPositionRight
		} else {
			i.animationPosCurr = animationPositionLeft
		}
	case animationPositionRight:
		i.animationPosCurr = animationPositionMiddle
		i.animationForward = false
	case animationPositionLeft:
		i.animationPosCurr = animationPositionMiddle
		i.animationForward = true
	}
}

// SetStatus changes the current status of the systray instance.
func (i *Instance) SetStatus(newStatus status) {
	i.statusPrev = i.statusCurr
	i.statusCurr = newStatus
}

// setTitle updates the systray title and tooltip based on the current status.
func (i *Instance) setTitle() {
	title := config.App.AppName

	switch i.statusCurr {
	case statusUnloaded:
		title += " - Model not loaded"
	case statusLoading:
		title += " - Loading model..."
	case statusLoaded:
		title += " - Model loaded"
	case statusListening:
		title += " - Listening..."
	case statusTranscribing:
		title += " - Transcribing..."
	case statusPostProcessing:
		title += " - Post-processing..."
	}

	systray.SetTitle(title)
	systray.SetTooltip(title)
}

// setIcon updates the systray icon based on the current status and animation position.
func (i *Instance) setIcon() {
	res := logo.LogoBlackGray.PNG.Size32

	switch i.statusCurr {
	case statusUnloaded:
		res = logo.LogoBlackGray.PNG.Size32
	case statusLoading:
		res = logo.LogoBlackAmber.PNG.Size32
	case statusLoaded:
		res = logo.LogoBlackWhite.PNG.Size32
	case statusListening:
		res = logo.LogoBlackPink.PNG.Size32
	case statusTranscribing:
		res = logo.LogoBlackBlue.PNG.Size32
	case statusPostProcessing:
		res = logo.LogoBlackGreen.PNG.Size32
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

		if i.statusPrev != i.statusCurr {
			i.setTitle()
		}

		if i.statusPrev != i.statusCurr || i.animationPosPrev != i.animationPosCurr {
			i.setIcon()
		}

		i.setNextAnimationPosition()
		i.animationTimer.Reset(animationFrameDuration)
	}
}
