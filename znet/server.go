package znet

import (
	"fmt"
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
	//当前的Server添加一个router，server注册的链接对应的处理业务
	Router ziface.IRouter
}

//Start 服务器启动
func (s *Server) Start() {
	fmt.Printf("[START] Server Listener at IP :%s, Port %d, is starting\n", s.IP, s.Port)
	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}

		fmt.Println("start Zinx server, ", s.Name, "succ, Listening...")

		var cid uint32
		cid = 0
		for {
			conn, error := listener.AcceptTCP()

			if error != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//echo函数,最大512字节
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

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

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router success!")
}

//NewServer 实例化服务器
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
	}
	return s
}
