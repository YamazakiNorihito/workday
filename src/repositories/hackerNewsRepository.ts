import { SchemaFieldTypes, RedisClientType } from 'redis';
import { inject, singleton } from "tsyringe";

export interface IHackerNewsRepository {
    save(itemId: number, item: HackerNewsItem): Promise<void>;
    get(itemId: number): Promise<HackerNewsItem | null>;
}
@singleton()
export class HackerNewsRepository implements IHackerNewsRepository {
    constructor(
        @inject("RedisClient") private readonly redisClient: RedisClientType
    ) {
        (async () => {
            await this.createIndex();
        })();
    }

    private async createIndex(): Promise<void> {
        try {
            await this.redisClient.ft.create('idx:hackernews', {
                '$.id': { type: SchemaFieldTypes.NUMERIC, SORTABLE: true },
                '$.type': { type: SchemaFieldTypes.TEXT },
                '$.by': { type: SchemaFieldTypes.TEXT },
                '$.time': { type: SchemaFieldTypes.NUMERIC, SORTABLE: true },
                '$.deleted': { type: SchemaFieldTypes.NUMERIC },
                '$.dead': { type: SchemaFieldTypes.NUMERIC },
                '$.kids[*]': { type: SchemaFieldTypes.NUMERIC },
                '$.url': { type: SchemaFieldTypes.TEXT },
                '$.score': { type: SchemaFieldTypes.NUMERIC },
                '$.title': { type: SchemaFieldTypes.TEXT },
                '$.descendants': { type: SchemaFieldTypes.NUMERIC },
                '$.text': { type: SchemaFieldTypes.TEXT },
                '$.parent': { type: SchemaFieldTypes.NUMERIC },
                '$.parts[*]': { type: SchemaFieldTypes.NUMERIC }
            }, {
                ON: 'JSON',
                PREFIX: 'hackernews:'
            });
        } catch (e: any) {
            if (e.message !== 'Index already exists') {
                throw e;
            }
        }
    }

    public async save(itemId: number, item: HackerNewsItem): Promise<void> {
        const itemKey = `hackernews:${itemId}`;
        let itemJsonString = JSON.stringify(item);
        const result = await this.redisClient.json.set(itemKey, '$', itemJsonString);
    }

    public async get(itemId: number): Promise<HackerNewsItem | null> {
        const itemKey = `hackernews:${itemId}`;
        const document = await this.redisClient.json.get(itemKey);

        const json = (document as unknown) as string;

        if (!json) {
            return null;
        }
        return JSON.parse(json as string) as HackerNewsItem
    }
}

export type ItemType = 'story' | 'comment' | 'job' | 'poll' | 'pollopt';

export type HackerNewsItem = StoryHackerNewsItem | JobHackerNewsItem | CommentHackerNewsItem | PollHackerNewsItem | PollOptionHackerNewsItem;

export interface BaseHackerNewsItem {
    id: number;
    deleted?: boolean;
    type: ItemType;
    by: string;
    time: number;
    dead?: boolean;
    kids?: number[];
}

export interface StoryHackerNewsItem extends BaseHackerNewsItem {
    url: string;
    score: number;
    title: string;
    descendants: number;
}

export interface JobHackerNewsItem extends BaseHackerNewsItem {
    text: string;
    url: string;
    title: string;
}

export interface PollHackerNewsItem extends BaseHackerNewsItem {
    text: string;
    score: number;
    title: string;
    parts: number[];
    descendants: number;
}

export interface CommentHackerNewsItem extends BaseHackerNewsItem {
    text: string;
    parent: number;
}

export interface PollOptionHackerNewsItem extends BaseHackerNewsItem {
    parent: number;
    score: number;
}