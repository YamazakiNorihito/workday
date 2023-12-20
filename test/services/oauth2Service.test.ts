import 'reflect-metadata';
import axios from 'axios';
import { ILoginAuthenticationService, CognitoOAuth2Service } from '../../src/services/oauth2Service';

jest.mock('axios');

describe('IOAuth2Service', () => {
    describe('ILoginAuthenticationService', () => {
        let mockedAxios: jest.Mocked<typeof axios>;

        beforeEach(() => {
            mockedAxios = axios as jest.Mocked<typeof axios>;
        });

        describe('getAuthorizationUrl', () => {
            it(`should return the correct authorization URL with necessary query parameters`, async () => {
                // Arrange
                const cognitoDomain = 'https://example.com';
                const cognitoUserPoolURL = 'https://example-user-pool.com';
                const clientId = 'test-client-id';
                const clientSecret = 'test-client-secret';
                const callbackUrl = 'https://callback.com/authorize/callback';

                // Act
                const loginAuthenticationService: ILoginAuthenticationService =
                    new CognitoOAuth2Service(cognitoDomain, cognitoUserPoolURL, clientId, clientSecret);
                const act = loginAuthenticationService.getAuthorizationUrl(callbackUrl);

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