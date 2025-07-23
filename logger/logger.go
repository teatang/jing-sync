package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// LoggerInit函数用于初始化日志记录器
func LoggerInit() {
	// 创建一个新的日志记录器
	Log = logrus.New()

	// 设置日志记录器的报告调用者信息
	Log.SetReportCaller(true)
	// 设置控制台日志级别为Info级别
	Log.SetLevel(logrus.InfoLevel)
	// 控制台输出配置（文本格式）
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		// 设置时间戳格式
		FullTimestamp:   true,
		// 强制使用颜色
		ForceColors:     true,
		// 设置时间戳格式
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 文件输出配置（JSON格式+按天分割）
	writer, _ := rotatelogs.New(
		// 设置日志文件名格式
		filepath.Join("data/logs", "app_%Y%m%d.json"),
		// 设置软链接名
		rotatelogs.WithLinkName("data/logs/app.json"),
		// 设置日志文件轮转时间
		rotatelogs.WithRotationTime(24*time.Hour),
		// 设置日志文件最大保存时间
		rotatelogs.WithMaxAge(7*24*time.Hour),
	)

	// 创建JSON格式化器
	fileFormatter := &logrus.JSONFormatter{
		// 设置时间戳格式
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		// 设置字段映射
		FieldMap: logrus.FieldMap{
			// 设置时间字段名
			logrus.FieldKeyTime:  "timestamp",
			// 设置日志级别字段名
			logrus.FieldKeyLevel: "level",
			// 设置日志消息字段名
			logrus.FieldKeyMsg:   "message",
		},
	}

	// 添加文件输出钩子
	Log.AddHook(&fileHook{
		// 设置文件写入器
		Writer:    writer,
		// 设置格式化器
		Formatter: fileFormatter,
	})
}

type fileHook struct {
	Writer    *rotatelogs.RotateLogs
	Formatter logrus.Formatter
}

// 返回fileHook的日志级别
func (h *fileHook) Levels() []logrus.Level {
	// 返回Info级别及以上的日志级别
	return []logrus.Level{
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

// Fire方法用于将日志条目写入文件
func (h *fileHook) Fire(entry *logrus.Entry) error {
	// 使用Formatter格式化日志条目
	msg, _ := h.Formatter.Format(entry)
	// 将格式化后的日志条目写入文件
	_, err := h.Writer.Write(msg)
	// 返回写入文件时发生的错误
	return err
}
