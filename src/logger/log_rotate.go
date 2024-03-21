package logger

import (
	"bytes"
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"path/filepath"
	"time"
)

func newLfsHook(logPath string, levels []logrus.Level, formatter logrus.Formatter, maxAge time.Duration) logrus.Hook {
	writer, err := rotatelogs.New(
		logPath+"-%Y%m%d%H.log",
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithMaxAge(maxAge),           // 文件最大保存时间
		rotatelogs.WithLinkName(logPath+".log"), // 生成软链，指向最新日志文件
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	lfsHook := lfshook.NewHook(bindWriter(writer, levels), formatter)
	return lfsHook
}

func bindWriter(writer io.Writer, levels []logrus.Level) lfshook.WriterMap {
	writerMap := lfshook.WriterMap{}
	for i := 0; i < len(levels); i++ {
		writerMap[levels[i]] = writer
	}
	return writerMap
}

//MyFormatter 自定义的format
type MyFormatter struct{}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006/01/02 15:04:05")
	var newLog string

	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("[%s] [%s] <%s:%d> %s\n",
			//timestamp, entry.Level, fName, entry.Caller.Line, entry.Caller.Function, entry.Message)
			timestamp, entry.Level, fName, entry.Caller.Line, entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}
