import container from './diContainer';
import express from 'express';
import { asyncHandler } from './middlewares/asyncHandler';
import { TopController } from './controllers/top';
import { FreeeController } from './controllers/freee';
import { validateWorkRecords } from './validators/freeeValidators'
import { NewsController } from './controllers/news';
import { NewsScheduler } from './schedulers/newsScheduler';
import { LoginController } from './controllers/login';

// routing
const router = express.Router();

const topController = container.resolve(TopController);
const freeeController = container.resolve(FreeeController);
const newsController = container.resolve(NewsController);
const oauthController = container.resolve(LoginController);

router.get('/', asyncHandler((req, res) => topController.get(req, res)));
router.get('/authorize/callback', asyncHandler((req, res) => oauthController.authCallback(req, res)));

router.get('/freee', asyncHandler((req, res) => freeeController.getUserInfo(req, res)));
router.get('/freee/authorize', asyncHandler((req, res) => freeeController.redirectToAuth(req, res)));
router.get('/freee/authorize/callback', asyncHandler((req, res) => freeeController.authCallback(req, res)));
router.get('/freee/work-records', asyncHandler((req, res) => freeeController.getWorkRecords(req, res)));
router.post('/freee/work-records', validateWorkRecords, asyncHandler((req, res) => freeeController.registerWorkRecords(req, res)));
router.delete('/freee/work-records', asyncHandler((req, res) => freeeController.deleteWorkRecords(req, res)));

router.get('/news/hacker-news/top', asyncHandler((req, res) => newsController.getHackerNewsTopStories(req, res)));
router.get('/news/hacker-news/new', asyncHandler((req, res) => newsController.getHackerNewsNewStories(req, res)));
router.get('/news/hacker-news/best', asyncHandler((req, res) => newsController.getHackerNewsBestStories(req, res)));

router.get('/news/:category', asyncHandler((req, res) => newsController.getRSSFeed(req, res)));
router.get('/news', asyncHandler((req, res) => newsController.getAllRSSFeed(req, res)));

const newsScheduler = container.resolve(NewsScheduler);
newsScheduler.start();


export default router;
