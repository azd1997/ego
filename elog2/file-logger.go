package elog2

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"time"
)

var _ Interface = &FileLogger{}

type FileLoggerOption struct {
	Level int
	RecordChanSize int
	LogPath string
	LogName string
	SplitType string
	SplitTimeIntervalOrFileSize int64
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

	// 检查日志切分参数
	switch opt.SplitType {
	case SPLIT_TIME:
		if opt.SplitTimeIntervalOrFileSize <= 0 {
			opt.SplitTimeIntervalOrFileSize = SPLIT_TIME_HOUR
		}
	case SPLIT_SIZE:
		if opt.SplitTimeIntervalOrFileSize <= 0 {
			opt.SplitTimeIntervalOrFileSize = SPLIT_SIZE_100M
		}
	case SPLIT_NONE:
		opt.SplitTimeIntervalOrFileSize = 0
	default:
		opt.SplitType = SPLIT_TIME
		opt.SplitTimeIntervalOrFileSize = SPLIT_TIME_HOUR
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

	// 日志切分方式
	// 日志堆在一起，不方便查阅；另一点是日志文件大小也有限制，需要在一定量以后分文件
	// 因此这里预置了两种切分方式：按时间切分(time)；按大小切分(size)
	// 相应的，有两个参数：时间间隔和切分大小
	splitType string
	splitTimeIntervalOrFileSize int64		// 对于时间而言，单位为s；对于大小而言，单位为Byte

	// 为了实现切分，还必须记录日志的切分时间
	// 这对于按时间切分是必须的
	lastSplitTime int64
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
		recordChan:make(chan *Record, DEFAULT_RECORD_CHAN_SIZE),
		splitType:opt.SplitType,
		splitTimeIntervalOrFileSize:opt.SplitTimeIntervalOrFileSize,
		lastSplitTime:time.Now().Unix(),
	}

	logger.init()

	fmt.Printf("logger splittype = %s, splitinterval = %d\n", logger.splitType, logger.splitTimeIntervalOrFileSize)

	return logger, nil
}

func (l *FileLogger) init() {
	filename := fmt.Sprintf("%s/%s.log", l.logPath, l.logName)
	file, err := openFile0755(filename)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed: %s", filename, err))
	}
	errorFilename := fmt.Sprintf("%s/%s.log.error", l.logPath, l.logName)
	errorFile, err := openFile0755(errorFilename)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed: %s", filename, err))
	}

	l.file, l.errorFile = file, errorFile

	// 另起协程，后台写日志
	go l.writeLogBackground()
}

// 在每次写日志之前，先检查日志文件是否达到切分条件
// 将isErrorFile传入而不是直接在写日志时同时检查两种日志文件，这样效率更好
func (l *FileLogger) checkFileAndSplit(isErrorFile bool) {
	switch l.splitType {
	case SPLIT_TIME:
		l.checkFileAndSplitByTime(isErrorFile)
	case SPLIT_SIZE:
		l.checkFileAndSplitBySize(isErrorFile)
	case SPLIT_NONE:
		// DO NOTHING
	}
}

func (l *FileLogger) checkFileAndSplitByTime(isErrorFile bool) {
	// 切分日志，生成备份文件
	if time.Now().Unix() - l.lastSplitTime >= l.splitTimeIntervalOrFileSize {
		curLogFile, backupLogFile := l.logFilenameAndBackupFilename(isErrorFile)
		l.backupLog(curLogFile, backupLogFile, isErrorFile)
	}
}

func (l *FileLogger) checkFileAndSplitBySize(isErrorFile bool) {
	// 获取两种文件名
	curLogFile, backupLogFile := l.logFilenameAndBackupFilename(isErrorFile)
	file := l.file
	if isErrorFile {file = l.errorFile}
	fileStat, _ := file.Stat()
	// 比较Byte数是否达到切分要求
	if fileStat.Size() >= l.splitTimeIntervalOrFileSize {
		l.backupLog(curLogFile, backupLogFile, isErrorFile)
	}
}

func (l *FileLogger) logFilenameAndBackupFilename(isErrorFile bool) (string, string) {
	now := time.Now()
	var backuplogFileName, curLogFileName string

	if isErrorFile {
		curLogFileName = fmt.Sprintf("%s/%s.log.error", l.logPath, l.logName)
		backuplogFileName = fmt.Sprintf("%s/%s.log.error_%04d%02d%02d%02d",
			l.logPath, l.logName, now.Year(), now.Month(), now.Day(), now.Hour())
	} else {
		curLogFileName = fmt.Sprintf("%s/%s.log", l.logPath, l.logName)
		backuplogFileName = fmt.Sprintf("%s/%s.log_%04d%02d%02d%02d",
			l.logPath, l.logName, now.Year(), now.Month(), now.Day(), now.Hour())
	}

	return curLogFileName, backuplogFileName
}

func (l *FileLogger) backupLog(curLogFile, backupLogFile string, isErrorFile bool) {
	if isErrorFile {
		// 关闭当前日志文件
		_ = l.errorFile.Close()
		// 重命名为备份文件
		_ = os.Rename(curLogFile, backupLogFile)
		// 生成新文件。由于这种情形下不太可能出现致命错误(初始化时能创建，现在应该也能创建)，
		// 而且这时不应该让程序崩溃，所以忽略错误
		l.errorFile, _ = openFile0755(curLogFile)
	} else {
		_ = l.file.Close()
		_ = os.Rename(curLogFile, backupLogFile)
		l.file, _ = openFile0755(curLogFile)
	}
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
			// 先检查Log是否需要切分
			l.checkFileAndSplit(false)
			_, _ = fmt.Fprint(l.file, rec.String())
		} else {
			l.checkFileAndSplit(true)
			_, _ = fmt.Fprint(l.errorFile, rec.String())
		}
	}
}


