package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrxacker/go-to-do-app/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to create application: %v", err)
	}

	if err := a.Start(ctx); err != nil {
		log.Fatalf("failed to start application: %v", err)
	}
}
