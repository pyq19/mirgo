package mirtcp

import (
	"encoding/binary"
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/golog"
	"io"
	"strconv"
	"strings"
)

var (
	ErrMaxPacket  = errors.New("packet over size")
	ErrMinPacket  = errors.New("packet short size")
	ErrShortMsgID = errors.New("short msgid")
)

const (
	bodySize  = 2 // 包体大小字段
	msgIDSize = 2 // 消息ID字段
)

var log = golog.New("proc.mirtcp.util")

func String(bytes []byte) string {
	strSlice := []string{}
	for _, b := range bytes {
		strSlice = append(strSlice, strconv.Itoa(int(b)))
	}
	res := strings.Join(strSlice, ", ")
	return "[" + res + "]"
}

// 接收Length-Type-Value格式的封包流程
func ClientRecvLTVPacket(reader io.Reader, maxPacketSize int) (msg interface{}, err error) {

	// Size为uint16，占2字节
	var sizeBuffer = make([]byte, bodySize)

	// 持续读取Size直到读到为止
	_, err = io.ReadFull(reader, sizeBuffer)

	// 发生错误时返回
	if err != nil {
		return
	}

	if len(sizeBuffer) < bodySize {
		return nil, ErrMinPacket
	}

	// 用小端格式读取Size
	size := binary.LittleEndian.Uint16(sizeBuffer)

	if maxPacketSize > 0 && size >= uint16(maxPacketSize) {
		return nil, ErrMaxPacket
	}

	// 分配包体大小
	body := make([]byte, size)

	// 读取包体数据
	_, err = io.ReadFull(reader, body)

	allBytes := append(sizeBuffer, body...)
	log.Debugln("<--- 客户端收到 Length: " + strconv.Itoa(len(allBytes)) + ", AllBytes: " + String(allBytes))

	// 发生错误时返回
	if err != nil {
		return
	}

	if len(body) < msgIDSize {
		return nil, ErrShortMsgID
	}

	msgid := binary.LittleEndian.Uint16(body)

	// FIXME 客户端接收到的是服务器的包, ID + 2000
	msgid = msgid + 2000

	msgData := body[msgIDSize:]

	// 将字节数组和消息ID用户解出消息
	msg, _, err = codec.DecodeMessage(int(msgid), msgData)
	if err != nil {
		// TODO 接收错误时，返回消息
		return nil, err
	}

	return
}

// 发送Length-Type-Value格式的封包流程
func ClientSendLTVPacket(writer io.Writer, ctx cellnet.ContextSet, data interface{}) error {

	var (
		msgData []byte
		msgID   int
		meta    *cellnet.MessageMeta
	)

	switch m := data.(type) {
	case *cellnet.RawPacket: // 发裸包
		msgData = m.MsgData
		msgID = m.MsgID
	default: // 发普通编码包
		var err error

		// 将用户数据转换为字节数组和消息ID
		msgData, meta, err = codec.EncodeMessage(data, ctx)

		if err != nil {
			return err
		}

		msgID = meta.ID
	}

	pkt := make([]byte, bodySize+msgIDSize+len(msgData))

	// Length
	binary.LittleEndian.PutUint16(pkt, uint16(msgIDSize+len(msgData)))

	// FIXME 客户端发送的是客户端的包, ID - 1000
	msgID = msgID - 1000

	// Type
	binary.LittleEndian.PutUint16(pkt[bodySize:], uint16(msgID))

	// Value
	copy(pkt[bodySize+msgIDSize:], msgData)

	// 将数据写入Socket
	err := WriteFull(writer, pkt)

	log.Debugln("---> 客户端发送: " + String(pkt))

	// Codec中使用内存池时的释放位置
	if meta != nil {
		codec.FreeCodecResource(meta.Codec, msgData, ctx)
	}

	return err
}

// 接收Length-Type-Value格式的封包流程
func ServerRecvLTVPacket(reader io.Reader, maxPacketSize int) (msg interface{}, err error) {

	// Size为uint16，占2字节
	var sizeBuffer = make([]byte, bodySize)

	// 持续读取Size直到读到为止
	_, err = io.ReadFull(reader, sizeBuffer)

	// 发生错误时返回
	if err != nil {
		return
	}

	if len(sizeBuffer) < bodySize {
		return nil, ErrMinPacket
	}

	// 用小端格式读取Size
	size := binary.LittleEndian.Uint16(sizeBuffer)

	if maxPacketSize > 0 && size >= uint16(maxPacketSize) {
		return nil, ErrMaxPacket
	}

	// 分配包体大小
	body := make([]byte, size)

	// 读取包体数据
	_, err = io.ReadFull(reader, body)

	allBytes := append(sizeBuffer, body...)
	log.Debugln("<--- 服务端收到 Length: " + strconv.Itoa(len(allBytes)) + ", AllBytes: " + String(allBytes))

	// 发生错误时返回
	if err != nil {
		return
	}

	if len(body) < msgIDSize {
		return nil, ErrShortMsgID
	}

	msgid := binary.LittleEndian.Uint16(body)

	// FIXME 服务端接收到的是客户端的包, ID + 1000
	msgid = msgid + 1000

	msgData := body[msgIDSize:]

	// 将字节数组和消息ID用户解出消息
	msg, _, err = codec.DecodeMessage(int(msgid), msgData)
	if err != nil {
		// TODO 接收错误时，返回消息
		return nil, err
	}

	return
}

// 发送Length-Type-Value格式的封包流程
func ServerSendLTVPacket(writer io.Writer, ctx cellnet.ContextSet, data interface{}) error {

	var (
		msgData []byte
		msgID   int
		meta    *cellnet.MessageMeta
	)

	switch m := data.(type) {
	case *cellnet.RawPacket: // 发裸包
		msgData = m.MsgData
		msgID = m.MsgID
	default: // 发普通编码包
		var err error

		// 将用户数据转换为字节数组和消息ID
		msgData, meta, err = codec.EncodeMessage(data, ctx)

		if err != nil {
			return err
		}

		msgID = meta.ID
	}

	pkt := make([]byte, bodySize+msgIDSize+len(msgData))

	// Length
	binary.LittleEndian.PutUint16(pkt, uint16(msgIDSize+len(msgData)))

	// FIXME 服务端发送的是服务端的包, ID - 2000
	msgID = msgID - 2000

	// Type
	binary.LittleEndian.PutUint16(pkt[bodySize:], uint16(msgID))

	// Value
	copy(pkt[bodySize+msgIDSize:], msgData)

	// 将数据写入Socket
	err := WriteFull(writer, pkt)

	log.Debugln("---> 服务端发送: " + String(pkt))

	// Codec中使用内存池时的释放位置
	if meta != nil {
		codec.FreeCodecResource(meta.Codec, msgData, ctx)
	}

	return err
}

func WriteFull(writer io.Writer, buf []byte) error {

	total := len(buf)

	for pos := 0; pos < total; {

		n, err := writer.Write(buf[pos:])

		if err != nil {
			return err
		}

		pos += n
	}

	return nil

}
