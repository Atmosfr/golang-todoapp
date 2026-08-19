package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/Atmosfr/golang-todoapp/internal/core/config"
	core_logger "github.com/Atmosfr/golang-todoapp/internal/core/logger"
	core_pgx_pool "github.com/Atmosfr/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/Atmosfr/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/Atmosfr/golang-todoapp/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/Atmosfr/golang-todoapp/internal/features/statistics/repository/postgres"
	statistics_service "github.com/Atmosfr/golang-todoapp/internal/features/statistics/service"
	statistics_transport_http "github.com/Atmosfr/golang-todoapp/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/Atmosfr/golang-todoapp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/Atmosfr/golang-todoapp/internal/features/tasks/service"
	tasks_transport_http "github.com/Atmosfr/golang-todoapp/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/Atmosfr/golang-todoapp/internal/features/users/repository/postgres"
	users_service "github.com/Atmosfr/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/Atmosfr/golang-todoapp/internal/features/users/transport/http"
	web_fs_repository "github.com/Atmosfr/golang-todoapp/internal/features/web/repository/file_system"
	web_service "github.com/Atmosfr/golang-todoapp/internal/features/web/service"
	web_transport_http "github.com/Atmosfr/golang-todoapp/internal/features/web/transport/http"
	"go.uber.org/zap"

	_ "github.com/Atmosfr/golang-todoapp/docs"
)

// @title Golang Todo API
// @version 1.0
// @description Todo Application REST-API scheme
// @host 127.0.0.1:5050
// @BasePath /api/v1
func main() {
	cfg := core_config.NewConfigMust()

	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	loggerConfig := core_logger.NewLoggerConfigMust()
	logger, err := core_logger.NewLogger(loggerConfig)
	if err != nil {
		fmt.Println("failed to init app logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("app time zone", zap.Any("zone", time.Local))

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)

	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}

	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTaskHTTPHandler(tasksService)

	logger.Debug("initializing feature statistics", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("initializing feature web", zap.String("feature", "web"))
	webRepository := web_fs_repository.NewWebRepository()
	webService := web_service.NewWebService(webRepository)
	webTransportHTTP := web_transport_http.NewWebHTTPHandler(webService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.CORSMiddleware(),
		core_http_middleware.RequestIDMiddleware(),
		core_http_middleware.LoggerMiddleware(logger),
		core_http_middleware.TraceMiddleware(),
		core_http_middleware.RecoverMiddleware(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)
	httpServer.RegisterSwagger()
	httpServer.RegisterRoutes(webTransportHTTP.Routes()...)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server runtime error: %w", zap.Error(err))
	}
}
