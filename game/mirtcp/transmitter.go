package gametcp

import (
	"io"
	"net"

	"github.com/davyxu/cellnet"
)

type ClientTCPMessageTransmitter struct{}

type ServerTCPMessageTransmitter struct{}

type socketOpt interface {
	MaxPacketSize() int
	ApplySocketReadTimeout(conn net.Conn, callback func())
	ApplySocketWriteTimeout(conn net.Conn, callback func())
}

func (ClientTCPMessageTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	reader, ok := ses.Raw().(io.Reader)

	// 转换错误，或者连接已经关闭时退出
	if !ok || reader == nil {
		return nil, nil
	}

	opt := ses.Peer().(socketOpt)

	if conn, ok := reader.(net.Conn); ok {

		// 有读超时时，设置超时
		opt.ApplySocketReadTimeout(conn, func() {

			msg, err = ClientRecvLTVPacket(reader, opt.MaxPacketSize())

		})
	}

	return
}

func (ClientTCPMessageTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) (err error) {

	writer, ok := ses.Raw().(io.Writer)

	// 转换错误，或者连接已经关闭时退出
	if !ok || writer == nil {
		return nil
	}

	opt := ses.Peer().(socketOpt)

	// 有写超时时，设置超时
	opt.ApplySocketWriteTimeout(writer.(net.Conn), func() {

		err = ClientSendLTVPacket(writer, ses.(cellnet.ContextSet), msg)

	})

	return
}

func (ServerTCPMessageTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	reader, ok := ses.Raw().(io.Reader)

	// 转换错误，或者连接已经关闭时退出
	if !ok || reader == nil {
		return nil, nil
	}

	opt := ses.Peer().(socketOpt)

	if conn, ok := reader.(net.Conn); ok {

		// 有读超时时，设置超时
		opt.ApplySocketReadTimeout(conn, func() {

			msg, err = ServerRecvLTVPacket(reader, opt.MaxPacketSize())

		})
	}

	return
}

func (ServerTCPMessageTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) (err error) {

	writer, ok := ses.Raw().(io.Writer)

	// 转换错误，或者连接已经关闭时退出
	if !ok || writer == nil {
		return nil
	}

	opt := ses.Peer().(socketOpt)

	// 有写超时时，设置超时
	opt.ApplySocketWriteTimeout(writer.(net.Conn), func() {

		err = ServerSendLTVPacket(writer, ses.(cellnet.ContextSet), msg)

	})

	return
}
