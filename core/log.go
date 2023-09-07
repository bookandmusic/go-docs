package core

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/utils"
)

type ColoredFormatter struct {
	SimpleFormatter logrus.Formatter
}

func (f *ColoredFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor *color.Color

	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = color.New(color.FgCyan)
	case logrus.InfoLevel:
		levelColor = color.New(color.FgGreen)
	case logrus.WarnLevel:
		levelColor = color.New(color.FgYellow)
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = color.New(color.FgRed)
	default:
		levelColor = color.New(color.FgWhite)
	}

	levelColor.Set()
	defer color.Unset()

	return f.SimpleFormatter.Format(entry)
}

func InitLog() *logrus.Logger {
	// 创建一个 logrus 的日志记录器
	logger := logrus.New()

	// 解析日志级别
	parsedLogLevel, err := logrus.ParseLevel(global.GVA_CONFIG.Log.Level)
	if err != nil {
		// 如果指定的日志级别无效，使用默认级别：info
		logger.Warn("error log level, set default level: info")
		parsedLogLevel = logrus.InfoLevel
	}

	// 设置日志输出格式为JSON格式
	logger.SetFormatter(&logrus.JSONFormatter{})
	// 设置日志级别
	logger.SetLevel(parsedLogLevel)

	if err := utils.PathExists(global.GVA_CONFIG.Log.LogFile); err != nil {
		panic(fmt.Errorf("create log file error: %v", err))
	}

	// 创建一个 lumberjack.Logger，设置日志文件路径和分割条件
	logFile := &lumberjack.Logger{
		Filename:   global.GVA_CONFIG.Log.LogFile, // 日志文件名
		MaxSize:    5,                             // 日志文件大小上限（以MB为单位）
		MaxBackups: 3,                             // 最大保留的旧日志文件数量
		MaxAge:     28,                            // 旧日志文件最大保存天数
		LocalTime:  true,                          // 使用本地时间命名旧日志文件
		Compress:   false,                         // 是否压缩旧日志文件
	}

	// 设置日志输出到 logrus 和文件
	// 将日志同时输出到控制台和文件
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	logger.SetOutput(multiWriter)

	// 设置日志的日期和时间格式
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// 替换logrus的默认Formatter
	logger.SetFormatter(&ColoredFormatter{
		SimpleFormatter: &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		},
	})

	return logger
}
