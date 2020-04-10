package mircodec

import (
	"reflect"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/golog"
	"github.com/yenkeia/mirgo/game/cm"
	"github.com/yenkeia/mirgo/game/proto/server"
)

var log = golog.New("codec.mircodec")

func init() {
	codec.RegisterCodec(new(MirCodec))
}

// MirCodec 编码解码
type MirCodec struct{}

// Name 返回名字
func (m *MirCodec) Name() string {
	return "MirCodec"
}

// MimeType 我也不知道是干嘛的
func (m *MirCodec) MimeType() string {
	return "application/binary"
}

// Encode 将结构体转换为字节数组
func (*MirCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {
	switch res := msgObj.(type) {
	case *server.UserInformation:
		return encodeUserInformation(res)
	case *server.SplitItem:
		return encodeSplitItem(res)
	case *server.PlayerInspect:
		return encodePlayerInspect(res)
	case *server.ObjectPlayer:
		return encodeObjectPlayer(res)
	case *server.ObjectNPC:
		return encodeObjectNPC(res)
	case *server.NPCResponse:
		return encodeNPCResponse(res)
	case *server.TradeItem:
		return encodeTradeItem(res)
	default:
		return encode(msgObj)
	}
}

func encodeTradeItem(msgObj *server.TradeItem) (data interface{}, err error) {
	var bytes []byte
	writer := &BytesWrapper{Bytes: &bytes}
	length := len(msgObj.TradeItems)
	writer.Write(length)
	for i := 0; i < length; i++ {
		ui := msgObj.TradeItems[i]
		if ui == nil {
			writer.Write(false)
		} else {
			writer.Write(true)
			writer.Write(ui)
		}
	}
	return *writer.Bytes, nil
}

func encodeSplitItem(msgObj *server.SplitItem) (data interface{}, err error) {
	var bytes []byte
	msg := msgObj
	writer := &BytesWrapper{Bytes: &bytes}
	if msg.Item != nil {
		writer.Write(true)
		writer.Write(msg.Item)
	} else {
		writer.Write(false)
	}
	writer.Write(msg.Grid)
	return *writer.Bytes, nil
}

func encodePlayerInspect(msgObj *server.PlayerInspect) (data interface{}, err error) {
	var bytes []byte
	pi := msgObj
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

func encodeNPCResponse(msgObj *server.NPCResponse) (data interface{}, err error) {
	var bytes []byte
	res := msgObj
	writer := &BytesWrapper{Bytes: &bytes}
	count := len(res.Page)
	writer.Write(count)
	for i := 0; i < count; i++ {
		writer.Write(res.Page[i])
	}
	return *writer.Bytes, nil
}

func encodeObjectNPC(msgObj *server.ObjectNPC) (data interface{}, err error) {
	var bytes []byte
	on := msgObj
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

func encodeObjectPlayer(msgObj *server.ObjectPlayer) (data interface{}, err error) {
	var bytes []byte
	op := msgObj
	log.Debugln("发送玩家信息", op)
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

func encodeUserInformation(msgObj *server.UserInformation) (data interface{}, err error) {
	//return encode(msgObj)
	var bytes []byte
	ui := msgObj
	writer := &BytesWrapper{Bytes: &bytes}
	writer.Write(ui.ObjectID)
	writer.Write(ui.RealID)
	writer.Write(ui.Name)
	writer.Write(ui.GuildName)
	writer.Write(ui.GuildRank)
	writer.Write(ui.NameColor)
	writer.Write(ui.Class)
	writer.Write(ui.Gender)
	writer.Write(ui.Level)
	writer.Write(ui.Location.X)
	writer.Write(ui.Location.Y)
	writer.Write(ui.Direction)
	writer.Write(ui.Hair)
	writer.Write(ui.HP)
	writer.Write(ui.MP)
	writer.Write(ui.Experience)
	writer.Write(ui.MaxExperience)
	writer.Write(ui.LevelEffect)

	// Inventory
	hasInventory := true
	if ui.Inventory == nil || len(ui.Inventory) == 0 {
		hasInventory = false
	}
	writer.Write(hasInventory)
	if hasInventory {
		l := len(ui.Inventory)
		//l := 46
		writer.Write(int32(l))
		for i := 0; i < l; i++ {
			hasUserItem := !IsNull(ui.Inventory[i])
			writer.Write(hasUserItem)
			if !hasUserItem {
				continue
			}
			writer.Write(ui.Inventory[i])
		}
	}

	// Equipment
	hasEquipment := true
	if ui.Equipment == nil || len(ui.Equipment) == 0 {
		hasEquipment = false
	}
	writer.Write(hasEquipment)
	if hasEquipment {
		l := len(ui.Equipment)
		//l := 14
		writer.Write(int32(l))
		for i := 0; i < l; i++ {
			hasUserItem := !IsNull(ui.Equipment[i])
			writer.Write(hasUserItem)
			if !hasUserItem {
				continue
			}
			writer.Write(ui.Equipment[i])
		}
	}

	// QuestInventory
	hasQuestInventory := true
	if ui.QuestInventory == nil || len(ui.QuestInventory) == 0 {
		hasQuestInventory = false
	}
	writer.Write(hasQuestInventory)
	if hasQuestInventory {
		l := len(ui.QuestInventory)
		//l := 40
		writer.Write(int32(l))
		for i := 0; i < l; i++ {
			hasUserItem := !IsNull(ui.QuestInventory[i])
			writer.Write(hasUserItem)
			if !hasUserItem {
				continue
			}
			writer.Write(ui.QuestInventory[i])
		}
	}
	writer.Write(ui.Gold)
	writer.Write(ui.Credit)
	writer.Write(ui.HasExpandedStorage)
	writer.Write(ui.ExpandedStorageExpiryTime)

	count := len(ui.ClientMagics)
	writer.Write(count)
	for i := range ui.ClientMagics {
		writer.Write(ui.ClientMagics[i])
	}
	return *writer.Bytes, nil
}

func IsNull(ui *cm.UserItem) bool {
	if ui == nil || (ui.ID == 0 && ui.ItemID == 0) {
		return true
	}
	return false
}

// Decode 将字节数组转换为结构体
func (*MirCodec) Decode(data interface{}, msgObj interface{}) error {
	bytes := data.([]byte)
	switch res := msgObj.(type) {
	/*
		case *server.UserInformation:
			return decodeUserInformation(res, bytes)
		case *server.SplitItem:
			return decodeSplitItem(res, bytes)
		case *server.PlayerInspect:
			return decodePlayerInspect(res, bytes)
		case *server.ObjectPlayer:
			return decodeObjectPlayer(res, bytes)
		case *server.ObjectNPC:
			return decodeObjectNPC(res, bytes)
		case *server.NPCResponse:
			return decodeNPCResponse(res, bytes)
	*/
	default:
		return decode(res, bytes)
	}
}

// Decode 将字节数组转换为数据
func decodePlayerInspect(data interface{}, msgObj *server.PlayerInspect) error {
	pi := msgObj
	bytes := data.([]byte)
	reader := &BytesWrapper{Bytes: &bytes}
	pi.Name = reader.ReadString()
	pi.GuildName = reader.ReadString()
	pi.GuildRank = reader.ReadString()
	count := reader.ReadInt32()
	pi.Equipment = make([]*cm.UserItem, count)
	for i := 0; i < int(count); i++ {
		if reader.ReadBoolean() {
			last := reader.Last()
			item := &pi.Equipment[i]
			*reader.Bytes = decodeValue(reflect.ValueOf(item), last)
		}
	}
	pi.Class = cm.MirClass(reader.ReadByte())
	pi.Gender = cm.MirGender(reader.ReadByte())
	pi.Hair = uint8(reader.ReadByte())
	pi.Level = reader.ReadUInt16()
	pi.LoverName = reader.ReadString()
	return nil
}

func decodeUserInformation(data interface{}, msgObj *server.UserInformation) error {
	ui := msgObj
	bytes := data.([]byte)
	reader := &BytesWrapper{Bytes: &bytes}
	ui.ObjectID = reader.ReadUInt32()
	ui.RealID = reader.ReadUInt32()
	ui.Name = reader.ReadString()
	ui.GuildName = reader.ReadString()
	ui.GuildRank = reader.ReadString()
	ui.NameColor = reader.ReadInt32()
	ui.Class = cm.MirClass(reader.ReadByte())
	ui.Gender = cm.MirGender(reader.ReadByte())
	ui.Level = reader.ReadUInt16()
	x := reader.ReadInt32()
	y := reader.ReadInt32()
	ui.Location = cm.Point{X: uint32(x), Y: uint32(y)}
	ui.Direction = cm.MirDirection(reader.ReadByte())
	ui.Hair = reader.ReadUInt8()
	ui.HP = reader.ReadUInt16()
	ui.MP = reader.ReadUInt16()
	ui.Experience = reader.ReadInt64()
	ui.MaxExperience = reader.ReadInt64()
	ui.LevelEffect = cm.LevelEffects(reader.ReadUInt8())

	// Inventory
	if reader.ReadBoolean() {
		count := reader.ReadInt32()
		ui.Inventory = make([]*cm.UserItem, count)
		for i := 0; i < int(count); i++ {
			if reader.ReadBoolean() {
				last := reader.Last()
				// FIXME
				item := &ui.Inventory[i]
				*reader.Bytes = decodeValue(reflect.ValueOf(item), last)
			}
		}
	}

	// Equipment
	if reader.ReadBoolean() {
		count := reader.ReadInt32()
		ui.Equipment = make([]*cm.UserItem, count)
		for i := 0; i < int(count); i++ {
			if reader.ReadBoolean() {
				last := reader.Last()
				// FIXME
				item := &ui.Equipment[i]
				*reader.Bytes = decodeValue(reflect.ValueOf(item), last)
			}
		}
	}

	// QuestInventory
	if reader.ReadBoolean() {
		count := reader.ReadInt32()
		ui.QuestInventory = make([]*cm.UserItem, count)
		for i := 0; i < int(count); i++ {
			if reader.ReadBoolean() {
				last := reader.Last()
				// FIXME
				item := &ui.QuestInventory[i]
				*reader.Bytes = decodeValue(reflect.ValueOf(item), last)
			}
		}
	}
	ui.Gold = reader.ReadUInt32()
	ui.Credit = reader.ReadUInt32()
	ui.HasExpandedStorage = reader.ReadBoolean()
	ui.ExpandedStorageExpiryTime = reader.ReadInt64()

	count := reader.ReadInt32()
	clientMagics := make([]*cm.ClientMagic, 0)
	for i := 0; i < int(count); i++ {
		last := reader.Last()
		magic := new(cm.ClientMagic)
		*reader.Bytes = decodeValue(reflect.ValueOf(magic), last)
		clientMagics = append(clientMagics, magic)
	}
	ui.ClientMagics = clientMagics
	return nil
}

func decodeObjectPlayer(data interface{}, msgObj *server.ObjectPlayer) error {
	op := msgObj
	bytes := data.([]byte)
	reader := &BytesWrapper{Bytes: &bytes}
	op.ObjectID = reader.ReadUInt32()
	op.Name = reader.ReadString()
	op.GuildName = reader.ReadString()
	op.GuildRankName = reader.ReadString()
	op.NameColor = reader.ReadInt32()
	op.Class = cm.MirClass(reader.ReadByte())
	op.Gender = cm.MirGender(reader.ReadByte())
	op.Level = reader.ReadUInt16()
	op.Location = cm.Point{X: uint32(reader.ReadInt32()), Y: uint32(reader.ReadInt32())}
	op.Direction = cm.MirDirection(reader.ReadByte())
	op.Hair = uint8(reader.ReadByte())
	op.Light = uint8(reader.ReadByte())
	op.Weapon = reader.ReadInt16()
	op.WeaponEffect = reader.ReadInt16()
	op.Armour = reader.ReadInt16()
	op.Poison = cm.PoisonType(reader.ReadUInt16())
	op.Dead = reader.ReadBoolean()
	op.Hidden = reader.ReadBoolean()
	op.Effect = cm.SpellEffect(reader.ReadByte())
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
	op.Buffs = make([]cm.BuffType, bc)
	for i := 0; i < int(bc); i++ {
		op.Buffs[i] = cm.BuffType(reader.ReadByte())
	}

	op.LevelEffects = cm.LevelEffects(reader.ReadByte())
	return nil
}

func decodeObjectNPC(data interface{}, msgObj *server.ObjectNPC) error {
	on := msgObj
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

func decodeNPCResponse(data interface{}, msgObj *server.NPCResponse) error {
	res := msgObj
	bytes := data.([]byte)
	reader := &BytesWrapper{Bytes: &bytes}
	count := int(reader.ReadInt32())
	res.Page = make([]string, count)
	for i := 0; i < count; i++ {
		res.Page[i] = reader.ReadString()
	}
	return nil
}
