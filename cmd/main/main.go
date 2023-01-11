package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"zip_service/internal/app"
)

// @title   Zip service API documentation
// @version 1.0.0

// @BasePath /

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	a := app.NewApp()
	err := a.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
