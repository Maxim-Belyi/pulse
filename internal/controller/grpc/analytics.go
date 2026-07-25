package grpc

import (
	"context"
	"pulse/internal/usecase"
	"pulse/pkg/logger"
)

type AnalyticsController struct {
	logger *logger.Logger
	useCase *usecase.AnalyticsUseCase
}

func NewAnalyticsController(logger *logger.Logger, useCase *usecase.AnalyticsUseCase) *AnalyticsController {
	return &AnalyticsController {
		logger: logger,
		useCase: useCase,
	}
}

func (a *AnalyticsController) GetTopSources(ctx context.Context, req *pb.Get)