package models

import (
	"database/sql"
	"time"
)

type TUser struct {
	ID        int64  `json:"id" remark:"自增id" gorm:"primary_key"`
	UserID    sql.NullInt64  `json:"user_id" remark:"知晓云用户id" gorm:"type:real"`
	Openid    sql.NullString `json:"openid" remark:"微信openid" gorm:"type:varchar"`
	Unionid   sql.NullString `json:"unionid" remark:"微信openid" gorm:"type:varchar"`
	StuID     sql.NullString `json:"stu_id" remark:"学号" gorm:"type:varchar"`
	Avatar    sql.NullString `json:"avatar" remark:"头像" gorm:"type:varchar"`
	Nickname  sql.NullString `json:"nickname" remark:"昵称" gorm:"type:varchar"`
	City      sql.NullString `json:"city" remark:"城市" gorm:"type:varchar"`
	Province  sql.NullString `json:"province" remark:"省份" gorm:"type:varchar"`
	Country   sql.NullString `json:"country" remark:"国家" gorm:"type:varchar"`
	Gender    sql.NullInt64  `json:"gender" remark:"性别" gorm:"type:real"`
	Language  sql.NullString `json:"language" remark:"语言" gorm:"type:varchar"`
	Phone     sql.NullString `json:"_phone" remark:"手机号码" gorm:"type:varchar"`
	CreatedAt time.Time      `json:"created_at" remark:"创建时间" gorm:"default:current_timestamp"`
	UpdatedAt time.Time      `json:"updated_at" remark:"更新时间" gorm:"default:current_timestamp"`
}
