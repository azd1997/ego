package elog2

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestNow(t *testing.T) {
	str := Now()
	log.Println(str)
}

func TestGetLineInfo(t *testing.T) {
	fmt.Println(GetLineInfo())
}

func TestWriteLog(t *testing.T) {
	writeLog(os.Stdout, DEBUG, "%s\n", "天下第一")
	writeLog(os.Stdout, WARN, "%s\n", "天下第一")
	writeLog(os.Stdout, FATAL, "something is occurred\n")
}
