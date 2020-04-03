package mirtcp

import (
	"bytes"

	"github.com/davyxu/cellnet"
	"github.com/gorilla/websocket"
)

type WSMessageTransmitter struct {
}

func (WSMessageTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	conn, ok := ses.Raw().(*websocket.Conn)

	// 转换错误，或者连接已经关闭时退出
	if !ok || conn == nil {
		return nil, nil
	}

	var messageType int
	var raw []byte
	messageType, raw, err = conn.ReadMessage()

	if err != nil {
		return
	}

	// if len(raw) < MsgIDSize {
	// 	return nil, util.ErrMinPacket
	// }

	switch messageType {
	case websocket.BinaryMessage:
		reader := bytes.NewReader(raw)
		msg, err = ServerRecvLTVPacket(reader, 1024*1024)
	}

	return
}

func (WSMessageTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) error {

	conn, ok := ses.Raw().(*websocket.Conn)

	// 转换错误，或者连接已经关闭时退出
	if !ok || conn == nil {
		return nil
	}

	writer := &bytes.Buffer{}

	err := ServerSendLTVPacket(writer, nil, msg)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.BinaryMessage, writer.Bytes())
}
