package mircodec

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/yenkeia/mirgo/common"
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
	op := msgObj.(*server.ObjectPlayer)
	log.Debugln(op)
	writer := &BytesWrapper{Bytes: &bytes}
	writer.Write(op.ObjectID)
	writer.Write(op.Name)
	writer.Write(op.GuildName)
	writer.Write(op.GuildRankName)
	writer.Write(op.NameColor)
	writer.Write(op.Class)
	writer.Write(op.Gender)
	writer.Write(op.Level)
	writer.Write(op.Location.X)
	writer.Write(op.Location.Y)
	writer.Write(op.Direction)
	writer.Write(op.Hair)
	writer.Write(op.Light)
	writer.Write(op.Weapon)
	writer.Write(op.WeaponEffect)
	writer.Write(op.Armour)
	writer.Write(op.Poison)
	writer.Write(op.Dead)
	writer.Write(op.Hidden)
	writer.Write(op.Effect)
	writer.Write(op.WingEffect)
	writer.Write(op.Extra)
	writer.Write(op.MountType)
	writer.Write(op.RidingMount)
	writer.Write(op.Fishing)
	writer.Write(op.TransformType)
	writer.Write(op.ElementOrbEffect)
	writer.Write(op.ElementOrbLvl)
	writer.Write(op.ElementOrbMax)
	bc := len(op.Buffs)
	writer.Write(bc)
	for i := range op.Buffs {
		b := op.Buffs[i]
		writer.Write(b)
	}
	writer.Write(op.LevelEffects)
	b := *writer.Bytes
	_ = b
	return *writer.Bytes, nil
}

// Decode 将字节数组转换为数据
func (*MirObjectPlayerCodec) Decode(data interface{}, msgObj interface{}) error {
	op := msgObj.(*server.ObjectPlayer)
	bytes := data.([]byte)
	reader := &BytesWrapper{Bytes: &bytes}
	op.ObjectID = reader.ReadUInt32()
	op.Name = reader.ReadString()
	op.GuildName = reader.ReadString()
	op.GuildRankName = reader.ReadString()
	op.NameColor = reader.ReadInt32()
	op.Class = common.MirClass(reader.ReadByte())
	op.Gender = common.MirGender(reader.ReadByte())
	op.Level = reader.ReadUInt16()
	op.Location = common.Point{X: uint32(reader.ReadInt32()), Y: uint32(reader.ReadInt32())}
	op.Direction = common.MirDirection(reader.ReadByte())
	op.Hair = uint8(reader.ReadByte())
	op.Light = uint8(reader.ReadByte())
	op.Weapon = reader.ReadInt16()
	op.WeaponEffect = reader.ReadInt16()
	op.Armour = reader.ReadInt16()
	op.Poison = common.PoisonType(reader.ReadUInt16())
	op.Dead = reader.ReadBoolean()
	op.Hidden = reader.ReadBoolean()
	op.Effect = common.SpellEffect(reader.ReadByte())
	op.WingEffect = uint8(reader.ReadByte())
	op.Extra = reader.ReadBoolean()
	op.MountType = reader.ReadInt16()
	op.RidingMount = reader.ReadBoolean()
	op.Fishing = reader.ReadBoolean()
	op.TransformType = reader.ReadInt16()
	op.ElementOrbEffect = reader.ReadUInt32()
	op.ElementOrbLvl = reader.ReadUInt32()
	op.ElementOrbMax = reader.ReadUInt32()

	bc := reader.ReadInt32()
	op.Buffs = make([]common.BuffType, bc)
	for i := 0; i < int(bc); i++ {
		op.Buffs[i] = common.BuffType(reader.ReadByte())
	}

	op.LevelEffects = common.LevelEffects(reader.ReadByte())
	return nil
}
