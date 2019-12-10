package etcp

// Option 服务器配置选项
type Option struct {

	// IPVersion 服务器绑定的IP版本
	IPVersion string

	// Host 当前服务器主机IP
	Host string

	// TcpPort 当前服务器主机监听端口号
	TcpPort int

	// Name 当前服务器名称
	Name string

	// Version 当前版本号
	Version string

	// MaxPacketSize 数据包的最大值
	MaxPacketSize uint32

	// MaxConn 当前服务器主机允许的最大链接个数
	MaxConn int

	// WorkerPoolSize 工作池大小
	WorkerPoolSize uint32

	// MaxWorkerTaskLen 业务工作Worker对应负责的任务队列最大任务存储数量
	MaxWorkerTaskLen uint32
}

func DefaultOption() *Option {
	return &Option{
		IPVersion:        "tcp4",
		Host:             "localhost",
		TcpPort:          8000,
		Name:             "ego-tcp",
		Version:          "v0.1",
		MaxPacketSize:    4096,
		MaxConn:          100,
		WorkerPoolSize:   20,
		MaxWorkerTaskLen: 1024,
	}
}

func (op *Option) SetHost(host string) *Option {
	op.Host = host
	return op
}

func (op *Option) SetPort(port int) *Option {
	op.TcpPort = port
	return op
}

func (op *Option) SetName(name string) *Option {
	op.Name = name
	return op
}

func (op *Option) SetMaxPacketSize(size int) *Option {
	op.MaxPacketSize = uint32(size)
	return op
}

func (op *Option) SetMaxConn(num int) *Option {
	op.MaxConn = num
	return op
}

func (op *Option) SetWorkerPoolSize(size int) *Option {
	op.WorkerPoolSize = uint32(size)
	return op
}

func (op *Option) SetMax(size int) *Option {
	op.MaxWorkerTaskLen = uint32(size)
	return op
}
