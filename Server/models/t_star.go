package models

import (
	"database/sql"
	"time"
)

//用户与主题的关系记录 可以用以点赞、参与等
type TRelation struct {
	ID       int64  `json:"id" remark:"自增id" gorm:"primary_key"`
	ObjectID string `json:"object_id" remark:"主题对象记录ID" gorm:"type:real;not null"`
	// star 点赞记录  claim 认领记录
	Type sql.NullString `json:"type" remark:"关系类型" gorm:"type:varchar"`

	CreatedBy sql.NullInt64 `json:"created_by" remark:"创建者" gorm:"type:real"`
	CreatedAt time.Time     `json:"created_at" remark:"创建时间" gorm:"default:current_timestamp"`
}
