package ziface

import "net"

//定义链接模块的抽象层

type IConnection interface {
	//Start 启动链接，让当前的链接准备开始工作
	Start()

	//Stop 停止链接，结束当前链接的工作
	Stop()

	//GetTCPConnection 获取当前连接绑定的socket connect
	GetTCPConnection() *net.TCPConn

	//GetConnID 获取当前模块的链接id
	GetConnID() uint32

	//RemoteAddr 获取客户端的TCP状态 Ip sort
	RemoteAddr() net.Addr

	//Send 发送数据，将数据发送给客户端
	SendMsg(msgId uint32,data []byte) error
}


//HandleFunc 定义一个处理链接业务的方法，回调函数
type HandleFunc func(*net.TCPConn,[]byte,int) error