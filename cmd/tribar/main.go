package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/varavelio/tribar/internal/clipboard"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/engine"
	"github.com/varavelio/tribar/internal/eventbus"
	"github.com/varavelio/tribar/internal/history"
	"github.com/varavelio/tribar/internal/instance"
	"github.com/varavelio/tribar/internal/logger"
	"github.com/varavelio/tribar/internal/notify"
	"github.com/varavelio/tribar/internal/onnx"
	"github.com/varavelio/tribar/internal/postprocess"
	"github.com/varavelio/tribar/internal/record"
	"github.com/varavelio/tribar/internal/server"
	"github.com/varavelio/tribar/internal/shortcut"
	"github.com/varavelio/tribar/internal/sound"
	"github.com/varavelio/tribar/internal/state"
	"github.com/varavelio/tribar/internal/systray"
	"github.com/varavelio/tribar/internal/toggle"
	"github.com/varavelio/tribar/internal/transcribe"
	"golang.design/x/hotkey/mainthread"
)

// singleShutdownTimeout is the maximum time to wait for graceful
// shutdown of a single component. After this duration, the process
// exits forcefully to prevent system freezes.
const singleShutdownTimeout = 3 * time.Second

func main() {
	flags := config.ParseFlags()
	log := logger.NewSlogLogger(flags.Debug)

	mainthread.Init(func() {
		if err := run(log, flags); err != nil {
			log.Error(context.Background(), "error while running the app", "err", err)
			os.Exit(1)
		}
	})
}

func run(logger logger.Logger, flags config.Flags) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	state.InitRuntime()
	logger.Debug(
		ctx, "runtime info initialized",
		"os", state.RuntimeInfo.OS,
		"arch", state.RuntimeInfo.Arch,
		"display_server", state.RuntimeInfo.DisplayServer,
	)

	if err := config.EnsureDirectories(logger); err != nil {
		return fmt.Errorf("error ensuring app directories: %w", err)
	}

	if flags.Toggle {
		if err := toggle.Execute(); err != nil {
			return fmt.Errorf("error toggling recording: %w", err)
		}
		fmt.Println("recording toggled successfully")
		return nil
	}

	inst := instance.NewManager()
	if err := inst.AcquireLock(); err != nil {
		if errors.Is(err, instance.ErrAlreadyRunning) {
			fmt.Println("tribar is already running")
			return nil
		}
		return fmt.Errorf("error acquiring lock: %w", err)
	}
	defer inst.Cleanup()

	listener, err := inst.CreateListener(flags.Host, flags.Port, flags.PortExplicit)
	if err != nil {
		return fmt.Errorf("error binding server port: %w", err)
	}

	if err := inst.WritePortFile(); err != nil {
		return fmt.Errorf("error writing port file: %w", err)
	}

	if err := onnx.EnsureSharedLibrary(logger); err != nil {
		return fmt.Errorf("error ensuring ONNX Runtime shared library: %w", err)
	}

	eventBus := eventbus.New()

	settingsManager, err := config.NewSettingsManager(eventBus)
	if err != nil {
		return fmt.Errorf("error loading settings: %w", err)
	}

	historyManager := history.NewManager(logger, settingsManager)
	historyManager.LoadAsync(ctx)

	appState := state.New(eventBus, historyManager)

	recorder, err := record.NewRecorder(logger, settingsManager)
	if err != nil {
		return fmt.Errorf("error creating recorder: %w", err)
	}
	defer shutdownWithTimeout(recorder.Shutdown)

	transcriber, err := transcribe.New(logger)
	if err != nil {
		return fmt.Errorf("error creating transcriber: %w", err)
	}
	defer shutdownWithTimeout(transcriber.Shutdown)

	notifier := notify.New(logger, settingsManager)

	soundPlayer := sound.New(logger, settingsManager)
	defer shutdownWithTimeout(soundPlayer.Shutdown)

	cpb := clipboard.New(logger)

	postProcessor := postprocess.New(logger, settingsManager)

	eng := engine.New(engine.Dependencies{
		Logger:          logger,
		SettingsManager: settingsManager,
		HistoryManager:  historyManager,
		State:           appState,
		Recorder:        recorder,
		Transcriber:     transcriber,
		PostProcess:     postProcessor,
		Writer:          cpb,
		Notifier:        notifier,
		Sound:           soundPlayer,
	})
	defer shutdownWithTimeout(eng.Shutdown)

	go func() {
		if err := eng.EnsureModelsDownloaded(); err != nil {
			logger.Error(ctx, "failed to download models", "err", err)
		}
	}()

	shortcutMgr := shortcut.New(logger, settingsManager, eng.ToggleRecording)
	if err := shortcutMgr.Start(); err != nil {
		logger.Warn(ctx, "hotkey registration failed", "err", err)
	}
	defer shutdownWithTimeout(shortcutMgr.Stop)

	srv := server.NewServer(logger, settingsManager, appState, eng, shortcutMgr, eventBus)
	go func() {
		addr := fmt.Sprintf("%s:%d", flags.Host, inst.Port())
		logger.Info(ctx, "Starting server", "address", "http://"+addr+"/", "version", config.AppVersion)
		if err := srv.Server.Serve(listener); err != nil && err != http.ErrServerClosed {
			logger.Error(ctx, "server error", "err", err)
			stop()
		}
	}()

	stray := systray.New(appState, eng, inst.Port(), stop)
	go stray.Start()
	defer shutdownWithTimeout(stray.Shutdown)

	<-ctx.Done()
	stop()
	logger.Info(ctx, "shutting down gracefully...")
	return nil
}

// shutdownWithTimeout runs fn in a goroutine and waits up to shutdownTimeout for it
// to complete. This prevents any single cleanup operation from blocking the
// entire shutdown process indefinitely.
func shutdownWithTimeout(fn func()) {
	done := make(chan any)
	go func() {
		defer close(done)
		fn()
	}()

	select {
	case <-done:
	case <-time.After(singleShutdownTimeout):
	}
}
