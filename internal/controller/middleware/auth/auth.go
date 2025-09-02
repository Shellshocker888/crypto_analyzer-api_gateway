package auth

import (
	"context"
	authpb "crypto_analyzer-api_gateway/gen/go/auth"
	"crypto_analyzer-api_gateway/internal/domain/portfolio"
	"crypto_analyzer-api_gateway/internal/infrastructure/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"strings"
)

var _ AuthMiddlewareVerifierContract = (*grpcVerifierAdapter)(nil)

type AuthMiddlewareVerifierContract interface {
	Verify(ctx context.Context, in *authpb.VerifyRequest) (*authpb.VerifyResponse, error)
}

type grpcVerifierAdapter struct {
	grpcClient authpb.AuthServiceClient
}

type AuthMiddlewareVerifier struct {
	authClient AuthMiddlewareVerifierContract
}

func NewAuthMiddlewareVerifier(authClient authpb.AuthServiceClient) *AuthMiddlewareVerifier {
	adapter := &grpcVerifierAdapter{grpcClient: authClient}
	return &AuthMiddlewareVerifier{authClient: adapter}
}

func (m *grpcVerifierAdapter) Verify(ctx context.Context, in *authpb.VerifyRequest) (*authpb.VerifyResponse, error) {
	return m.grpcClient.Verify(ctx, in)
}

func (m *AuthMiddlewareVerifier) AuthVerify(c *fiber.Ctx) error {
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

	res, err := m.authClient.Verify(mdCTX, &authpb.VerifyRequest{})
	if err != nil {
		log.Warn("failed to verify token", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "failed to verify"})
	}

	user := &portfolio.User{
		Id:       res.UserId,
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
