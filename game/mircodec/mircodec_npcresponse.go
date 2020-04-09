package mircodec

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/yenkeia/mirgo/game/proto/server"
)

func init() {
	codec.RegisterCodec(new(MirNPCResponseCodec))
}

/*
MirNPCResponseCodec
*/
type MirNPCResponseCodec struct{}

// Name 编码器的名字
func (*MirNPCResponseCodec) Name() string {
	return "MirNPCResponseCodec"
}

// MimeType 兼容http类型
func (*MirNPCResponseCodec) MimeType() string {
	return "application/binary"
}

// Encode 将数据转换为字节数组
func (*MirNPCResponseCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {
	var bytes []byte
	res := msgObj.(*server.NPCResponse)
	writer := &BytesWrapper{Bytes: &bytes}
	count := len(res.Page)
	writer.Write(count)
	for i := 0; i < count; i++ {
		writer.Write(res.Page[i])
	}
	return *writer.Bytes, nil
}

// Decode 将字节数组转换为数据
func (*MirNPCResponseCodec) Decode(data interface{}, msgObj interface{}) error {
	res := msgObj.(*server.NPCResponse)
	bytes := data.([]byte)
	reader := &BytesWrapper{Bytes: &bytes}
	count := int(reader.ReadInt32())
	res.Page = make([]string, count)
	for i := 0; i < count; i++ {
		res.Page[i] = reader.ReadString()
	}
	return nil
}
