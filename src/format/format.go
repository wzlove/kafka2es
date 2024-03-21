package format

import "whoops/kafka2es/src/model"

const (
	formatDemo = "demo"
)

type FuncFormat func(data string) string

//CreateFormats 创建格式化方法
func CreateFormats(handlerConf *model.Handler) []FuncFormat {
	formats := make([]FuncFormat, 0)
	formats = append(formats, baseFormat())
	for _, f := range handlerConf.FormatConf {
		switch f.Action {
		case formatDemo:
			formats = append(formats, demoFormat())
		default:
			return nil
		}
	}
	return formats
}
