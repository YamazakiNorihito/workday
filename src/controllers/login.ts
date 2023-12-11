
import { NextFunction, Request, Response } from 'express';
import { inject, singleton } from 'tsyringe';
import { ILoginAuthenticationService } from '../services/oauth2Service';
import jwt, { JwtPayload } from 'jsonwebtoken';

@singleton()
export class LoginController {
    private readonly _callback_url: string;

    constructor(
        @inject("LoginAuthenticationService") private readonly loginAuthenticationService: ILoginAuthenticationService
        , private _appDomainURL: string) {
        this._callback_url = `${this._appDomainURL}/authorize/callback`;
    }

    public async login(req: Request, res: Response, next: NextFunction) {
        try {
            if (req.path === '/authorize/callback') {
                return next();
            }
            if (req.session.oauthToken) {
                const currentTime = new Date();
                if (currentTime > req.session.oauthToken.expiresAt) {
                    await this.refreshToken(req, res);
                }
                return next();
            }
            const loginUrl = this.loginAuthenticationService.getAuthorizationUrl(this._callback_url);
            console.log("loginUrl:", loginUrl)
            req.session.returnTo = req.originalUrl;
            res.redirect(loginUrl);
        } catch (e) {
            console.log(e);
            delete req.session.oauthToken;
            res.status(500).json({
                error: {
                    message: e,
                }
            });
        }
    }

    public async authCallback(req: Request, res: Response): Promise<void> {
        const authCode = req.query.code as string;

        if (!authCode) {
            res.status(400).json({
                error: {
                    message: 'Error during authentication.',
                }
            });
        }
        const token = await this.loginAuthenticationService.getAccessToken(authCode, this._callback_url);

        const idTokenPayload = await this.verifyToken(token.id_token);
        const accessTokenPayload = await this.verifyToken(token.access_token);

        req.session.oauthToken = {
            idToken: token.id_token,
            accessToken: token.access_token,
            refreshToken: token.refresh_token,
            expiresIn: token.expires_in,
            tokenType: token.token_type,
            username: accessTokenPayload["username"],
            sub: idTokenPayload.sub!,
            expiresAt: new Date(idTokenPayload.exp! * 1000)
        }
        const redirectUrl = req.session.returnTo || '/';
        delete req.session.returnTo;
        res.redirect(redirectUrl);
    }

    private async refreshToken(req: Request, res: Response): Promise<void> {
        const refreshedToken = await this.loginAuthenticationService.refreshToken(req.session.oauthToken!.refreshToken!);
        const newIdTokenPayload = await this.verifyToken(refreshedToken.id_token);
        const newAccessTokenPayload = await this.verifyToken(refreshedToken.access_token);

        req.session.oauthToken = {
            idToken: refreshedToken.id_token,
            accessToken: refreshedToken.access_token,
            refreshToken: refreshedToken.refresh_token,
            expiresIn: refreshedToken.expires_in,
            tokenType: refreshedToken.token_type,
            username: newAccessTokenPayload["username"],
            sub: newIdTokenPayload.sub!,
            expiresAt: new Date(newIdTokenPayload.exp! * 1000)
        }
    }

    private async verifyToken(token: string): Promise<JwtPayload> {
        const decodedToken = jwt.decode(token, { complete: true });
        if (!decodedToken || typeof decodedToken !== 'object' || !decodedToken.header.kid) {
            throw new Error('Invalid token');
        }
        const jwk = await this.loginAuthenticationService.getPublicKey(decodedToken.header.kid);
        if (!jwk) {
            throw new Error('Public key not found');
        }
        return jwt.verify(token, jwk) as JwtPayload;
    }
}
/*
    idToken: {
        at_hash: '99zzI6AZ7Gjbk50JDMp64A',
        sub: '45982b36-3580-457b-80d4-4fa79851fed2',
        email_verified: true,
        iss: 'https://cognito-idp.ap-northeast-1.amazonaws.com/ap-northeast-1_69ZwVZJZy',
        'cognito:username': 'yamazakinorihito',
        origin_jti: '94e53f0a-9b47-4e0d-a042-121874a5a8e8',
        aud: '1bb8ik2kufie8fnlat1ik11dgn',
        token_use: 'id',
        auth_time: 1701216154,
        exp: 1701219754,
        iat: 1701216155,
        jti: '151df7d0-5328-4f04-823d-55564d08ecb1',
        email: 'n.yamazaki@frontierfield.co.jp'
    }
    accessToken: {
        sub: '45982b36-3580-457b-80d4-4fa79851fed2',
        iss: 'https://cognito-idp.ap-northeast-1.amazonaws.com/ap-northeast-1_69ZwVZJZy',
        version: 2,
        client_id: '1bb8ik2kufie8fnlat1ik11dgn',
        origin_jti: '94e53f0a-9b47-4e0d-a042-121874a5a8e8',
        token_use: 'access',
        scope: 'phone openid email',
        auth_time: 1701216154,
        exp: 1701219754,
        iat: 1701216155,
        jti: 'ed8bdc07-2542-41bd-8c33-8b9960e12688',
        username: 'yamazakinorihito'
    }
*/