package lru

import (
	"container/list"
)

// 使用map和container/list(go标准库提供的双链表)实现LRU(Least Recently Used)最近最少使用淘汰算法

// Cache 代表使用了LRU淘汰算法的缓存
type Cache struct {
	// 通常在实现数据结构的时候，往往会使用cap和size来标记
	// 容量和已使用情况。但是作为缓存，内存是最关心的问题，
	// 由于存储的数据单元大小不确定，因此使用字节数作为空间占用标准更为合适
	// 此外，这里的内存限制忽略了接口和结构体包裹数据所占的空间(这部分空间不论具体存的什么消耗都一致，可以不考虑)
	maxBytes int64	// 每个缓存允许使用的最大内存
	nBytes int64	// 当前已使用的内存

	// 实现LRU策略的数据结构: 哈希表 + 双链表
	ll *list.List
	m map[string]*list.Element

	// 移除某条记录时的回调函数
	OnEvicted func(key string, value Value)
}

// Entry 代表了每一条存入双链表ll的实体(记录)
type Entry struct {
	key string
	value Value
}

// 任何实现了Len()方法的类型都可以作为值存进缓存
type Value interface {
	Len() int
}

// 构造方法
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		m:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 查
// 从字典查对应的双链表节点，得到之后将节点移到双链表尾部
// 为了描述方便，双链表一端为头，一端为尾，头端添加数据，
// 尾端移除(也就是扮演队列)
// 如果该节点不存在，则要将该节点插入双链表(这部分逻辑见添加操作)(目前这里并没这么做)
func (c *Cache) Get(key string) (value Value, ok bool) {
	if elem, ok := c.m[key]; ok {
		c.ll.MoveToFront(elem)
		kv := elem.Value.(*Entry)
		return kv.value, true
	}

	return
}

// 删
// 缓存淘汰，移除最近最少访问的节点
// 由于这不是通用的数据结构，因此移除时不返回移除的那个节点值
// remove对于外部使用者并不需要，因此设为私有
func (c *Cache) remove() {
	// 获取待删除的那个节点
	elem := c.ll.Back()		// 取出链表尾部节点，也就是最近最少使用的节点
	if elem != nil {
		c.ll.Remove(elem)	// 从双链表移除
		kv := elem.Value.(*Entry)
		delete(c.m, kv.key)	// 从哈希表删除
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())	// 更新可用字节数
		// 执行删除节点的回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 增/改
// 若节点存在，则把节点挪至链表头部，并更新值
// 若不存在，则是新增操作。
// 		1.1)若缓存已满或者未满但再把当前元素放入就超过maxBytes，
// 		就先删除链表尾部节点直至容量足够；
//		1.2)若缓存容量足够，直接插到链表头
// 		2)同时更新哈希表
func (c *Cache) Add(key string, value Value) {
	// 尝试获取节点，如果存在，则修改值并拉到链表头
	if elem, ok := c.m[key]; ok {
		c.ll.MoveToFront(elem)
		kv := elem.Value.(*Entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
		return
	}

	// 否则的话，加入队头
	elem := c.ll.PushFront(&Entry{key, value})
	c.m[key] = elem
	c.nBytes += int64(len(key)) + int64(value.Len())

	// 看下容量是否不足，不足的话删队尾
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.remove()		// 移除末尾最近最少使用的节点
	}
}

// 返回缓存数据节点数
func (c *Cache) Len() int {
	return c.ll.Len()
}