package clickhouse

import (
	"fmt"
	"context"
	"pulse/internal/entity"
	"pulse/internal/usecase"
	"pulse/pkg/logger"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type EventRepository struct {
	logger *logger.Logger
	dbConn driver.Conn
}

func NewEventRepository(logger *logger.Logger, dbConn driver.Conn) *EventRepository {
	return &EventRepository{
		logger: logger,
		dbConn: dbConn,
	}
}

func (r *EventRepository) SaveBatch(ctx context.Context, events []*entity.Event) error {
	batch, err := r.dbConn.PrepareBatch(ctx, "INSERT INTO events")
	if err != nil {
		return fmt.Errorf("ошибка батча %w", err)
	}

	for _, event := range events {
		batch.Append(
			event.ID,
			string(event.Source),
			string(event.Type),
			event.ExternalID,
			event.Title,
			string(event.Payload),
			event.CollectedAt,
			event.OccurredAt,
		)
		if err != nil {
			return fmt.Errorf("ошибка добавления события %s в батч %w", event.ID, err)
		}
	}
	return batch.Send()
}

var _ usecase.EventRepository = (*EventRepository)(nil)
