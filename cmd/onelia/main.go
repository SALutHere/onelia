package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/SALutHere/onelia/internal/core/config"
	core_logger "github.com/SALutHere/onelia/internal/core/logger"
	core_pgx_pool "github.com/SALutHere/onelia/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/SALutHere/onelia/internal/core/transport/http/middleware"
	core_http_server "github.com/SALutHere/onelia/internal/core/transport/http/server"
	routes_calculator "github.com/SALutHere/onelia/internal/routes/calculator"
	routes_postgres_repository "github.com/SALutHere/onelia/internal/routes/repository/postgres"
	routes_service "github.com/SALutHere/onelia/internal/routes/service"
	routes_transport_http "github.com/SALutHere/onelia/internal/routes/transport/http"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		return err
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "routes"))
	routesRepository := routes_postgres_repository.NewRoutesRepository(pool)
	routesCalculator, err := routes_calculator.NewRoutesCalculator(ctx, routesRepository)
	if err != nil {
		return err
	}
	routesService := routes_service.NewRoutesService(routesCalculator)
	routesTransportHTTP := routes_transport_http.NewRoutesHTTPHandler(routesService)

	logger.Debug("initializing HTTP server")
	httpConfig := core_http_server.NewConfigMust()
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
		core_http_middleware.CORS(httpConfig.AllowedOrigins),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(routesTransportHTTP.Routes()...)
	httpServer.RegisterAPIRouters(
		apiVersionRouterV1,
	)
	if err := httpServer.Run(ctx); err != nil {
		logger.Fatal("HTTP server run error", zap.Error(err))
	}

	return nil
}
