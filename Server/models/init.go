package models

import (
	"GZHU-Pi/env"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"log"
)

const (
	v_topic = `
			create or replace VIEW v_topic as
			select t.*,
				   u.id                                                            as uid,
				   u.gender,
				   (CASE WHEN anonymous = true
							THEN 'https://shaw-1256261760.cos.ap-guangzhou.myqcloud.com/gzhu-pi/images/icon/anonmous_avatar.png'
						ELSE u.avatar END)                                         as avatar,
				   (CASE WHEN anonymous = true THEN anonymity ELSE u.nickname END) as nickname,
			
				   (select count(*) from t_comment where object_id = t.id)         as comment_count,
				   (select count(*) from t_relation where object_id = t.id)        as licked
			
			from t_topic as t, t_user as u
			where t.created_by = u.id;
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
		log.Print(err)
		return err
	}
	logs.Info("数据库：%s:%d", d.Host, d.Port)
	//关闭复数表名
	db.SingularTable(true)

	//自动迁移 只会 创建表、缺失的列、缺失的索引，不会 更改现有列的类型或删除未使用的列
	db.AutoMigrate(&TStuInfo{}, &TGrade{}, &TApiRecord{}, &TUser{},
		&TTopic{}, &TComment{}, &TRelation{})

	db.Exec(v_topic)
	return nil
}

func GetGorm() *gorm.DB {
	return db
}
