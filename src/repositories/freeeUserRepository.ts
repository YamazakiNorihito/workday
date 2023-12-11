import { SchemaFieldTypes, RedisClientType } from 'redis';
import { inject, singleton } from "tsyringe";
import { FreeeOAuthTokenResponse } from '../services/freeeHrService';

@singleton()
export class FreeeUserRepository {
    constructor(
        @inject("RedisClient") private readonly redisClient: RedisClientType
    ) {
        (async () => {
            await this.createIndex();
        })();
    }

    public isReady(): boolean {
        return this.redisClient.isReady;
    }

    public async save(userId: string, user: User): Promise<void> {
        const userKey = `freeeuser:${userId}`;
        let userJsonString = JSON.stringify(user);
        const result = await this.redisClient.set(userKey, userJsonString);
    }

    public async get(userId: string): Promise<User | null> {
        const userKey = `freeeuser:${userId}`;
        const userJsonString = await this.redisClient.get(userKey);

        if (!userJsonString) {
            return null;
        }
        return JSON.parse(userJsonString) as User;
    }

    private async createIndex(): Promise<void> {
        try {
            await this.redisClient.ft.create('idx:freeeusers', {
                '$.id': {
                    type: SchemaFieldTypes.NUMERIC,
                    SORTABLE: true
                },
                '$.updated_at': {
                    type: SchemaFieldTypes.NUMERIC,
                    SORTABLE: true
                },
                '$.companies[*].id': {
                    type: SchemaFieldTypes.NUMERIC
                },
                '$.companies[*].name': {
                    type: SchemaFieldTypes.TEXT
                },
                '$.companies[*].role': {
                    type: SchemaFieldTypes.TEXT
                },
                '$.companies[*].external_cid': {
                    type: SchemaFieldTypes.NUMERIC
                },
                '$.companies[*].employee_id': {
                    type: SchemaFieldTypes.NUMERIC
                },
                '$.companies[*].display_name': {
                    type: SchemaFieldTypes.TEXT
                },
                '$.oauth.access_token': {
                    type: SchemaFieldTypes.TEXT
                },
                '$.oauth.token_type': {
                    type: SchemaFieldTypes.TEXT
                },
                '$.oauth.expires_in': {
                    type: SchemaFieldTypes.NUMERIC,
                    SORTABLE: true
                },
                '$.oauth.refresh_token': {
                    type: SchemaFieldTypes.TEXT
                },
                '$.oauth.scope': {
                    type: SchemaFieldTypes.TEXT
                },
                '$.oauth.created_at': {
                    type: SchemaFieldTypes.NUMERIC,
                    SORTABLE: true
                },
                '$.oauth.company_id': {
                    type: SchemaFieldTypes.NUMERIC
                }
            }, {
                ON: 'JSON',
                PREFIX: 'freeeuser:'
            });
        } catch (e: any) {
            if (e.message === 'Index already exists') {
                //console.log('Index exists already, skipped creation.');
            } else {
                throw e;
            }
        }
    }
}

export interface OAuth {
    access_token: string;
    token_type: string;
    expires_in: number;
    refresh_token: string;
    scope: string;
    created_at: number;
    company_id: number;
}

export interface Company {
    id: number;
    name: string;
    role: string;
    external_cid: number;
    employee_id?: number | null;
    display_name?: string | null;
}

export interface User {
    id: number;
    companies: Company[];
    oauth: FreeeOAuthTokenResponse;
    updated_at: number;
}