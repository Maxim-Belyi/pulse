package rabbitmq

import (
	"context"
	"pulse/pkg/logger"

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

func (c *Consumer) Start(ctx context.Context, workersCount int) error {
	c.ch.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}
