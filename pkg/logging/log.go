package logging

import (
	"github.com/gin-gonic/gin"
	"go-gin-demo/pkg/setting"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

//Log配置结构体
type Log struct {
	Gin         string
	App         string
	Http        string
	ServiceName string
}

var LogSetting = &Log{}
var AppLogger *zap.Logger
var HTTPLogger *zap.Logger

//定制日志
func Setup() {
	setting.MapTo("log", LogSetting)
	//记录Gin日志
	f, _ := os.Create(LogSetting.Gin)
	// Use the following code if you need to write the logs to file and console at the same time.
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	//定制日志
	AppLogger = NewLogger(LogSetting.App, zapcore.InfoLevel, 128, 30, 7, true, LogSetting.ServiceName)
	HTTPLogger = NewLogger(LogSetting.Http, zapcore.InfoLevel, 128, 30, 7, true, LogSetting.ServiceName)
}
