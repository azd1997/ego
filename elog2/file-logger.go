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
	warnFile *os.File
}

func NewFileLogger(level int, logPath, logName string) *FileLogger {
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
	f.warnFile = openFile0755(warnfilename)
}

func openFile0755(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed: %s", filename, err))
	}
	return file
}

func (f *FileLogger) SetLevel(level int) {
	if level < DEBUG || level > FATAL {
		level = DEFAULT_LOG_LEVEL
	}
	f.level = level
}

func (f *FileLogger) Debug(format string, args ...interface{}) {
	if f.level > DEBUG {
		return
	}
	_, _ = fmt.Fprintf(f.file, format, args...)
}

func (f *FileLogger) Trace(format string, args ...interface{}) {
	if f.level > DEBUG {
		return
	}
	_, _ = fmt.Fprintf(f.file, format, args...)
}

func (f *FileLogger) Info(format string, args ...interface{}) {
	if f.level > DEBUG {
		return
	}
	_, _ = fmt.Fprintf(f.file, format, args...)
}

func (f *FileLogger) Warn(format string, args ...interface{}) {
	if f.level > DEBUG {
		return
	}
	_, _ = fmt.Fprintf(f.file, format, args...)
}

func (f *FileLogger) Error(format string, args ...interface{}) {
	if f.level > DEBUG {
		return
	}
	_, _ = fmt.Fprintf(f.warnFile, format, args...)
}

func (f *FileLogger) Fatal(format string, args ...interface{}) {
	if f.level > DEBUG {
		return
	}
	_, _ = fmt.Fprintf(f.warnFile, format, args...)
}

func (f *FileLogger) Close() {
	_ = f.file.Close()
	_ = f.warnFile.Close()
}

