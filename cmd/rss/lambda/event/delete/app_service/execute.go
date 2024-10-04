package app_service

import (
	"context"

	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/internal/infrastructure"
)

func Execute(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, source string) error {
	err := Delete(ctx, logger, rssRepository, source)
	if err != nil {
		return err
	}

	logger.Info("RSS entry deleted successfully", "source", source)
	return nil
}

func Delete(ctx context.Context, logger infrastructure.Logger, rssRepository rss.IRssRepository, source string) error {
	rssEntry, err := rssRepository.FindBySource(ctx, source)
	if err != nil {
		logger.Error("Failed to retrieve RSS entry", "error", err, "source", source)
		return err
	}

	err = rssRepository.Delete(ctx, rssEntry)
	if err != nil {
		logger.Error("Failed to delete RSS entry", "error", err, "source", source)
		return err
	}
	return nil
}
