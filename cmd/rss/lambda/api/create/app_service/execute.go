package app_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/message"
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

type CreateCommand struct {
	FeedURL            string `validate:"required,url,startswith=http"`
	SourceLanguageCode string `validate:"required,oneof=af sq am ar hy az bn bs bg ca zh zh-TW hr cs da fa-AF nl en et fa tl fi fr fr-CA ka de el gu ht ha he hi hu is id ga it ja kn kk ko lv lt mk ms ml mt mr mn no ps pl pt pt-PT pa ro ru sr si sk sl so es es-MX sw sv ta te th tr uk ur uz vi cy"`
	ItemFilter         struct {
		IncludeKeywords []string `json:"include_keywords"`
		ExcludeKeywords []string `json:"exclude_keywords"`
	} `json:"item_filter"`
}

func (c *CreateCommand) Validation(ctx context.Context) error {
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

func Execute(ctx context.Context, logger infrastructure.Logger, publisher publisher.SubscribeMessagePublisher, command CreateCommand) error {
	err := Trigger(ctx, logger, publisher, command)
	if err != nil {
		return err
	}

	logger.Info("Message Trigger successfully")
	return nil
}

func Trigger(ctx context.Context, logger infrastructure.Logger, publisher publisher.SubscribeMessagePublisher, command CreateCommand) error {
	err := command.Validation(ctx)

	if err != nil {
		return err
	}

	message := message.Subscribe{
		FeedURL:    command.FeedURL,
		Language:   command.SourceLanguageCode,
		ItemFilter: rss.NewItemFilter(command.ItemFilter.IncludeKeywords, command.ItemFilter.ExcludeKeywords),
	}

	return publisher.Publish(ctx, message)
}
