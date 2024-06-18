package rss

import (
	"context"

	"github.com/google/uuid"
)

func Exists(ctx context.Context, repo IRssRepository, rss Rss) (bool, Rss) {
	targetRss, err := repo.FindBySource(ctx, rss.Source)

	if err == nil || targetRss.ID != uuid.Nil {
		return true, targetRss
	}
	return false, Rss{}
}
