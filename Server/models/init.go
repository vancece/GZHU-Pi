package models

import (
	"GZHU-Pi/env"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	vTopic = `
			create or replace VIEW v_topic as
			select t.*, u.open_id, u.gender,
				   (CASE WHEN anonymous = true
							THEN u.profile_pic
						ELSE u.avatar END)                                         as avatar,
				   (CASE WHEN anonymous = true THEN anonymity ELSE u.nickname END) as nickname,
				   (select count(*) from t_discuss where object_id = t.id)         as discussed,
				   --点赞数量
				   (select count(*) from t_relation where object_id = t.id and object = 't_topic'
					  and type = 'star')                                           as liked,
				   --查询当前主题有关用户的点赞记录
				   (select json_agg(result)
					from (select r.*, t_user.nickname, t_user.avatar
						  from t_relation r, t_user where r.created_by = t_user.id
							and r.object = 't_topic' and r.type = 'star' ) result
					where object_id = t.id)                                        as star_list,
				   --查询当前主题有关用户的认领记录
				   (select json_agg(result)
					from (select r.*, t_user.nickname, t_user.avatar
						  from t_relation r, t_user where r.created_by = t_user.id
							and r.object = 't_topic' and r.type = 'claim' ) result
					where object_id = t.id)                                        as claim_list
			from t_topic as t, t_user as u where t.created_by = u.id;
			comment on view v_topic is '主题/帖子视图';
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
			select d.*, u.open_id, u.gender,
				   (CASE WHEN anonymous = true
							THEN u.profile_pic
						ELSE u.avatar END)                                         as avatar,
				   (CASE WHEN anonymous = true THEN anonymity ELSE u.nickname END) as nickname
			from t_discuss as d, t_user as u where d.created_by = u.id;
			comment on view v_discuss is '评论视图';
	`
	vUser = `
			create or replace view v_user as (
			select u.*, s.stu_name, s.admit_year, s.class_id, s.college, s.college_id, 
				s.major, s.major_class, s.major_id
			from t_user u left join t_stu_info s on u.stu_id = s.stu_id);
			comment on view v_discuss is '学生用户视图';
	`
)

var db *gorm.DB
var sqlxDB *sqlx.DB

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

	sqlxDB = sqlx.MustOpen("postgres", dbInfo)
	err = sqlxDB.Ping()
	if err != nil {
		logs.Error(err)
		return err
	}

	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	db.DB().SetMaxIdleConns(10)

	// SetMaxOpenCons 设置数据库的最大连接数量。
	db.DB().SetMaxOpenConns(5)

	// SetConnMaxLifetime 设置连接的最大可复用时间。
	db.DB().SetConnMaxLifetime(time.Hour)

	//关闭复数表名
	db.SingularTable(true)

	if env.Conf.App.InitModels {
		t := time.Now()
		modelsInit()
		logs.Info("init models in:", time.Since(t))
	}

	return nil
}

func GetGorm() *gorm.DB {
	if db == nil {
		InitDb()
	}
	return db
}

func GetSqlx() *sqlx.DB {
	if sqlxDB == nil {
		InitDb()
	}
	return sqlxDB
}

func modelsInit() {
	logs.Info("models initializing ...")

	//自动迁移 只会 创建表、缺失的列、缺失的索引，不会 更改现有列的类型或删除未使用的列
	db.AutoMigrate(&TStuInfo{}, &TGrade{}, &TApiRecord{}, &TUser{},
		&TTopic{}, &TDiscuss{}, &TRelation{})

	db.Model(&TGrade{}).AddUniqueIndex("t_grade_stu_id_course_id_jxb_id_idx",
		"stu_id", "course_id", "jxb_id")

	db.Model(&TRelation{}).AddUniqueIndex("t_relation_object_object_id_type_created_by_idx",
		"object", "object_id", "type", "created_by")

	db.Exec(vTopic)
	db.Exec(vDiscuss)
	db.Exec(vUser)

	t := time.Now()
	db.Exec(vGrade)
	logs.Debug("init view vGrade in", time.Since(t))

}
