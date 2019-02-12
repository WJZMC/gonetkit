package interfacer

import (
	"time"
)

type Servicer interface {
	Start()
	Stop()
	Serve()

	GetConnectionMgr() ConnManager
	GetMsgHandler() MsgManager

	GetConnectionQueue() chan interface{}

	AddRouter(name uint32, router Routerer)
	CallLater(duration time.Duration, f func(args ...interface{}), args ...interface{})
	CallWhen(ts string, f func(args ...interface{}), args ...interface{})
	CallLoop(duration time.Duration, f func(args ...interface{}), args ...interface{})
}
