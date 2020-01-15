package mircodec

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/yenkeia/mirgo/proto/server"
)

func init() {
	codec.RegisterCodec(new(MirObjectPlayerCodec))
}

/*
MirPlayerInspectCodec
*/
type MirObjectPlayerCodec struct{}

// Name 编码器的名字
func (*MirObjectPlayerCodec) Name() string {
	return "MirObjectPlayerCodec"
}

// MimeType 兼容http类型
func (*MirObjectPlayerCodec) MimeType() string {
	return "application/binary"
}

// Encode 将数据转换为字节数组
func (*MirObjectPlayerCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {
	var bytes []byte
	pi := msgObj.(*server.ObjectPlayer)
	writer := &BytesWrapper{Bytes: &bytes}
	return *writer.Bytes, nil
}

// Decode 将字节数组转换为数据
func (*MirObjectPlayerCodec) Decode(data interface{}, msgObj interface{}) error {
	pi := msgObj.(*server.ObjectPlayer)
	bytes := data.([]byte)
	reader := &BytesWrapper{Bytes: &bytes}
	return nil
}
