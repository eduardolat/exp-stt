// Package executor implements the workflow execution engine.
package executor

import (
	"github.com/varavelio/tribar/internal/clipboard"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/history"
	"github.com/varavelio/tribar/internal/logger"
	"github.com/varavelio/tribar/internal/notify"
	"github.com/varavelio/tribar/internal/postprocess"
	"github.com/varavelio/tribar/internal/sound"
	"github.com/varavelio/tribar/internal/transcribe"
)

// ServiceContainer holds all shared services for node execution.
// Nodes receive this container instead of instantiating services.
// It implements the nodes.ServiceProvider interface.
type ServiceContainer struct {
	logger          logger.Logger
	settingsManager *config.SettingsManager
	transcriber     *transcribe.Instance
	postProcessor   *postprocess.Instance
	notifier        *notify.Instance
	sound           *sound.Instance
	clipboard       *clipboard.Instance
	historyManager  *history.Manager
}

// NewServiceContainer creates a new service container with all dependencies.
func NewServiceContainer(
	log logger.Logger,
	settingsManager *config.SettingsManager,
	transcriber *transcribe.Instance,
	postProcessor *postprocess.Instance,
	notifier *notify.Instance,
	snd *sound.Instance,
	cpb *clipboard.Instance,
	historyManager *history.Manager,
) *ServiceContainer {
	return &ServiceContainer{
		logger:          log,
		settingsManager: settingsManager,
		transcriber:     transcriber,
		postProcessor:   postProcessor,
		notifier:        notifier,
		sound:           snd,
		clipboard:       cpb,
		historyManager:  historyManager,
	}
}

// GetLogger returns the logger instance.
func (s *ServiceContainer) GetLogger() logger.Logger {
	return s.logger
}

// GetSettingsManager returns the settings manager.
func (s *ServiceContainer) GetSettingsManager() *config.SettingsManager {
	return s.settingsManager
}

// GetTranscriber returns the transcriber instance.
func (s *ServiceContainer) GetTranscriber() *transcribe.Instance {
	return s.transcriber
}

// GetPostProcessor returns the post-processor instance.
func (s *ServiceContainer) GetPostProcessor() *postprocess.Instance {
	return s.postProcessor
}

// GetNotifier returns the notifier instance.
func (s *ServiceContainer) GetNotifier() *notify.Instance {
	return s.notifier
}

// GetSound returns the sound player instance.
func (s *ServiceContainer) GetSound() *sound.Instance {
	return s.sound
}

// GetClipboard returns the clipboard instance.
func (s *ServiceContainer) GetClipboard() *clipboard.Instance {
	return s.clipboard
}

// GetHistoryManager returns the history manager.
func (s *ServiceContainer) GetHistoryManager() *history.Manager {
	return s.historyManager
}
