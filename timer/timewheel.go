package timer

import (
	"sync"
)

type TimeWheel struct {
	Name string
	Interal int64  //时间轮间隔
	Scale int      //时间轮单位
	Curent int		//当前时间轮个数
	MaxCaps int		//最大时间轮个数
	TimerQueue map[int]map[uint32]*Timer  //当前时间轮map
	NextTimeWheel *TimeWheel	//下一个时间轮
	sync.RWMutex	//时间轮专用锁
}

func NewTimeWheel(name string,interal int64,scale ,curent,maxCaps int) *TimeWheel  {
	t:=&TimeWheel{
		Name:name,
		Interal:interal,
		Scale:scale,
		Curent:curent,
		MaxCaps:maxCaps,
		TimerQueue:make(map[int]map[uint32]*Timer),
	}
	for i:=0;i<t.MaxCaps;i++{
		t.TimerQueue[i]=make(map[uint32]*Timer,t.MaxCaps)
	}
	return t
}

func (t *TimeWheel)AddTimer(tid int64,timer *Timer) error {
	t.Lock()
	defer t.Unlock()

	return t.addTimer(tid,timer,false)
}

func (t *TimeWheel)addTimer(tid int64,timer *Timer,isWheelScroll bool) error {

	////异常捕获
	//defer func() error {
	//	if err:=recover();err!=nil{
	//		errmsg:=fmt.Sprintf("addTimer err:%v",err)
	//		return errors.New(errmsg)
	//	}
	//	return nil
	//}()
	//
	//offsetTime:=timer.tx-util.GetTimer()
	//
	//if offsetTime>t.Interal




	return nil

}
