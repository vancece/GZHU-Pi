package env

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"log"
	"os"
	"strings"
)

func InitLogger(path string) error {
	if path == "" {
		path = "/tmp/log/"
	}
	//创建日志目录
	path = strings.TrimRight(path, "/") + "/"
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal("创建日志目录失败 ", err)
		return err
	}
	fileName := path + "gzhupi.log"
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("创建日志文件失败 ", err)
		return err
	}
	defer f.Close()

	logs.SetLogFuncCallDepth(3)    //调用层级
	logs.EnableFuncCallDepth(true) //输出文件名和行号
	//logs.Async()                   //提升性能, 可以设置异步输出

	config := make(map[string]interface{})
	config["filename"] = fileName

	logs.SetLevel(logs.LevelDebug)

	configStr, err := json.Marshal(config)
	if err != nil {
		log.Fatal("initLogger failed, marshal err:", err)
		return err
	}
	err = logs.SetLogger(logs.AdapterConsole, "") //控制台输出
	if err != nil {
		log.Fatal("SetLogger failed, err:", err)
		return err
	}
	err = logs.SetLogger(logs.AdapterFile, string(configStr)) //文件输出
	if err != nil {
		log.Fatal("SetLogger failed, err:", err)
		return err
	}
	//err = logs.SetLogger(logs.AdapterEs, `{"dsn":"http://localhost:9200/","level":1}`)
	//if err != nil {
	//	log.Fatal("SetLogger failed, err:", err)
	//	return err
	//}
	return nil
}
