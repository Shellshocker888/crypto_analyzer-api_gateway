package auth

import "crypto_analyzer-api_gateway/internal/usecase/auth"

type AuthServiceController struct {
	authUsecaseObj *auth.AuthServiceUsecase
}

func NewAuthController(authUsecaseObj *auth.AuthServiceUsecase) *AuthServiceController {
	return &AuthServiceController{authUsecaseObj: authUsecaseObj}
}
