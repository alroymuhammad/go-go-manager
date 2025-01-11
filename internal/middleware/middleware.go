package middleware

import (
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logging request method and URL
		log.Printf("%s %s", r.Method, r.URL.Path)

		// Set Content-Type header for response
		w.Header().Set("Content-Type", "application/json")

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
