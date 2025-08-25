package config

import (
	"crypto_analyzer-api_gateway/internal/config/model"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func getEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("failed to load env %s", key)
	}

	return val, nil
}

func LoadConfig() (*model.Config, error) {
	env := ".env"

	if os.Getenv("APP_ENV") == "test" {
		env = ".env.test"
	}

	err := godotenv.Load(env)
	if err != nil {
		return nil, fmt.Errorf("failed to load env: %w", err)
	}

	port, err := getEnv("PORT")
	if err != nil {
		return nil, fmt.Errorf("failed to load port config: %w", err)
	}

	authServiceURL, err := getEnv("AUTH_SERVICE_URL")
	if err != nil {
		return nil, fmt.Errorf("failed to load authServiceURL config: %w", err)
	}

	portfolioServiceURL, err := getEnv("PORTFOLIO_SERVICE_URL")
	if err != nil {
		return nil, fmt.Errorf("failed to load portfolioServiceURL config: %w", err)
	}

	alertServiceURL, err := getEnv("ALERT_SERVICE_URL")
	if err != nil {
		return nil, fmt.Errorf("failed to load alertServiceURL config: %w", err)
	}

	return &model.Config{
		Port:                port,
		AuthServiceURL:      authServiceURL,
		PortfolioServiceURL: portfolioServiceURL,
		AlertServiceURL:     alertServiceURL,
	}, nil
}
