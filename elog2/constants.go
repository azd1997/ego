package elog2

const (
	DEBUG = iota
	TRACE
	INFO
	WARN
	ERROR
	FATAL
)

const DEFAULT_LOG_LEVEL = DEBUG

var LevelMap = map[int]string{
	DEBUG: "[DEBUG]",
	TRACE: "[TRACE]",
	INFO: "[INFO]",
	WARN: "[WARN]",
	ERROR: "[ERROR]",
	FATAL: "[FATAL]",
}

const DEFAULT_RECORD_CHAN_SIZE = 50000


// 日志切分
const (
	DEFAULT_SPLIT_TYPE = "time"
	SPLIT_TIME = "time"
	SPLIT_SIZE = "size"
	SPLIT_NONE = "none"

	// 1h = 60min = 3600s
	SPLIT_TIME_HOUR = 3600

	// 100MB = 1e8 B
	SPLIT_SIZE_100M = 1e8
)