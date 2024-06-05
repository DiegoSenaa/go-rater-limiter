package main

import (
	"log"
	"net/http"

	"github.com/DiegoSenaa/go-rater-limiter/internal/middleware"
	"github.com/DiegoSenaa/go-rater-limiter/internal/redisclient"
	"github.com/go-chi/chi/v5"
)

func main() {
	redisclient.InitRedisClient()

	r := chi.NewRouter()
	r.Use(middleware.RateLimitMiddleware)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
