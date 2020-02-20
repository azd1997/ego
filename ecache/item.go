package ecache

// Item 为缓存的数据存储单元
// data存储真实的缓存数据，字节数组方便转换成其他各种类型的数据
// 对data做封装的目的是保证其只读，不能被修改
type Item struct {
	data []byte
}

// 实现Value接口
func (v Item) Len() int {
	return len(v.data)
}

// ByteSlice 返回数据的拷贝
func (v Item) ByteSlice() []byte {
	return cloneBytes(v.data)
}

// String 返回数据的字符串形式
func (v Item) String() string {
	return string(v.data)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}