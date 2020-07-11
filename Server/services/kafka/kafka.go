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
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
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

		//提交offset的间隔时间，每秒提交一次给kafka
		config.Consumer.Offsets.CommitInterval = 1 * time.Second
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
	if h == nil || k.consumer == nil {
		return
	}
	if h.Offset == 0 {
		h.Offset = sarama.OffsetNewest
	}
	k.handler <- h
	return
}

func (k *Kafka) SendData(data *ProduceData) (err error) {
	if data == nil || k.producer == nil {
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
				continue
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
			pc, err := k.consumer.ConsumePartition(h.Topic, h.Partition, sarama.OffsetNewest)
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

func SaramaConsumer() {

	fmt.Println("start consume")
	config := sarama.NewConfig()

	//提交offset的间隔时间，每秒提交一次给kafka
	config.Consumer.Offsets.CommitInterval = 1 * time.Second

	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	config.Version = sarama.V0_10_0_1

	//consumer新建的时候会新建一个client，这个client归属于这个consumer，并且这个client不能用作其他的consumer
	consumer, err := sarama.NewConsumer([]string{"182.61.9.153:6667", "182.61.9.154:6667", "182.61.9.155:6667"}, config)
	if err != nil {
		panic(err)
	}

	//新建一个client，为了后面offsetManager做准备
	client, err := sarama.NewClient([]string{"182.61.9.153:6667", "182.61.9.154:6667", "182.61.9.155:6667"}, config)
	if err != nil {
		panic("client create error")
	}
	defer client.Close()

	//新建offsetManager，为了能够手动控制offset
	offsetManager, err := sarama.NewOffsetManagerFromClient("group111", client)
	if err != nil {
		panic("offsetManager create error")
	}
	defer offsetManager.Close()

	//创建一个第2分区的offsetManager，每个partition都维护了自己的offset
	partitionOffsetManager, err := offsetManager.ManagePartition("0606_test", 2)
	if err != nil {
		panic("partitionOffsetManager create error")
	}
	defer partitionOffsetManager.Close()

	fmt.Println("consumer init success")

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	//sarama提供了一些额外的方法，以便我们获取broker那边的情况
	//topics, _ := consumer.Topics()
	//fmt.Println(topics)
	//partitions, _ := consumer.Partitions("0606_test")
	//fmt.Println(partitions)

	//第一次的offset从kafka获取(发送OffsetFetchRequest)，之后从本地获取，由MarkOffset()得来
	nextOffset, _ := partitionOffsetManager.NextOffset()
	fmt.Println(nextOffset)

	//创建一个分区consumer，从上次提交的offset开始进行消费
	partitionConsumer, err := consumer.ConsumePartition("0606_test", 2, nextOffset+1)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	fmt.Println("start consume really")

ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d\n message:%s", msg.Offset, string(msg.Value))
			//拿到下一个offset
			nextOffset, offsetString := partitionOffsetManager.NextOffset()
			fmt.Println(nextOffset+1, "...", offsetString)
			//提交offset，默认提交到本地缓存，每秒钟往broker提交一次（可以设置）
			partitionOffsetManager.MarkOffset(nextOffset+1, "modified metadata")

		case <-signals:
			break ConsumerLoop
		}
	}
}
