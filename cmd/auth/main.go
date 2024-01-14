package main

import (
	"context"
	"github.com/patyukin/bs-auth/internal/config"
	"log"
	"log/slog"
	"os"

	"github.com/patyukin/bs-auth/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// set config init
	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("failed to init config: %s", err.Error())
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// logger init

	a, err := app.NewApp(ctx, cfg, logger)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
