package rss

import (
	"context"

	"github.com/google/uuid"
)

func Exists(ctx context.Context, repo IRssRepository, rss Rss) (bool, Rss) {
	targetRss, err := repo.FindBySource(ctx, rss.Source)

	if err == nil && targetRss.ID != uuid.Nil {
		return true, targetRss
	}
	return false, Rss{}
}

func GetItems(ctx context.Context, repo IRssRepository, rss Rss) (Rss, error) {
	targetRss, err := repo.FindItems(ctx, rss)
	return targetRss, err
}

func GetItem(ctx context.Context, repo IRssRepository, rss Rss, guid Guid) (Rss, error) {
	targetRss, err := repo.FindItemsByPk(ctx, rss, guid)
	return targetRss, err
}

func Delete(ctx context.Context, repo IRssRepository, rss Rss) error {
	return repo.Delete(ctx, rss)
}
