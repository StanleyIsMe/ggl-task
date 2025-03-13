package api

import (
	"context"
	"time"

	taskHTTP "ggltask/internal/task/delivery/http"
	taskRepo "ggltask/internal/task/repository/memory"
	taskUseCase "ggltask/internal/task/usecase"

	pkgMiddleware "ggltask/pkg/transport/middleware"
)

const (
	defaultTimeout = 10 * time.Second
)

func (a *API) registerHTTPSvc(_ context.Context) {
	a.server.SetupHTTPServer()
	httpRouter := a.server.HTTPRouter()

	taskRepository := taskRepo.NewTaskRepository()

	taskUseCase := taskUseCase.NewTaskUseCaseImpl(taskRepository)

	httpRouter.Use(
		pkgMiddleware.GinRecover(),
		pkgMiddleware.GinContextLogger(a.logger), //nolint:contextcheck
		pkgMiddleware.GinTimeout(defaultTimeout), //nolint:contextcheck
	)

	taskHTTP.RegisterTaskRoutes(httpRouter, taskUseCase)
}
