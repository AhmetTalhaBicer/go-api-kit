package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/username/go-api-kit/internal/config"
	"github.com/username/go-api-kit/internal/database"
	"github.com/username/go-api-kit/internal/domain1"
	"github.com/username/go-api-kit/internal/domain2"
	"github.com/username/go-api-kit/internal/health"
	"github.com/username/go-api-kit/internal/logger"
	"github.com/username/go-api-kit/internal/middleware"
)

// TODO: Replace these stubs with real repository implementations.
// Example after running `make sqlc`:
//
//	queries  := sqlcdb.New(pool)
//	d1Repo   := domain1postgres.New(queries)
//	d2Repo   := domain2postgres.New(queries)

type domain1Stub struct{}

func (s *domain1Stub) FindAll(ctx context.Context) ([]domain1.Domain1, error) {
	return []domain1.Domain1{}, nil
}
func (s *domain1Stub) FindByID(ctx context.Context, id int64) (*domain1.Domain1, error) {
	return &domain1.Domain1{ID: id, Name: "Stub 1", CreatedAt: time.Now()}, nil
}
func (s *domain1Stub) Create(ctx context.Context, d *domain1.Domain1) error { return nil }
func (s *domain1Stub) Update(ctx context.Context, d *domain1.Domain1) error { return nil }
func (s *domain1Stub) Delete(ctx context.Context, id int64) error            { return nil }

type domain2Stub struct{}

func (s *domain2Stub) FindAll(ctx context.Context) ([]domain2.Domain2, error) {
	return []domain2.Domain2{}, nil
}
func (s *domain2Stub) FindByID(ctx context.Context, id int64) (*domain2.Domain2, error) {
	return &domain2.Domain2{ID: id, Name: "Stub 2", CreatedAt: time.Now()}, nil
}
func (s *domain2Stub) Create(ctx context.Context, d *domain2.Domain2) error { return nil }
func (s *domain2Stub) Update(ctx context.Context, d *domain2.Domain2) error { return nil }
func (s *domain2Stub) Delete(ctx context.Context, id int64) error            { return nil }

func main() {
	// 1. Load configuration from environment variables
	cfg := config.Load()

	// 2. Initialize structured logger
	log := logger.New(cfg.Log.Level, cfg.Log.Format)
	log.Info("starting application",
		"name", cfg.App.Name,
		"env", cfg.App.Env,
		"version", cfg.App.Version,
	)

	ctx := context.Background()

	// 3. Connect to PostgreSQL
	//
	// pool, err := database.Connect(ctx, cfg.DB.DSN())
	// if err != nil {
	// 	log.Error("failed to connect to database", "error", err)
	// 	os.Exit(1)
	// }
	// defer pool.Close()
	// log.Info("database connected", "host", cfg.DB.Host, "name", cfg.DB.Name)

	// 4. Run database migrations (goose — embedded SQL files)
	//
	// if err := database.Migrate(ctx, cfg.DB.DSN()); err != nil {
	// 	log.Error("migration failed", "error", err)
	// 	os.Exit(1)
	// }
	// log.Info("migrations applied")

	// Suppress unused import warnings while the DB block is commented out.
	_ = database.Connect
	_ = ctx

	// 5. Wire dependencies for handlers
	d1Handler := domain1.NewHandler(log, domain1.NewService(&domain1Stub{}))
	d2Handler := domain2.NewHandler(log, domain2.NewService(&domain2Stub{}))
	healthHandler := health.NewHandler()

	// 6. Register routes
	mux := http.NewServeMux()
	healthHandler.RegisterRoutes(mux)
	d1Handler.RegisterRoutes(mux)
	d2Handler.RegisterRoutes(mux)

	// 7. Apply middleware chain (RequestID -> Timeout -> Logger -> Recover -> CORS)
	chain := middleware.Chain(mux,
		middleware.RequestID,
		middleware.Timeout(cfg.Server.ReadTimeout),
		middleware.Logger(log),
		middleware.Recover(log),
		middleware.CORS,
	)

	// 8. Configure the HTTP server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      chain,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 9. Listen for OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("server listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	<-quit
	log.Info("shutdown signal received")

	// Graceful shutdown: allow up to 30s for active connections to finish
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	log.Info("server shut down gracefully")
}
