package handler

import (
	portfoliopb "crypto_analyzer-api_gateway/gen/go/portfolio"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	PortfolioClient portfoliopb.PortfolioServiceClient
}

func NewGRPCClient(conn *grpc.ClientConn) *GRPCClient {
	return &GRPCClient{PortfolioClient: portfoliopb.NewPortfolioServiceClient(conn)}
}
