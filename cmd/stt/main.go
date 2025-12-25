package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/eduardolat/exp-stt/internal/systray"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("error while running the app: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	st := systray.New()
	go st.Start()
	defer st.Shutdown()

	<-ctx.Done()
	stop()
	return nil
}
