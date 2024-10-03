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

	if !shouldUpdateRssEntry(existingRss, rssEntry) {
		logger.Info("RSS entry is up-to-date, no update needed", "source", rssEntry.Source)
		return existingRss, nil
	}

	savedRss, err := rssRepository.Save(ctx, rssEntry, metadata.UserMeta{ID: rssEntry.Source, Name: rssEntry.Source})
	return savedRss, err
}

func shouldUpdateRssEntry(existingRss rss.Rss, newRss rss.Rss) bool {
	if !existingRss.LastBuildDate.Equal(newRss.LastBuildDate) {
		return true
	}

	if !existingRss.ItemFilter.Equal(newRss.ItemFilter) {
		return true
	}

	return false
}
