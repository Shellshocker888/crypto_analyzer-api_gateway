package portfolio

import (
	"crypto_analyzer-api_gateway/internal/usecase/portfolio"
)

type PortfolioController struct {
	portfolioUsecaseObj *portfolio.PortfolioUsecase
}

func NewPortfolioController(portfolioUsecaseObj *portfolio.PortfolioUsecase) *PortfolioController {
	return &PortfolioController{portfolioUsecaseObj: portfolioUsecaseObj}
}
