package elog2

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

// 测试用的输出源

type testOutput struct {
	io.WriteCloser
	size int64
	name string		// 规定这个name为 name+time(年月日时共10位)形式
}

func (o *testOutput) Size() int64 {
	return o.size
}

func (o *testOutput) Name() string {
	return o.name
}

func (o *testOutput) Rename(newname string) {
	o.name = newname
}

func (o *testOutput) BackupName() string {
	l := len(o.name)
	basename := o.name[:l-10]
	now := time.Now()
	nextname := fmt.Sprintf("%s%04d%02d%02d%02d",
		basename, now.Year(), now.Month(), now.Day(), now.Hour())
	return nextname
}

// 这个size增加的方法由于不同输出源不同，还是交给使用着自己决定是否要添加
// 要的话自己在elog2提供的日志函数之上再包装一次，增加AddSize方法
func (o *testOutput) AddSize(delta int64) {
	o.size += delta
}

func newTestOutput(name string) OutputSource {
	return &testOutput{
		WriteCloser: os.Stdout,
		size:        0,
		name:        name,
	}
}



func TestLogger(t *testing.T) {

	now := time.Now()
	basename := "testlog-"
	logname := fmt.Sprintf("%s%04d%02d%02d%02d",
		basename, now.Year(), now.Month(), now.Day(), now.Hour())

	opt := &LoggerOption{
		Level:   DEBUG,
		WriteTo: newTestOutput(logname),
		SplitType:SPLIT_SIZE,
		SplitTimeIntervalOrFileSize:1000,	// 测试用1000B，快一点
		NewOutputSource:newTestOutput,
	}
	logger, _ := NewLogger(opt)
	defer logger.Close()

	for i:=0; i<1000; i++ {
		// 打印当前输出源的信息
		fmt.Printf("now, testOutput : size = %d, name = %s\n\n", logger.writeTo.Size(), logger.writeTo.Name())

		logger.log(FATAL, "fatal error!\n")
		logger.Debug("user id %d logged in\n", 37419)
		logger.Warn("something is occurred\n")
		logger.Fatal("fatal error!\n")

		// 虚拟测试，假设增加了100B
		logger.writeTo.(*testOutput).AddSize(100)
	}

}

func TestFileLogger(t *testing.T) {
	opt := &FileLoggerOption{
		Level:   DEBUG,
		LogPath: "./tmp",
		LogName: "tes",
		RecordChanSize:10000,	// 这里缓冲区要大些，不然会阻塞住，因为下边日志产出太快
		SplitType:SPLIT_SIZE,
		SplitTimeIntervalOrFileSize:1000,
	}
	logger, _ := NewFileLogger(opt)

	for i:=0; i<2000; i++ {
		logger.Debug("user id %d logged in\n", 37419)
		logger.Warn("something is occurred\n")
		logger.Fatal("fatal error!\n")
	}

}

func TestConsoleLogger(t *testing.T) {
	opt := &ConsoleLoggerOption{Level:DEBUG}
	logger, _ := NewConsoleLogger(opt)
	logger.Debug("user id %d logged in\n", 37419)
	logger.Warn("something is occurred\n")
	logger.Fatal("fatal error!\n")
	logger.Debug("user id %d logged in\n", 37419)

	for i:=0; i<10; i++ {
		//logger.Debug("user id %d logged in\n", 37419)
		logger.Warn("something is occurred%s\n", "eeeeeeee")
		//logger.Warn("something is occurred\n")
	}
}

func TestELog(t *testing.T) {
	name := "console"
	opt := &ConsoleLoggerOption{Level:DEBUG}
	InitELog(name, opt)
	ELog.Warn("something is occurred\n")
}
