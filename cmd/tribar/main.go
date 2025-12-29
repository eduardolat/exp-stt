package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/varavelio/tribar/internal/clipboard"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/engine"
	"github.com/varavelio/tribar/internal/history"
	"github.com/varavelio/tribar/internal/instance"
	"github.com/varavelio/tribar/internal/logger"
	"github.com/varavelio/tribar/internal/notify"
	"github.com/varavelio/tribar/internal/onnx"
	"github.com/varavelio/tribar/internal/postprocess"
	"github.com/varavelio/tribar/internal/record"
	"github.com/varavelio/tribar/internal/server"
	"github.com/varavelio/tribar/internal/sound"
	"github.com/varavelio/tribar/internal/state"
	"github.com/varavelio/tribar/internal/systray"
	"github.com/varavelio/tribar/internal/toggle"
	"github.com/varavelio/tribar/internal/transcribe"
)

func main() {
	flags := config.ParseFlags()
	logger := logger.NewSlogLogger(flags.Debug)

	if err := run(logger, flags); err != nil {
		logger.Error(context.Background(), "error while running the app", "err", err)
		os.Exit(1)
	}
}

func run(logger logger.Logger, flags config.Flags) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Debug(
		ctx, "operating system info",
		"os", runtime.GOOS,
		"arch", runtime.GOARCH,
	)

	if err := config.EnsureDirectories(logger); err != nil {
		return fmt.Errorf("error ensuring app directories: %w", err)
	}

	if flags.Toggle {
		if err := toggle.Execute(); err != nil {
			return err
		}
		fmt.Println("Recording toggled successfully.")
		return nil
	}

	inst := instance.NewManager()
	if err := inst.AcquireLock(); err != nil {
		if errors.Is(err, instance.ErrAlreadyRunning) {
			fmt.Println("Tribar is already running.")
			return nil
		}
		return err
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

	settingsManager, err := config.NewSettingsManager()
	if err != nil {
		return fmt.Errorf("error loading settings: %w", err)
	}

	historyManager := history.NewManager(logger, settingsManager)
	historyManager.LoadAsync(ctx)

	appState := state.New(historyManager)

	recorder, err := record.NewRecorder(settingsManager)
	if err != nil {
		return fmt.Errorf("error creating recorder: %w", err)
	}

	transcriber, err := transcribe.New()
	if err != nil {
		return fmt.Errorf("error creating transcriber: %w", err)
	}
	defer func() { _ = transcriber.Shutdown() }()

	notifier := notify.New(logger, settingsManager)

	soundPlayer := sound.New(logger, settingsManager)
	defer soundPlayer.Shutdown()

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
	defer eng.Shutdown()

	go func() {
		if err := eng.EnsureModelsDownloaded(); err != nil {
			logger.Error(ctx, "failed to download models", "err", err)
		}
	}()

	srv := server.NewServer(logger, settingsManager, appState, eng)
	go func() {
		addr := fmt.Sprintf("%s:%d", flags.Host, inst.Port())
		logger.Info(ctx, "Starting server", "address", "http://"+addr+"/", "version", config.AppVersion)
		if err := srv.Server.Serve(listener); err != nil && err != http.ErrServerClosed {
			logger.Error(ctx, "server error", "err", err)
			stop()
		}
	}()

	stray := systray.New(appState, eng, stop)
	go stray.Start()
	defer stray.Shutdown()

	<-ctx.Done()
	stop()
	logger.Info(ctx, "shutting down gracefully...")
	return nil
}
