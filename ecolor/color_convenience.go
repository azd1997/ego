package ecolor

import "fmt"

// 彩色字符串常用的只是在打印到终端（os.StdOut）时设置字体前景色，这里给出几种常用颜色的打印

// PrintCF 仅设置字体颜色打印
func PrintCF(color FColor, a ...interface{}) {
	fmt.Printf("\x1b[%dm%s\x1b[0m", color,
		fmt.Sprint(a...))
}

