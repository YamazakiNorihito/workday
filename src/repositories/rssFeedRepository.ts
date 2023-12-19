import { SchemaFieldTypes, RedisClientType } from 'redis';
import { inject, singleton } from "tsyringe";

export interface IRSSFeedRepository {
    save(feedId: string, feed: RSSFeed): Promise<void>;
    get(feedId: string): Promise<RSSFeed | null>;
}

@singleton()
export class RSSFeedRepository implements IRSSFeedRepository {
    constructor(
        @inject("RedisClient") private readonly redisClient: RedisClientType
    ) {
        (async () => {
            await this.createIndex();
        })();
    }

    private async createIndex(): Promise<void> {
        try {
            await this.redisClient.ft.create('idx:rssfeed', {
                '$.title': { type: SchemaFieldTypes.TEXT },
                '$.link': { type: SchemaFieldTypes.TEXT },
                '$.pubDate': { type: SchemaFieldTypes.TEXT },
            }, {
                ON: 'JSON',
                PREFIX: 'rssfeed:'
            });
        } catch (e: any) {
            if (e.message !== 'Index already exists') {
                throw e;
            }
        }
    }

    public async save(feedId: string, feed: RSSFeed): Promise<void> {
        const feedKey = `rssfeed:${feedId}`;
        let feedJsonString = JSON.stringify(feed);
        await this.redisClient.json.set(feedKey, '$', feedJsonString);
    }

    public async get(feedId: string): Promise<RSSFeed | null> {
        const feedKey = `rssfeed:${feedId}`;
        const document = await this.redisClient.json.get(feedKey);

        const json = (document as unknown) as string;

        if (!json) {
            return null;
        }
        return JSON.parse(json as string) as RSSFeed;
    }
}

export interface RSSFeedItem {
    title?: string;
    link?: string;
    pubDate?: Date;
    description?: string;
    contentSnippet?: string;
    categories?: string[];
}

export interface RSSFeed {
    title?: string;
    description?: string;
    link?: string;
    lastBuildDate?: Date;
    items: RSSFeedItem[]
}