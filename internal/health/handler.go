package health

import (
	"net/http"

	"github.com/username/go-api-kit/internal/response"
)

// Handler holds dependencies for health endpoints.
type Handler struct{}

// NewHandler creates a new Health Handler.
func NewHandler() *Handler {
	return &Handler{}
}

// RegisterRoutes registers health endpoints on the provided ServeMux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /healthz", h.Health)
}

// Health returns HTTP 200 OK with status: ok.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{
		"status": "ok",
	})
}
