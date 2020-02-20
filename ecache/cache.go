package ecache

import (
	"sync"

	"github.com/azd1997/ego/ecache/strategy/lru"
)

// 对lru.Cache进行封装
type cache struct {
	sync.Mutex
	lru        *lru.Cache	// LRU缓存
	cacheBytes int64		// 缓存字节容量
}

func (c *cache) add(key string, value Item) {
	c.Lock()
	defer c.Unlock()

	// 在 add 方法中，判断了 c.lru 是否为 nil，如果不等于 nil 再创建实例。
	// 这种方法称之为延迟初始化(Lazy Initialization)，
	// 一个对象的延迟初始化意味着该对象的创建将会延迟至第一次使用该对象时。
	// 主要用于提高性能，并减少程序内存要求。
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}

	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value Item, ok bool) {
	c.Lock()
	defer c.Unlock()

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(Item), ok
	}

	return
}
