// Package ecolor 彩色字符串打印
package ecolor

import (
	"fmt"
	"io"
)

// Color 颜色
type Color uint8

// BColor 背景颜色
// Linux终端背景颜色设置只有 40~47。设置其它数字无效
type BColor Color

// FColor 前景颜色
// Linux终端前景颜色（字体颜色）设置只有 30~37。但是设置90+也可以。 （有待研究）
type FColor Color

// 显示方式
type Display uint8

// 背景颜色
const (
	BBlack BColor = iota + 40
	BRed
	BGreen
	BYellow
	BBlue
	BMagenta // 紫红色
	BCyan      // 青蓝色
	BWhite

	B_Red BColor = 91
	B_Green	BColor = 92// 92
	B_Yellow BColor = 93	 // 93
	B_Blue BColor = 94	// 94
	B_MAGENTA BColor = 95		// 95
)

// 前景颜色
const (
	FBlack FColor = iota + 30
	FRed
	FGreen
	FYellow
	FBlue
	FMagenta // 紫红色
	FCyan      // 青蓝色
	FWhite

	F_Red FColor = 91
	F_Green	FColor = 92// 92
	F_Yellow FColor = 93	 // 93
	F_Blue FColor = 94	// 94
	F_MAGENTA FColor = 95		// 95
)

// 终端显示方式
const (
	Default   Display = iota
	Highlight // 1 高亮
	_
	_
	Underline // 4 下划线
	Twinkle   // 5 闪烁
	_
	ReverseDisplay // 7 反白显示
	Invisible      // 8 不可见
)

// 0x1B或者\x1b是标记; [开始颜色定义

// Print Printf Println
// SPrint SPrintf SPrintln
// FPrint FPrintf FPrintln

// Printf 将字符串转为彩色字符串打印
func Printf(format string, d Display, b BColor, f FColor, a ...interface{}) {
	fmt.Printf("\x1b[%d;%d;%dm%s\x1b[0m", d, b, f,
		fmt.Sprintf(format, a...))
}

// Println 将字符串转为彩色字符串打印
func Println(d Display, b BColor, f FColor, a ...interface{}) {
	fmt.Printf("\x1b[%d;%d;%dm%s\x1b[0m\n", d, b, f,
		fmt.Sprint(a...))
}

// Print 将字符串转为彩色字符串打印
func Print(d Display, b BColor, f FColor, a ...interface{}) {
	fmt.Printf("\x1b[%d;%d;%dm%s\x1b[0m", d, b, f,
		fmt.Sprint(a...))
}

// Sprintf 将字符串转为彩色字符串，当它重新被直接打印在控制台时，将是彩色输出
func Sprintf(format string, d Display, b BColor, f FColor, a ...interface{}) string {
	return fmt.Sprintf("\x1b[%d;%d;%dm%s\x1b[0m", d, b, f,
		fmt.Sprintf(format, a...))
}

// Sprintln 将字符串转为彩色字符串，当它重新被直接打印在控制台时，将是彩色输出
func Sprintln(d Display, b BColor, f FColor, a ...interface{}) string {
	return fmt.Sprintf("\x1b[%d;%d;%dm%s\x1b[0m\n", d, b, f,
		fmt.Sprint(a...))
}

// Sprint 将字符串转为彩色字符串，当它重新被直接打印在控制台时，将是彩色输出
func Sprint(d Display, b BColor, f FColor, a ...interface{}) string {
	return fmt.Sprintf("\x1b[%d;%d;%dm%s\x1b[0m", d, b, f,
		fmt.Sprint(a...))
}

// Fprintf 将字符串转为彩色字符串，当它重新被直接打印在控制台时，将是彩色输出
func Fprintf(w io.Writer, format string, d Display, b BColor, f FColor, a ...interface{}) (n int, err error) {
	n, err = fmt.Fprint(w, Sprintf(format, d, b, f, a...))
	return n-14, err	// -14是附带的颜色等信息占用
}

// Fprintln 将字符串转为彩色字符串，当它重新被直接打印在控制台时，将是彩色输出
func Fprintln(w io.Writer, d Display, b BColor, f FColor, a ...interface{}) (n int, err error) {
	n, err = fmt.Fprint(w, Sprintln(d, b, f, a...))
	return n-14, err	// -14是附带的颜色等信息占用
}

// Fprintf 将字符串转为彩色字符串，当它重新被直接打印在控制台时，将是彩色输出
func Fprint(w io.Writer, d Display, b BColor, f FColor, a ...interface{}) (n int, err error) {
	n, err = fmt.Fprint(w, Sprint(d, b, f, a...))
	return n-14, err	// -14是附带的颜色等信息占用
}