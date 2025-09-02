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

func (con PortfolioController) UpsertAsset(c *fiber.Ctx) error {
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

	userId := user.Id

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("user_id", userId))

	var upsertAssetObj dto.UpsertAssetObject
	if err := c.BodyParser(&upsertAssetObj); err != nil {
		log.Warn("failed to parse upsert data", zap.Error(err))
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "wrong upsert data",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	portfolioId := c.Params("portfolio_id")
	if portfolioId == "" {
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "portfolio_id is required",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	portfolioIdInt, err := strconv.Atoi(portfolioId)
	if err != nil {
		log.Warn("failed to convert portfolio_id", zap.Error(err))
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "wrong portfolio_id",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	amount := upsertAssetObj.Amount
	symbol := upsertAssetObj.Symbol

	err = con.portfolioUsecaseObj.UpsertAsset(ctx, portfolioIdInt, symbol, amount)
	if err != nil {
		httpErr := mapper.GrpcCodeToHTTPError(codes.Unknown, "failed to upsert asset")
		if st, ok := status.FromError(err); ok {
			httpErr = mapper.GrpcCodeToHTTPError(st.Code(), "failed to upsert asset")
		}

		log.Error("failed to upsert asset",
			zap.String("user_id", userId),
			zap.String("portfolio_id", portfolioId),
			zap.String("symbol", symbol),
			zap.Error(err),
		)

		return c.Status(httpErr.Status).JSON(httpErr)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "asset upserted successfully",
	})
}
