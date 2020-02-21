package elog2

import (
	"fmt"
	"os"
)

var _ Interface = &FileLogger{}


type FileLogger struct {
	level int
	logPath string
	logName string
	file *os.File
	errorFile *os.File
}

func NewFileLogger(level int, logPath, logName string) *FileLogger {
	// 传入的文件名不能为空
	if logPath == "" || logName == "" {
		panic("logPath and logName is empty")
	}

	if level < DEBUG || level > FATAL {
		level = DEFAULT_LOG_LEVEL
	}

	logger := &FileLogger{
		level: level,
		logPath:logPath,
		logName:logName,
	}

	logger.init()

	return logger
}

func (f *FileLogger) init() {
	filename := fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	f.file = openFile0755(filename)

	warnfilename := fmt.Sprintf("%s/%s.log.error", f.logPath, f.logName)
	f.errorFile = openFile0755(warnfilename)
}


func (f *FileLogger) SetLevel(level int) {
	if level < DEBUG || level > FATAL {
		level = DEFAULT_LOG_LEVEL
	}
	f.level = level
}

func (f *FileLogger) log(writeto *os.File, logLevel int, format string, args ...interface{}) {
	if f.level > logLevel {
		return
	}
	writeLog(writeto, logLevel, format, args)
}

func (f *FileLogger) Debug(format string, args ...interface{}) {
	f.log(f.file, DEBUG, format, args...)
}

func (f *FileLogger) Trace(format string, args ...interface{}) {
	f.log(f.file, TRACE, format, args...)
}

func (f *FileLogger) Info(format string, args ...interface{}) {
	f.log(f.file, INFO, format, args...)
}

func (f *FileLogger) Warn(format string, args ...interface{}) {
	f.log(f.file, WARN, format, args...)
}

func (f *FileLogger) Error(format string, args ...interface{}) {
	f.log(f.errorFile, ERROR, format, args...)
}

func (f *FileLogger) Fatal(format string, args ...interface{}) {
	f.log(f.errorFile, FATAL, format, args...)
}

func (f *FileLogger) Close() {
	_ = f.file.Close()
	_ = f.errorFile.Close()
}

