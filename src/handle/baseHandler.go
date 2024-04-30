package handle

import (
	"github.com/go-basic/uuid"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"runtime"
	"whoops/kafka2es/src/es_proxy"
	"whoops/kafka2es/src/format"
)

const (
	//默认缓冲区大小
	defaultBulkCount    = 2048
	defaultTaskPoolSize = 8
	defaultIndexName    = "test"
)

type (
	//BaseHandler 公共处理器
	BaseHandler struct {
		messageQueue chan string //消息通道
		handlerOption
	}
	HandlerOption func(options *handlerOption)

	handlerOption struct {
		taskPool   *ants.Pool          //异步处理任务的工作池
		esIndex    string              //es的索引
		formatFunc []format.FuncFormat //格式化数据的函数
		bulkCount  int32               //批处理大小
	}
)

func WithBulkCount(bulkCount int32) HandlerOption {
	return func(options *handlerOption) {
		options.bulkCount = bulkCount
	}
}

func WithEsIndex(indexName string) HandlerOption {
	return func(options *handlerOption) {
		if indexName != "" {
			options.esIndex = indexName
		}
	}
}

func WithTaskPool(taskPool *ants.Pool) HandlerOption {
	return func(options *handlerOption) {
		options.taskPool = taskPool
	}
}

func WithFormatFun(formatFunc []format.FuncFormat) HandlerOption {
	return func(options *handlerOption) {
		options.formatFunc = formatFunc
	}
}

// initHandler 初始化handler处理器
func initHandler(messageQueue chan string, opts ...HandlerOption) {
	options := newHandlerOptions()
	for _, opt := range opts {
		opt(&options)
	}
	baseHandler := &BaseHandler{
		messageQueue,
		options,
	}
	//启动handle的处理程序
	go baseHandler.Start()
}

func newHandlerOptions() handlerOption {
	taskPool, _ := ants.NewPool(defaultTaskPoolSize)
	return handlerOption{
		bulkCount: defaultBulkCount,
		taskPool:  taskPool,
		esIndex:   defaultIndexName,
	}
}

// Start 开始执行handler
func (handler *BaseHandler) Start() {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			l := runtime.Stack(buf, false)
			logrus.Errorf("|runForever | baseHandler err,%s:", err)
			logrus.Errorf("|runForever | baseHandler err,%s:", buf[:l])
		}
	}()
	//注意此处因为go1.22之前 短声明形式定义的循环变量从整个循环定义和共享一个
	for body := range handler.messageQueue {
		handler.putTaskBodyToEsChannel(body)
	}
}

// 格式化数据后放到EsChannel中
func (handler *BaseHandler) putTaskBodyToEsChannel(body string) {
	if err := handler.taskPool.Submit(func() {
		esData := &es_proxy.EsData{
			OpType: es_proxy.OpTypeIndex,
			Index:  handler.esIndex,
			Id:     uuid.New(),
		}
		//格式化消息
		for _, formatFunc := range handler.formatFunc {
			body = formatFunc(body)
		}
		esData.Body = body
		//若body不为空才写入
		if body != "" {
			es_proxy.EsDataChannel <- esData
		}
	}); err != nil {
		logrus.Errorf("submit convert message err,%s:", err.Error())
	}
}
