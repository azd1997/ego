package elog2

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"time"
)

var _ Interface = &Logger{}

// 作为日志输出源，需要实现OutputSource
type OutputSource interface {
	io.WriteCloser
	Size() int64	// 返回已使用字节数
	Name() string	// 标识符
	Rename(string)	// 设置标识符
	BackupName() string	// 生成备份文件的名字
}

// // 外部提供一个根据名字转换成新的日志输出源的方法
type NewOutputSource func(name string) OutputSource

type LoggerOption struct {
	Level int
	RecordChanSize int

	WriteTo OutputSource

	// 日志切分参数
	SplitType string
	SplitTimeIntervalOrFileSize int64
	// 外部提供一个根据名字转换成新的日志输出源的方法
	NewOutputSource NewOutputSource
}

func (opt *LoggerOption) Check() error {
	// 传入的接口不能为空
	if opt.WriteTo == nil {
		return errors.New("WriteTo can't be empty!")
	}

	if opt.Level < DEBUG || opt.Level > FATAL {
		opt.Level = DEFAULT_LOG_LEVEL
	}

	if opt.RecordChanSize <= 0 {
		opt.RecordChanSize = DEFAULT_RECORD_CHAN_SIZE
	}

	// 检查日志切分参数
	switch opt.SplitType {
	case SPLIT_TIME:
		if opt.SplitTimeIntervalOrFileSize <= 0 {
			opt.SplitTimeIntervalOrFileSize = SPLIT_TIME_HOUR
		}
		// 检查NewOutputSource
		if opt.NewOutputSource == nil {
			return errors.New("NewOutputSource can't be empty!")
		}
	case SPLIT_SIZE:
		if opt.SplitTimeIntervalOrFileSize <= 0 {
			opt.SplitTimeIntervalOrFileSize = SPLIT_SIZE_100M
		}
		// 检查BackupIOFunc
		if opt.NewOutputSource == nil {
			return errors.New("NewOutputSource can't be empty!")
		}
	case SPLIT_NONE:
		opt.SplitTimeIntervalOrFileSize = 0
		opt.NewOutputSource = nil
	default:
		opt.SplitType = SPLIT_TIME
		opt.SplitTimeIntervalOrFileSize = SPLIT_TIME_HOUR
	}

	return nil
}

// 通用的日志记录器
type Logger struct {
	level int
	writeTo OutputSource
	recordChan chan *Record

	// 日志切分用
	splitType string
	splitTimeIntervalOrFileSize int64
	lastSplitTime int64
	newOutputSource NewOutputSource
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
		recordChan:make(chan *Record, DEFAULT_RECORD_CHAN_SIZE),
		splitType:opt.SplitType,
		splitTimeIntervalOrFileSize:opt.SplitTimeIntervalOrFileSize,
		lastSplitTime:time.Now().Unix(),
		newOutputSource:opt.NewOutputSource,
	}

	// 另起协程，后台写日志
	go logger.writeLogBackground()

	return logger, nil
}

// 在每次写日志之前，先检查日志文件是否达到切分条件
// 将isErrorFile传入而不是直接在写日志时同时检查两种日志文件，这样效率更好
func (l *Logger) checkLogAndSplit() {
	switch l.splitType {
	case SPLIT_TIME:
		l.checkLogAndSplitByTime()
	case SPLIT_SIZE:
		l.checkLogAndSplitBySize()
	case SPLIT_NONE:
		// DO NOTHING
	}
}

func (l *Logger) checkLogAndSplitByTime() {
	// 切分日志，生成备份文件
	if time.Now().Unix() - l.lastSplitTime >= l.splitTimeIntervalOrFileSize {
		l.backupLog()
	}
}

func (l *Logger) checkLogAndSplitBySize() {
	// 比较Byte数是否达到切分要求
	if l.writeTo.Size() >= l.splitTimeIntervalOrFileSize {
		l.backupLog()
	}
}

func (l *Logger) backupLog() {
	backupLogName := l.writeTo.BackupName()
	curLogName := l.writeTo.Name()
	l.writeTo.Rename(backupLogName)
	_ = l.writeTo.Close()
	l.writeTo = l.newOutputSource(curLogName)
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
	_ = l.writeTo.Close()
	close(l.recordChan)
}

// 后台写日志
func (l *Logger) writeLogBackground() {
	// 不停从channel取record
	for rec := range l.recordChan {
		// 先检查Log是否需要切分
		l.checkLogAndSplit()
		_, _ = fmt.Fprint(l.writeTo, rec.String())
	}
}
