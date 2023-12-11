import axios, { AxiosInstance } from "axios";
import { singleton } from "tsyringe";

@singleton()
export class WeekdayService {
    private holidaysClient: AxiosInstance;

    constructor() {
        this.holidaysClient = axios.create({
            baseURL: 'https://holidays-jp.github.io'
        });
    }

    public async getWorkDays(from: Date, to: Date): Promise<Date[]> {
        const weekdays: Date[] = [];
        let currentDate = this.toJST(from);
        const targetEnd = this.toJST(to);
        while (currentDate <= targetEnd) {
            // 週末 (0: Sunday, 6: Saturday) を除外
            if (currentDate.getUTCDay() !== 0 && currentDate.getUTCDay() !== 6) {
                weekdays.push(new Date(currentDate));
            }
            currentDate.setUTCDate(currentDate.getUTCDate() + 1);
        }

        const holidays = await this.getHolidays();

        // 祝日を除外
        return weekdays.filter(day => !holidays.some(holiday => this.equals(day, holiday)));
    }

    private toJST(date: Date): Date {
        const offset = 9 * 60; // JST is UTC+9
        const newDate = new Date(date.getTime() + offset * 60 * 1000);
        newDate.setUTCHours(0, 0, 0, 0);
        return newDate;
    }

    private async getHolidays(): Promise<Date[]> {
        const response = await this.holidaysClient.get('/api/v1/date.json');
        return Object.keys(response.data).map(key => new Date(key));
    }

    private equals(d1: Date, d2: Date): boolean {
        return d1.getFullYear() === d2.getFullYear() &&
            d1.getMonth() === d2.getMonth() &&
            d1.getDate() === d2.getDate();
    }
}

type HolidaysResponse = {
    [date: string]: string;
};
