package app_service

import (
	"context"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
)

func Execute(ctx context.Context, logger infrastructure.Logger, translator shared.Translator, publisher publisher.WriterMessagePublisher, rssEntry rss.Rss) error {
	translateRss, err := Translate(ctx, logger, translator, rssEntry)
	if err != nil {
		return err
	}

	err = publisher.Publish(ctx, translateRss)
	if err != nil {
		return err
	}
	logger.Info("Message published successfully", "feedURL", translateRss.Source)
	return nil
}

func Translate(ctx context.Context, logger infrastructure.Logger, translator shared.Translator, rssEntry rss.Rss) (rss.Rss, error) {
	if rssEntry.Language == "ja" || rssEntry.Language == "" {
		logger.Warn("翻訳対象外のためSkipします", "rssEntry.Source", rssEntry.Source)
		return rssEntry, nil
	}

	logger.Info("Source language found", "sourceLanguageCode", rssEntry.Language)
	for guid, item := range rssEntry.Items {
		if len(item.Description) > 0 {
			translatedText, err := translator.TranslateText(ctx, rssEntry.Language, "ja", item.Description)
			if err != nil {
				logger.Warn("変換に失敗しました。原文のまま処理します。", "item.Title", item.Title, "error", err)
				continue
			}
			logger.Info("Translation succeeded", "item.Title", item.Title, "before", item.Description, "after", translatedText)
			item.Description = translatedText
			rssEntry.Items[guid] = item
		}
	}
	return rssEntry, nil
}
