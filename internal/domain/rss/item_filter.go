package rss

type ItemFilter struct {
	IncludeKeywords []string `json:"include_keywords"`
	ExcludeKeywords []string `json:"exclude_keywords"`
}

func NewItemFilter(includeKeywords, excludeKeywords []string) ItemFilter {
	if includeKeywords == nil {
		includeKeywords = []string{}
	}
	if excludeKeywords == nil {
		excludeKeywords = []string{}
	}
	return ItemFilter{
		IncludeKeywords: includeKeywords,
		ExcludeKeywords: excludeKeywords,
	}
}

func (f *ItemFilter) GetIncludeKeywords() []string {
	return f.IncludeKeywords
}

func (f *ItemFilter) GetExcludeKeywords() []string {
	return f.ExcludeKeywords
}
