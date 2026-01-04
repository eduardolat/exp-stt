package api

import (
	"fmt"
	"time"

	"github.com/varavelio/tribar/internal/server/api/uforpc"
)

const pingInterval = 30 * time.Second

func (h *handlers) registerStreamListenForEvents() {
	h.uforpcServer.Streams.ListenForEvents.Handle(func(c *uforpc.ListenForEventsHandlerContext[urpcProps], emit uforpc.ListenForEventsEmitFunc[urpcProps]) error {
		// Send initial state
		historyEntries := h.appState.GetHistory(c.Context)
		initialState := h.buildState(historyEntries)
		if err := emit(c, uforpc.ListenForEventsOutput{
			EventType:    "stateUpdated",
			StateUpdated: uforpc.Optional[uforpc.State]{Present: true, Value: initialState},
		}); err != nil {
			return fmt.Errorf("failed to emit initial state: %w", err)
		}

		// Channels for receiving event signals
		stateCh := make(chan any, 1)
		settingsCh := make(chan any, 1)

		// Subscribe to state changes
		stateHandler := func() {
			select {
			case stateCh <- nil:
			default: // drop if channel is full to avoid blocking
			}
		}
		_ = h.eventBus.SubscribeStateChanged(stateHandler)
		defer func() { _ = h.eventBus.UnsubscribeStateChanged(stateHandler) }()

		// Subscribe to settings changes
		settingsHandler := func() {
			select {
			case settingsCh <- nil:
			default: // drop if channel is full to avoid blocking
			}
		}
		_ = h.eventBus.SubscribeSettingsChanged(settingsHandler)
		defer func() { _ = h.eventBus.UnsubscribeSettingsChanged(settingsHandler) }()

		pingTicker := time.NewTicker(pingInterval)
		defer pingTicker.Stop()

		for {
			select {
			case <-c.Context.Done():
				return nil

			case <-pingTicker.C:
				if err := emit(c, uforpc.ListenForEventsOutput{
					EventType: "ping",
				}); err != nil {
					return fmt.Errorf("failed to emit ping: %w", err)
				}

			case <-stateCh:
				historyEntries := h.appState.GetHistory(c.Context)
				newState := h.buildState(historyEntries)
				if err := emit(c, uforpc.ListenForEventsOutput{
					EventType:    "stateUpdated",
					StateUpdated: uforpc.Optional[uforpc.State]{Present: true, Value: newState},
				}); err != nil {
					return fmt.Errorf("failed to emit state update: %w", err)
				}

			case <-settingsCh:
				settings := h.settingsManager.Get()
				apiSettings := buildSettings(settings)
				if err := emit(c, uforpc.ListenForEventsOutput{
					EventType:       "settingsUpdated",
					SettingsUpdated: uforpc.Optional[uforpc.Settings]{Present: true, Value: apiSettings},
				}); err != nil {
					return fmt.Errorf("failed to emit settings update: %w", err)
				}
			}
		}
	})
}
