package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/Atmosfr/golang-todoapp/internal/core/logger"
	core_pgx_pool "github.com/Atmosfr/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/Atmosfr/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/Atmosfr/golang-todoapp/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/Atmosfr/golang-todoapp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/Atmosfr/golang-todoapp/internal/features/tasks/service"
	tasks_transport_http "github.com/Atmosfr/golang-todoapp/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/Atmosfr/golang-todoapp/internal/features/users/repository/postgres"
	users_service "github.com/Atmosfr/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/Atmosfr/golang-todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

var timeZone = time.UTC

func main() {
	time.Local = timeZone

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	loggerConfig := core_logger.NewLoggerConfigMust()
	logger, err := core_logger.NewLogger(loggerConfig)
	if err != nil {
		fmt.Println("failed to init app logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("app time zone", zap.Any("zone", timeZone))

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
	logger.Debug("initializing HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestIDMiddleware(),
		core_http_middleware.LoggerMiddleware(logger),
		core_http_middleware.TraceMiddleware(),
		core_http_middleware.RecoverMiddleware(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server runtime error: %w", zap.Error(err))
	}
}
