package portfolio

import (
	"context"
)

type User struct {
	Id       string
	Username string
	Email    string
}

type Portfolio struct {
	Id       int32
	Name     string
	IsPublic bool
}

type PortfolioContent struct {
	Assets map[string]float64
}

type PricePoint struct {
	Timestamp string
	Value     float64
}

type PortfolioHistory struct {
	History map[string][]PricePoint
}

type PublicPortfolio struct {
	PortfolioId int32
	Name        string
	Assets      map[string]float64
}

type PortfolioServiceContract interface {
	CreateNewPortfolio(ctx context.Context, name string, isPublic bool) (Portfolio, error)
	GetPortfolioContentById(ctx context.Context, portfolioID int) (PortfolioContent, error)
	UpsertAsset(ctx context.Context, portfolioId int, symbol string, amount float64) error
	DeleteAsset(ctx context.Context, portfolioId int, symbol string) error
	GetAllPortfolios(ctx context.Context) ([]Portfolio, error)
	GetPortfolioHistory(ctx context.Context, id, page, pageSize int32) (PortfolioHistory, error)
	GetPublicPortfolios(ctx context.Context, userId int) ([]PublicPortfolio, error)
}
