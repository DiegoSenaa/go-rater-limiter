package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/DiegoSenaa/go-rater-limiter/internal/middleware"
	"github.com/DiegoSenaa/go-rater-limiter/internal/ratelimiter"
	"github.com/DiegoSenaa/go-rater-limiter/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	rateLimitIP, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_IP"))
	rateLimitToken, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_TOKEN"))

	redisStorage := storage.NewRedisStorage(redisAddr, redisPassword)
	rateLimiter := ratelimiter.NewRateLimiter(redisStorage, rateLimitIP, rateLimitToken)

	r := chi.NewRouter()
	r.Use(middleware.RateLimitMiddleware(rateLimiter))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})

	log.Fatal(http.ListenAndServe(":8080", r))

	log.Fatal(http.ListenAndServe(":8080", r))
}
