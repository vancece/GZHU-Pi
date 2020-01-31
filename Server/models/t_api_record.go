package models

import (
	"time"
)

type TApiRecord struct {
	ID        int64      `json:"id,omitempty" gorm:"primary_key"`
	Username  string    `json:"username,omitempty" gorm:"type:varchar"`
	Uri       string    `json:"uri,omitempty" gorm:"type:varchar"`
	Duration  int64     `json:"duration,omitempty" gorm:"type:real"` //耗时统计：毫秒
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"default:current_timestamp"`
}

func SaveApiRecord(r *TApiRecord) {
	db.Create(&r)
}
