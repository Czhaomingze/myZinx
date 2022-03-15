package znet

import "zinx/ziface"

// BaseRouter 实现router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类进行重写
type BaseRouter struct {

}

// BaseRouter的方法都为空
// 是因为有的Router不希望有PreHandle、PostHandle这两个业务
// 所以Router全部继承BaseRouter的好处就是 无需实现PreHandle PostHandle

func (b *BaseRouter) PreHandle(request ziface.IRequest) {}

func (b *BaseRouter) Handle(request ziface.IRequest) {}

func (b *BaseRouter) PostHandle(request ziface.IRequest) {}