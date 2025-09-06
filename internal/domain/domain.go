package domain

import (
	"context"
	"time"
)

type RateLimiterContract interface {
	TryAllow(ctx context.Context, key string, cost int) (allowed bool, tryAfter time.Duration, err error)
}
