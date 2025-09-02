package mapper

import (
	portfoliopb "crypto_analyzer-api_gateway/gen/go/portfolio"
	"crypto_analyzer-api_gateway/internal/domain/portfolio"
)

func MapPortfoliosToDomainPortfolios(portfolios []*portfoliopb.AllUserPortfolio) []portfolio.Portfolio {
	portfoliosSlice := make([]portfolio.Portfolio, 0, len(portfolios))

	for _, v := range portfolios {
		portfoliosSlice = append(portfoliosSlice, portfolio.Portfolio{
			Id:       v.Id,
			Name:     v.Name,
			IsPublic: v.IsPublic.GetValue(),
		})
	}

	return portfoliosSlice
}

func MapProtoToDomainHistory(portfolioHistory map[string]*portfoliopb.PricePoints) portfolio.PortfolioHistory {
	history := portfolio.PortfolioHistory{History: make(map[string][]portfolio.PricePoint, len(portfolioHistory))}

	for key, v := range portfolioHistory {
		if v == nil {
			continue
		}

		history.History[key] = make([]portfolio.PricePoint, 0, len(v.Points))

		for _, m := range v.Points {
			history.History[key] = append(history.History[key], portfolio.PricePoint{
				Timestamp: m.Timestamp,
				Value:     m.Value,
			})
		}
	}

	return history
}

func MapProtoToDomainPublicPortfolios(portfolios []*portfoliopb.PublicPortfolio) []portfolio.PublicPortfolio {
	publicPortfolios := make([]portfolio.PublicPortfolio, 0, len(portfolios))

	for _, v := range portfolios {
		publicPortfolios = append(publicPortfolios, portfolio.PublicPortfolio{
			PortfolioId: v.PortfolioId,
			Name:        v.Name,
			Assets:      v.Assets,
		})
	}

	return publicPortfolios
}
