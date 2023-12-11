export class DateOnly {
    public readonly year: number;
    public readonly month: number;
    public readonly day: number;

    constructor(year: number, month: number, day: number) {
        if (!DateOnly.isValidDate(year, month, day)) {
            throw new Error('Invalid date provided.');
        }
        this.year = year;
        this.month = month;
        this.day = day;
    }

    toString(): string {
        return `${this.year}-${this.month.toString().padStart(2, '0')}-${this.day.toString().padStart(2, '0')}`;
    }

    toFullDateString(): string {
        const dayNames = ["日", "月", "火", "水", "木", "金", "土"];
        const jsDate = new Date(this.toString());
        const dayName = dayNames[jsDate.getDay()];
        return `${this.toString()} (${dayName})`;
    }

    static fromDateString(dateString: string): DateOnly {
        const date = this.tryDateParse(dateString);

        if (!date) {
            throw new Error('Invalid date string.');
        }

        const [yearString, monthString, dayString] = dateString.split('-');
        return new DateOnly(Number(yearString), Number(monthString), Number(dayString));
    }

    static fromDate(date: Date): DateOnly {
        return new DateOnly(date.getFullYear(), date.getMonth() + 1, date.getDate());
    }

    private static tryDateParse(dateString: string): Date | null {
        const date = new Date(dateString);
        return date.toString() !== "Invalid Date" && !isNaN(date.getTime()) ? date : null;
    }

    private static isValidDate(year: number, month: number, day: number): boolean {
        const dateString = `${year}-${month.toString().padStart(2, '0')}-${day.toString().padStart(2, '0')}`;
        const date = new Date(dateString);
        return date.toString() !== "Invalid Date" && !isNaN(date.getTime());
    }
}
