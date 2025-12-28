package systray

import (
	"runtime"
	"time"

	"fyne.io/systray"
	"github.com/varavelio/tribar/assets/logo"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/state"
)

const animationFrameDuration = time.Millisecond * 200

type animationPosition int

const (
	animationPositionUnknown animationPosition = iota
	animationPositionMiddle
	animationPositionRight
	animationPositionLeft
)

// Engine defines the interface for engine actions that systray can trigger.
type Engine interface {
	ToggleRecording()
}

type Instance struct {
	appState *state.Instance
	engine   Engine
	onQuit   func()

	systrayStart func()
	systrayEnd   func()

	animationBackward bool // True if animation is moving backward, false if forward
	animationPosCurr  animationPosition
	animationPosPrev  animationPosition
	animationTimer    *time.Timer

	isShuttingDown bool

	menuRecord *systray.MenuItem
	menuQuit   *systray.MenuItem
}

func New(appState *state.Instance, engine Engine, onQuit func()) *Instance {
	i := &Instance{
		appState:         appState,
		engine:           engine,
		onQuit:           onQuit,
		animationPosCurr: animationPositionMiddle,
		animationTimer:   time.NewTimer(0),
	}

	start, end := systray.RunWithExternalLoop(i.onReady, func() {})
	i.systrayStart = start
	i.systrayEnd = end

	return i
}

func (i *Instance) onReady() {
	i.setIcon()
	i.setTitle()

	versionItem := systray.AddMenuItem(config.AppName+" v"+config.AppVersion, "")
	versionItem.Disable()
	systray.AddSeparator()

	i.menuRecord = systray.AddMenuItem("Toggle Recording", "Start or stop recording")
	systray.AddSeparator()
	i.menuQuit = systray.AddMenuItem("Quit", "Exit the application")

	go i.handleMenuClicks()
	go i.animate()
}

func (i *Instance) handleMenuClicks() {
	for {
		select {
		case <-i.menuRecord.ClickedCh:
			if i.engine != nil {
				i.engine.ToggleRecording()
			}
		case <-i.menuQuit.ClickedCh:
			if i.onQuit != nil {
				i.onQuit()
			}
			return
		}
	}
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
	statusCurrent, _ := i.appState.GetStatus()
	i.animationPosPrev = i.animationPosCurr

	if statusCurrent == state.StatusUnloaded || statusCurrent == state.StatusLoaded {
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
	statusCurrent, _ := i.appState.GetStatus()
	title := config.AppName

	switch statusCurrent {
	case state.StatusUnloaded:
		title += " - Model not loaded"
	case state.StatusLoading:
		title += " - Loading model..."
	case state.StatusLoaded:
		title += " - Model loaded"
	case state.StatusListening:
		title += " - Listening..."
	case state.StatusTranscribing:
		title += " - Transcribing..."
	case state.StatusPostProcessing:
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
	statusCurrent, _ := i.appState.GetStatus()

	switch statusCurrent {
	case state.StatusUnloaded:
		res = pngOrIco(logo.LogoBlackGray)
	case state.StatusLoading:
		res = pngOrIco(logo.LogoBlackAmber)
	case state.StatusLoaded:
		res = pngOrIco(logo.LogoBlackWhite)
	case state.StatusListening:
		res = pngOrIco(logo.LogoBlackPink)
	case state.StatusTranscribing:
		res = pngOrIco(logo.LogoBlackBlue)
	case state.StatusPostProcessing:
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

		statusCurrent, statusPrevious := i.appState.GetStatus()

		if statusPrevious != statusCurrent {
			i.setTitle()
		}

		if statusPrevious != statusCurrent || i.animationPosPrev != i.animationPosCurr {
			i.setIcon()
		}

		i.setNextAnimationPosition()
		i.animationTimer.Reset(animationFrameDuration)
	}
}
