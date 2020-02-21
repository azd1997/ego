package elog2

const (
	DEBUG = iota
	TRACE
	INFO
	WARN
	ERROR
	FATAL
)

var DEFAULT_LOG_LEVEL = DEBUG

var LevelMap = map[int]string{
	DEBUG: "[DEBUG]",
	TRACE: "[TRACE]",
	INFO: "[INFO]",
	WARN: "[WARN]",
	ERROR: "[ERROR]",
	FATAL: "[FATAL]",
}