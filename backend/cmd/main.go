package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/pressly/goose/v3"
	"github.com/tmcnulty387/LaundryStatus/backend/internal/config"
	repo "github.com/tmcnulty387/LaundryStatus/backend/internal/repository/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func gooseUp(cfg *config.Config, ctx context.Context) {
	if !cfg.RunMigrations {
		log.Println("Skipping goose migrations")
		return
	}

	gooseDB, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to open goose connection: %v", err)
	}
	defer gooseDB.Close()

	if err := gooseDB.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping goose connection: %v", err)
	}

	if err := goose.SetDialect(cfg.GooseDriver); err != nil {
		log.Fatalf("Failed to set goose dialect: %v", err)
	}

	if err := goose.Up(gooseDB, cfg.GooseMigrationDir); err != nil {
		log.Fatalf("Failed to run goose up: %v", err)
	}
	log.Printf("Goose migrations applied successfully from directory: %s", cfg.GooseMigrationDir)
}

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

	log.Printf("Running migrations...")
	gooseUp(cfg, ctx)

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
