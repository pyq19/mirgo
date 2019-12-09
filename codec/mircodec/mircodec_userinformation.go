package mircodec

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/yenkeia/mirgo/proto/server"
)

// MirUserInformationCodec ...
type MirUserInformationCodec struct{}

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

// TODO Decode 将字节数组转换为数据
func (*MirUserInformationCodec) Decode(data interface{}, msgObj interface{}) error {
	ui := msgObj.(*server.UserInformation)
	bytes := data.([]byte)
	reader := &BytesWrapper{Bytes: &bytes}
	ui.ObjectID = reader.ReadUInt32()
	ui.RealId = reader.ReadUInt32()
	ui.Name = reader.ReadString()
	ui.GuildName = reader.ReadString()
	ui.GuildRank = reader.ReadString()
	//ui.NameColour = Color.FromArgb(reader.ReadInt32());
	//ui.Class = (MirClass)reader.ReadByte();
	//ui.Gender = (MirGender)reader.ReadByte();
	//ui.Level = reader.ReadUInt16();
	//ui.Location = new Point(reader.ReadInt32(), reader.ReadInt32());
	//ui.Direction = (MirDirection)reader.ReadByte();
	//ui.Hair = reader.ReadByte();
	//ui.HP = reader.ReadUInt16();
	//ui.MP = reader.ReadUInt16();
	//
	//ui.Experience = reader.ReadInt64();
	//ui.MaxExperience = reader.ReadInt64();
	//
	//ui.LevelEffect = nil
	return decode(msgObj, bytes)
}

func init() {
	codec.RegisterCodec(new(MirUserInformationCodec))
}
