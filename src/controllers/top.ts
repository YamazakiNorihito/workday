
import ejs from 'ejs';
import { Request, Response } from 'express';
import path from 'path';
import { singleton } from 'tsyringe';

@singleton()
export class TopController {
    constructor() { }

    public async get(req: Request, res: Response): Promise<void> {
        const indexPath = path.join(__dirname, './../views/top/index.ejs')
        const renderedBody = await ejs.renderFile(indexPath, { contents: `Hello! ${req.session.oauthToken?.username}  ${req.session.oauthToken?.sub} ` });
        res.render('layout', {
            title: 'Top',
            body: renderedBody
        });
    }
}