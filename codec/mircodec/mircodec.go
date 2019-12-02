package mircodec

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/protoplus/proto"
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
	return proto.Marshal(msgObj)
}

func (m *MirCodec) Decode(data interface{}, msgObj interface{}) error {
	return proto.Unmarshal(data.([]byte), msgObj)
}

func init() {
	codec.RegisterCodec(new(MirCodec))
}
