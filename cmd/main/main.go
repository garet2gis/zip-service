package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"zip_service/internal/app"
	"zip_service/internal/config"
)

// @title   Zip service API documentation
// @version 1.0.0

// @BasePath /

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Println("config initializing")
	cfg := config.GetConfig()
	log.Printf("Config: %+v", cfg)

	a := app.NewApp(cfg)
	err := a.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
