import 'reflect-metadata';
import { IRSSFeedService, RSSFeedService } from '../../src/services/rssFeedService';
import { IRSSFeedRepository } from '../../src/repositories/rssFeedRepository';
import Parser from 'rss-parser';


jest.mock('rss-parser', () => {
    return jest.fn().mockImplementation(() => ({
        parseURL: jest.fn()
    }));
});

describe('IRSSFeedService', () => {
    let rssFeedService: IRSSFeedService;
    let mockRSSFeedRepository: jest.Mocked<IRSSFeedRepository>;

    beforeEach(() => {
        mockRSSFeedRepository = {
            save: jest.fn(),
            get: jest.fn()
        };
        rssFeedService = new RSSFeedService(mockRSSFeedRepository)
    });

    describe('getRSSFeed', () => {
        it('Should return the RSS feed matching a registered FeedURL', async () => {
            // Arrange
            const feedUrl = 'http://example.com/rss';
            mockRSSFeedRepository.get.mockResolvedValue({
                title: 'Example Feed',
                description: 'This is a test feed',
                link: feedUrl,
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
            });

            // Act
            const actual = await rssFeedService.getRSSFeed(feedUrl);

            // Assert
            expect(actual).toEqual({
                title: 'Example Feed',
                description: 'This is a test feed',
                link: feedUrl,
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
            });
            expect(mockRSSFeedRepository.get).toHaveBeenCalledWith('027efd3619b74e8e26de89e62e4b651f5ca1122f2cc397f4d0d4a86e536a1b4c');
        });
        it('Should return NULL for an unregistered FeedURL', async () => {
            // Arrange
            const feedUrl = 'http://example.com/rss';
            mockRSSFeedRepository.get.mockResolvedValue(null);

            // Act
            const actual = await rssFeedService.getRSSFeed(feedUrl);

            // Assert
            expect(actual).toBeNull();
            expect(mockRSSFeedRepository.get).toHaveBeenCalledWith('027efd3619b74e8e26de89e62e4b651f5ca1122f2cc397f4d0d4a86e536a1b4c');
        });
        it('Should propagate the exception when an error occurs during data retrieval from the repository', async () => {
            // Arrange
            const feedUrl = 'http://example.com/rss';
            mockRSSFeedRepository.get.mockRejectedValue(new Error('Network Error'));

            // Act&Assert
            await expect(rssFeedService.getRSSFeed(feedUrl))
                .rejects
                .toThrow('Network Error');
            expect(mockRSSFeedRepository.get).toHaveBeenCalledWith('027efd3619b74e8e26de89e62e4b651f5ca1122f2cc397f4d0d4a86e536a1b4c');
        });
    })

    describe('getAndSaveRSSFeed', () => {
        beforeEach(() => {
            (Parser as jest.MockedClass<typeof Parser>).mockClear();
        });

        const testCases = [
            {
                description: 'Standard feed with multiple items',
                feedUrl: 'http://example.com/rss',
                mockFeedData: {
                    title: "Example RSS Feed",
                    link: "http://example.com",
                    description: "This is an example RSS feed",
                    lastBuildDate: "2023-01-01T00:00:00.000Z",
                    dummyProperty: "dummy Value",
                    items: []
                },
                expectedSaveArg:
                {
                    title: "Example RSS Feed",
                    link: "http://example.com",
                    description: "This is an example RSS feed",
                    lastBuildDate: new Date("2023-01-01T00:00:00.000Z"),
                    items: []
                }
            },
            {
                description: 'Should correctly process a feed containing a single item',
                feedUrl: 'http://example.com/rss',
                mockFeedData: {
                    title: "Example RSS Feed",
                    link: "http://example.com",
                    description: "This is an example RSS feed",
                    lastBuildDate: "2023-01-01T00:00:00.000Z",
                    dummyProperty: "dummy Value",
                    items: [
                        {
                            title: "Example Article1",
                            link: "http://example.com/article-1",
                            pubDate: "2023-01-01T00:00:00.000Z",
                            contentSnippet: "This is an example article1",
                            categories: ["Category1", "Category2"],
                            dummyPropertyItem1: "dummy Item Value1",
                        }
                    ]
                },
                expectedSaveArg:
                {
                    title: "Example RSS Feed",
                    link: "http://example.com",
                    description: "This is an example RSS feed",
                    lastBuildDate: new Date("2023-01-01T00:00:00.000Z"),
                    items: [
                        {
                            title: "Example Article1",
                            link: "http://example.com/article-1",
                            pubDate: new Date("2023-01-01T00:00:00.000Z"),
                            contentSnippet: "This is an example article1",
                            categories: ["Category1", "Category2"],
                        },
                    ]
                }
            },
            {
                description: 'Should accurately handle a feed with multiple items',
                feedUrl: 'http://example.com/rss',
                mockFeedData: {
                    title: "Example RSS Feed",
                    link: "http://example.com",
                    description: "This is an example RSS feed",
                    lastBuildDate: "2023-01-01T00:00:00.000Z",
                    dummyProperty: "dummy Value",
                    items: [
                        {
                            title: "Example Article1",
                            link: "http://example.com/article-1",
                            pubDate: "2023-01-01T00:00:00.000Z",
                            contentSnippet: "This is an example article1",
                            categories: ["Category1", "Category2"],
                            dummyPropertyItem1: "dummy Item Value1",
                        },
                        {
                            title: "Example Article2",
                            link: "http://example.com/article-2",
                            pubDate: "2023-01-02T00:00:00.000Z",
                            contentSnippet: "This is an example article2",
                            categories: ["Category3", "Category4"],
                            dummyPropertyItem1: "dummy Item Value2",
                        }
                    ]
                },
                expectedSaveArg:
                {
                    title: "Example RSS Feed",
                    link: "http://example.com",
                    description: "This is an example RSS feed",
                    lastBuildDate: new Date("2023-01-01T00:00:00.000Z"),
                    items: [
                        {
                            title: "Example Article1",
                            link: "http://example.com/article-1",
                            pubDate: new Date("2023-01-01T00:00:00.000Z"),
                            contentSnippet: "This is an example article1",
                            categories: ["Category1", "Category2"],
                        },
                        {
                            title: "Example Article2",
                            link: "http://example.com/article-2",
                            pubDate: new Date("2023-01-02T00:00:00.000Z"),
                            contentSnippet: "This is an example article2",
                            categories: ["Category3", "Category4"],
                        },
                    ]
                }
            },
            {
                description: 'Should process a feed correctly even when the title is missing',
                feedUrl: 'http://example.com/rss',
                mockFeedData: {
                    title: "",
                    link: "http://example.com",
                    description: "This is an example RSS feed",
                    lastBuildDate: "2023-01-01T00:00:00.000Z",
                    dummyProperty: "dummy Value",
                    items: []
                },
                expectedSaveArg:
                {
                    title: "",
                    link: "http://example.com",
                    description: "This is an example RSS feed",
                    lastBuildDate: new Date("2023-01-01T00:00:00.000Z"),
                    items: []
                }
            },
            {
                description: 'Should process a feed correctly even when the link is missing',
                feedUrl: 'http://example.com/rss',
                mockFeedData: {
                    title: "Example RSS Feed",
                    link: "",
                    description: "This is an example RSS feed",
                    lastBuildDate: "2023-01-01T00:00:00.000Z",
                    dummyProperty: "dummy Value",
                    items: []
                },
                expectedSaveArg:
                {
                    title: "Example RSS Feed",
                    link: "",
                    description: "This is an example RSS feed",
                    lastBuildDate: new Date("2023-01-01T00:00:00.000Z"),
                    items: []
                }
            },
            {
                description: 'Should process a feed correctly even when the description is missing',
                feedUrl: 'http://example.com/rss',
                mockFeedData: {
                    title: "Example RSS Feed",
                    link: "http://example.com",
                    description: "",
                    lastBuildDate: "2023-01-01T00:00:00.000Z",
                    dummyProperty: "dummy Value",
                    items: []
                },
                expectedSaveArg:
                {
                    title: "Example RSS Feed",
                    link: "http://example.com",
                    description: "",
                    lastBuildDate: new Date("2023-01-01T00:00:00.000Z"),
                    items: []
                }
            },
            {
                description: 'Should correctly handle an item with an empty title in the items array',
                feedUrl: 'http://example.com/rss',
                mockFeedData: {
                    title: "Example RSS Feed",
                    link: "http://example.com",
                    description: "This is an example RSS feed",
                    lastBuildDate: "2023-01-01T00:00:00.000Z",
                    dummyProperty: "dummy Value",
                    items: [
                        {
                            title: "",
                            link: "http://example.com/article-1",
                            pubDate: "2023-01-01T00:00:00.000Z",
                            contentSnippet: "This is an example article1",
                            categories: ["Category1", "Category2"],
                            dummyPropertyItem1: "dummy Item Value1",
                        }
                    ]
                },
                expectedSaveArg:
                {
                    title: "Example RSS Feed",
                    link: "http://example.com",
                    description: "This is an example RSS feed",
                    lastBuildDate: new Date("2023-01-01T00:00:00.000Z"),
                    items: [
                        {
                            title: "",
                            link: "http://example.com/article-1",
                            pubDate: new Date("2023-01-01T00:00:00.000Z"),
                            contentSnippet: "This is an example article1",
                            categories: ["Category1", "Category2"],
                        },
                    ]
                }
            },
            {
                description: 'Should correctly handle an item with an empty link in the items array',
                feedUrl: 'http://example.com/rss',
                mockFeedData: {
                    title: "Example RSS Feed",
                    link: "http://example.com",
                    description: "This is an example RSS feed",
                    lastBuildDate: "2023-01-01T00:00:00.000Z",
                    dummyProperty: "dummy Value",
                    items: [
                        {
                            title: "Example Article1",
                            link: "",
                            pubDate: "2023-01-01T00:00:00.000Z",
                            contentSnippet: "This is an example article1",
                            categories: ["Category1", "Category2"],
                            dummyPropertyItem1: "dummy Item Value1",
                        }
                    ]
                },
                expectedSaveArg:
                {
                    title: "Example RSS Feed",
                    link: "http://example.com",
                    description: "This is an example RSS feed",
                    lastBuildDate: new Date("2023-01-01T00:00:00.000Z"),
                    items: [
                        {
                            title: "Example Article1",
                            link: "",
                            pubDate: new Date("2023-01-01T00:00:00.000Z"),
                            contentSnippet: "This is an example article1",
                            categories: ["Category1", "Category2"],
                        },
                    ]
                }
            },
            {
                description: 'Should correctly handle an item with an empty pubDate in the items array',
                feedUrl: 'http://example.com/rss',
                mockFeedData: {
                    title: "Example RSS Feed",
                    link: "http://example.com",
                    description: "This is an example RSS feed",
                    lastBuildDate: "2023-01-01T00:00:00.000Z",
                    dummyProperty: "dummy Value",
                    items: [
                        {
                            title: "Example Article1",
                            link: "http://example.com/article-1",
                            pubDate: "",
                            contentSnippet: "This is an example article1",
                            categories: ["Category1", "Category2"],
                            dummyPropertyItem1: "dummy Item Value1",
                        }
                    ]
                },
                expectedSaveArg:
                {
                    title: "Example RSS Feed",
                    link: "http://example.com",
                    description: "This is an example RSS feed",
                    lastBuildDate: new Date("2023-01-01T00:00:00.000Z"),
                    items: [
                        {
                            title: "Example Article1",
                            link: "http://example.com/article-1",
                            contentSnippet: "This is an example article1",
                            categories: ["Category1", "Category2"],
                        },
                    ]
                }
            },
            {
                description: 'Should correctly handle an item with an empty categories in the items array',
                feedUrl: 'http://example.com/rss',
                mockFeedData: {
                    title: "Example RSS Feed",
                    link: "http://example.com",
                    description: "This is an example RSS feed",
                    lastBuildDate: "2023-01-01T00:00:00.000Z",
                    dummyProperty: "dummy Value",
                    items: [
                        {
                            title: "Example Article1",
                            link: "http://example.com/article-1",
                            pubDate: "2023-01-01T00:00:00.000Z",
                            contentSnippet: "This is an example article1",
                            dummyPropertyItem1: "dummy Item Value1",
                        }
                    ]
                },
                expectedSaveArg:
                {
                    title: "Example RSS Feed",
                    link: "http://example.com",
                    description: "This is an example RSS feed",
                    lastBuildDate: new Date("2023-01-01T00:00:00.000Z"),
                    items: [
                        {
                            title: "Example Article1",
                            link: "http://example.com/article-1",
                            pubDate: new Date("2023-01-01T00:00:00.000Z"),
                            contentSnippet: "This is an example article1",
                        },
                    ]
                }
            },
        ];

        testCases.forEach(({ description, feedUrl, mockFeedData, expectedSaveArg }) => {
            it(`Should fetch and save the RSS feed successfully - ${description}`, async () => {
                // Arrange
                const mockParseURL = jest.fn().mockResolvedValue(mockFeedData);
                (Parser as jest.MockedClass<typeof Parser>).mockImplementation(() => ({
                    parseURL: mockParseURL,
                    parseString: jest.fn()
                }));

                // Act
                await rssFeedService.getAndSaveRSSFeed(feedUrl);

                // Assert
                expect(mockParseURL).toHaveBeenCalledWith(feedUrl);
                expect(mockRSSFeedRepository.save).toHaveBeenCalledWith(
                    '027efd3619b74e8e26de89e62e4b651f5ca1122f2cc397f4d0d4a86e536a1b4c',
                    expectedSaveArg
                );
            });
        });
    })
})