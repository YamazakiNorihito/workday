import 'reflect-metadata';
import axios from "axios";
import redis from 'redis';
import { Employee, FreeeAuthenticationService, FreeeService, IFreeeAuthenticationService, IFreeeService } from "../../src/services/freeeHrService";
import { FreeeHrHttpApiClient, IFreeeHrHttpApiClient } from '../../src/httpClients/freeeHttpClient';
import { FreeeUserRepository, IFreeeUserRepository, User } from '../../src/repositories/freeeUserRepository';
import { RedisClientType } from 'redis';
import { DateOnly } from '../../src/types/dateOnly';
import { TimeOnly } from '../../src/types/timeOnly';


jest.mock('axios');
jest.mock('redis', () => require('redis-mock'));
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


describe('IFreeeService', () => {
    let mockedFreeeAuthenticationService: jest.Mocked<IFreeeAuthenticationService>;
    let mockedFreeeHrHttpApiClient: jest.Mocked<IFreeeHrHttpApiClient>;
    let mockedFreeeUserRepository: jest.Mocked<IFreeeUserRepository>;
    let mockedRedisClient: jest.Mocked<RedisClientType>;

    beforeEach(() => {
        mockedFreeeAuthenticationService = {
            getAuthorizationUrl: jest.fn(),
            getAccessToken: jest.fn(),
            getPublicKey: jest.fn(),
            refreshToken: jest.fn(),
        };
        mockedFreeeHrHttpApiClient = {
            get: jest.fn(),
            post: jest.fn(),
            put: jest.fn(),
            delete: jest.fn(),
        };
        mockedFreeeUserRepository = {
            isReady: jest.fn(),
            save: jest.fn(),
            get: jest.fn(),
        };
        mockedRedisClient = redis.createClient() as jest.Mocked<RedisClientType>;
    });

    describe('getAuthorizeUrl', () => {
        beforeEach(() => {
            mockedFreeeAuthenticationService.getAuthorizationUrl.mockClear();
        });

        it(`should return the correct authorization URL with necessary query parameters`, () => {
            // Arrange
            const freeeService: IFreeeService =
                new FreeeService(mockedFreeeAuthenticationService, mockedFreeeHrHttpApiClient, mockedFreeeUserRepository, mockedRedisClient);
            mockedFreeeAuthenticationService.getAuthorizationUrl = jest.fn((callbackUrl: string) => {
                return `https://test.com/public_api/authorize?client_id=clientId&redirect_uri=${callbackUrl}&response_type=code&prompt=select_company`
            });
            const redirectUri = 'https://callback.com/authorize/callback';

            // Act
            const actual = freeeService.getAuthorizeUrl(redirectUri);

            // Assert
            expect(actual).toBe(`https://test.com/public_api/authorize?client_id=clientId&redirect_uri=https://callback.com/authorize/callback&response_type=code&prompt=select_company`)
        });

        it('should propagate the exception to the caller in case of an error', async () => {
            // Arrange
            const freeeService: IFreeeService =
                new FreeeService(mockedFreeeAuthenticationService, mockedFreeeHrHttpApiClient, mockedFreeeUserRepository, mockedRedisClient);
            mockedFreeeAuthenticationService.getAuthorizationUrl.mockImplementation(() => {
                throw new Error('Network Error');
            });

            // Act & Assert
            expect(() => {
                freeeService.getAuthorizeUrl('https://callback.com/authorize/callback')
            }).toThrow('Network Error');
        });
    })

    describe('handleAuthCallback', () => {
        beforeEach(() => {
            mockedFreeeAuthenticationService.getAccessToken.mockClear();
            mockedFreeeHrHttpApiClient.get.mockClear();
            jest.spyOn(Date, 'now').mockImplementation(() => new Date('2024-01-01T00:00:00.000Z').getTime());
        });

        it(`Should Save Token and Retrieve 'Me' Information for Frontier Field Corporation`, async () => {
            // Arrange
            const mockFreeeOAuthTokenResponse = {
                access_token: 'test-access-token',
                id_token: 'test-id-token',
                refresh_token: 'test-refresh-token',
                expires_in: 3600,
                token_type: 'Bearer',
                scope: 'test-scope',
                created_at: 1234,
                company_id: 5678
            };
            mockedFreeeAuthenticationService.getAccessToken.mockResolvedValueOnce(mockFreeeOAuthTokenResponse);

            const meResponse = {
                id: 1,
                companies: [
                    {
                        id: 103,
                        name: "株式会社フロンティア・フィールド",
                        role: 'self_only',
                        external_cid: 3456789012,
                        employee_id: 3001,
                        display_name: "鈴木一郎"
                    }
                ]
            };
            mockedFreeeHrHttpApiClient.get.mockResolvedValueOnce(meResponse);

            let actual_userId: string | null = null;
            let actual_user: User | null = null;
            mockedFreeeUserRepository.save.mockImplementation((userId, user) => {
                actual_userId = userId;
                actual_user = user;
                return Promise.resolve();
            })

            const freeeService: IFreeeService =
                new FreeeService(mockedFreeeAuthenticationService, mockedFreeeHrHttpApiClient, mockedFreeeUserRepository, mockedRedisClient);

            const redirectUri = 'https://callback.com/authorize/callback';
            const userId = 'test-123456';
            const authCode = 'test-auth-code';
            // Act
            const actual = await freeeService.handleAuthCallback(userId, authCode, redirectUri);

            // Assert
            expect(actual).toEqual({
                employee_id: 3001,
                employee_name: "鈴木一郎",
                company_id: 103,
                company_name: "株式会社フロンティア・フィールド",
                external_cid: 3456789012,
                role: 'self_only',
                updated_at: new Date('2024-01-01T00:00:00.000Z')
            });

            // 保存データが正しいことを検証する
            expect(actual_userId).toBe('test-123456');
            expect(actual_user).toEqual({
                id: 1,
                companies: [
                    {
                        id: 103,
                        name: "株式会社フロンティア・フィールド",
                        role: 'self_only',
                        external_cid: 3456789012,
                        employee_id: 3001,
                        display_name: "鈴木一郎"
                    }
                ],
                oauth: {
                    access_token: 'test-access-token',
                    id_token: 'test-id-token',
                    refresh_token: 'test-refresh-token',
                    expires_in: 3600,
                    token_type: 'Bearer',
                    scope: 'test-scope',
                    created_at: 1234,
                    company_id: 5678
                },
                updated_at: (new Date('2024-01-01T00:00:00.000Z')).getTime()
            });
        });

        it(`Should Save Token But Not Retrieve 'Me' Information When Frontier Field Corporation's Info is Missing`, async () => {
            // Arrange
            const mockFreeeOAuthTokenResponse = {
                access_token: 'test-access-token',
                id_token: 'test-id-token',
                refresh_token: 'test-refresh-token',
                expires_in: 3600,
                token_type: 'Bearer',
                scope: 'test-scope',
                created_at: 1234,
                company_id: 5678
            };
            mockedFreeeAuthenticationService.getAccessToken.mockResolvedValueOnce(mockFreeeOAuthTokenResponse);

            const meResponse = {
                id: 1,
                companies: [
                    {
                        id: 103,
                        name: "株式会社日経HR",
                        role: 'self_only',
                        external_cid: 3456789012,
                        employee_id: 3001,
                        display_name: "鈴木一郎"
                    }
                ]
            };
            mockedFreeeHrHttpApiClient.get.mockResolvedValueOnce(meResponse);

            let actual_userId: string | null = null;
            let actual_user: User | null = null;
            mockedFreeeUserRepository.save.mockImplementation((userId, user) => {
                actual_userId = userId;
                actual_user = user;
                return Promise.resolve();
            })

            const freeeService: IFreeeService =
                new FreeeService(mockedFreeeAuthenticationService, mockedFreeeHrHttpApiClient, mockedFreeeUserRepository, mockedRedisClient);

            const redirectUri = 'https://callback.com/authorize/callback';
            const userId = 'test-123456';
            const authCode = 'test-auth-code';
            // Act
            const actual = await freeeService.handleAuthCallback(userId, authCode, redirectUri);

            // Assert
            expect(actual).toEqual({
                employee_id: 0,
                employee_name: "",
                company_id: 0,
                company_name: "",
                external_cid: 0,
                role: "",
                updated_at: null
            });

            // 保存データが正しいことを検証する
            expect(actual_userId).toBe('test-123456');
            expect(actual_user).toEqual({
                id: 1,
                companies: [
                    {
                        id: 103,
                        name: "株式会社日経HR",
                        role: 'self_only',
                        external_cid: 3456789012,
                        employee_id: 3001,
                        display_name: "鈴木一郎"
                    }
                ],
                oauth: {
                    access_token: 'test-access-token',
                    id_token: 'test-id-token',
                    refresh_token: 'test-refresh-token',
                    expires_in: 3600,
                    token_type: 'Bearer',
                    scope: 'test-scope',
                    created_at: 1234,
                    company_id: 5678
                },
                updated_at: (new Date('2024-01-01T00:00:00.000Z')).getTime()
            });
        });

        it(`Should Save Token But Not Retrieve 'Me' Information When 'Me' Data is Absent`, async () => {
            // Arrange
            const mockFreeeOAuthTokenResponse = {
                access_token: 'test-access-token',
                id_token: 'test-id-token',
                refresh_token: 'test-refresh-token',
                expires_in: 3600,
                token_type: 'Bearer',
                scope: 'test-scope',
                created_at: 1234,
                company_id: 5678
            };
            mockedFreeeAuthenticationService.getAccessToken.mockResolvedValueOnce(mockFreeeOAuthTokenResponse);
            mockedFreeeHrHttpApiClient.get.mockResolvedValueOnce(null);

            let actual_userId: string | null = null;
            let actual_user: User | null = null;
            mockedFreeeUserRepository.save.mockImplementation((userId, user) => {
                actual_userId = userId;
                actual_user = user;
                return Promise.resolve();
            })

            const freeeService: IFreeeService =
                new FreeeService(mockedFreeeAuthenticationService, mockedFreeeHrHttpApiClient, mockedFreeeUserRepository, mockedRedisClient);

            const redirectUri = 'https://callback.com/authorize/callback';
            const userId = 'test-123456';
            const authCode = 'test-auth-code';
            // Act
            const actual = await freeeService.handleAuthCallback(userId, authCode, redirectUri);

            // Assert
            expect(actual).toEqual({
                employee_id: 0,
                employee_name: "",
                company_id: 0,
                company_name: "",
                external_cid: 0,
                role: "",
                updated_at: null
            });

            // 保存データが正しいことを検証する
            expect(actual_userId).toBe('test-123456');
            expect(actual_user).toEqual({
                id: undefined,
                companies: undefined,
                oauth: {
                    access_token: 'test-access-token',
                    id_token: 'test-id-token',
                    refresh_token: 'test-refresh-token',
                    expires_in: 3600,
                    token_type: 'Bearer',
                    scope: 'test-scope',
                    created_at: 1234,
                    company_id: 5678
                },
                updated_at: (new Date('2024-01-01T00:00:00.000Z')).getTime()
            });
        });

        it('should propagate the exception to the caller in case of an error', async () => {
            // Arrange
            mockedFreeeAuthenticationService.getAccessToken.mockImplementation(() => {
                throw new Error('Network Error');
            });
            const freeeService: IFreeeService =
                new FreeeService(mockedFreeeAuthenticationService, mockedFreeeHrHttpApiClient, mockedFreeeUserRepository, mockedRedisClient);

            // Act & Assert
            await expect(freeeService.handleAuthCallback('test-123456', 'test-auth-code', 'https://callback.com/authorize/callback'))
                .rejects
                .toThrow('Network Error');
        });
    })

    describe('getWorkRecords', () => {
        beforeEach(() => {
            mockedFreeeAuthenticationService.getAccessToken.mockClear();
            mockedFreeeHrHttpApiClient.get.mockClear();
            jest.spyOn(Date, 'now').mockImplementation(() => new Date('2024-01-01T00:00:00.000Z').getTime());
        });

        it(`should verify that the API call is made correctly`, async () => {
            // Arrange
            mockedFreeeHrHttpApiClient.get.mockResolvedValue(
                makeMockDataWorkRecordSummaries([
                    {
                        "break_records": [
                            {
                                "clock_in_at": "2024-01-10T12:00:00",
                                "clock_out_at": "2024-01-10T13:00:00"
                            }
                        ],
                        "clock_in_at": "2024-01-10T09:00:00",
                        "clock_out_at": "2024-01-10T17:00:00",
                        "date": "2024-01-10",
                    }
                ]));

            const freeeService: IFreeeService =
                new TestFreeeService(mockedFreeeAuthenticationService, mockedFreeeHrHttpApiClient, mockedFreeeUserRepository, mockedRedisClient);

            const userId = 'test-123456';
            // Act
            const actual = await freeeService.getWorkRecords(userId, 2023, 12);

            // Assert
            expect(mockedFreeeHrHttpApiClient.get).toHaveBeenCalledWith(
                `/api/v1/employees/3001/work_record_summaries/2023/12?company_id=103&work_records=true`,
                "test-token"
            );
            expect(actual).toEqual([
                {
                    workDay: new DateOnly(2024, 1, 10),
                    breakRecords: [{
                        clockInAt: new TimeOnly(12, 0, 0),
                        clockOutAt: new TimeOnly(13, 0, 0)
                    }],
                    clockInAt: new TimeOnly(9, 0, 0),
                    clockOutAt: new TimeOnly(17, 0, 0),
                }
            ]);
        });

        it(`should correctly retrieve workdays`, async () => {
            // Arrange
            mockedFreeeHrHttpApiClient.get.mockResolvedValue(
                makeMockDataWorkRecordSummaries([
                    {
                        "break_records": [
                            {
                                "clock_in_at": "2024-01-10T10:00:00",
                                "clock_out_at": "2024-01-10T11:00:00"
                            },
                            {
                                "clock_in_at": "2024-01-10T12:00:00",
                                "clock_out_at": "2024-01-10T13:00:00"
                            }
                        ],
                        "clock_in_at": "2024-01-10T09:00:00",
                        "clock_out_at": "2024-01-10T17:00:00",
                        "date": "2024-01-10",
                    }, {

                        "break_records": [
                            {
                                "clock_in_at": "2024-01-15T14:00:00",
                                "clock_out_at": "2024-01-15T15:00:00"
                            }
                        ],
                        "clock_in_at": "2024-01-15T08:30:00",
                        "clock_out_at": "2024-01-15T18:30:00",
                        "date": "2024-01-15",
                    }, {

                        "break_records": [],
                        "clock_in_at": "2024-01-20T09:00:00",
                        "clock_out_at": "2024-01-20T17:00:00",
                        "date": "2024-01-20",
                    }
                ]));

            const freeeService: IFreeeService =
                new TestFreeeService(mockedFreeeAuthenticationService, mockedFreeeHrHttpApiClient, mockedFreeeUserRepository, mockedRedisClient);

            const userId = 'test-123456';
            // Act
            const actual = await freeeService.getWorkRecords(userId, 2023, 12);

            // Assert
            expect(actual).toEqual([
                {
                    workDay: new DateOnly(2024, 1, 10),
                    breakRecords: [{
                        clockInAt: new TimeOnly(10, 0, 0),
                        clockOutAt: new TimeOnly(11, 0, 0)
                    }, {
                        clockInAt: new TimeOnly(12, 0, 0),
                        clockOutAt: new TimeOnly(13, 0, 0)
                    }],
                    clockInAt: new TimeOnly(9, 0, 0),
                    clockOutAt: new TimeOnly(17, 0, 0),
                },
                {
                    workDay: new DateOnly(2024, 1, 15),
                    breakRecords: [{
                        clockInAt: new TimeOnly(14, 0, 0),
                        clockOutAt: new TimeOnly(15, 0, 0)
                    }],
                    clockInAt: new TimeOnly(8, 30, 0),
                    clockOutAt: new TimeOnly(18, 30, 0),
                },
                {
                    workDay: new DateOnly(2024, 1, 20),
                    breakRecords: [],
                    clockInAt: new TimeOnly(9, 0, 0),
                    clockOutAt: new TimeOnly(17, 0, 0),
                }
            ]);
        });

        it(`should not retrieve workdays with no work time entered`, async () => {
            // Arrange
            mockedFreeeHrHttpApiClient.get.mockResolvedValue(
                makeMockDataWorkRecordSummaries([
                    {
                        "break_records": [
                            {
                                "clock_in_at": "2024-01-10T10:00:00",
                                "clock_out_at": "2024-01-10T11:00:00"
                            },
                            {
                                "clock_in_at": "2024-01-10T12:00:00",
                                "clock_out_at": "2024-01-10T13:00:00"
                            }
                        ],
                        "clock_in_at": null,
                        "clock_out_at": null,
                        "date": "2024-01-10",
                    }, {
                        "break_records": [
                            {
                                "clock_in_at": "2024-01-15T14:00:00",
                                "clock_out_at": "2024-01-15T15:00:00"
                            }
                        ],
                        "clock_in_at": undefined,
                        "clock_out_at": undefined,
                        "date": "2024-01-15",
                    }
                ]));

            const freeeService: IFreeeService =
                new TestFreeeService(mockedFreeeAuthenticationService, mockedFreeeHrHttpApiClient, mockedFreeeUserRepository, mockedRedisClient);

            const userId = 'test-123456';
            // Act
            const actual = await freeeService.getWorkRecords(userId, 2023, 12);

            // Assert
            expect(actual).toHaveLength(0)
        });

        it('should propagate the exception to the caller in case of an error', async () => {
            // Arrange
            mockedFreeeAuthenticationService.getAccessToken.mockImplementation(() => {
                throw new Error('Network Error');
            });
            const freeeService: IFreeeService =
                new FreeeService(mockedFreeeAuthenticationService, mockedFreeeHrHttpApiClient, mockedFreeeUserRepository, mockedRedisClient);

            // Act & Assert
            await expect(freeeService.handleAuthCallback('test-123456', 'test-auth-code', 'https://callback.com/authorize/callback'))
                .rejects
                .toThrow('Network Error');
        });

        function makeMockDataWorkRecordSummaries(workRecords: {
            date: string,
            clock_in_at: string | null | undefined,
            clock_out_at: string | null | undefined,
            break_records: {
                clock_in_at: string,
                clock_out_at: string
            }[]
        }[]) {
            const mockDataWorkRecords = workRecords.map(workRecord => {
                return {
                    "break_records": workRecord.break_records.map(break_record => ({
                        "clock_in_at": break_record.clock_in_at,
                        "clock_out_at": break_record.clock_out_at,
                    })),
                    "clock_in_at": workRecord.clock_in_at,
                    "clock_out_at": workRecord.clock_out_at,
                    "date": workRecord.date,
                    "day_pattern": "normal_day",
                    "early_leaving_mins": 0,
                    "hourly_paid_holiday_mins": 0,
                    "is_absence": false,
                    "is_editable": true,
                    "lateness_mins": 0,
                    "normal_work_clock_in_at": `${workRecord.date}T08:30:00`,
                    "normal_work_clock_out_at": `${workRecord.date}T18:30:00`,
                    "normal_work_mins": 480,
                    "normal_work_mins_by_paid_holiday": 0,
                    "note": "Standard workday.",
                    "paid_holiday": 0,
                    "use_attendance_deduction": false,
                    "use_default_work_pattern": true,
                    "total_overtime_work_mins": 0,
                    "total_holiday_work_mins": 0,
                    "total_latenight_work_mins": 0,
                    "not_auto_calc_work_time": false,
                    "total_excess_statutory_work_mins": 0,
                    "total_latenight_excess_statutory_work_mins": 0,
                    "total_overtime_except_normal_work_mins": 0,
                    "total_latenight_overtime_except_normal_work_min": 0
                };
            });

            return {
                year: 2024,
                month: 1,
                start_date: "2024-01-01",
                end_date: "2024-01-31",
                work_days: 20,
                total_work_mins: 9600,
                total_normal_work_mins: 8000,
                total_excess_statutory_work_mins: 400,
                total_overtime_except_normal_work_mins: 200,
                total_overtime_within_normal_work_mins: 100,
                total_holiday_work_mins: 300,
                total_latenight_work_mins: 50,
                num_absences: 1,
                num_paid_holidays: 2,
                num_paid_holidays_and_hours: {
                    days: 2,
                    hours: 16
                },
                num_paid_holidays_left: 10,
                num_paid_holidays_and_hours_left: {
                    days: 10,
                    hours: 80
                },
                num_substitute_holidays_used: 1,
                num_compensatory_holidays_used: 1,
                num_special_holidays_used: 1,
                num_special_holidays_and_hours_used: {
                    days: 1,
                    hours: 8
                },
                total_lateness_and_early_leaving_mins: 30,
                multi_hourly_wages: [
                    {
                        name: "Regular",
                        total_normal_time_mins: 8000
                    },
                    {
                        name: "Overtime",
                        total_normal_time_mins: 500
                    }
                ],
                work_records: mockDataWorkRecords,
                total_shortage_work_mins: 0,
                total_deemed_paid_excess_statutory_work_mins: 0,
                total_deemed_paid_overtime_except_normal_work_mins: 0
            };
        }
    })
})

// MeやTokenを取得するメソッドが複雑なので、簡素化したServiceクラス
class TestFreeeService extends FreeeService {

    protected getInternalMe(user: User): Employee {
        return {
            employee_id: 3001,
            employee_name: "鈴木一郎",
            company_id: 103,
            company_name: "株式会社フロンティア・フィールド",
            external_cid: 3456789012,
            role: 'self_only',
            updated_at: new Date('2024-01-01T00:00:00.000Z')
        }
    }

    protected async getAccessToken(userId: string, retryCount: number = 0): Promise<string> {
        return "test-token";
    }

}