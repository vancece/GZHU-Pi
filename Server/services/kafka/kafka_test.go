/**
 * @File: kafka_test.go
 * @Author: Shaw
 * @Date: 2020/5/3 3:42 PM
 * @Desc

 */

package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"testing"
	"time"
)

func TestDefaultKafka(t *testing.T) {

	//开启生产者和消费者实例
	k, err := DefaultKafka([]string{"localhost:9092"})
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	h := &CustomHandler{
		Topic:     "mykafka",
		CustomFun: customMsg,
		ErrorFun:  errHandler,
		Offset:    0,
	}

	err = k.AddCustomer(h)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	for i := 0; i < 5; i++ {
		data := &ProduceData{
			Topic: "mykafka",
			Data:  []byte("I am a message: " + fmt.Sprint(i)),
		}

		err = k.SendData(data)
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		time.Sleep(1 * time.Second)
	}

	time.Sleep(10 * time.Second)

}

func customMsg(msg *sarama.ConsumerMessage) (err error) {

	if msg == nil {
		return
	}

	logs.Info("消费了 %s:%d \n %s", msg.Topic, msg.Offset, string(msg.Value))

	return nil

}

func errHandler(h *CustomHandler, err error) {

	if err == nil {
		return
	}
	logs.Error("%s发生错误：%s", h.Topic, err.Error())
}
