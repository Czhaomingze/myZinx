package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

// MsgHandle 消息处理模块的实现
type MsgHandle struct {
	// 存放每个MsgID对应的处理方法
	APIs map[uint32]ziface.IRouter
	//负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作worker池的worker数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		APIs:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}
func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1 从request中找到msgID
	handler, ok := m.APIs[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is NOT FOUND! Need Register!")
		return
	}
	// 2 根据MsgID调度router对应的业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := m.APIs[msgID]; ok {
		panic("repeat api ,msgID=" + strconv.Itoa(int(msgID)))
	}
	// 2 添加msg和API的绑定关系
	m.APIs[msgID] = router
	fmt.Println("Add api MsgID", msgID, "success!")
}

//StartWorkerPool 启动一个worker工作池
func (m *MsgHandle) StartWorkerPool() {
	//根据workerPoolSize，分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 一个 worker 被启动
		// 1 当前的 worker 对应的 channel 消息队列，开辟空间 第0个 worker 就用第0个channel...
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2 启动当前的 worker, 阻塞等待消息从 channel 传递进来
		go m.startOneWorker(i, m.TaskQueue[i])
	}
}

//startOneWorker 启动一个worker工作流程
func (m *MsgHandle) startOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started ...")
	// 不断的阻塞等待对应的消息队列的消息
	for {
		select {
		// 如果有消息过来，出列的就是一个客户端的request，执行当前request绑定的业务
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

func (m *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1 将消息分配给不同的 worker
	// 根据客户端建立的 ConnID 来分配
	workerId := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		", request MsgId = ", request.GetMsgID(), " to WorkerID = ", workerId)
	// 2 将消息发送给对应的 worker 的 TaskQueue
	m.TaskQueue[workerId] <- request
}
