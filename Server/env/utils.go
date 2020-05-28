/**
 * @File: utils
 * @Author: Shaw
 * @Date: 2020/5/28 5:33 PM
 * @Desc

 */

package env

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron/v3"
	"io"
)

//对字符串进行MD5哈希
func StringMD5(data string) string {
	t := md5.New()
	_, _ = io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//对字符串进行SHA1哈希
func StringSha1(data string) string {
	t := sha1.New()
	_, _ = io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func CornTask(spec string, task func()) {

	cronTab := cron.New()
	// 添加定时任务, "* 0/5 7-21 * *" 是 cronTab,表示7-21点，每五分钟
	_, err := cronTab.AddFunc(spec, task)
	if err != nil {
		logs.Error(err)
		return
	}
	// 启动定时器
	cronTab.Start()
	// 定时任务是另起协程执行的,这里使用 select 简单阻塞.实际开发中需要根据实际情况进行控制
	select {}
}
