package middleware

import (
	"crypto_analyzer-api_gateway/internal/infrastructure/metrics"
	"github.com/gofiber/fiber/v2"
	"time"
)

func MetricsMiddleware(c *fiber.Ctx) error {
	if c.Path() == "/metrics" {
		return c.Next()
	}

	start := time.Now()
	err := c.Next()

	path := string(c.Request().URI().Path())
	status := c.Response().StatusCode()

	metrics.IncRequest(path, status)
	metrics.ObserveDuration(path, time.Since(start))

	return err
}
