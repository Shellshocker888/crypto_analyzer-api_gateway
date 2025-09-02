package portfolio

import (
	"context"
	"crypto_analyzer-api_gateway/internal/domain/portfolio"
)

type PortfolioUsecase struct {
	portfolioService portfolio.PortfolioServiceContract
}

func NewPortfolioServiceUsecase(portfolioService portfolio.PortfolioServiceContract) *PortfolioUsecase {
	return &PortfolioUsecase{portfolioService: portfolioService}
}

func (u PortfolioUsecase) CreateNewPortfolio(ctx context.Context, name string, isPublic bool) (portfolio.Portfolio, error) {
	return u.portfolioService.CreateNewPortfolio(ctx, name, isPublic)
}

func (u PortfolioUsecase) GetPortfolioContentById(ctx context.Context, portfolioID int) (portfolio.PortfolioContent, error) {
	return u.portfolioService.GetPortfolioContentById(ctx, portfolioID)
}

func (u PortfolioUsecase) UpsertAsset(ctx context.Context, portfolioId int, symbol string, amount float64) error {
	return u.portfolioService.UpsertAsset(ctx, portfolioId, symbol, amount)
}

func (u PortfolioUsecase) DeleteAsset(ctx context.Context, portfolioId int, symbol string) error {
	return u.portfolioService.DeleteAsset(ctx, portfolioId, symbol)
}

func (u PortfolioUsecase) GetAllPortfolios(ctx context.Context) ([]portfolio.Portfolio, error) {
	return u.portfolioService.GetAllPortfolios(ctx)
}

func (u PortfolioUsecase) GetPortfolioHistory(ctx context.Context, id, page, pageSize int32) (portfolio.PortfolioHistory, error) {
	return u.portfolioService.GetPortfolioHistory(ctx, id, page, pageSize)
}

func (u PortfolioUsecase) GetPublicPortfolios(ctx context.Context, userId int) ([]portfolio.PublicPortfolio, error) {
	return u.portfolioService.GetPublicPortfolios(ctx, userId)
}
