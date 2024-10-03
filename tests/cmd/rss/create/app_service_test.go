package create

import (
	"context"
	"testing"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/create/app_service"
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

func TestCreateCommand_Validation(t *testing.T) {
	t.Run("Should Pass All Validations", func(t *testing.T) {
		testCases := []struct {
			name    string
			command app_service.CreateCommand
		}{
			{
				name: "Valid URL and language code",
				command: app_service.CreateCommand{
					FeedURL:            "http://validurl.com",
					SourceLanguageCode: "en",
				},
			},
		}

		ctx := context.Background()
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Act
				err := tc.command.Validation(ctx)
				// Assert
				assert.NoError(t, err)
			})
		}
	})

	t.Run("Should Fail Validations Where Appropriate", func(t *testing.T) {
		testCases := []struct {
			name    string
			command app_service.CreateCommand
		}{
			{
				name: "Invalid URL with valid language code",
				command: app_service.CreateCommand{
					FeedURL:            "htp://invalid-url",
					SourceLanguageCode: "en",
				},
			},
			{
				name: "Valid URL with invalid language code",
				command: app_service.CreateCommand{
					FeedURL:            "http://validurl.com",
					SourceLanguageCode: "fake-lang",
				},
			},
			{
				name: "Valid URL with empty language code",
				command: app_service.CreateCommand{
					FeedURL:            "http://validurl.com",
					SourceLanguageCode: "",
				},
			},
			{
				name: "Empty URL with valid language code",
				command: app_service.CreateCommand{
					FeedURL:            "",
					SourceLanguageCode: "en",
				},
			},
		}

		ctx := context.Background()
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Act
				err := tc.command.Validation(ctx)
				// Assert
				assert.Error(t, err)
			})
		}
	})
}

func TestAppService_Trigger(t *testing.T) {
	t.Run("should handle event successfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		logger := helper.MockLogger{}
		messageClient := spyMessageClient{}
		subscribeMessagePublisher := publisher.NewSubscribeMessagePublisher(&messageClient)

		command := app_service.CreateCommand{
			FeedURL:            "https://azure.microsoft.com/ja-jp/blog/feed",
			SourceLanguageCode: "en",
			ItemFilter:         rss.NewItemFilter([]string{"Azure", "Cloud", "Microsoft"}, []string{"AWS", "Google Cloud"}),
		}

		// Act
		err := app_service.Trigger(ctx, &logger, *subscribeMessagePublisher, command)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, messageClient.Messages, 1)
		assert.ElementsMatch(t, messageClient.Messages, []string{
			"{\"feed_url\":\"https://azure.microsoft.com/ja-jp/blog/feed\",\"language\":\"en\",\"item_filter\":{\"include_keywords\":[\"Azure\",\"Cloud\",\"Microsoft\"],\"exclude_keywords\":[\"AWS\",\"Google Cloud\"]}}",
		})
	})
}
