package metadata

import "time"

type UserMeta struct {
	ID   string `json:"id" dynamodbav:"id"`
	Name string `json:"name" dynamodbav:"name"`
}

type CreateBy UserMeta
type CreateAt = time.Time
type UpdateBy UserMeta
type UpdateAt = time.Time
