package dto

type CreatePortfolioObject struct {
	Name     string `json:"name"`
	IsPublic bool   `json:"is_public"`
}

type UpsertAssetObject struct {
	Symbol string  `json:"symbol"`
	Amount float64 `json:"amount"`
}

type HTTPError struct {
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type Portfolio struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsPublic bool   `json:"isPublic"`
}

type PortfolioHistoryData struct {
	Id   int `json:"id"`
	Page int `json:"page"`
	Size int `json:"size"`
}

type PricePoint struct {
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
}

type PortfolioHistory struct {
	History map[string][]PricePoint `json:"history"`
}

type PublicPortfolio struct {
	PortfolioId int32              `json:"portfolioId"`
	Name        string             `json:"name"`
	Assets      map[string]float64 `json:"assets"`
}
