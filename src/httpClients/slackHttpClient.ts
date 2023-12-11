import axios, { AxiosInstance } from 'axios';
import { singleton } from 'tsyringe';

@singleton()
export class SlackHttpApiClient {
    private static baseURL = 'https://slack.com/api';
    private httpClient: AxiosInstance;
    private readonly _accessToken: string;

    constructor(accessToken: string) {
        this._accessToken = accessToken;
        this.httpClient = axios.create({
            baseURL: SlackHttpApiClient.baseURL
        });

        this.httpClient.interceptors.request.use((request) => {
            return request;
        });
        this.httpClient.interceptors.response.use((response) => {
            return response;
        }, (error) => {
            console.log('Response Error:', error);
            if (error.response) {
                console.log('Error Response Data:', JSON.stringify(error.response.data, null, 2));
            }
            return Promise.reject(error);
        });
    }

    public async post<T = any>(path: string, data: any): Promise<T> {
        const response = await this.httpClient.post(
            path,
            data,
            {
                headers: {
                    'accept': 'application/json',
                    'content-type': 'application/json',
                    'Authorization': `Bearer ${this._accessToken}`,
                },
            }
        );
        return response.data;
    }
}