package elog2

import (
	"github.com/pkg/errors"
	"os"
)

var _ Interface = &ConsoleLogger{}

type ConsoleLoggerOption struct {
	Level int
}

func (opt *ConsoleLoggerOption) Check() error {
	if opt.Level < DEBUG || opt.Level > FATAL {
		opt.Level = DEFAULT_LOG_LEVEL
	}
	return nil
}


type ConsoleLogger struct {
	level int
}

// 总是返回空错误
func NewConsoleLogger(op Option) (*ConsoleLogger, error) {
	var opt *ConsoleLoggerOption
	opt, ok := op.(*ConsoleLoggerOption)
	if !ok {
		return nil, errors.New("wrong option type")
	}
	if err := opt.Check(); err != nil {
		return nil, err
	}

	logger := &ConsoleLogger{
		level: opt.Level,
	}

	return logger, nil
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

