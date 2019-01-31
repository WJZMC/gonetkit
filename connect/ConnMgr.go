package connect

import (
	"gonetkit/interfacer"
	"sync"
	"errors"
)

type ConnMgr struct {
	connections map[uint32]interfacer.Conner
	connectionMutex sync.RWMutex
}

func NewConnMgr() interfacer.ConnManager  {
	connManager:=&ConnMgr{
		connections:make(map[uint32]interfacer.Conner),
	}
	return connManager
}
//添加连接
func (c *ConnMgr)Add(conner interfacer.Conner){
	c.connectionMutex.Lock()
	defer c.connectionMutex.Unlock()
	c.connections[conner.GetSessionId()]=conner
}

//移除连接
func (c *ConnMgr)Remove(conner interfacer.Conner){
	c.connectionMutex.Lock()
	defer c.connectionMutex.Unlock()
	delete(c.connections,conner.GetSessionId())
}

//利用sessionId获取连接
func (c *ConnMgr)Get(sessionId uint32) (interfacer.Conner, error){
	c.connectionMutex.RLock()
	defer c.connectionMutex.RUnlock()
	conn,ok:=c.connections[sessionId]
	if ok{
		return conn,nil
	}
	return nil,errors.New("connection not exist!!!")
}
//获取所有网络连接的个数
func (c *ConnMgr)Len() uint32{
	c.connectionMutex.RLock()
	defer c.connectionMutex.RUnlock()
	return uint32(len(c.connections))
}