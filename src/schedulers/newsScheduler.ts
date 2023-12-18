import { CronJob } from 'cron';
import { inject, singleton } from 'tsyringe';
import { IRSSFeedService, RSSFeedService } from '../services/rssFeedService';
import { PostMessageService } from '../services/postMessageService';
import { RSSFeed } from '../repositories/rssFeedRepository';

@singleton()
export class NewsScheduler {
    private fetchJob: CronJob;
    private notificationJob: CronJob;

    constructor(@inject(RSSFeedService) private readonly _rssFeedService: IRSSFeedService,
        @inject(PostMessageService) private readonly _postMessageService: PostMessageService) {
        this.fetchJob = new CronJob(
            '*/20 * * * *',
            () => {
                console.log('Starting to fetch all feeds...');
                this.fetchAllFeeds();
            },
            null,
            false,
            'Asia/Tokyo'
        );
        this.notificationJob = new CronJob(
            '0 * * * *',
            () => {
                console.log('Starting to fetch all feeds...');
                this.notificationAllFeeds();
            },
            null,
            false,
            'Asia/Tokyo'
        );
    }

    start() {
        console.log('NewsScheduler started.');
        this.fetchJob.start();
        this.notificationJob.start();
    }

    private async fetchAllFeeds() {
        console.log('Starting to fetch RSS feeds...');
        for (const category in categoryFeeds) {
            console.log(`Fetching feed for category: ${category}`);
            const categoryFeed = categoryFeeds[category];
            if (!categoryFeed) {
                console.log(`Category not found: ${category}`);
                return;
            }

            const feedUrl = categoryFeed.url;
            try {
                console.log(`Fetching RSS feed from URL: ${feedUrl}`);
                await this._rssFeedService.getAndSaveRSSFeed(feedUrl);
                console.log(`RSS feed fetched and saved for category: ${category}`);
            } catch (error) {
                console.error(`Error while fetching RSS feed for category '${category}':`, error);
            }
        }
        console.log('All RSS feeds fetched successfully.');
    }

    private async notificationAllFeeds() {

        const bufferTime = 1 * 5 * 60 * 1000;
        let notificationStartTime: Date;
        const currentTime = new Date();
        notificationStartTime = new Date(currentTime.getTime() - (60 * 60 * 1000) - bufferTime);
        console.log('Starting to send notifications for all feeds...:', notificationStartTime);
        for (const category in categoryFeeds) {
            const categoryFeed = categoryFeeds[category];
            if (!categoryFeed) {
                console.log(`Category not found for notification: ${category}`);
                continue;
            }

            const feedUrl = categoryFeed.url;
            try {
                console.log(`Fetching RSS feed for notification from URL: ${feedUrl}`);
                const feed = await this.getRSSFeed(feedUrl, notificationStartTime);

                if (feed && feed.items.length > 0) {
                    await this._postMessageService.postFeedToSlack(feed, '#色々通知', category);
                    console.log(`Notification sent for category: ${category}`);
                } else {
                    console.log(`No recent items to notify for category: ${category}`);
                }
            } catch (error) {
                console.error(`Error while sending notification for category '${category}':`, error);
            }
        }
        console.log('All notifications sent successfully.');
    }

    private async getRSSFeed(feedUrl: string, notificationStartTime: Date): Promise<RSSFeed | null> {
        const feed = await this._rssFeedService.getRSSFeed(feedUrl);

        if (!feed || !feed.items) {
            return null;
        }

        const recentItems = feed.items.filter(item => {
            const pubDate = item.pubDate ? new Date(item.pubDate) : null;
            return pubDate && pubDate >= notificationStartTime;
        });

        if (recentItems.length === 0) {
            return null;
        }

        const feedWithRecentItems = { ...feed, items: recentItems };
        return feedWithRecentItems;
    }
}
export interface ICategoryFeed {
    url: string;
    lang: string;
    label: string;
}

export interface ICategoryFeeds {
    [category: string]: ICategoryFeed;
}

export const categoryFeeds: ICategoryFeeds = {
    azure: {
        url: 'https://azurecomcdn.azureedge.net/ja-jp/updates/feed/',
        lang: 'ja-jp',
        label: 'Azure News'
    },
    azure_blog: {
        url: 'https://azure.microsoft.com/ja-jp/blog/feed/',
        lang: 'ja-jp',
        label: 'Azure Blog News'
    },
    dotnet: {
        url: 'https://devblogs.microsoft.com/dotnet/feed/',
        lang: 'en-us',
        label: '.NET News'
    },
    typescript: {
        url: 'https://devblogs.microsoft.com/typescript/feed/',
        lang: 'en-us',
        label: 'TypeScript News'
    },
    nhn_techorus: {
        url: 'https://techblog.nhn-techorus.com/feed',
        lang: 'ja-jp',
        label: 'NHN Techorus'
    },
    gunosy: {
        url: 'https://tech.gunosy.io/rss',
        lang: 'ja-jp',
        label: 'Gunosy Tech'
    },
    sansan: {
        url: 'https://buildersbox.corp-sansan.com/rss',
        lang: 'ja-jp',
        label: 'Sansan Builders Box'
    },
    askul: {
        url: 'https://tech.askul.co.jp/feed',
        lang: 'ja-jp',
        label: 'Askul Tech Blog'
    },
    aws: {
        url: 'https://aws.amazon.com/jp/blogs/news/feed/',
        lang: 'ja-jp',
        label: 'AWS News'
    },
    google_blog: {
        url: 'https://developers-jp.googleblog.com/atom.xml',
        lang: 'ja-jp',
        label: 'Google Developers Blog'
    },
    sakura: {
        url: 'https://knowledge.sakura.ad.jp/rss/',
        lang: 'ja-jp',
        label: 'Sakura Internet Knowledge'
    },
    hatena: {
        url: 'https://developer.hatenastaff.com/rss',
        lang: 'ja-jp',
        label: 'Hatena Developer Blog'
    },
    feature: {
        url: 'https://future-architect.github.io/atom.xml',
        lang: 'ja-jp',
        label: 'Future Architect'
    },
    yahoo_line: {
        url: 'https://techblog.lycorp.co.jp/ja/feed/index.xml',
        lang: 'ja-jp',
        label: 'Yahoo & LINE Tech Blog'
    },
    openai: {
        url: 'https://jamesg.blog/openai.xml',
        lang: 'en-us',
        label: 'OpenAI Blog'
    },
    shibayan: {
        url: 'https://blog.shibayan.jp/rss',
        lang: 'ja-jp',
        label: 'しばやん雑記'
    },
    ufcpp: {
        url: 'https://ufcpp.net/rss',
        lang: 'ja-jp',
        label: '++C++; // 未確認飛行 C'
    },
    ipaSecurity: {
        url: 'https://www.ipa.go.jp/security/alert-rss.rdf',
        lang: 'ja-jp',
        label: '重要なセキュリティ情報  IPA'
    },
    jpcert: {
        url: 'https://www.jpcert.or.jp/rss/jpcert-all.rdf',
        lang: 'ja-jp',
        label: '重要なセキュリティ情報  JPCERT'
    }
};
