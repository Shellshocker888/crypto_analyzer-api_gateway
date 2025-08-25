package mapper

import (
	portfoliopb "crypto_analyzer-api_gateway/gen/go/portfolio"
	"crypto_analyzer-api_gateway/internal/transport/http/portfolio/dto"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
)

func GrpcCodeToHTTPError(code codes.Code, msg string) *dto.HTTPError {
	switch code {
	case codes.Canceled:
		return &dto.HTTPError{fiber.StatusRequestTimeout, "canceled", "request was canceled"}
	case codes.Unknown, codes.Internal, codes.DataLoss:
		return &dto.HTTPError{fiber.StatusInternalServerError, "internal_error", "internal server error"}
	case codes.InvalidArgument, codes.OutOfRange:
		return &dto.HTTPError{fiber.StatusBadRequest, "invalid_argument", msg}
	case codes.DeadlineExceeded:
		return &dto.HTTPError{fiber.StatusGatewayTimeout, "timeout", "request timeout"}
	case codes.NotFound:
		return &dto.HTTPError{fiber.StatusNotFound, "not_found", msg}
	case codes.AlreadyExists, codes.Aborted:
		return &dto.HTTPError{fiber.StatusConflict, "conflict", msg}
	case codes.PermissionDenied:
		return &dto.HTTPError{fiber.StatusForbidden, "forbidden", msg}
	case codes.ResourceExhausted:
		return &dto.HTTPError{fiber.StatusTooManyRequests, "rate_limit", "rate limit exceeded"}
	case codes.FailedPrecondition:
		return &dto.HTTPError{fiber.StatusPreconditionFailed, "failed_precondition", msg}
	case codes.Unimplemented:
		return &dto.HTTPError{fiber.StatusNotImplemented, "not_implemented", msg}
	case codes.Unavailable:
		return &dto.HTTPError{fiber.StatusServiceUnavailable, "unavailable", "service unavailable"}
	case codes.Unauthenticated:
		return &dto.HTTPError{fiber.StatusUnauthorized, "unauthenticated", msg}
	default:
		return &dto.HTTPError{fiber.StatusInternalServerError, "unexpected_error", "unexpected error"}
	}
}

func MapPortfolios(portfolios []*portfoliopb.AllUserPortfolio) []dto.Portfolio {
	portfoliosSlice := make([]dto.Portfolio, 0, len(portfolios))

	for _, v := range portfolios {
		portfoliosSlice = append(portfoliosSlice, dto.Portfolio{
			Id:       int(v.Id),
			Name:     v.Name,
			IsPublic: v.IsPublic.GetValue(),
		})
	}

	return portfoliosSlice
}

func MapHistory(portfolioHistory map[string]*portfoliopb.PricePoints) map[string][]*dto.PortfolioHistory {
	history := make(map[string][]*dto.PortfolioHistory)

	for key, v := range portfolioHistory {
		if v == nil {
			continue
		}
		for _, m := range v.Points {
			history[key] = append(history[key], &dto.PortfolioHistory{
				Time:  m.Timestamp,
				Price: m.Value,
			})
		}
	}

	return history
}

func MapPublicPortfolios(portfolios []*portfoliopb.PublicPortfolio) []dto.PublicPortfolio {
	publicPortfolios := make([]dto.PublicPortfolio, 0, len(portfolios))

	for _, v := range portfolios {
		publicPortfolios = append(publicPortfolios, dto.PublicPortfolio{
			PortfolioId: v.PortfolioId,
			Name:        v.Name,
			Assets:      v.Assets,
		})
	}

	return publicPortfolios
}
