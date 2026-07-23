package shutdown

import (
	"context"
	"sync"
	"time"

	"pulse/pkg/logger"
)

type Operation func(ctx context.Context) error

type GracefulShutdown struct {
	operations map[string]Operation
	log *logger.Logger
}

func New(log *logger.Logger) *GracefulShutdown {
	return &GracefulShutdown{
		operations: make(map[string]Operation),
		log: log,
	}
}

func (gs *GracefulShutdown) Add(name string, op Operation) {
	gs.operations[name] = op
}

func (gs *GracefulShutdown) Wait(ctx context.Context, timeout time.Duration) {
	<- ctx.Done()
	gs.log.Info("получен сигнал остановки, запуск Graceful Shutdown...")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup
	done := make(chan struct{})

	for name, op := range gs.operations {
		wg.Add(1)

		go func(n string, operation Operation){
			defer wg.Done()

			gs.log.Info("Начинаем закрытие: " + n)
			if err := operation(timeoutCtx); err != nil {
				gs.log.Error(err, "Ошибка при закрытии ресурса" + n)
			}
		}(name, op)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <- done:
		gs.log.Info("Все процессы завершены успешно")

	case <- timeoutCtx.Done():
		gs.log.Info("Таймаут Graseful Shutdown! Принудительное завершение")
	}
}
