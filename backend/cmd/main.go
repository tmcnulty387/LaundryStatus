package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/tmcnulty387/LaundryStatus/backend/internal/config"
	repo "github.com/tmcnulty387/LaundryStatus/backend/internal/repository/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	// log working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	log.Println("Working directory:", wd)

	// Load .env file into environment (only useful when not using docker-compose)
	if err := godotenv.Load(".env", "../.env"); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to database
	timeout := time.AfterFunc(10*time.Second, func() {
		log.Fatal("Failed to connect to database: timed out after 10 seconds")
	})
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	timeout.Stop()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	timeout = time.AfterFunc(10*time.Second, func() {
		log.Fatal("Failed to ping database: timed out after 10 seconds")
	})
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	timeout.Stop()

	defer pool.Close()
	log.Println("Connected to database at", cfg.DatabaseURL)

	queries := repo.New(pool)

	api := &application{
		config:  cfg,
		queries: queries,
		pool:    pool,
	}

	// TODO: Start timers for any existing reservations

	if err := api.run(api.mount()); err != nil {
		log.Fatal("Server has failed to start: ", err)
	}
}
