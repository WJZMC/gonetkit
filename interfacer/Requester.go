package interfacer


type Requester interface {
	GetConnection() Conner	//获取请求连接信息
	GetData() []byte			//获取请求消息的数据
	GetMsgId() uint32			//获取请求消息的ID
}
