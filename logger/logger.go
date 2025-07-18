package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func LoggerInit() {
	Log = logrus.New()

	// 控制台输出配置（文本格式）
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 文件输出配置（JSON格式+按天分割）
	writer, _ := rotatelogs.New(
		filepath.Join("data/logs", "app_%Y%m%d.json"),
		rotatelogs.WithLinkName("data/logs/app.json"),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(7*24*time.Hour),
	)

	fileFormatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyMsg:   "@message",
		},
	}

	Log.AddHook(&fileHook{
		Writer:    writer,
		Formatter: fileFormatter,
	})
}

type fileHook struct {
	Writer    *rotatelogs.RotateLogs
	Formatter logrus.Formatter
}

func (h *fileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *fileHook) Fire(entry *logrus.Entry) error {
	msg, _ := h.Formatter.Format(entry)
	_, err := h.Writer.Write(msg)
	return err
}
