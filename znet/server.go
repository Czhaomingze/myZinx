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

			// 3.2 TODO Server.Start() 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			// 3.3 TODO Server.Start() 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			//echo函数,最大512字节
			// 将处理新连接的业务方法和conn进行绑定，得到我们定义的连接模块

			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

//Stop 停止服务器
func (s *Server) Stop() {
	//TODO: 将一些服务器的资源状态停止回收

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

//NewServer 实例化服务器
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
	}
	return s
}
