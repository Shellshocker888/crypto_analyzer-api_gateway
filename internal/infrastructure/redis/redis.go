package redis

import (
	"context"
	"crypto_analyzer-api_gateway/internal/config/model"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func InitRedisClient(ctx context.Context, cfg *model.RedisConfig) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.SessionDB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return client, nil
}
