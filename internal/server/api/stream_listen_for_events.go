package api

import (
	"reflect"
	"time"

	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/server/api/uforpc"
	"github.com/varavelio/tribar/internal/state"
)

const (
	pollInterval = 500 * time.Millisecond
	pingInterval = 30 * time.Second
)

func (h *handlers) registerStreamListenForEvents() {
	h.uforpcServer.Streams.ListenForEvents.Handle(func(c *uforpc.ListenForEventsHandlerContext[urpcProps], emit uforpc.ListenForEventsEmitFunc[urpcProps]) error {
		// Track previous state to detect changes
		var prevStatus state.Status
		var prevHistoryCount int
		var prevSettings config.Settings

		// Initialize with current state
		prevStatus, _ = h.appState.GetStatus()
		prevHistoryCount = h.appState.HistoryCount()
		prevSettings = h.settingsManager.Get()

		// Send initial state
		entries := h.appState.GetHistory(c.Context)
		initialState := h.buildState(entries)
		if err := emit(c, uforpc.ListenForEventsOutput{
			EventType:    "stateUpdated",
			StateUpdated: uforpc.Optional[uforpc.State]{Present: true, Value: initialState},
		}); err != nil {
			return err
		}

		pollTicker := time.NewTicker(pollInterval)
		defer pollTicker.Stop()

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
					return err
				}

			case <-pollTicker.C:
				// Check for state changes
				currentStatus, _ := h.appState.GetStatus()
				currentHistoryCount := h.appState.HistoryCount()

				stateChanged := currentStatus != prevStatus || currentHistoryCount != prevHistoryCount

				if stateChanged {
					prevStatus = currentStatus
					prevHistoryCount = currentHistoryCount
					entries := h.appState.GetHistory(c.Context)
					newState := h.buildState(entries)
					if err := emit(c, uforpc.ListenForEventsOutput{
						EventType:    "stateUpdated",
						StateUpdated: uforpc.Optional[uforpc.State]{Present: true, Value: newState},
					}); err != nil {
						return err
					}
				}

				// Check for settings changes
				currentSettings := h.settingsManager.Get()
				if !reflect.DeepEqual(currentSettings, prevSettings) {
					prevSettings = currentSettings
					apiSettings := buildSettings(currentSettings)
					if err := emit(c, uforpc.ListenForEventsOutput{
						EventType:       "settingsUpdated",
						SettingsUpdated: uforpc.Optional[uforpc.Settings]{Present: true, Value: apiSettings},
					}); err != nil {
						return err
					}
				}
			}
		}
	})
}
