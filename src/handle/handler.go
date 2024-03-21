package handle

import (
	"whoops/kafka2es/src/es_proxy"
	"whoops/kafka2es/src/format"
	"whoops/kafka2es/src/kafka_proxy"
	"whoops/kafka2es/src/model"
)

func Init(handlerConf *model.Handler) {
	//a.初始化kafka消息队列
	messageQueue := kafka_proxy.InitConsumerMessageQueue(handlerConf.KafkaConf)
	//b.初始化handler处理器
	initHandler(messageQueue,
		WithEsIndex(handlerConf.ElasticConf.Index),
		WithFormatFun(format.CreateFormats(handlerConf)))
	//c.初始化es输出组件
	es_proxy.Init(handlerConf)
}
