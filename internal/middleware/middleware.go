package middleware

import (
	"net/http"
	"strings"

	"github.com/DiegoSenaa/go-rater-limiter/internal/ratelimiter"
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := strings.Split(r.RemoteAddr, ":")[0]
		apiKey := r.Header.Get("API_KEY")

		if apiKey != "" {
			if !ratelimiter.AllowRequest(apiKey, "token") {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
		} else {
			if !ratelimiter.AllowRequest(clientIP, "ip") {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
