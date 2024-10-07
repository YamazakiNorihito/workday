package app_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	errors map[string]string
}

func (ve *ValidationError) Error() string {
	var errMessages []string
	for field, message := range ve.errors {
		errMessages = append(errMessages, fmt.Sprintf("%s: %s", field, message))
	}
	return fmt.Sprintf("Validation failed: %s", strings.Join(errMessages, ", "))
}

func (ve *ValidationError) Errors() map[string]string {
	return ve.errors
}

func (c *DeleteCommand) Validation(ctx context.Context) error {
	validate := validator.New()
	errMap := make(map[string]string)

	if err := validate.StructCtx(ctx, c); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			tag := err.Tag()
			param := err.Param()

			var values string
			switch tag {
			case "min", "max":
				values = fmt.Sprintf("value must be %s %s", tag, param)
			case "oneof":
				values = fmt.Sprintf("value must be one of [%s]", strings.ReplaceAll(param, " ", ", "))
			default:
				if param == "" {
					values = "invalid value"
				} else {
					values = param
				}
			}

			message := fmt.Sprintf("%s is %s: %s", fieldName, tag, values)
			errMap[fieldName] = message
		}
		if len(errMap) > 0 {
			return &ValidationError{errors: errMap}
		}
	}
	return nil
}

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
	err := command.Validation(ctx)

	if err != nil {
		return err
	}

	return publisher.Publish(ctx, command.Source)
}
