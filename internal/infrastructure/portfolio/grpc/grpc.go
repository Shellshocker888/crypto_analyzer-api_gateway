package grpc

import (
	"context"
	portfoliopb "crypto_analyzer-api_gateway/gen/go/portfolio"
	"crypto_analyzer-api_gateway/internal/domain/portfolio"
	"crypto_analyzer-api_gateway/internal/infrastructure/logger"
	"crypto_analyzer-api_gateway/internal/infrastructure/portfolio/mapper"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type portfolioServiceClient struct {
	GRPCClient portfoliopb.PortfolioServiceClient
}

func NewPortfolioServiceClient(GRPCClient portfoliopb.PortfolioServiceClient) portfolio.PortfolioServiceContract {
	return &portfolioServiceClient{GRPCClient: GRPCClient}
}

func (c *portfolioServiceClient) CreateNewPortfolio(ctx context.Context, name string, isPublic bool) (portfolio.Portfolio, error) {
	log := logger.FromContext(ctx)

	res, err := c.GRPCClient.CreateNewPortfolio(ctx, &portfoliopb.CreateNewPortfolioRequest{
		Name:     name,
		IsPublic: &wrapperspb.BoolValue{Value: isPublic},
	})
	if err != nil {
		st, _ := status.FromError(err)
		log.Error("failed to create portfolio via gRPC",
			zap.String("grpc_code", st.Code().String()),
			zap.Error(err),
		)

		return portfolio.Portfolio{}, err
	}

	return portfolio.Portfolio{
		Id:       res.Id,
		Name:     res.Name,
		IsPublic: res.IsPublic.GetValue(),
	}, nil
}

func (c *portfolioServiceClient) GetPortfolioContentById(ctx context.Context, portfolioID int) (portfolio.PortfolioContent, error) {
	log := logger.FromContext(ctx)

	res, err := c.GRPCClient.GetPortfolioContentById(ctx, &portfoliopb.GetPortfolioContentByIdRequest{
		Id: int32(portfolioID),
	})
	if err != nil {
		st, _ := status.FromError(err)
		log.Error("failed to get portfolio content via gRPC",
			zap.String("grpc_code", st.Code().String()),
			zap.Error(err),
		)

		return portfolio.PortfolioContent{}, err
	}

	return portfolio.PortfolioContent{
		Assets: res.Assets,
	}, nil
}

func (c *portfolioServiceClient) UpsertAsset(ctx context.Context, portfolioId int, symbol string, amount float64) error {
	log := logger.FromContext(ctx)

	_, err := c.GRPCClient.UpsertAsset(ctx, &portfoliopb.UpsertAssetRequest{
		PortfolioId: int32(portfolioId),
		Symbol:      symbol,
		Amount:      amount,
	})
	if err != nil {
		st, _ := status.FromError(err)
		log.Error("failed to upsert asset via gRPC",
			zap.String("grpc_code", st.Code().String()),
			zap.Error(err),
		)

		return err
	}

	return nil
}

func (c *portfolioServiceClient) DeleteAsset(ctx context.Context, portfolioId int, symbol string) error {
	log := logger.FromContext(ctx)

	_, err := c.GRPCClient.DeleteAsset(ctx, &portfoliopb.DeleteAssetRequest{
		PortfolioId: int32(portfolioId),
		Symbol:      symbol,
	})
	if err != nil {
		st, _ := status.FromError(err)
		log.Error("failed to delete asset via gRPC",
			zap.String("grpc_code", st.Code().String()),
			zap.Error(err),
		)

		return err
	}

	return nil
}

func (c *portfolioServiceClient) GetAllPortfolios(ctx context.Context) ([]portfolio.Portfolio, error) {
	log := logger.FromContext(ctx)

	res, err := c.GRPCClient.GetAllPortfolios(ctx, &emptypb.Empty{})
	if err != nil {
		st, _ := status.FromError(err)
		log.Error("failed to get all portfolios via gRPC",
			zap.String("grpc_code", st.Code().String()),
			zap.Error(err),
		)

		return nil, err
	}

	portfolios := mapper.MapPortfoliosToDomainPortfolios(res.Portfolios)
	return portfolios, nil
}

func (c *portfolioServiceClient) GetPortfolioHistory(ctx context.Context, id, page, pageSize int32) (portfolio.PortfolioHistory, error) {
	log := logger.FromContext(ctx)

	res, err := c.GRPCClient.GetPortfolioHistory(ctx, &portfoliopb.GetPortfolioHistoryRequest{
		Id:       id,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		st, _ := status.FromError(err)
		log.Error("failed to get portfolio history via gRPC",
			zap.String("grpc_code", st.Code().String()),
			zap.Error(err),
		)

		return portfolio.PortfolioHistory{}, err
	}

	portfolioHistory := mapper.MapProtoToDomainHistory(res.History)
	return portfolioHistory, nil
}

func (c *portfolioServiceClient) GetPublicPortfolios(ctx context.Context, userId int) ([]portfolio.PublicPortfolio, error) {
	log := logger.FromContext(ctx)

	res, err := c.GRPCClient.GetPublicPortfolios(ctx, &portfoliopb.GetPublicPortfoliosRequest{UserId: int32(userId)})
	if err != nil {
		st, _ := status.FromError(err)
		log.Error("failed to get public portfolios via gRPC",
			zap.String("grpc_code", st.Code().String()),
			zap.Error(err),
		)

		return nil, err
	}

	portfolios := mapper.MapProtoToDomainPublicPortfolios(res.Portfolios)
	return portfolios, nil
}
