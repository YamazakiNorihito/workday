package helper

import (
	"context"

	"github.com/YamazakiNorihito/workday/internal/domain/metadata"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
)

type SpyRssRepository struct {
	FindBySourceFunc  func(ctx context.Context, source string) (rss.Rss, error)
	FindAllFunc       func(ctx context.Context) ([]rss.Rss, error)
	FindItemsFunc     func(ctx context.Context, rss rss.Rss) (rss.Rss, error)
	FindItemsByPkFunc func(ctx context.Context, rss rss.Rss, guid rss.Guid) (rss.Rss, error)
	SaveFunc          func(ctx context.Context, rss rss.Rss, updateBy metadata.UserMeta) (rss.Rss, error)
	DeleteFunc        func(ctx context.Context, rss rss.Rss) error
}

func (r *SpyRssRepository) FindBySource(ctx context.Context, source string) (rss.Rss, error) {
	if r.FindBySourceFunc != nil {
		return r.FindBySourceFunc(ctx, source)
	}
	panic("FindBySourceFunc is not implemented")
}

func (r *SpyRssRepository) FindAll(ctx context.Context) ([]rss.Rss, error) {
	if r.FindAllFunc != nil {
		return r.FindAllFunc(ctx)
	}
	panic("FindAllFunc is not implemented")
}

func (r *SpyRssRepository) FindItems(ctx context.Context, rss rss.Rss) (rss.Rss, error) {
	if r.FindItemsFunc != nil {
		return r.FindItemsFunc(ctx, rss)
	}
	panic("FindItemsFunc is not implemented")
}

func (r *SpyRssRepository) FindItemsByPk(ctx context.Context, rss rss.Rss, guid rss.Guid) (rss.Rss, error) {
	if r.FindItemsByPkFunc != nil {
		return r.FindItemsByPkFunc(ctx, rss, guid)
	}
	panic("FindItemsByPkFunc is not implemented")
}

func (r *SpyRssRepository) Save(ctx context.Context, rss rss.Rss, updateBy metadata.UserMeta) (rss.Rss, error) {
	if r.SaveFunc != nil {
		return r.SaveFunc(ctx, rss, updateBy)
	}
	panic("SaveFunc is not implemented")
}

func (r *SpyRssRepository) Delete(ctx context.Context, rss rss.Rss) error {
	if r.DeleteFunc != nil {
		return r.DeleteFunc(ctx, rss)
	}
	panic("DeleteFunc is not implemented")
}
