package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"time"
)

type TStuInfo struct {
	ID         int64  `json:"id,omitempty" remark:"id" gorm:"primary_key"`
	StuID      string `json:"stu_id,omitempty" remark:"学号" gorm:"type:varchar;unique_index;not null"`
	StuName    string `json:"stu_name,omitempty" remark:"姓名" gorm:"type:varchar"`
	AdmitYear  string `json:"admit_year,omitempty" remark:"年级" gorm:"type:varchar"`
	ClassID    string `json:"class_id,omitempty" remark:"班级id" gorm:"type:varchar"`
	College    string `json:"college,omitempty" remark:"学院" gorm:"type:varchar"`
	CollegeID  string `json:"college_id,omitempty" remark:"学院id" gorm:"type:varchar"`
	Major      string `json:"major,omitempty" remark:"专业" gorm:"type:varchar"`
	MajorClass string `json:"major_class,omitempty" remark:"专业班级" gorm:"type:varchar"`
	MajorID    string `json:"major_id,omitempty" remark:"专业id" gorm:"type:varchar"`

	CreatedAt time.Time `json:"created_at,omitempty" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"default:current_timestamp"`
}

func SaveStuInfo(info *TStuInfo) {
	//获取匹配的第一条记录, 否则根据给定的条件创建一个新的记录 (仅支持 struct 和 map 条件)
	//db.FirstOrCreate(info, &info)
	var stu TStuInfo
	err := db.Where("stu_id = ?", info.StuID).First(&stu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = db.FirstOrCreate(info, &info).Error
			if err != nil {
				logs.Error(err)
				return
			}
			return
		}
		logs.Error(err)
		return
	}
	if stu.ClassID != info.ClassID || stu.College != info.College {
		logs.Info("%s 更新信息", stu.StuID)
		err = db.Model(&stu).Where("stu_id = ?", info.StuID).Update(info).Error
		if err != nil {
			logs.Error(err)
			return
		}
	}
}
