package log

import (
	"encoding/json"
	"io"
)

// Logger 日志记录器
// Logger绑定了一系列的日志记录方法，必须含有一个Run方法，运行中开辟协程接收日志并写入到存储
// 尝试实现结构化日志
type Logger struct {

	// storage 日志存储
	// TODO: 在没实现Storage接口之前，先使用io.ReadWriter顶替
	storage io.ReadWriter

	// logChan 日志channel。日志记录被调用时日志被写入channel，并被接收
	logChan chan ILog

	// colors 日志级别组对应的颜色
	colors map[Level]Color

}

// NewLogger 新建日志器，负责资源的分配初始化。比如说map、chan的创建
func NewLogger() ILogger {
	return &Logger{
		storage: nil,
		logChan: nil,
		colors: nil ,
	}
}

// Run 运行日志记录器，在 NewLogger() 及配置设置 之后开辟goroutine调用
func (l *Logger) Run() {
	go func() {

		for {
			oneLog := <- l.logChan
			// 由于ILog接口，json并不适合序列化
			// 若使用gob序列化，需要提前注册好实现ILog接口的结构体或结构体指针
			// TODO:
			oneLogData, _ := json.Marshal(oneLog)
			_, _= l.storage.Write(oneLogData)

			// TODO: 对于接收过程发生的错误，应跳过
			// TODO: 对于调用方法发生的致命错误（需要退出线程），应想办法把该致命错误消息先传到这个再退出。
			//  为了实现这样的效果，需要手动控制return，而不是调用Fatal这类API
		}

	}()
}



// SetColors 为全部日志级别统一设置颜色组
func (l *Logger) SetColors(colors map[Level]Color) ILogger {
	l.colors = colors
	return l
}

// SetColor 为不同的日志级别定制颜色
func (l *Logger) SetColor(loglevel Level, color Color) ILogger {
	l.colors[loglevel] = color
	return l
}