package es_proxy

import (
	"context"
	"errors"
	"github.com/olivere/elastic/v7"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
	"whoops/kafka2es/src/model"
)

// EsService ES服务
type EsService struct {
	esClient       *elastic.Client    //es客户端
	esDataObjPool  sync.Pool          //esData的对象池
	esDataTaskPool *ants.PoolWithFunc //异步处理es的数据的协程池
	workNum        int32              //并发执行的work数量
	bulkActionNum  int32              //flush的批量条数
}

var (
	//EsDataChannel Es数据的Data通道,默认无缓冲
	EsDataChannel = make(chan *EsData)
)

// Init 初始化es服务
func Init(handlerConf *model.Handler) {
	//初始化es数据的chanel缓冲区大小
	if bulkActions := handlerConf.ElasticConf.BulkActions; bulkActions > 0 {
		EsDataChannel = make(chan *EsData, bulkActions)
	}
	//初始化es客户端
	esClient, err := initEsClient(handlerConf)
	if err != nil {
		logrus.Panicf("initEsClient err %s", err.Error())
	}
	esService := &EsService{
		esClient: esClient,
		esDataObjPool: sync.Pool{
			New: func() interface{} {
				return make([]*EsData, 0, handlerConf.ElasticConf.BulkActions)
			},
		},
		workNum:       handlerConf.ElasticConf.Workers,
		bulkActionNum: handlerConf.ElasticConf.BulkActions,
	}
	esService.initEsDataTaskPool(esClient)
	//开启协程去异步将数据放入到ES通道中
	go esService.handleDataToEs()
}

// initEsClient 初始化es客户端
func initEsClient(handlerConf *model.Handler) (*elastic.Client, error) {
	if len(handlerConf.ElasticConf.Hosts) == 0 {
		return nil, errors.New("es hosts is empty")
	}

	//create a  client and connect to the Elasticsearch
	opts := make([]elastic.ClientOptionFunc, 0)

	opts = []elastic.ClientOptionFunc{
		elastic.SetURL(handlerConf.ElasticConf.Hosts...),
		elastic.SetGzip(true),
		elastic.SetErrorLog(logrus.New()),
	}
	//若为单机情况
	if len(handlerConf.ElasticConf.Hosts) == 1 {
		opts = append(opts, elastic.SetSniff(false))
	}

	client, err := elastic.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	// Ping the Elasticsearch server
	//info, code, err := client.Ping("http://127.0.0.1:9200").Do(es_proxy.ctx)
	//create a context to execute each service
	ctx := context.Background()
	for _, addr := range handlerConf.ElasticConf.Hosts {
		_, _, pingErr := client.Ping(addr).Do(ctx)
		if pingErr != nil {
			// Handle error
			return nil, pingErr
		}
		// Getting the ES version number
		esVersion, versionErr := client.ElasticsearchVersion(addr)
		if versionErr != nil {
			// Handle error
			return nil, pingErr
		}
		logrus.Infof("ElasticSearch addr:%s version %s", addr, esVersion)
	}
	return client, nil
}

/**
 * @Description: 初始化写入Es数据协程池
 */
func (es *EsService) initEsDataTaskPool(esClient *elastic.Client) {
	logrus.Infof("init esData task pool success")
	esDataTaskPool, _ := ants.NewPoolWithFunc(int(es.workNum), func(esDataListObj interface{}) {
		esDataList := esDataListObj.([]*EsData)
		err := es.bulkOperate(esClient, esDataList)
		if err != nil {
			logrus.Errorf("bulk insert esDataList err,err is %s", err.Error())
		}
	})
	es.esDataTaskPool = esDataTaskPool
}

// 批量操作数据到es
func (es *EsService) bulkOperate(esClient *elastic.Client, list []*EsData) error {
	if len(list) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bulkReqList := make([]elastic.BulkableRequest, 0, len(list))
	for _, val := range list {
		switch val.OpType {
		case OpTypeIndex:
			bulkReqTmp := elastic.NewBulkIndexRequest().Index(val.Index).Id(val.Id).Doc(val.Body)
			bulkReqList = append(bulkReqList, bulkReqTmp)
		}
	}
	bulkProcessorService := esClient.BulkProcessor()
	if es.workNum > 0 {
		bulkProcessorService.Workers(int(es.workNum))
	}
	if es.bulkActionNum > 0 {
		bulkProcessorService.BulkActions(int(es.bulkActionNum))
	}
	p, err := bulkProcessorService.Do(ctx)
	if err != nil {
		return err
	}
	for _, bulkReq := range bulkReqList {
		p.Add(bulkReq)
	}
	return p.Close()
}

/**
 * @Description: 写数据到ES中
 */
func (es *EsService) handleDataToEs() {
	//从对象池中拿对象
	esDatalist := es.esDataObjPool.Get().([]*EsData)
	for esData := range EsDataChannel {
		esDatalist = append(esDatalist, esData)
		if len(esDatalist) >= int(es.bulkActionNum) {
			if err := es.esDataTaskPool.Invoke(esDatalist); err != nil {
				logrus.Errorf("write es data err,%s", err.Error())
			}
			//清空数据并放会对象池
			esDatalist = esDatalist[:0]
			es.esDataObjPool.Put(esDatalist)
			esDatalist = es.esDataObjPool.Get().([]*EsData)
		}
	}
}
