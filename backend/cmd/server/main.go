package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/app"
	"github.com/Wei-Shaw/sub2api/internal/config"
)

func main() {
	migrateOnly := flag.Bool("migrate-only", false, "run database migration and seed data, then exit")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("initialize app: %v%s", err, databaseHint(cfg))
	}
	defer application.Close()

	if *migrateOnly {
		if err := application.Bootstrap(context.Background()); err != nil {
			log.Fatalf("migrate starter data: %v%s", err, databaseHint(cfg))
		}
		log.Println("database migration and seed completed")
		return
	}

	if err := application.Bootstrap(context.Background()); err != nil {
		log.Fatalf("bootstrap app: %v%s", err, databaseHint(cfg))
	}

	go func() {
		if err := application.Start(); err != nil {
			log.Fatalf("start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown server: %v", err)
	}
}

func databaseHint(cfg *config.Config) string {
	return fmt.Sprintf(
		"\n\nDatabase target: %s:%d/%s\nIf PostgreSQL is not running yet, start it first. Example:\n  docker compose up -d postgres redis\nOr run a local PostgreSQL instance and set DATABASE_HOST / DATABASE_PORT / DATABASE_USER / DATABASE_PASSWORD / DATABASE_NAME.",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
}
