package ratelimiter

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DiegoSenaa/go-rater-limiter/internal/redisclient"
	"github.com/go-redis/redis/v8"
)

func AllowRequest(identifier, idType string) bool {
	var limit int
	if idType == "token" {
		limit = redisclient.RateLimitToken
	} else {
		limit = redisclient.RateLimitIP
	}

	key := fmt.Sprintf("rate_limiter:%s:%s", idType, identifier)
	val, err := redisclient.RedisClient.Get(redisclient.Ctx, key).Result()

	if err == redis.Nil {
		fmt.Printf("Setting initial value for key: %s\n", key)
		err = redisclient.RedisClient.Set(redisclient.Ctx, key, "1", 10*time.Second).Err()
		if err != nil {
			fmt.Println("Error setting value in Redis:", err)
		}
		return true
	} else if err != nil {
		fmt.Println("Error getting value from Redis:", err)
		return false
	}

	requests, _ := strconv.Atoi(val)
	fmt.Printf("Current requests for key %s: %d\n", key, requests)

	if requests >= limit {
		fmt.Printf("Limit reached for key %s\n", key)
		return false
	}

	err = redisclient.RedisClient.Incr(redisclient.Ctx, key).Err()
	if err != nil {
		fmt.Println("Error incrementing value in Redis:", err)
	}
	return true
}
