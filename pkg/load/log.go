package load

import (
	"cloud-platform/pkg/load/tlog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createLumberWriter(filePath string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    20,   // 一个文件最大可达 20M。
		MaxBackups: 5,    // 最多同时保存 5 个文件。
		MaxAge:     10,   // 一个文件最多可以保存 10 天。
		Compress:   true, // 用 gzip 压缩。
		// LocalTime:  true,
	}
}

func checkPathExist(logPath string) (string, string, string) {
	os.MkdirAll(logPath, os.ModePerm)
	logFileName := time.Now().Format("2006-01-02") + ".log"
	warnFileName := time.Now().Format("2006-01-02") + "-warn.log"
	errorFileName := time.Now().Format("2006-01-02") + "-error.log"
	file := path.Join(logPath, logFileName)
	warnfile := path.Join(logPath, warnFileName)
	errfile := path.Join(logPath, errorFileName)
	checkPath(file)
	checkPath(warnfile)
	checkPath(errfile)
	return file, warnfile, errfile
}

func checkPath(file string) {
	if _, err := os.Stat(file); err != nil {
		if _, err := os.Create(file); err != nil {
			panic(err.Error())
		}
	}
}

func InitLog(logPath string) {
	infofile, warnfile, errfile := checkPathExist(logPath)
	zapOptions := []zap.Option{zap.AddCallerSkip(1),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel)}
	zapLogger := zap.New(
		buildZapCores(infofile, warnfile, errfile),
		zapOptions...)
	tlog.InitLogger(zapLogger)
	// logger := hertzzap.NewLogger(hertzzap.WithZapOptions(zapOptions...))
	// logger.SetOutput(createLumberWriter(file))
	// logger.SetLevel(tlog.LevelInfo)
	// tlog.SetLogger(logger)
}

func buildZapCores(infoFile, warnFile, errFile string) zapcore.Core {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:  "msg",                       //结构化（json）输出：msg的key
		LevelKey:    "level",                     //结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:     "ts",                        //结构化（json）输出：时间的key（INFO，WARN，ERROR等）
		CallerKey:   "file",                      //结构化（json）输出：打印日志的文件对应的Key
		EncodeLevel: zapcore.CapitalLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			pwd, _ := os.Getwd()
			path, _ := filepath.Rel(pwd, caller.File)
			if strings.HasPrefix(path, "../") {
				path = caller.TrimmedPath()
			}
			enc.AppendString(path)
		},
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05")) // 时间格式
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl >= zap.InfoLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	infoWriter := createLumberWriter(infoFile)
	warnWriter := createLumberWriter(warnFile)
	errWriter := createLumberWriter(errFile)

	return zapcore.NewTee(
		// 将info及以下写入logPath，NewConsoleEncoder 是非结构化输出
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(warnWriter), warnLevel),
		// warn及以上写入errPath
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(errWriter), errorLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(infoWriter), zapcore.AddSync(os.Stdout)), zap.InfoLevel),
	)
}
