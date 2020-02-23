package models

import (
	"GZHU-Pi/env"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	vTopic = `
			create or replace VIEW v_topic as
			select t.*, u.gender,
				   (CASE WHEN anonymous = true
							THEN 'https://shaw-1256261760.cos.ap-guangzhou.myqcloud.com/gzhu-pi/images/icon/anonmous_avatar.png'
						ELSE u.avatar END)                                         as avatar,
				   (CASE WHEN anonymous = true THEN anonymity ELSE u.nickname END) as nickname,
			
				   (select count(*) from t_discuss where object_id = t.id)         as discussed,
				   (select count(*) from t_relation where object_id = t.id)        as liked
			
			from t_topic as t, t_user as u
			where t.created_by = u.id;
			comment on view v_grade is '主题/帖子视图';
	`
	vGrade = `
			create or replace view v_grade (stu_id, class_id, major_class, major_id, major, stu_name, college_id, college, admit_year, year,
						semester, course_id, course_name, credit, grade_value, grade, course_gpa, course_type, exam_type, invalid, jxb_id,
						teacher, year_sem, created_at) as SELECT s.stu_id, s.class_id, s.major_class, s.major_id, s.major, s.stu_name,
						  s.college_id, s.college,s.admit_year, g.year, g.semester, g.course_id, g.course_name, g.credit,g.grade_value,
						  g.grade, g.course_gpa, g.course_type, g.exam_type, g.invalid,g.jxb_id, g.teacher, g.year_sem, g.created_at
			FROM t_stu_info s, t_grade g WHERE ((s.stu_id)::text = (g.stu_id)::text);
			comment on view v_grade is '学生成绩视图';
	`
	vDiscuss = `
			create or replace VIEW v_discuss as
			select d.*, u.gender,
				   (CASE WHEN anonymous = true
							THEN 'https://shaw-1256261760.cos.ap-guangzhou.myqcloud.com/gzhu-pi/images/icon/anonmous_avatar.png'
						ELSE u.avatar END)                                         as avatar,
				   (CASE WHEN anonymous = true THEN anonymity ELSE u.nickname END) as nickname
			from t_discuss as d, t_user as u where d.created_by = u.id;
			comment on view v_grade is '评论视图';
	`
)

var db *gorm.DB

func InitDb() error {
	d := env.Conf.Db
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Dbname, d.SslMode)

	var err error
	db, err = gorm.Open("postgres", dbInfo)
	if err != nil {
		logs.Error(err)
		return err
	}
	logs.Info("数据库：%s:%d/%s", d.Host, d.Port, d.Dbname)

	modelsInit()

	return nil
}

func GetGorm() *gorm.DB {
	return db
}

func modelsInit() {
	//关闭复数表名
	db.SingularTable(true)

	//自动迁移 只会 创建表、缺失的列、缺失的索引，不会 更改现有列的类型或删除未使用的列
	db.AutoMigrate(&TStuInfo{}, &TGrade{}, &TApiRecord{}, &TUser{},
		&TTopic{}, &TDiscuss{}, &TRelation{})

	db.Model(&TGrade{}).AddUniqueIndex("t_grade_stu_id_course_id_jxb_id_idx",
		"stu_id", "course_id", "jxb_id")

	db.Exec(vTopic)
	db.Exec(vGrade)
	db.Exec(vDiscuss)
}
