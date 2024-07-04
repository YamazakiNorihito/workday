package app_service

import (
	"context"
	"net/http"

	"github.com/mmcdole/gofeed"
)

type FeedRepository struct {
	goParser *gofeed.Parser
	feedURL  string
}

func NewFeedRepository(httpClient *http.Client, feedURL string) FeedRepository {
	fp := gofeed.NewParser()
	fp.Client = httpClient

	return FeedRepository{goParser: fp, feedURL: feedURL}
}

func (r *FeedRepository) FeedURL() string {
	return r.feedURL
}

func (r *FeedRepository) GetFeed(ctx context.Context) (feed *gofeed.Feed, err error) {
	fp := r.goParser
	return fp.ParseURLWithContext(r.feedURL, ctx)
}
