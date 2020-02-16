package queue

type Queue interface {
	Enqueue(v interface{})
	Dequeue() (v interface{}, err error)
	Front() (v interface{}, err error)
	Size() int
	IsEmpty() bool
	String() string
}
