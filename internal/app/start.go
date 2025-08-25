package app

import (
	"context"
	"crypto_analyzer-api_gateway/internal/config"
	"crypto_analyzer-api_gateway/internal/infrastructure/logger"
	auth2 "crypto_analyzer-api_gateway/internal/transport/http/auth/middleware"
	"crypto_analyzer-api_gateway/internal/transport/http/portfolio/handler"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	logStd "log"
)

func Start(ctx context.Context) error {
	err := logger.InitLogger()
	if err != nil {
		logStd.Printf("failed to init logger: %v", err)
		return fmt.Errorf("failed to init logger: %w", err)
	}
	defer logger.SyncLogger()

	log := logger.Log

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("failed to load config", zap.Error(err))
		return fmt.Errorf("failed to load config: %w", err)
	}

	authConn, err := grpc.NewClient(cfg.AuthServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect auth service", zap.Error(err))
		return fmt.Errorf("failed to connect auth service: %w", err)
	}
	defer authConn.Close()

	portfolioConn, err := grpc.NewClient(cfg.PortfolioServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect portfolio service", zap.Error(err))
		return fmt.Errorf("failed to connect portfolio service: %w", err)
	}
	defer portfolioConn.Close()

	app := fiber.New()

	app.Use(auth2.LoggerMiddleware())
	app.Use(auth2.TraceMiddleware())

	authClient := auth2.NewGRPCClient(authConn)
	portfolioClient := handler.NewGRPCClient(portfolioConn)

	app.Get("/auth/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Post("/portfolios", authClient.WithAuth(), portfolioClient.CreatePortfolioHandler)

	log.Info("Starting API Gateway", zap.String("port", cfg.Port))
	if err = app.Listen(":" + cfg.Port); err != nil {
		log.Error("failed to start gateway", zap.Error(err))
	}

	return nil
}
