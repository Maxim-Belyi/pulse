package usecase

import (
	"context"
	"errors"
	"fmt"
	"pulse/internal/entity"
	"pulse/pkg/logger"
)

type IngestionUseCase struct {
	logger    *logger.Logger
	publisher EventPublisher
}

func NewIngestionUseCase(log *logger.Logger, publisher EventPublisher) *IngestionUseCase {
	return &IngestionUseCase{
		logger:    log,
		publisher: publisher,
	}
}

func (i *IngestionUseCase) ProcessEvent(ctx context.Context, event *entity.Event) error {

	if event.ID == "" || event.Type == "" {
		return errors.New("пустой ID или Type")
	}

	err := i.publisher.Publish(ctx, event)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("событие отправлено (Title: %s, ID: %s, Source: %s, Type: %s)",event.Title, event.ID, event.Source, event.Type)
	i.logger.Info(msg)
	return nil
}
