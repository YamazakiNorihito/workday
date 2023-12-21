import axios, { AxiosInstance } from "axios";
import { inject, singleton } from "tsyringe";
import { BaseHackerNewsItem, HackerNewsItem, HackerNewsRepository, IHackerNewsRepository } from "../repositories/hackerNewsRepository";

export interface IHackerNewsService {
    getTopStories(): Promise<BaseHackerNewsItem[]>;
    getNewStories(): Promise<BaseHackerNewsItem[]>;
    getBestStories(): Promise<BaseHackerNewsItem[]>;
    getAskHNStories(): Promise<BaseHackerNewsItem[]>;
    getShowHNStories(): Promise<BaseHackerNewsItem[]>;
    getJobStories(): Promise<BaseHackerNewsItem[]>;
    getMaxItem(): Promise<number>;
    getItem(itemId: number): Promise<HackerNewsItem>;
}

@singleton()
export class HackerNewsService implements IHackerNewsService {
    protected hackerNewsHttpClient: AxiosInstance;

    constructor(
        @inject(HackerNewsRepository) private readonly _hackerNewsRepository: IHackerNewsRepository) {
        this.hackerNewsHttpClient = axios.create({
            baseURL: 'https://hacker-news.firebaseio.com/v0'
        });
    }
    public async getTopStories(): Promise<BaseHackerNewsItem[]> {
        return this.getStories('/topstories.json');
    }

    public async getNewStories(): Promise<BaseHackerNewsItem[]> {
        return this.getStories('/newstories.json');
    }

    public async getBestStories(): Promise<BaseHackerNewsItem[]> {
        return this.getStories('/beststories.json');
    }

    public async getAskHNStories(): Promise<BaseHackerNewsItem[]> {
        return this.getStories('/askstories.json');
    }

    public async getShowHNStories(): Promise<BaseHackerNewsItem[]> {
        return this.getStories('/showstories.json');
    }

    public async getJobStories(): Promise<BaseHackerNewsItem[]> {
        return this.getStories('/jobstories.json');
    }

    public async getMaxItem(): Promise<number> {
        const response = await this.hackerNewsHttpClient.get('/maxitem.json');
        return response.data;
    }

    public async getItem(itemId: number): Promise<HackerNewsItem> {
        const response = await this.hackerNewsHttpClient.get<HackerNewsItem>(`/item/${itemId}.json`);
        return response.data;
    }

    private async getStories(endpoint: string): Promise<BaseHackerNewsItem[]> {
        const response = await this.hackerNewsHttpClient.get(endpoint);
        let itemIds = response.data as number[];

        const items = await Promise.all(itemIds.map(itemId => this._hackerNewsRepository.get(itemId)));
        const missingItemIds = itemIds.filter((id, index) => !items[index]);
        await Promise.all(missingItemIds.map(async itemId => {
            try {
                const fetchNews = await this.getItem(itemId);
                await this._hackerNewsRepository.save(fetchNews.id, fetchNews);
                return fetchNews;
            } catch (error) {
                console.error(`Error fetching item ${itemId}:`, error);
                return null;
            }
        })).then(newItems => {
            newItems.forEach(item => {
                if (item) items[itemIds.indexOf(item.id)] = item;
            });
        });

        return items as HackerNewsItem[];
    }
}