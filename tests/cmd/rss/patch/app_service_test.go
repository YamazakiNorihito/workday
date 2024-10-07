package create

import (
	"context"
	"testing"
	"time"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/patch/app_service"
	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/shared/validation_error"
	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/shared/validator"
	"github.com/YamazakiNorihito/workday/internal/domain/rss"
	"github.com/YamazakiNorihito/workday/pkg/rss/publisher"
	"github.com/YamazakiNorihito/workday/tests/helper"
	"github.com/stretchr/testify/assert"
)

type spyMessageClient struct{ Messages []string }

func (r *spyMessageClient) Publish(ctx context.Context, message string) error {
	r.Messages = append(r.Messages, message)
	return nil
}

func TestPatchCommand_Validation(t *testing.T) {
	t.Run("Should Pass All Validations", func(t *testing.T) {
		testCases := []struct {
			name    string
			command app_service.PatchCommand
		}{
			{
				name: "Valid Source and Language Code",
				command: app_service.PatchCommand{
					Source:             "Some Source",
					SourceLanguageCode: "en",
				},
			},
		}

		ctx := context.Background()
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Act
				err := validator.Validate(ctx, tc.command)
				// Assert
				assert.NoError(t, err)
			})
		}
	})

	t.Run("Should Fail Validations Where Appropriate", func(t *testing.T) {
		testCases := []struct {
			name    string
			command app_service.PatchCommand
		}{
			{
				name: "Missing Source",
				command: app_service.PatchCommand{
					SourceLanguageCode: "en",
				},
			},
			{
				name: "Invalid Language Code",
				command: app_service.PatchCommand{
					Source:             "Valid Source",
					SourceLanguageCode: "fake-lang",
				},
			},
			{
				name: "Empty Language Code",
				command: app_service.PatchCommand{
					Source:             "Valid Source",
					SourceLanguageCode: "",
				},
			},
		}

		ctx := context.Background()
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Act
				err := validator.Validate(ctx, tc.command)
				// Assert
				assert.Error(t, err)
				if ve, ok := err.(*validation_error.ValidationError); ok {
					assert.Greater(t, len(ve.Errors()), 0)
				}
			})
		}
	})
}

func TestAppService_Update(t *testing.T) {
	t.Run("should handle event successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		logger := helper.MockLogger{}
		repo := helper.SpyRssRepository{
			FindBySourceFunc: func(ctx context.Context, source string) (rss.Rss, error) {
				dummy_rss1, err := rss.New(
					"ダミーニュースのフィード1",
					"connpass.com",
					"https://connpass.com/explore/ja.atom",
					"このフィードはダミーニュース1を提供します。",
					"ja",
					time.Date(2024, time.July, 3, 13, 0, 0, 0, time.UTC),
				)
				return dummy_rss1, err
			},
		}
		messageClient := spyMessageClient{}
		subscribeMessagePublisher := publisher.NewSubscribeMessagePublisher(&messageClient)

		command := app_service.PatchCommand{
			Source:             "connpass.com",
			SourceLanguageCode: "en",
			ItemFilter: struct {
				IncludeKeywords []string
				ExcludeKeywords []string
			}{
				IncludeKeywords: []string{"Azure", "Cloud", "Microsoft"},
				ExcludeKeywords: []string{"AWS", "Google Cloud"},
			},
		}

		// Act
		err := app_service.Update(ctx, &logger, &repo, *subscribeMessagePublisher, command)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, messageClient.Messages, 1)
		assert.ElementsMatch(t, messageClient.Messages, []string{
			"{\"feed_url\":\"https://connpass.com/explore/ja.atom\",\"language\":\"en\",\"item_filter\":{\"include_keywords\":[\"Azure\",\"Cloud\",\"Microsoft\"],\"exclude_keywords\":[\"AWS\",\"Google Cloud\"]}}",
		})
	})
}
