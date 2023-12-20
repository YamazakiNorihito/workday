import 'reflect-metadata';
import axios from 'axios';
import { ILoginAuthenticationService, CognitoOAuth2Service } from '../../src/services/oauth2Service';

jest.mock('axios');

describe('IOAuth2Service', () => {
    describe('ILoginAuthenticationService', () => {
        let mockedAxios: jest.Mocked<typeof axios>;

        beforeEach(() => {
            mockedAxios = axios as jest.Mocked<typeof axios>;
            mockedAxios.post.mockClear()
            mockedAxios.get.mockClear()
        });

        describe('getAuthorizationUrl', () => {
            it(`should return the correct authorization URL with necessary query parameters`, () => {
                // Arrange
                const cognitoDomain = 'https://example.com';
                const cognitoUserPoolURL = 'https://example-user-pool.com';
                const clientId = 'test-client-id';
                const clientSecret = 'test-client-secret';
                const callbackUrl = 'https://callback.com/authorize/callback';
                const loginAuthenticationService: ILoginAuthenticationService =
                    new CognitoOAuth2Service(cognitoDomain, cognitoUserPoolURL, clientId, clientSecret);

                // Act
                const actual = loginAuthenticationService.getAuthorizationUrl(callbackUrl);

                // Assert
                expect(actual.startsWith('https://example.com/oauth2/authorize'));

                const queryParams = new URLSearchParams(actual.split('?')[1]);
                expect(queryParams.get('response_type')).toBe('code');
                expect(queryParams.get('client_id')).toBe('test-client-id');
                expect(queryParams.get('redirect_uri')).toBe('https://callback.com/authorize/callback');
            });
        })

        describe('getAccessToken', () => {
            it(`should retrieve an access token given a valid authorization code and callback URL`, async () => {
                // Arrange
                const cognitoDomain = 'https://example.com';
                const cognitoUserPoolURL = 'https://example-user-pool.com';
                const clientId = 'test-client-id';
                const clientSecret = 'test-client-secret';
                const redirectUri = 'https://callback.com/authorize/callback';
                const authorizationCode = 'valid-authorization-code';

                const mockResponse = {
                    data: {
                        access_token: 'test-access-token',
                        id_token: 'test-id-token',
                        refresh_token: 'test-refresh-token',
                        expires_in: 3600,
                        token_type: 'Bearer'
                    }
                };

                mockedAxios.post.mockResolvedValue(mockResponse);
                mockedAxios.create.mockReturnThis();

                const loginAuthenticationService: ILoginAuthenticationService =
                    new CognitoOAuth2Service(cognitoDomain, cognitoUserPoolURL, clientId, clientSecret);

                // Act
                const actual = await loginAuthenticationService.getAccessToken(authorizationCode, redirectUri);

                // Assert
                expect(mockedAxios.post).toHaveBeenCalledWith(
                    `/oauth2/token`,
                    new URLSearchParams({
                        grant_type: 'authorization_code',
                        client_id: 'test-client-id',
                        redirect_uri: 'https://callback.com/authorize/callback',
                        code: 'valid-authorization-code'
                    }),
                    {
                        headers: {
                            'Authorization': 'Basic dGVzdC1jbGllbnQtaWQ6dGVzdC1jbGllbnQtc2VjcmV0',
                            'Content-Type': 'application/x-www-form-urlencoded'
                        }
                    }
                );
                expect(actual).toEqual({
                    id_token: 'test-id-token',
                    access_token: 'test-access-token',
                    refresh_token: 'test-refresh-token',
                    expires_in: 3600,
                    token_type: 'Bearer'
                });
            });

            it('should propagate the exception to the caller in case of an error', async () => {
                // Arrange
                mockedAxios.post.mockRejectedValue(new Error('Network Error'));
                mockedAxios.create.mockReturnThis();

                const loginAuthenticationService: ILoginAuthenticationService =
                    new CognitoOAuth2Service("", "", "", "");

                // Act&Assert
                await expect(loginAuthenticationService.getAccessToken("", ""))
                    .rejects
                    .toThrow('Network Error');
                expect(mockedAxios.post).toHaveBeenCalledTimes(1);
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