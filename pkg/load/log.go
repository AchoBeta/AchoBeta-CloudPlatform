package load

import (
	"bytes"
	"cloud-platform/pkg/load/log"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func InitLogrus(logPath string) {
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

	infoHook := createHook(file)
	errHook := createHook(errfile)
	witers := []io.Writer{infoHook, errHook, os.Stdout}
	logrus.SetReportCaller(true)
	logrus.SetOutput(io.MultiWriter(witers...))
	logrus.SetFormatter(&customFormatter{})
	logrus.AddHook(&log.TraceIdHook{})
	// logrus.SetLevel(logrus.InfoLevel)
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

type customFormatter struct{}

func (m *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	fmt.Printf("%+v\n", entry)
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string
	if entry.HasCaller() {
		dir, _ := os.Getwd()
		fName, _ := filepath.Rel(dir, entry.Caller.File)
		newLog = fmt.Sprintf("[%s] [%s] [%s:%d] [trace:%s] %s\n",
			timestamp, entry.Level, fName, entry.Caller.Line, entry.Data["trace"], entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}
