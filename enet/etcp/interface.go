package etcp

import "net"

type IServer interface {

	// Start 启动服务器
	Start()

	// Stop 停止服务器
	Stop()

	// Serve 运行服务器
	Serve()

	// AddRouter 路由功能，给当前的服务器注册一个路由方法，供客户端的连接处理使用
	AddRouter(msgId uint32, router IRouter)

	//得到链接管理
	GetConnMgr() IConnManager

	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func(IConnection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(IConnection))
	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConnection)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConnection)

	// 获取配置的拷贝
	Option() Option
}

// IConnection 连接的抽象接口
type IConnection interface {
	// Start 启动连接
	Start()

	// Stop 停止连接
	Stop()

	// GetTCPConn 获取当前连接绑定的socket
	GetTCPConn() *net.TCPConn

	// GetConnID 获取当前连接的ID
	GetConnID() uint32

	// RemoteAddr 获取远程客户端的TCP状态 IP Port
	RemoteAddr() net.Addr

	// SendMsg 发送数据
	SendMsg(msgId uint32, data []byte) error

	// SendBuffMsg 带缓冲发送数据
	SendBuffMsg(msgId uint32, data []byte) error

	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string)(interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
}

// 路由的抽象接口，路由里的数据都是IRequest
type IRouter interface {

	// PreHandle 处理conn业务之前的钩子方法(Hook)
	PreHandle(r IRequest)

	// Handle 处理conn业务的主方法(Hook)
	Handle(r IRequest)

	// PostHandle 处理conn业务之后的钩子方法(Hook)
	PostHandle(r IRequest)
}

// IMessage 将请求的一个消息封装到message中，定义抽象层接口
type IMessage interface {

	// GetMsgId 获取消息ID
	GetMsgId() uint32
	// GetDataLen 获取消息数据段长度
	GetDataLen() uint32
	// GetData 获取消息内容
	GetData() []byte

	// SetMsgId 设置消息ID
	SetMsgId(uint32)
	// 设置消息数据段长度
	SetDataLen(uint32)
	// 设置消息内容
	SetData([]byte)
}

// IRequest 将客户端请求的连接信息和请求的数据 包装到一个Request中
type IRequest interface {

	// GetConn 获取当前连接
	GetConn() IConnection

	// GetData 获取当前连接的请求数据
	GetData() []byte

	// GetMsgId 获取消息Id
	GetMsgId() uint32
}

// IMsgHandle 消息管理抽象层
type IMsgHandler interface{
	DoMsgHandler(request IRequest)          //马上以非阻塞方式处理消息
	AddRouter(msgId uint32, router IRouter) //为消息添加具体的处理逻辑
	StartWorkerPool()                       //启动worker工作池
	SendMsgToTaskQueue(request IRequest)    //将消息交给TaskQueue,由worker进行处理
}

// IConnManager 连接管理抽象层
type IConnManager interface {
	Add(conn IConnection)                   //添加链接
	Remove(conn IConnection)                //删除连接
	Get(connID uint32) (IConnection, error) //利用ConnID获取链接
	Len() int                               //获取当前连接
	ClearConn()                             //删除并停止所有链接
}

// IDataPack 封包数据和拆包数据
// 直接面向TCP连接中的数据流,为传输数据添加头部信息，用于处理TCP粘包问题。
type IDataPack interface{

	MaxPacketSize() int

	// GetHeadLen 获取包头长度方法
	GetHeadLen() uint32

	// Pack 封包方法
	Pack(msg IMessage)([]byte, error)

	// Unpack 拆包
	Unpack([]byte)(IMessage, error)
}




