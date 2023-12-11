import { Response } from 'express';

/*
 考えるのがだるかったのでControllersにおいている。
*/

interface SuccessResponse<T> {
    data: T;
    message: string;
    status: number;
}

interface ErrorObject {
    field: string;
    message: string;
}

interface ErrorResponse {
    errors: ErrorObject[];
    message: string;
    status: number;
}

export const sendSuccessResponse = <T>(res: Response, data: T, message: string = 'Success', statusCode: number = 200): void => {
    const apiResponse: SuccessResponse<T> = {
        data,
        message,
        status: statusCode,
    };
    res.status(statusCode).json(apiResponse);
};

export const sendErrorResponse = (
    res: Response,
    errors: any,
    message: string = 'Error',
    statusCode: number = 400
): void => {
    const errorResponse: ErrorResponse = {
        errors: errors.array().map((error: any) => ({
            field: error.path,
            message: error.msg,
        })),
        message,
        status: statusCode,
    };
    res.status(statusCode).json(errorResponse);
};
