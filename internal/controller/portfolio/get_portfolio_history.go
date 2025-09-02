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

func (con PortfolioController) GetPortfolioHistory(c *fiber.Ctx) error {
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

	var portfolioHistoryData dto.PortfolioHistoryData
	if err := c.BodyParser(&portfolioHistoryData); err != nil {
		log.Warn("failed to parse portfolio history request data", zap.Error(err))
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "wrong portfolio history request data",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	portfolioId := portfolioHistoryData.Id
	if portfolioId == 0 {
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "portfolio_id is required",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	size := portfolioHistoryData.Size
	if size == 0 {
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "size is required",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	page := portfolioHistoryData.Page
	if page == 0 {
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "page is required",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	res, err := con.portfolioUsecaseObj.GetPortfolioHistory(ctx, int32(portfolioId), int32(page), int32(size))
	if err != nil {
		httpErr := mapper.GrpcCodeToHTTPError(codes.Unknown, "failed to get portfolio history")

		if st, ok := status.FromError(err); ok {
			httpErr = mapper.GrpcCodeToHTTPError(st.Code(), "failed to get portfolio history")
		}

		log.Error("failed to get portfolio history",
			zap.String("user_id", userId),
			zap.Int("portfolio_id", portfolioId),
			zap.Error(err),
		)

		return c.Status(httpErr.Status).JSON(httpErr)
	}

	history := mapper.MapDomainToDTOHistory(res)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"portfolio_history": history,
	})
}
