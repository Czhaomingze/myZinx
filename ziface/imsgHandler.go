package ziface

type IMsgHandle interface {
	//DoMsgHandler 调度执行对应Router的处理方法
	DoMsgHandler(IRequest)
	//AddRouter 为消息添加具体的处理逻辑
	AddRouter(msgID uint32,router IRouter)
}
