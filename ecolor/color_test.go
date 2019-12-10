package ecolor

import (
	"fmt"
	"os"
	"testing"
)

var bColors = []BColor{BBlack, BRed, BGreen, BYellow, BBlue, BMagenta, BCyan, BWhite, B_Red, B_Green, B_Yellow, B_Blue, B_MAGENTA}
var fColors = []FColor{FBlack, FRed, FGreen, FYellow, FBlue, FMagenta, FCyan, FWhite, F_Red, F_Green, F_Yellow, F_Blue, F_MAGENTA}
var disPlays = []Display{Default, Highlight, Underline, Twinkle, ReverseDisplay, Invisible}

func TestPrintf(t *testing.T) {
	 for _, b := range bColors {
	   for _, f := range fColors {
	     for _, d := range disPlays { // 显示方式 = 0,1,4,5,7,8
	       Printf("%s", d, b, f, "Hello Eiger")
	     }
	     fmt.Println("")
	   }
	   fmt.Println("")
	 }
}

func TestPrintln(t *testing.T) {
	for _, b := range bColors {
		for _, f := range fColors {
			for _, d := range disPlays { // 显示方式 = 0,1,4,5,7,8
				Println(d, b, f, "Hello Eiger")
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func TestPrint(t *testing.T) {
	for _, b := range bColors {
		for _, f := range fColors {
			for _, d := range disPlays { // 显示方式 = 0,1,4,5,7,8
				Print(d, b, f, "Hello Eiger")
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func TestSprintf(t *testing.T) {
	for _, b := range bColors {
		for _, f := range fColors {
			for _, d := range disPlays { // 显示方式 = 0,1,4,5,7,8
				fmt.Print(" ", Sprintf("%s", d, b, f, "Hello Eiger"), " ")
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func TestSprintln(t *testing.T) {
	for _, b := range bColors {
		for _, f := range fColors {
			for _, d := range disPlays { // 显示方式 = 0,1,4,5,7,8
				fmt.Print(" ", Sprintln(d, b, f, "Hello Eiger"), " ")
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func TestSprint(t *testing.T) {
	for _, b := range bColors {
		for _, f := range fColors {
			for _, d := range disPlays { // 显示方式 = 0,1,4,5,7,8
				fmt.Print(" ", Sprint(d, b, f, "Hello Eiger"), " ")
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func TestFprint(t *testing.T) {
	for _, b := range bColors {
		for _, f := range fColors {
			for _, d := range disPlays { // 显示方式 = 0,1,4,5,7,8
				file1, _ := os.Create("./fprintf1.txt")
				file2, _ := os.Create("./fprintf2.txt")
				cnt1, _ := fmt.Fprint(file1, "Hello Eiger")
				cnt2, _ := Fprint(file2, d, b, f, "Hello Eiger")
				fmt.Println("cnt1=", cnt1, "cnt2=", cnt2)
				// 比较两个文件内容，可知实际内容并不一致，彩色字符串内容两端包含了颜色信息。一般情况下不需要写入到文件
				// 而相应的读取出来的长度也相应地减去 14 （这是设置了全参数的彩色打印信息，如果不是设置这么多项需要调整该数字）
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func TestPrintRed(t *testing.T) {
	PrintCF(FRed, "Hello Eiger")
}

// ------------------------------------------
// 测试彩色打印

func TestPrintWithColor(t *testing.T) {
	for _, b := range bColors {
		for _, f := range fColors {
			for _, d := range disPlays { // 显示方式 = 0,1,4,5,7,8
				printWithColor5(d, b, f)
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func printWithColor1(d Display, b BColor, f FColor) {
	fmt.Printf(" %c[%d;%d;%dm%sf=%d,b=%d,d=%d%c[0m ", 0x1B, d, b, f, "", f, b, d, 0x1B)
}
func printWithColor2(d Display, b BColor, f FColor) {
	fmt.Printf(" %c[%d;%d;%dm%s(f=%d,b=%d,d=%d)%c[0m ", 0x1B, d, b, f, "", f, b, d, 0x1B)
}

func printWithColor3(d Display, b BColor, f FColor) {
	format := "(f=%d,b=%d,d=%d)"
	format = " %c[%d;%d;%dm%s" + format + "%c[0m "
	fmt.Printf(format, 0x1B, d, b, f, "", f, b, d, 0x1B)
}

func printWithColor4(d Display, b BColor, f FColor) {
	format := "(f=%d,b=%d,d=%d)"
	format = " \x1b[%d;%d;%dm%s" + format + "\x1b[0m "
	fmt.Printf(format, d, b, f, "", f, b, d)
}

func printWithColor5(d Display, b BColor, f FColor) {
	format := "%s"
	format = " \x1b[%d;%d;%dm" + format + "\x1b[0m "
	fmt.Printf(format, d, b, f, "Hello Eiger")
}

// PrintfAny 根据格式输出任意类型值，但是输出的字符串会被方括号包裹
//func PrintfAny(format string, d Display, b BColor, f FColor, a ...interface{}) {
//	format = "\x1b[%d;%d;%dm" + format + "\x1b[0m"
//	fmt.Printf(format, d, b, f, a)
//}