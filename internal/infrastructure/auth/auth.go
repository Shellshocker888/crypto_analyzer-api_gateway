package auth

import (
	authpb "crypto_analyzer-api_gateway/gen/go/auth"
	domain "crypto_analyzer-api_gateway/internal/domain/auth"
)

type AuthServiceClient struct {
	grpcClient authpb.AuthServiceClient
}

func NewAuthServiceClient(grpcClient authpb.AuthServiceClient) domain.AuthServiceContract {
	return AuthServiceClient{grpcClient: grpcClient}
}
