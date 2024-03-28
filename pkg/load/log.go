package load

import (
	"cloud-platform/global"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

func InitLog(logPath string) {
	hook := lumberjack.Logger{
		Filename: logPath,
	}
	w := zapcore.AddSync(&hook)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		zap.InfoLevel,
	)

	logger = zap.New(core, zap.AddCaller())
	global.Logger = logger.Sugar()
}

func Sync() {
	logger.Sync()
}
