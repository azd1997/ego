package elog2

import (
	"fmt"
)

// Record是日志写入文件、数据库、网络流等的基本数据单位
// 除了控制台日志以外，日志的写应该是异步的，否则当日志量大时，容易阻塞，成为性能瓶颈
// 因此Record用于异步写日志的传输单位，当然控制台也能用
type Record struct {
	Time string
	Level int
	LineNo int
	FileName string
	FuncName string
	Msg string
}

func record(logLevel int, format string, args ...interface{}) *Record {
	filename, funcname, lineno := GetLineInfo()
	now := Now()
	msg := fmt.Sprintf(format, args...)


	rec := &Record{
		Time:     now,
		Level:    logLevel,
		FileName: filename,
		FuncName: funcname,
		LineNo:   lineno,
		Msg:msg,
	}

	return rec
}

func (r *Record) String() string {
	return fmt.Sprintf("%s %s (%s:%s:%d) %s",
		r.Time,
		LevelMap[r.Level],
		r.FileName,
		r.FuncName,
		r.LineNo,
		r.Msg)
}