package elog2

// elog2.Interface
type Interface interface {
	SetLevel(level int)
	Debug(format string, args ...interface{})
	Trace(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Close()
}

type Option interface {
	Check() error
}
