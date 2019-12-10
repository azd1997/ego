package elog

import "fmt"

// Level 日志级别
// 指定Level对响应级别的日志的配置； 指定上限可以限定日志级别上限
type Level uint8

const (
	L_INFO Level = iota
	L_TRAC
	L_ERRO
	L_WARN
	L_SUCC
)

// Prefix 日志前缀
type Prefix string

const (
	P_INFO Prefix = "[INFO]"
	P_TRAC Prefix = "[TRAC]"
	P_ERRO Prefix = "[ERRO]"
	P_WARN Prefix = "[WARN]"
	P_SUCC Prefix = "[SUCC]"
)

// Color 颜色
type Color uint8

func (c *Color) Color(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", c, s)
}

const (
	RED = Color(iota + 91)
	GREEN	// 92
	YELLOW	// 93
	BLUE	// 94
	MAGENTA		// 95
)

// Format 日志打印格式
type Format string

const (
	// TODO: 预置一些日志打印格式

	DefaultFormat Format = ""
)

