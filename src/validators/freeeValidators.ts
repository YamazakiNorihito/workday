import { body, validationResult } from 'express-validator';

export const validateWorkRecords = [
    // 日付のバリデーション
    body('workFromDate').isDate().withMessage('workFromDate must be a valid date.'),
    body('workToDate').isDate().withMessage('workToDate must be a valid date.'),

    // 時間のバリデーション（HH:MM形式を想定）
    body('workStartHours').matches(/^([01]?[0-9]|2[0-3]):[0-5][0-9]$/).withMessage('workStartHours must be in HH:MM format.'),
    body('workEndHours').matches(/^([01]?[0-9]|2[0-3]):[0-5][0-9]$/).withMessage('workEndHours must be in HH:MM format.'),
    body('workBreakStartHours').matches(/^([01]?[0-9]|2[0-3]):[0-5][0-9]$/).withMessage('workBreakStartHours must be in HH:MM format.'),
    body('workBreakEndHours').matches(/^([01]?[0-9]|2[0-3]):[0-5][0-9]$/).withMessage('workBreakEndHours must be in HH:MM format.'),
];