package app_service

import (
	"context"
	"time"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/shared/validation_error"
	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/shared/validator"
	"github.com/YamazakiNorihito/workday/internal/domain/metadata"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/google/uuid"
)

type GetCommand struct {
	Source string `validate:"required"`
}

type RssResponse struct {
	ID            uuid.UUID         `json:"id"`
	Source        string            `json:"source"`
	Title         string            `json:"title"`
	Link          string            `json:"link"`
	Description   string            `json:"description"`
	Language      string            `json:"language"`
	LastBuildDate time.Time         `json:"last_build_date"`
	ItemFilter    rss.ItemFilter    `json:"item_filter"`
	CreatedBy     metadata.CreateBy `json:"create_by"`
	CreatedAt     metadata.CreateAt `json:"create_at"`
	UpdatedBy     metadata.UpdateBy `json:"update_by"`
	UpdatedAt     metadata.UpdateAt `json:"update_at"`
}

func Execute(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, command GetCommand) (RssResponse, error) {
	rssFeed, err := GetRssFeed(ctx, logger, rssRepository, command)
	if err != nil {
		return RssResponse{}, err
	}

	logger.Info("Message FetchAllRssFeeds successfully")
	return rssFeed, nil
}

func GetRssFeed(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, command GetCommand) (RssResponse, error) {
	err := validator.Validate(ctx, command)

	if err != nil {
		return RssResponse{}, err
	}

	feed, err := rssRepository.FindBySource(ctx, command.Source)
	if err != nil {
		return RssResponse{}, err
	}

	if feed.ID == uuid.Nil {
		validationErr := validation_error.New(map[string]string{
			"source": "not found source: " + command.Source,
		})
		return RssResponse{}, validationErr
	}

	response := RssResponse{
		ID:            feed.ID,
		Source:        feed.Source,
		Title:         feed.Title,
		Link:          feed.Link,
		Description:   feed.Description,
		Language:      feed.Language,
		LastBuildDate: feed.LastBuildDate,
		ItemFilter:    feed.ItemFilter,
		CreatedBy:     feed.CreatedBy,
		CreatedAt:     feed.CreatedAt,
		UpdatedBy:     feed.UpdatedBy,
		UpdatedAt:     feed.UpdatedAt,
	}

	return response, nil
}
