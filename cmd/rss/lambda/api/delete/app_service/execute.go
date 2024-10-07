package app_service

import (
	"context"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/shared/validator"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
)

type DeleteCommand struct {
	Source string `validate:"required"`
}

func Execute(ctx context.Context, logger infrastructure.Logger, publisher publisher.DeleteMessagePublisher, command DeleteCommand) error {
	err := Delete(ctx, logger, publisher, command)
	if err != nil {
		return err
	}

	logger.Info("Message Delete successfully")
	return nil
}

func Delete(ctx context.Context, logger infrastructure.Logger, publisher publisher.DeleteMessagePublisher, command DeleteCommand) error {
	err := validator.Validate(ctx, command)

	if err != nil {
		return err
	}

	return publisher.Publish(ctx, command.Source)
}
