package common

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

// Character 角色
type Character struct {
	ID               int32 `gorm:"primary_key"`
	Name             string
	Level            uint16
	Class            MirClass
	Gender           MirGender
	Hair             uint8
	CurrentMapID     int32
	CurrentLocationX int32
	CurrentLocationY int32
	Direction        MirDirection
	HP               uint16
	MP               uint16
	Experience       int64
	AttackMode       AttackMode
	PetMode          PetMode
	Gold             uint64 `codec:"-"` // 编码时，忽略这个字段，只用在数据库查询
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
type CharacterUserMagic struct {
	ID          int `gorm:"primary_key"`
	CharacterID int
	UserMagicID int
	//Character   Character `gorm:"-"` // orm ignore
	//UserMagic   UserMagic `gorm:"-"` // orm ignore
}

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
	Bind           int16
	Reflect        uint8
	HpDrainRate    uint8
	UniqueItem     int16
	RandomStatsId  uint8
	CanFastRun     bool
	CanAwakening   bool
	IsToolTip      bool
	ToolTip        string
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
	ID            int `gorm:"primary_key"`
	MapID         int
	SourceX       int `gorm:"Column:source_x"`
	SourceY       int `gorm:"Column:source_y"`
	DestinationX  int `gorm:"Column:destination_x"`
	DestinationY  int `gorm:"Column:destination_y"`
	NeedHole      int
	NeedMove      int
	ConquestIndex int
}

type NpcInfo struct {
	ID            int `gorm:"primary_key"`
	MapID         int
	Filename      string `gorm:"Column:file_name"`
	Name          string
	LocationX     int `gorm:"Column:location_x"`
	LocationY     int `gorm:"Column:location_y"`
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
}

type UserMagic struct {
	ID          int `gorm:"primary_key"`
	CharacterID int
	MagicID     int
	Spell       Spell
	Level       int // byte
	Key         int // byte
	Experience  int // uint16
	IsTempSpell bool
	CastTime    int // int64
	//Magic       MagicInfo `gorm:"-"` // orm ignore
}

// TODO
func (um *UserMagic) GetDamage(damageBase int) int {
	return 0
}

// ClientMagic 客户端显示技能
type ClientMagic struct {
	Spell      Spell
	BaseCost   uint8
	LevelCost  uint8
	Icon       uint8
	Level1     uint8
	Level2     uint8
	Level3     uint8
	Need1      uint8
	Need2      uint8
	Need3      uint8
	Level      uint8
	Key        uint8
	Experience uint16
	Delay      int64
	Range      uint8
	CastTime   int64
}
