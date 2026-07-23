package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"pulse/internal/entity"
	"pulse/internal/usecase"
	"pulse/pkg/logger"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	logger    *logger.Logger
	ch        *amqp091.Channel
	queueName string
}

func NewConsumer(logger *logger.Logger, ch *amqp091.Channel, queueName string) *Consumer {
	return &Consumer{
		logger:    logger,
		ch:        ch,
		queueName: queueName,
	}
}

func (c *Consumer) Start(ctx context.Context, workersCount int, outchan chan<- usecase.ProcessingMessage) error {
	delivieries, err := c.ch.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("ошибка коньсюмера: %w", err)
	}
	var wg sync.WaitGroup
	for i := 0; i < workersCount; i++ {
		wg.Add(1)

		go func(workerId int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					c.logger.Error(fmt.Errorf("%v", r), fmt.Sprintf("воркер %d поймал панику", workerId))
				}
			}()

			for {
				select {
				case <-ctx.Done():
					return

				case msg, ok := <-delivieries:
					if !ok {
						return
					}

					event := &entity.Event{}
					err := json.Unmarshal(msg.Body, event)
					if err != nil {
						msg.Nack(false, false)
						text := fmt.Sprintf("ошибка парсинга ID: %s, RoutingKey: %s", msg.MessageId, msg.RoutingKey)
						c.logger.Error(err, text)
						continue
					}
					msgForUseCase := usecase.ProcessingMessage{
						Event: event,
						Ack:   func() error { return msg.Ack(false) },
						Nack:  func() error { return msg.Nack(false, true) },
					}
					select {
					case outchan <- msgForUseCase:
					case <-ctx.Done():
						return

					}
					text := fmt.Sprintf("консьюмер получил сообщение! ID: %s, RoutingKey: %s", msg.MessageId, msg.RoutingKey)
					c.logger.Info(text)
				}
			}
		}(i)
	}
	wg.Wait()
	return nil
}
