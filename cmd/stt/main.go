package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/eduardolat/exp-stt/internal/app"
	"github.com/eduardolat/exp-stt/internal/config"
	"github.com/eduardolat/exp-stt/internal/logger"
	"github.com/eduardolat/exp-stt/internal/onnx"
	"github.com/eduardolat/exp-stt/internal/systray"
)

func main() {
	logger := logger.NewSlogLogger(false)
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

	a := app.New()

	st := systray.New(a)
	go st.Start()
	defer st.Shutdown()

	<-ctx.Done()
	stop()
	return nil
}
