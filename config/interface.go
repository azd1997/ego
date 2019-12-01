package config

// IConfig 公共配置接口。将由json/yml/toml/ini配置器、数据库类配置器等实现
type IConfig interface {

	// Get 根据键查值
	Get(key string) (value interface{})

	// Put 修改配置文件
	Put(key string, value interface{}) error

	// Reload 重新加载配置
	Reload() error
}
