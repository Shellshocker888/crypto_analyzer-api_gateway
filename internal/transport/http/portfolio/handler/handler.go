package handler

import (
	portfoliopb "crypto_analyzer-api_gateway/gen/go/portfolio"
	"crypto_analyzer-api_gateway/internal/domain"
	"crypto_analyzer-api_gateway/internal/infrastructure/logger"
	"crypto_analyzer-api_gateway/internal/transport/http/portfolio/dto"
	"crypto_analyzer-api_gateway/internal/transport/http/portfolio/mapper"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strconv"
)

func (g *GRPCClient) CreatePortfolioHandler(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := logger.FromContext(ctx)

	userVal := c.Locals("user")
	user, ok := userVal.(*domain.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.HTTPError{
			Status:  fiber.StatusUnauthorized,
			Error:   "unauthorized",
			Message: "user not found in context",
		})
	}

	userID := user.ID
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("user_id", userID))

	var createPortfolioObj dto.CreatePortfolioObject
	if err := c.BodyParser(&createPortfolioObj); err != nil {
		log.Warn("failed to parse portfolio data", zap.Error(err),
			zap.String("user_id", user.ID))

		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "wrong portfolio data",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	res, err := g.PortfolioClient.CreateNewPortfolio(ctx, &portfoliopb.CreateNewPortfolioRequest{
		Name:     createPortfolioObj.Name,
		IsPublic: wrapperspb.Bool(createPortfolioObj.IsPublic),
	})
	if err != nil {
		httpErr := mapper.GrpcCodeToHTTPError(codes.Unknown, "failed to create new portfolio")

		st, ok := status.FromError(err)
		grpcCode := "unknown"
		if ok {
			httpErr = mapper.GrpcCodeToHTTPError(st.Code(), "failed to create new portfolio")
			grpcCode = st.Code().String()
		}

		log.Error("failed to create portfolio",
			zap.String("grpc_code", grpcCode),
			zap.String("user_id", user.ID),
			zap.String("username", user.Username),
			zap.String("email", user.Email),
			zap.Error(err),
		)

		return c.Status(httpErr.Status).JSON(httpErr)
	}

	log.Info("portfolio created successfully",
		zap.String("user_id", user.ID),
		zap.String("username", user.Username),
		zap.String("email", user.Email),
		zap.Int32("portfolio_id", res.Id),
		zap.String("portfolio_name", res.Name),
	)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":        res.Id,
		"user_id":   userID,
		"name":      res.Name,
		"is_public": res.IsPublic.GetValue(),
	})
}

func (g *GRPCClient) PortfolioContentHandler(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := logger.FromContext(ctx)

	userVal := c.Locals("user")
	user, ok := userVal.(*domain.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.HTTPError{
			Status:  fiber.StatusUnauthorized,
			Error:   "unauthorized",
			Message: "user not found in context",
		})
	}

	userID := user.ID

	mdCTX := metadata.NewOutgoingContext(ctx, metadata.Pairs("user_id", userID))

	portfolioID := c.Params("portfolio_id")
	if portfolioID == "" {
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "portfolio_id is required",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	portfolioIDInt, err := strconv.Atoi(portfolioID)
	if err != nil {
		log.Warn("failed to convert portfolio_id", zap.Error(err),
			zap.String("user_id", userID),
			zap.String("portfolio_id", portfolioID),
		)
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "wrong portfolio_id",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	res, err := g.PortfolioClient.GetPortfolioContentById(mdCTX, &portfoliopb.GetPortfolioContentByIdRequest{
		Id: int32(portfolioIDInt),
	})
	if err != nil {
		httpErr := mapper.GrpcCodeToHTTPError(codes.Unknown, "failed to get portfolio content")
		grpcCode := "unknown"

		if st, ok := status.FromError(err); ok {
			httpErr = mapper.GrpcCodeToHTTPError(st.Code(), "failed to get portfolio content")
			grpcCode = st.Code().String()
		}

		log.Error("failed to get portfolio content",
			zap.String("grpc_code", grpcCode),
			zap.String("user_id", userID),
			zap.String("portfolio_id", portfolioID),
			zap.Error(err),
		)

		return c.Status(httpErr.Status).JSON(httpErr)
	}

	log.Info("portfolio content retrieved successfully",
		zap.String("user_id", userID),
		zap.String("portfolio_id", portfolioID),
		zap.Int("assets_count", len(res.Assets)),
	)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":      user.ID,
		"portfolio_id": portfolioID,
		"assets":       res.Assets,
	})
}

func (g *GRPCClient) UpsertAssetHandler(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := logger.FromContext(ctx)

	userVal := c.Locals("user")
	user, ok := userVal.(*domain.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.HTTPError{
			Status:  fiber.StatusUnauthorized,
			Error:   "unauthorized",
			Message: "user not found in context",
		})
	}

	userID := user.ID

	mdCTX := metadata.NewOutgoingContext(ctx, metadata.Pairs("user_id", userID))

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

	portfolioID := c.Params("portfolio_id")
	if portfolioID == "" {
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "portfolio_id is required",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	portfolioIDInt, err := strconv.Atoi(portfolioID)
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

	_, err = g.PortfolioClient.UpsertAsset(mdCTX, &portfoliopb.UpsertAssetRequest{
		PortfolioId: int32(portfolioIDInt),
		Symbol:      symbol,
		Amount:      amount,
	})
	if err != nil {
		httpErr := mapper.GrpcCodeToHTTPError(codes.Unknown, "failed to upsert asset")
		grpcCode := "unknown"
		if st, ok := status.FromError(err); ok {
			httpErr = mapper.GrpcCodeToHTTPError(st.Code(), "failed to upsert asset")
			grpcCode = st.Code().String()
		}

		log.Error("failed to upsert asset",
			zap.String("grpc_code", grpcCode),
			zap.String("portfolio_id", portfolioID),
			zap.String("symbol", symbol),
			zap.Float64("amount", amount),
			zap.Error(err),
		)

		return c.Status(httpErr.Status).JSON(httpErr)
	}

	log.Info("upsert asset successful",
		zap.String("portfolio_id", portfolioID),
		zap.String("symbol", symbol),
		zap.Float64("amount", amount),
	)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "asset upserted successfully",
	})
}

func (g *GRPCClient) DeleteAssetHandler(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := logger.FromContext(ctx)

	userVal := c.Locals("user")
	user, ok := userVal.(*domain.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.HTTPError{
			Status:  fiber.StatusUnauthorized,
			Error:   "unauthorized",
			Message: "user not found in context",
		})
	}

	userID := user.ID

	mdCTX := metadata.NewOutgoingContext(ctx, metadata.Pairs("user_id", userID))

	portfolioID := c.Params("portfolio_id")
	if portfolioID == "" {
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "portfolio_id is required",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	symbol := c.Params("asset")
	if symbol == "" {
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "symbol is required",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	portfolioIDInt, err := strconv.Atoi(portfolioID)
	if err != nil {
		log.Warn("failed to convert portfolio_id", zap.Error(err))
		httpErr := &dto.HTTPError{
			Status:  fiber.StatusBadRequest,
			Error:   "bad_request",
			Message: "wrong portfolio_id",
		}
		return c.Status(httpErr.Status).JSON(httpErr)
	}

	_, err = g.PortfolioClient.DeleteAsset(mdCTX, &portfoliopb.DeleteAssetRequest{
		PortfolioId: int32(portfolioIDInt),
		Symbol:      symbol,
	})
	if err != nil {
		httpErr := mapper.GrpcCodeToHTTPError(codes.Unknown, "failed to delete asset")
		grpcCode := "unknown"
		if st, ok := status.FromError(err); ok {
			httpErr = mapper.GrpcCodeToHTTPError(st.Code(), "failed to delete asset")
			grpcCode = st.Code().String()
		}

		log.Error("failed to delete asset",
			zap.String("grpc_code", grpcCode),
			zap.String("portfolio_id", portfolioID),
			zap.String("symbol", symbol),
			zap.Error(err),
		)

		return c.Status(httpErr.Status).JSON(httpErr)
	}

	log.Info("delete asset successful",
		zap.String("portfolio_id", portfolioID),
		zap.String("symbol", symbol),
	)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "asset deleted successfully",
	})
}

func (g *GRPCClient) GetUserPortfoliosHandler(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := logger.FromContext(ctx)

	userVal := c.Locals("user")
	user, ok := userVal.(*domain.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.HTTPError{
			Status:  fiber.StatusUnauthorized,
			Error:   "unauthorized",
			Message: "user not found in context",
		})
	}

	userID := user.ID

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("user_id", userID))

	res, err := g.PortfolioClient.GetAllPortfolios(ctx, &emptypb.Empty{})
	if err != nil {
		httpError := mapper.GrpcCodeToHTTPError(codes.Unknown, "failed to get user portfolios")
		grpcCode := "unknown"

		if st, ok := status.FromError(err); ok {
			httpError = mapper.GrpcCodeToHTTPError(st.Code(), "failed to get user portfolios")
			grpcCode = st.Code().String()
		}

		log.Error("failed to get user portfolios",
			zap.String("grpcCode", grpcCode),
			zap.String("user_id", userID),
			zap.Error(err))

		return c.Status(httpError.Status).JSON(httpError)
	}

	log.Info("get user portfolios successful", zap.String("user_id", userID))

	portfolios := mapper.MapPortfolios(res.Portfolios)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"portfolios": portfolios,
	})

}

func (g *GRPCClient) GetPortfolioHistoryHandler(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := logger.FromContext(ctx)

	userVal := c.Locals("user")
	user, ok := userVal.(*domain.User)
	if !ok || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.HTTPError{
			Status:  fiber.StatusUnauthorized,
			Error:   "unauthorized",
			Message: "user not found in context",
		})
	}

	userID := user.ID

	mdCTX := metadata.NewOutgoingContext(ctx, metadata.Pairs("user_id", userID))

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

	portfolioID := portfolioHistoryData.Id
	if portfolioID == 0 {
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

	res, err := g.PortfolioClient.GetPortfolioHistory(mdCTX, &portfoliopb.GetPortfolioHistoryRequest{
		Id:       int32(portfolioID),
		Page:     int32(page),
		PageSize: int32(size),
	})
	if err != nil {
		httpErr := mapper.GrpcCodeToHTTPError(codes.Unknown, "failed to get portfolio history")
		grpcCode := "unknown"
		if st, ok := status.FromError(err); ok {
			httpErr = mapper.GrpcCodeToHTTPError(st.Code(), "failed to get portfolio history")
			grpcCode = st.Code().String()
		}

		log.Error("failed to get portfolio history",
			zap.String("grpc_code", grpcCode),
			zap.String("user_id", userID),
			zap.Int("portfolio_id", portfolioID),
			zap.Error(err),
		)

		return c.Status(httpErr.Status).JSON(httpErr)
	}

	log.Info("get portfolio history successful",
		zap.String("user_id", userID),
		zap.Int("portfolio_id", portfolioID))

	history := mapper.MapHistory(res.History)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"historyPortfolios": history,
	})
}

func (g *GRPCClient) GetUserPublicPortfoliosHandler(c *fiber.Ctx) error {
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

	res, err := g.PortfolioClient.GetPublicPortfolios(ctx, &portfoliopb.GetPublicPortfoliosRequest{
		UserId: int32(userIDInt),
	})
	if err != nil {
		httpErr := mapper.GrpcCodeToHTTPError(codes.Unknown, "failed to get public portfolios")
		grpcCode := "unknown"

		if st, ok := status.FromError(err); ok {
			httpErr = mapper.GrpcCodeToHTTPError(st.Code(), "failed to get public portfolios")
			grpcCode = st.Code().String()
		}

		log.Error("failed to get public portfolios",
			zap.String("grpc_code", grpcCode),
			zap.String("user_id", userID),
			zap.Error(err),
		)

		return c.Status(httpErr.Status).JSON(httpErr)
	}

	log.Info("public portfolios retrieved successfully", zap.String("user_id", userID))

	publicPortfolios := mapper.MapPublicPortfolios(res.Portfolios)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":           userID,
		"public_portfolios": publicPortfolios,
	})
}
