package main

import (
	"log"
	"net/http"
	"time"

	"github.com/tmcnulty387/LaundryStatus/backend/internal/config"
	repo "github.com/tmcnulty387/LaundryStatus/backend/internal/repository/sqlc"
	"github.com/tmcnulty387/LaundryStatus/backend/internal/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	config  *config.Config
	queries *repo.Queries
	pool    *pgxpool.Pool
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK\n"))
	})

	// GET /api/rooms/{room_slug}/machines
	// PUT /api/rooms/{room_slug}/machines/{machine_id}/reserve
	// POST /api/rooms/{room_slug}/machines/{machine_id}/available
	// POST /api/rooms/{room_slug}/machines/{machine_id}/out-of-order

	s := routes.NewService(app.queries, app.pool, app.config)
	h := routes.NewHandler(s)
	r.Route("/api/rooms/{room_slug}/machines", func(r chi.Router) {
		r.Get("/", h.GetMachines)
		r.Put("/{machine_id}/reserve", h.CreateReservation)
		r.Post("/{machine_id}/available", h.SetMachineAvailable)
		r.Post("/{machine_id}/out-of-order", h.SetMachineOutOfOrder)
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.Addr(),
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at %s", app.config.Addr())

	return srv.ListenAndServe()
}
