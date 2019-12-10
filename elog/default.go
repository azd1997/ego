package elog

import "io"

// DefaultColors 默认颜色组
var DefaultColors = map[Level]Color{
	L_INFO: BLUE,
	L_TRAC: YELLOW,
	L_WARN: MAGENTA,
	L_ERRO: RED,
	L_SUCC: GREEN,
}

// DefaultStorage 默认日志存储
var DefaultStorage io.ReadWriter

// DefaultLogger 默认日志记录器
var DefaultLogger = &Logger{
	storage: nil,
	logChan: make(chan ILog),
	colors:  DefaultColors,
}

// DefaultFormat 默认日志打印格式
// 见types.go


