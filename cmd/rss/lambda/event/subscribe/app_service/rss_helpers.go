package app_service

import (
	"net/url"
	"time"

	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/mmcdole/gofeed"
)

func getFQDN(uri string) string {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	return parsedURL.Host
}

func getLastBuildDate(feed gofeed.Feed) (lastBuildDate time.Time) {
	for _, item := range feed.Items {
		if item.PublishedParsed != nil && lastBuildDate.Before(*item.PublishedParsed) {
			lastBuildDate = *item.PublishedParsed
		}
	}

	if lastBuildDate.IsZero() && feed.UpdatedParsed != nil {
		lastBuildDate = *feed.UpdatedParsed
	}
	return lastBuildDate
}

func getGuid(item gofeed.Item) (rss.Guid, error) {
	guid := rss.Guid{Value: item.GUID}

	if guid.Value == "" {
		link, err := url.Parse(item.Link)
		if err != nil {
			return rss.Guid{}, err
		}
		link.RawQuery = ""
		guid = rss.Guid{Value: link.String()}
	}

	return guid, nil
}

func getAuthor(item gofeed.Item) string {
	if item.Author == nil {
		return ""
	}

	if email := item.Author.Email; email != "" {
		return email
	}
	return item.Author.Name
}

func getDescription(item gofeed.Item) (string, error) {
	description := item.Description
	if description == "" {
		description = item.Content
	}
	return extractTextFromHTML(description)
}
