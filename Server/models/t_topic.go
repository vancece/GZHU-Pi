package models

import (
	"github.com/jmoiron/sqlx/types"
	"gopkg.in/guregu/null.v3"
	"time"
)

type TTopic struct {
	ID      int64       `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	Type    null.String `json:"type,omitempty" remark:"主题类型" gorm:"type:varchar"`
	Title   null.String `json:"title,omitempty" remark:"标题" gorm:"type:varchar"`
	Content null.String `json:"content,omitempty" remark:"主体内容" gorm:"type:varchar"`
	//Mention   null.String `json:"mention,omitempty" remark:"提及@用户" gorm:"type:varchar"`
	Category  null.String    `json:"category,omitempty" remark:"归属类别" gorm:"type:varchar"`
	Image     types.JSONText `json:"image,omitempty" remark:"图片地址" gorm:"type:varchar[]"`
	Label     types.JSONText `json:"label,omitempty" remark:"标签" gorm:"type:varchar[]"`
	Viewed    null.Int       `json:"viewed,omitempty" remark:"浏览量" gorm:"type:int;default:0"`
	Anonymous null.Bool      `json:"anonymous,omitempty" remark:"是否匿名" gorm:"type:bool;default:false"`
	Anonymity null.String    `json:"anonymity,omitempty" remark:"匿名/化名" gorm:"type:varchar"`

	Addi      types.JSONText `json:"addi,omitempty" remark:"附加信息" gorm:"type:jsonb"`
	Status    null.Int       `json:"status,omitempty" remark:"状态" gorm:"type:smallint;default:0"`
	CreatedBy null.Int       `json:"created_by,omitempty" remark:"创建者" gorm:"type:bigint"`
	CreatedAt time.Time      `json:"created_at,omitempty" remark:"创建时间" gorm:"default:current_timestamp"`
	UpdatedAt time.Time      `json:"updated_at,omitempty" remark:"更新时间" gorm:"default:current_timestamp"`
}
