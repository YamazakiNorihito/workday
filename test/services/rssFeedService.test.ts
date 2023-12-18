import 'reflect-metadata';
import { IRSSFeedService, RSSFeedService } from '../../src/services/rssFeedService';
import { IRSSFeedRepository, RSSFeed } from '../../src/repositories/rssFeedRepository';


jest.mock('rss-parser');

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
                        pubDate: "2023-12-19T12:00:00Z",
                        description: "最初のニュースアイテムの説明",
                        contentSnippet: "最初のニュースアイテムのスニペット...",
                        categories: ["最新ニュース", "サンプルカテゴリ"]
                    },
                    {
                        title: "アップデート情報: 例2",
                        link: "http://example.com/news2",
                        pubDate: "2023-12-18T15:30:00Z",
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
                        pubDate: "2023-12-19T12:00:00Z",
                        description: "最初のニュースアイテムの説明",
                        contentSnippet: "最初のニュースアイテムのスニペット...",
                        categories: ["最新ニュース", "サンプルカテゴリ"]
                    },
                    {
                        title: "アップデート情報: 例2",
                        link: "http://example.com/news2",
                        pubDate: "2023-12-18T15:30:00Z",
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

    describe('getAndSaveRSSFeed', () => { })
})