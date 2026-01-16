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
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/link-tracker/backlink-service/internal/config"
	"github.com/link-tracker/backlink-service/internal/handler"
	"github.com/link-tracker/backlink-service/internal/repository"
	"github.com/link-tracker/backlink-service/internal/service"
	"github.com/link-tracker/shared/pkg/middleware"
)

func main() {
	cfg := config.Load()

	// Database connection
	dbPool, err := pgxpool.New(context.Background(), cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Verify database connection
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database")

	// Initialize repositories
	projectRepo := repository.NewProjectRepository(dbPool)
	backlinkRepo := repository.NewBacklinkRepository(dbPool)

	// Initialize services
	projectService := service.NewProjectService(projectRepo)
	backlinkService := service.NewBacklinkService(backlinkRepo, projectRepo)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	projectHandler := handler.NewProjectHandler(projectService)
	backlinkHandler := handler.NewBacklinkHandler(backlinkService)

	// JWT middleware config
	jwtMiddleware := middleware.JWTAuth(middleware.JWTConfig{
		Secret: cfg.JWT.Secret,
	})

	// Setup router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Timeout(30 * time.Second))
	r.Use(corsMiddleware)

	// Health endpoints (public)
	r.Get("/health", healthHandler.Health)
	r.Get("/ready", healthHandler.Ready)

	// API routes (protected)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(jwtMiddleware)

		// Projects
		r.Route("/projects", func(r chi.Router) {
			r.Get("/", projectHandler.List)
			r.Post("/", projectHandler.Create)
			r.Get("/{id}", projectHandler.Get)
			r.Put("/{id}", projectHandler.Update)
			r.Delete("/{id}", projectHandler.Delete)
		})

		// Backlinks
		r.Route("/backlinks", func(r chi.Router) {
			r.Get("/", backlinkHandler.List)
			r.Post("/", backlinkHandler.Create)
			r.Post("/bulk", backlinkHandler.BulkCreate)
			r.Delete("/bulk", backlinkHandler.BulkDelete)
			r.Post("/import", backlinkHandler.Import)
			r.Get("/{id}", backlinkHandler.Get)
			r.Put("/{id}", backlinkHandler.Update)
			r.Delete("/{id}", backlinkHandler.Delete)
		})
	})

	// Server setup
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Starting backlink-service on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Request-ID")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
