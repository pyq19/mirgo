package mircodec

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

type MirCodec struct {
}

func (m *MirCodec) Name() string {
	return "MirCodec"
}

func (m *MirCodec) MimeType() string {
	return "application/binary"
}

func (m *MirCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {
	// TODO 将数据转换为字节数组
	return nil, nil
}

func (m *MirCodec) Decode(data interface{}, msgObj interface{}) error {
	// TODO 将字节数组转换为数据
	return nil
}

func init() {
	codec.RegisterCodec(new(MirCodec))
}
