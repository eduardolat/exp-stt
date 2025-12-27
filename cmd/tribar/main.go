package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/varavelio/tribar/internal/app"
	"github.com/varavelio/tribar/internal/config"
	"github.com/varavelio/tribar/internal/logger"
	"github.com/varavelio/tribar/internal/onnx"
	"github.com/varavelio/tribar/internal/systray"
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

	a := app.New()

	st := systray.New(a)
	go st.Start()
	defer st.Shutdown()

	<-ctx.Done()
	stop()
	return nil
}

func parseFlags() cliFlags {
	debugPtr := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	return cliFlags{
		Debug: *debugPtr,
	}
}
