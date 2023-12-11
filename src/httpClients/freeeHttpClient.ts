import axios, { AxiosInstance } from 'axios';
import { singleton } from 'tsyringe';

@singleton()
export class FreeeHrHttpApiClient {
    private baseURL = 'https://api.freee.co.jp/hr';
    private httpClient: AxiosInstance;

    constructor() {
        this.httpClient = axios.create({
            baseURL: this.baseURL
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

    public async get<T = any>(path: string, accessToken: string): Promise<T> {
        const response = await this.httpClient.get(
            `${path}`,
            {
                headers: {
                    accept: 'application/json',
                    Authorization: `Bearer ${accessToken}`,
                },
            }
        );
        return response.data;
    }

    public async post<T = any>(path: string, accessToken: string, data: any): Promise<T> {
        const response = await this.httpClient.post(
            `${path}`,
            data,
            {
                headers: {
                    'accept': 'application/json',
                    'content-type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`,
                },
            }
        );
        return response.data;
    }

    public async put<T = any>(path: string, accessToken: string, data: any): Promise<T> {
        const response = await this.httpClient.put(
            `${path}`,
            data,
            {
                headers: {
                    'accept': 'application/json',
                    'content-type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`,
                },
            }
        );
        return response.data;
    }

    public async delete<T = any>(path: string, accessToken: string): Promise<T> {
        const response = await this.httpClient.delete(
            `${path}`,
            {
                headers: {
                    'accept': 'application/json',
                    'content-type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`,
                },
            }
        );
        return response.data;
    }
}