package timer

import (
	"time"
)

type Timer struct {
	f  *DelayFunc
	tx int64
}

func NewTimerAt(fun func(...interface{}), args []interface{}, tx int64) *Timer {
	t := &Timer{
		NewDelayFunc(fun, args),
		int64(tx),
	}
	t.run()
	return t
}
func NewTimerAfter(fun func(...interface{}), args []interface{}, tx time.Duration) *Timer {
	return NewTimerAt(fun, args, time.Now().UnixNano()+int64(tx))
}

func (t *Timer) run() {
	go func() {
		now := time.Now().UnixNano()
		if now > t.tx {
			time.Sleep(time.Duration(t.tx - now))
		}
		t.f.Call()
	}()
}
