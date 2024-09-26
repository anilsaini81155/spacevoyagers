package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs the details of incoming requests and responses
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now() // Record the start time

		log.Printf("Started %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		duration := time.Since(startTime) // Calculate the time taken
		log.Printf("Completed %s %s in %v", r.Method, r.RequestURI, duration)

	})
}
