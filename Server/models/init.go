package models

import (
	"GZHU-Pi/env"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"log"
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
	db.AutoMigrate(&TStuInfo{}, &TGrade{}, &TApiRecord{}, &TUser{}, &TTopic{})
	return nil
}
