package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *zap.Logger
)

func Init() {
	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "@timestamp"
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder

	enc := zapcore.NewJSONEncoder(cfg)

	globalLogger = zap.New(
		zapcore.NewCore(
			enc,
			os.Stdout,
			zapcore.DebugLevel,
		),
		zap.AddCaller(),
	)
}

func Sync() {
	globalLogger.Sync()
}

func Named(name string) *zap.Logger {
	return globalLogger.Named(name)
}
