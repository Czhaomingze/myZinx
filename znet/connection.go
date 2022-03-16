package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

//链接模块

type Connection struct {
	//socket TCP 套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前链接的状态（是否已经关闭）
	isClosed bool

	//等待链接被动退出的channel
	ExitChan chan bool

	//该链接处理的方法router
	Router ziface.IRouter
}

//NewConnection 实例化一个链接
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}
	return c
}

//StartReader 启动读数据业务
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, ", Reader is exit, remote addr is ", c.RemoteAddr().String()) // 2
	defer c.Stop()                                                                                         // 1

	for {
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf error ", err)
			continue
		}

		req := Request{
			conn: c,
			data: buf,
		}

		//从路由中找到注册绑定的Conn对应的router调用
		//执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

//Start 启动链接，让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn start().. ConnID=", c.ConnID)
	// 启动从当前连接读数据的业务
	go c.StartReader()

	//TODO 启动从当前连接写数据的业务
	//go c.StartWriter()
}

//Stop 停止链接，结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn stop().. ConnID=", c.ConnID)

	if c.isClosed {
		return
	}

	c.isClosed = true
	//关闭链接
	c.Conn.Close()
	//关闭管道
	close(c.ExitChan)
}

//GetTCPConnection 获取当前连接绑定的socket connect
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//GetConnID 获取当前模块的链接id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//RemoteAddr 获取客户端的TCP状态 Ip sort
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//Send 发送数据，将数据发送给客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
