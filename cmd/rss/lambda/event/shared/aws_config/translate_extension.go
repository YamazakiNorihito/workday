package awsConfig

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/translate"
)

type EasyTranslateClient struct {
	client *translate.Client
}

func NewTranslateClient(client *translate.Client) *EasyTranslateClient {
	return &EasyTranslateClient{client}
}

func (c *EasyTranslateClient) TranslateText(ctx context.Context, sourceLanguageCode string, targetLanguageCode string, text string) (translatedText string, err error) {

	input := &translate.TranslateTextInput{
		SourceLanguageCode: aws.String(sourceLanguageCode),
		TargetLanguageCode: aws.String(targetLanguageCode),
		Text:               aws.String(text),
	}

	result, err := c.client.TranslateText(ctx, input)
	if err != nil {
		return "", err
	}

	translatedText = *result.TranslatedText
	return translatedText, nil
}
