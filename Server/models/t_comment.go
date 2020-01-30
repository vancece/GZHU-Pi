package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx/types"
	"time"
)

type TComment struct {
	ID       int64          `json:"id" remark:"自增id" gorm:"primary_key"`
	ObjectID string         `json:"object_id" remark:"主题对象记录ID" gorm:"type:real;not null"`
	Content  sql.NullString `json:"content" remark:"主体内容" gorm:"type:varchar"`
	//Mention  sql.NullString `json:"mention" remark:"提及@用户" gorm:"type:varchar"`
	ReplyID sql.NullInt64 `json:"reply_id" remark:"回复id" gorm:"type:real"` // reference to t_comment(id)

	Image     types.JSONText `json:"image" remark:"图片地址[string]" gorm:"type:varchar[]"`
	Anonymous sql.NullBool   `json:"anonymous" remark:"是否匿名" gorm:"type:bool;default:false"`
	Anonymity sql.NullString `json:"anonymity" remark:"匿名/化名" gorm:"type:varchar"`

	Type      sql.NullString `json:"type" remark:"留言类型" gorm:"type:varchar"`
	Addi      types.JSONText `json:"addi" remark:"附加信息" gorm:"type:jsonb"`
	Status    sql.NullInt64  `json:"status" remark:"状态" gorm:"type:real;default:0"`
	CreatedBy sql.NullInt64  `json:"created_by" remark:"创建者" gorm:"type:real"`
	CreatedAt time.Time      `json:"created_at" remark:"创建时间" gorm:"default:current_timestamp"`
}
