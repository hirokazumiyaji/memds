package memds

import (
	"time"

	"github.com/uber-go/zap"
)

var (
	defaultLogger zap.Logger
)

func init() {
	enc := zap.NewJSONEncoder(
		zap.TimeFormatter(func(t time.Time) zap.Field {
			return zap.String("time", t.Local().Format(time.RFC3339))
		}),
	)
	defaultLogger = zap.New(
		enc,
	)
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
