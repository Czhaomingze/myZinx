package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)
// MsgHandle 消息处理模块的实现
type MsgHandle struct {
	// 存放每个MsgID对应的处理方法
	APIs map[uint32]ziface.IRouter
}

func NewMsgHandle() *MsgHandle{
	return &MsgHandle{
		APIs: make(map[uint32]ziface.IRouter),
	}
}
func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1 从request中找到msgID
	handler,ok:=m.APIs[request.GetMsgID()]
	if !ok{
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
