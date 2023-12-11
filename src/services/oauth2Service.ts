import axios, { AxiosInstance } from "axios";
import { singleton } from "tsyringe";
import jwkToPem from "jwk-to-pem";

export interface OAuthTokenResponse {
    id_token: string;
    access_token: string;
    refresh_token: string;
    expires_in: number;
    token_type: string;
}
export interface IOAuth2Service<T extends OAuthTokenResponse> {
    getAuthorizationUrl(callbackUrl: string): string;
    getAccessToken(authorizationCode: string, redirect_uri: string): Promise<T>;
    getPublicKey(kid: string): Promise<string | null>;
    refreshToken(refreshToken: string): Promise<T>;
}

export interface ILoginAuthenticationService extends IOAuth2Service<OAuthTokenResponse> {
}

@singleton()
export class CognitoOAuth2Service implements ILoginAuthenticationService {
    private httpClient: AxiosInstance;
    private cognitoIdpHttpClient: AxiosInstance;

    private basicBase64EncodeCredentials: string;

    constructor(private cognitoDomain: string, private cognitoUserPoolURL: string, private clientId: string, private clientSecret: string) {
        this.httpClient = axios.create({
            baseURL: cognitoDomain
        });
        this.cognitoIdpHttpClient = axios.create({
            baseURL: cognitoUserPoolURL
        });
        this.basicBase64EncodeCredentials = Buffer.from(`${clientId}:${clientSecret}`).toString('base64');
    }

    public getAuthorizationUrl(callbackUrl: string): string {
        return `${this.cognitoDomain}/oauth2/authorize?response_type=code&client_id=${this.clientId}&redirect_uri=${callbackUrl}`
    }

    public async getAccessToken(authorizationCode: string, redirectUri: string): Promise<OAuthTokenResponse> {
        const response = await this.httpClient.post(
            `/oauth2/token`,
            new URLSearchParams({
                grant_type: 'authorization_code',
                client_id: this.clientId,
                redirect_uri: redirectUri,
                code: authorizationCode
            }),
            {
                headers: {
                    'Authorization': `Basic ${this.basicBase64EncodeCredentials}`,
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            }
        );
        return response.data;
    }

    public async getPublicKey(kid: string): Promise<string | null> {
        const response = await this.cognitoIdpHttpClient.get(
            `/.well-known/jwks.json`
        );
        const jwks = response.data.keys;

        const key = jwks.find((k: any) => k.kid === kid);
        if (!key) {
            return null;
        }

        return jwkToPem(key);
    }

    public async refreshToken(refreshToken: string): Promise<OAuthTokenResponse> {
        const response = await this.httpClient.post(
            `/oauth2/token`,
            new URLSearchParams({
                grant_type: 'refresh_token',
                client_id: this.clientId,
                refresh_token: refreshToken
            }),
            {
                headers: {
                    'Authorization': `Basic ${this.basicBase64EncodeCredentials}`,
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            }
        );
        return response.data;
    }
}
