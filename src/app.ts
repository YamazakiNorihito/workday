import 'reflect-metadata';
import express from 'express';
import bodyParser from 'body-parser';
import dotenv from 'dotenv';
import path from 'path';
import { errorHandler } from './middlewares/errorHandler';
import methodOverride from 'method-override'

if (process.env.NODE_ENV === 'development') {
    dotenv.config({ path: '.env.development' });
} else {
    dotenv.config();
}
import appRoutes from './routes';
import container from './diContainer';
import session from 'express-session';
import { LoginController } from './controllers/login';
import { categoryFeeds } from './schedulers/newsScheduler';

declare module 'express-session' {
    export interface SessionData {
        oauthToken: {
            idToken: string;
            accessToken: string;
            refreshToken: string;
            expiresIn: number;
            tokenType: string;
            username: string;
            sub: string;
            expiresAt: Date;
        }
        returnTo?: string;
    }
}

const app = express();

// ビューエンジンとして 'ejs' を設定
app.set('view engine', 'ejs');
// ビューのディレクトリを設定
app.set('views', path.join(__dirname, 'views'));

app.use(express.json());
app.use(bodyParser.urlencoded({ extended: true }));
// override with the X-HTTP-Method-Override header in the request
// https://github.com/expressjs/method-override
app.use(methodOverride('_method'))
app.use(methodOverride('X-HTTP-Method')) //          Microsoft
app.use(methodOverride('X-HTTP-Method-Override')) // Google/GData
app.use(methodOverride('X-Method-Override')) //      IBM

app.use((req, res, next) => {
    res.locals.techNewsCategories = Object.entries(categoryFeeds).map(([key, value]) => {
        return {
            category: key,
            menuName: value.label
        };
    });
    next();
});

app.use(session({
    secret: 'session_secret',
    resave: false,
    saveUninitialized: false,
    cookie: { secure: false }
}));


app.use((req, res, next) => {
    const loginController = container.resolve(LoginController);
    return loginController.login(req, res, next);
});

app.use('/', appRoutes);

app.use(errorHandler);

const LISTEN_PORT = process.env.PORT || 3000;
app.listen(LISTEN_PORT, () => {
    console.log(`Server is running on port ${LISTEN_PORT} ${process.env.NODE_ENV}`);
});

export default app;