package middleware

import (
	"crypto_analyzer-api_gateway/internal/infrastructure/logger"
	"github.com/gofiber/fiber/v2"
)

func LoggerMiddleware(c *fiber.Ctx) error {
	ctx := logger.WithLogger(c.UserContext(), logger.Log)
	c.SetUserContext(ctx)
	return c.Next()
}

func TraceMiddleware(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := logger.FromContext(ctx)
	log = logger.WithTraceID(ctx, log)
	ctx = logger.WithLogger(ctx, log)
	c.SetUserContext(ctx)
	return c.Next()
}
