package models

import (
	"github.com/astaxie/beego/logs"
	"math"
	"reflect"
	"time"
)

type TGrade struct {
	StuID    string `json:"stu_id,omitempty" remark:"学号" gorm:"primary_key;type:varchar"`
	CourseID string `json:"course_id,omitempty" remark:"课程ID" gorm:"primary_key;type:varchar"`
	JxbID    string `json:"jxb_id,omitempty" remark:"教学班ID" gorm:"primary_key;type:varchar"`

	Credit     float64 `json:"credit,omitempty" remark:"学分" gorm:"type:real"`
	CourseGpa  float64 `json:"course_gpa,omitempty" remark:"课程绩点" gorm:"type:real"`
	GradeValue float64 `json:"grade_value,omitempty" remark:"成绩分数" gorm:"type:real"`
	Grade      string  `json:"grade,omitempty" remark:"成绩" gorm:"type:varchar"`
	CourseName string  `json:"course_name,omitempty" remark:"课程名称" gorm:"type:varchar"`
	CourseType string  `json:"course_type,omitempty" remark:"课程类型" gorm:"type:varchar"`
	ExamType   string  `json:"exam_type,omitempty" remark:"考试类型" gorm:"type:varchar"`
	Invalid    string  `json:"invalid,omitempty" remark:"是否作废" gorm:"type:varchar"`
	Semester   string  `json:"semester,omitempty" remark:"学期" gorm:"type:varchar"`
	Teacher    string  `json:"teacher,omitempty" remark:"教师" gorm:"type:varchar"`
	Year       string  `json:"year,omitempty" remark:"学年如2018-2019" gorm:"type:varchar"`
	YearSem    string  `json:"year_sem,omitempty" remark:"学年学期" gorm:"type:varchar"`

	CreatedAt time.Time `json:"created_at,omitempty" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"default:current_timestamp"`
}

func SaveOrUpdateGrade(grades []*TGrade) {

	for _, v := range grades {
		//根据主键查询
		res := &TGrade{StuID: v.StuID, CourseID: v.CourseID, JxbID: v.JxbID}
		db.First(res)
		//不存在记录则插入
		if res.CourseGpa == 0 &&
			res.GradeValue == 0 &&
			res.Grade == "" &&
			res.Invalid == "" {
			db.Create(v)
			continue
		}
		//存在记录但没有变动，跳过
		if math.Round(res.CourseGpa*10)/10 == v.CourseGpa &&
			res.GradeValue == v.GradeValue &&
			res.Grade == v.Grade &&
			res.Invalid == v.Invalid {
			continue
		}
		//更新记录 结构体转换为map
		m := make(map[string]interface{})
		elem := reflect.ValueOf(v).Elem()
		relType := elem.Type()
		for i := 0; i < relType.NumField(); i++ {
			m[relType.Field(i).Name] = elem.Field(i).Interface()
		}
		delete(m, "CreatedAt")
		delete(m, "UpdatedAt")

		db.Model(&res).Updates(m)
		logs.Debug("update record: %s %s %s ", v.StuID, v.CourseID, v.JxbID)
	}
}
