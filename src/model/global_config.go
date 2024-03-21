package model

type (
	//GlobalConfig 配置信息
	GlobalConfig struct {
		Handler []*Handler //handler相关配置
		Log     *LogConf   //日志相关配置
		Port    int        //http端口号
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
		Brokers   []string //kafka的 brokers
		Topics    []string //消费的topic
		GroupID   string   //消费者的组id
		MinBytes  int32    //一次请求抓取的最小bytes
		Consumers int      //消费者的数量
		QueueSize int      //消息队列大小
	}

	ElasticConf struct {
		Hosts       []string //ES的地址
		BulkActions int      //批量写ES的个数
		Index       string   //ES写入的index索引名
		Workers     int      //写ES的工作数量
	}

	//LogConf 日志相关配置
	LogConf struct {
		Level   string
		LogPath string
	}
)
