package interfacer

type MsgManager interface {
	DeliverToMsgQueue(request Requester)          //把消息发送至消息队列
	DoMsg(request Requester)                      //马上以非阻塞方式处理消息
	AddRouter(routerType uint32, router Routerer) //为消息添加具体的处理逻辑
	StartWorker(poolSize uint32)                  //开启worker，循环处理消息
}
