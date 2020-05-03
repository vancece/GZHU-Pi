/**
 * @File: msg-queue
 * @Author: Shaw
 * @Date: 2020/5/3 6:03 PM
 * @Desc: 创建消息队列

 */

package env

import (
	"GZHU-Pi/services/kafka"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"math"
	"reflect"
	"strconv"
)

const (
	QueueTopicGrade   = "kafka-queue-grade"
	QueueTopicStuInfo = "kafka-queue-student-info"

	keyQueueGrade   = "gzhupi:queue:" + QueueTopicGrade
	keyQueueStuInfo = "gzhupi:queue:" + QueueTopicStuInfo
)

//成绩保存更新队列
func GradeQueue() (err error) {

	val, err := RedisCli.Get(keyQueueGrade).Result()
	if err != nil && err != redis.Nil {
		logs.Error(err)
		return
	}
	var offset int64
	if val != "" {
		offset, _ = strconv.ParseInt(val, 10, 64)
		offset = offset + 1
	}
	logs.Info("%s offset: %d", QueueTopicGrade, offset)

	h := &kafka.CustomHandler{
		Topic:     QueueTopicGrade,
		CustomFun: CustomGradeMsg,
		ErrorFun:  ErrHandler,
		Offset:    offset,
	}

	err = Kafka.AddCustomer(h)
	if err != nil {
		logs.Error(err)
		return
	}
	return
}

//学生信息存储队列
func StuInfoQueue() (err error) {

	val, err := RedisCli.Get(keyQueueStuInfo).Result()
	if err != nil && err != redis.Nil {
		logs.Error(err)
		return
	}
	var offset int64
	if val != "" {
		offset, _ = strconv.ParseInt(val, 10, 64)
		offset = offset + 1
	}
	logs.Info("%s offset: %d", QueueTopicStuInfo, offset)

	h := &kafka.CustomHandler{
		Topic:     QueueTopicStuInfo,
		CustomFun: CustomInfo,
		ErrorFun:  ErrHandler,
		Offset:    offset,
	}

	err = Kafka.AddCustomer(h)
	if err != nil {
		logs.Error(err)
		return
	}
	return
}

func CustomGradeMsg(msg *sarama.ConsumerMessage) (err error) {
	if msg == nil {
		return
	}
	var grades []*TGrade
	err = json.Unmarshal(msg.Value, &grades)
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Info("消费成功 topic: %s offset: %d", msg.Topic, msg.Offset)
	SaveOrUpdateGrade(grades)

	//消费成功 记录更新offset
	err = RedisCli.Set(keyQueueGrade, fmt.Sprint(msg.Offset), 0).Err()
	if err != nil {
		logs.Error(err)
		return
	}
	return
}

func CustomInfo(msg *sarama.ConsumerMessage) (err error) {
	if msg == nil {
		return
	}
	var info *TStuInfo
	err = json.Unmarshal(msg.Value, &info)
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Info("消费成功 topic: %s offset: %d", msg.Topic, msg.Offset)
	SaveStuInfo(info)

	//消费成功 记录更新offset
	err = RedisCli.Set(keyQueueStuInfo, fmt.Sprint(msg.Offset), 0).Err()
	if err != nil {
		logs.Error(err)
		return
	}
	return
}

func ErrHandler(h *kafka.CustomHandler, err error) {
	if err != nil && h != nil {
		logs.Error("%s 消息发生错误：%s", h.Topic, err)
	}
}

func SaveOrUpdateGrade(grades []*TGrade) {
	db := GetGorm()
	for _, v := range grades {
		//根据主键查询
		var res = TGrade{}
		result := db.Where("stu_id = ? and course_id = ? and jxb_id = ?",
			v.StuID, v.CourseID, v.JxbID).First(&res)

		//不存在记录则插入
		if result.Error == gorm.ErrRecordNotFound {
			logs.Debug("create record for course_id %s", v.CourseID)
			db.Create(v)
			continue
		}
		//存在记录但没有变动，跳过
		if math.Round(res.CourseGpa*10)/10 == v.CourseGpa &&
			res.GradeValue == v.GradeValue &&
			res.Grade == v.Grade &&
			res.Credit == v.Credit &&
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

		result = db.Model(&res).Where("stu_id = ? and course_id = ? and jxb_id = ?",
			v.StuID, v.CourseID, v.JxbID).Updates(m)
		if result.Error != nil {
			logs.Error(result.Error)
			continue
		}
		logs.Debug("update record: %s %s %s ", v.StuID, v.CourseID, v.JxbID)
	}
}

func SaveStuInfo(info *TStuInfo) {
	db := GetGorm()
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
