package interfacer
type ConnManager interface {
	Add(conner Conner)					//添加连接
	Remove(conner Conner)				//移除连接
	Get(sessionId uint32) (Conner, error)	//利用sessionId获取连接
	Len() uint32									//获取所有网络连接的个数
}
