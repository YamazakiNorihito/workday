package message

import "github.com/YamazakiNorihito/workday/internal/domain/rss"

const MaxMessageSize = 256 * 1024

type Subscribe struct {
	FeedURL string `json:"feed_url"`
}

type Write struct {
	RssFeed    rss.Rss `json:"rss,omitempty"`
	Compressed bool    `json:"compressed"`
	Data       []byte  `json:"data,omitempty"`
}
