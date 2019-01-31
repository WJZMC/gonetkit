package timer

import (
	"fmt"
	"reflect"
)

type DelayFunc struct {
	f    func(...interface{})
	args []interface{}
}

func NewDelayFunc(f func(...interface{}), args []interface{}) *DelayFunc {
	return &DelayFunc{
		f,
		args,
	}
}

func (delay *DelayFunc) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(delay.Call))
}
func (delay *DelayFunc) Call() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("delay err:", err)
		}
	}()

	delay.f(delay.args...)
}
