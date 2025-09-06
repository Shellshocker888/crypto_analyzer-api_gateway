package ratelimiter

import (
	"context"
	"crypto_analyzer-api_gateway/internal/infrastructure/logger"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

//go:embed token_bucket.lua
var luaTokenBucket string

type RateLimiter struct {
	client *redis.Client
	script *redis.Script
	limit  int
	rate   float64
}

func NewRateLimiter(client *redis.Client, limit int, rate float64) *RateLimiter {
	return &RateLimiter{
		client: client,
		script: redis.NewScript(luaTokenBucket),
		limit:  limit,
		rate:   rate,
	}
}

func (l *RateLimiter) TryAllow(ctx context.Context, key string, cost int) (bool, time.Duration, error) {
	log := logger.FromContext(ctx)
	now := time.Now().UnixMilli()

	res, err := l.script.Run(ctx, l.client, []string{key}, l.limit, l.rate, now, cost).Result()
	if err != nil {
		log.Error("error to complete rate limiter script",
			zap.String("key", key),
			zap.Error(err))

		return false, 0, err
	}

	array := res.([]interface{})

	// allowed
	allowed, ok := array[0].(int64)
	if !ok {
		return false, 0, fmt.Errorf("unexpected type for allowed: %T", array[0])
	}

	// tokens
	var tokens float64
	switch v := array[1].(type) {
	case int64:
		tokens = float64(v)
	case float64:
		tokens = v
	default:
		return false, 0, fmt.Errorf("unexpected type for tokens: %T", array[1])
	}

	// lastRefill
	lastRefill, ok := array[2].(int64)
	if !ok {
		return false, 0, fmt.Errorf("unexpected type for lastRefill: %T", array[2])
	}

	if allowed != 1 {
		retryAfter := time.Duration(float64(cost) / l.rate * float64(time.Second))
		return false, retryAfter, nil
	}

	_ = tokens
	_ = lastRefill

	return true, 0, nil

}
