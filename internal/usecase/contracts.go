package usecase

import (
	"context"
	"pulse/internal/entity"
)

type EventPublisher interface {
	Publish(ctx context.Context, event *entity.Event) (error)
}

type EventRepository interface {
	SaveBatch(ctx context.Context, events []*entity.Event) (error)
}

type CacheRepository interface {
	IncSourceCount(ctx context.Context, source entity.SourceType) (error)
}