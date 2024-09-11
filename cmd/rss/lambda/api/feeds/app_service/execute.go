package app_service

import (
	"context"
	"time"

	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/google/uuid"
)

type RssResponse struct {
	ID            uuid.UUID `json:"id"`
	Source        string    `json:"source"`
	Title         string    `json:"title"`
	Link          string    `json:"link"`
	LastBuildDate time.Time `json:"lastBuildDate"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func Execute(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository) ([]RssResponse, error) {
	rssFeeds, err := AllRssFeeds(ctx, logger, rssRepository)
	if err != nil {
		return nil, err
	}

	logger.Info("Message FetchAllRssFeeds successfully")
	return rssFeeds, nil
}

func AllRssFeeds(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository) ([]RssResponse, error) {
	rssFeeds, err := rssRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]RssResponse, len(rssFeeds))
	for i, feed := range rssFeeds {
		response[i] = RssResponse{
			ID:            feed.ID,
			Source:        feed.Source,
			Title:         feed.Title,
			Link:          feed.Link,
			LastBuildDate: feed.LastBuildDate,
			CreatedAt:     feed.CreatedAt,
			UpdatedAt:     feed.UpdatedAt,
		}
	}

	return response, nil
}
