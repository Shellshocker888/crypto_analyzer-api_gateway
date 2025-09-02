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
)

func (con PortfolioController) CreateNewPortfolio(c *fiber.Ctx) error {
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

	var createPortfolioObj dto.CreatePortfolioObject
	if err := c.BodyParser(&createPortfolioObj); err != nil {
		log.Warn("failed to parse portfolio data", zap.Error(err),
			zap.String("user_id", user.Id))

		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "wrong portfolio data",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	res, err := con.portfolioUsecaseObj.CreateNewPortfolio(ctx, createPortfolioObj.Name, createPortfolioObj.IsPublic)
	if err != nil {
		httpErr := mapper.GrpcCodeToHTTPError(codes.Unknown, "failed to create new portfolio")

		st, ok := status.FromError(err)

		if ok {
			httpErr = mapper.GrpcCodeToHTTPError(st.Code(), "failed to create new portfolio")
		}

		log.Error("failed to create portfolio",
			zap.String("user_id", user.Id),
			zap.Error(err),
		)

		return c.Status(httpErr.Status).JSON(httpErr)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":   userId,
		"id":        res.Id,
		"name":      res.Name,
		"is_public": res.IsPublic,
	})
}
