
import ejs from 'ejs';
import { Request, Response } from 'express';
import path from 'path';
import { inject, singleton } from 'tsyringe';

import { HackerNewsService } from '../services/hackerNewsService';
import { RSSFeedService } from '../services/rssFeedService';
import { RSSFeed, RSSFeedItem } from '../repositories/rssFeedRepository';
import { categoryFeeds } from '../schedulers/newsScheduler';

@singleton()
export class NewsController {
    constructor(
        @inject(HackerNewsService) private readonly hackerNewsService: HackerNewsService,
        @inject(RSSFeedService) private readonly rssFeedService: RSSFeedService) { }

    public async getHackerNewsTopStories(req: Request, res: Response): Promise<void> {

        const stories = await this.hackerNewsService.getTopStories();
        const indexPath = path.join(__dirname, './../views/news/index.ejs')
        const pageTitle = "HackerNews TopStories";
        const renderedBody = await ejs.renderFile(indexPath, { stories, pageTitle });
        res.render('layout', {
            title: pageTitle,
            body: renderedBody
        });
    }

    public async getHackerNewsNewStories(req: Request, res: Response): Promise<void> {

        const stories = await this.hackerNewsService.getNewStories();
        const indexPath = path.join(__dirname, './../views/news/index.ejs')
        const pageTitle = "HackerNews NewStories";
        const renderedBody = await ejs.renderFile(indexPath, { stories, pageTitle });
        res.render('layout', {
            title: pageTitle,
            body: renderedBody
        });
    }

    public async getHackerNewsBestStories(req: Request, res: Response): Promise<void> {

        const stories = await this.hackerNewsService.getBestStories();
        const indexPath = path.join(__dirname, './../views/news/index.ejs');
        const pageTitle = "HackerNews BestStories";
        const renderedBody = await ejs.renderFile(indexPath, { stories, pageTitle });
        res.render('layout', {
            title: pageTitle,
            body: renderedBody
        });
    }

    private categoryFeeds = categoryFeeds;

    public async getRSSFeed(req: Request, res: Response): Promise<void> {
        const category = req.params.category;
        const categoryFeed = this.categoryFeeds[category];
        if (!categoryFeed) {
            res.status(404).send('Category not found');
            return;
        }

        const feedUrl = categoryFeed.url;
        const feed = await this.rssFeedService.getRSSFeed(feedUrl);

        const indexPath = path.join(__dirname, './../views/news/rss.ejs');
        const renderedBody = await ejs.renderFile(indexPath, { feed });
        res.render('layout', {
            title: feed?.title,
            body: renderedBody
        });
    }
    public async getAllRSSFeed(req: Request, res: Response): Promise<void> {
        let allFeeds: RSSFeedItem[] = [];
        let latestDates: Date[] = [new Date(0)];

        const feedPromises = Object.entries(this.categoryFeeds).map(async ([category, feedInfo]) => {
            try {
                const feed = await this.rssFeedService.getRSSFeed(feedInfo.url);
                if (!feed) {
                    return [];
                }
                if (feed?.lastBuildDate) {
                    latestDates.push(new Date(feed.lastBuildDate));
                }
                return feed.items.map(item => ({ ...item }));
            } catch (error) {
                console.error(`Error fetching RSS feed for category ${category}:`, error);
                return [];
            }
        });

        const results = await Promise.all(feedPromises);
        allFeeds = results.flat();

        allFeeds.sort((a, b) => {
            const dateA = a.pubDate ? new Date(a.pubDate).getTime() : 0;
            const dateB = b.pubDate ? new Date(b.pubDate).getTime() : 0;
            return dateB - dateA;
        });

        // すべてのフィードの中で最新の日付を取得
        let latestDate = new Date(Math.max(...latestDates.map(date => date.getTime())));

        const feed: RSSFeed = {
            title: "ALL Category RSS Feeds",
            lastBuildDate: latestDate,
            items: allFeeds
        };

        const indexPath = path.join(__dirname, './../views/news/rss.ejs');
        const renderedBody = await ejs.renderFile(indexPath, { feed: feed });
        res.render('layout', {
            title: "ALL Category RSS Feeds",
            body: renderedBody
        });
    }
}
