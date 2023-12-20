import 'reflect-metadata';
import axios from 'axios';

jest.mock('axios');

describe('IOAuth2Service', () => {
    describe('ILoginAuthenticationService', () => {
        let mockedAxios: jest.Mocked<typeof axios>;

        beforeEach(() => {
            mockedAxios = axios as jest.Mocked<typeof axios>;
        });

        describe('getAuthorizationUrl', () => {
            it(``, async () => {
                // Arrange

                // Act

                // Assert
            });

            it('should propagate the exception to the caller in case of an error', async () => {
                // Arrange
                //mockSlackHttpApiClient.post.mockRejectedValue(new Error('Network Error'));

                // Act&Assert
                //await expect(postMessageService.postFeedToSlack(rss, channelId, fromName))
                //.rejects
                //.toThrow('Network Error');
                //(mockSlackHttpApiClient.post).toHaveBeenCalledTimes(1);
            });
        })

        describe('getAccessToken', () => {
            it(``, async () => {
                // Arrange

                // Act

                // Assert
            });

            it('should propagate the exception to the caller in case of an error', async () => {
                // Arrange
                //mockSlackHttpApiClient.post.mockRejectedValue(new Error('Network Error'));

                // Act&Assert
                //await expect(postMessageService.postFeedToSlack(rss, channelId, fromName))
                //.rejects
                //.toThrow('Network Error');
                //(mockSlackHttpApiClient.post).toHaveBeenCalledTimes(1);
            });
        })

        describe('getPublicKey', () => {
            it(``, async () => {
                // Arrange

                // Act

                // Assert
            });

            it('should propagate the exception to the caller in case of an error', async () => {
                // Arrange
                //mockSlackHttpApiClient.post.mockRejectedValue(new Error('Network Error'));

                // Act&Assert
                //await expect(postMessageService.postFeedToSlack(rss, channelId, fromName))
                //.rejects
                //.toThrow('Network Error');
                //(mockSlackHttpApiClient.post).toHaveBeenCalledTimes(1);
            });
        })

        describe('refreshToken', () => {
            it(``, async () => {
                // Arrange

                // Act

                // Assert
            });

            it('should propagate the exception to the caller in case of an error', async () => {
                // Arrange
                //mockSlackHttpApiClient.post.mockRejectedValue(new Error('Network Error'));

                // Act&Assert
                //await expect(postMessageService.postFeedToSlack(rss, channelId, fromName))
                //.rejects
                //.toThrow('Network Error');
                //(mockSlackHttpApiClient.post).toHaveBeenCalledTimes(1);
            });
        })
    })
})