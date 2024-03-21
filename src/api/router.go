package api

import (
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"time"
)

//LoadRoute 加载路由信息
func LoadRoute(r *mux.Router) {
	r.HandleFunc("/health", Health)
}

//Health 健康监测服务
func Health(w http.ResponseWriter, r *http.Request) {
	ansMap := map[string]string{
		"currentTime": time.Now().Format("2006-01-02 15:04:05"),
	}
	response, _ := jsoniter.Marshal(ansMap)
	_, _ = w.Write(response)
}
