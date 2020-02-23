package elog

import "time"

// Log 日志消息
// 格式信息由于通常字符串表示，并且固定，不太应该写到Log内，而是应该给到Logger或者其他方式外部传入.
type Log struct {

	// Level 日志级别
	Level Level

	// Time 日志发生时间
	Time time.Time

	// Caller 发生日志记录时的调用方法或函数
	Caller string

	// Msg 日志正常消息
	// TODO: Msg可能会修改为结构体，实现Msg部分定制
	Msg string

	// Err 日志错误消息
	// TODO: Err可能会实现为Err结构体
	Err string
}

// String 实现fmt.Stringer接口
func (l *Log) String() string {
	return ""
}

// FormatString 根据格式信息输出
// 难点在于如何指定Log结构体内字段进行打印。
// TODO: 暂时放弃
func (l *Log) FormatString(format string) string {
	return ""
}