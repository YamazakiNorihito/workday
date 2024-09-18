package app_service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/YamazakiNorihito/workday/internal/domain/metadata"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type ValidationError struct {
	errors map[string]string
}

func (ve *ValidationError) Error() string {
	var errMessages []string
	for field, message := range ve.errors {
		errMessages = append(errMessages, fmt.Sprintf("%s: %s", field, message))
	}
	return fmt.Sprintf("Validation failed: %s", strings.Join(errMessages, ", "))
}

func (ve *ValidationError) Errors() map[string]string {
	return ve.errors
}

func (c *GetCommand) Validation(ctx context.Context) error {
	validate := validator.New()
	errMap := make(map[string]string)

	if err := validate.StructCtx(ctx, c); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			tag := err.Tag()
			param := err.Param()

			var values string
			switch tag {
			case "min", "max":
				values = fmt.Sprintf("value must be %s %s", tag, param)
			case "oneof":
				values = fmt.Sprintf("value must be one of [%s]", strings.ReplaceAll(param, " ", ", "))
			default:
				if param == "" {
					values = "invalid value"
				} else {
					values = param
				}
			}

			message := fmt.Sprintf("%s is %s: %s", fieldName, tag, values)
			errMap[fieldName] = message
		}
		if len(errMap) > 0 {
			return &ValidationError{errors: errMap}
		}
	}
	return nil
}

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
	err := command.Validation(ctx)

	if err != nil {
		return RssResponse{}, err
	}

	feed, err := rssRepository.FindBySource(ctx, command.Source)
	if err != nil {
		return RssResponse{}, err
	}

	if feed.ID == uuid.Nil {
		validationErr := ValidationError{
			errors: map[string]string{
				"source": "not found source: " + command.Source,
			},
		}
		return RssResponse{}, &validationErr
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
