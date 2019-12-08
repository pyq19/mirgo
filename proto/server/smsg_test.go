package server

import (
	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/codec/mircodec"
	"reflect"
	"testing"
)

func TestConst(t *testing.T) {
	t.Log(CONNECTED)
	t.Log(CLIENT_VERSION)
	t.Log(reflect.TypeOf(CONNECTED))
}

func TestMapInformation(t *testing.T) {
	//[2019/12/8 16:14:38]: --->发送服务端包信息: ServerPackets.MapInformation
	//[2019/12/8 16:14:38]: --->发送服务端包字节信息: 30, 0, 17, 0, 1, 48, 14, 66, 105, 99, 104, 111, 110, 80, 114, 111, 118, 105, 110, 99, 101, 101, 0, 135, 0, 0, 0, 0, 0, 0,

	bytes := []byte{
		1, 48,
		14, 66, 105, 99, 104, 111, 110, 80, 114, 111, 118, 105, 110, 99, 101,
		101, 0,
		135, 0,
		0, 0,
		0,
		0,
		0}
	t.Log(bytes)
	msg := new(MapInformation)

	codec := new(mircodec.MirCodec)
	if err := codec.Decode(bytes, msg); err != nil {
		panic(err)
	}
	t.Log(msg)

	cs := new(cellnet.ContextSet)
	en, err := codec.Encode(msg, *cs)
	if err != nil {
		panic(err)
	}
	t.Log(en)
	//[1 48 14 66 105 99 104 111 110 80 114 111 118 105 110 99 101 101 0 135 0 0 0 0 0 0]
}
