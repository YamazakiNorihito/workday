package awsConfig

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type EasyTranslateClient struct {
	url string
}

func NewTranslateClient(url string) *EasyTranslateClient {
	return &EasyTranslateClient{url: url}
}

type translateRequest struct {
	SourceLanguageCode string `json:"sourceLanguageCode"`
	TargetLanguageCode string `json:"targetLanguageCode"`
	Text               string `json:"text"`
}

type translateResponse struct {
	Code  int    `json:"code"`
	Text  string `json:"text"`
	Error string `json:"error"`
}

func (c *EasyTranslateClient) TranslateText(ctx context.Context, sourceLanguageCode string, targetLanguageCode string, text string) (translatedText string, err error) {
	reqBody := translateRequest{
		SourceLanguageCode: sourceLanguageCode,
		TargetLanguageCode: targetLanguageCode,
		Text:               text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var translateResp translateResponse
	if err := json.NewDecoder(resp.Body).Decode(&translateResp); err != nil {
		return "", err
	}

	if translateResp.Code != 200 {
		return "", errors.New(translateResp.Error)
	}

	return translateResp.Text, nil
}
