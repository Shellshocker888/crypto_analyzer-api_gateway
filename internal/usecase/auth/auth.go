package auth

import domain "crypto_analyzer-api_gateway/internal/domain/auth"

type AuthServiceUsecase struct {
	authService domain.AuthServiceContract
}

func NewAuthServiceUsecase(authService domain.AuthServiceContract) *AuthServiceUsecase {
	return &AuthServiceUsecase{authService: authService}
}
