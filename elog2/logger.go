package elog2

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
)

var _ Interface = &Logger{}

type LoggerOption struct {
	Level int
	RecordChanSize int
	WriteTo io.WriteCloser
}

func (opt *LoggerOption) Check() error {
	// 传入的接口不能为空
	if opt.WriteTo == nil {
		return errors.New("writeTo is empty")
	}

	if opt.Level < DEBUG || opt.Level > FATAL {
		opt.Level = DEFAULT_LOG_LEVEL
	}

	if opt.RecordChanSize <= 0 {
		opt.RecordChanSize = DEFAULT_RECORD_CHAN_SIZE
	}

	return nil
}

// 通用的日志记录器
type Logger struct {
	level int
	writeTo io.WriteCloser
	recordChan chan *Record
}

func NewLogger(op Option) (*Logger, error) {
	var opt *LoggerOption
	opt, ok := op.(*LoggerOption)
	if !ok {
		return nil, errors.New("wrong option type")
	}
	if err := opt.Check(); err != nil {
		return nil, err
	}

	logger := &Logger{
		level: opt.Level,
		writeTo:opt.WriteTo,

	}

	// 另起协程，后台写日志
	go logger.writeLogBackground()

	return logger, nil
}

func (l *Logger) SetLevel(level int) {
	if level < DEBUG || level > FATAL {
		level = DEFAULT_LOG_LEVEL
	}
	l.level = level
}

func (l *Logger) log(level int, format string, args ...interface{}) {
	// l.level是日志器设置的最高日志级别，level则是具体调用的log级别
	// 级别数字越低越高，因此只有当 l.level <= level 才允许记录
	if l.level > level {
		return
	}
	rec := record(level, format, args...)
	l.recordChan <- rec
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Trace(format string, args ...interface{}) {
	l.log(TRACE, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

func (l *Logger) Close() {
	if l.writeTo != os.Stdout && l.writeTo != os.Stdin && l.writeTo != os.Stderr {
		_ = l.writeTo.Close()
	}
	close(l.recordChan)
}

// 后台写日志
func (l *Logger) writeLogBackground() {
	// 不停从channel取record
	for rec := range l.recordChan {
		_, _ = fmt.Fprint(l.writeTo, rec.String())
	}
}
