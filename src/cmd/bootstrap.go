package cmd

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"whoops/kafka2es/src/api"
	"whoops/kafka2es/src/handle"
	"whoops/kafka2es/src/logger"
	"whoops/kafka2es/src/model"
)

var (
	//配置文件信息
	cfgFileInfo string
)

//InitService 初始化服务
func InitService(globalConf *model.GlobalConfig) {
	initDependent(globalConf)
}

//加载依赖相关配置
func initDependent(globalConf *model.GlobalConfig) {
	//1.初始化日志信息
	initLog(globalConf)
	//2.初始化handler
	initHandler(globalConf)
	//3.初始化http服务
	initHttpServer(globalConf)
}

//初始化日志
func initLog(globalConf *model.GlobalConfig) {
	logger.Init(globalConf.Log)
}

//初始化handler
func initHandler(globalConf *model.GlobalConfig) {
	//初始化handler
	for _, conf := range globalConf.Handler {
		//初始化handle组件
		handle.Init(conf)
	}
}

//initHttpServer 初始化httpServer
func initHttpServer(gConfig *model.GlobalConfig) {
	if port := gConfig.Port; port == 0 {
		panic(errors.New("init http port is empty"))
	} else {
		r := mux.NewRouter()
		api.LoadRoute(r)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", gConfig.Port), r); err != nil {
			logrus.Errorf("init server err,err is %s", err.Error())
		}
	}
}
