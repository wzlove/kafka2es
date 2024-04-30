package model

type (
	//GlobalConfig 配置信息
	GlobalConfig struct {
		Log     *LogConf   //日志相关配置
		Handler []*Handler //handler相关配置
		Port    int32      //http端口号
	}

	Handler struct {
		KafkaConf   *KafkaConf    //输入流 kafka相关配置
		ElasticConf *ElasticConf  //输出流 es相关配置
		FormatConf  []*FormatConf //中间流 格式化相关配置
	}

	FormatConf struct {
		Action string `json:",options=demo"` //格式化的行为
	}

	KafkaConf struct {
		GroupID   string   //消费者的组id
		Brokers   []string //kafka的 brokers
		Topics    []string //消费的topic
		Consumers int32    //消费者的数量
		QueueSize int32    //消息队列大小
		MinBytes  int32    //一次请求抓取的最小bytes
	}

	ElasticConf struct {
		Index       string   //ES写入的index索引名
		Hosts       []string //ES的地址
		BulkActions int32    //批量写ES的个数
		Workers     int32    //写ES的工作数量
	}

	//LogConf 日志相关配置
	LogConf struct {
		Level   string
		LogPath string
	}
)
