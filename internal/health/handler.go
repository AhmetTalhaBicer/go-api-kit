package health

import (
	"net/http"

	"github.com/username/go-api-kit/internal/platform/http/response"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/health", h.HealthCheck)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, "go-api-kit is running", map[string]string{
		"status": "ok",
	})
}
