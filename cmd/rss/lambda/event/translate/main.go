package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared"
	awsConfig "github.com/YamazakiNorihito/workday/cmd/rss/lambda/event/shared/aws_config"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/message"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.SNSEvent) error {
	cfg := awsConfig.LoadConfig(ctx)
	snsClient := cfg.NewSnsClient()
	translateClient := cfg.NewTranslateClient()

	snsTopicClient := awsConfig.NewSnsTopicClient(snsClient, os.Getenv("RSS_WRITE_ARN"))
	easyTranslateClient := awsConfig.NewTranslateClient(translateClient)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("SNS Event", "event", shared.SnsEventToJson(event))

	for _, record := range event.Records {
		recordLogger := logger.With("messageID", record.SNS.MessageID)
		err := processRecord(ctx, recordLogger, easyTranslateClient, snsTopicClient, record)

		if err != nil {
			recordLogger.Error("Failed", "error", err)
		}
		logger.Info("finish")
	}

	return nil
}

func processRecord(ctx context.Context, logger infrastructure.Logger, translator shared.Translator, rssWritePublisher shared.Publisher, record events.SNSEventRecord) error {
	receiveMessage, err := getMessage(record)
	if err != nil {
		return err
	}

	logger.Info("Processing command", receiveMessage.RssFeed.Source)
	entryRss, err := Core(ctx, logger, translator, receiveMessage.RssFeed)
	if err != nil {
		return err
	}

	err = Publish(ctx, rssWritePublisher, entryRss)
	if err != nil {
		return err
	}
	logger.Info("Message published successfully", "feedURL", receiveMessage.RssFeed.Source)
	return nil
}

func Core(ctx context.Context, logger infrastructure.Logger, translator shared.Translator, rssEntry rss.Rss) (rss.Rss, error) {
	languageFeedMap := map[string]string{
		"go.dev":         "en",
		"feed.infoq.com": "en",
		"techcrunch.com": "en",
	}

	sourceLanguageCode, ok := languageFeedMap[rssEntry.Source]
	if ok == false {
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

func getMessage(record events.SNSEventRecord) (receiveMessage message.Write, err error) {
	err = json.Unmarshal([]byte(record.SNS.Message), &receiveMessage)
	if err != nil {
		return message.Write{}, err
	}

	if receiveMessage.Compressed {
		decompressedRssData, err := message.DecodeAndDecompressData(receiveMessage.Data)
		if err != nil {
			return message.Write{}, err
		}
		receiveMessage.RssFeed = decompressedRssData
		receiveMessage.Data = nil
		receiveMessage.Compressed = false
	}

	return receiveMessage, nil
}

func main() {
	lambda.Start(handler)
}
