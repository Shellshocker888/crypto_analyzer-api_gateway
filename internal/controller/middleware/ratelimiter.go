package middleware

import (
	"crypto_analyzer-api_gateway/internal/domain"
	"crypto_analyzer-api_gateway/internal/infrastructure/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strconv"
)

type RateLimiterMiddleware struct {
	limiter domain.RateLimiterContract
	keyFn   func(c *fiber.Ctx) string
	limit   int
	rate    float64
}

func NewRateLimiterMiddleware(limiter domain.RateLimiterContract, keyFn func(c *fiber.Ctx) string,
	limit int, rate float64) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		limiter: limiter,
		keyFn:   keyFn,
		limit:   limit,
		rate:    rate,
	}
}

func (l *RateLimiterMiddleware) Handler(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := logger.FromContext(ctx)
	key := l.keyFn(c)

	var cost int
	switch c.Method() {
	case fiber.MethodGet:
		cost = 1
	case fiber.MethodPost, fiber.MethodPut:
		cost = 2
	case fiber.MethodDelete:
		cost = 3
	default:
		cost = 1
	}

	if cost > l.limit {
		log.Warn("forbidden operation", zap.String("key", key))
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "you dont have permission for this operation",
		})
	}

	c.Set("X-RateLimit-Limit", strconv.Itoa(l.limit))
	c.Set("X-RateLimit-Rate", strconv.FormatFloat(l.rate, 'f', 2, 64))

	allowed, retry, err := l.limiter.TryAllow(ctx, key, cost)
	if err != nil {
		log.Error("rate limiter error", zap.String("key", key), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	if !allowed {
		log.Warn("too many requests", zap.String("key", key))
		c.Set("Retry-After", strconv.Itoa(int(retry.Seconds())))
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "too many requests",
		})
	}

	return c.Next()
}
