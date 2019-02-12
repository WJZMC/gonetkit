package msg

import (
	"fmt"
	"gonetkit/interfacer"
	"gonetkit/util"
	"log"
)

type MsgMgr struct {
	PoorSize  uint32
	TaskQueue []chan interfacer.Requester
	Routers   map[uint32]interfacer.Routerer
}

func NewMsgMgr() interfacer.MsgManager {
	return &MsgMgr{
		PoorSize:  util.GBConfig.WorkPoolSize,
		TaskQueue: make([]chan interfacer.Requester, util.GBConfig.WorkPoolSize),
		Routers:   make(map[uint32]interfacer.Routerer),
	}
}

//把消息发送至消息队列
func (m *MsgMgr) DeliverToMsgQueue(request interfacer.Requester) {
	index := request.GetMsgId() % m.PoorSize
	m.TaskQueue[index] <- request
}

//马上以非阻塞方式处理消息
func (m *MsgMgr) DoMsg(request interfacer.Requester) {
	go m.do(request)
}

//为消息添加具体的处理逻辑
func (m *MsgMgr) AddRouter(routerTypeName uint32, router interfacer.Routerer) {
	m.Routers[routerTypeName] = router
}

//开启worker，循环处理消息
func (m *MsgMgr) StartWorker(poolSize uint32) {
	for i := 0; i < int(m.PoorSize); i++ {
		m.TaskQueue[i] = make(chan interfacer.Requester, util.GBConfig.WorkGoChanCaps)
		go func(index int, task chan interfacer.Requester) {
			//todo 时间轮
			//case 中执行延时操作
			for {
				select {
				case request := <-task:
					m.do(request)

				}
			}
		}(i, m.TaskQueue[i])
	}
}

func (m *MsgMgr) do(request interfacer.Requester) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	router, ok := m.Routers[request.GetMsgId()]
	if !ok {
		log.Panic("路由不存在，回调失败")
	}

	router.Handle(request)

}
