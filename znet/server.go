package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)
import "net"

//Server 实现IServer接口，定义一个Server服务器模块
type Server struct {
	//Name 服务器名称
	Name string
	//IPVersion 服务器绑定的IP版本
	IPVersion string
	//IP 服务器监听的IP地址
	IP string
	//Port 服务器监听的端口
	Port int
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
	// 连接管理器
	ConnMgr ziface.IConnManager
	// 创建连接后自动调用的 Hook 函数 OnConnStart
	OnConnStart func(conn ziface.IConnection)
	// 销毁连接前自动调用的 Hook 函数 OnConnStop
	OnConnStop func(conn ziface.IConnection)
}

//Start 服务器启动
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, listenner at IP: %s, Port: %d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)

	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	fmt.Printf("[START] Server Listener at IP :%s, Port %d, is starting\n", s.IP, s.Port)
	go func() {
		// 0 开启一个worker工作池
		s.MsgHandler.StartWorkerPool()
		// 1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}
		// 2 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}

		fmt.Println("start Zinx server, ", s.Name, "succ, Listening...")

		var cid uint32
		cid = 0
		// 3 阻塞的等待客户端连接，处理客户端连接业务（读写）
		for {
			// 3.1 阻塞等待客户端建立连接请求
			conn, error := listener.AcceptTCP()

			if error != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			fmt.Println(">>>>>>ConnMgr Len = ", s.ConnMgr.Len(), ", MAX = ", utils.GlobalObject.MaxConn)
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println("====> Too Many Connections MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// 3.3 TODO Server.Start() 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			// 将处理新连接的业务方法和conn进行绑定，得到我们定义的连接模块

			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

//Stop 停止服务器
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server, name ", s.Name)
	// 将一些服务器的资源、状态、或者 已经开辟的连接信息进行停止或回收
	s.ConnMgr.ClearConn()
}

//Serve 服务器进行服务
func (s *Server) Serve() {
	//启动服务器
	s.Start()

	//TODO: 做一些额外业务

	//阻塞状态
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router success!")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

//NewServer 实例化服务器
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// Hook 函数的设置与调用

func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}
func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("-----> Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("-----> Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}
