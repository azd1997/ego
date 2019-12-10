package etcp

import (
	"fmt"
	"net"
)

// Server iface.IServer接口的一个实现
type Server struct {

	// Opts 参数配置
	Opts *Option

	// MsgHandler 服务端注册的连接对应的消息管理模块（多路由）
	MsgHandler IMsgHandler

	//当前Server的链接管理器
	ConnMgr IConnManager

	//该Server的连接创建时Hook函数
	OnConnStart func(conn IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn IConnection)
}

// NewServer 新建一个Server
func NewServer(opts *Option) IServer {

	return &Server{
		Opts:opts,
		MsgHandler: NewMsgHandler(int(opts.WorkerPoolSize)),
		ConnMgr:NewConnManager(),
	}
}

// Start 启动
func (s *Server) Start() {
	fmt.Printf("[Start] Server starts listening at IP: %s, Port: %d\n", s.Opts.Host, s.Opts.TcpPort)
	fmt.Printf("[ETCP] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		s.Opts.Version,
		s.Opts.MaxConn,
		s.Opts.MaxPacketSize)

	// 异步启动，避免阻塞主线程
	go func() {

		// 0. 启动 worker 工作池机制
		s.MsgHandler.StartWorkerPool()

		// 1. 获取一个TCP Addr
		addr, err := net.ResolveTCPAddr(s.Opts.IPVersion, fmt.Sprintf("%s:%d", s.Opts.Host, s.Opts.TcpPort))
		if err != nil {
			fmt.Printf("Start: resolving tcp addr err: %s\n", err)
			return
		}

		// 2. 监听服务器地址
		listener, err := net.ListenTCP(s.Opts.IPVersion, addr)
		if err != nil {
			fmt.Printf("Start listen %s err: %s\n", s.Opts.IPVersion, err)
			return
		}

		fmt.Printf("[Start] Server starts listening at IP: %s, Port: %d successful\n", s.Opts.Host, s.Opts.TcpPort)

		// TODO: 添加一个自动生成ID的方法
		var connID uint32 = 0

		// 3. 阻塞等待客户端连接，处理客户端连接业务
		for {
			// 如果有客户端连接，则阻塞会返回，往下执行
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("Start: Listener AcceptTCP: %s\n", err)
				continue
			}

			//=============
			//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= s.Opts.MaxConn {
				conn.Close()
				continue
			}
			//=============

			// 将处理当前连接的业务方法和conn进行绑定，得到我们的连接模块
			// dealConn := NewConnection(conn, connID, s.MsgHandler)
			dealConn := NewConnection(s, conn, connID, s.MsgHandler)
			connID++
			// 尝试启动连接模块
			go dealConn.Start()
		}
	}()
}

// Stop 停止
func (s *Server) Stop() {
	fmt.Printf("[STOP] server %s\n", s.Opts.Name)

	// TODO: 将服务器的资源、状态或者已经建立的连接等等进行停止或回收

	s.ConnMgr.ClearConn()
}

// Serve 运行
func (s *Server) Serve() {
	// 启动server服务功能
	s.Start()

	// TODO: 可以做一些服务器启动之后的额外业务

	// 阻塞状态
	select {}
}

// AddRouter 添加路由
func (s *Server) AddRouter(msgId uint32, router IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router Successfully!")
}

//得到链接管理
func (s *Server) GetConnMgr() IConnManager {
	return s.ConnMgr
}

//设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func (IConnection)) {
	s.OnConnStart = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func (IConnection)) {
	s.OnConnStop = hookFunc
}

//调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}

// 获取拷贝的配置
func (s *Server) Option() Option {
	return *(s.Opts)
}

