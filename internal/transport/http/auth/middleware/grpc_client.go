package middleware

import (
	authpb "crypto_analyzer-api_gateway/gen/go/auth"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	AuthClient authpb.AuthServiceClient
}

func NewGRPCClient(conn *grpc.ClientConn) *GRPCClient {
	return &GRPCClient{
		AuthClient: authpb.NewAuthServiceClient(conn),
	}
}
