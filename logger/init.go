package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
)

var logger = logrus.New()

// 全局的 logger 配置

func init() {
	logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)        // 输出到标准输出
	logger.SetLevel(logrus.DebugLevel) // level ： debug level
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			function = ""
			file = fmt.Sprintf(" %s:%d ", frame.File, frame.Line)
			return
		}},
	)
}

// 返回全局的 logger 句柄
func Lg() *logrus.Logger {
	return logger
}
