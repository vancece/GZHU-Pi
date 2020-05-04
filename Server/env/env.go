/**
 * @File: env
 * @Author: Shaw
 * @Date: 2020/5/3 4:52 PM
 * @Desc

 */

package env

import (
	"github.com/astaxie/beego/logs"
	"log"
)

func EnvInit() {

	//------1、初始化配置文件------
	err := InitConfigure()
	if err != nil {
		log.Fatal(err)
	}

	//------2、初始化日志------
	err = InitLogger("/tmp/log/")
	if err != nil {
		log.Fatal(err)
	}

	//------3、初始化数据库------
	err = InitDb()
	if err != nil {
		log.Fatal(err)
	}

	//------4、初始化redis------
	err = InitRedis()
	if err != nil {
		log.Fatal(err)
	}

	//------5、初始化kafka------
	err = InitKafka()
	if err != nil {
		log.Fatal(err)
	}

	//------6、初始化性能监控------
	go MemoryCollector()

	logs.Info("======== finish to init app env ========")
}
