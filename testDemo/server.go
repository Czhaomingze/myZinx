package main

// 基于Zinx框架开发的服务器端应用程序
import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// PingRouter ping test 自定义路由
//模板方法模式
type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")

	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println("call back ping...ping...ping error")
	}
}

func main() {
	//创建s服务器
	s := znet.NewServer("[Zinx V0.2]")

	//给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	//启动服务
	s.Serve()
}
