package interfacer

type Routerer interface {
	//PreHandle(request Requester)//预处理
	Handle(request Requester) //处理
	//PostHandle(request Requester)//请求处理
}
