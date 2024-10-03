package rss

import "regexp"

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

func (f *ItemFilter) IsMatch(item Item) bool {
	if len(f.IncludeKeywords) > 0 {
		matched := false
		for _, pattern := range f.IncludeKeywords {
			re, err := regexp.Compile(pattern)
			if err != nil {
				continue
			}
			if re.MatchString(item.Title) || re.MatchString(item.Description) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	if len(f.ExcludeKeywords) > 0 {
		for _, pattern := range f.ExcludeKeywords {
			re, err := regexp.Compile(pattern)
			if err != nil {
				continue
			}
			if re.MatchString(item.Title) || re.MatchString(item.Description) {
				return false
			}
		}
	}

	return true
}

func (f *ItemFilter) Equal(other ItemFilter) bool {
	if len(f.IncludeKeywords) != len(other.IncludeKeywords) {
		return false
	}
	for i, keyword := range f.IncludeKeywords {
		if keyword != other.IncludeKeywords[i] {
			return false
		}
	}

	if len(f.ExcludeKeywords) != len(other.ExcludeKeywords) {
		return false
	}
	for i, keyword := range f.ExcludeKeywords {
		if keyword != other.ExcludeKeywords[i] {
			return false
		}
	}

	return true
}
