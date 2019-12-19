package models

import (
	"time"
)

type TApiRecord struct {
	ID        uint      `gorm:"primary_key"`
	Username  string    `json:"username" gorm:"type:varchar"`
	Uri       string    `json:"uri" gorm:"type:varchar"`
	Duration  int64     `json:"duration" gorm:"type:real"` //耗时统计：毫秒
	CreatedAt time.Time `json:"created_at" gorm:"default:current_timestamp"`
}

func SaveApiRecord(r *TApiRecord) {
	db.Create(&r)
}
