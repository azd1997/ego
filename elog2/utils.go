package elog2

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

func openFile0755(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed: %s", filename, err))
	}
	return file
}



var timeFormat1 = "Mon Jan 2 15:04:05 -0700 MST 2006"
var timeFormat2 = "2006-01-02 15:04:05.999"		// .999说明精确到ms级别
// time包有预定义的一些时间格式,见time/format.go首部

// 获取当前时间格式化后的字符串
func Now() string {
	// 极度奇葩!!!
	// 不取[:0]直接打印得到的就是
	// 2020/02/20 10:04:17 2020-02-20T10:04:17-05:00
	// return time.Now().Format(time.RFC3339)[:0]
	//2020/02/20 10:03:33

	// 但是这么做传出去之后打印出的却又是空字符串了...太恶趣味了啊

	return time.Now().Format(timeFormat2)
}


func GetLineInfo() (filename string, funcName string, lineNo int) {
	//pc, file, line, ok := runtime.Caller(0)		// 0表示当前函数的直接调用者

	// 真正打印时希望获取的是Debug()这些函数的调用位置等信息
	// runtime.Caller <- GetLineInfo <- record <- FileLogger.log <- FileLogger.Debug/Warn/...
	pc, file, line, ok := runtime.Caller(4)		// 要获取Debug/Info/...的调用者的信息，就是输4

	if ok {
		// path.Base(path) 用于去除路径，只留下文件名
		// runtime.FuncForPC(pc)获取调用者函数信息
		filename, lineNo = path.Base(file), line
		funcName = path.Base(runtime.FuncForPC(pc).Name())
	}
	return
}