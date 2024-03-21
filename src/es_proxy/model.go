package es_proxy

//EsData 存入es队列结构
type EsData struct {
	// op:index, create, update, delete
	OpType string `json:"op_type"`
	Index  string `json:"index"`
	Id     string `json:"id"`
	Body   string `json:"body"`
}

const (
	//OpTypeIndex 如果数据存在，使用create操作失败，会提示文档已经存在，使用index则可以成功执行
	OpTypeIndex = "index"
)
