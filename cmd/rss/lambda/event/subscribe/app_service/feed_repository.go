package app_service

import (
	"context"
	"net/http"

	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/mmcdole/gofeed"
)

type FeedRepository struct {
	goParser   *gofeed.Parser
	feedURL    string
	language   string
	itemFilter rss.ItemFilter
}

func NewFeedRepository(httpClient *http.Client, feedURL, language string, itemFilter rss.ItemFilter) FeedRepository {
	fp := gofeed.NewParser()
	fp.Client = httpClient

	return FeedRepository{goParser: fp, feedURL: feedURL, language: language, itemFilter: itemFilter}
}

func (r *FeedRepository) FeedURL() string {
	return r.feedURL
}

func (r *FeedRepository) Language() string {
	return r.language
}

func (r *FeedRepository) ItemFilter() rss.ItemFilter {
	return r.itemFilter
}

func (r *FeedRepository) GetFeed(ctx context.Context) (feed *gofeed.Feed, err error) {
	fp := r.goParser
	return fp.ParseURLWithContext(r.feedURL, ctx)
}
