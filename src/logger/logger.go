package logger

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
	"whoops/kafka2es/src/model"
)

var (
	logFormatter = nested.Formatter{
		NoColors:        true,
		NoFieldsColors:  true,
		TimestampFormat: "2006-01-02 15:04:05.000|",
		HideKeys:        true,
		CallerFirst:     true,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			s := strings.Split(frame.Function, ".")
			funcName := s[len(s)-1]
			return fmt.Sprintf(" [%s:%d][%s()]", path.Base(frame.File), frame.Line, funcName)
		},
	}
)

//Init 读取配置文件配置logger
func Init(config *model.LogConf) {
	//运行日志
	logrus.SetOutput(os.Stdout)
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		logrus.Warnf("setupLogger parse log level err %s, use default level: debug", err.Error())
		level = logrus.DebugLevel
	}
	levels := logrus.AllLevels[:level+1]
	logrus.SetFormatter(&logFormatter)
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)
	lfsLogHook := newLfsHook(config.LogPath, levels, &logFormatter, 2*24*time.Hour)
	logrus.AddHook(lfsLogHook)
}
