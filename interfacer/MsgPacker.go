package interfacer

//解决数据粘包
type MsgPacker interface {
	GetMsgHeadLen() uint32
	Pack(msger Msger)([]byte,error)
	UnPack([]byte)(Msger,error)
}