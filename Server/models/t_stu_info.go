package models

import "time"

type TStuInfo struct {
	StuID      string `json:"stu_id" remark:"学号" gorm:"primary_key;type:varchar"`
	StuName    string `json:"stu_name" remark:"姓名" gorm:"type:varchar"`
	AdmitYear  string `json:"admit_year" remark:"年级" gorm:"type:varchar"`
	ClassID    string `json:"class_id" remark:"班级id" gorm:"type:varchar"`
	College    string `json:"college" remark:"学院" gorm:"type:varchar"`
	CollegeID  string `json:"college_id" remark:"学院id" gorm:"type:varchar"`
	Major      string `json:"major" remark:"专业" gorm:"type:varchar"`
	MajorClass string `json:"major_class" remark:"专业班级" gorm:"type:varchar"`
	MajorID    string `json:"major_id" remark:"专业id" gorm:"type:varchar"`

	CreatedAt time.Time `json:"-" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"-" gorm:"default:current_timestamp"`
}

func SaveStuInfo(info *TStuInfo) {
	//获取匹配的第一条记录, 否则根据给定的条件创建一个新的记录 (仅支持 struct 和 map 条件)
	db.FirstOrCreate(info, &info)
}
