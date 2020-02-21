package elog2

import (
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	opt := &LoggerOption{
		Level:   DEBUG,
		WriteTo: os.Stdout,
	}
	logger, _ := NewLogger(opt)
	defer logger.Close()
	logger.log(FATAL, "fatal error!\n")
	logger.Debug("user id %d logged in\n", 37419)
	logger.Warn("something is occurred\n")
	logger.Fatal("fatal error!\n")
}

func TestFileLogger(t *testing.T) {
	opt := &FileLoggerOption{
		Level:   DEBUG,
		LogPath: "./tmp",
		LogName: "tes",
	}
	logger, _ := NewFileLogger(opt)
	logger.Debug("user id %d logged in\n", 37419)
	logger.Warn("something is occurred\n")
	logger.Fatal("fatal error!\n")
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
