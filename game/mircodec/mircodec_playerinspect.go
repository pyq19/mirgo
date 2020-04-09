package mircodec

import (
	"reflect"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/yenkeia/mirgo/game/proto/server"
	"github.com/yenkeia/mirgo/util"
)

func init() {
	codec.RegisterCodec(new(MirPlayerInspectCodec))
}

/*
MirPlayerInspectCodec
*/
type MirPlayerInspectCodec struct{}

// Name 编码器的名字
func (*MirPlayerInspectCodec) Name() string {
	return "MirPlayerInspectCodec"
}

// MimeType 兼容http类型
func (*MirPlayerInspectCodec) MimeType() string {
	return "application/binary"
}

// Encode 将数据转换为字节数组
func (*MirPlayerInspectCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {
	var bytes []byte
	pi := msgObj.(*server.PlayerInspect)
	writer := &BytesWrapper{Bytes: &bytes}
	writer.Write(pi.Name)
	writer.Write(pi.GuildName)
	writer.Write(pi.GuildRank)
	// Equipment
	l := len(pi.Equipment)
	if l != 14 {
		panic("equipment != 14")
	}
	//l := 14
	writer.Write(int32(l))
	for i := 0; i < l; i++ {
		hasUserItem := !IsNull(pi.Equipment[i])
		writer.Write(hasUserItem)
		if !hasUserItem {
			continue
		}
		writer.Write(pi.Equipment[i])
	}

	writer.Write(pi.Class)
	writer.Write(pi.Gender)
	writer.Write(pi.Hair)
	writer.Write(pi.Level)
	writer.Write(pi.LoverName)
	return *writer.Bytes, nil
}

// Decode 将字节数组转换为数据
func (*MirPlayerInspectCodec) Decode(data interface{}, msgObj interface{}) error {
	pi := msgObj.(*server.PlayerInspect)
	bytes := data.([]byte)
	reader := &BytesWrapper{Bytes: &bytes}
	pi.Name = reader.ReadString()
	pi.GuildName = reader.ReadString()
	pi.GuildRank = reader.ReadString()
	count := reader.ReadInt32()
	pi.Equipment = make([]*util.UserItem, count)
	for i := 0; i < int(count); i++ {
		if reader.ReadBoolean() {
			last := reader.Last()
			item := &pi.Equipment[i]
			*reader.Bytes = decodeValue(reflect.ValueOf(item), last)
		}
	}
	pi.Class = util.MirClass(reader.ReadByte())
	pi.Gender = util.MirGender(reader.ReadByte())
	pi.Hair = uint8(reader.ReadByte())
	pi.Level = reader.ReadUInt16()
	pi.LoverName = reader.ReadString()
	return nil
}
