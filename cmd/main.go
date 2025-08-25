package main

import (
	"context"
	"crypto_analyzer-api_gateway/internal/app"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Start(ctx); err != nil {
		log.Fatalf("start app error: %v", err)
	}
}
