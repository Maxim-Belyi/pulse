package grpc

import (
	"context"
	pb "pulse/api/pb/event_v1"
	"pulse/internal/entity"
	"pulse/internal/usecase"
	"pulse/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventController struct {
	logger  *logger.Logger
	useCase *usecase.IngestionUseCase
}

func NewEventController(log *logger.Logger, useCase *usecase.IngestionUseCase) *EventController {
	return &EventController{
		logger:  log,
		useCase: useCase,
	}
}

func (e *EventController) PublishEvent(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {
	event := &entity.Event{
		ID:          req.Id,
		ExternalID:  req.ExternalId,
		Title:       req.Title,
		Source:      entity.SourceType(req.Source),
		Type:        entity.EventType(req.Type),
		Payload:     req.Payload,
		CollectedAt: req.CollectedAt.AsTime(),
		OccuredAt:   req.OccurredAt.AsTime(),
	}

	err := e.useCase.ProcessEvent(ctx, event)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.EventResponse{Id: event.ID}, nil
}
