package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
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

	//当前Server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
}

//NewConnection 实例化一个链接
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		MsgHandler:   msgHandler,
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
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf error ", err)
		//	continue
		//}

		//进行拆包解包操作
		dp := DataPack{}
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			fmt.Println("read msg head error!", err)
			break
		}

		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error!", err)
			break
		}

		//再次读取data部分
		if msg.GetMsgLen() > 0 {
			data := make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read data error!", err)
				break
			}
			msg.SetData(data)
		}
		req := Request{
			conn: c,
			msg:  msg,
		}

		//从路由中找到注册绑定的Conn对应的router调用
		//执行注册的路由方法
		go c.MsgHandler.DoMsgHandler(&req)
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

//SendMsg 封包，发送数据，将数据发送给客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}
	msg := NewMsgPackage(msgId, data)
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("pack error!", err)
		return errors.New("pack error msg")
	}
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("write error!", err)
		return errors.New("conn write error")
	}
	return nil
}
