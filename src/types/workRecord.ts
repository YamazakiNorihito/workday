import { DateOnly } from "./dateOnly";
import { TimeOnly } from "./timeOnly";

export interface BreakRecord {
    clockInAt: TimeOnly;
    clockOutAt: TimeOnly;
}

export interface WorkRecord {
    workDay: DateOnly;
    breakRecords: BreakRecord[];
    clockInAt: TimeOnly;
    clockOutAt: TimeOnly;
}
