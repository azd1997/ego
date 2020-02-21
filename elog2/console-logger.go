package elog2

import (
	"os"
)

var _ Interface = &ConsoleLogger{}


type ConsoleLogger struct {
	level int
}

func NewConsoleLogger(level int) *ConsoleLogger {
	if level < DEBUG || level > FATAL {
		level = DEFAULT_LOG_LEVEL
	}

	logger := &ConsoleLogger{
		level: level,
	}

	return logger
}


func (f *ConsoleLogger) SetLevel(level int) {
	if level < DEBUG || level > FATAL {
		level = DEFAULT_LOG_LEVEL
	}
	f.level = level
}

func (f *ConsoleLogger) log(logLevel int, format string, args ...interface{}) {
	if f.level > logLevel {
		return
	}
	writeLog(os.Stdout, logLevel, format, args)
}

func (f *ConsoleLogger) Debug(format string, args ...interface{}) {
	f.log(DEBUG, format, args...)
}

func (f *ConsoleLogger) Trace(format string, args ...interface{}) {
	f.log(TRACE, format, args...)
}

func (f *ConsoleLogger) Info(format string, args ...interface{}) {
	f.log(INFO, format, args...)
}

func (f *ConsoleLogger) Warn(format string, args ...interface{}) {
	f.log(WARN, format, args...)
}

func (f *ConsoleLogger) Error(format string, args ...interface{}) {
	f.log(ERROR, format, args...)
}

func (f *ConsoleLogger) Fatal(format string, args ...interface{}) {
	f.log(FATAL, format, args...)
}

func (f *ConsoleLogger) Close() {}

