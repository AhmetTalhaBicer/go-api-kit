package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/username/go-api-kit/internal/config"
	"github.com/username/go-api-kit/internal/domain1"
	"github.com/username/go-api-kit/internal/domain2"
	"github.com/username/go-api-kit/internal/health"
	"github.com/username/go-api-kit/internal/logger"
	"github.com/username/go-api-kit/internal/middleware"
)

// --- Domain1 in-memory stub ---

type domain1Stub struct{}

func (s *domain1Stub) FindAll(_ context.Context) ([]domain1.Domain1, error) {
	return []domain1.Domain1{{ID: 1, Name: "Stub 1", CreatedAt: time.Now()}}, nil
}
func (s *domain1Stub) FindByID(_ context.Context, id int64) (*domain1.Domain1, error) {
	return &domain1.Domain1{ID: id, Name: "Stub 1", CreatedAt: time.Now()}, nil
}
func (s *domain1Stub) Create(_ context.Context, _ *domain1.Domain1) error { return nil }
func (s *domain1Stub) Update(_ context.Context, _ *domain1.Domain1) error { return nil }
func (s *domain1Stub) Delete(_ context.Context, _ int64) error            { return nil }

// --- Domain2 in-memory stub ---

type domain2Stub struct{}

func (s *domain2Stub) FindAll(_ context.Context) ([]domain2.Domain2, error) {
	return []domain2.Domain2{{ID: 1, Name: "Stub 2", CreatedAt: time.Now()}}, nil
}
func (s *domain2Stub) FindByID(_ context.Context, id int64) (*domain2.Domain2, error) {
	return &domain2.Domain2{ID: id, Name: "Stub 2", CreatedAt: time.Now()}, nil
}
func (s *domain2Stub) Create(_ context.Context, _ *domain2.Domain2) error { return nil }
func (s *domain2Stub) Update(_ context.Context, _ *domain2.Domain2) error { return nil }
func (s *domain2Stub) Delete(_ context.Context, _ int64) error            { return nil }

// --- Test helpers ---

func newTestApp(t *testing.T) http.Handler {
	t.Helper()

	cfg := config.Load()
	log := logger.New(cfg.Log.Level, cfg.Log.Format)

	mux := http.NewServeMux()
	health.NewHandler().RegisterRoutes(mux)
	domain1.NewHandler(log, domain1.NewService(&domain1Stub{})).RegisterRoutes(mux)
	domain2.NewHandler(log, domain2.NewService(&domain2Stub{})).RegisterRoutes(mux)

	return middleware.Chain(mux,
		middleware.RequestID,
		middleware.Timeout(5*time.Second),
		middleware.Logger(log),
		middleware.Recover(log),
		middleware.CORS,
	)
}

func TestHealthEndpoint(t *testing.T) {
	app := newTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	reqID := rec.Header().Get("X-Request-ID")
	if reqID == "" {
		t.Error("expected X-Request-ID header to be set")
	}
}

func TestDomain1List(t *testing.T) {
	app := newTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/domain1", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestDomain2List(t *testing.T) {
	app := newTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/domain2", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
