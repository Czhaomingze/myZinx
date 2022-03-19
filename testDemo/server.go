package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// 基于Zinx框架开发的服务器端应用程序

// PingRouter ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Handle Test
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	// 先读取客户端的数据，再回写 ping..ping..ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))
	if err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping")); err != nil {
		fmt.Println(err)
	}
}

// HelloZinxRouter hello Zinx test 自定义路由
type HelloZinxRouter struct {
	znet.BaseRouter
}

// Handle Test
func (h *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle...")
	// 先读取客户端的数据，再回写 ping..ping..ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))
	if err := request.GetConnection().SendMsg(201, []byte("Hello! Welcome to Zinx!")); err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 1 创建一个server句柄, 使用Zinx的api
	s := znet.NewServer("[Zinx V0.6]")
	// 2 给当前zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	// 3 启动server
	s.Serve()
}
