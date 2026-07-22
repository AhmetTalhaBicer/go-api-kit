package docs

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "github.com/username/go-api-kit/api" // Import swagger docs
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("/docs/doc.json"),
	)

	mux.Handle("/docs/", swaggerHandler)
	mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/index.html", http.StatusMovedPermanently)
	})
}
