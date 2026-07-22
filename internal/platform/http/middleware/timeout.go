package middleware

import (
	"net/http"
	"time"
)

// Timeout returns a Middleware that bounds request execution duration.
func Timeout(d time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, d, `{"error":"request timeout"}`)
	}
}
