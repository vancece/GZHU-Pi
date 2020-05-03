/**
 * @File: models
 * @Author: Shaw
 * @Date: 2020/5/3 5:22 PM
 * @Desc

 */

package env

import (
	"github.com/jmoiron/sqlx/types"
	"gopkg.in/guregu/null.v3"
	"time"
)

type TApiRecord struct {
	ID        int64     `json:"id,omitempty" gorm:"primary_key"`
	Username  string    `json:"username,omitempty" gorm:"type:varchar"`
	Uri       string    `json:"uri,omitempty" gorm:"type:varchar"`
	Duration  int64     `json:"duration,omitempty" gorm:"type:real"` //耗时统计：毫秒
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"default:current_timestamp"`
}

type TDiscuss struct {
	ID       int64       `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	ObjectID int64       `json:"object_id,omitempty" remark:"主题对象记录ID" gorm:"type:bigint;not null"`
	Content  null.String `json:"content,omitempty" remark:"主体内容" gorm:"type:varchar"`
	ReplyID  null.Int    `json:"reply_id,omitempty" remark:"回复留言id" gorm:"type:bigint"` // reference to t_discuss(id)

	//Mark      null.Int       `json:"mark,omitempty" remark:"打星、评分" gorm:"type:smallint"`
	Image     types.JSONText `json:"image,omitempty" remark:"图片地址[string]" gorm:"type:varchar[]"`
	Anonymous null.Bool      `json:"anonymous,omitempty" remark:"是否匿名" gorm:"type:bool;default:false"`
	Anonymity null.String    `json:"anonymity,omitempty" remark:"匿名/化名" gorm:"type:varchar"`

	//Type      null.String    `json:"type,omitempty" remark:"留言类型(普通/互动)" gorm:"type:varchar"`
	Addi      types.JSONText `json:"addi,omitempty" remark:"附加信息" gorm:"type:jsonb"`
	Status    null.Int       `json:"status,omitempty" remark:"状态" gorm:"type:real;default:0"`
	CreatedBy null.Int       `json:"created_by,omitempty" remark:"创建者" gorm:"type:real"`
	CreatedAt time.Time      `json:"created_at,omitempty" remark:"创建时间" gorm:"default:current_timestamp"`
}

type TGrade struct {
	ID       int64  `json:"id,omitempty" remark:"id" gorm:"primary_key"`
	StuID    string `json:"stu_id" remark:"学号" gorm:"type:varchar;not null"`
	CourseID string `json:"course_id" remark:"课程ID" gorm:"type:varchar;not null"`
	JxbID    string `json:"jxb_id" remark:"教学班ID" gorm:"type:varchar;not null"`

	Credit     float64 `json:"credit" remark:"学分" gorm:"type:numeric(5,2)"`
	CourseGpa  float64 `json:"course_gpa" remark:"课程绩点" gorm:"type:numeric(5,2)"`
	GradeValue float64 `json:"grade_value" remark:"成绩分数" gorm:"type:numeric(5,2)"`
	Grade      string  `json:"grade" remark:"成绩" gorm:"type:varchar"`
	CourseName string  `json:"course_name" remark:"课程名称" gorm:"type:varchar"`
	CourseType string  `json:"course_type" remark:"课程类型" gorm:"type:varchar"`
	ExamType   string  `json:"exam_type" remark:"考试类型" gorm:"type:varchar"`
	Invalid    string  `json:"invalid" remark:"是否作废" gorm:"type:varchar"`
	Semester   string  `json:"semester" remark:"学期" gorm:"type:varchar"`
	Teacher    string  `json:"teacher" remark:"教师" gorm:"type:varchar"`
	Year       string  `json:"year" remark:"学年如2018-2019" gorm:"type:varchar"`
	YearSem    string  `json:"year_sem" remark:"学年学期" gorm:"type:varchar"`

	CreatedAt time.Time `json:"created_at,omitempty" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"default:current_timestamp"`
}

//用户与主题的关系记录 可以用以点赞、参与等
type TRelation struct {
	ID       int64       `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	Object   null.String `json:"object,omitempty" remark:"实体表名"`
	ObjectID int64       `json:"object_id,omitempty" remark:"主题对象记录ID" gorm:"type:bigint;not null"`
	//star点赞 claim认领 favourite收藏
	Type null.String `json:"type,omitempty" remark:"关系类型" gorm:"type:varchar"`

	CreatedBy null.Int  `json:"created_by,omitempty" remark:"创建者" gorm:"type:bigint"`
	CreatedAt time.Time `json:"created_at,omitempty" remark:"创建时间" gorm:"default:current_timestamp"`
}

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

type TUser struct {
	ID         int64          `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	MinappID   null.Int       `json:"minapp_id,omitempty" remark:"知晓云用户id" gorm:"type:bigint;unique_index"`
	OpenID     null.String    `json:"open_id,omitempty" remark:"微信openid" gorm:"type:varchar;unique_index;not null"`
	UnionID    null.String    `json:"union_id,omitempty" remark:"微信unionid" gorm:"type:varchar;unique"`
	StuID      null.String    `json:"stu_id,omitempty" remark:"学号" gorm:"type:varchar"`
	RoleID     null.Int       `json:"role_id,omitempty" remark:"用户角色id" gorm:"type:smallint"`
	Avatar     null.String    `json:"avatar,omitempty" remark:"微信头像" gorm:"type:varchar"`
	ProfilePic null.String    `json:"profile_pic,omitempty" remark:"系统随机头像" gorm:"type:varchar"`
	Nickname   null.String    `json:"nickname,omitempty" remark:"昵称" gorm:"type:varchar"`
	City       null.String    `json:"city,omitempty" remark:"城市" gorm:"type:varchar"`
	Province   null.String    `json:"province,omitempty" remark:"省份" gorm:"type:varchar"`
	Country    null.String    `json:"country,omitempty" remark:"国家" gorm:"type:varchar"`
	Gender     null.Int       `json:"gender,omitempty" remark:"性别" gorm:"type:smallint"`
	Language   null.String    `json:"language,omitempty" remark:"语言" gorm:"type:varchar"`
	Phone      null.String    `json:"phone,omitempty" remark:"手机号码" gorm:"type:varchar"`
	Tag        types.JSONText `json:"tag,omitempty" remark:"身份标签" gorm:"type:varchar[]"`
	CreatedAt  time.Time      `json:"created_at,omitempty" remark:"创建时间" gorm:"default:current_timestamp"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty" remark:"更新时间" gorm:"default:current_timestamp"`
}

type VUser struct {
	ID         int64          `json:"id,omitempty" remark:"自增id" gorm:"primary_key"`
	MinappID   null.Int       `json:"minapp_id,omitempty" remark:"知晓云用户id" gorm:"type:bigint;unique_index"`
	OpenID     null.String    `json:"open_id,omitempty" remark:"微信openid" gorm:"type:varchar;unique_index;not null"`
	UnionID    null.String    `json:"union_id,omitempty" remark:"微信unionid" gorm:"type:varchar;unique"`
	StuID      null.String    `json:"stu_id,omitempty" remark:"学号" gorm:"type:varchar"`
	RoleID     null.Int       `json:"role_id,omitempty" remark:"用户角色id" gorm:"type:smallint"`
	Avatar     null.String    `json:"avatar,omitempty" remark:"头像" gorm:"type:varchar"`
	ProfilePic null.String    `json:"profile_pic,omitempty" remark:"系统随机头像" gorm:"type:varchar"`
	Nickname   null.String    `json:"nickname,omitempty" remark:"昵称" gorm:"type:varchar"`
	City       null.String    `json:"city,omitempty" remark:"城市" gorm:"type:varchar"`
	Province   null.String    `json:"province,omitempty" remark:"省份" gorm:"type:varchar"`
	Country    null.String    `json:"country,omitempty" remark:"国家" gorm:"type:varchar"`
	Gender     null.Int       `json:"gender,omitempty" remark:"性别" gorm:"type:smallint"`
	Language   null.String    `json:"language,omitempty" remark:"语言" gorm:"type:varchar"`
	Phone      null.String    `json:"phone,omitempty" remark:"手机号码" gorm:"type:varchar"`
	Tag        types.JSONText `json:"tag,omitempty" remark:"身份标签" gorm:"type:varchar[]"`
	CreatedAt  time.Time      `json:"created_at,omitempty" remark:"创建时间" gorm:"default:current_timestamp"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty" remark:"更新时间" gorm:"default:current_timestamp"`

	StuName    string `json:"stu_name,omitempty" remark:"姓名" gorm:"type:varchar"`
	AdmitYear  string `json:"admit_year,omitempty" remark:"年级" gorm:"type:varchar"`
	ClassID    string `json:"class_id,omitempty" remark:"班级id" gorm:"type:varchar"`
	College    string `json:"college,omitempty" remark:"学院" gorm:"type:varchar"`
	CollegeID  string `json:"college_id,omitempty" remark:"学院id" gorm:"type:varchar"`
	Major      string `json:"major,omitempty" remark:"专业" gorm:"type:varchar"`
	MajorClass string `json:"major_class,omitempty" remark:"专业班级" gorm:"type:varchar"`
	MajorID    string `json:"major_id,omitempty" remark:"专业id" gorm:"type:varchar"`
}
