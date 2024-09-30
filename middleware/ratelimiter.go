package middleware

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

// RateLimiterMiddleware applies rate limiting for incoming HTTP requests.
func RateLimiterMiddleware(limiter *rate.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If the request exceeds the rate limit, return a 429 error
			if !limiter.Allow() {

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)

				// Custom message with retry information
				message := fmt.Sprintf(`{
					"error": "Request limit exceeded",
					"message": "You have exceeded the request limit. Please wait 60 seconds before trying again."
				}`)

				w.Write([]byte(message))

				// http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			// Call the next handler if the rate limit is not exceeded
			next.ServeHTTP(w, r)
		})
	}
}
