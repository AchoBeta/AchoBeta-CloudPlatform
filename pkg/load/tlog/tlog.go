package tlog

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logKey string

const loggerKey logKey = "logger"

var logger *zap.Logger

// 给指定的context添加字段
func NewContext(ctx context.Context, fields ...zapcore.Field) context.Context {
	return context.WithValue(ctx, loggerKey, withContext(ctx).With(fields...))
}

func InitLogger(zapLogger *zap.Logger) {
	logger = zapLogger
}

// 从指定的context返回一个zap实例
func withContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return logger
	}
	if ctxLogger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return ctxLogger
	}
	return logger
}

func Infof(format string, v ...interface{}) {
	logger.Info(fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...interface{}) {
	logger.Error(fmt.Sprintf(format, v...))
}

func Warnf(format string, v ...interface{}) {
	logger.Warn(fmt.Sprintf(format, v...))
}

func Debugf(format string, v ...interface{}) {
	logger.Debug(fmt.Sprintf(format, v...))
}

func Panicf(format string, v ...interface{}) {
	logger.Panic(fmt.Sprintf(format, v...))
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatal(fmt.Sprintf(format, v...))
}

// 下面的logger方法会携带trace id

func CtxInfof(ctx context.Context, format string, v ...interface{}) {
	withContext(ctx).Info(fmt.Sprintf(format, v...))
}

func CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	withContext(ctx).Error(fmt.Sprintf(format, v...))
}

func CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	withContext(ctx).Warn(fmt.Sprintf(format, v...))
}

func CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	withContext(ctx).Debug(fmt.Sprintf(format, v...))
}

func CtxPanicf(ctx context.Context, format string, v ...interface{}) {
	withContext(ctx).Panic(fmt.Sprintf(format, v...))
}

func CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	withContext(ctx).Fatal(fmt.Sprintf(format, v...))
}
