type Hour = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 | 21 | 22 | 23;
type Minute = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 | 21 | 22 | 23 | 24 | 25 | 26 | 27 | 28 | 29 | 30 | 31 | 32 | 33 | 34 | 35 | 36 | 37 | 38 | 39 | 40 | 41 | 42 | 43 | 44 | 45 | 46 | 47 | 48 | 49 | 50 | 51 | 52 | 53 | 54 | 55 | 56 | 57 | 58 | 59;
type Second = Minute;

export class TimeOnly {
    hour: Hour;
    minute: Minute;
    second: Second;

    constructor(hour: Hour, minute: Minute, second: Second) {
        this.hour = hour;
        this.minute = minute;
        this.second = second;
    }

    toString(): string {
        return `${this.hour.toString().padStart(2, '0')}:${this.minute.toString().padStart(2, '0')}`;
    }

    static fromTimeString(timeString: string): TimeOnly {
        let [datePart, timePart] = timeString.split('T');
        if (!timePart) {
            timePart = timeString;
        }
        const [hourString, minuteString, secondString] = timePart.split(':');
        const hour: Hour = +hourString as Hour;
        const minute: Minute = +minuteString as Minute;
        const timeOnly = new TimeOnly(hour, minute, 0);
        return timeOnly;
    }
}