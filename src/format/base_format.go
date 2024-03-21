package format

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

//基础format
func baseFormat() FuncFormat {
	return func(val string) string {
		//写入到ES的数据必须为json格式,因此做个校验check
		var m map[string]interface{}
		if err := jsoniter.Unmarshal([]byte(val), &m); err != nil {
			logrus.Errorf("Unmarshal err,%s", err.Error())
			return ""
		}
		return val
	}
}
