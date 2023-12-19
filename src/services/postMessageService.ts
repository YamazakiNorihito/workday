import { inject, singleton } from "tsyringe";
import { SlackHttpApiClient } from "../httpClients/slackHttpClient";
import { RSSFeed } from "../repositories/rssFeedRepository";

export interface IPostMessageService {
    postFeedToSlack(feed: RSSFeed, channelId: string, fromName: string): Promise<void>;
}

@singleton()
export class PostMessageService implements IPostMessageService {
    constructor(
        @inject(SlackHttpApiClient) private slackHttpApiClient: SlackHttpApiClient
    ) { }

    public async postFeedToSlack(feed: RSSFeed, channelId: string, fromName: string): Promise<void> {
        const message = this.formatFeedMessage(feed);
        await this.slackHttpApiClient.post(
            '/chat.postMessage', {
            channel: channelId,
            username: fromName,
            text: message
        });
    }

    private formatFeedMessage(feed: RSSFeed): string {
        let message = `*フィードタイトル:* <${feed.link}|${feed.title}>\n*フィード詳細:* ${feed.description}`;
        if (feed.lastBuildDate) {
            const lastBuildDate = new Date(feed.lastBuildDate);
            message += `\n*最終更新日:* ${lastBuildDate.toISOString()}`;
        }
        const view_length = 74 * 2;
        message += `\n\n*最新の記事:*\n`;
        feed.items.forEach((item, index) => {
            const truncatedDescription = item.contentSnippet?.substring(0, view_length) + (item.contentSnippet && item.contentSnippet.length > view_length ? '...' : '');
            message += `${index + 1}. *記事タイトル:* <${item.link}|${item.title}>\n    *公開日:* ${item.pubDate}\n    *概要:* ${truncatedDescription}\n    *カテゴリ:* ${item.categories?.join(", ")}\n\n`;
        });

        return message;
    }
}