import 'reflect-metadata';
import axios from "axios";
import { FreeeAuthenticationService, IFreeeAuthenticationService } from "../../src/services/freeeHrService";


jest.mock('axios');
describe('IFreeeAuthenticationService', () => {
    let mockedAxios: jest.Mocked<typeof axios>;

    beforeEach(() => {
        mockedAxios = axios as jest.Mocked<typeof axios>;
        mockedAxios.create.mockClear()
        mockedAxios.post.mockClear()
        mockedAxios.get.mockClear()

        mockedAxios.create.mockReturnThis()
    });


    describe('getAuthorizationUrl', () => {
        it(`should return the correct authorization URL with necessary query parameters`, () => {
            // Arrange
            const clientId = 'test-client-id';
            const clientSecret = 'test-client-secret';
            const callbackUrl = 'https://callback.com/authorize/callback';
            const freeeAuthenticationService: IFreeeAuthenticationService =
                new FreeeAuthenticationService(clientId, clientSecret);

            // Act
            const actual = freeeAuthenticationService.getAuthorizationUrl(callbackUrl);

            // Assert
            expect(actual).toMatch(new RegExp(`^(https://accounts.secure.freee.co.jp/public_api/authorize)`));

            const queryParams = new URLSearchParams(actual.split('?')[1]);
            expect(queryParams.get('response_type')).toBe('code');
            expect(queryParams.get('client_id')).toBe('test-client-id');
            expect(queryParams.get('redirect_uri')).toBe('https://callback.com/authorize/callback');
        });
    })


    describe('getAccessToken', () => {
        it(`should retrieve an access token given a valid authorization code and callback URL`, async () => {
            // Arrange
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
                    token_type: 'Bearer',
                    scope: 'test-scope',
                    created_at: 1234,
                    company_id: 5678
                }
            };
            mockedAxios.post.mockResolvedValue(mockResponse);

            const freeeAuthenticationService: IFreeeAuthenticationService =
                new FreeeAuthenticationService(clientId, clientSecret);

            // Act
            const actual = await freeeAuthenticationService.getAccessToken(authorizationCode, redirectUri);

            // Assert
            expect(actual).toEqual({
                id_token: 'test-id-token',
                access_token: 'test-access-token',
                refresh_token: 'test-refresh-token',
                expires_in: 3600,
                token_type: 'Bearer',
                scope: 'test-scope',
                created_at: 1234,
                company_id: 5678
            });
        });

        it('should send a correctly formatted request for getting an access token', async () => {
            // Axiosの実装に依存するため、HttpClientが変わるごとに修正しないといけないUTです。
            // Arrange
            const clientId = 'test-client-id';
            const clientSecret = 'test-client-secret';
            const redirectUri = 'https://callback.com/authorize/callback';
            const authorizationCode = 'valid-authorization-code';

            const freeeAuthenticationService: IFreeeAuthenticationService =
                new FreeeAuthenticationService(clientId, clientSecret);

            // Act
            await freeeAuthenticationService.getAccessToken(authorizationCode, redirectUri);

            // Assert
            expect(mockedAxios.post).toHaveBeenCalledWith(
                `/public_api/token`,
                new URLSearchParams({
                    grant_type: 'authorization_code',
                    client_id: 'test-client-id',
                    client_secret: 'test-client-secret',
                    code: 'valid-authorization-code',
                    redirect_uri: 'https://callback.com/authorize/callback'
                }),
                {
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded'
                    }
                }
            );
        });

        it('should propagate the exception to the caller in case of an error', async () => {
            // Arrange
            mockedAxios.post.mockRejectedValue(new Error('Network Error'));

            const freeeAuthenticationService: IFreeeAuthenticationService =
                new FreeeAuthenticationService("", "");

            // Act&Assert
            await expect(freeeAuthenticationService.getAccessToken("", ""))
                .rejects
                .toThrow('Network Error');
            expect(mockedAxios.post).toHaveBeenCalledTimes(1);
        });
    })

    describe('getPublicKey', () => {
        it('throws an exception when the public key is not published', async () => {
            // Arrange
            const freeeAuthenticationService: IFreeeAuthenticationService =
                new FreeeAuthenticationService("", "");

            // Act&Assert
            await expect(freeeAuthenticationService.getPublicKey(""))
                .rejects
                .toThrow('公開鍵が公開されいません。');
        });
    })

    describe('refreshToken', () => {
        it('should retrieve a new set of tokens given a valid refresh token', async () => {
            // Arrange
            const clientId = 'test-client-id';
            const clientSecret = 'test-client-secret';
            const refreshToken = 'valid-refresh-token';

            const mockResponse = {
                data: {
                    access_token: 'test-access-token',
                    id_token: 'test-id-token',
                    refresh_token: 'test-refresh-token',
                    expires_in: 3600,
                    token_type: 'Bearer',
                    scope: 'test-scope',
                    created_at: 1234,
                    company_id: 5678
                }
            };
            mockedAxios.post.mockResolvedValue(mockResponse);

            const freeeAuthenticationService: IFreeeAuthenticationService =
                new FreeeAuthenticationService(clientId, clientSecret);

            // Act
            const actual = await freeeAuthenticationService.refreshToken(refreshToken);

            // Assert
            expect(actual).toEqual({
                id_token: 'test-id-token',
                access_token: 'test-access-token',
                refresh_token: 'test-refresh-token',
                expires_in: 3600,
                token_type: 'Bearer',
                scope: 'test-scope',
                created_at: 1234,
                company_id: 5678
            });
        });

        it('should correctly send a request to refresh the token', async () => {
            // Axiosの実装に依存するため、HttpClientが変わるごとに修正しないといけないUTです。
            // Arrange
            const clientId = 'test-client-id';
            const clientSecret = 'test-client-secret';
            const refreshToken = 'valid-refresh-token';

            const freeeAuthenticationService: IFreeeAuthenticationService =
                new FreeeAuthenticationService(clientId, clientSecret);

            // Act
            const actual = await freeeAuthenticationService.refreshToken(refreshToken);

            // Assert
            expect(mockedAxios.post).toHaveBeenCalledWith(
                `/oauth2/token`,
                new URLSearchParams({
                    grant_type: 'refresh_token',
                    client_id: 'test-client-id',
                    client_secret: 'test-client-secret',
                    refresh_token: 'valid-refresh-token'
                }),
                {
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded'
                    }
                }
            );
        });

        it('should propagate the exception to the caller in case of an error', async () => {
            // Arrange
            mockedAxios.post.mockRejectedValue(new Error('Network Error'));

            const freeeAuthenticationService: IFreeeAuthenticationService =
                new FreeeAuthenticationService("", "");

            // Act & Assert
            await expect(freeeAuthenticationService.refreshToken("invalid-token"))
                .rejects
                .toThrow('Network Error');
            expect(mockedAxios.post).toHaveBeenCalledTimes(1);
        });
    })
})


describe('IFreeeAuthenticationService', () => {
    let mockedAxios: jest.Mocked<typeof axios>;

    beforeEach(() => {
        mockedAxios = axios as jest.Mocked<typeof axios>;
        mockedAxios.create.mockClear()
        mockedAxios.post.mockClear()
        mockedAxios.get.mockClear()

        mockedAxios.create.mockReturnThis()
    });
})