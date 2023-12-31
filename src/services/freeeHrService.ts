import { inject, singleton } from "tsyringe";
import { FreeeHrHttpApiClient, IFreeeHrHttpApiClient } from "../httpClients/freeeHttpClient";
import { FreeeUserRepository, IFreeeUserRepository, OAuth, User as UserModel } from '../repositories/freeeUserRepository';
import { WorkRecord } from '../types/workRecord';
import { RedisClientType } from "redis";
import { DateOnly } from "../types/dateOnly";
import { TimeOnly } from "../types/timeOnly";
import { IOAuth2Service, OAuthTokenResponse } from "./oauth2Service";
import axios, { AxiosInstance } from "axios";

export interface FreeeOAuthTokenResponse extends OAuthTokenResponse {
    scope: string;
    created_at: number;
    company_id: number;
}

export interface IFreeeAuthenticationService extends IOAuth2Service<FreeeOAuthTokenResponse> { }

export interface IFreeeService {
    getAuthorizeUrl(redirectUri: string): string;
    handleAuthCallback(userId: string, authCode: string, redirectUri: string): Promise<Employee>;
    getWorkRecords(userId: string, year: number, month: number): Promise<WorkRecord[]>;
    getMe(userId: string): Promise<Employee>;
    updateWorkRecord(userId: string, workRecord: WorkRecord): Promise<void>;
    deleteWorkRecord(userId: string, workDay: DateOnly): Promise<void>;
}

@singleton()
export class FreeeAuthenticationService implements IFreeeAuthenticationService {
    private httpClient: AxiosInstance;
    private readonly freeeDomain = 'https://accounts.secure.freee.co.jp';

    constructor(private clientId: string, private clientSecret: string) {
        this.httpClient = axios.create({ baseURL: this.freeeDomain });
    }

    public getAuthorizationUrl(callbackUrl: string): string {
        return `${this.freeeDomain}/public_api/authorize?client_id=${this.clientId}&redirect_uri=${encodeURIComponent(callbackUrl)}&response_type=code&prompt=select_company`;
    }

    public async getAccessToken(authorizationCode: string, redirectUri: string): Promise<FreeeOAuthTokenResponse> {
        const response = await this.httpClient.post(
            `/public_api/token`,
            new URLSearchParams({
                grant_type: 'authorization_code',
                client_id: this.clientId,
                client_secret: this.clientSecret,
                code: authorizationCode,
                redirect_uri: redirectUri,
            }),
            {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            }
        );
        return response.data;
    }

    public async getPublicKey(kid: string): Promise<string | null> {
        throw new Error("公開鍵が公開されいません。");
    }

    public async refreshToken(refreshToken: string): Promise<FreeeOAuthTokenResponse> {
        const response = await this.httpClient.post(
            `/oauth2/token`,
            new URLSearchParams({
                grant_type: 'refresh_token',
                client_id: this.clientId,
                client_secret: this.clientSecret,
                refresh_token: refreshToken,
            }),
            {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            }
        );
        return response.data;
    }
}

@singleton()
export class FreeeService implements IFreeeService {
    constructor(
        @inject("IFreeeAuthenticationService") private freeeAuthenticationService: IFreeeAuthenticationService,
        @inject(FreeeHrHttpApiClient) private freeeHrHttpApiClient: IFreeeHrHttpApiClient,
        @inject(FreeeUserRepository) private readonly freeeUserRepository: IFreeeUserRepository,
        @inject("RedisClient") private readonly redisClient: RedisClientType
    ) { }

    public getAuthorizeUrl(redirect_uri: string): string {
        return this.freeeAuthenticationService.getAuthorizationUrl(redirect_uri);
    }

    public async handleAuthCallback(userId: string, authCode: string, redirect_uri: string): Promise<Employee> {
        const oauthResponse = await this.freeeAuthenticationService.getAccessToken(authCode, redirect_uri);
        const meResponse = await this.freeeHrHttpApiClient.get<User>('/api/v1/users/me', oauthResponse.access_token);

        const freeeUser: UserModel = {
            id: meResponse?.id,
            companies: meResponse?.companies,
            oauth: oauthResponse,
            updated_at: Date.now(),
        };
        await this.freeeUserRepository.save(userId, freeeUser);
        return this.getInternalMe(freeeUser);
    }

    public async getWorkRecords(userId: string, year: number, month: number): Promise<WorkRecord[]> {
        const me = await this.getMe(userId);
        const accessToken = await this.getAccessToken(userId);
        const response = await this.freeeHrHttpApiClient.get<FreeeWorkSummary>(
            `/api/v1/employees/${me.employee_id}/work_record_summaries/${year}/${month}?company_id=${me.company_id}&work_records=true`
            , accessToken);

        const mappedWorkRecords: WorkRecord[] = response.work_records
            .filter(workRecord => workRecord.clock_in_at !== null && workRecord.clock_out_at !== null)
            .map((workRecord) => ({
                workDay: DateOnly.fromDateString(workRecord.date),
                breakRecords: workRecord.break_records.map((breakRecord) => ({
                    clockInAt: TimeOnly.fromTimeString(breakRecord.clock_in_at),
                    clockOutAt: TimeOnly.fromTimeString(breakRecord.clock_out_at)
                })),
                clockInAt: TimeOnly.fromTimeString(workRecord.clock_in_at!),
                clockOutAt: TimeOnly.fromTimeString(workRecord.clock_out_at!),
            }));
        return mappedWorkRecords;
    }

    public async getMe(userId: string): Promise<Employee> {
        let freeeUser = await this.freeeUserRepository.get(userId);
        return this.getInternalMe(freeeUser!);
    }

    public async updateWorkRecord(userId: string, workRecord: WorkRecord) {
        const me = await this.getMe(userId);
        const updateWorkRecord: FreeeWorkRecord =
        {
            company_id: me.company_id,
            break_records: workRecord.breakRecords.map((p) => ({
                clock_in_at: `${workRecord.workDay.toString()} ${p.clockInAt.toString()}`,
                clock_out_at: `${workRecord.workDay.toString()} ${p.clockOutAt.toString()}`
            })),
            clock_in_at: `${workRecord.workDay.toString()} ${workRecord.clockInAt.toString()}`,
            clock_out_at: `${workRecord.workDay.toString()} ${workRecord.clockOutAt.toString()}`
        }

        const accessToken = await this.getAccessToken(userId);
        await this.freeeHrHttpApiClient.put(
            `/api/v1/employees/${me.employee_id}/work_records/${workRecord.workDay.toString()}`
            , accessToken
            , updateWorkRecord);
    }

    public async deleteWorkRecord(userId: string, workDay: DateOnly) {
        const me = await this.getMe(userId);
        const accessToken = await this.getAccessToken(userId);
        await this.freeeHrHttpApiClient.delete(
            `/api/v1/employees/${me.employee_id}/work_records/${workDay.toString()}?company_id=${me.company_id}`
            , accessToken);
    }

    private getInternalMe(user: UserModel): Employee {
        if (!user || !user.companies) {
            return {
                employee_id: 0,
                employee_name: "",
                company_id: 0,
                company_name: "",
                external_cid: 0,
                role: "",
                updated_at: null
            };
        }
        const targetCompany = user
            .companies
            .find(company => company.name === '株式会社フロンティア・フィールド');

        if (!targetCompany) {
            return {
                employee_id: 0,
                employee_name: "",
                company_id: 0,
                company_name: "",
                external_cid: 0,
                role: "",
                updated_at: null
            };
        }

        return {
            employee_id: targetCompany.employee_id ? targetCompany.employee_id : 0,
            employee_name: targetCompany.display_name ? targetCompany.display_name : "",
            company_id: targetCompany.id,
            company_name: targetCompany.name,
            external_cid: targetCompany.external_cid,
            role: targetCompany.role as string,
            updated_at: new Date(user.updated_at)
        };
    }

    protected async getAccessToken(userId: string, retryCount: number = 0): Promise<string> {
        const MAX_RETRIES = 3;

        const freeeUser = await this.freeeUserRepository.get(userId);

        if (!freeeUser) {
            throw new Error(`No User with ID: ${userId}`);
        }

        if (this.isTokenExpired(freeeUser.oauth)) {
            const lockKey = `lock:tokenrefresh:${userId}`;

            const lock = await this.redisClient.set(lockKey, "1", { NX: true, EX: 10 });
            if (!lock) {
                if (retryCount < MAX_RETRIES) {
                    // wait 1sec
                    await new Promise(res => setTimeout(res, 3000));
                    return this.getAccessToken(userId, retryCount + 1);
                } else {
                    throw new Error("Failed to get access token after multiple retries");
                }
            }

            try {
                if (this.isTokenExpired(freeeUser.oauth)) {
                    const oauthResponse = await this.freeeAuthenticationService.refreshToken(freeeUser.oauth.refresh_token);
                    freeeUser.oauth = oauthResponse;
                    freeeUser.updated_at = Date.now();
                    await this.freeeUserRepository.save(userId, freeeUser);
                }
            } finally {
                await this.redisClient.del(lockKey);
            }
        }

        return freeeUser.oauth.access_token;
    }


    private isTokenExpired(userOAuth: FreeeOAuthTokenResponse): boolean {
        const BUFFER_TIME = 30; // 30 seconds buffer
        const now = Math.floor(Date.now() / 1000); // 現在の時間（秒単位）
        const tokenExpiryTime = userOAuth.created_at + userOAuth.expires_in - BUFFER_TIME;
        return now > tokenExpiryTime;
    }
}

enum UserRole {
    COMPANY_ADMIN = 'company_admin',
    SELF_ONLY = 'self_only',
    CLERK = 'clerk'
}

interface Company {
    id: number;             // 事業所ID
    name: string;           // 事業所名
    role: UserRole;         // 事業所におけるロール
    external_cid: number;   // 事業所番号(半角数字10桁)
    employee_id?: number | null;   // 事業所に所属する従業員としての従業員ID、従業員情報が未登録の場合はnullになります。
    display_name?: string | null;  // 事業所に所属する従業員の表示名
}

interface User {
    id: number;            // ユーザーID
    companies: Company[];  // ユーザーが属する事業所の一覧
}

interface FreeeBreakRecord {
    clock_in_at: string;
    clock_out_at: string;
}

interface FreeeWorkRecord {
    company_id: number;
    break_records: FreeeBreakRecord[];
    clock_in_at: string;
    clock_out_at: string;
}
interface FreeeWorkSummary {
    year: number;
    month: number;
    start_date: string;
    end_date: string;
    work_days: number;
    total_work_mins: number;
    total_normal_work_mins: number;
    total_excess_statutory_work_mins: number;
    total_overtime_except_normal_work_mins: number;
    total_overtime_within_normal_work_mins: number;
    total_holiday_work_mins: number;
    total_latenight_work_mins: number;
    num_absences: number;
    num_paid_holidays: number;
    num_paid_holidays_and_hours: {
        days: number;
        hours: number;
    };
    num_paid_holidays_left: number;
    num_paid_holidays_and_hours_left: {
        days: number;
        hours: number;
    };
    num_substitute_holidays_used: number;
    num_compensatory_holidays_used: number;
    num_special_holidays_used: number;
    num_special_holidays_and_hours_used: {
        days: number;
        hours: number;
    };
    total_lateness_and_early_leaving_mins: number;
    multi_hourly_wages?: {
        name: string;
        total_normal_time_mins: number;
    }[];
    work_records: {
        break_records: {
            clock_in_at: string;
            clock_out_at: string;
        }[];
        clock_in_at?: string;
        clock_out_at?: string;
        date: string;
        day_pattern: 'normal_day' | 'prescribed_holiday' | 'legal_holiday';
        schedule_pattern?: '' | 'substitute_holiday_work' | 'substitute_holiday' | 'compensatory_holiday_work' | 'compensatory_holiday' | 'special_holiday';
        early_leaving_mins: number;
        hourly_paid_holiday_mins: number;
        is_absence: boolean;
        is_editable: boolean;
        lateness_mins: number;
        normal_work_clock_in_at?: string;
        normal_work_clock_out_at?: string;
        normal_work_mins: number;
        normal_work_mins_by_paid_holiday: number;
        note: string;
        paid_holiday: number;
        use_attendance_deduction: boolean;
        use_default_work_pattern: boolean;
        total_overtime_work_mins: number;
        total_holiday_work_mins: number;
        total_latenight_work_mins: number;
        not_auto_calc_work_time: boolean;
        total_excess_statutory_work_mins: number;
        total_latenight_excess_statutory_work_mins: number;
        total_overtime_except_normal_work_mins: number;
        total_latenight_overtime_except_normal_work_min: number;
    }[];
    total_shortage_work_mins?: number;
    total_deemed_paid_excess_statutory_work_mins?: number;
    total_deemed_paid_overtime_except_normal_work_mins?: number;
}

export interface Employee {
    employee_id: number;
    employee_name: string;
    company_id: number;
    company_name: string;
    external_cid: number;
    role: string;
    updated_at: Date | null;
}