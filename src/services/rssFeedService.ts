import { inject, singleton } from "tsyringe";
import Parser from 'rss-parser';
import { createHash } from "crypto";
import { IRSSFeedRepository, RSSFeed, RSSFeedItem, RSSFeedRepository } from "../repositories/rssFeedRepository";

export interface IRSSFeedService {
    getRSSFeed(feedUrl: string): Promise<RSSFeed | null>;
    getAndSaveRSSFeed(feedUrl: string): Promise<void>;
}

@singleton()
export class RSSFeedService {

    constructor(@inject(RSSFeedRepository) private readonly _rssFeedRepository: IRSSFeedRepository) { }

    public async getRSSFeed(feedUrl: string): Promise<RSSFeed | null> {

        const feedId = this.generateIdFromUrl(feedUrl);

        const feedRss = await this._rssFeedRepository.get(feedId);

        return feedRss;

        /*
        let parser = new Parser();
        const feed = await parser.parseURL(feedUrl);

        const rssFeedItems: RSSFeedItem[] = feed.items.map(item => ({
            title: item.title,
            link: item.link,
            pubDate: item.pubDate,
            contentSnippet: item.contentSnippet,
            categories: item.categories
        }));

        const feedRss: RSSFeed = {
            title: feed.title,
            description: feed.description,
            link: feed.link,
            lastBuildDate: new Date(feed.lastBuildDate),
            items: rssFeedItems
        }

        return feedRss;*/
        //console.log("feed", feed);
        //const item = feed.items[0];
        //console.log("Title:", item.title);
        //console.log("Link:", item.link);
        //console.log("Publication Date:", item.pubDate);
        //console.log("Creator:", item.creator);
        //console.log("Content:", item['content:encoded']);
        //console.log("Content Snippet:", item['content:encodedSnippet']);
        //console.log("DC Creator:", item['dc:creator']);
        //console.log("Comments Link:", item.comments);
        //console.log("Simple Content:", item.content);
        //console.log("Simple Content Snippet:", item.contentSnippet);
        //console.log("GUID:", item.guid);
        //console.log("Categories:", item.categories);
        //console.log("ISO Date:", item.isoDate);
    }

    public async getAndSaveRSSFeed(feedUrl: string): Promise<void> {
        let parser = new Parser();
        const feed = await parser.parseURL(feedUrl);

        const rssFeedItems: RSSFeedItem[] = feed.items.map(item => ({
            title: item.title,
            link: item.link,
            pubDate: item.pubDate ? new Date(item.pubDate) : undefined,
            contentSnippet: item.contentSnippet,
            categories: item.categories
        }));

        const rssFeed: RSSFeed = {
            title: feed.title,
            description: feed.description,
            link: feed.link,
            lastBuildDate: new Date(feed.lastBuildDate),
            items: rssFeedItems
        }

        const feedId = this.generateIdFromUrl(feedUrl);

        await this._rssFeedRepository.save(feedId, rssFeed);
    }

    private generateIdFromUrl(url: string): string {
        return createHash('sha256').update(url).digest('hex');
    }
}