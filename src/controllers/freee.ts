
import ejs from 'ejs';
import { Request, Response } from 'express';
import path from 'path';
import { inject, singleton } from 'tsyringe';
import { Employee, FreeeService, IFreeeService } from '../services/freeeHrService';
import { IWeekdayService } from '../services/weekdayService';
import { BreakRecord, WorkRecord } from '../types/workRecord';
import { TimeOnly } from '../types/timeOnly';
import { DateOnly } from '../types/dateOnly';
import { Semaphore } from '../system/semaphore';
import { AxiosError } from 'axios';
import polly from 'polly-js';

import { validationResult } from 'express-validator';
import { sendErrorResponse } from './apiResponse';

@singleton()
export class FreeeController {
    private readonly _callback_url: string;

    constructor(
        @inject(FreeeService) private readonly freeeService: IFreeeService
        , @inject("IWeekdayService") private readonly weekdayService: IWeekdayService
        , private _appDomainURL: string) {
        this._callback_url = `${this._appDomainURL}/freee/authorize/callback`;
    }

    public async getUserInfo(req: Request, res: Response): Promise<void> {
        const meInfo = await this.freeeService.getMe(req.session.oauthToken?.sub!);
        const currentDate = meInfo?.updated_at;
        const formattedDate = currentDate ? this.formatDate(currentDate) : undefined;

        const data: SyncEmployeeInfo = {
            userInfo: meInfo,
            lastSyncDate: formattedDate,
            yearMonth: null,
            workRecords: []
        };

        const indexPath = path.join(__dirname, './../views/freee/index.ejs')
        const renderedBody = await ejs.renderFile(indexPath, { data });
        res.render('layout', {
            title: 'Top',
            body: renderedBody
        });
    }

    public async redirectToAuth(req: Request, res: Response): Promise<void> {
        res.redirect(this.freeeService.getAuthorizeUrl(this._callback_url));
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
        await this.freeeService.handleAuthCallback(req.session.oauthToken?.sub!, authCode, this._callback_url);
        res.redirect('/freee');
    }

    public async getWorkRecords(req: Request, res: Response): Promise<void> {
        const yearMonth = new Date(req.query.yearMonth + '-01');
        yearMonth.setMonth(yearMonth.getMonth() + 1);
        const year = yearMonth.getFullYear();
        const month = yearMonth.getMonth() + 1;

        const meInfo = await this.freeeService.getMe(req.session.oauthToken?.sub!);
        const currentDate = meInfo?.updated_at;
        const formattedDate = currentDate ? this.formatDate(currentDate) : undefined;

        const workRecords = await this.freeeService.getWorkRecords(req.session.oauthToken?.sub!, year, month);

        const data: SyncEmployeeInfo = {
            userInfo: meInfo,
            lastSyncDate: formattedDate,
            yearMonth: req.query.yearMonth as string,
            workRecords: workRecords
        };

        const indexPath = path.join(__dirname, './../views/freee/index.ejs')
        const renderedBody = await ejs.renderFile(indexPath, { data });
        res.render('layout', {
            title: 'Top',
            body: renderedBody
        });
    }

    public async registerWorkRecords(req: Request, res: Response): Promise<void> {
        const errors = validationResult(req);

        if (!errors.isEmpty()) {
            return sendErrorResponse(res, errors);
        }

        const workFromDate = new Date(req.body.workFromDate);
        const workToDate = new Date(req.body.workToDate);

        const workStartHours: string = req.body.workStartHours;
        const workEndHours: string = req.body.workEndHours;
        const workBreakStartHours: string = req.body.workBreakStartHours;
        const workBreakEndHours: string = req.body.workBreakEndHours;

        const workDays = await this.weekdayService.getWorkDays(workFromDate, workToDate);
        const workRecords: WorkRecord[] = workDays.map((workDay: Date) => {
            const breakRecord: BreakRecord = {
                clockInAt: TimeOnly.fromTimeString(workBreakStartHours),
                clockOutAt: TimeOnly.fromTimeString(workBreakEndHours),
            };

            return {
                workDay: DateOnly.fromDate(workDay),
                breakRecords: [breakRecord],
                clockInAt: TimeOnly.fromTimeString(workStartHours),
                clockOutAt: TimeOnly.fromTimeString(workEndHours),
            };
        });

        const sem = new Semaphore(30, 1000); // 同時に実行できるタスク数を2に制限
        const tasks = workRecords.map(async record => {
            await sem.acquire();
            try {
                await this.freeeService.updateWorkRecord(req.session.oauthToken?.sub!, record);
            } finally {
                sem.release();
            }
        });

        await Promise.all(tasks);
        res.redirect(`/freee/work-records?yearMonth=${workDays[0].getFullYear()}-${(workDays[0].getMonth() + 1)}`);
    }

    public async deleteWorkRecords(req: Request, res: Response): Promise<void> {
        const workFromDate = new Date(req.body.workFromDate);
        const workToDate = new Date(req.body.workToDate);


        const workDays = (await this.weekdayService.getWorkDays(workFromDate, workToDate))
            .map((workDay: Date) => {
                return DateOnly.fromDate(workDay);
            });

        const sem = new Semaphore(10, 2000);
        const tasks = workDays.map(async workDay => {
            await sem.acquire();
            try {
                await polly()
                    .handle(err => {
                        if (err instanceof AxiosError && err.response && err.response.status === 500) {
                            return true;
                        }
                        return false;
                    })
                    .waitAndRetry(2)
                    .executeForPromise(async () => {
                        await this.freeeService.deleteWorkRecord(req.session.oauthToken?.sub!, workDay);
                    });
            } catch (e) {
                if (e instanceof AxiosError && e.response && (e.response.status === 404 || e.response.status === 500)) {
                    console.log(`Record not found for workDay: ${workDay}, skipping...`);
                } else {
                    console.log(`${workDay}: ${e}`);
                    throw e;
                }
            } finally {
                sem.release();
            }
        });

        await Promise.all(tasks);

        res.redirect(`/freee/work-records?yearMonth=${workFromDate.getFullYear()}-${(workFromDate.getMonth() + 1)}`);
    }

    private formatDate(date: Date): string {
        return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:${String(date.getSeconds()).padStart(2, '0')}`;
    }
}

type SyncEmployeeInfo = {
    userInfo: Employee | undefined;
    lastSyncDate: string | undefined;
    yearMonth: string | null;
    workRecords: WorkRecord[] | undefined;
}
