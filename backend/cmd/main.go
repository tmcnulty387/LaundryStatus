package main

import (
	"context"
	"log"

	"github.com/tmcnulty387/LaundryStatus/backend/internal/config"
	repo "github.com/tmcnulty387/LaundryStatus/backend/internal/repository/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file into environment (only useful when not using docker-compose)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to database
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
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
