package api

import (
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/history"
	"github.com/varavelio/tribar/internal/server/api/uforpc"
	"github.com/varavelio/tribar/internal/state"
)

// statusToString converts a state.Status enum to its API string representation.
func statusToString(s state.Status) string {
	switch s {
	case state.StatusUnknown:
		return "unknown"
	case state.StatusUnloaded:
		return "unloaded"
	case state.StatusDownloading:
		return "downloading"
	case state.StatusLoading:
		return "loading"
	case state.StatusLoaded:
		return "loaded"
	case state.StatusListening:
		return "listening"
	case state.StatusTranscribing:
		return "transcribing"
	case state.StatusPostProcessing:
		return "post_processing"
	default:
		return "unknown"
	}
}

// buildState constructs the full API State from the current application state.
func (h *handlers) buildState(entries []history.Entry) uforpc.State {
	currentStatus, _ := h.appState.GetStatus()
	downloadProgress := h.appState.GetDownloadProgress()
	devices := h.appState.GetAvailableDevices()

	return uforpc.State{
		Status: statusToString(currentStatus),
		DownloadProgress: uforpc.DownloadProgress{
			FileName:   downloadProgress.FileName,
			Downloaded: int(downloadProgress.Downloaded),
			Total:      int(downloadProgress.Total),
			Percent:    downloadProgress.Percent,
		},
		Devices: uforpc.AvailableDevices{
			InputDevices:  mapAudioDevices(devices.InputDevices),
			OutputDevices: mapAudioDevices(devices.OutputDevices),
		},
		History:    mapHistoryEntries(entries),
		SystemInfo: systemInfoToAPI(state.RuntimeInfo),
	}
}

// buildSettings constructs the API Settings from config.Settings.
func buildSettings(s config.Settings) uforpc.Settings {
	return uforpc.Settings{
		SchemaVersion:          s.Version,
		NotifyOnError:          s.NotifyOnError,
		NotifyOnStart:          s.NotifyOnStart,
		NotifyOnFinish:         s.NotifyOnFinish,
		SoundFeedbackEnable:    s.SoundFeedbackEnable,
		SoundFeedbackRecordId:  s.SoundFeedbackRecordID,
		SoundFeedbackSuccessId: s.SoundFeedbackSuccessID,
		SoundFeedbackErrorId:   s.SoundFeedbackErrorID,
		SoundFeedbackVolume:    s.SoundFeedbackVolume,
		InputDevice:            s.InputDevice,
		OutputMode:             string(s.OutputMode),
		OutputTrailingSpace:    s.OutputTrailingSpace,
		PostProcessEnabled:     s.PostProcessEnabled,
		PostProcessBaseUrl:     s.PostProcessBaseURL,
		PostProcessApiKey:      s.PostProcessAPIKey,
		PostProcessModel:       s.PostProcessModel,
		PostProcessPromptId:    s.PostProcessPromptID,
		Prompts:                mapPrompts(s.Prompts),
		HistoryLimit:           s.HistoryLimit,
		ModelUnloadEnable:      s.ModelUnloadEnable,
		ModelUnloadSeconds:     s.ModelUnloadSeconds,
		ShortcutToggle:         shortcutToAPI(s.ShortcutToggle),
		PasteShortcut:          s.PasteShortcut,
	}
}

// settingsFromAPI converts API Settings to config.Settings.
func settingsFromAPI(s uforpc.Settings) config.Settings {
	return config.Settings{
		Version:                s.SchemaVersion,
		NotifyOnError:          s.NotifyOnError,
		NotifyOnStart:          s.NotifyOnStart,
		NotifyOnFinish:         s.NotifyOnFinish,
		SoundFeedbackEnable:    s.SoundFeedbackEnable,
		SoundFeedbackRecordID:  s.SoundFeedbackRecordId,
		SoundFeedbackSuccessID: s.SoundFeedbackSuccessId,
		SoundFeedbackErrorID:   s.SoundFeedbackErrorId,
		SoundFeedbackVolume:    s.SoundFeedbackVolume,
		InputDevice:            s.InputDevice,
		OutputMode:             config.OutputMode(s.OutputMode),
		OutputTrailingSpace:    s.OutputTrailingSpace,
		PostProcessEnabled:     s.PostProcessEnabled,
		PostProcessBaseURL:     s.PostProcessBaseUrl,
		PostProcessAPIKey:      s.PostProcessApiKey,
		PostProcessModel:       s.PostProcessModel,
		PostProcessPromptID:    s.PostProcessPromptId,
		Prompts:                promptsFromAPI(s.Prompts),
		HistoryLimit:           s.HistoryLimit,
		ModelUnloadEnable:      s.ModelUnloadEnable,
		ModelUnloadSeconds:     s.ModelUnloadSeconds,
		ShortcutToggle:         shortcutFromAPI(s.ShortcutToggle),
		PasteShortcut:          s.PasteShortcut,
	}
}

// mapAudioDevices converts state.AudioDevice slice to API AudioDevice slice.
func mapAudioDevices(devices []state.AudioDevice) []uforpc.AudioDevice {
	result := make([]uforpc.AudioDevice, len(devices))
	for i, d := range devices {
		result[i] = uforpc.AudioDevice{
			Id:        d.ID,
			Name:      d.Name,
			IsDefault: d.IsDefault,
		}
	}
	return result
}

// mapHistoryEntries converts history.Entry slice to API HistoryEntry slice.
func mapHistoryEntries(entries []history.Entry) []uforpc.HistoryEntry {
	result := make([]uforpc.HistoryEntry, len(entries))
	for i, e := range entries {
		result[i] = uforpc.HistoryEntry{
			Id:            e.ID,
			Timestamp:     e.StartedAt.Format("2006-01-02T15:04:05Z07:00"),
			TextRaw:       e.TranscriptionRaw,
			TextFinal:     e.TranscriptionFinal,
			DurationMs:    int(e.RecordingDurationMs),
			PostProcessed: e.PostProcessed,
		}
	}
	return result
}

// mapPrompts converts config.Prompt slice to API Prompt slice.
func mapPrompts(prompts []config.Prompt) []uforpc.Prompt {
	result := make([]uforpc.Prompt, len(prompts))
	for i, p := range prompts {
		result[i] = uforpc.Prompt{
			Id:   p.ID,
			Name: p.Name,
			Body: p.Body,
		}
	}
	return result
}

// promptsFromAPI converts API Prompt slice to config.Prompt slice.
func promptsFromAPI(prompts []uforpc.Prompt) []config.Prompt {
	result := make([]config.Prompt, len(prompts))
	for i, p := range prompts {
		result[i] = config.Prompt{
			ID:   p.Id,
			Name: p.Name,
			Body: p.Body,
		}
	}
	return result
}

// shortcutToAPI converts config.Shortcut to API Shortcut.
func shortcutToAPI(s config.Shortcut) uforpc.Shortcut {
	return uforpc.Shortcut{
		Modifiers: s.Modifiers,
		Key:       s.Key,
	}
}

// shortcutFromAPI converts API Shortcut to config.Shortcut.
func shortcutFromAPI(s uforpc.Shortcut) config.Shortcut {
	return config.Shortcut{
		Modifiers: s.Modifiers,
		Key:       s.Key,
	}
}

// systemInfoToAPI converts state.SystemInfo to API SystemInfo.
func systemInfoToAPI(s state.SystemInfo) uforpc.SystemInfo {
	return uforpc.SystemInfo{
		Os:            s.OS,
		Arch:          s.Arch,
		DisplayServer: s.DisplayServer,
	}
}
