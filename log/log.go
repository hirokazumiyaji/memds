package log

import (
	"time"

	"github.com/uber-go/zap"
)

var (
	defaultLogger zap.Logger
)

func init() {
	defaultLogger = New("info")
}

func New(level string) zap.Logger {
	enc := zap.NewJSONEncoder(
		zap.TimeFormatter(func(t time.Time) zap.Field {
			return zap.String("time", t.Local().Format(time.RFC3339))
		}),
	)
	l := zap.InfoLevel
	l.Set(level)
	logger := zap.New(
		enc,
		l,
	)
	return logger
}

func SetLevel(level string) {
	defaultLogger = New(level)
}

func Debug(msg string, fields ...zap.Field) {
	defaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	defaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	defaultLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	defaultLogger.Error(msg, fields...)
}
