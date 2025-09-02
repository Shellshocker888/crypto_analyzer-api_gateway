package portfolio

import (
	"crypto_analyzer-api_gateway/internal/controller/portfolio/dto"
	"crypto_analyzer-api_gateway/internal/controller/portfolio/mapper"
	"crypto_analyzer-api_gateway/internal/infrastructure/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

func (con PortfolioController) GetPublicPortfolios(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := logger.FromContext(ctx)

	userID := c.Params("user_id")
	if userID == "" {
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "user_id is required",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		log.Warn("failed to convert user_id", zap.Error(err))
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "wrong user_id",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	res, err := con.portfolioUsecaseObj.GetPublicPortfolios(ctx, userIDInt)
	if err != nil {
		httpErr := mapper.GrpcCodeToHTTPError(codes.Unknown, "failed to get public portfolios")

		if st, ok := status.FromError(err); ok {
			httpErr = mapper.GrpcCodeToHTTPError(st.Code(), "failed to get public portfolios")
		}

		log.Error("failed to get public portfolios",
			zap.String("user_id", userID),
			zap.Error(err),
		)

		return c.Status(httpErr.Status).JSON(httpErr)
	}

	publicPortfolios := mapper.MapPublicPortfolios(res)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":           userID,
		"public_portfolios": publicPortfolios,
	})
}
