package util

import "fmt"

const (
	MaxUint32   = 4294967295
	MaxChanSize = 512
)

type UUIDGenerator struct {
	Prefix       string
	Gid          uint32
	InternolChan chan uint32
}

func (u *UUIDGenerator) start() {
	go func() {
		for {
			if u.Gid == MaxUint32 {
				u.Gid = 1
			} else {
				u.Gid++
			}
			u.InternolChan <- u.Gid
		}
	}()
}

func (u *UUIDGenerator) Get() string {
	return fmt.Sprintf("%v%v", u.Prefix, <-u.InternolChan)
}
func (u *UUIDGenerator)GetGid() uint32 {
	return <-u.InternolChan
}
func NewUUIDGenerator(pre string) *UUIDGenerator {
	gen := &UUIDGenerator{
		pre,
		0,
		make(chan uint32, MaxChanSize),
	}
	gen.start()
	return gen
}
