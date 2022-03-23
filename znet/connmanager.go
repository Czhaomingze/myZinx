package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	// 管理的连接集合
	connections map[uint32]ziface.IConnection
	// 保护连接集合的读写锁
	connLock sync.RWMutex
}

// NewConnManager 创建连接
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	cm.connections[conn.GetConnID()]=conn

	fmt.Println("connID = ", conn.GetConnID(), " add to ConnManager successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections,conn.GetConnID())

	fmt.Println("connID = ", conn.GetConnID(), " remove from ConnManager successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND!")
	}
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	// 保护共享资源 map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除 conn 并停止 conn 的工作
	for connID, conn := range cm.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.connections, connID)
	}

	fmt.Println("Clear All connections succ! conn num = ", cm.Len())
}



