package logging

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"reflect"
)

const (
	PanicLevel = logrus.PanicLevel
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel = logrus.FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel = logrus.ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel = logrus.WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel = logrus.InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel = logrus.DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel = logrus.TraceLevel
)

var Logger = logrus.New()

// 设置格式化的格式
func SetFormat(format string) {
	Logger.SetFormatter(&TextFormatter{
		LogFormat: format,
	})
}

// 设置输出的位置
func SetOutput(writer interface{}) {
	if reflect.TypeOf(writer).String() == "string" {
		if filePath, ok := writer.(string); ok {
			file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err == nil {
				Logger.SetOutput(file)
			} else {
				Logger.Error("无法打开文件, 请检查文件是否存在和权限是否正确")
			}
		}
	} else if w, ok := writer.(io.Writer); ok {
		Logger.SetOutput(w)
	}
}

// 设置日志输出级别
func SetLevel(level interface{}) {
	if l, ok := level.(logrus.Level); ok {
		Logger.SetLevel(l)
	} else if l, ok := level.(int); ok {
		Logger.SetLevel(logrus.Level(l))
	}
}
