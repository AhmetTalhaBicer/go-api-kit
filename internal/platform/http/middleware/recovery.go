package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/username/go-api-kit/internal/platform/http/response"
)

// Recover returns a Middleware that catches panics and returns a 500 Internal Server Error.
func Recover(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error("panic recovered", "error", fmt.Sprintf("%v", err), "path", r.URL.Path)
					response.Error(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error", nil)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
