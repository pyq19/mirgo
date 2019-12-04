package mircodec

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

// MirCodec 编码解码
type MirCodec struct {
}

// Name 返回名字
func (m *MirCodec) Name() string {
	return "MirCodec"
}

// MimeType 我也不知道是干嘛的
func (m *MirCodec) MimeType() string {
	return "application/binary"
}

// Encode 将数据转换为字节数组
func (m *MirCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {
	return nil, nil
}

// Decode 将字节数组转换为数据
func (m *MirCodec) Decode(data interface{}, msgObj interface{}) error {
	return nil
}

func init() {
	codec.RegisterCodec(new(MirCodec))
}
