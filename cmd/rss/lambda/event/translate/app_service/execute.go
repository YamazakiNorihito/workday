package app_service

import (
	"context"
	"encoding/json"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/message"
)

type FeedLanguageProvider interface {
	GetSourceLanguage(source string) (string, bool)
}

func Execute(ctx context.Context, logger infrastructure.Logger, translator shared.Translator, feedLanguageProvider FeedLanguageProvider, rssWritePublisher shared.Publisher, rssEntry rss.Rss) error {
	translateRss, err := Translate(ctx, logger, translator, feedLanguageProvider, rssEntry)
	if err != nil {
		return err
	}

	err = Publish(ctx, rssWritePublisher, translateRss)
	if err != nil {
		return err
	}
	logger.Info("Message published successfully", "feedURL", translateRss.Source)
	return nil
}

func Translate(ctx context.Context, logger infrastructure.Logger, translator shared.Translator, feedLanguageProvider FeedLanguageProvider, rssEntry rss.Rss) (rss.Rss, error) {
	sourceLanguageCode, ok := feedLanguageProvider.GetSourceLanguage(rssEntry.Source)
	if ok == false || sourceLanguageCode == "ja" {
		logger.Warn("翻訳対象外のためSkipします", "rssEntry.Source", rssEntry.Source)
		return rssEntry, nil
	}

	logger.Info("Source language found", "sourceLanguageCode", sourceLanguageCode)
	for guid, item := range rssEntry.Items {
		if len(item.Description) > 0 {
			translatedText, err := translator.TranslateText(ctx, sourceLanguageCode, "ja", item.Description)
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

func Publish(ctx context.Context, rssWritePublisher shared.Publisher, rssEntry rss.Rss) error {
	writeMessage, err := message.NewWriteMessage(rssEntry)
	if err != nil {
		return err
	}

	rssJson, _ := json.Marshal(writeMessage)
	err = rssWritePublisher.Publish(ctx, string(rssJson))
	if err != nil {
		return err
	}
	return nil
}
