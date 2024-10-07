package app_service

import (
	"context"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/shared/validation_error"
	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/shared/validator"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/message"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
	"github.com/google/uuid"
)

type PatchCommand struct {
	Source             string `validate:"required"`
	SourceLanguageCode string `validate:"required,oneof=af sq am ar hy az bn bs bg ca zh zh-TW hr cs da fa-AF nl en et fa tl fi fr fr-CA ka de el gu ht ha he hi hu is id ga it ja kn kk ko lv lt mk ms ml mt mr mn no ps pl pt pt-PT pa ro ru sr si sk sl so es es-MX sw sv ta te th tr uk ur uz vi cy"`
	ItemFilter         struct {
		IncludeKeywords []string
		ExcludeKeywords []string
	}
}

func Execute(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, publisher publisher.SubscribeMessagePublisher, command PatchCommand) error {
	err := Update(ctx, logger, rssRepository, publisher, command)
	if err != nil {
		return err
	}

	logger.Info("Message Trigger successfully")
	return nil
}

func Update(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, publisher publisher.SubscribeMessagePublisher, command PatchCommand) error {
	err := validator.Validate(ctx, command)

	if err != nil {
		return err
	}

	feed, err := rssRepository.FindBySource(ctx, command.Source)
	if err != nil {
		return err
	}

	if feed.ID == uuid.Nil {
		return validation_error.New(map[string]string{
			"source": "not found source: " + command.Source,
		})
	}

	message := message.Subscribe{
		FeedURL:    feed.Link,
		Language:   command.SourceLanguageCode,
		ItemFilter: rss.NewItemFilter(command.ItemFilter.IncludeKeywords, command.ItemFilter.ExcludeKeywords),
	}

	return publisher.Publish(ctx, message)
}
