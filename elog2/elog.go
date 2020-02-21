package elog2

// 提供一个默认的日志库实例

var ELog Interface

func InitELog(name string, opt Option) (err error) {
	switch name {
	case "file":
		ELog, err = NewFileLogger(opt)
	case "console":
		ELog, err = NewConsoleLogger(opt)
	case "io":
		ELog, err = NewLogger(opt)
	default:
		ELog, err = NewConsoleLogger(opt)
	}

	// 出错时，至少让控制台日志生效
	if err != nil {
		ELog, _ = NewConsoleLogger(opt)
		return err
	}

	return nil
}

func Debug(format string, args ...interface{}) {
	ELog.Debug(format, args...)
}

func Trace(format string, args ...interface{}) {
	ELog.Trace(format, args...)
}

func Info(format string, args ...interface{}) {
	ELog.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	ELog.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	ELog.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	ELog.Fatal(format, args...)
}

func Close() {
	ELog.Close()
}