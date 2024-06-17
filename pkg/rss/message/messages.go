package message

import "github.com/YamazakiNorihito/workday/internal/domain/rss"

type Subscribe struct {
	FeedURL string `json:"feed_url"`
}

type Write struct {
	RssEntry rss.Rss `json:"rss"`
}
