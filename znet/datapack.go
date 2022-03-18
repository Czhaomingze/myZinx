package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

//拆包、封包的具体模块

type DataPack struct{}

//NewDataPack 实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

//GetHeadLen 获取包的头部长度方法
func (d DataPack) GetHeadLen() uint32 {
	// DataLen uint32（4字节）+ ID uint32（4字节）= 8字节
	return 8
}

func (d DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//将msgLen写进dataBuff中
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//将msgId写进dataBuff中
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//将data写进dataBuff中
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (d DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	dataBuff:=bytes.NewReader(binaryData)

	msg:=&Message{}

	// 读dataLen
	if err := binary.Read(dataBuff, binary.BigEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 读msgId
	if err := binary.Read(dataBuff, binary.BigEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen是否已经超出允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too Large msg data recv!")
	}

	return msg, nil
}
