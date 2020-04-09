package mircodec

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/yenkeia/mirgo/game/proto/server"
)

func init() {
	codec.RegisterCodec(new(MirSplitItemCodec))
}

// MirSplitItemCodec
type MirSplitItemCodec struct{}

// Name 编码器的名字
func (*MirSplitItemCodec) Name() string {
	return "MirSplitItemCodec"
}

// MimeType 兼容http类型
func (*MirSplitItemCodec) MimeType() string {
	return "application/binary"
}

// Encode 将数据转换为字节数组
func (*MirSplitItemCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {
	var bytes []byte
	msg := msgObj.(*server.SplitItem)
	writer := &BytesWrapper{Bytes: &bytes}
	if msg.Item != nil {
		writer.Write(true)
		writer.Write(msg.Item)
	} else {
		writer.Write(false)
	}
	writer.Write(msg.Grid)
	return *writer.Bytes, nil
}

// Decode 将字节数组转换为数据
func (*MirSplitItemCodec) Decode(data interface{}, msgObj interface{}) error {
	return nil
}
