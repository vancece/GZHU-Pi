package models

import (
	"gopkg.in/guregu/null.v3"
	"time"
)

//用户与主题的关系记录 可以用以点赞、参与等
type TRelation struct {
	ID       int64  `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	ObjectID string `json:"object_id,omitempty" remark:"主题对象记录ID" gorm:"type:bigint;not null"`
	// star 点赞记录  claim 认领记录
	Type null.String `json:"type,omitempty" remark:"关系类型" gorm:"type:varchar"`

	CreatedBy null.Int  `json:"created_by,omitempty" remark:"创建者" gorm:"type:bigint"`
	CreatedAt time.Time `json:"created_at,omitempty" remark:"创建时间" gorm:"default:current_timestamp"`
}
