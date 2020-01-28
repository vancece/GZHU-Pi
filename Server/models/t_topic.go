package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx/types"
	"time"
)

type TTopic struct {
	ID        int64          `json:"id" remark:"自增id" gorm:"primary_key"`
	Type      sql.NullString `json:"type" remark:"主题类型" gorm:"type:varchar"`
	Title     sql.NullString `json:"title" remark:"标题" gorm:"type:varchar"`
	Content   sql.NullString `json:"content" remark:"主体内容" gorm:"type:varchar"`
	Mention   sql.NullString `json:"mention" remark:"提及/高亮内容" gorm:"type:varchar"`
	Category  sql.NullString `json:"category" remark:"归属类别" gorm:"type:varchar"`
	Image     types.JSONText `json:"image" remark:"图片地址[string]" gorm:"type:varchar[]"`
	Label     types.JSONText `json:"label" remark:"标签[string]" gorm:"type:varchar[]"`
	Viewed    sql.NullInt64  `json:"viewed" remark:"浏览量" gorm:"type:real;default:0"`
	Anonymous sql.NullBool   `json:"anonymous" remark:"是否匿名" gorm:"type:bool;default:false"`
	Anonymity sql.NullString `json:"anonymity" remark:"匿名/化名" gorm:"type:varchar"`
	Addi      types.JSONText `json:"addi" remark:"附加信息" gorm:"type:jsonb"`
	Status    sql.NullInt64  `json:"status" remark:"状态" gorm:"type:real;default:0"`
	CreatedBy sql.NullInt64  `json:"created_by" remark:"创建者" gorm:"type:real"`
	CreatedAt time.Time      `json:"created_at" remark:"创建时间" gorm:"default:current_timestamp"`
	UpdatedAt time.Time      `json:"updated_at" remark:"更新时间" gorm:"default:current_timestamp"`
}
