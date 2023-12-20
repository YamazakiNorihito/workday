import 'reflect-metadata';
import axios from 'axios';
import { ILoginAuthenticationService, CognitoOAuth2Service } from '../../src/services/oauth2Service';

jest.mock('axios');

describe('IOAuth2Service', () => {
    describe('ILoginAuthenticationService', () => {
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

                const loginAuthenticationService: ILoginAuthenticationService =
                    new CognitoOAuth2Service(cognitoDomain, cognitoUserPoolURL, clientId, clientSecret);

                // Act
                const actual = await loginAuthenticationService.getAccessToken(authorizationCode, redirectUri);

                // Assert
                expect(actual).toEqual({
                    id_token: 'test-id-token',
                    access_token: 'test-access-token',
                    refresh_token: 'test-refresh-token',
                    expires_in: 3600,
                    token_type: 'Bearer'
                });
            });

            it('should send a correctly formatted request for getting an access token', async () => {
                // Axiosの実装に依存するため、HttpClientが変わるごとに修正しないといけないUTです。
                // Arrange
                const cognitoDomain = 'https://example.com';
                const cognitoUserPoolURL = 'https://example-user-pool.com';
                const clientId = 'test-client-id';
                const clientSecret = 'test-client-secret';
                const redirectUri = 'https://callback.com/authorize/callback';
                const authorizationCode = 'valid-authorization-code';

                const loginAuthenticationService: ILoginAuthenticationService =
                    new CognitoOAuth2Service(cognitoDomain, cognitoUserPoolURL, clientId, clientSecret);

                // Act
                await loginAuthenticationService.getAccessToken(authorizationCode, redirectUri);

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
            });

            it('should propagate the exception to the caller in case of an error', async () => {
                // Arrange
                mockedAxios.post.mockRejectedValue(new Error('Network Error'));

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
            it('should successfully retrieve the public key for a given kid', async () => {
                // Arrange
                // https://dev.classmethod.jp/articles/azure-ad-id-token-verify-jwt-to-pem/からとってきた
                const jwksResponse = {
                    "keys": [
                        {
                            "kty": "RSA",
                            "use": "sig",
                            "kid": "T1St-dLTvyWRgxB_676u8krXS-I",
                            "x5t": "T1St-dLTvyWRgxB_676u8krXS-I",
                            "n": "s2TCRTB0HKEfLBPi3_8CxCbWirz7rlvzcXnp_0j3jrmb_hst0iiHifSBwE0FV1WW79Kyw0AATkLfSLLyllyCuzgoUOgmXd3YMaqB8mQOBIecFQDAHkM1syzi_VwVdJt8H1yI0hOGcOktujDPHidVFtOuoDqAWlCs7kCGwlazK4Sfu_pnfJI4RmU8AvqO0auGcxg24ICbpP01G0PgbvW8uhWSWSSTXmfdIh567JOHsgvFr0m1AUQv7wbeRxgyiHwn29h6g1bwSYJB4I6TMG-cDygvU9lNWFzeYhtqG4Z_cA3khWIMmTq3dVzCsi4iU309-c0FopWacTHouHyMRcpJFQ",
                            "e": "AQAB",
                            "x5c": [
                                "MIIC/TCCAeWgAwIBAgIIUd7j/OIahkYwDQYJKoZIhvcNAQELBQAwLTErMCkGA1UEAxMiYWNjb3VudHMuYWNjZXNzY29udHJvbC53aW5kb3dzLm5ldDAeFw0yMzExMDExNjAzMjdaFw0yODExMDExNjAzMjdaMC0xKzApBgNVBAMTImFjY291bnRzLmFjY2Vzc2NvbnRyb2wud2luZG93cy5uZXQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCzZMJFMHQcoR8sE+Lf/wLEJtaKvPuuW/Nxeen/SPeOuZv+Gy3SKIeJ9IHATQVXVZbv0rLDQABOQt9IsvKWXIK7OChQ6CZd3dgxqoHyZA4Eh5wVAMAeQzWzLOL9XBV0m3wfXIjSE4Zw6S26MM8eJ1UW066gOoBaUKzuQIbCVrMrhJ+7+md8kjhGZTwC+o7Rq4ZzGDbggJuk/TUbQ+Bu9by6FZJZJJNeZ90iHnrsk4eyC8WvSbUBRC/vBt5HGDKIfCfb2HqDVvBJgkHgjpMwb5wPKC9T2U1YXN5iG2obhn9wDeSFYgyZOrd1XMKyLiJTfT35zQWilZpxMei4fIxFykkVAgMBAAGjITAfMB0GA1UdDgQWBBRNcCE3HDX+HOJOu/bKfLYoSX3/0jANBgkqhkiG9w0BAQsFAAOCAQEAExns169MDr1dDNELYNK0JDjPUA6GR50jqfc+xa2KOljeXErOdihSvKgDS/vnDN6fjNNZuOMDyr6jjLvRsT0jVWzf/B6v92FrPRa/rv3urGXvW5am3BZyVPipirbiolMTuork95G7y7imftK7117uHcMq3D8f4fxscDiDXgjEEZqjkuzYDGLaVWGJqpv5xE4w+K4o2uDwmEIeIX+rI1MEVucS2vsvraOrjqjHwc3KrzuVRSsOU7YVHyUhku+7oOrB4tYrVbYYgwd6zXnkdouVPqOX9wTkc9iTmbDP+rfkhdadLxU+hmMyMuCJKgkZbWKFES7ce23jfTMbpqoHB4pgtQ=="
                            ]
                        },
                        {
                            "kty": "RSA",
                            "use": "sig",
                            "kid": "5B3nRxtQ7ji8eNDc3Fy05Kf97ZE",
                            "x5t": "5B3nRxtQ7ji8eNDc3Fy05Kf97ZE",
                            "n": "37fxbBQ8eAP7znqk8B8kUwVFEdV7N8WXSflojHJf9tNCyqrMd4gsAu4RUzhwlBlAHwYLhmwMNgRs-B4gLsEUaWJUjx5O4NxcfokC-6TL1p_IoLMGeqwFOMSBtxa1OnnL3eAQD1D7O0pOjJBst1_SZPswhdVzTEGoWod_1vMFDu02d3ogP_tuv2zl5Jd92t17Yuqb61wDKLzHoCeUMTVEzQS44n0mXrWJniaY0-_zs8opwwmUHRW3JJ_m0i5B1m9lFmcVUUKh5VYZPaec5ddCO47J91nZi0tcgADhqPzhAJF20cBEvkCP9kSJiiv2ssedlEbTSZGQTuC7OlP9G8tvvQ",
                            "e": "AQAB",
                            "x5c": [
                                "MIIC/TCCAeWgAwIBAgIISlx9oAuA2/MwDQYJKoZIhvcNAQELBQAwLTErMCkGA1UEAxMiYWNjb3VudHMuYWNjZXNzY29udHJvbC53aW5kb3dzLm5ldDAeFw0yMzEyMDUxNzE2NTdaFw0yODEyMDUxNzE2NTdaMC0xKzApBgNVBAMTImFjY291bnRzLmFjY2Vzc2NvbnRyb2wud2luZG93cy5uZXQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDft/FsFDx4A/vOeqTwHyRTBUUR1Xs3xZdJ+WiMcl/200LKqsx3iCwC7hFTOHCUGUAfBguGbAw2BGz4HiAuwRRpYlSPHk7g3Fx+iQL7pMvWn8igswZ6rAU4xIG3FrU6ecvd4BAPUPs7Sk6MkGy3X9Jk+zCF1XNMQahah3/W8wUO7TZ3eiA/+26/bOXkl33a3Xti6pvrXAMovMegJ5QxNUTNBLjifSZetYmeJpjT7/OzyinDCZQdFbckn+bSLkHWb2UWZxVRQqHlVhk9p5zl10I7jsn3WdmLS1yAAOGo/OEAkXbRwES+QI/2RImKK/ayx52URtNJkZBO4Ls6U/0by2+9AgMBAAGjITAfMB0GA1UdDgQWBBR6Y4Oi5GGItIomQ0yZfH/woCAogzANBgkqhkiG9w0BAQsFAAOCAQEAaNbWUtHv3+ryZecDc7m6V1V1rWrVkUwC2QO78a2TprEN3owOeP0IHP42fbd/wcSsufTTtkk/J+fqL5dtsQ6zk2kDQfY5CgOyVCsaxVqHsg3t8fAWBkHiNScjZvRhLx4ll9QMOtLAwL4Os3Of0qtvP61zONP9sCJoUB6hkB33SRma1OyPZnYK/l3r0Y49+Ov0wahcdI4yZI72hFXlyyLnOT8dMbJDwZ9LNXA/BauEff4qTI4nIQk/lQKS6BjHzvXZbkHYEV/6M7r1g1syeahDmnaII+ZiBwp6tmAZKZC0Q0O7y3DmcPrHiZdv35AHadZY5cGWy1rw8NIMkaHWZ0mP6Q=="
                            ]
                        },
                        {
                            "kty": "RSA",
                            "use": "sig",
                            "kid": "fwNs8F_h9KrHMA_OcTP1pGFWwyc",
                            "x5t": "fwNs8F_h9KrHMA_OcTP1pGFWwyc",
                            "n": "6Jiu4AU4ZWHBFEbO1-41P6dxKgGx7J31i5wNzH5eoJlsNjWrWoGlZip8ey_ZppcNMY0GY330p8YwdazqRX24mPkyOxbYF1uGEGB_XtmMOtuG45WyPlbARQl8hok7y_hbydS8uyfm_ZQXN7MLgju0f4_cYo-dgic5OaR3W6CWfgOrNnf287ZZ2HtJ8DZNm-oHE2_Tg9FFTIIkpltNIZ4rJ0uwzuy7zkep_Pfxptzmpd0rwd0F87IneYu-jtKUvHVVPJQ7yQvgin0rZR8tXIp_IzComGipktu_AJ89z3atOEt0_vZPizQIMRpToHjUTNXuXaDWIvCIJYMkvvl0HJxf1Q",
                            "e": "AQAB",
                            "x5c": [
                                "MIIC6TCCAdGgAwIBAgIIV6K/4n2M5VAwDQYJKoZIhvcNAQELBQAwIzEhMB8GA1UEAxMYbG9naW4ubWljcm9zb2Z0b25saW5lLnVzMB4XDTIzMTEzMDAwMTAxNVoXDTI4MTEzMDAwMTAxNVowIzEhMB8GA1UEAxMYbG9naW4ubWljcm9zb2Z0b25saW5lLnVzMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA6Jiu4AU4ZWHBFEbO1+41P6dxKgGx7J31i5wNzH5eoJlsNjWrWoGlZip8ey/ZppcNMY0GY330p8YwdazqRX24mPkyOxbYF1uGEGB/XtmMOtuG45WyPlbARQl8hok7y/hbydS8uyfm/ZQXN7MLgju0f4/cYo+dgic5OaR3W6CWfgOrNnf287ZZ2HtJ8DZNm+oHE2/Tg9FFTIIkpltNIZ4rJ0uwzuy7zkep/Pfxptzmpd0rwd0F87IneYu+jtKUvHVVPJQ7yQvgin0rZR8tXIp/IzComGipktu/AJ89z3atOEt0/vZPizQIMRpToHjUTNXuXaDWIvCIJYMkvvl0HJxf1QIDAQABoyEwHzAdBgNVHQ4EFgQUTtiYd3S6DOacXBmYsKyr1EK67f4wDQYJKoZIhvcNAQELBQADggEBAGbSzomDLsU7BX6Vohf/VweoJ9TgYs4cYcdFwMrRQVMpMGKYU6HT7f8mDzRGqpursuJTsN9yIOk7s5xp+N7EL6XKauzo+VHOGJT1qbwzJXT1XY6DuzBrlhtY9C7AUHlpYAD4uWyt+JfuB+z5Qq5cbGr0dMS/EEKR/m0iBboUNDji6u9sUzGqUYn4tpBoE+y0J8UttankG/09PPHwQIjMxMXcBDmGi5VTp9eY5RFk9GQ4qdQJUp2hhdQZDVpz6lcPxhG92RPO/ca3P/9dvfI5aNaiSyV7vuK2NGCVGCTeo/okA+V5dm5jeuf0bupNEPnXSGyM8EHjcRjR+cHsby5pIGs="
                            ]
                        }
                    ]
                };
                mockedAxios.get.mockResolvedValue({ data: jwksResponse });

                const cognitoDomain = 'https://example.com';
                const cognitoUserPoolURL = 'https://example-user-pool.com';
                const clientId = 'test-client-id';
                const clientSecret = 'test-client-secret';
                const loginAuthenticationService = new CognitoOAuth2Service(
                    cognitoDomain, cognitoUserPoolURL, clientId, clientSecret
                );

                // Act
                const publicKey = await loginAuthenticationService.getPublicKey('T1St-dLTvyWRgxB_676u8krXS-I');

                // Assert
                const expectedPublicKey =
                    '-----BEGIN PUBLIC KEY-----\n' +
                    'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAs2TCRTB0HKEfLBPi3/8C\n' +
                    'xCbWirz7rlvzcXnp/0j3jrmb/hst0iiHifSBwE0FV1WW79Kyw0AATkLfSLLyllyC\n' +
                    'uzgoUOgmXd3YMaqB8mQOBIecFQDAHkM1syzi/VwVdJt8H1yI0hOGcOktujDPHidV\n' +
                    'FtOuoDqAWlCs7kCGwlazK4Sfu/pnfJI4RmU8AvqO0auGcxg24ICbpP01G0PgbvW8\n' +
                    'uhWSWSSTXmfdIh567JOHsgvFr0m1AUQv7wbeRxgyiHwn29h6g1bwSYJB4I6TMG+c\n' +
                    'DygvU9lNWFzeYhtqG4Z/cA3khWIMmTq3dVzCsi4iU309+c0FopWacTHouHyMRcpJ\n' +
                    'FQIDAQAB\n' +
                    '-----END PUBLIC KEY-----\n';
                expect(publicKey).toBe(expectedPublicKey);
            });

            it('should return null if no key matches the given kid', async () => {
                // Arrange
                // https://dev.classmethod.jp/articles/azure-ad-id-token-verify-jwt-to-pem/からとってきた
                const jwksResponse = {
                    "keys": [
                        {
                            "kty": "RSA",
                            "use": "sig",
                            "kid": "T1St-dLTvyWRgxB_676u8krXS-I",
                            "x5t": "T1St-dLTvyWRgxB_676u8krXS-I",
                            "n": "s2TCRTB0HKEfLBPi3_8CxCbWirz7rlvzcXnp_0j3jrmb_hst0iiHifSBwE0FV1WW79Kyw0AATkLfSLLyllyCuzgoUOgmXd3YMaqB8mQOBIecFQDAHkM1syzi_VwVdJt8H1yI0hOGcOktujDPHidVFtOuoDqAWlCs7kCGwlazK4Sfu_pnfJI4RmU8AvqO0auGcxg24ICbpP01G0PgbvW8uhWSWSSTXmfdIh567JOHsgvFr0m1AUQv7wbeRxgyiHwn29h6g1bwSYJB4I6TMG-cDygvU9lNWFzeYhtqG4Z_cA3khWIMmTq3dVzCsi4iU309-c0FopWacTHouHyMRcpJFQ",
                            "e": "AQAB",
                            "x5c": [
                                "MIIC/TCCAeWgAwIBAgIIUd7j/OIahkYwDQYJKoZIhvcNAQELBQAwLTErMCkGA1UEAxMiYWNjb3VudHMuYWNjZXNzY29udHJvbC53aW5kb3dzLm5ldDAeFw0yMzExMDExNjAzMjdaFw0yODExMDExNjAzMjdaMC0xKzApBgNVBAMTImFjY291bnRzLmFjY2Vzc2NvbnRyb2wud2luZG93cy5uZXQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCzZMJFMHQcoR8sE+Lf/wLEJtaKvPuuW/Nxeen/SPeOuZv+Gy3SKIeJ9IHATQVXVZbv0rLDQABOQt9IsvKWXIK7OChQ6CZd3dgxqoHyZA4Eh5wVAMAeQzWzLOL9XBV0m3wfXIjSE4Zw6S26MM8eJ1UW066gOoBaUKzuQIbCVrMrhJ+7+md8kjhGZTwC+o7Rq4ZzGDbggJuk/TUbQ+Bu9by6FZJZJJNeZ90iHnrsk4eyC8WvSbUBRC/vBt5HGDKIfCfb2HqDVvBJgkHgjpMwb5wPKC9T2U1YXN5iG2obhn9wDeSFYgyZOrd1XMKyLiJTfT35zQWilZpxMei4fIxFykkVAgMBAAGjITAfMB0GA1UdDgQWBBRNcCE3HDX+HOJOu/bKfLYoSX3/0jANBgkqhkiG9w0BAQsFAAOCAQEAExns169MDr1dDNELYNK0JDjPUA6GR50jqfc+xa2KOljeXErOdihSvKgDS/vnDN6fjNNZuOMDyr6jjLvRsT0jVWzf/B6v92FrPRa/rv3urGXvW5am3BZyVPipirbiolMTuork95G7y7imftK7117uHcMq3D8f4fxscDiDXgjEEZqjkuzYDGLaVWGJqpv5xE4w+K4o2uDwmEIeIX+rI1MEVucS2vsvraOrjqjHwc3KrzuVRSsOU7YVHyUhku+7oOrB4tYrVbYYgwd6zXnkdouVPqOX9wTkc9iTmbDP+rfkhdadLxU+hmMyMuCJKgkZbWKFES7ce23jfTMbpqoHB4pgtQ=="
                            ]
                        },
                        {
                            "kty": "RSA",
                            "use": "sig",
                            "kid": "5B3nRxtQ7ji8eNDc3Fy05Kf97ZE",
                            "x5t": "5B3nRxtQ7ji8eNDc3Fy05Kf97ZE",
                            "n": "37fxbBQ8eAP7znqk8B8kUwVFEdV7N8WXSflojHJf9tNCyqrMd4gsAu4RUzhwlBlAHwYLhmwMNgRs-B4gLsEUaWJUjx5O4NxcfokC-6TL1p_IoLMGeqwFOMSBtxa1OnnL3eAQD1D7O0pOjJBst1_SZPswhdVzTEGoWod_1vMFDu02d3ogP_tuv2zl5Jd92t17Yuqb61wDKLzHoCeUMTVEzQS44n0mXrWJniaY0-_zs8opwwmUHRW3JJ_m0i5B1m9lFmcVUUKh5VYZPaec5ddCO47J91nZi0tcgADhqPzhAJF20cBEvkCP9kSJiiv2ssedlEbTSZGQTuC7OlP9G8tvvQ",
                            "e": "AQAB",
                            "x5c": [
                                "MIIC/TCCAeWgAwIBAgIISlx9oAuA2/MwDQYJKoZIhvcNAQELBQAwLTErMCkGA1UEAxMiYWNjb3VudHMuYWNjZXNzY29udHJvbC53aW5kb3dzLm5ldDAeFw0yMzEyMDUxNzE2NTdaFw0yODEyMDUxNzE2NTdaMC0xKzApBgNVBAMTImFjY291bnRzLmFjY2Vzc2NvbnRyb2wud2luZG93cy5uZXQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDft/FsFDx4A/vOeqTwHyRTBUUR1Xs3xZdJ+WiMcl/200LKqsx3iCwC7hFTOHCUGUAfBguGbAw2BGz4HiAuwRRpYlSPHk7g3Fx+iQL7pMvWn8igswZ6rAU4xIG3FrU6ecvd4BAPUPs7Sk6MkGy3X9Jk+zCF1XNMQahah3/W8wUO7TZ3eiA/+26/bOXkl33a3Xti6pvrXAMovMegJ5QxNUTNBLjifSZetYmeJpjT7/OzyinDCZQdFbckn+bSLkHWb2UWZxVRQqHlVhk9p5zl10I7jsn3WdmLS1yAAOGo/OEAkXbRwES+QI/2RImKK/ayx52URtNJkZBO4Ls6U/0by2+9AgMBAAGjITAfMB0GA1UdDgQWBBR6Y4Oi5GGItIomQ0yZfH/woCAogzANBgkqhkiG9w0BAQsFAAOCAQEAaNbWUtHv3+ryZecDc7m6V1V1rWrVkUwC2QO78a2TprEN3owOeP0IHP42fbd/wcSsufTTtkk/J+fqL5dtsQ6zk2kDQfY5CgOyVCsaxVqHsg3t8fAWBkHiNScjZvRhLx4ll9QMOtLAwL4Os3Of0qtvP61zONP9sCJoUB6hkB33SRma1OyPZnYK/l3r0Y49+Ov0wahcdI4yZI72hFXlyyLnOT8dMbJDwZ9LNXA/BauEff4qTI4nIQk/lQKS6BjHzvXZbkHYEV/6M7r1g1syeahDmnaII+ZiBwp6tmAZKZC0Q0O7y3DmcPrHiZdv35AHadZY5cGWy1rw8NIMkaHWZ0mP6Q=="
                            ]
                        },
                        {
                            "kty": "RSA",
                            "use": "sig",
                            "kid": "fwNs8F_h9KrHMA_OcTP1pGFWwyc",
                            "x5t": "fwNs8F_h9KrHMA_OcTP1pGFWwyc",
                            "n": "6Jiu4AU4ZWHBFEbO1-41P6dxKgGx7J31i5wNzH5eoJlsNjWrWoGlZip8ey_ZppcNMY0GY330p8YwdazqRX24mPkyOxbYF1uGEGB_XtmMOtuG45WyPlbARQl8hok7y_hbydS8uyfm_ZQXN7MLgju0f4_cYo-dgic5OaR3W6CWfgOrNnf287ZZ2HtJ8DZNm-oHE2_Tg9FFTIIkpltNIZ4rJ0uwzuy7zkep_Pfxptzmpd0rwd0F87IneYu-jtKUvHVVPJQ7yQvgin0rZR8tXIp_IzComGipktu_AJ89z3atOEt0_vZPizQIMRpToHjUTNXuXaDWIvCIJYMkvvl0HJxf1Q",
                            "e": "AQAB",
                            "x5c": [
                                "MIIC6TCCAdGgAwIBAgIIV6K/x5c+41P6dxKgGx7J31i5wNzH5eoJlsNjWrWoGlZip8ey/ZppcNMY0GY330p8YwdazqRX24mPkyOxbYF1uGEGB/XtmMOtuG45WyPlbARQl8hok7y/hbydS8uyfm/ZQXN7MLgju0f4/cYo+dgic5OaR3W6CWfgOrNnf287ZZ2HtJ8DZNm+oHE2/Tg9FFTIIkpltNIZ4rJ0uwzuy7zkep/Pfxptzmpd0rwd0F87IneYu+jtKUvHVVPJQ7yQvgin0rZR8tXIp/IzComGipktu/AJ89z3atOEt0/vZPizQIMRpToHjUTNXuXaDWIvCIJYMkvvl0HJxf1QIDAQABoyEwHzAdBgNVHQ4EFgQUTtiYd3S6DOacXBmYsKyr1EK67f4wDQYJKoZIhvcNAQELBQADggEBAGbSzomDLsU7BX6Vohf/VweoJ9TgYs4cYcdFwMrRQVMpMGKYU6HT7f8mDzRGqpursuJTsN9yIOk7s5xp+N7EL6XKauzo+VHOGJT1qbwzJXT1XY6DuzBrlhtY9C7AUHlpYAD4uWyt+JfuB+z5Qq5cbGr0dMS/EEKR/m0iBboUNDji6u9sUzGqUYn4tpBoE+y0J8UttankG/09PPHwQIjMxMXcBDmGi5VTp9eY5RFk9GQ4qdQJUp2hhdQZDVpz6lcPxhG92RPO/ca3P/9dvfI5aNaiSyV7vuK2NGCVGCTeo/okA+V5dm5jeuf0bupNEPnXSGyM8EHjcRjR+cHsby5pIGs="
                            ]
                        }
                    ]
                };
                mockedAxios.get.mockResolvedValue({ data: jwksResponse });

                const cognitoDomain = 'https://example.com';
                const cognitoUserPoolURL = 'https://example-user-pool.com';
                const clientId = 'test-client-id';
                const clientSecret = 'test-client-secret';
                const loginAuthenticationService = new CognitoOAuth2Service(
                    cognitoDomain, cognitoUserPoolURL, clientId, clientSecret
                );

                // Act
                const publicKey = await loginAuthenticationService.getPublicKey('kid-not-registered');

                // Assert
                expect(publicKey).toBeNull();
            });

            it('should return null if the JWKS response is empty', async () => {
                // Arrange
                const emptyJwksResponse = {
                    keys: [] // JWKS response with an empty keys array
                };
                mockedAxios.get.mockResolvedValue({ data: emptyJwksResponse });

                const cognitoDomain = 'https://example.com';
                const cognitoUserPoolURL = 'https://example-user-pool.com';
                const clientId = 'test-client-id';
                const clientSecret = 'test-client-secret';
                const loginAuthenticationService = new CognitoOAuth2Service(
                    cognitoDomain, cognitoUserPoolURL, clientId, clientSecret
                );

                // Act
                const publicKey = await loginAuthenticationService.getPublicKey('any-kid');

                // Assert
                expect(publicKey).toBeNull();
            });

            it('should propagate the exception to the caller in case of an error', async () => {
                // Arrange
                mockedAxios.get.mockRejectedValue(new Error('Network Error'));

                const loginAuthenticationService: ILoginAuthenticationService =
                    new CognitoOAuth2Service("", "", "", "");

                // Act&Assert
                await expect(loginAuthenticationService.getPublicKey('any-kid'))
                    .rejects
                    .toThrow('Network Error');
                expect(mockedAxios.get).toHaveBeenCalledTimes(1);
            });
        })

        describe('refreshToken', () => {
            it('should retrieve a new set of tokens given a valid refresh token', async () => {
                // Arrange
                const cognitoDomain = 'https://example.com';
                const cognitoUserPoolURL = 'https://example-user-pool.com';
                const clientId = 'test-client-id';
                const clientSecret = 'test-client-secret';
                const refreshToken = 'valid-refresh-token';

                mockedAxios.post.mockResolvedValue({
                    data: {
                        access_token: 'new-access-token',
                        id_token: 'new-id-token',
                        refresh_token: 'new-refresh-token',
                        expires_in: 3600,
                        token_type: 'Bearer'
                    }
                });

                const loginAuthenticationService: ILoginAuthenticationService =
                    new CognitoOAuth2Service(cognitoDomain, cognitoUserPoolURL, clientId, clientSecret);

                // Act
                const actual = await loginAuthenticationService.refreshToken(refreshToken);

                // Assert
                expect(actual).toEqual({
                    access_token: 'new-access-token',
                    id_token: 'new-id-token',
                    refresh_token: 'new-refresh-token',
                    expires_in: 3600,
                    token_type: 'Bearer'
                });
            });

            it('should correctly send a request to refresh the token', async () => {
                // Axiosの実装に依存するため、HttpClientが変わるごとに修正しないといけないUTです。
                // Arrange
                const cognitoDomain = 'https://example.com';
                const cognitoUserPoolURL = 'https://example-user-pool.com';
                const clientId = 'test-client-id';
                const clientSecret = 'test-client-secret';
                const refreshToken = 'valid-refresh-token';

                const loginAuthenticationService: ILoginAuthenticationService =
                    new CognitoOAuth2Service(cognitoDomain, cognitoUserPoolURL, clientId, clientSecret);

                // Act
                await loginAuthenticationService.refreshToken(refreshToken);

                // Assert
                expect(mockedAxios.post).toHaveBeenCalledWith(
                    `/oauth2/token`,
                    new URLSearchParams({
                        grant_type: 'refresh_token',
                        client_id: 'test-client-id',
                        refresh_token: 'valid-refresh-token'
                    }),
                    {
                        headers: {
                            'Authorization': 'Basic dGVzdC1jbGllbnQtaWQ6dGVzdC1jbGllbnQtc2VjcmV0',
                            'Content-Type': 'application/x-www-form-urlencoded'
                        }
                    }
                );
            });

            it('should propagate the exception to the caller in case of an error', async () => {
                // Arrange
                mockedAxios.post.mockRejectedValue(new Error('Network Error'));

                const loginAuthenticationService: ILoginAuthenticationService =
                    new CognitoOAuth2Service("", "", "", "");

                // Act & Assert
                await expect(loginAuthenticationService.refreshToken("invalid-token"))
                    .rejects
                    .toThrow('Network Error');
                expect(mockedAxios.post).toHaveBeenCalledTimes(1);
            });
        })
    })
})