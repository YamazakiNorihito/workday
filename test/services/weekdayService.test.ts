
import 'reflect-metadata';
import axios from 'axios';
import { IWeekdayService, WeekdayService } from './../../src/services/weekdayService'

const mockedAxios = axios as jest.Mocked<typeof axios>;
jest.mock('axios');

describe('WeekdayService', () => {
    let weekdayService: IWeekdayService;

    beforeEach(() => {
        mockedAxios.create.mockReturnThis();
        mockedAxios.get.mockReset();
        weekdayService = new WeekdayService();
    });

    describe('getWorkDays', () => {
        it('should return workdays only when there are no holidays in the period', async () => {
            // Arrange
            mockedAxios.get.mockResolvedValue({ data: {} });
            const fromDate = new Date('2023-05-08');
            const toDate = new Date('2023-05-12');

            // Act
            const actual = await weekdayService.getWorkDays(fromDate, toDate);

            // Assert
            expect(actual).toEqual([
                new Date('2023-05-08'), // 月曜日
                new Date('2023-05-09'), // 火曜日
                new Date('2023-05-10'), // 水曜日
                new Date('2023-05-11'), // 木曜日
                new Date('2023-05-12'), // 金曜日
            ]);
        });
        it('should return workdays excluding weekends', async () => {
            // Arrange
            mockedAxios.get.mockResolvedValue({ data: {} });
            const fromDate = new Date('2023-05-08');
            const toDate = new Date('2023-05-14');

            // Act
            const actual = await weekdayService.getWorkDays(fromDate, toDate);

            // Assert
            expect(actual).toEqual([
                new Date('2023-05-08'), // 月曜日
                new Date('2023-05-09'), // 火曜日
                new Date('2023-05-10'), // 水曜日
                new Date('2023-05-11'), // 木曜日
                new Date('2023-05-12'), // 金曜日
            ]);
        });
        it('should return workdays excluding public holidays when the period includes holidays', async () => {
            // Arrange
            mockedAxios.get.mockResolvedValue({
                data: {
                    "2023-05-03": "憲法記念日",
                    "2023-05-04": "みどりの日",
                    "2023-05-05": "こどもの日",
                }
            });
            const fromDate = new Date('2023-05-01');
            const toDate = new Date('2023-05-05');

            // Act
            const actual = await weekdayService.getWorkDays(fromDate, toDate);

            // Assert
            expect(actual).toEqual([
                new Date('2023-05-01'), // 月曜日
                new Date('2023-05-02'), // 火曜日
            ]);
        });
        it('should return workdays excluding both weekends and holidays when the period includes both', async () => {
            // Arrange
            mockedAxios.get.mockResolvedValue({
                data: {
                    "2023-05-03": "憲法記念日",
                    "2023-05-04": "みどりの日",
                    "2023-05-05": "こどもの日",
                }
            });
            const fromDate = new Date('2023-05-01');
            const toDate = new Date('2023-05-14');

            // Act
            const actual = await weekdayService.getWorkDays(fromDate, toDate);

            // Assert
            expect(actual).toEqual([
                new Date('2023-05-01'), // 月曜日
                new Date('2023-05-02'), // 火曜日
                new Date('2023-05-08'), // 月曜日
                new Date('2023-05-09'), // 火曜日
                new Date('2023-05-10'), // 水曜日
                new Date('2023-05-11'), // 木曜日
                new Date('2023-05-12'), // 金曜日
            ]);
        });
        it('should return an empty list of workdays when the period consists only of a weekend', async () => {
            // Arrange
            mockedAxios.get.mockResolvedValue({ data: {} });
            const fromDate = new Date('2023-05-13');
            const toDate = new Date('2023-05-14');

            // Act
            const actual = await weekdayService.getWorkDays(fromDate, toDate);

            // Assert
            expect(actual).toHaveLength(0);
        });
        it('should return an empty list of workdays when fromDate and toDate are reversed', async () => {
            mockedAxios.get.mockResolvedValue({ data: {} });
            const fromDate = new Date('2023-05-01');
            const toDate = new Date('2022-04-28');

            // Act
            const actual = await weekdayService.getWorkDays(fromDate, toDate);

            // Assert
            expect(actual).toHaveLength(0);
        });
        it('should propagate the exception to the caller in case of an error', async () => {
            // Arrange
            mockedAxios.get.mockRejectedValue(new Error('Network Error'));
            const fromDate = new Date('2023-05-13');
            const toDate = new Date('2023-05-14');

            // Act&Assert
            await expect(weekdayService.getWorkDays(fromDate, toDate))
                .rejects
                .toThrow('Network Error');
            expect(mockedAxios.get).toHaveBeenCalledTimes(1);
        });
    });
});