package ratelimiter

import (
	"fmt"
	"strconv"

	"github.com/DiegoSenaa/go-rater-limiter/internal/storage"
	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	Storage        storage.Storage
	RateLimitIP    int
	RateLimitToken int
}

func NewRateLimiter(storage storage.Storage, rateLimitIP, rateLimitToken int) *RateLimiter {
	return &RateLimiter{
		Storage:        storage,
		RateLimitIP:    rateLimitIP,
		RateLimitToken: rateLimitToken,
	}
}

func (r *RateLimiter) AllowRequest(identifier, idType string) bool {
	var limit int
	if idType == "token" {
		limit = r.RateLimitToken
	} else {
		limit = r.RateLimitIP
	}

	key := fmt.Sprintf("rate_limiter:%s:%s", idType, identifier)
	val, err := r.Storage.Get(key)

	if err != nil && err != redis.Nil {
		fmt.Println("Error getting value from storage:", err)
		return false
	}

	if err == redis.Nil {
		fmt.Printf("Setting initial value for key: %s\n", key)
		err = r.Storage.Set(key, "1", 10)
		if err != nil {
			fmt.Println("Error setting value in storage:", err)
		}
		return true
	}

	requests, _ := strconv.Atoi(val)
	fmt.Printf("Current requests for key %s: %d\n", key, requests)

	if requests >= limit {
		fmt.Printf("Limit reached for key %s\n", key)
		return false
	}

	err = r.Storage.Incr(key)
	if err != nil {
		fmt.Println("Error incrementing value in storage:", err)
	}
	return true
}
