package queue

import (
	"fmt"
	"testing"
)

func TestLoopQueue(t *testing.T) {
	queue := NewLoopQueue(5)
	for i:=0; i<7; i++ {
		queue.Enqueue(i)
		fmt.Println(queue)
	}
}
