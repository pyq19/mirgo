package cm

import (
	"fmt"
	"math"

	"github.com/yenkeia/mirgo/game/util"
)

type DropInfo struct {
	Low           int // 1/3 中的 3，2/10 中的 2
	High          int // 1/3 中的 3，2/10 中的 10
	Count         int
	ItemName      string // ItemInfo.Name
	QuestRequired bool
}

// Account 账号
type Account struct {
	ID       int `gorm:"primary_key"`
	Username string
	Password string
}

// AccountCharacter 账号角色关系
type AccountCharacter struct {
	ID          int `gorm:"primary_key"`
	AccountID   int
	CharacterID int
}

/*
type Basic struct {
	ID            int `gorm:"primary_key"`
	GameVersion   int
	CustomVersion int
	MapIndex      int
	MonsterIndex  int
	NPCIndex      int
	QuestIndex    int
	GameShopIndex int
	ConquestIndex int
	RespawnIndex  int
}
*/

// Basic 保存上一次关闭服务端时候的自增 ID
// TODO MonsterID ItemID 分开
type Basic struct {
	ID       int32 `gorm:"primary_key"`
	ObjectID uint32
}

// Character 角色
type Character struct {
	ID               int32 `gorm:"primary_key"`
	Name             string
	Level            uint16
	Class            MirClass
	Gender           MirGender
	Hair             uint8
	CurrentMapID     int
	CurrentLocationX int
	CurrentLocationY int
	BindMapID        int
	BindLocationX    int
	BindLocationY    int
	Direction        MirDirection
	HP               uint16
	MP               uint16
	Experience       int64
	AttackMode       AttackMode
	PetMode          PetMode
	Gold             uint64 `codec:"-"` // 编码时，忽略这个字段，只用在数据库查询
	AllowGroup       bool   `codec:"-"`
}

// CharacterUserItem 角色物品关系
type CharacterUserItem struct {
	ID          int `gorm:"primary_key"`
	CharacterID int
	UserItemID  int
	Type        int // 	类型: Inventory / Equipment / QuestInventory
	Index       int //	所在类型格子的索引，比如在 Inventory 的第几个格子
}

// CharacterUserMagic 角色魔法关系
// type CharacterUserMagic struct {
// 	ID          int `gorm:"primary_key"`
// 	CharacterID int
// 	UserMagicID int
//	Character   Character `gorm:"-"` // orm ignore
//	UserMagic   UserMagic `gorm:"-"` // orm ignore
// }

// GameShopItem 游戏内商城物品
type GameShopItem struct {
	ID          int `gorm:"primary_key"`
	ItemID      int
	GoldPrice   int
	CreditPrice int
	Count       int
	Class       string
	Category    string
	Stock       int
	IStock      int
	Deal        int
	TopItem     int
	//CreateDate
}

type ItemInfo struct {
	ID             int32 `gorm:"primary_key"`
	Name           string
	Type           ItemType
	Grade          ItemGrade
	RequiredType   RequiredType
	RequiredClass  RequiredClass
	RequiredGender RequiredGender
	ItemSet        ItemSet
	Shape          int16
	Weight         uint8
	Light          uint8
	RequiredAmount uint8
	Image          uint16
	Durability     uint16
	StackSize      uint32
	Price          uint32
	MinAC          uint8
	MaxAC          uint8
	MinMAC         uint8
	MaxMAC         uint8
	MinDC          uint8
	MaxDC          uint8
	MinMC          uint8
	MaxMC          uint8
	MinSC          uint8
	MaxSC          uint8
	HP             uint16
	MP             uint16
	Accuracy       uint8
	Agility        uint8
	Luck           int8
	AttackSpeed    int8
	StartItem      bool
	BagWeight      uint8
	HandWeight     uint8
	WearWeight     uint8
	Effect         uint8
	Strong         uint8
	MagicResist    uint8
	PoisonResist   uint8
	HealthRecovery uint8
	SpellRecovery  uint8
	PoisonRecovery uint8
	HpRate         uint8 // C# HRate
	MpRate         uint8 // C# MRate
	CriticalRate   uint8
	CriticalDamage uint8
	Bools          uint8
	MaxAcRate      uint8
	MaxMacRate     uint8
	Holy           uint8
	Freezing       uint8
	PoisonAttack   uint8
	Bind           uint16
	Reflect        uint8
	HpDrainRate    uint8
	UniqueItem     int16
	RandomStatsId  uint8
	CanFastRun     bool
	CanAwakening   bool
	IsToolTip      bool
	ToolTip        string
	ClassBased     bool `gorm:"-" codec:"-"`
	LevelBased     bool `gorm:"-" codec:"-"`
}

type MagicInfo struct {
	ID              int `gorm:"primary_key"`
	Name            string
	Spell           int
	BaseCost        int
	LevelCost       int
	Icon            int
	Level1          int
	Level2          int
	Level3          int
	Need1           int
	Need2           int
	Need3           int
	DelayBase       int
	DelayReduction  int
	PowerBase       int
	PowerBonus      int
	MPowerBase      int
	MPowerBonus     int
	MagicRange      int
	MultiplierBase  float32
	MultiplierBonus float32
}

type MapInfo struct {
	ID              int    `gorm:"primary_key"`
	Filename        string `gorm:"Column:file_name"`
	Title           string
	MiniMap         int
	BigMap          int
	Music           int
	Light           int
	MapDarkLight    int
	MineIndex       int
	NoTeleport      int
	NoReconnect     int
	NoRandom        int
	NoEscape        int
	NoRecall        int
	NoDrug          int
	NoPosition      int
	NoFight         int
	NoThrowItem     int
	NoDropPlayer    int
	NoDropMonster   int
	NoNames         int
	NoMount         int
	NeedBridle      int
	Fight           int
	Fire            int
	Lightning       int
	NoTownTeleport  int
	NoReincarnation int
	NoReconnectMap  string
	FireDamage      int
	LightningDamage int
}

//type MineZone struct {
//	MapIndex  int `gorm:"primary_key"`
//	Mine      int
//	LocationX int `gorm:"Column:location_x"`
//	LocationY int `gorm:"Column:location_y"`
//	Size      int
//}

type MonsterInfo struct {
	ID          int `gorm:"primary_key"`
	Name        string
	Image       int
	AI          int `gorm:"Column:ai"`
	Effect      int
	Level       int
	ViewRange   int
	CoolEye     int
	HP          int `gorm:"Column:hp"`
	MinAC       int
	MaxAC       int
	MinMAC      int
	MaxMAC      int
	MinDC       int
	MaxDC       int
	MinMC       int
	MaxMC       int
	MinSC       int
	MaxSC       int
	Accuracy    int
	Agility     int
	Light       int
	AttackSpeed int
	MoveSpeed   int
	Experience  int
	CanPush     int
	CanTame     int
	AutoRev     int
	Undead      int
}

type MovementInfo struct {
	ID             int `gorm:"primary_key"`
	SourceMap      int `gorm:"Column:source_map"`
	SourceX        int `gorm:"Column:source_x"`
	SourceY        int `gorm:"Column:source_y"`
	DestinationMap int `gorm:"Column:destination_map"`
	DestinationX   int `gorm:"Column:destination_x"`
	DestinationY   int `gorm:"Column:destination_y"`
	NeedHole       int
	NeedMove       int
	ConquestIndex  int
}

type NpcInfo struct {
	ID            int `gorm:"primary_key"`
	MapID         int
	Filename      string `gorm:"Column:file_name"`
	Name          string
	ChineseName   string `gorm:"Column:chinese_name"`
	LocationX     int    `gorm:"Column:location_x"`
	LocationY     int    `gorm:"Column:location_y"`
	Rate          int
	Image         int
	TimeVisible   int
	HourStart     int
	MinuteStart   int
	HourEnd       int
	MinuteEnd     int
	MinLev        int
	MaxLev        int
	DayOfWeek     string
	ClassRequired string
	FlagNeeded    int
	Conquest      int
}

type QuestInfo struct {
	ID               int `gorm:"primary_key"`
	Name             string
	QuestGroup       string
	Filename         string `gorm:"Column:file_name"`
	RequiredMinLevel int
	RequiredMaxLevel int
	RequiredQuest    int
	RequiredClass    int
	QuestType        int
	GotoMessage      string
	KillMessage      string
	ItemMessage      string
	FlagMessage      string
}

type RespawnInfo struct {
	ID              int `gorm:"primary_key"`
	MapID           int
	MonsterID       int
	LocationX       int
	LocationY       int
	Count           int
	Spread          int
	Delay           int
	RandomDelay     int
	Direction       int
	RoutePath       string
	RespawnIndex    int
	SaveRespawnTime int
	RespawnTicks    int
}

type SafeZoneInfo struct {
	ID         int `gorm:"primary_key"`
	MapID      int
	LocationX  int
	LocationY  int
	Size       int
	StartPoint int
}

type UserItem struct {
	ID             uint64 `gorm:"primary_key"` // UniqueID
	ItemID         int32
	CurrentDura    uint16
	DuraChanged    bool `gorm:"-" codec:"-"`
	MaxDura        uint16
	Count          uint32
	AC             uint8
	MAC            uint8
	DC             uint8
	MC             uint8
	SC             uint8
	Accuracy       uint8
	Agility        uint8
	HP             uint8
	MP             uint8
	AttackSpeed    int8
	Luck           int8
	SoulBoundId    uint32
	Bools          uint8
	Strong         uint8
	MagicResist    uint8
	PoisonResist   uint8
	HealthRecovery uint8
	ManaRecovery   uint8
	PoisonRecovery uint8
	CriticalRate   uint8
	CriticalDamage uint8
	Freezing       uint8
	PoisonAttack   uint8
	Info           *ItemInfo `gorm:"-" codec:"-"`
}

func (u *UserItem) Price() uint64 {
	if u.Info == nil {
		return 0
	}

	var p float64

	if u.Info.Durability > 0 {
		var r = float64(u.Info.Price) / 2.0 / float64(u.Info.Durability)

		p = float64(u.MaxDura) * r

		if u.MaxDura > 0 {
			r = float64(u.CurrentDura) / float64(u.MaxDura)
		}

		p = math.Floor(p/2.0+((p/2.0)*r)) + float64(u.Info.Price)/2.0
	}

	v := int(u.AC + u.MAC + u.DC + u.MC + u.SC + u.Accuracy + u.Agility + u.HP + u.MP)
	v += int(u.AttackSpeed + u.Luck)
	v += int(u.Strong + u.MagicResist + u.PoisonResist + u.HealthRecovery + u.ManaRecovery + u.PoisonRecovery + u.CriticalRate + u.CriticalDamage + u.Freezing + u.PoisonAttack)

	return uint64(p * (float64(v)*0.1 + 1))
}

func (u *UserItem) Clone(id uint32) *UserItem {
	return &UserItem{
		ID:             uint64(id),
		ItemID:         u.ItemID,
		CurrentDura:    u.CurrentDura,
		MaxDura:        u.MaxDura,
		Count:          u.Count,
		AC:             u.AC,
		MAC:            u.MAC,
		DC:             u.DC,
		MC:             u.MC,
		SC:             u.SC,
		Accuracy:       u.Accuracy,
		Agility:        u.Agility,
		HP:             u.HP,
		MP:             u.MP,
		AttackSpeed:    u.AttackSpeed,
		Luck:           u.Luck,
		SoulBoundId:    u.SoulBoundId,
		Bools:          u.Bools,
		Strong:         u.Strong,
		MagicResist:    u.MagicResist,
		PoisonResist:   u.PoisonResist,
		HealthRecovery: u.HealthRecovery,
		ManaRecovery:   u.ManaRecovery,
		PoisonRecovery: u.PoisonRecovery,
		CriticalRate:   u.CriticalRate,
		CriticalDamage: u.CriticalDamage,
		Freezing:       u.Freezing,
		PoisonAttack:   u.PoisonAttack,
		Info:           u.Info,
	}
}

// RepairPrice 计算修理费
func (u *UserItem) RepairPrice() uint32 {
	if u.Info == nil || u.Info.Durability == 0 {
		return 0
	}
	p := u.Info.Price
	if u.Info.Durability > 0 {
		// TODO p = (uint)Math.Floor(MaxDura * ((Info.Price / 2F) / Info.Durability) + Info.Price / 2F);
		// FIXME AttackSpeed 和 Luck 都是 int8 类型，转成 uint8 可能会有溢出问题
		p = uint32(float32(p) * (float32(u.AC+u.MAC+u.DC+u.MC+u.SC+u.Accuracy+u.Agility+u.HP+u.MP+uint8(u.AttackSpeed)+uint8(u.Luck)+u.Strong+u.MagicResist+u.PoisonResist+u.HealthRecovery+u.ManaRecovery+u.PoisonRecovery+u.CriticalRate+u.CriticalDamage+u.Freezing+u.PoisonAttack)*0.1 + 1.0))
	}
	cost := p*u.Count - uint32(u.Price())
	/* FIXME
	if u.RentalInformation == nil {
		return cost
	}
	*/
	return cost * 2
}

func (u UserItem) String() string {
	return fmt.Sprintf("UserItem ID: %d, ItemID: %d, Count: %d", u.ID, u.ItemID, u.Count)
}

type UserMagic struct {
	ID          int `gorm:"primary_key"`
	CharacterID int
	MagicID     int
	Spell       Spell
	Level       int // byte
	Key         int `gorm:"Column:magic_key"` // byte
	Experience  int // uint16
	IsTempSpell bool
	CastTime    int        // int64
	Info        *MagicInfo `gorm:"-"` // orm ignore
}

func (um *UserMagic) GetDamage(damageBase int) int {
	return damageBase + um.GetPower()
}

func (um *UserMagic) GetPower() int {
	return um.GetPower1(um.MPower())
}
func (um *UserMagic) GetPower1(power int) int {
	return int(math.Round((float64(power)/4.0)*float64(um.Level+1) + float64(um.DefPower())))
}

func (um *UserMagic) MPower() int {
	if um.Info.MPowerBonus > 0 {
		return util.RandomNext2(um.Info.MPowerBase, um.Info.MPowerBonus+um.Info.MPowerBase)
	} else {
		return um.Info.MPowerBase
	}
}
func (um *UserMagic) DefPower() int {
	if um.Info.MPowerBonus > 0 {
		return util.RandomNext2(um.Info.PowerBase, um.Info.PowerBonus+um.Info.PowerBase)
	} else {
		return um.Info.MPowerBase
	}
}

func (um *UserMagic) GetDelay() int {
	return um.Info.DelayBase - (um.Level * um.Info.DelayReduction)
}

func (um *UserMagic) GetClientMagic(info *MagicInfo) *ClientMagic {
	delay := info.DelayBase - (um.Level * info.DelayReduction)
	//castTime := (CastTime != 0) && (SMain.Envir.Time > CastTime) ? SMain.Envir.Time - CastTime : 0
	castTime := 0
	return &ClientMagic{
		Name:       info.Name,
		Spell:      um.Spell,
		BaseCost:   uint8(info.BaseCost),
		LevelCost:  uint8(info.LevelCost),
		Icon:       uint8(info.Icon),
		Level1:     uint8(info.Level1),
		Level2:     uint8(info.Level2),
		Level3:     uint8(info.Level3),
		Need1:      uint16(info.Need1),
		Need2:      uint16(info.Need2),
		Need3:      uint16(info.Need3),
		Level:      uint8(um.Level),
		Key:        uint8(um.Key),
		Experience: uint16(um.Experience),
		Delay:      int64(delay),
		Range:      uint8(info.MagicRange),
		CastTime:   int64(castTime),
	}
}

// ClientMagic 客户端显示技能
type ClientMagic struct {
	Name       string
	Spell      Spell
	BaseCost   uint8
	LevelCost  uint8
	Icon       uint8
	Level1     uint8
	Level2     uint8
	Level3     uint8
	Need1      uint16
	Need2      uint16
	Need3      uint16
	Level      uint8
	Key        uint8
	Experience uint16
	Delay      int64
	Range      uint8
	CastTime   int64
}
