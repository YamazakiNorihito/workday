package app_service

import (
	"context"

	"github.com/YamazakiNorihito/workday/internal/domain/metadata"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
)

func Execute(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, rssEntry rss.Rss) error {
	_, err := Write(ctx, logger, rssRepository, rssEntry)
	if err != nil {
		return err
	}

	logger.Info("Message Saved successfully", "feedURL", rssEntry.Source)
	return nil
}

func Write(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, rssEntry rss.Rss) (rss.Rss, error) {
	exists, existingRss := rss.Exists(ctx, rssRepository, rssEntry)
	logger.Info("Checking existence of RSS entry", "exists", exists, "source", rssEntry.Source)

	if exists && existingRss.LastBuildDate.Equal(rssEntry.LastBuildDate) {
		logger.Info("RSS entry is up-to-date, no update needed", "source", rssEntry.Source)
		return existingRss, nil
	}

	savedRss, err := rssRepository.Save(ctx, rssEntry, metadata.UserMeta{ID: rssEntry.Source, Name: rssEntry.Source})
	return savedRss, err
}
