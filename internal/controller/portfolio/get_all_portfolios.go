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

func (con PortfolioController) GetAllPortfolios(c *fiber.Ctx) error {
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

	portfolios, err := con.portfolioUsecaseObj.GetAllPortfolios(ctx)
	if err != nil {
		httpError := mapper.GrpcCodeToHTTPError(codes.Unknown, "failed to get all portfolios")
		if st, ok := status.FromError(err); ok {
			httpError = mapper.GrpcCodeToHTTPError(st.Code(), "failed to get all portfolios")
		}

		log.Error("failed to get all portfolios",
			zap.String("user_id", userId),
			zap.Error(err))

		return c.Status(httpError.Status).JSON(httpError)
	}

	portfoliosDTO := mapper.MapPortfolios(portfolios)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"portfolios": portfoliosDTO,
	})

}
