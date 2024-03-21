package kafka_proxy

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"whoops/kafka2es/src/model"
)

type ConsumerClient struct {
	//消息处理器
	handler *MessageHandler
	topics  []string
	client  sarama.ConsumerGroup
	//就绪通道，判断是否退出
	ready chan bool
}

//InitConsumerMessageQueue 初始化消费者消息队列
func InitConsumerMessageQueue(conf *model.KafkaConf) chan string {
	messageQueue := withMessageQueue(conf.QueueSize)
	consumerClient := &ConsumerClient{
		handler: &MessageHandler{
			MessageQueue: messageQueue,
		},
		topics: conf.Topics,
		client: newConsumerGroup(conf),
	}
	//开启协程去消费消息
	go consumerClient.StartConsume()
	return messageQueue
}

//WithMessageQueue 获取缓存消息队列
func withMessageQueue(size int) chan string {
	if size <= 0 {
		return make(chan string)
	}
	return make(chan string, size)
}

//新建消费者组
func newConsumerGroup(conf *model.KafkaConf) sarama.ConsumerGroup {
	kafkaConf := sarama.NewConfig()
	kafkaConf.Consumer.Return.Errors = true
	kafkaConf.Consumer.Offsets.Initial = sarama.OffsetNewest
	//一次请求抓取的最小bytes
	if conf.MinBytes != 0 {
		kafkaConf.Consumer.Fetch.Min = conf.MinBytes
		kafkaConf.Consumer.Fetch.Max = conf.MinBytes << 3
	}
	consumerGroup, err := sarama.NewConsumerGroup(conf.Brokers, conf.GroupID, kafkaConf)
	if err != nil {
		logrus.Panicf("new consumer group err,%s", err.Error())
		return nil
	} else {
		return consumerGroup
	}
}

//StartConsume 消息消息
func (cc *ConsumerClient) StartConsume() {
	go cc.captureErrs()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := cc.client.Consume(ctx, cc.topics, cc.handler); err != nil {
				logrus.Panicf("Error from consumer: %s", err.Error())
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			cc.ready = make(chan bool)
		}
	}()

	<-cc.ready // Await till the consumer has been set up
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		logrus.Errorf("terminating: context cancelled")
	case <-sigterm:
		logrus.Errorf("terminating: via signal")
	}

	wg.Wait()
	if err := cc.client.Close(); err != nil {
		logrus.Panicf("Error closing client: %v", err)
	}
}

//捕捉错误消息,并打印日志
func (cc *ConsumerClient) captureErrs() {
	// Track errors
	for err := range cc.client.Errors() {
		logrus.Errorf("%v Inverted Listener error, %s", cc.topics, err.Error())
	}
}
