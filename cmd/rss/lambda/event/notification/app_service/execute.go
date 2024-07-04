package app_service

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
)

type RssConditions struct {
	Target     func(rss.Rss) bool
	ItemFilter func(item rss.Item) bool
}

type SlackSender interface {
	PostMessageContext(ctx context.Context, text string, username string) (respChannel string, respTimestamp string, err error)
}

func Execute(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, slackSender SlackSender, rssConditions RssConditions, source string) error {
	err := Notification(ctx, logger, rssRepository, slackSender, rssConditions, source)
	if err != nil {
		return err
	}

	logger.Info("Message published successfully", "feedURL", source)
	return nil
}

func Notification(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, slackSender SlackSender, rssConditions RssConditions, source string) error {
	newImageRss, err := rssRepository.FindBySource(ctx, source)
	if err != nil {
		return err
	}

	if rssConditions.Target(newImageRss) == false {
		return nil
	}

	modifyRss, err := rss.GetItems(ctx, rssRepository, newImageRss)
	if err != nil {
		return err
	}

	postMessage := makeMessage(modifyRss, rssConditions.ItemFilter)

	if postMessage == "" {
		logger.Info("No message to send to Slack", "source", source)
		return nil
	}

	respChannel, respTimestamp, err := slackSender.PostMessageContext(ctx, postMessage, source)
	if err != nil {
		return err
	}

	logger.Info("Successfully sent message to Slack", "source", source, "response channel", respChannel, "timestamp", respTimestamp)
	return nil
}

func makeMessage(r rss.Rss, itemFilter func(item rss.Item) bool) string {
	filteredItems := filterMap(r.Items, itemFilter)
	if len(filteredItems) == 0 {
		return ""
	}

	var messageBuilder strings.Builder
	messageBuilder.WriteString(fmt.Sprintf("*フィードタイトル:* <%s|%s>\n*フィード詳細:* %s", r.Link, r.Title, r.Description))
	if !r.LastBuildDate.IsZero() {
		messageBuilder.WriteString(fmt.Sprintf("\n*最終更新日:* %s", r.LastBuildDate.Format(time.RFC3339)))
	}
	messageBuilder.WriteString("\n\n*最新の記事:*\n")

	keys := make([]string, 0, len(filteredItems))
	for k := range filteredItems {
		keys = append(keys, k.Value)
	}
	sort.Strings(keys)

	i := 1
	for _, key := range keys {
		item := filteredItems[rss.Guid{Value: key}]
		truncatedDescription := truncate(item.Description)
		messageBuilder.WriteString(fmt.Sprintf("%d. *記事タイトル:* <%s|%s>\n    *公開日:* %s\n    *概要:* %s\n    *カテゴリ:* %s\n\n",
			i, item.Link, item.Title, item.PubDate.Format(time.RFC3339), truncatedDescription, strings.Join(item.Tags, ", ")))
		i++
	}

	return messageBuilder.String()
}

func filterMap[K comparable, V any](m map[K]V, filterFunc func(V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if filterFunc(v) {
			result[k] = v
		}
	}
	return result
}

func truncate(s string) string {
	const viewLength = 50 * 4

	re := regexp.MustCompile(`\s+|\n|\r|\t`)
	trimmedString := re.ReplaceAllString(s, " ")

	runes := []rune(trimmedString)
	if len(runes) > viewLength {
		return string(runes[:viewLength]) + "..."
	}
	return string(runes)
}
