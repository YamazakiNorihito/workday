/**
 * Semaphoreクラスは並行性を制限するメカニズムを提供します。
 * 同時に実行できるタスクの数を制御するために使用します。
 * オプションとして、セマフォの取得を遅延させるための待ち時間も設定できます。
 */
export class Semaphore {
    private waitTimeMs: number;

    /**
     * Semaphoreのインスタンスを作成します。
     * @param {number} max - 同時に許可されるタスクの最大数。
     * @param {number} [waitTimeMs=0] - セマフォを取得する前に待機する時間（ミリ秒）。
     */
    constructor(private max: number, waitTimeMs: number = 0) {
        this.queue = [];
        this.count = 0;
        this.waitTimeMs = waitTimeMs;
    }

    private queue: (() => void)[];
    private count: number;

    /**
     * セマフォを取得しようとします。同時のタスクの最大数に達した場合、
     * このメソッドはセマフォを取得できるまで待機します。
     * 待ち時間が設定されている場合、セマフォを取得する前に指定された時間だけ待機します。
     * タスクが待機列に追加され、その後実行されると、指定された待ち時間が経過した後にPromiseが解決されます。
     * @returns {Promise<void>}
     */
    async acquire(): Promise<void> {
        if (this.count < this.max) {
            this.count++;
            return Promise.resolve();
        }

        return new Promise<void>(resolve => {
            this.queue.push(async () => {
                // 待ち時間が設定されている場合、指定された時間だけ待機します。
                if (this.waitTimeMs > 0) {
                    await new Promise(res => setTimeout(res, this.waitTimeMs));
                }
                resolve();
            });
        });
    }

    /**
     * セマフォをリリースします。これにより、セマフォの待機中の他のタスクがそれを取得できるようになります。
     * 待機列の先頭のタスクが実行されると、指定された待ち時間だけ待った後、そのタスクのPromiseが解決されます。
     */
    release(): void {
        if (this.queue.length > 0) {
            const next = this.queue.shift();
            if (next) {
                next();
            }
        } else {
            this.count--;
        }
    }
}

/**
 * 1. 初期設定:
 *    Semaphoreを最大同時実行数4（max）として、待ち時間5000ms（waitTimeMs）で初期化します。
 * 
 * 2. タスクの追加:
 *    workDaysから5つのタスクをマップします。それぞれのタスクはsem.acquire();を呼び出して、
 *    セマフォの許可を取得しようとします。
 * 
 * 3. セマフォの動作:
 *    - 最初の4つのタスク:
 *      - これらのタスクはセマフォを即座に取得できます（this.countが4未満なので）。
 *      - this.countはタスクごとに1ずつ増加します。4つのタスクが実行された後、this.countは4になります。
 *      - これらのタスクは、セマフォを取得した直後にthis._freeeService.deleteWorkRecordを呼び出します。
 * 
 *    - 5つ目のタスク:
 *      - このタスクは、this.countがすでに4になっているため、セマフォを即座に取得することはできません。
 *      - したがって、5つ目のタスクのresolve関数がqueueに追加されます。
 *      - この時点で、このタスクは待機状態になり、resolveが呼び出されるのを待ちます。
 * 
 * 4. タスクの終了:
 *    最初の4つのタスクのいずれかが終了すると、sem.release();が呼び出されます。
 * 
 * 5. 5つ目のタスクの開始:
 *    - sem.release();が呼び出されると、queueから待機中の関数（タスク）が取り出されます（この場合、5つ目のタスク）。
 *    - 取り出された関数は、5000ms（waitTimeMs）の待機時間を経てから実行されます。
 *    - 待機時間が終了すると、関数内のresolve();が呼び出され、5つ目のタスクのsem.acquire();が完了します。
 *    - これにより、5つ目のタスクはthis._freeeService.deleteWorkRecordを呼び出すことができます。
 */
