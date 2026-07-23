package usecase

import (
	"context"
	"pulse/internal/entity"
	"pulse/pkg/logger"
	"sync"
	"time"
)

type Collector struct {
	logger  *logger.Logger
	sources []EventSource
}

func NewCollector(log *logger.Logger, sources ...EventSource) *Collector {
	return &Collector{
		logger:  log,
		sources: sources,
	}
}

func (c *Collector) Start(ctx context.Context, interval time.Duration) <-chan *entity.Event {
	out := make(chan *entity.Event, 100)
	var wg sync.WaitGroup

	for _, source := range c.sources {
		wg.Add(1)

		go func(e EventSource) {
			defer wg.Done()

			ticker := time.NewTicker(interval)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					events, err := e.Fetch(ctx)
					if err != nil {
						c.logger.Error(err, "ошибка при сборе данных")
						continue
					}
					for _, event := range events {
						select {
						case out <- event:
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}(source)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
