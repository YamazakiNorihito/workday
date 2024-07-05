package app_service

import (
	"context"

	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
)

func Execute(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, publisher publisher.WriterMessagePublisher, rssEntry rss.Rss) error {
	cleansingRss, err := Clean(ctx, logger, rssRepository, rssEntry)
	if err != nil {
		return err
	}

	err = publisher.Publish(ctx, cleansingRss)
	if err != nil {
		return err
	}
	logger.Info("Message published successfully", "feedURL", cleansingRss.Source)
	return nil
}

func Clean(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, rssEntry rss.Rss) (cleansingRss rss.Rss, err error) {
	exists, existingRss := rss.Exists(ctx, rssRepository, rssEntry)
	logger.Info("Checking existence of RSS entry", "exists", exists, "source", rssEntry.Source)

	if exists == false {
		return rssEntry, nil
	}

	existingRss.SetLastBuildDate(rssEntry.LastBuildDate)
	for _, item := range rssEntry.Items {
		existingRss.AddOrUpdateItem(item)
	}

	cleansingRss = existingRss
	cleansingRss.Items = map[rss.Guid]rss.Item{}

	for key, item := range rssEntry.Items {
		findItem, err := rss.GetItem(ctx, rssRepository, rssEntry, key)
		if err != nil {
			logger.Error("Error retrieving item", "error", err, "source", rssEntry.Source, "guid", key)
			continue
		}

		if len(findItem.Items) == 0 {
			cleansingRss.Items[key] = item
		} else {
			logger.Info("Item already exists and will not be added", "source", rssEntry.Source, "guid", key)
		}
	}

	return cleansingRss, nil
}
