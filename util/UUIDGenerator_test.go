package util

import (
	"testing"
	"fmt"
)

func Test_NewUUIDGenerator(t *testing.T)  {
	UUID:=NewUUIDGenerator("jack")

	for i:=0;i<10 ;i++{
		fmt.Println(UUID.Get())
	}


}