package rss

import (
	"errors"
	"time"

	"github.com/YamazakiNorihito/workday/internal/domain/metadata"
	"github.com/google/uuid"
)

type Rss struct {
	ID            uuid.UUID         `json:"id"`
	Source        string            `json:"source"`
	Title         string            `json:"title"`
	Link          string            `json:"link"`
	Description   string            `json:"description"`
	Language      string            `json:"language"`
	LastBuildDate time.Time         `json:"last_build_date"`
	Items         map[Guid]Item     `json:"items"`
	ItemFilter    ItemFilter        `json:"item_filter"`
	CreatedBy     metadata.CreateBy `json:"create_by"`
	CreatedAt     metadata.CreateAt `json:"create_at"`
	UpdatedBy     metadata.UpdateBy `json:"update_by"`
	UpdatedAt     metadata.UpdateAt `json:"update_at"`
}

func New(title, source, link, description, language string, lastBuildDate time.Time) (Rss, error) {
	if title == "" || source == "" || link == "" || lastBuildDate.IsZero() {
		return Rss{}, errors.New("missing required fields: title, source, link, lastBuildDate must be provided")
	}

	return Rss{
		ID:            uuid.New(),
		Source:        source,
		Title:         title,
		Link:          link,
		Description:   description,
		Language:      language,
		LastBuildDate: lastBuildDate,
		Items:         make(map[Guid]Item),
		ItemFilter:    NewItemFilter(nil, nil),
	}, nil
}

func (r *Rss) SetLastBuildDate(lastBuildDate time.Time) error {
	if lastBuildDate.IsZero() {
		return errors.New("missing required fields: lastBuildDate must be provided")
	}
	r.LastBuildDate = lastBuildDate
	return nil
}

func (r *Rss) SetLanguage(language string) error {
	r.Language = language
	return nil
}

func (r *Rss) AddOrUpdateItem(item Item) {
	if r.ItemFilter.IsMatch(item) {
		r.Items[item.Guid] = item
	}
}

func (r *Rss) SetItemFilter(includeKeywords, excludeKeywords []string) {
	r.ItemFilter = NewItemFilter(includeKeywords, excludeKeywords)
}
