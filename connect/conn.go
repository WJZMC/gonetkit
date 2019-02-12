package connect

import (
	"errors"
	"fmt"
	"gonetkit/interfacer"
	"gonetkit/msg"
	"gonetkit/request"
	"gonetkit/util"
	"io"
	"net"
	"sync"
	"time"
)

type Conn struct {
	Connect   *net.TCPConn
	isClose   bool //是否已关闭
	sessionId uint32

	dataPack  interfacer.MsgPacker
	msgHandel interfacer.MsgManager

	buffChan chan []byte //缓冲消息队列
	writeCh  chan []byte //无缓冲消息队列
	exitCh   chan bool

	//自定义连接属性
	propertys     map[string]interface{}
	propertyMutex sync.RWMutex
}

func NewConn(conn *net.TCPConn, sessionId uint32, msgHandel interfacer.MsgManager) interfacer.Conner {
	return &Conn{
		Connect:   conn,
		isClose:   false,
		sessionId: sessionId,
		dataPack:  msg.NewPack(),
		msgHandel: msgHandel,
		buffChan:  make(chan []byte, util.GBConfig.WorkMsgChanCaps),
		writeCh:   make(chan []byte),
		exitCh:    make(chan bool),
		propertys: make(map[string]interface{}),
	}
}

//设置限制访问频率
func (c *Conn) setFrequency() {

	num, suf := util.GBConfig.FrequencyFormat()
	if num == 0 {
		return
	}
	c.SetProperty("gnk_fqy_count", 0)
	c.SetProperty("gnk_fqy_num", num)
	switch suf {
	case "h":
		c.SetProperty("gnk_fqy_ts", time.Now().UnixNano()/1e6+int64(3600*1e3))
	case "m":
		c.SetProperty("gnk_fqy_ts", time.Now().UnixNano()/1e6+int64(60*1e3))
	case "s":
		c.SetProperty("gnk_fqy_ts", time.Now().UnixNano()/1e6+int64(1e3))

	default:
	}

}
func (c *Conn) doFrequency() error {
	ts, err := c.GetProperty("gnk_fqy_ts")
	if err != nil {
		return err
	}
	if ts.(int64) > time.Now().UnixNano() {

		c.setFrequency()
		return nil
	}

	coutTmp, err := c.GetProperty("gnk_fqy_count")
	if err != nil {
		return err
	}
	numTmp, err := c.GetProperty("gnk_fqy_num")
	if err != nil {
		return err
	}
	count := coutTmp.(int)
	num := numTmp.(int)
	count++
	if num < count {
		return errors.New("消息发送过于频繁！！")
	}
	c.SetProperty("gnk_fqy_count", count)
	return nil

}
func (c *Conn) startWritingGoroutine() {
	fmt.Println("startWritingGoroutine----2----", c.Connect.RemoteAddr().String())
	defer fmt.Println("startWritingGoroutine  close")
	defer c.Stop()
	for {
		select {
		case data := <-c.buffChan:
			c.Connect.Write(data)
		case data := <-c.writeCh:
			c.Connect.Write(data)
		case <-c.exitCh:
			return
		}
	}
}

func (c *Conn) startReadingGoroutine() {
	fmt.Println("startReadingGoroutine----2----", c.Connect.RemoteAddr().String())
	defer fmt.Println("startReadingGoroutine  close")
	defer c.Stop()
	for {
		headBuf := make([]byte, c.dataPack.GetMsgHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headBuf); err != nil {
			fmt.Println("startReadingGoroutine  read head:", err)
			return
		}

		msgHead, err := c.dataPack.UnPack(headBuf)
		if err != nil {
			fmt.Println("unpack msg head:", err)
			return
		}
		if msgHead.GetLen() <= 0 {
			fmt.Println("msg head==0")
			return
		}

		data := make([]byte, msgHead.GetLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
			fmt.Println("read msg body:", err)
			return
		}

		msgHead.SetMSgData(data)

		//频率控制
		if err = c.doFrequency(); err != nil {
			fmt.Println("doFrequency:", err)
			return
		}

		if util.GBConfig.WorkPoolSize > 0 {
			c.msgHandel.DeliverToMsgQueue(&request.Request{
				Conn: c,
				Msg:  msgHead,
			})
		} else {
			c.msgHandel.DoMsg(&request.Request{
				Conn: c,
				Msg:  msgHead,
			})
		}

	}
}

//启动连接
func (c *Conn) Start() {
	c.setFrequency() //读取数据时间

	go c.startWritingGoroutine()

	go c.startReadingGoroutine()

	if util.GBConfig.OnConnectioned != nil {
		util.GBConfig.OnConnectioned(c)
	}
}

//断开连接
func (c *Conn) Stop() {
	if c.isClose {
		return
	}

	c.isClose = true

	c.Connect.Close()

	close(c.exitCh)

	util.GBConfig.Server.GetConnectionMgr().Remove(c)

	close(c.buffChan)

	close(c.writeCh)
}

//获取TCP连接
func (c *Conn) GetTCPConnection() *net.TCPConn {
	return c.Connect
}

//获取session id
func (c *Conn) GetSessionId() uint32 {
	return c.sessionId
}

//获取远程机器地址
func (c *Conn) RemoteAddr() net.Addr {
	return c.Connect.RemoteAddr()
}

//直接发送数据至TCP连接对方
func (c *Conn) SendMsg(msgId uint32, data []byte) error {
	if c.isClose {
		return errors.New("send fail ,conn closed ")
	}
	if msg, err := c.dataPack.Pack(msg.NewMsg(msgId, data)); err != nil {
		return err
	} else {
		c.writeCh <- msg
	}
	return nil
}

//把数据发送至缓冲区
func (c *Conn) SendMsgBuff(msgId uint32, data []byte) error {
	if c.isClose {
		return errors.New("send fail ,conn closed ")
	}
	if msg, err := c.dataPack.Pack(msg.NewMsg(msgId, data)); err != nil {
		return err
	} else {
		//当消息队列写满的时候，超过2秒回超时
		select {
		case <-time.After(2 * time.Second):
			return errors.New("发送消息超时")
		case c.buffChan <- msg:
			return nil
		}
	}
	return nil
}

//设置连接属性
func (c *Conn) SetProperty(k string, v interface{}) {
	c.propertyMutex.Lock()
	defer c.propertyMutex.Unlock()
	c.propertys[k] = v
}

//获取连接属性
func (c *Conn) GetProperty(k string) (interface{}, error) {
	c.propertyMutex.RLock()
	defer c.propertyMutex.RUnlock()
	v, ok := c.propertys[k]
	if ok {
		return v, nil
	}

	return nil, errors.New("not exist!!!")
}

//移除连接属性
func (c *Conn) RemoveProperty(k string) {
	c.propertyMutex.Lock()
	defer c.propertyMutex.Unlock()
	delete(c.propertys, k)
}
