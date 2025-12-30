package shortcut

import (
	"context"
	"sync"

	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/logger"
	"github.com/varavelio/tribar/internal/state"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

// Manager handles global hotkey registration and events.
type Manager struct {
	logger          logger.Logger
	settingsManager *config.SettingsManager
	onToggle        func()

	mu       sync.Mutex
	hk       *hotkey.Hotkey
	stopCh   chan struct{}
	active   bool
	disabled bool // true if Wayland or unsupported
}

// New creates a new shortcut Manager.
func New(logger logger.Logger, settingsManager *config.SettingsManager, onToggle func()) *Manager {
	displayServer := state.RuntimeInfo.DisplayServer
	disabled := displayServer == "wayland" || displayServer == "unknown"

	if disabled {
		logger.Warn(
			context.Background(), "hotkeys disabled because your display server does not support them",
			"display_server", displayServer,
		)
	}

	return &Manager{
		logger:          logger,
		settingsManager: settingsManager,
		onToggle:        onToggle,
		disabled:        disabled,
	}
}

// IsActive returns whether a hotkey is currently registered.
func (m *Manager) IsActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.active
}

// Start registers the hotkey from settings and begins listening.
// This function must be called from within mainthread.Init context.
func (m *Manager) Start() error {
	if m.disabled {
		return nil
	}

	settings := m.settingsManager.Get()
	return m.register(settings.ShortcutToggle)
}

// Stop unregisters the current hotkey and stops listening.
func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.stopCh != nil {
		close(m.stopCh)
		m.stopCh = nil
	}

	if m.hk != nil {
		_ = m.hk.Unregister()
		m.hk = nil
	}

	m.active = false
}

// Update changes the hotkey to a new combination, saves to settings, and re-registers.
func (m *Manager) Update(shortcut config.Shortcut) error {
	if m.disabled {
		return nil
	}

	m.Stop()

	settings := m.settingsManager.Get()
	settings.ShortcutToggle = shortcut
	if err := m.settingsManager.Update(settings); err != nil {
		return err
	}

	return m.register(shortcut)
}

// register creates and registers a hotkey, then starts listening in a goroutine.
func (m *Manager) register(shortcut config.Shortcut) error {
	mods, key, err := ParseShortcut(shortcut)
	if err != nil {
		m.logger.Error(context.Background(), "failed to parse shortcut", "err", err)
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Use mainthread.Call to register since hotkey requires main thread on macOS
	var regErr error
	mainthread.Call(func() {
		m.hk = hotkey.New(mods, key)
		regErr = m.hk.Register()
	})

	if regErr != nil {
		m.logger.Error(context.Background(), "failed to register hotkey", "err", regErr)
		m.hk = nil
		return regErr
	}

	m.active = true
	m.stopCh = make(chan struct{})

	m.logger.Info(context.Background(), "hotkey registered", "shortcut", shortcut)

	// Start listening in background
	go m.listen(m.stopCh)

	return nil
}

// listen waits for hotkey events and calls the toggle callback.
func (m *Manager) listen(stopCh chan struct{}) {
	for {
		select {
		case <-stopCh:
			return
		case <-m.hk.Keydown():
			m.logger.Debug(context.Background(), "hotkey triggered")
			m.onToggle()
		}
	}
}
