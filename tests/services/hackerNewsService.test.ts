import 'reflect-metadata';
import axios from "axios";
import { HackerNewsItem, IHackerNewsRepository } from "../../src/repositories/hackerNewsRepository";
import { HackerNewsService, IHackerNewsService } from "../../src/services/hackerNewsService";

jest.mock('axios');

describe('IHackerNewsService', () => {
    describe('getTopStories', () => {
        let hackerNewsService: IHackerNewsService;
        let mockHackerNewsRepository: jest.Mocked<IHackerNewsRepository>;
        let mockedAxios: jest.Mocked<typeof axios>;
        beforeEach(() => {
            mockedAxios = axios as jest.Mocked<typeof axios>;
            mockedAxios.create.mockClear();
            mockedAxios.post.mockClear();
            mockedAxios.get.mockClear();
            mockedAxios.create.mockReturnThis();

            mockHackerNewsRepository = {
                save: jest.fn(),
                get: jest.fn(),
            };
            hackerNewsService = new MockHackerNewsService(mockHackerNewsRepository);
        });

        it('should Return All  Stories From DB When Already Stored', async () => {
            // Arrange
            const mockTopStoryIds = [1001, 1002];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getTopStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/topstories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1001);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1002);
            expect(mockHackerNewsRepository.save).toBeCalledTimes(0)
            expect(actual).toEqual([
                {
                    id: 1001,
                    type: 'story',
                    by: 'author1',
                    time: 1618000000,
                    url: 'https://example.com/story1',
                    score: 120,
                    title: 'Example Story 1',
                    descendants: 15,
                },
                {
                    id: 1002,
                    type: 'story',
                    by: 'author2',
                    time: 1618000050,
                    url: 'https://example.com/story2',
                    score: 75,
                    title: 'Example Story 2',
                    descendants: 8,
                },]);
        });
        it('should Fetch And Save Unregistered  Stories From API', async () => {
            // Arrange
            const mockTopStoryIds = [1011, 1012];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getTopStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/topstories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1011);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1012);
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1011,
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
            );
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1012,
                {
                    id: 1012,
                    type: 'story',
                    by: 'new_author2',
                    time: 1618001050,
                    url: 'https://example.com/new_story2',
                    score: 60,
                    title: 'New Example Story 2',
                    descendants: 3,
                },
            );
            expect(actual).toEqual([
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
                {
                    id: 1012,
                    type: 'story',
                    by: 'new_author2',
                    time: 1618001050,
                    url: 'https://example.com/new_story2',
                    score: 60,
                    title: 'New Example Story 2',
                    descendants: 3,
                },]);
        });
        it('should Return Only Story Type Items Regardless of Registration Status', async () => {
            // Arrange
            const mockTopStoryIds = [1001, 1011];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getTopStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/topstories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1001);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1011);
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1011,
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
            );
            expect(actual).toEqual([
                {
                    id: 1001,
                    type: 'story',
                    by: 'author1',
                    time: 1618000000,
                    url: 'https://example.com/story1',
                    score: 120,
                    title: 'Example Story 1',
                    descendants: 15,
                },
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                }]);
        });
        it('should Return Only Story Type Items From Stories', async () => {
            // Arrange
            const mockTopStoryIds = [
                1001 // story
                , 2021  // job
                , 3021 // comment
                , 4021 // poll
                , 5021 // pollopt
            ];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getTopStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/topstories.json');
            expect(actual).toEqual([
                {
                    id: 1001,
                    type: 'story',
                    by: 'author1',
                    time: 1618000000,
                    url: 'https://example.com/story1',
                    score: 120,
                    title: 'Example Story 1',
                    descendants: 15,
                },
                {
                    id: 2021,
                    type: 'job',
                    by: 'new_company1',
                    time: 1618001100,
                    text: 'New job description here for job 1',
                    url: 'https://example.com/new_job1',
                    title: 'New Example Job 1',
                },
                {
                    id: 3021,
                    type: 'comment',
                    by: 'new_commenter1',
                    time: 1618001300,
                    text: 'New Example comment 1',
                    parent: 1011,
                },
                {
                    id: 4021,
                    type: 'poll',
                    by: 'new_poll_creator1',
                    time: 1618001500,
                    text: 'New Poll description for poll 1',
                    score: 30,
                    title: 'New Example Poll 1',
                    parts: [5021],
                    descendants: 1,
                },
                {
                    id: 5021,
                    type: 'pollopt',
                    by: 'new_pollopt_author1',
                    time: 1618001700,
                    parent: 4021,
                    score: 10,
                },]);
        });

        it('should propagate the exception to the caller in case of an error', async () => {
            // Arrange
            mockedAxios.get.mockRejectedValue(new Error('Network Error'));

            // Act&Assert
            await expect(hackerNewsService.getTopStories())
                .rejects
                .toThrow('Network Error');
            expect(mockedAxios.get).toHaveBeenCalledTimes(1);
        });
    })

    describe('getNewStories', () => {
        let hackerNewsService: IHackerNewsService;
        let mockHackerNewsRepository: jest.Mocked<IHackerNewsRepository>;
        let mockedAxios: jest.Mocked<typeof axios>;
        beforeEach(() => {
            mockedAxios = axios as jest.Mocked<typeof axios>;
            mockedAxios.create.mockClear();
            mockedAxios.post.mockClear();
            mockedAxios.get.mockClear();
            mockedAxios.create.mockReturnThis();

            mockHackerNewsRepository = {
                save: jest.fn(),
                get: jest.fn(),
            };
            hackerNewsService = new MockHackerNewsService(mockHackerNewsRepository);
        });

        it('should Return All  Stories From DB When Already Stored', async () => {
            // Arrange
            const mockTopStoryIds = [1001, 1002];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id => {
                const result = mockHackerNewsItems.find(story => story.id === id) ?? null
                return Promise.resolve(result)
            });

            // Act
            const actual = await hackerNewsService.getNewStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/newstories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1001);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1002);
            expect(mockHackerNewsRepository.save).toBeCalledTimes(0)
            expect(actual).toEqual([
                {
                    id: 1001,
                    type: 'story',
                    by: 'author1',
                    time: 1618000000,
                    url: 'https://example.com/story1',
                    score: 120,
                    title: 'Example Story 1',
                    descendants: 15,
                },
                {
                    id: 1002,
                    type: 'story',
                    by: 'author2',
                    time: 1618000050,
                    url: 'https://example.com/story2',
                    score: 75,
                    title: 'Example Story 2',
                    descendants: 8,
                },
            ]);
        });
        it('should Fetch And Save Unregistered  Stories From API', async () => {
            // Arrange
            const mockTopStoryIds = [1011, 1012];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getNewStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/newstories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1011);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1012);
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1011,
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
            );
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1012,
                {
                    id: 1012,
                    type: 'story',
                    by: 'new_author2',
                    time: 1618001050,
                    url: 'https://example.com/new_story2',
                    score: 60,
                    title: 'New Example Story 2',
                    descendants: 3,
                },
            );
            expect(actual).toEqual([
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
                {
                    id: 1012,
                    type: 'story',
                    by: 'new_author2',
                    time: 1618001050,
                    url: 'https://example.com/new_story2',
                    score: 60,
                    title: 'New Example Story 2',
                    descendants: 3,
                },]);
        });
        it('should Return Only Story Type Items Regardless of Registration Status', async () => {
            // Arrange
            const mockTopStoryIds = [1001, 1011];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getNewStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/newstories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1001);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1011);
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1011,
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
            );
            expect(actual).toEqual([
                {
                    id: 1001,
                    type: 'story',
                    by: 'author1',
                    time: 1618000000,
                    url: 'https://example.com/story1',
                    score: 120,
                    title: 'Example Story 1',
                    descendants: 15,
                },
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                }]);
        });
        it('should Return Only Story Type Items From Stories', async () => {
            // Arrange
            const mockTopStoryIds = [
                1001 // story
                , 2021  // job
                , 3021 // comment
                , 4021 // poll
                , 5021 // pollopt
            ];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getNewStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/newstories.json');
            expect(actual).toEqual([
                {
                    id: 1001,
                    type: 'story',
                    by: 'author1',
                    time: 1618000000,
                    url: 'https://example.com/story1',
                    score: 120,
                    title: 'Example Story 1',
                    descendants: 15,
                },
                {
                    id: 2021,
                    type: 'job',
                    by: 'new_company1',
                    time: 1618001100,
                    text: 'New job description here for job 1',
                    url: 'https://example.com/new_job1',
                    title: 'New Example Job 1',
                },
                {
                    id: 3021,
                    type: 'comment',
                    by: 'new_commenter1',
                    time: 1618001300,
                    text: 'New Example comment 1',
                    parent: 1011,
                },
                {
                    id: 4021,
                    type: 'poll',
                    by: 'new_poll_creator1',
                    time: 1618001500,
                    text: 'New Poll description for poll 1',
                    score: 30,
                    title: 'New Example Poll 1',
                    parts: [5021],
                    descendants: 1,
                },
                {
                    id: 5021,
                    type: 'pollopt',
                    by: 'new_pollopt_author1',
                    time: 1618001700,
                    parent: 4021,
                    score: 10,
                },]);
        });

        it('should propagate the exception to the caller in case of an error', async () => {
            // Arrange
            mockedAxios.get.mockRejectedValue(new Error('Network Error'));

            // Act&Assert
            await expect(hackerNewsService.getNewStories())
                .rejects
                .toThrow('Network Error');
            expect(mockedAxios.get).toHaveBeenCalledTimes(1);
        });
    })

    describe('getBestStories', () => {
        let hackerNewsService: IHackerNewsService;
        let mockHackerNewsRepository: jest.Mocked<IHackerNewsRepository>;
        let mockedAxios: jest.Mocked<typeof axios>;
        beforeEach(() => {
            mockedAxios = axios as jest.Mocked<typeof axios>;
            mockedAxios.create.mockClear();
            mockedAxios.post.mockClear();
            mockedAxios.get.mockClear();
            mockedAxios.create.mockReturnThis();

            mockHackerNewsRepository = {
                save: jest.fn(),
                get: jest.fn(),
            };
            hackerNewsService = new MockHackerNewsService(mockHackerNewsRepository);
        });

        it('should Return All  Stories From DB When Already Stored', async () => {
            // Arrange
            const mockTopStoryIds = [1001, 1002];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id => {
                const result = mockHackerNewsItems.find(story => story.id === id) ?? null
                return Promise.resolve(result)
            });

            // Act
            const actual = await hackerNewsService.getBestStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/beststories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1001);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1002);
            expect(mockHackerNewsRepository.save).toBeCalledTimes(0)
            expect(actual).toEqual([
                {
                    id: 1001,
                    type: 'story',
                    by: 'author1',
                    time: 1618000000,
                    url: 'https://example.com/story1',
                    score: 120,
                    title: 'Example Story 1',
                    descendants: 15,
                },
                {
                    id: 1002,
                    type: 'story',
                    by: 'author2',
                    time: 1618000050,
                    url: 'https://example.com/story2',
                    score: 75,
                    title: 'Example Story 2',
                    descendants: 8,
                },
            ]);
        });
        it('should Fetch And Save Unregistered  Stories From API', async () => {
            // Arrange
            const mockTopStoryIds = [1011, 1012];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getBestStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/beststories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1011);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1012);
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1011,
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
            );
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1012,
                {
                    id: 1012,
                    type: 'story',
                    by: 'new_author2',
                    time: 1618001050,
                    url: 'https://example.com/new_story2',
                    score: 60,
                    title: 'New Example Story 2',
                    descendants: 3,
                },
            );
            expect(actual).toEqual([
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
                {
                    id: 1012,
                    type: 'story',
                    by: 'new_author2',
                    time: 1618001050,
                    url: 'https://example.com/new_story2',
                    score: 60,
                    title: 'New Example Story 2',
                    descendants: 3,
                },]);
        });
        it('should Return Only Story Type Items Regardless of Registration Status', async () => {
            // Arrange
            const mockTopStoryIds = [1001, 1011];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getBestStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/beststories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1001);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1011);
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1011,
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
            );
            expect(actual).toEqual([
                {
                    id: 1001,
                    type: 'story',
                    by: 'author1',
                    time: 1618000000,
                    url: 'https://example.com/story1',
                    score: 120,
                    title: 'Example Story 1',
                    descendants: 15,
                },
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                }]);
        });
        it('should Return Only Story Type Items From Stories', async () => {
            // Arrange
            const mockTopStoryIds = [
                1001 // story
                , 2021  // job
                , 3021 // comment
                , 4021 // poll
                , 5021 // pollopt
            ];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getBestStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/beststories.json');
            expect(actual).toEqual([
                {
                    id: 1001,
                    type: 'story',
                    by: 'author1',
                    time: 1618000000,
                    url: 'https://example.com/story1',
                    score: 120,
                    title: 'Example Story 1',
                    descendants: 15,
                },
                {
                    id: 2021,
                    type: 'job',
                    by: 'new_company1',
                    time: 1618001100,
                    text: 'New job description here for job 1',
                    url: 'https://example.com/new_job1',
                    title: 'New Example Job 1',
                },
                {
                    id: 3021,
                    type: 'comment',
                    by: 'new_commenter1',
                    time: 1618001300,
                    text: 'New Example comment 1',
                    parent: 1011,
                },
                {
                    id: 4021,
                    type: 'poll',
                    by: 'new_poll_creator1',
                    time: 1618001500,
                    text: 'New Poll description for poll 1',
                    score: 30,
                    title: 'New Example Poll 1',
                    parts: [5021],
                    descendants: 1,
                },
                {
                    id: 5021,
                    type: 'pollopt',
                    by: 'new_pollopt_author1',
                    time: 1618001700,
                    parent: 4021,
                    score: 10,
                },]);
        });

        it('should propagate the exception to the caller in case of an error', async () => {
            // Arrange
            mockedAxios.get.mockRejectedValue(new Error('Network Error'));

            // Act&Assert
            await expect(hackerNewsService.getBestStories())
                .rejects
                .toThrow('Network Error');
            expect(mockedAxios.get).toHaveBeenCalledTimes(1);
        });
    })

    describe('getAskHNStories', () => {
        let hackerNewsService: IHackerNewsService;
        let mockHackerNewsRepository: jest.Mocked<IHackerNewsRepository>;
        let mockedAxios: jest.Mocked<typeof axios>;
        beforeEach(() => {
            mockedAxios = axios as jest.Mocked<typeof axios>;
            mockedAxios.create.mockClear();
            mockedAxios.post.mockClear();
            mockedAxios.get.mockClear();
            mockedAxios.create.mockReturnThis();

            mockHackerNewsRepository = {
                save: jest.fn(),
                get: jest.fn(),
            };
            hackerNewsService = new MockHackerNewsService(mockHackerNewsRepository);
        });

        it('should Return All  Stories From DB When Already Stored', async () => {
            // Arrange
            const mockTopStoryIds = [1001, 1002];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id => {
                const result = mockHackerNewsItems.find(story => story.id === id) ?? null
                return Promise.resolve(result)
            });

            // Act
            const actual = await hackerNewsService.getAskHNStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/askstories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1001);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1002);
            expect(mockHackerNewsRepository.save).toBeCalledTimes(0)
            expect(actual).toEqual([
                {
                    id: 1001,
                    type: 'story',
                    by: 'author1',
                    time: 1618000000,
                    url: 'https://example.com/story1',
                    score: 120,
                    title: 'Example Story 1',
                    descendants: 15,
                },
                {
                    id: 1002,
                    type: 'story',
                    by: 'author2',
                    time: 1618000050,
                    url: 'https://example.com/story2',
                    score: 75,
                    title: 'Example Story 2',
                    descendants: 8,
                },
            ]);
        });
        it('should Fetch And Save Unregistered  Stories From API', async () => {
            // Arrange
            const mockTopStoryIds = [1011, 1012];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getAskHNStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/askstories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1011);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1012);
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1011,
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
            );
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1012,
                {
                    id: 1012,
                    type: 'story',
                    by: 'new_author2',
                    time: 1618001050,
                    url: 'https://example.com/new_story2',
                    score: 60,
                    title: 'New Example Story 2',
                    descendants: 3,
                },
            );
            expect(actual).toEqual([
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
                {
                    id: 1012,
                    type: 'story',
                    by: 'new_author2',
                    time: 1618001050,
                    url: 'https://example.com/new_story2',
                    score: 60,
                    title: 'New Example Story 2',
                    descendants: 3,
                },]);
        });
        it('should Return Only Story Type Items Regardless of Registration Status', async () => {
            // Arrange
            const mockTopStoryIds = [1001, 1011];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getAskHNStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/askstories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1001);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1011);
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1011,
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
            );
            expect(actual).toEqual([{
                id: 1001,
                type: 'story',
                by: 'author1',
                time: 1618000000,
                url: 'https://example.com/story1',
                score: 120,
                title: 'Example Story 1',
                descendants: 15,
            },
            {
                id: 1011,
                type: 'story',
                by: 'new_author1',
                time: 1618001000,
                url: 'https://example.com/new_story1',
                score: 90,
                title: 'New Example Story 1',
                descendants: 5,
            },]);
        });
        it('should Return Only Story Type Items From Stories', async () => {
            // Arrange
            const mockTopStoryIds = [
                1001 // story
                , 2021  // job
                , 3021 // comment
                , 4021 // poll
                , 5021 // pollopt
            ];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getAskHNStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/askstories.json');
            expect(actual).toEqual([{
                id: 1001,
                type: 'story',
                by: 'author1',
                time: 1618000000,
                url: 'https://example.com/story1',
                score: 120,
                title: 'Example Story 1',
                descendants: 15,
            },
            {
                id: 2021,
                type: 'job',
                by: 'new_company1',
                time: 1618001100,
                text: 'New job description here for job 1',
                url: 'https://example.com/new_job1',
                title: 'New Example Job 1',
            },
            {
                id: 3021,
                type: 'comment',
                by: 'new_commenter1',
                time: 1618001300,
                text: 'New Example comment 1',
                parent: 1011,
            },
            {
                id: 4021,
                type: 'poll',
                by: 'new_poll_creator1',
                time: 1618001500,
                text: 'New Poll description for poll 1',
                score: 30,
                title: 'New Example Poll 1',
                parts: [5021],
                descendants: 1,
            },
            {
                id: 5021,
                type: 'pollopt',
                by: 'new_pollopt_author1',
                time: 1618001700,
                parent: 4021,
                score: 10,
            },]);
        });

        it('should propagate the exception to the caller in case of an error', async () => {
            // Arrange
            mockedAxios.get.mockRejectedValue(new Error('Network Error'));

            // Act&Assert
            await expect(hackerNewsService.getAskHNStories())
                .rejects
                .toThrow('Network Error');
            expect(mockedAxios.get).toHaveBeenCalledTimes(1);
        });
    })

    describe('getShowHNStories', () => {
        let hackerNewsService: IHackerNewsService;
        let mockHackerNewsRepository: jest.Mocked<IHackerNewsRepository>;
        let mockedAxios: jest.Mocked<typeof axios>;
        beforeEach(() => {
            mockedAxios = axios as jest.Mocked<typeof axios>;
            mockedAxios.create.mockClear();
            mockedAxios.post.mockClear();
            mockedAxios.get.mockClear();
            mockedAxios.create.mockReturnThis();

            mockHackerNewsRepository = {
                save: jest.fn(),
                get: jest.fn(),
            };
            hackerNewsService = new MockHackerNewsService(mockHackerNewsRepository);
        });

        it('should Return All  Stories From DB When Already Stored', async () => {
            // Arrange
            const mockTopStoryIds = [1001, 1002];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id => {
                const result = mockHackerNewsItems.find(story => story.id === id) ?? null
                return Promise.resolve(result)
            });

            // Act
            const actual = await hackerNewsService.getShowHNStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/showstories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1001);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1002);
            expect(mockHackerNewsRepository.save).toBeCalledTimes(0)
            expect(actual).toEqual([
                {
                    id: 1001,
                    type: 'story',
                    by: 'author1',
                    time: 1618000000,
                    url: 'https://example.com/story1',
                    score: 120,
                    title: 'Example Story 1',
                    descendants: 15,
                },
                {
                    id: 1002,
                    type: 'story',
                    by: 'author2',
                    time: 1618000050,
                    url: 'https://example.com/story2',
                    score: 75,
                    title: 'Example Story 2',
                    descendants: 8,
                },
            ]);
        });
        it('should Fetch And Save Unregistered  Stories From API', async () => {
            // Arrange
            const mockTopStoryIds = [1011, 1012];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getShowHNStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/showstories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1011);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1012);
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1011,
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
            );
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1012,
                {
                    id: 1012,
                    type: 'story',
                    by: 'new_author2',
                    time: 1618001050,
                    url: 'https://example.com/new_story2',
                    score: 60,
                    title: 'New Example Story 2',
                    descendants: 3,
                },
            );
            expect(actual).toEqual([
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
                {
                    id: 1012,
                    type: 'story',
                    by: 'new_author2',
                    time: 1618001050,
                    url: 'https://example.com/new_story2',
                    score: 60,
                    title: 'New Example Story 2',
                    descendants: 3,
                },]);
        });
        it('should Return Only Story Type Items Regardless of Registration Status', async () => {
            // Arrange
            const mockTopStoryIds = [1001, 1011];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getShowHNStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/showstories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1001);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1011);
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1011,
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
            );
            expect(actual).toEqual([
                {
                    id: 1001,
                    type: 'story',
                    by: 'author1',
                    time: 1618000000,
                    url: 'https://example.com/story1',
                    score: 120,
                    title: 'Example Story 1',
                    descendants: 15,
                },
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                }]);
        });
        it('should Return Only Story Type Items From Stories', async () => {
            // Arrange
            const mockTopStoryIds = [
                1001 // story
                , 2021  // job
                , 3021 // comment
                , 4021 // poll
                , 5021 // pollopt
            ];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getShowHNStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/showstories.json');
            expect(actual).toEqual([{
                id: 1001,
                type: 'story',
                by: 'author1',
                time: 1618000000,
                url: 'https://example.com/story1',
                score: 120,
                title: 'Example Story 1',
                descendants: 15,
            },
            {
                id: 2021,
                type: 'job',
                by: 'new_company1',
                time: 1618001100,
                text: 'New job description here for job 1',
                url: 'https://example.com/new_job1',
                title: 'New Example Job 1',
            },
            {
                id: 3021,
                type: 'comment',
                by: 'new_commenter1',
                time: 1618001300,
                text: 'New Example comment 1',
                parent: 1011,
            },
            {
                id: 4021,
                type: 'poll',
                by: 'new_poll_creator1',
                time: 1618001500,
                text: 'New Poll description for poll 1',
                score: 30,
                title: 'New Example Poll 1',
                parts: [5021],
                descendants: 1,
            },
            {
                id: 5021,
                type: 'pollopt',
                by: 'new_pollopt_author1',
                time: 1618001700,
                parent: 4021,
                score: 10,
            },]);
        });

        it('should propagate the exception to the caller in case of an error', async () => {
            // Arrange
            mockedAxios.get.mockRejectedValue(new Error('Network Error'));

            // Act&Assert
            await expect(hackerNewsService.getShowHNStories())
                .rejects
                .toThrow('Network Error');
            expect(mockedAxios.get).toHaveBeenCalledTimes(1);
        });
    })

    describe('getJobStories', () => {
        let hackerNewsService: IHackerNewsService;
        let mockHackerNewsRepository: jest.Mocked<IHackerNewsRepository>;
        let mockedAxios: jest.Mocked<typeof axios>;
        beforeEach(() => {
            mockedAxios = axios as jest.Mocked<typeof axios>;
            mockedAxios.create.mockClear();
            mockedAxios.post.mockClear();
            mockedAxios.get.mockClear();
            mockedAxios.create.mockReturnThis();

            mockHackerNewsRepository = {
                save: jest.fn(),
                get: jest.fn(),
            };
            hackerNewsService = new MockHackerNewsService(mockHackerNewsRepository);
        });

        it('should Return All  Stories From DB When Already Stored', async () => {
            // Arrange
            const mockTopStoryIds = [1001, 1002];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id => {
                const result = mockHackerNewsItems.find(story => story.id === id) ?? null
                return Promise.resolve(result)
            });

            // Act
            const actual = await hackerNewsService.getJobStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/jobstories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1001);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1002);
            expect(mockHackerNewsRepository.save).toBeCalledTimes(0)
            expect(actual).toEqual([
                {
                    id: 1001,
                    type: 'story',
                    by: 'author1',
                    time: 1618000000,
                    url: 'https://example.com/story1',
                    score: 120,
                    title: 'Example Story 1',
                    descendants: 15,
                },
                {
                    id: 1002,
                    type: 'story',
                    by: 'author2',
                    time: 1618000050,
                    url: 'https://example.com/story2',
                    score: 75,
                    title: 'Example Story 2',
                    descendants: 8,
                },
            ]);
        });
        it('should Fetch And Save Unregistered  Stories From API', async () => {
            // Arrange
            const mockTopStoryIds = [1011, 1012];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getJobStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/jobstories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1011);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1012);
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1011,
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
            );
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1012,
                {
                    id: 1012,
                    type: 'story',
                    by: 'new_author2',
                    time: 1618001050,
                    url: 'https://example.com/new_story2',
                    score: 60,
                    title: 'New Example Story 2',
                    descendants: 3,
                },
            );
            expect(actual).toEqual([
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
                {
                    id: 1012,
                    type: 'story',
                    by: 'new_author2',
                    time: 1618001050,
                    url: 'https://example.com/new_story2',
                    score: 60,
                    title: 'New Example Story 2',
                    descendants: 3,
                },]);
        });
        it('should Return Only Story Type Items Regardless of Registration Status', async () => {
            // Arrange
            const mockTopStoryIds = [1001, 1011];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getJobStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/jobstories.json');
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1001);
            expect(mockHackerNewsRepository.get).toHaveBeenCalledWith(1011);
            expect(mockHackerNewsRepository.save).toHaveBeenCalledWith(
                1011,
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                },
            );
            expect(actual).toEqual([
                {
                    id: 1001,
                    type: 'story',
                    by: 'author1',
                    time: 1618000000,
                    url: 'https://example.com/story1',
                    score: 120,
                    title: 'Example Story 1',
                    descendants: 15,
                },
                {
                    id: 1011,
                    type: 'story',
                    by: 'new_author1',
                    time: 1618001000,
                    url: 'https://example.com/new_story1',
                    score: 90,
                    title: 'New Example Story 1',
                    descendants: 5,
                }]);
        });
        it('should Return Only Story Type Items From Stories', async () => {
            // Arrange
            const mockTopStoryIds = [
                1001 // story
                , 2021  // job
                , 3021 // comment
                , 4021 // poll
                , 5021 // pollopt
            ];
            mockedAxios.get.mockResolvedValueOnce({ data: mockTopStoryIds });

            mockHackerNewsRepository.get.mockImplementation(id =>
                Promise.resolve(mockHackerNewsItems.find(story => story.id === id) ?? null)
            );

            // Act
            const actual = await hackerNewsService.getJobStories();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/jobstories.json');
            expect(actual).toEqual([{
                id: 1001,
                type: 'story',
                by: 'author1',
                time: 1618000000,
                url: 'https://example.com/story1',
                score: 120,
                title: 'Example Story 1',
                descendants: 15,
            },
            {
                id: 2021,
                type: 'job',
                by: 'new_company1',
                time: 1618001100,
                text: 'New job description here for job 1',
                url: 'https://example.com/new_job1',
                title: 'New Example Job 1',
            },
            {
                id: 3021,
                type: 'comment',
                by: 'new_commenter1',
                time: 1618001300,
                text: 'New Example comment 1',
                parent: 1011,
            },
            {
                id: 4021,
                type: 'poll',
                by: 'new_poll_creator1',
                time: 1618001500,
                text: 'New Poll description for poll 1',
                score: 30,
                title: 'New Example Poll 1',
                parts: [5021],
                descendants: 1,
            },
            {
                id: 5021,
                type: 'pollopt',
                by: 'new_pollopt_author1',
                time: 1618001700,
                parent: 4021,
                score: 10,
            },]);
        });

        it('should propagate the exception to the caller in case of an error', async () => {
            // Arrange
            mockedAxios.get.mockRejectedValue(new Error('Network Error'));

            // Act&Assert
            await expect(hackerNewsService.getJobStories())
                .rejects
                .toThrow('Network Error');
            expect(mockedAxios.get).toHaveBeenCalledTimes(1);
        });
    })

    describe('getMaxItem', () => {
        let hackerNewsService: IHackerNewsService;
        let mockHackerNewsRepository: jest.Mocked<IHackerNewsRepository>;
        let mockedAxios: jest.Mocked<typeof axios>;
        beforeEach(() => {
            mockedAxios = axios as jest.Mocked<typeof axios>;
            mockedAxios.create.mockClear();
            mockedAxios.post.mockClear();
            mockedAxios.get.mockClear();
            mockedAxios.create.mockReturnThis();

            mockHackerNewsRepository = {
                save: jest.fn(),
                get: jest.fn(),
            };
            hackerNewsService = new HackerNewsService(mockHackerNewsRepository);
        });

        it('should Return Max Item ID Successfully', async () => {
            // Arrange
            const mockMaxStoryId = 9999;
            mockedAxios.get.mockResolvedValueOnce({ data: mockMaxStoryId });

            // Act
            const actual = await hackerNewsService.getMaxItem();

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/maxitem.json');
            expect(actual).toEqual(9999);
        });

        it('should propagate the exception to the caller in case of an error', async () => {
            // Arrange
            mockedAxios.get.mockRejectedValue(new Error('Network Error'));

            // Act&Assert
            await expect(hackerNewsService.getMaxItem())
                .rejects
                .toThrow('Network Error');
            expect(mockedAxios.get).toHaveBeenCalledTimes(1);
        });
    })

    describe('getItem', () => {
        let hackerNewsService: IHackerNewsService;
        let mockHackerNewsRepository: jest.Mocked<IHackerNewsRepository>;
        let mockedAxios: jest.Mocked<typeof axios>;
        beforeEach(() => {
            mockedAxios = axios as jest.Mocked<typeof axios>;
            mockedAxios.create.mockClear();
            mockedAxios.post.mockClear();
            mockedAxios.get.mockClear();
            mockedAxios.create.mockReturnThis();

            mockHackerNewsRepository = {
                save: jest.fn(),
                get: jest.fn(),
            };
            hackerNewsService = new HackerNewsService(mockHackerNewsRepository);
        });

        it('Jobの取得ができること', async () => {
            // Arrange
            mockedAxios.get.mockResolvedValueOnce({
                data:
                {
                    id: 3001,
                    type: 'comment',
                    by: 'commenter1',
                    time: 1618000300,
                    text: 'Example comment 1',
                    parent: 1001,
                }
            });

            // Act
            const actual = await hackerNewsService.getItem(3001);

            // Assert
            expect(mockedAxios.get).toHaveBeenCalledWith('/item/3001.json');
            expect(actual).toEqual(
                {
                    id: 3001,
                    type: 'comment',
                    by: 'commenter1',
                    time: 1618000300,
                    text: 'Example comment 1',
                    parent: 1001,
                });
        });

        it('should propagate the exception to the caller in case of an error', async () => {
            // Arrange
            mockedAxios.get.mockRejectedValue(new Error('Network Error'));

            // Act&Assert
            await expect(hackerNewsService.getItem(0))
                .rejects
                .toThrow('Network Error');
            expect(mockedAxios.get).toHaveBeenCalledTimes(1);
        });
    })
})


class MockHackerNewsService extends HackerNewsService {
    constructor(mockHackerNewsRepository: IHackerNewsRepository) {
        super(mockHackerNewsRepository);
    }

    public async getItem(itemId: number): Promise<HackerNewsItem> {
        const response = mockUnRegisteredHackerNewsItems.find(o => o.id === itemId);
        if (!response) {
            // アイテムが見つからない場合の処理
            // 例：デフォルトのHackerNewsItemを返す、またはエラーを投げる
            throw new Error(`Item with id ${itemId} not found`);
            // またはデフォルト値を返す
            // return Promise.resolve(/* デフォルトのHackerNewsItem */);
        }

        return Promise.resolve(response);
    }
}

const mockHackerNewsItems: HackerNewsItem[] = [
    // StoryHackerNewsItems
    {
        id: 1001,
        type: 'story',
        by: 'author1',
        time: 1618000000,
        url: 'https://example.com/story1',
        score: 120,
        title: 'Example Story 1',
        descendants: 15,
    },
    {
        id: 1002,
        type: 'story',
        by: 'author2',
        time: 1618000050,
        url: 'https://example.com/story2',
        score: 75,
        title: 'Example Story 2',
        descendants: 8,
    },

    // JobHackerNewsItems
    {
        id: 2001,
        type: 'job',
        by: 'company1',
        time: 1618000100,
        text: 'Job description here for job 1',
        url: 'https://example.com/job1',
        title: 'Example Job 1',
    },
    {
        id: 2002,
        type: 'job',
        by: 'company2',
        time: 1618000200,
        text: 'Job description here for job 2',
        url: 'https://example.com/job2',
        title: 'Example Job 2',
    },

    // CommentHackerNewsItems
    {
        id: 3001,
        type: 'comment',
        by: 'commenter1',
        time: 1618000300,
        text: 'Example comment 1',
        parent: 1001,
    },
    {
        id: 3002,
        type: 'comment',
        by: 'commenter2',
        time: 1618000400,
        text: 'Example comment 2',
        parent: 1002,
    },

    // PollHackerNewsItems
    {
        id: 4001,
        type: 'poll',
        by: 'poll_creator1',
        time: 1618000500,
        text: 'Poll description for poll 1',
        score: 60,
        title: 'Example Poll 1',
        parts: [5001],
        descendants: 3,
    },
    {
        id: 4002,
        type: 'poll',
        by: 'poll_creator2',
        time: 1618000600,
        text: 'Poll description for poll 2',
        score: 45,
        title: 'Example Poll 2',
        parts: [5002, 5003],
        descendants: 4,
    },

    // PollOptionHackerNewsItems
    {
        id: 5001,
        type: 'pollopt',
        by: 'pollopt_author1',
        time: 1618000700,
        parent: 4001,
        score: 30,
    },
    {
        id: 5002,
        type: 'pollopt',
        by: 'pollopt_author2',
        time: 1618000800,
        parent: 4002,
        score: 20,
    }
];

const mockUnRegisteredHackerNewsItems: HackerNewsItem[] = [
    // Unregistered StoryHackerNewsItems
    {
        id: 1011,
        type: 'story',
        by: 'new_author1',
        time: 1618001000,
        url: 'https://example.com/new_story1',
        score: 90,
        title: 'New Example Story 1',
        descendants: 5,
    },
    {
        id: 1012,
        type: 'story',
        by: 'new_author2',
        time: 1618001050,
        url: 'https://example.com/new_story2',
        score: 60,
        title: 'New Example Story 2',
        descendants: 3,
    },

    // Unregistered JobHackerNewsItems
    {
        id: 2021,
        type: 'job',
        by: 'new_company1',
        time: 1618001100,
        text: 'New job description here for job 1',
        url: 'https://example.com/new_job1',
        title: 'New Example Job 1',
    },
    {
        id: 2022,
        type: 'job',
        by: 'new_company2',
        time: 1618001200,
        text: 'New job description here for job 2',
        url: 'https://example.com/new_job2',
        title: 'New Example Job 2',
    },

    // Unregistered CommentHackerNewsItems
    {
        id: 3021,
        type: 'comment',
        by: 'new_commenter1',
        time: 1618001300,
        text: 'New Example comment 1',
        parent: 1011,
    },
    {
        id: 3022,
        type: 'comment',
        by: 'new_commenter2',
        time: 1618001400,
        text: 'New Example comment 2',
        parent: 1012,
    },

    // Unregistered PollHackerNewsItems
    {
        id: 4021,
        type: 'poll',
        by: 'new_poll_creator1',
        time: 1618001500,
        text: 'New Poll description for poll 1',
        score: 30,
        title: 'New Example Poll 1',
        parts: [5021],
        descendants: 1,
    },
    {
        id: 4022,
        type: 'poll',
        by: 'new_poll_creator2',
        time: 1618001600,
        text: 'New Poll description for poll 2',
        score: 25,
        title: 'New Example Poll 2',
        parts: [5022],
        descendants: 2,
    },

    // Unregistered PollOptionHackerNewsItems
    {
        id: 5021,
        type: 'pollopt',
        by: 'new_pollopt_author1',
        time: 1618001700,
        parent: 4021,
        score: 10,
    },
    {
        id: 5022,
        type: 'pollopt',
        by: 'new_pollopt_author2',
        time: 1618001800,
        parent: 4022,
        score: 5,
    }
];
