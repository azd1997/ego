package elog2

import (
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := NewLogger(DEBUG, os.Stdout)
	defer logger.Close()
	logger.log(FATAL, "fatal error!\n")
	logger.Debug("user id %d logged in\n", 37419)
	logger.Warn("something is occurred\n")
	logger.Fatal("fatal error!\n")
}

func TestFileLogger(t *testing.T) {
	logger := NewFileLogger(DEBUG, "./tmp", "tes")
	logger.Debug("user id %d logged in\n", 37419)
	logger.Warn("something is occurred\n")
	logger.Fatal("fatal error!\n")
}

func TestConsoleLogger(t *testing.T) {
	logger := NewConsoleLogger(DEBUG)
	logger.Debug("user id %d logged in\n", 37419)
	logger.Warn("something is occurred\n")
	logger.Fatal("fatal error!\n")
	logger.Debug("user id %d logged in\n", 37419)
}
