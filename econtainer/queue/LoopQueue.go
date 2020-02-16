package queue

import (
	"errors"
	"fmt"
	"strings"
)

// 基于切片的循环队列
// 循环队列需要浪费一个实际空间，来使得front和tail来判空判满
// 如果不的话，会发现front==tail既会出现在队列空的情况，也会出现在队列满的情况
// 因此需要浪费一个单元空间，来产生明确的判空判满条件
// 判空 front == tail
// 判满 (tail+1) % len(data) == front

type LoopQueue struct {
	data []interface{}
	size int
	front, tail int
	resizeFactor int
}

func NewLoopQueue(initcap int) *LoopQueue {
	data := make([]interface{}, initcap+1, initcap+1)
	return &LoopQueue{
		data:         data,
		resizeFactor: 2,
	}
}

func (q *LoopQueue) Enqueue(v interface{}) {
	if q.IsFull() {	// 满则扩容
		q.resize(q.Cap() * q.resizeFactor)
	}
	q.data[q.tail] = v
	q.tail = (q.tail + 1) % len(q.data)
	q.size++
}

func (q *LoopQueue) Dequeue() (v interface{}, err error) {
	if q.size == 0 {
		return nil, errors.New("queue is empty")
	}

	ret := q.data[q.front]
	q.front = (q.front + 1) % len(q.data)
	q.size--

	// 判断是否缩容
	if q.size == q.Cap() / (2 * q.resizeFactor) && q.Cap() / 2 != 0 {
		q.resize(q.Cap() / q.resizeFactor)
	}

	return ret, nil
}

func (q *LoopQueue) Front() (v interface{}, err error) {
	if q.size == 0 {
		return nil, errors.New("queue is empty")
	}
	return q.data[q.front], nil
}

func (q *LoopQueue) Size() int {
	return q.size
}

func (q *LoopQueue) IsEmpty() bool {
	return q.front==q.tail
}

func (q *LoopQueue) String() string {
	res := strings.Builder{}
	res.WriteString(fmt.Sprintf("LoopQueue: size = %d, cap = %d\n", q.size, q.Cap()))
	res.WriteString("front [")
	for i := 0; i < q.size; i++ {
		res.WriteString(fmt.Sprintf("%d", q.data[(i + q.front) % len(q.data)]))
		if (i+q.front+1) % len(q.data) != q.tail {	// 没到最后一个元素
			res.WriteString(", ")
		}
	}
	res.WriteString("] tail")
	return res.String()
}

func (q *LoopQueue) IsFull() bool {
	return (q.tail + 1) % len(q.data) == q.front
}

// 实际数据容量，没太大意义，会扩容
func (q *LoopQueue) Cap() int {
	return len(q.data) - 1
}

func (q *LoopQueue) resize(newcap int) {
	newdata := make([]interface{}, newcap+1)
	for i := 0; i < q.size; i++ {
		newdata[i] = q.data[(i + q.front) % len(q.data)]
	}
	q.data = newdata
	q.front, q.tail = 0, q.size
}