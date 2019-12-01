package log

// 使用接口作一次抽象的目的是：允许用户定制具体的实现

// ILogger 接口
type ILogger interface {

	// Run 运行日志记录器
	Run()

	// SetColors 为全部日志级别统一设置颜色组
	SetColors(colors map[Level]Color) ILogger

	// SetColor 为不同的日志级别定制颜色
	SetColor(loglevel Level, color Color) ILogger
}

// ILog 日志消息接口
// Log的定制化相对复杂，这里ILog暂时是弃用的
type ILog interface {

	// 定制 相关方法

	// Format 根据格式定制日志消息
	Format(format string) ILog


	// 功能 相关方法

	// String 实现fmt.Stringer接口，根据打印格式输出
	String() string


}
