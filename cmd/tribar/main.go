package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/varavelio/tribar/internal/clipboard"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/engine"
	"github.com/varavelio/tribar/internal/logger"
	"github.com/varavelio/tribar/internal/notify"
	"github.com/varavelio/tribar/internal/onnx"
	"github.com/varavelio/tribar/internal/postprocess"
	"github.com/varavelio/tribar/internal/record"
	"github.com/varavelio/tribar/internal/sound"
	"github.com/varavelio/tribar/internal/state"
	"github.com/varavelio/tribar/internal/systray"
	"github.com/varavelio/tribar/internal/transcribe"
)

type cliFlags struct {
	Debug bool
}

func main() {
	flags := parseFlags()
	logger := logger.NewSlogLogger(flags.Debug)
	if err := run(logger); err != nil {
		logger.Error(context.Background(), "error while running the app", "err", err)
		os.Exit(1)
	}
}

func run(logger logger.Logger) error {
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

	if err := onnx.EnsureSharedLibrary(logger); err != nil {
		return fmt.Errorf("error ensuring ONNX Runtime shared library: %w", err)
	}

	settingsManager, err := config.NewSettingsManager()
	if err != nil {
		return fmt.Errorf("error loading settings: %w", err)
	}
	settings := settingsManager.Get()

	appState := state.New(settings.HistoryLimit)

	recorder, err := record.NewRecorder()
	if err != nil {
		return fmt.Errorf("error creating recorder: %w", err)
	}

	transcriber, err := transcribe.New()
	if err != nil {
		return fmt.Errorf("error creating transcriber: %w", err)
	}
	defer func() { _ = transcriber.Shutdown() }()

	notifier := notify.New(logger, notify.Settings{
		NotifyOnError:  settings.NotifyOnError,
		NotifyOnStart:  settings.NotifyOnStart,
		NotifyOnFinish: settings.NotifyOnFinish,
	})

	soundPlayer := sound.New(logger, sound.Settings{
		SoundOnStart:  settings.SoundOnStart,
		SoundOnFinish: settings.SoundOnFinish,
	})
	defer soundPlayer.Shutdown()

	cpb := clipboard.New(logger)

	postProcessor := postprocess.New(logger, settingsManager)

	eng := engine.New(engine.Dependencies{
		Logger:          logger,
		SettingsManager: settingsManager,
		State:           appState,
		Recorder:        recorder,
		Transcriber:     transcriber,
		PostProcess:     postProcessor,
		Writer:          cpb,
		Notifier:        notifier,
		Sound:           soundPlayer,
	})
	defer eng.Shutdown()

	go loadModelsAsync(ctx, logger, eng)

	stray := systray.New(appState, eng, stop)
	go stray.Start()
	defer stray.Shutdown()

	<-ctx.Done()
	stop()
	logger.Info(ctx, "shutting down gracefully...")
	return nil
}

func loadModelsAsync(ctx context.Context, logger logger.Logger, eng *engine.Engine) {
	progressCallback := func(filename string, downloaded, total int64, percent float64) {
		logger.Info(ctx, "downloading model",
			"file", filename,
			"progress", fmt.Sprintf("%.1f%%", percent),
		)
	}

	if err := eng.LoadModels(progressCallback); err != nil {
		logger.Error(ctx, "failed to load models", "err", err)
	}
}

func parseFlags() cliFlags {
	debugPtr := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	return cliFlags{
		Debug: *debugPtr,
	}
}
