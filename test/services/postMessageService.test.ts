import 'reflect-metadata';
import { IPostMessageService, PostMessageService } from '../../src/services/postMessageService';
import { SlackHttpApiClient } from '../../src/httpClients/slackHttpClient';

jest.mock('../../src/httpClients/slackHttpClient')

describe('IRSSFeedService', () => {
    let postMessageService: IPostMessageService;
    let mockSlackHttpApiClient: jest.Mocked<SlackHttpApiClient>;

    beforeEach(() => {
        mockSlackHttpApiClient = new SlackHttpApiClient('dummy_token') as jest.Mocked<SlackHttpApiClient>;
        mockSlackHttpApiClient.post.mockImplementation(() => Promise.resolve({}));

        postMessageService = new PostMessageService(mockSlackHttpApiClient);
    });

    describe('postFeedToSlack', () => {
        const testCases = [
            {
                description: 'Empty Feed - No Items',
                mockFeedData: {
                    title: 'Example Feed',
                    description: 'This is a test feed',
                    link: 'http://example.com/rss',
                    lastBuildDate: new Date('2023-12-19'),
                    items: []
                },
                expectedSaveArg:
                    "*フィードタイトル:* <http://example.com/rss|Example Feed>\n" +
                    "*フィード詳細:* This is a test feed\n" +
                    "*最終更新日:* 2023-12-19T00:00:00.000Z\n\n" +
                    "*最新の記事:*\n" +
                    "なし\n"
            },
            {
                description: 'Single Item Feed - One News Item',
                mockFeedData: {
                    title: 'Example Feed',
                    description: 'This is a test feed',
                    link: 'http://example.com/rss',
                    lastBuildDate: new Date('2023-12-19'),
                    items: [
                        {
                            title: "最新ニュース: 例1",
                            link: "http://example.com/news1",
                            pubDate: new Date("2023-12-19T12:00:00Z"),
                            description: "最初のニュースアイテムの説明",
                            contentSnippet: "最初のニュースアイテムのスニペット...",
                            categories: ["最新ニュース", "サンプルカテゴリ"]
                        }]
                },
                expectedSaveArg:
                    "*フィードタイトル:* <http://example.com/rss|Example Feed>\n" +
                    "*フィード詳細:* This is a test feed\n" +
                    "*最終更新日:* 2023-12-19T00:00:00.000Z\n\n" +
                    "*最新の記事:*\n" +
                    "1. *記事タイトル:* <http://example.com/news1|最新ニュース: 例1>\n" +
                    "    *公開日:* Tue Dec 19 2023 21:00:00 GMT+0900 (日本標準時)\n" +
                    "    *概要:* 最初のニュースアイテムのスニペット...\n" +
                    "    *カテゴリ:* 最新ニュース, サンプルカテゴリ\n\n"
            },
            {
                description: 'Multiple Items Feed - Two News Items',
                mockFeedData: {
                    title: 'Example Feed',
                    description: 'This is a test feed',
                    link: 'http://example.com/rss',
                    lastBuildDate: new Date('2023-12-19'),
                    items: [
                        {
                            title: "最新ニュース: 例1",
                            link: "http://example.com/news1",
                            pubDate: new Date("2023-12-19T12:00:00Z"),
                            description: "最初のニュースアイテムの説明",
                            contentSnippet: "最初のニュースアイテムのスニペット...",
                            categories: ["最新ニュース", "サンプルカテゴリ"]
                        },
                        {
                            title: "アップデート情報: 例2",
                            link: "http://example.com/news2",
                            pubDate: new Date("2023-12-18T15:30:00Z"),
                            description: "2番目のニュースアイテムの説明",
                            contentSnippet: "2番目のニュースアイテムのスニペット...",
                            categories: ["アップデート情報", "テクノロジー"]
                        }]
                },
                expectedSaveArg:
                    "*フィードタイトル:* <http://example.com/rss|Example Feed>\n" +
                    "*フィード詳細:* This is a test feed\n" +
                    "*最終更新日:* 2023-12-19T00:00:00.000Z\n\n" +
                    "*最新の記事:*\n" +
                    "1. *記事タイトル:* <http://example.com/news1|最新ニュース: 例1>\n" +
                    "    *公開日:* Tue Dec 19 2023 21:00:00 GMT+0900 (日本標準時)\n" +
                    "    *概要:* 最初のニュースアイテムのスニペット...\n" +
                    "    *カテゴリ:* 最新ニュース, サンプルカテゴリ\n\n" +
                    "2. *記事タイトル:* <http://example.com/news2|アップデート情報: 例2>\n" +
                    "    *公開日:* Tue Dec 19 2023 00:30:00 GMT+0900 (日本標準時)\n" +
                    "    *概要:* 2番目のニュースアイテムのスニペット...\n" +
                    "    *カテゴリ:* アップデート情報, テクノロジー\n\n"
            },
        ]

        testCases.forEach(({ description, mockFeedData, expectedSaveArg }) => {
            it(`should post message to Slack - ${description}`, async () => {
                // Arrange
                const channelId = "ut_channel"
                const fromName = "ut_name"

                // Act (Call the method you want to test)
                await postMessageService.postFeedToSlack(mockFeedData, channelId, fromName)

                // Assert
                expect(mockSlackHttpApiClient.post).toHaveBeenCalledWith(
                    '/chat.postMessage',
                    {
                        channel: channelId,
                        username: fromName,
                        text: expectedSaveArg
                    }
                );
            });
        })

        it('should propagate the exception to the caller in case of an error', async () => {
            // Arrange
            const rss = {
                title: 'Example Feed',
                description: 'This is a test feed',
                link: 'http://example.com/rss',
                lastBuildDate: new Date('2023-12-19'),
                items: []
            }

            const channelId = "ut_channel"
            const fromName = "ut_name"
            mockSlackHttpApiClient.post.mockRejectedValue(new Error('Network Error'));

            // Act&Assert
            await expect(postMessageService.postFeedToSlack(rss, channelId, fromName))
                .rejects
                .toThrow('Network Error');
            expect(mockSlackHttpApiClient.post).toHaveBeenCalledTimes(1);
        });
    })
})