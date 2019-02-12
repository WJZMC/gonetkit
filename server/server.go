package server

import (
	"fmt"
	"gonetkit/connect"
	"gonetkit/interfacer"
	"gonetkit/msg"
	"gonetkit/util"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Serve struct {
	IP      string
	Prot    uint32
	MaxConn uint32

	MsgMgr  interfacer.MsgManager
	ConnMgr interfacer.ConnManager

	GenMgr *util.UUIDGenerator

	signalChan chan os.Signal
}

func NewServe() interfacer.Servicer {
	s := &Serve{
		IP:         util.GBConfig.Host,
		Prot:       util.GBConfig.Port,
		MaxConn:    util.GBConfig.MaxConn,
		ConnMgr:    connect.NewConnMgr(),
		MsgMgr:     msg.NewMsgMgr(),
		GenMgr:     util.NewUUIDGenerator(""),
		signalChan: make(chan os.Signal),
	}
	util.GBConfig.Server = s
	return s
}
func (s *Serve) Start() {
	go func() {
		//初始化工作池
		s.MsgMgr.StartWorker(util.GBConfig.WorkPoolSize)
		//启动监听
		tcpAddr, err := net.ResolveTCPAddr(util.GBConfig.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Prot))
		if err != nil {
			log.Panic("tcp addr err:", err)
		}
		listener, err := net.ListenTCP("tcp4", tcpAddr)
		if err != nil {
			log.Panic("tcp addr err:", err)
		}

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("conn accept err:", err)
				continue
			}
			if s.ConnMgr.Len() >= util.GBConfig.MaxConn {
				conn.Write([]byte("服务器繁忙！！"))
				conn.Close()
			} else {
				go s.HandelConn(conn)
			}

		}
	}()
}
func (s *Serve) Stop() {
	fmt.Println("server stop!!")
	if util.GBConfig.OnServerStop != nil {
		util.GBConfig.OnServerStop()
	}

}
func (s *Serve) Serve() {
	s.Start()
	s.waitSignal()
}
func (s *Serve) GetConnectionMgr() interfacer.ConnManager {
	return s.ConnMgr
}
func (s *Serve) GetMsgHandler() interfacer.MsgManager {
	return s.MsgMgr
}

//todo
func (s *Serve) GetConnectionQueue() chan interface{} {
	return nil
}

func (s *Serve) AddRouter(name uint32, router interfacer.Routerer) {
	s.MsgMgr.AddRouter(name, router)

}

//todo
func (s *Serve) CallLater(duration time.Duration, f func(args ...interface{}), args ...interface{}) {

}

//todo
func (s *Serve) CallWhen(ts string, f func(args ...interface{}), args ...interface{}) {

}

//todo
func (s *Serve) CallLoop(duration time.Duration, f func(args ...interface{}), args ...interface{}) {

}

//处理请求
func (s *Serve) HandelConn(conn *net.TCPConn) {
	conn.SetNoDelay(true)
	conn.SetKeepAlive(true)
	connect := connect.NewConn(conn, s.GenMgr.GetGid(), s.MsgMgr)
	s.ConnMgr.Add(connect)
	connect.Start()
}

func (s *Serve) waitSignal() {
	signal.Notify(s.signalChan, os.Kill, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	sign := <-s.signalChan
	fmt.Println("server closed", sign)
	s.Stop()
}

//SIGHUP     终止进程     终端线路挂断
//SIGINT     终止进程     中断进程
//SIGQUIT   建立CORE文件终止进程，并且生成core文件
//SIGILL   建立CORE文件       非法指令
//SIGTRAP   建立CORE文件       跟踪自陷
//SIGBUS   建立CORE文件       总线错误
//SIGSEGV   建立CORE文件       段非法错误
//SIGFPE   建立CORE文件       浮点异常
//SIGIOT   建立CORE文件       执行I/O自陷
//SIGKILL   终止进程     杀死进程
//SIGPIPE   终止进程     向一个没有读进程的管道写数据
//SIGALARM   终止进程     计时器到时
//SIGTERM   终止进程     软件终止信号
//SIGSTOP   停止进程     非终端来的停止信号
//SIGTSTP   停止进程     终端来的停止信号
//SIGCONT   忽略信号     继续执行一个停止的进程
//SIGURG   忽略信号     I/O紧急信号
//SIGIO     忽略信号     描述符上可以进行I/O
//SIGCHLD   忽略信号     当子进程停止或退出时通知父进程
//SIGTTOU   停止进程     后台进程写终端
//SIGTTIN   停止进程     后台进程读终端
//SIGXGPU   终止进程     CPU时限超时
//SIGXFSZ   终止进程     文件长度过长
//SIGWINCH   忽略信号     窗口大小发生变化
//SIGPROF   终止进程     统计分布图用计时器到时
//SIGUSR1   终止进程     用户定义信号1
//SIGUSR2   终止进程     用户定义信号2
//SIGVTALRM 终止进程     虚拟计时器到时
