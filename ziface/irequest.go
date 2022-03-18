package ziface

//IRequest 接口：实际上是把客户端请求的链接信息和请求的数据包装到了一个request中
type IRequest interface {
	//GetConnection 获取当前请求的链接
	GetConnection() IConnection

	GetData() []byte

	GetMsgID() uint32
}
