/**
 * @File: kafka
 * @Author: Shaw
 * @Date: 2020/5/3 4:39 PM
 * @Desc

 */

package env

import (
	"GZHU-Pi/services/kafka"
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

var Kafka *kafka.Kafka

func InitKafka() (err error) {

	broker := []string{fmt.Sprintf("ifeel.vip:9092")}

	logs.Info("connecting to broker: %v", broker)

	Kafka, err = kafka.DefaultKafka(broker)
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Info("kafka初始化成功：%v", broker)

	go func() {
		time.Sleep(3 * time.Second)
		err = StuInfoQueue()
		if err != nil {
			return
		}
		err = GradeQueue()
		if err != nil {
			return
		}
	}()

	return
}
