package msg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"gonetkit/interfacer"
	"gonetkit/util"
)

var MAX_PACKET_SIZE = 1024 * 1024

type Msg struct {
	msgId   uint32
	msgLen  uint32
	msgData []byte
}

func NewMsg(msgId uint32, msgData []byte) interfacer.Msger {
	return &Msg{
		msgId:   msgId,
		msgData: msgData,
		msgLen:  uint32(len(msgData)),
	}
}
func (m *Msg) SetMsgId(msgId uint32) {
	m.msgId = msgId
}
func (m *Msg) SetMsgLen(msgLen uint32) {
	m.msgLen = msgLen
}
func (m *Msg) SetMSgData(data []byte) {
	m.msgData = data
}

//获取消息数据段长度
func (m *Msg) GetLen() uint32 {
	return m.msgLen
}

//获取消息id
func (m *Msg) GetMsgId() uint32 {
	return m.msgId
}

//获取消息内容
func (m *Msg) GetData() []byte {
	return m.msgData
}

const MaxHeadLen = 8

type MsgPack struct {
}

func NewPack() interfacer.MsgPacker {
	return &MsgPack{}
}
func (m *MsgPack) GetMsgHeadLen() uint32 {
	return MaxHeadLen
}
func (m *MsgPack) Pack(msger interfacer.Msger) ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})

	err := binary.Write(buff, binary.LittleEndian, msger.GetLen())
	if err != nil {
		return nil, err
	}

	err = binary.Write(buff, binary.LittleEndian, msger.GetMsgId())
	if err != nil {
		return nil, err
	}

	err = binary.Write(buff, binary.LittleEndian, msger.GetData())
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil

}

func (m *MsgPack) UnPack(data []byte) (interfacer.Msger, error) {
	buffRead := bytes.NewReader(data)

	var msg *Msg= new(Msg)
	//向固定字节的空间解压缩，写完为止，这里只写4字节
	if err := binary.Read(buffRead, binary.LittleEndian, &msg.msgLen); err != nil {
		return nil, err
	}

	if err := binary.Read(buffRead, binary.LittleEndian, &msg.msgId); err != nil {
		return nil, err
	}

	if util.GBConfig.MaxMsgPackSize > 0 && msg.GetLen() > util.GBConfig.MaxMsgPackSize || util.GBConfig.MaxMsgPackSize == 0 && msg.GetLen() > uint32(MAX_PACKET_SIZE) {
		return nil, errors.New("package lager than maxSize!!!")
	}
	return msg, nil
}
