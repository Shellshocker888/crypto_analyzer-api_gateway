package portfolio

import (
	"crypto_analyzer-api_gateway/internal/controller/portfolio/dto"
	"crypto_analyzer-api_gateway/internal/controller/portfolio/mapper"
	"crypto_analyzer-api_gateway/internal/domain/portfolio"
	"crypto_analyzer-api_gateway/internal/infrastructure/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
)

func (con PortfolioController) GetPortfolioContentById(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := logger.FromContext(ctx)

	userVal := c.Locals("user")
	user, ok := userVal.(*portfolio.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.HTTPError{
			Status:  fiber.StatusUnauthorized,
			Error:   "unauthorized",
			Message: "user not found in context",
		})
	}

	userID := user.Id

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("user_id", userID))

	portfolioId := c.Params("portfolio_id")
	if portfolioId == "" {
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "portfolio_id is required",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	portfolioIDInt, err := strconv.Atoi(portfolioId)
	if err != nil {
		log.Warn("failed to convert portfolio_id", zap.Error(err),
			zap.String("user_id", userID),
			zap.String("portfolio_id", portfolioId),
		)
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "wrong portfolio_id",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	res, err := con.portfolioUsecaseObj.GetPortfolioContentById(ctx, portfolioIDInt)
	if err != nil {
		httpErr := mapper.GrpcCodeToHTTPError(codes.Unknown, "failed to get portfolio content")

		if st, ok := status.FromError(err); ok {
			httpErr = mapper.GrpcCodeToHTTPError(st.Code(), "failed to get portfolio content")
		}

		log.Error("failed to get portfolio content",
			zap.String("user_id", userID),
			zap.String("portfolio_id", portfolioId),
			zap.Error(err),
		)

		return c.Status(httpErr.Status).JSON(httpErr)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":      user.Id,
		"portfolio_id": portfolioId,
		"assets":       res.Assets,
	})
}
