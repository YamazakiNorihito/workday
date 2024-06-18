package rss

import (
	"errors"
	"time"
)

type Item struct {
	Guid        Guid      `json:"guid"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	PubDate     time.Time `json:"pubDate"`
	Tags        []string  `json:"tags"`
}

func NewItem(guid Guid, title, link, description, author string, pubDate time.Time) (Item, error) {
	if title == "" || link == "" || guid.Value == "" {
		return Item{}, errors.New("title, link, and guid cannot be empty")
	}

	return Item{
		Title:       title,
		Link:        link,
		Description: description,
		Author:      author,
		Guid:        guid,
		PubDate:     pubDate,
		Tags:        []string{},
	}, nil
}

func (i *Item) AddTag(tag string) {
	for _, t := range i.Tags {
		if t == tag {
			return
		}
	}
	i.Tags = append(i.Tags, tag)
}
