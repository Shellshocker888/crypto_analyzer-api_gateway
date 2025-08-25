package middleware

import (
	authpb "crypto_analyzer-api_gateway/gen/go/auth"
	"crypto_analyzer-api_gateway/internal/domain"
	"crypto_analyzer-api_gateway/internal/infrastructure/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"strings"
)

func (m *GRPCClient) WithAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := c.UserContext()
		log := logger.FromContext(ctx)

		// Проверка токена
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Warn("missing auth header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		mdCTX := metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", token))

		res, err := m.AuthClient.Verify(mdCTX, &authpb.VerifyRequest{})
		if err != nil {
			log.Warn("failed to verify token", zap.Error(err))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "failed to verify"})
		}

		user := &domain.User{
			ID:       res.UserId,
			Username: res.Username,
			Email:    res.Email,
		}

		c.Locals("user", user)

		log = logger.WithTraceID(ctx, log).With(
			zap.String("userID", res.UserId),
			zap.String("username", res.Username),
		)

		return c.Next()
	}
}

func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := logger.WithLogger(c.UserContext(), logger.Log)
		c.SetUserContext(ctx)
		return c.Next()
	}
}

func TraceMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		log := logger.FromContext(ctx)

		log = logger.WithTraceID(ctx, log)

		ctx = logger.WithLogger(ctx, log)
		c.SetUserContext(ctx)

		return c.Next()
	}
}
