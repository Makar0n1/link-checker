package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/link-tracker/index-service/internal/config"
	"github.com/link-tracker/index-service/internal/handler"
	"github.com/link-tracker/index-service/internal/repository"
	"github.com/link-tracker/index-service/internal/service"
	"github.com/link-tracker/shared/pkg/middleware"
)

func main() {
	cfg := config.Load()

	// Database connection
	dbPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Verify connection
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	log.Println("Connected to database")

	// Initialize layers
	platformRepo := repository.NewPlatformRepository(dbPool)
	platformService := service.NewPlatformService(platformRepo)
	platformHandler := handler.NewPlatformHandler(platformService)
	healthHandler := handler.NewHealthHandler(dbPool)

	// JWT middleware config
	jwtConfig := middleware.JWTConfig{
		Secret:     cfg.JWTSecret,
		ContextKey: "claims",
	}

	// Router setup
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health endpoints (public)
	r.Get("/health", healthHandler.Health)
	r.Get("/ready", healthHandler.Ready)

	// API routes (protected)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.JWTAuth(jwtConfig))

		r.Route("/platforms", func(r chi.Router) {
			r.Get("/", platformHandler.List)
			r.Post("/", platformHandler.Create)
			r.Post("/bulk", platformHandler.BulkCreate)
			r.Get("/{id}", platformHandler.GetByID)
			r.Put("/{id}", platformHandler.Update)
			r.Delete("/{id}", platformHandler.Delete)
			r.Post("/{id}/check", platformHandler.CheckIndex)
		})
	})

	// Server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Index service starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
