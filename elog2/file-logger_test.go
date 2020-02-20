package elog2

import "testing"

func TestFileLogger(t *testing.T) {
	logger := NewFileLogger(DEBUG, "./tmp", "tes")
	logger.Debug("user id %d logged in\n", 37419)
	logger.Warn("something is occurred\n")
	logger.Fatal("fatal error!\n")
}
