package znet

import "zinx/ziface"

type Request struct {
	//conn 已经和客户端建立好的链接
	conn ziface.IConnection

	//data 客户端请求的数据
	msg ziface.IMessage
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32{
	return r.msg.GetMsgId()
}


