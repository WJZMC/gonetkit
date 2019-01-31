package request

import "gonetkit/interfacer"

type Request struct {
	Conn interfacer.Conner
	Msg  interfacer.Msger
}

//获取请求连接信息
func (r *Request) GetConnection() interfacer.Conner {
	return r.Conn
}

//获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.Msg.GetData()
}

//获取请求消息的ID
func (r *Request) GetMsgId() uint32 {
	return r.Msg.GetMsgId()
}
