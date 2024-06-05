package redisclient

import (
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var Ctx = context.Background()
var RedisClient *redis.Client
var RateLimitIP, RateLimitToken, BlockDuration int

func InitRedisClient() {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	RateLimitIP, _ = strconv.Atoi(os.Getenv("RATE_LIMIT_IP"))
	RateLimitToken, _ = strconv.Atoi(os.Getenv("RATE_LIMIT_TOKEN"))
	BlockDuration, _ = strconv.Atoi(os.Getenv("BLOCK_DURATION"))

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})
}
