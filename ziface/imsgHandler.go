package ziface

type IMsgHandle interface {
	//DoMsgHandler 调度执行对应Router的处理方法
	DoMsgHandler(IRequest)
	//AddRouter 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
	//StartWorkerPool 启动一个worker工作池
	StartWorkerPool()
	//SendMsgToTaskQueue 将消息发送给消息任务队列处理
	SendMsgToTaskQueue(request IRequest)
}
