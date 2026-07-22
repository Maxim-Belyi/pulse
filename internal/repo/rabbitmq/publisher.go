package rabbitmq

import (
	"context"
	"encoding/json"
	"pulse/internal/entity"
	"pulse/pkg/logger"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	logger *logger.Logger
	ch *amqp091.Channel
	exchangeName string
}

func NewPublisher(log *logger.Logger, ch *amqp091.Channel, exchangeName string) *Publisher {
	return &Publisher{
		logger: log,
		ch: ch,
		exchangeName: exchangeName,
	}
}

func (p *Publisher) Publish(ctx context.Context, event *entity.Event) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		p.logger.Error(err, "ошибка сериализации")
		return err
	}

	err = p.ch.PublishWithContext(
		ctx,
		p.exchangeName,
		"",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			DeliveryMode: amqp091.Persistent,
			Body: eventBytes,
		},
	)
	if err != nil {
		p.logger.Error(err, "ошибка публикации")
		return fmt.Errorf("ошибка публикации: %w", err)
	}
	p.logger.Info("сообщение отправлено")
	return fmt.Errorf("amqp: сообщение отправлено: %w", err)
}