package models

import (
	"github.com/jmoiron/sqlx/types"
	"gopkg.in/guregu/null.v3"
	"time"
)

type TUser struct {
	ID        int64          `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	MinappID  null.Int       `json:"minapp_id,omitempty" remark:"知晓云用户id" gorm:"type:bigint;unique_index"`
	OpenID    null.String    `json:"open_id,omitempty" remark:"微信openid" gorm:"type:varchar;unique_index;not null"`
	UnionID   null.String    `json:"union_id,omitempty" remark:"微信unionid" gorm:"type:varchar;unique"`
	StuID     null.String    `json:"stu_id,omitempty" remark:"学号" gorm:"type:varchar"`
	RoleID    null.Int       `json:"role_id,omitempty" remark:"用户角色id" gorm:"type:smallint"`
	Avatar    null.String    `json:"avatar,omitempty" remark:"头像" gorm:"type:varchar"`
	Nickname  null.String    `json:"nickname,omitempty" remark:"昵称" gorm:"type:varchar"`
	City      null.String    `json:"city,omitempty" remark:"城市" gorm:"type:varchar"`
	Province  null.String    `json:"province,omitempty" remark:"省份" gorm:"type:varchar"`
	Country   null.String    `json:"country,omitempty" remark:"国家" gorm:"type:varchar"`
	Gender    null.Int       `json:"gender,omitempty" remark:"性别" gorm:"type:smallint"`
	Language  null.String    `json:"language,omitempty" remark:"语言" gorm:"type:varchar"`
	Phone     null.String    `json:"phone,omitempty" remark:"手机号码" gorm:"type:varchar"`
	Tag       types.JSONText `json:"tag,omitempty" remark:"身份标签" gorm:"type:varchar[]"`
	CreatedAt time.Time      `json:"created_at,omitempty" remark:"创建时间" gorm:"default:current_timestamp"`
	UpdatedAt time.Time      `json:"updated_at,omitempty" remark:"更新时间" gorm:"default:current_timestamp"`
}
