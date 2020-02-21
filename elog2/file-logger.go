package elog2

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
)

var _ Interface = &FileLogger{}

type FileLoggerOption struct {
	Level int
	RecordChanSize int
	LogPath string
	LogName string
}

func (opt *FileLoggerOption) Check() error {
	// 传入的文件名不能为空
	if opt.LogPath == "" || opt.LogName == "" {
		return errors.New("logPath or logName is empty")
	}

	if opt.Level < DEBUG || opt.Level > FATAL {
		opt.Level = DEFAULT_LOG_LEVEL
	}

	if opt.RecordChanSize <= 0 {
		opt.RecordChanSize = DEFAULT_RECORD_CHAN_SIZE
	}

	return nil
}

type FileLogger struct {
	level int
	logPath string
	logName string
	file *os.File
	errorFile *os.File

	// 异步写日志通道
	// 这里起的作用是消息缓冲队列，带缓冲chan是原生的队列
	recordChan chan *Record
}

func NewFileLogger(op Option) (*FileLogger, error) {
	var opt *FileLoggerOption
	opt, ok := op.(*FileLoggerOption)
	if !ok {
		return nil, errors.New("wrong option type")
	}
	if err := opt.Check(); err != nil {
		return nil, err
	}

	logger := &FileLogger{
		level: opt.Level,
		logPath:opt.LogPath,
		logName:opt.LogName,
		recordChan:make(chan *Record),
	}

	logger.init()

	return logger, nil
}

func (l *FileLogger) init() {
	filename := fmt.Sprintf("%s/%s.log", l.logPath, l.logName)
	l.file = openFile0755(filename)

	warnfilename := fmt.Sprintf("%s/%s.log.error", l.logPath, l.logName)
	l.errorFile = openFile0755(warnfilename)

	// 另起协程，后台写日志
	go l.writeLogBackground()
}


func (l *FileLogger) SetLevel(level int) {
	if level < DEBUG || level > FATAL {
		level = DEFAULT_LOG_LEVEL
	}
	l.level = level
}

func (l *FileLogger) log(logLevel int, format string, args ...interface{}) {
	if l.level > logLevel {
		return
	}
	// 文件日志，把日志异步提交到另一个go程去写
	rec := record(logLevel, format, args...) // 记得要跟...，否则会有些小bug
	l.recordChan <- rec                      // 丢到缓冲队列去
}

func (l *FileLogger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *FileLogger) Trace(format string, args ...interface{}) {
	l.log(TRACE, format, args...)
}

func (l *FileLogger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *FileLogger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

func (l *FileLogger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *FileLogger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

func (l *FileLogger) Close() {
	_ = l.file.Close()
	_ = l.errorFile.Close()
	close(l.recordChan)
}

// 后台写日志
func (l *FileLogger) writeLogBackground() {
	// 不停从channel取record
	for rec := range l.recordChan {
		if rec.Level <= WARN {
			_, _ = fmt.Fprint(l.file, rec.String())
		} else {
			_, _ = fmt.Fprint(l.errorFile, rec.String())
		}
	}
}


