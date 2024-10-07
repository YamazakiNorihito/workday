package app_service

import (
	"context"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/shared/validator"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/message"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
)

type CreateCommand struct {
	FeedURL            string `validate:"required,url,startswith=http"`
	SourceLanguageCode string `validate:"required,oneof=af sq am ar hy az bn bs bg ca zh zh-TW hr cs da fa-AF nl en et fa tl fi fr fr-CA ka de el gu ht ha he hi hu is id ga it ja kn kk ko lv lt mk ms ml mt mr mn no ps pl pt pt-PT pa ro ru sr si sk sl so es es-MX sw sv ta te th tr uk ur uz vi cy"`
	ItemFilter         struct {
		IncludeKeywords []string `json:"include_keywords"`
		ExcludeKeywords []string `json:"exclude_keywords"`
	} `json:"item_filter"`
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
	err := validator.Validate(ctx, command)

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
