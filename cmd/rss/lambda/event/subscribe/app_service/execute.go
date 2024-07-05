package app_service

import (
	"context"
	"fmt"

	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
	"github.com/mmcdole/gofeed"
)

type FeedParser interface {
	ParseURLWithContext(feedURL string, ctx context.Context) (feed *gofeed.Feed, err error)
}

func Execute(ctx context.Context, logger infrastructure.Logger, feedRepository *FeedRepository, publisher publisher.WriterMessagePublisher) error {
	entryRss, err := Subscribe(ctx, logger, feedRepository)
	if err != nil {
		return err
	}

	err = publisher.Publish(ctx, entryRss)
	if err != nil {
		return err
	}
	logger.Info("Message published successfully", "feedURL", feedRepository.FeedURL())
	return nil
}

func Subscribe(ctx context.Context, logger infrastructure.Logger, feedRepository *FeedRepository) (rssEntry rss.Rss, err error) {
	source := getFQDN(feedRepository.FeedURL())
	if source == "" {
		return rss.Rss{}, fmt.Errorf("invalid Feed URL: %s", feedRepository.FeedURL())
	}

	feed, err := feedRepository.GetFeed(ctx)
	if err != nil {
		logger.Error("Failed to retrieve RSS feed", "URL", feedRepository.FeedURL(), "error", err)
		return rss.Rss{}, err
	}

	lastBuildDate := getLastBuildDate(*feed)
	rssEntry, err = rss.New(feed.Title, source, feedRepository.FeedURL(), feed.Description, feed.Language, lastBuildDate.UTC())
	if err != nil {
		return rss.Rss{}, err
	}

	for _, item := range feed.Items {
		guid, err := getGuid(*item)
		if err != nil {
			logger.Error("Failed to create GUID from RSS item link", "error", err, "item", item.Title, "link", item.Link)
			continue
		}

		author := getAuthor(*item)
		textDescription, err := getDescription(*item)
		if err != nil {
			logger.Error("Failed to extract text from HTML content in RSS item description", "error", err, "itemTitle", item.Title)
			continue
		}

		entryItem, err := rss.NewItem(guid, item.Title, item.Link, textDescription, author, *item.PublishedParsed)
		if err != nil {
			logger.Error("Validation error when creating RSS item", "error", err, "item", item.Title)
			continue
		}
		rssEntry.AddOrUpdateItem(entryItem)
	}

	return rssEntry, nil
}
