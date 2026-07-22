package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/username/go-api-kit/config"
	"github.com/username/go-api-kit/internal/docs"
	"github.com/username/go-api-kit/internal/health"
	"github.com/username/go-api-kit/internal/platform/database"
	"github.com/username/go-api-kit/internal/platform/http/middleware"
	"github.com/username/go-api-kit/internal/platform/logger"
)

func main() {
	startTime := time.Now()

	// 1. Load configuration from environment variables
	cfg := config.Load()

	// 2. Initialize structured logger
	log := logger.New(cfg.Log.Level, cfg.Log.Format)

	if err := cfg.Validate(); err != nil {
		log.Error("configuration validation error", "error", err)
		os.Exit(1)
	}

	log.Info("starting application")
	log.Info("app info", "name", cfg.App.Name)
	log.Info("app info", "env", cfg.App.Env)
	log.Info("app info", "version", cfg.App.Version)

	ctx := context.Background()

	// 3. Connect to PostgreSQL
	pool, err := database.Connect(ctx, cfg.DB.DSN())
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	log.Info("database connected", "host", cfg.DB.Host, "name", cfg.DB.Name)

	// 4. Run database migrations (goose)
	if err := database.Migrate(ctx, cfg.DB.DSN()); err != nil {
		log.Error("migrations failed", "error", err)
		os.Exit(1)
	}
	log.Info("migrations applied successfully")

	// 5. Wire dependencies & Initialize Handlers
	healthHandler := health.NewHandler()
	docsHandler := docs.NewHandler()

	// 6. Register routes
	mux := http.NewServeMux()
	healthHandler.RegisterRoutes(mux)
	docsHandler.RegisterRoutes(mux)

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
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	startupDuration := time.Since(startTime)
	log.Info("server listening", "addr", addr)
	log.Info("application ready", "startup_duration", startupDuration.String())

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
