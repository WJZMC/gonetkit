package util

import (
	"encoding/json"
	"fmt"
	"gonetkit/interfacer"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

const (
	RouterTypeNone = iota
)

type Global struct {
	//server
	Server     interfacer.Servicer
	ServerName string
	IPVersion  string //tcp ,tcp4,tcp6
	Host       string
	Port       uint32
	MaxConn    uint32 //最大连接数

	//连接建立时的回调函数
	OnConnectioned func(fconn interfacer.Conner)
	//连接断开时的回调函数
	OnClosed func(fconn interfacer.Conner)
	//服务器停止时的回调函数
	OnServerStop func() //服务器停服回调

	//msg handel
	//处理消息的go程总数
	WorkPoolSize uint32
	//每个go程对应的任务队列容量
	WorkGoChanCaps uint32
	//每个连接队列缓存的最大消息数量
	WorkMsgChanCaps uint32
	//数据包最大大小
	MaxMsgPackSize uint32

	//发送消息的时间间隔 100/h, 100/m, 100/s
	Frequency string

	//log
	//todo

	//todo
	//时间轮

}

func (g *Global) FrequencyFormat() (int, string) {
	tmp := strings.Split(g.Frequency, "/")
	if len(tmp) == 2 {
		num, err := strconv.Atoi(tmp[0])
		if err != nil {
			return 0, ""
		}
		return num, tmp[1]
	}
	return 0, ""
}

func (g *Global) reload() {
	config, err := ioutil.ReadFile("conf/conf.json")
	if err != nil {
		fmt.Println("conf file not exist!!!")
		return
	}

	err = json.Unmarshal(config, &GBConfig)
	if err != nil {
		log.Panicln("config file err:", err)
	}
}

var GBConfig *Global

func init() {
	GBConfig = &Global{
		ServerName: "gonetkit",
		Host:       "0.0.0.0",
		Port:       9000,
		MaxConn:    12000,

		WorkPoolSize:    8,
		WorkGoChanCaps:  100,
		WorkMsgChanCaps: 1000,

		//发送消息的时间间隔 100/h, 100/m, 100/s
		Frequency: "100/s",
	}
	GBConfig.reload()
	fmt.Println(GBConfig)
}
