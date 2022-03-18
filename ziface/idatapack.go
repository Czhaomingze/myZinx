package ziface

//拆包 封包模块
//直接面向TCP链接中的数据流，用于处理TCP粘包问题

type IDataPack interface {
	// GetHeadLen 获取包的头的长度
	GetHeadLen() uint32
	// Pack 封包
	Pack(message IMessage) ([]byte,error)
	// Unpack 拆包
	Unpack([]byte) (IMessage,error)
}
