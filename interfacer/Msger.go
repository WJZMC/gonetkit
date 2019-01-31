package interfacer

type Msger interface {
	SetMsgId(uint32)
	SetMsgLen(uint32)
	SetMSgData([]byte)

	GetLen() uint32		//获取消息数据段长度
	GetMsgId() uint32	//获取消息id
	GetData() []byte	//获取消息内容
}