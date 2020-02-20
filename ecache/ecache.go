package ecache

import (
	"fmt"
	"log"
	"sync"
)

// Group是最核心的数据结构，负责用户交互，控制缓存值存储与获取的流程
//                            是
//接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
//                |  否                         是
//                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
//                            |  否
//                            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶


// 如果缓存不存在，应从数据源（文件，数据库等）获取数据并添加到缓存中。
// ECache 是否应该支持多种数据源的配置呢？
// 不应该，一是数据源的种类太多，没办法一一实现；而是扩展性不好。
// 如何从源头获取数据，应该是用户决定的事情，我们就把这件事交给用户好了。
// 因此，我们设计了一个回调函数(callback)，在缓存不存在时，调用这个函数，得到源数据。

// 定义一个函数类型 F，并且实现接口 A 的方法，然后在这个方法中调用自己。
// 这是 Go 语言中将其他函数（参数返回值定义与 F 一致）转换为接口 A 的常用技巧。

// Getter 用于给外部调用者自定义不同的数据源加载
type Getter interface {
	Get(key string) ([]byte, error)
}

// 实现Getter接口
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}


// 一个 Group 可以认为是一个缓存的命名空间，每个 Group 拥有一个唯一的名称 name。
// 比如可以创建三个 Group，缓存学生的成绩命名为 scores，
// 缓存学生信息的命名为 info，缓存学生课程的命名为 courses
type Group struct {
	name      string	// Group的标识
	getter    Getter	// 缓存未命中时获取源数据的回调(callback)
	cache cache		// 支持并发安全的缓存
}

// 定义一个Group组，所有实例化的Group都会记录到这里边
var (
	mu     sync.RWMutex		// 读保护groups
	groups = make(map[string]*Group)
)

// NewGroup 创建一个Group实例
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		cache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup 根据Group名字在groups中查询并返回Group实例
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get 从cache中获取值
func (g *Group) Get(key string) (Item, error) {
	// 参数检验
	if key == "" {
		return Item{}, fmt.Errorf("key is required")
	}

	// 缓存命中(就是要查的键存在)
	if v, ok := g.cache.get(key); ok {
		log.Println("[ECache] hit")
		return v, nil
	}

	// 键不存在，则从本地或远程获取获取
	return g.load(key)
}

// 从本地或远程获取，暂且只从本地获取
func (g *Group) load(key string) (value Item, err error) {
	return g.getFromLocal(key)
}

// 从本地获取数据
func (g *Group) getFromLocal(key string) (Item, error) {
	// 调用getter.Get()回调函数(用户自己定义如何获取)
	bytes, err := g.getter.Get(key)
	if err != nil {
		return Item{}, err

	}
	value := Item{data: cloneBytes(bytes)}
	// 将从本地获取的键值对添加到Group中
	g.addItem(key, value)
	return value, nil
}

func (g *Group) addItem(key string, value Item) {
	g.cache.add(key, value)
}