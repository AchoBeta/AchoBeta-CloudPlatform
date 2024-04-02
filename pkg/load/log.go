package load

import (
	"os"
	"path"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLog(logPath string) {
	logFileName := time.Now().Format("2006-01-02") + ".log"
	errorFileName := time.Now().Format("2006-01-02") + "-error.log"
	file := path.Join(logPath, logFileName)
	errfile := path.Join(logPath, errorFileName)
	if _, err := os.Stat(file); err != nil {
		if _, err := os.Create(file); err != nil {
			panic(err.Error())
		}
	}

	if _, err := os.Stat(errfile); err != nil {
		if _, err := os.Create(errfile); err != nil {
			panic(err.Error())
		}
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "msg",                       //结构化（json）输出：msg的key
		LevelKey:     "level",                     //结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:      "ts",                        //结构化（json）输出：时间的key（INFO，WARN，ERROR等）
		CallerKey:    "file",                      //结构化（json）输出：打印日志的文件对应的Key
		EncodeLevel:  zapcore.CapitalLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: zapcore.ShortCallerEncoder,  //采用短文件路径编码输出（test/main.go:14	）
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05")) // 时间格式
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}

	logLevel := zap.InfoLevel
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel && lvl >= logLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl >= logLevel
	})

	infoHook := createHook(file)
	errHook := createHook(errfile)

	logger := hertzzap.NewLogger(hertzzap.WithCores(
		hertzzap.CoreConfig{
			//将info及以下写入logPath，NewConsoleEncoder 是非结构化输出
			Enc: zapcore.NewConsoleEncoder(encoderConfig),
			Ws:  zapcore.AddSync(infoHook),
			Lvl: infoLevel,
		}, hertzzap.CoreConfig{
			//warn及以上写入errPath
			Enc: zapcore.NewConsoleEncoder(encoderConfig),
			Ws:  zapcore.AddSync(errHook),
			Lvl: warnLevel,
		}, hertzzap.CoreConfig{
			Enc: zapcore.NewJSONEncoder(encoderConfig),
			Ws:  zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
			Lvl: logLevel,
		}))

	logger.SetOutput(infoHook)
	logger.SetLevel(hlog.LevelDebug)
	hlog.SetLogger(logger)
}

func createHook(filePath string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    20,   // 一个文件最大可达 20M。
		MaxBackups: 5,    // 最多同时保存 5 个文件。
		MaxAge:     10,   // 一个文件最多可以保存 10 天。
		Compress:   true, // 用 gzip 压缩。
		// LocalTime:  true,
	}
}
