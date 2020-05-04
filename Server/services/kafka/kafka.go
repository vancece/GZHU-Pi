/**
 * @File: kafka
 * @Author: Shaw
 * @Date: 2020/5/3 2:45 AM
 * @Desc: 消息队列服务

 */

package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"strings"
)

//消费者读取目标并处理的参数
type CustomHandler struct {
	Topic     string
	Partition int32
	Offset    int64
	CustomFun func(*sarama.ConsumerMessage) error //自定义的数据消费操作函数
	ErrorFun  func(*CustomHandler, error)         //自定义错误处理函数
}

//提供给生产者的数据
type ProduceData struct {
	Topic     string
	Partition int32
	Data      []byte //写入队列的数据
}

//kafka 生产者消费者实例封装
type Kafka struct {
	//生产者数据
	data chan *ProduceData

	//消费者消费实例
	handler chan *CustomHandler

	Config *sarama.Config

	//kafka集群 broker连接地址
	Broker []string

	//同步生产者
	producer sarama.SyncProducer
	consumer sarama.Consumer

	//记录已经初始化的分区消费者，防止重复初始化 key := fmt.Sprintf("%s:%d", h.Topic, h.Partition)
	pc map[string]sarama.PartitionConsumer
}

func DefaultKafka(broker []string) (k *Kafka, err error) {
	return NewKafka(broker, nil)
}

func NewKafka(broker []string, config *sarama.Config) (k *Kafka, err error) {

	if config == nil {
		config = sarama.NewConfig()
		// 等待服务器所有副本都保存成功后的响应
		config.Producer.RequiredAcks = sarama.WaitForAll
		// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
		config.Producer.Partitioner = sarama.NewRandomPartitioner
		// 是否等待成功和失败后的响应
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true
	}

	if len(broker) == 0 {
		return nil, fmt.Errorf("broker is empty")
	}

	k = &Kafka{
		Config:  config,
		Broker:  broker,
		handler: make(chan *CustomHandler),
		data:    make(chan *ProduceData),
		pc:      make(map[string]sarama.PartitionConsumer),
	}

	// 使用给定代理地址和配置创建一个同步生产者
	k.producer, err = sarama.NewSyncProducer(k.Broker, k.Config)
	if err != nil {
		logs.Error(err)
		return
	}

	k.consumer, err = sarama.NewConsumer(k.Broker, k.Config)
	if err != nil {
		logs.Error(err)
		return
	}

	//开启生产消费实例，持续监听管道
	go k.produce()
	go k.consume()

	return
}

func (k *Kafka) AddCustomer(h *CustomHandler) (err error) {
	if h == nil {
		return
	}
	if h.Offset == 0 {
		h.Offset = sarama.OffsetNewest
	}
	k.handler <- h
	return
}

func (k *Kafka) SendData(data *ProduceData) (err error) {
	if data == nil {
		return
	}
	k.data <- data
	return
}

func (k *Kafka) produce() {

	for {
		select {
		case p := <-k.data:
			//构建发送的消息，
			msg := &sarama.ProducerMessage{
				Topic:     p.Topic,
				Partition: p.Partition,
				Value:     sarama.ByteEncoder(p.Data),
			}

			partition, offset, err := k.producer.SendMessage(msg)
			if err != nil {
				logs.Error(err)
				return
			}
			logs.Info("消息加入队列成功 topic: %s, partition: %d, offset: %d", p.Topic, partition, offset)
		}
	}
}

func (k *Kafka) consume() {

	for {
		select {
		case h := <-k.handler:

			//检查是否已经创建过 ConsumePartition
			key := fmt.Sprintf("%s:%d", h.Topic, h.Partition)
			_, ok := k.pc[key]
			if ok {
				//跳过，继续监听
				logs.Warn(key, " ConsumePartition已经初始化")
				continue
			}

			var pc sarama.PartitionConsumer
			pc, err := k.consumer.ConsumePartition(h.Topic, h.Partition, h.Offset)
			if err != nil {
				if strings.Contains(err.Error(), "outside") {
					pc, err = k.consumer.ConsumePartition(h.Topic, h.Partition, sarama.OffsetNewest)
					if err != nil {
						logs.Error(err)
						h.ErrorFun(h, err)
						break
					}
				} else {
					logs.Error(err)
					h.ErrorFun(h, err)
					break
				}
			}
			k.pc[key] = pc
			// 等待 生产者产生对应数据 然后消费
			go func() {
				for {
					select {
					case msg := <-pc.Messages():
						//调用处理函数
						err = h.CustomFun(msg)
						if err != nil {
							h.ErrorFun(h, err)
							break
						}
					case err = <-pc.Errors():
						if err != nil {
							logs.Error(err)
							h.ErrorFun(h, err)
							break
						}
					}
				}
			}()
		}
	}
}
