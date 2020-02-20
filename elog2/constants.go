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