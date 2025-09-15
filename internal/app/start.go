package app

import (
	"context"
	authpb "crypto_analyzer-api_gateway/gen/go/auth"
	portfoliopb "crypto_analyzer-api_gateway/gen/go/portfolio"
	"crypto_analyzer-api_gateway/internal/config"
	"crypto_analyzer-api_gateway/internal/controller/middleware"
	"crypto_analyzer-api_gateway/internal/controller/middleware/auth"
	portfolioController "crypto_analyzer-api_gateway/internal/controller/portfolio"
	"crypto_analyzer-api_gateway/internal/infrastructure/logger"
	"crypto_analyzer-api_gateway/internal/infrastructure/metrics"
	portfolioGRPC "crypto_analyzer-api_gateway/internal/infrastructure/portfolio/grpc"
	"crypto_analyzer-api_gateway/internal/infrastructure/ratelimiter"
	"crypto_analyzer-api_gateway/internal/infrastructure/redis"
	"crypto_analyzer-api_gateway/internal/usecase/portfolio"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	redisClient, err := redis.InitRedisClient(ctx, cfg.RedisCfg)
	if err != nil {
		log.Error("failed to init redis client", zap.Error(err))
		return fmt.Errorf("failed to init redis client: %w", err)
	}
	defer redisClient.Close()

	ratelimiterObj := ratelimiter.NewRateLimiter(redisClient, 5, 0.3)
	rlMw := middleware.NewRateLimiterMiddleware(ratelimiterObj, func(c *fiber.Ctx) string {
		if uid := c.Get("X-User-ID"); uid != "" {
			return "user:" + uid
		}
		return "ip:" + c.IP()
	}, 10, 5)

	authConn, err := grpc.NewClient(cfg.AuthServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect auth service", zap.Error(err))
		return fmt.Errorf("failed to connect auth service: %w", err)
	}
	defer authConn.Close()

	authClientProto := authpb.NewAuthServiceClient(authConn)

	authMiddlewareVerifier := auth.NewAuthMiddlewareVerifier(authClientProto)

	portfolioConn, err := grpc.NewClient(cfg.PortfolioServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect portfolio service", zap.Error(err))
		return fmt.Errorf("failed to connect portfolio service: %w", err)
	}
	defer portfolioConn.Close()

	portfolioServiceClientProto := portfoliopb.NewPortfolioServiceClient(portfolioConn)
	portfolioServiceClientContracted := portfolioGRPC.NewPortfolioServiceClient(portfolioServiceClientProto)

	portfolioServiceClient := portfolio.NewPortfolioServiceUsecase(portfolioServiceClientContracted)
	portfolioServiceController := portfolioController.NewPortfolioController(portfolioServiceClient)

	app := fiber.New()

	// Инициализируем метрики один раз
	metrics.InitMetrics()

	app.Use(middleware.LoggerMiddleware)
	app.Use(middleware.TraceMiddleware)
	app.Use(middleware.MetricsMiddleware)
	app.Use(rlMw.Handler)

	// Остальные эндпоинты
	app.Get("/limitertest", func(c *fiber.Ctx) error { return c.SendString("OK, not limited") })

	// Регистрируем только один раз и используем глобальный registry
	app.Get("/metrics", adaptor.HTTPHandler(
		promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{}),
	))
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Post("/portfolios", authMiddlewareVerifier.AuthVerify, portfolioServiceController.CreateNewPortfolio)
	app.Get("/portfolio/:id", authMiddlewareVerifier.AuthVerify, portfolioServiceController.GetPortfolioContentById)
	app.Post("/portfolio/:id/asset", authMiddlewareVerifier.AuthVerify, portfolioServiceController.UpsertAsset)
	app.Delete("/portfolio/:id/asset", authMiddlewareVerifier.AuthVerify, portfolioServiceController.DeleteAsset)
	app.Get("/portfolios", authMiddlewareVerifier.AuthVerify, portfolioServiceController.GetAllPortfolios)
	app.Get("/portfolio/:id/history", authMiddlewareVerifier.AuthVerify, portfolioServiceController.GetPortfolioHistory)
	app.Get("/portfolio/public/:username", portfolioServiceController.GetPublicPortfolios)

	log.Info("Starting API Gateway", zap.String("port", cfg.Port))
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Error("failed to start gateway", zap.Error(err))
		return fmt.Errorf("gateway failed: %w", err)
	}

	return nil
}
