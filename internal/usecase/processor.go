package usecase

import (
	"context"
	"pulse/internal/entity"
	"pulse/pkg/logger"
	"time"
)



type Processor struct {
	logger        *logger.Logger
	repo          EventRepository
	batchSize     int
	flushInterval time.Duration
	cache         CacheRepository
}

func NewProcessor(logger *logger.Logger, repo EventRepository, batchSize int, flushInterval time.Duration, cache CacheRepository) *Processor {
	return &Processor{
		logger:        logger,
		repo:          repo,
		batchSize:     batchSize,
		flushInterval: flushInterval,
		cache:         cache,
	}
}

func (p *Processor) flush(ctx context.Context, batch []ProcessingMessage) error {
	eventsToSave := make([]*entity.Event, 0, len(batch))

	for _, item := range batch {
		eventsToSave = append(eventsToSave, item.Event)
	}
	err := p.repo.SaveBatch(ctx, eventsToSave)
	if err != nil {
		for _, item := range batch {
			item.Nack()
		}
		return err
	}
	for _, item := range batch {
		p.cache.IncSourceCount(ctx, item.Event.Source)
		item.Ack()
	}
	return err
}

func (p *Processor) Start(ctx context.Context, inChan <-chan ProcessingMessage) error {
	ticker := time.NewTicker(p.flushInterval)
	defer ticker.Stop()
	var batch []ProcessingMessage
	for {
		select {
		case <-ctx.Done():
			{
				if len(batch) > 0 {
					p.flush(ctx, batch)
				}
				return ctx.Err()
			}
		case <-ticker.C:
			if len(batch) > 0 {
				p.flush(ctx, batch)
				batch = batch[:0]
			}
		case delivery, ok := <-inChan:
			if !ok {
				if len(batch) > 0 {
					p.flush(ctx, batch)
					return nil
				}
				return ctx.Err()
			}
			batch = append(batch, delivery)
			if len(batch) >= p.batchSize {
				p.flush(ctx, batch)
				batch = batch[:0]
			}
		}

	}
}
