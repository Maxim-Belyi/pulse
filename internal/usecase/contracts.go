package usecase

import (
	"context"
	"pulse/internal/entity"
	"time"
)
type ProcessingMessage struct {
	Event *entity.Event
	Ack   func() error
	Nack  func() error
}

type TrendStat struct {
	HourBacket time.Time
	TotalEvents uint64
}

type SourceStats struct {
	Source string
	TotalEvents uint64
}

type EventPublisher interface {
	Publish(ctx context.Context, event *entity.Event) (error)
}

type EventRepository interface {
	SaveBatch(ctx context.Context, events []*entity.Event) (error)
}

type CacheRepository interface {
	IncSourceCount(ctx context.Context, source entity.SourceType) (error)
}

type EventSource interface {
	Fetch(ctx context.Context) ([]*entity.Event, error)
}

type AnalyticsRepository interface {
	GetSources(ctx context.Context, since time.Time) ([]SourceStats, error)
	GetHourlyTrends(ctx context.Context, since time.Time) ([]TrendStat, error)
}