package shared

import "context"

type Translator interface {
	TranslateText(ctx context.Context, sourceLanguageCode string, targetLanguageCode string, text string) (translatedText string, err error)
}
