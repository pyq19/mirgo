package mircodec

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

// MirUserInformationCodec ...
type MirUserInformationCodec struct {}

// Name 返回名字
func (m *MirUserInformationCodec) Name() string {
	return "MirUserInformationCodec"
}

// MimeType 我也不知道是干嘛的
func (m *MirUserInformationCodec) MimeType() string {
	return "application/binary"
}

// Encode 将数据转换为字节数组
func (*MirUserInformationCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {
	return encode(msgObj)
}

// Decode 将字节数组转换为数据
func (*MirUserInformationCodec) Decode(data interface{}, msgObj interface{}) error {
	bytes := data.([]byte)
	return decode(msgObj, bytes)
}

func init() {
	codec.RegisterCodec(new(MirUserInformationCodec))
}
