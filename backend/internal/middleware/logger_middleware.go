package middleware

import (
	"log"
	"net/http"
	"time"
)

// RequestLoggerMiddleware logs the method and path of each incoming request and
// the time taken to complete it.
func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[REQUEST] %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("[DONE]    %s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}
