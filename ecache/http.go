package ecache


const defaultBasePath = "/_ecache/"

// HTTPPool 节点池
type HTTPPool struct {
	// 自己的地址 ip:port
	self string
	// 节点间通讯地址的前缀，例如 http://example.com/_ecache/开头的请求
	// 就用于节点间的访问。因为一个主机上还可能承载其他的服务，加一段 Path 是一个好习惯。
	// 比如，大部分网站的 API 接口，一般以 /api 作为前缀。
	basePath string
}

// NewHTTPPool 初始化一个HTTP节点池
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}