package elog2

import (
	"fmt"
	"log"
	"testing"
)

func TestNow(t *testing.T) {
	str := Now()
	log.Println(str)
}

func TestGetLineInfo(t *testing.T) {
	fmt.Println(GetLineInfo())
}
