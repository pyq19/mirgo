package mircodec

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/yenkeia/mirgo/game/cm"
	"github.com/yenkeia/mirgo/game/proto/server"
)

func init() {
	codec.RegisterCodec(new(MirObjectNPCCodec))
}

/*
MirPlayerInspectCodec
*/
type MirObjectNPCCodec struct{}

// Name 编码器的名字
func (*MirObjectNPCCodec) Name() string {
	return "MirObjectNPCCodec"
}

// MimeType 兼容http类型
func (*MirObjectNPCCodec) MimeType() string {
	return "application/binary"
}

// Encode 将数据转换为字节数组
func (*MirObjectNPCCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {
	var bytes []byte
	on := msgObj.(*server.ObjectNPC)
	writer := &BytesWrapper{Bytes: &bytes}
	writer.Write(on.ObjectID)
	writer.Write(on.Name)
	writer.Write(on.NameColor)
	writer.Write(on.Image)
	writer.Write(on.Color)
	writer.Write(on.Location.X)
	writer.Write(on.Location.Y)
	writer.Write(uint8(on.Direction))
	qc := len(on.QuestIDs)
	writer.Write(qc)
	for i := range on.QuestIDs {
		writer.Write(on.QuestIDs[i])
	}
	return *writer.Bytes, nil
}

// Decode 将字节数组转换为数据
func (*MirObjectNPCCodec) Decode(data interface{}, msgObj interface{}) error {
	on := msgObj.(*server.ObjectNPC)
	bytes := data.([]byte)
	reader := &BytesWrapper{Bytes: &bytes}
	on.ObjectID = reader.ReadUInt32()
	on.Name = reader.ReadString()
	on.NameColor = reader.ReadInt32()
	on.Image = reader.ReadUInt16()
	on.Color = reader.ReadInt32()
	on.Location = cm.Point{X: uint32(reader.ReadInt32()), Y: uint32(reader.ReadInt32())}
	on.Direction = cm.MirDirection(reader.ReadByte())

	nc := reader.ReadInt32()
	on.QuestIDs = make([]int32, nc)
	for i := 0; i < int(nc); i++ {
		on.QuestIDs[i] = reader.ReadInt32()
	}
	return nil
}
