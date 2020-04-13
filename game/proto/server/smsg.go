package server

import (
	"fmt"

	// 使用binary协议，因此匿名引用这个包，底层会自动注册

	_ "github.com/davyxu/cellnet/codec/binary"
	"github.com/yenkeia/mirgo/game/cm"
)

const (
	CONNECTED = 2000 + iota
	CLIENT_VERSION
	DISCONNECT
	KEEP_ALIVE
	NEW_ACCOUNT
	CHANGE_PASSWORD
	CHANGE_PASSWORD_BANNED
	LOGIN
	LOGIN_BANNED
	LOGIN_SUCCESS
	NEW_CHARACTER
	NEW_CHARACTER_SUCCESS
	DELETE_CHARACTER
	DELETE_CHARACTER_SUCCESS
	START_GAME
	START_GAME_BANNED
	START_GAME_DELAY
	MAP_INFORMATION
	USER_INFORMATION
	USER_LOCATION
	OBJECT_PLAYER
	OBJECT_REMOVE
	OBJECT_TURN
	OBJECT_WALK
	OBJECT_RUN
	CHAT
	OBJECT_CHAT
	NEW_ITEM_INFO
	MOVE_ITEM
	EQUIP_ITEM
	MERGE_ITEM
	REMOVE_ITEM
	REMOVE_SLOT_ITEM
	TAKE_BACK_ITEM
	STORE_ITEM
	SPLIT_ITEM
	SPLIT_ITEM1
	DEPOSIT_REFINE_ITEM
	RETRIEVE_REFINE_ITEM
	REFINE_CANCEL
	REFINE_ITEM
	DEPOSIT_TRADE_ITEM
	RETRIEVE_TRADE_ITEM
	USE_ITEM
	DROP_ITEM
	PLAYER_UPDATE
	PLAYER_INSPECT
	LOG_OUT_SUCCESS
	LOG_OUT_FAILED
	TIME_OF_DAY
	CHANGE_A_MODE
	CHANGE_P_MODE
	OBJECT_ITEM
	OBJECT_GOLD
	GAINED_ITEM
	GAINED_GOLD
	LOSE_GOLD
	GAINED_CREDIT
	LOSE_CREDIT
	OBJECT_MONSTER
	OBJECT_ATTACK
	STRUCK
	OBJECT_STRUCK
	DAMAGE_INDICATOR
	DURA_CHANGED
	HEALTH_CHANGED
	DELETE_ITEM
	DEATH
	OBJECT_DIED
	COLOUR_CHANGED
	OBJECT_COLOUR_CHANGED
	OBJECT_GUILD_NAME_CHANGED
	GAIN_EXPERIENCE
	LEVEL_CHANGED
	OBJECT_LEVELED
	OBJECT_HARVEST
	OBJECT_HARVESTED
	OBJECT_NPC
	NPC_RESPONSE
	OBJECT_HIDE
	OBJECT_SHOW
	POISONED
	OBJECT_POISONED
	MAP_CHANGED
	OBJECT_TELEPORT_OUT
	OBJECT_TELEPORT_IN
	TELEPORT_IN
	NPC_GOODS
	NPC_SELL
	NPC_REPAIR
	NPC_S_REPAIR
	NPC_REFINE
	NPC_CHECK_REFINE
	NPC_COLLECT_REFINE
	NPC_REPLACE_WED_RING
	NPC_STORAGE
	SELL_ITEM
	CRAFT_ITEM
	REPAIR_ITEM
	ITEM_REPAIRED
	NEW_MAGIC
	REMOVE_MAGIC
	MAGIC_LEVELED
	MAGIC
	MAGIC_DELAY
	MAGIC_CAST
	OBJECT_MAGIC
	OBJECT_EFFECT
	RANGE_ATTACK
	PUSHED
	OBJECT_PUSHED
	OBJECT_NAME
	USER_STORAGE
	SWITCH_GROUP
	DELETE_GROUP
	DELETE_MEMBER
	GROUP_INVITE
	ADD_MEMBER
	REVIVED
	OBJECT_REVIVED
	SPELL_TOGGLE
	OBJECT_HEALTH
	MAP_EFFECT
	OBJECT_RANGE_ATTACK
	ADD_BUFF
	REMOVE_BUFF
	OBJECT_HIDDEN
	REFRESH_ITEM
	OBJECT_SPELL
	USER_DASH
	OBJECT_DASH
	USER_DASH_FAIL
	OBJECT_DASH_FAIL
	NPC_CONSIGN
	NPC_MARKET
	NPC_MARKET_PAGE
	CONSIGN_ITEM
	MARKET_FAIL
	MARKET_SUCCESS
	OBJECT_SIT_DOWN
	IN_TRAP_ROCK
	BASE_STATS_INFO
	USER_NAME
	CHAT_ITEM_STATS
	GUILD_NOTICE_CHANGE
	GUILD_MEMBER_CHANGE
	GUILD_STATUS
	GUILD_INVITE
	GUILD_EXP_GAIN
	GUILD_NAME_REQUEST
	GUILD_STORAGE_GOLD_CHANGE
	GUILD_STORAGE_ITEM_CHANGE
	GUILD_STORAGE_LIST
	GUILD_REQUEST_WAR
	DEFAULT_NPC
	NPC_UPDATE
	NPC_IMAGE_UPDATE
	MARRIAGE_REQUEST
	DIVORCE_REQUEST
	MENTOR_REQUEST
	TRADE_REQUEST
	TRADE_ACCEPT
	TRADE_GOLD
	TRADE_ITEM
	TRADE_CONFIRM
	TRADE_CANCEL
	MOUNT_UPDATE
	EQUIP_SLOT_ITEM
	FISHING_UPDATE
	CHANGE_QUEST
	COMPLETE_QUEST
	SHARE_QUEST
	NEW_QUEST_INFO
	GAINED_QUEST_ITEM
	DELETE_QUEST_ITEM
	CANCEL_REINCARNATION
	REQUEST_REINCARNATION
	USER_BACK_STEP
	OBJECT_BACK_STEP
	USER_DASH_ATTACK
	OBJECT_DASH_ATTACK
	USER_ATTACK_MOVE
	COMBINE_ITEM
	ITEM_UPGRADED
	SET_CONCENTRATION
	SET_OBJECT_CONCENTRATION
	SET_ELEMENTAL
	SET_OBJECT_ELEMENTAL
	REMOVE_DELAYED_EXPLOSION
	OBJECT_DECO
	OBJECT_SNEAKING
	OBJECT_LEVEL_EFFECTS
	SET_BINDING_SHOT
	SEND_OUTPUT_MESSAGE
	NPC_AWAKENING
	NPC_DISASSEMBLE
	NPC_DOWNGRADE
	NPC_RESET
	AWAKENING_NEED_MATERIALS
	AWAKENING_LOCKED_ITEM
	AWAKENING
	RECEIVE_MAIL
	MAIL_LOCKED_ITEM
	MAIL_SEND_REQUEST
	MAIL_SENT
	PARCEL_COLLECTED
	MAIL_COST
	RESIZE_INVENTORY
	RESIZE_STORAGE
	NEW_INTELLIGENT_CREATURE
	UPDATE_INTELLIGENT_CREATURElIST
	INTELLIGENT_CREATURE_ENABLE_RENAME
	INTELLIGENT_CREATURE_PICKUP
	NPC_PEARL_GOODS
	TRANSFORM_UPDATE
	FRIEND_UPDATE
	LOVER_UPDATE
	MENTOR_UPDATE
	GUILD_BUFF_LIST
	NPC_REQUEST_INPUT
	GAME_SHOP_INFO
	GAME_SHOP_STOCK
	RANKINGS
	OPENDOOR
	GET_RENTED_ITEMS
	ITEM_RENTAL_REQUEST
	ITEM_RENTAL_FEE
	ITEM_RENTAL_PERIOD
	DEPOSIT_RENTAL_ITEM
	RETRIEVE_RENTAL_ITEM
	UPDATE_RENTAL_ITEM
	CANCEL_ITEM_RENTAL
	ITEM_RENTAL_LOCK
	ITEM_RENTAL_PARTNER_LOCK
	CAN_CONFIRM_ITEM_RENTAL
	CONFIRM_ITEM_RENTAL
	NEW_RECIPE_INFO
	OPEN_BROWSER
)

type Connected struct{}

type ClientVersion struct {
	Result uint8
}

type Disconnect struct {
	Reason uint8
	/*
	 * 0: Server Closing.
	 * 1: Another User.
	 * 2: Packet Error.
	 * 3: Server Crashed.
	 */
}

type KeepAlive struct {
	Time int64
}

type NewAccount struct {
	Result uint8
}

type ChangePassword struct {
	Result uint8
}

// TODO
type ChangePasswordBanned struct {
	//public string Reason = string.Empty;
	//public DateTime ExpiryDate;
}

type Login struct {
	Result uint8
}

// TODO
type LoginBanned struct {
	//public string Reason = string.Empty;
	//public DateTime ExpiryDate
}

type SelectInfo struct {
	Index      uint32
	Name       string
	Level      uint16
	Class      cm.MirClass
	Gender     cm.MirGender
	LastAccess int64
}

type LoginSuccess struct {
	Characters []SelectInfo
}

type NewCharacter struct {
	Result uint8
}

type NewCharacterSuccess struct {
	CharInfo SelectInfo
}

type DeleteCharacter struct {
	Result uint8
	/*
	 * 0: Disabled.
	 * 1: Character Not Found
	 * */
}

type DeleteCharacterSuccess struct {
	CharacterIndex int32
}

type StartGame struct {
	Result uint8
	/*
	 * 0: Disabled.
	 * 1: Not logged in
	 * 2: Character not found.
	 * 3: Start Game Error
	 * */

	Resolution int32
}

type StartGameBanned struct {
	Reason     string
	ExpiryDate int64 // DateTime
}

type StartGameDelay struct {
	Milliseconds int64
}

type MapInformation struct {
	FileName     string
	Title        string
	MiniMap      uint16
	BigMap       uint16
	Lights       cm.LightSetting
	Lightning    bool
	MapDarkLight uint8
	Music        uint16
}

type UserInformation struct {
	ObjectID                  uint32
	RealID                    uint32
	Name                      string
	GuildName                 string
	GuildRank                 string
	NameColor                 int32
	Class                     cm.MirClass
	Gender                    cm.MirGender
	Level                     uint16
	Location                  cm.Point
	Direction                 cm.MirDirection
	Hair                      uint8
	HP                        uint16
	MP                        uint16
	Experience                int64
	MaxExperience             int64
	LevelEffect               cm.LevelEffects
	Inventory                 []*cm.UserItem
	Equipment                 []*cm.UserItem
	QuestInventory            []*cm.UserItem
	Gold                      uint32
	Credit                    uint32
	HasExpandedStorage        bool
	ExpandedStorageExpiryTime int64
	ClientMagics              []*cm.ClientMagic
}

func (msg *UserInformation) String() string {
	return fmt.Sprintf("UserInformation: %s(%d)\n", msg.Name, msg.ObjectID)
}

type UserLocation struct {
	Location  cm.Point
	Direction cm.MirDirection
}

type ObjectPlayer struct {
	ObjectID         uint32
	Name             string
	GuildName        string
	GuildRankName    string
	NameColor        int32 // = Color.FromArgb(reader.ReadInt32());
	Class            cm.MirClass
	Gender           cm.MirGender
	Level            uint16
	Location         cm.Point
	Direction        cm.MirDirection
	Hair             uint8
	Light            uint8
	Weapon           int16
	WeaponEffect     int16
	Armour           int16
	Poison           cm.PoisonType // = (PoisonType)reader.ReadUInt16()
	Dead             bool
	Hidden           bool
	Effect           cm.SpellEffect // = (SpellEffect)reader.ReadByte()
	WingEffect       uint8
	Extra            bool
	MountType        int16
	RidingMount      bool
	Fishing          bool
	TransformType    int16
	ElementOrbEffect uint32
	ElementOrbLvl    uint32
	ElementOrbMax    uint32
	Buffs            []cm.BuffType
	LevelEffects     cm.LevelEffects
}

func (msg *ObjectPlayer) String() string {
	return fmt.Sprintf("ObjectPlayer Name: %s(%d), Location: %s, Class: %d, Gender: %d, Location: %s, Direction: %d, Hair: %d, Light: %d, Weapon: %d, WeaponEffect: %d, Armour: %d", msg.Name, msg.ObjectID, msg.Location, msg.Class, msg.Gender, msg.Location, msg.Direction, msg.Hair, msg.Light, msg.Weapon, msg.WeaponEffect, msg.Armour)
}

type ObjectRemove struct {
	ObjectID uint32
}

type ObjectTurn struct {
	ObjectID  uint32
	Location  cm.Point
	Direction cm.MirDirection
}

type ObjectWalk struct {
	ObjectID  uint32
	Location  cm.Point
	Direction cm.MirDirection
}

type ObjectRun struct {
	ObjectID  uint32
	Location  cm.Point
	Direction cm.MirDirection
}

type Chat struct {
	Message string
	Type    cm.ChatType
}

type ObjectChat struct {
	ObjectID uint32
	Text     string
	Type     cm.ChatType
}

type NewItemInfo struct {
	Info *cm.ItemInfo
}

type MoveItem struct {
	Grid    cm.MirGridType
	From    int32
	To      int32
	Success bool
}

type EquipItem struct {
	Grid     cm.MirGridType
	UniqueID uint64
	To       int32
	Success  bool
}

type MergeItem struct {
	GridFrom cm.MirGridType
	GridTo   cm.MirGridType
	IDFrom   uint64
	IDTo     uint64
	Success  bool
}

type RemoveItem struct {
	Grid     cm.MirGridType
	UniqueID uint64
	To       int32
	Success  bool
}

type RemoveSlotItem struct {
	Grid     cm.MirGridType
	GridTo   cm.MirGridType
	UniqueID uint64
	To       int32
	Success  bool
}

type TakeBackItem struct {
	From    int32
	To      int32
	Success bool
}

type StoreItem struct {
	From    int32
	To      int32
	Success bool
}

type SplitItem struct {
	Item *cm.UserItem
	Grid cm.MirGridType
}

type SplitItem1 struct {
	Grid     cm.MirGridType
	UniqueID uint64
	Count    uint32
	Success  bool
}

type DepositRefineItem struct {
	From    int32
	To      int32
	Success bool
}

type RetrieveRefineItem struct {
	From    int32
	To      int32
	Success bool
}

type RefineCancel struct {
	Unlock bool
}

type RefineItem struct {
	UniqueID uint64
}

type DepositTradeItem struct {
	From    int32
	To      int32
	Success bool
}

type RetrieveTradeItem struct {
	From    int32
	To      int32
	Success bool
}

type UseItem struct {
	UniqueID uint64
	Success  bool
}

type DropItem struct {
	UniqueID uint64
	Count    uint32
	Success  bool
}

type PlayerUpdate struct {
	ObjectID     uint32
	Light        uint8
	Weapon       int16
	WeaponEffect int16
	Armour       int16
	WingEffect   uint8
}

type PlayerInspect struct {
	Name      string
	GuildName string
	GuildRank string
	Equipment []*cm.UserItem
	Class     cm.MirClass
	Gender    cm.MirGender
	Hair      uint8
	Level     uint16
	LoverName string
}

type LogOutSuccess struct {
	Characters []SelectInfo
}

// TODO
type LogOutFailed struct {
}

type TimeOfDay struct {
	Lights cm.LightSetting
}

type ChangeAMode struct {
	Mode cm.AttackMode
}

type ChangePMode struct {
	Mode cm.PetMode
}

type ObjectItem struct {
	ObjectID  uint32
	Name      string
	NameColor int32
	LocationX int32
	LocationY int32
	Image     uint16
	Grade     cm.ItemGrade
}

type ObjectGold struct {
	ObjectID  uint32
	Gold      uint32
	LocationX int32
	LocationY int32
}

type GainedItem struct {
	Item *cm.UserItem
}

type GainedGold struct {
	Gold uint32
}

type LoseGold struct {
	Gold uint32
}

type GainedCredit struct {
	Credit uint32
}

type LoseCredit struct {
	Credit uint32
}

type ObjectMonster struct {
	ObjectID          uint32
	Name              string
	NameColor         int32
	Location          cm.Point
	Image             cm.Monster
	Direction         cm.MirDirection
	Effect            uint8
	AI                uint8
	Light             uint8
	Dead              bool
	Skeleton          bool
	Poison            cm.PoisonType
	Hidden            bool
	ShockTime         int64
	BindingShotCenter bool
	Extra             bool
	ExtraByte         uint8
}

func (msg *ObjectMonster) String() string {
	return fmt.Sprintf("ObjectMonster Name: %s, ObjectID: %d, Location: %s\n", msg.Name, msg.ObjectID, msg.Location)
}

type ObjectAttack struct {
	ObjectID  uint32
	LocationX int32
	LocationY int32
	Direction cm.MirDirection
	Spell     cm.Spell
	Level     uint8
	Type      uint8
}

type Struck struct {
	AttackerID uint32
}

type ObjectStruck struct {
	ObjectID   uint32
	AttackerID uint32
	LocationX  int32
	LocationY  int32
	Direction  cm.MirDirection
}

type DamageIndicator struct {
	Damage   int32
	Type     cm.DamageType
	ObjectID uint32
}

type DuraChanged struct {
	UniqueID    uint64
	CurrentDura uint16
}

type HealthChanged struct {
	HP uint16
	MP uint16
}

type DeleteItem struct {
	UniqueID uint64
	Count    uint32
}

type Death struct {
	LocationX int32
	LocationY int32
	Direction cm.MirDirection
}

type ObjectDied struct {
	ObjectID  uint32
	LocationX int32
	LocationY int32
	Direction cm.MirDirection
	Type      uint8
}

type ColourChanged struct {
	NameColor cm.Color
}

type ObjectColourChanged struct {
	ObjectID  uint32
	NameColor cm.Color
}

type ObjectGuildNameChanged struct {
	ObjectID  uint32
	GuildName string
}

type GainExperience struct {
	Amount uint32
}

type LevelChanged struct {
	Level         uint16
	Experience    int64
	MaxExperience int64
}

type ObjectLeveled struct {
	ObjectID uint32
}

type ObjectHarvest struct {
	ObjectID  uint32
	Location  cm.Point
	Direction cm.MirDirection
}

type ObjectHarvested struct {
	ObjectID  uint32
	Location  cm.Point
	Direction cm.MirDirection
}

type ObjectNPC struct {
	ObjectID  uint32
	Name      string
	NameColor int32
	Image     uint16
	Color     int32
	Location  cm.Point
	Direction cm.MirDirection
	QuestIDs  []int32
}

func (msg *ObjectNPC) String() string {
	return fmt.Sprintf("\nObjectNPC: ID(%d) Name(%s) NameColor(%d) Image(%d) Color(%d) Location(%s) Direction(%d) QuestIDs(%v)\n",
		msg.ObjectID, msg.Name, msg.NameColor, msg.Image, msg.Color, msg.Location, msg.Direction, msg.QuestIDs)
}

type NPCResponse struct {
	Page []string
}

type ObjectHide struct {
	ObjectID uint32
}

type ObjectShow struct {
	ObjectID uint32
}

type Poisoned struct {
	Poison cm.PoisonType
}

type ObjectPoisoned struct {
	ObjectID uint32
	Poison   cm.PoisonType
}

type MapChanged struct {
	FileName     string
	Title        string
	MiniMap      uint16
	BigMap       uint16
	Lights       cm.LightSetting
	Location     cm.Point
	Direction    cm.MirDirection
	MapDarkLight uint8
	Music        uint16
}

type ObjectTeleportOut struct {
	ObjectID uint32
	Type     uint8
}

type ObjectTeleportIn struct {
	ObjectID uint32
	Type     uint8
}

type TeleportIn struct{}

type NPCGoods struct {
	Goods []*cm.UserItem
	Rate  float32
	Type  cm.PanelType
}

func (n *NPCGoods) String() string {
	return fmt.Sprintf("Goods: %s, Rate: %f, Type: %d", n.Goods, n.Rate, n.Type)
}

type NPCSell struct{}

type NPCRepair struct {
	Rate float32
}

type NPCSRepair struct {
	Rate float32
}

type NPCRefine struct {
	Rate     float32
	Refining bool
}

type NPCCheckRefine struct{}

type NPCCollectRefine struct {
	Success bool
}

type NPCReplaceWedRing struct {
	Rate float32
}

type NPCStorage struct{}

type SellItem struct {
	UniqueID uint64
	Count    uint32
	Success  bool
}

type CraftItem struct {
	Success bool
}

type RepairItem struct {
	UniqueID uint64
}

type ItemRepaired struct {
	UniqueID    uint64
	MaxDura     uint16
	CurrentDura uint16
}

type NewMagic struct {
	Magic *cm.ClientMagic
}

type RemoveMagic struct {
	PlaceID int32
}

type MagicLeveled struct {
	Spell      cm.Spell
	Level      uint8
	Experience uint16
}

type Magic struct {
	Spell    cm.Spell
	TargetID uint32
	TargetX  int32
	TargetY  int32
	Cast     bool
	Level    uint8
}

type MagicDelay struct {
	Spell cm.Spell
	Delay int64
}

type MagicCast struct {
	Spell cm.Spell
}

type ObjectMagic struct {
	ObjectID      uint32
	LocationX     int32
	LocationY     int32
	Direction     cm.MirDirection
	Spell         cm.Spell
	TargetID      uint32
	TargetX       int32
	TargetY       int32
	Cast          bool
	Level         uint8
	SelfBroadcast bool
}

type ObjectEffect struct {
	ObjectID   uint32
	Effect     cm.SpellEffect
	EffectType uint32
	DelayTime  uint32
	Time       uint32
}

type RangeAttack struct {
	TargetID uint32
	Target   cm.Point
	Spell    cm.Spell
}

type Pushed struct {
	Location  cm.Point
	Direction cm.MirDirection
}

type ObjectPushed struct {
	ObjectID  uint32
	Location  cm.Point
	Direction cm.MirDirection
}

type ObjectName struct {
	ObjectID uint32
	Name     string
}

type UserStorage struct {
	Storage []*cm.UserItem `codec:"emptyflag"`
}

type SwitchGroup struct {
	AllowGroup bool
}

type DeleteGroup struct{}

type DeleteMember struct {
	Name string
}

type GroupInvite struct {
	Name string
}

type AddMember struct {
	Name string
}

type Revived struct{}

type ObjectRevived struct {
	ObjectID uint32
	Effect   bool
}

type SpellToggle struct {
	Spell  cm.Spell
	CanUse bool
}

type ObjectHealth struct {
	ObjectID uint32
	Percent  uint8
	Expire   uint8
}

type MapEffect struct {
	Location cm.Point
	Effect   cm.SpellEffect
	Value    uint8
}

type ObjectRangeAttack struct {
	ObjectID  uint32
	Location  cm.Point
	Direction cm.MirDirection
	TargetID  uint32
	Target    cm.Point
	Type      uint8
	Spell     cm.Spell
}

type AddBuff struct {
	Type     cm.BuffType
	Caster   string
	ObjectID uint32
	Visible  bool
	Expire   int64
	Values   []int32
	Infinite bool
}

type RemoveBuff struct {
	Type     cm.BuffType
	ObjectID uint32
}

type ObjectHidden struct {
	ObjectID uint32
	Hidden   bool
}

type RefreshItem struct {
	Item cm.UserItem
}

type ObjectSpell struct {
	ObjectID  uint32
	Location  cm.Point
	Spell     cm.Spell
	Direction cm.MirDirection
	Param     bool
}

type UserDash struct {
	Location  cm.Point
	Direction cm.MirDirection
}

type ObjectDash struct {
	ObjectID  uint32
	Location  cm.Point
	Direction cm.MirDirection
}

type UserDashFail struct {
	Location  cm.Point
	Direction cm.MirDirection
}

type ObjectDashFail struct {
	ObjectID  uint32
	Location  cm.Point
	Direction cm.MirDirection
}

// TODO
type NPCConsign struct{}

// TODO
type NPCMarket struct{}

// TODO
type NPCMarketPage struct{}

type ConsignItem struct{}

type MarketFail struct{}

type MarketSuccess struct{}

type ObjectSitDown struct{}

type InTrapRock struct{}

type BaseStatsInfo struct{}

type UserName struct {
	ID   uint32
	Name string
}

type ChatItemStats struct{}

type GuildNoticeChange struct {
	Update int32
	Notice []string
}

type GuildMemberChange struct {
	Name      string
	Status    uint8
	RankIndex uint8
	Ranks     []*cm.Rank
}

type GuildStatus struct{}

type GuildInvite struct {
	Name string
}

type GuildExpGain struct{}

type GuildNameRequest struct{}

type GuildStorageGoldChange struct{}

type GuildStorageItemChange struct{}

type GuildStorageList struct{}

type GuildRequestWar struct{}

type DefaultNPC struct{}

type NPCUpdate struct {
	NPCID uint32
}

type NPCImageUpdate struct{}

type MarriageRequest struct{}

type DivorceRequest struct{}

type MentorRequest struct{}

type TradeRequest struct {
	Name string
}

type TradeAccept struct {
	Name string
}

type TradeGold struct {
	Amount uint32
}

type TradeItem struct {
	TradeItems []*cm.UserItem
}

type TradeConfirm struct{}

type TradeCancel struct {
	Unlock bool
}

type MountUpdate struct{}

type EquipSlotItem struct{}

type FishingUpdate struct{}

type ChangeQuest struct{}

type CompleteQuest struct{}

type ShareQuest struct{}

type NewQuestInfo struct{}

type GainedQuestItem struct{}

type DeleteQuestItem struct{}

type CancelReincarnation struct{}

type RequestReincarnation struct{}

type UserBackStep struct{}

type ObjectBackStep struct{}

type UserDashAttack struct{}

type ObjectDashAttack struct{}

type UserAttackMove struct{}

type CombineItem struct{}

type ItemUpgraded struct{}

type SetConcentration struct {
	ObjectID    uint32
	Enabled     bool
	Interrupted bool
}

type SetObjectConcentration struct {
	ObjectID    uint32
	Enabled     bool
	Interrupted bool
}

type SetElemental struct{}
type SetObjectElemental struct{}
type RemoveDelayedExplosion struct{}
type ObjectDeco struct{}
type ObjectSneaking struct{}
type ObjectLevelEffects struct{}
type SetBindingShot struct{}
type SendOutputMessage struct{}
type NPCAwakening struct{}
type NPCDisassemble struct{}
type NPCDowngrade struct{}
type NPCReset struct{}
type AwakeningNeedMaterials struct{}
type AwakeningLockedItem struct{}
type Awakening struct{}
type ReceiveMail struct{}
type MailLockedItem struct{}
type MailSendRequest struct{}
type MailSent struct{}
type ParcelCollected struct{}
type MailCost struct{}
type ResizeInventory struct{}
type ResizeStorage struct{}
type NewIntelligentCreature struct{}
type UpdateIntelligentCreatureList struct{}
type IntelligentCreatureEnableRename struct{}
type IntelligentCreaturePickup struct{}
type NPCPearlGoods struct{}
type TransformUpdate struct{}
type FriendUpdate struct{}
type LoverUpdate struct{}
type MentorUpdate struct{}
type GuildBuffList struct{}
type NPCRequestInput struct{}
type GameShopInfo struct{}
type GameShopStock struct{}
type Rankings struct{}
type Opendoor struct {
	DoorIndex byte
	Close     bool
}
type GetRentedItems struct{}
type ItemRentalRequest struct{}
type ItemRentalFee struct{}
type ItemRentalPeriod struct{}
type DepositRentalItem struct{}
type RetrieveRentalItem struct{}
type UpdateRentalItem struct{}
type CancelItemRental struct{}
type ItemRentalLock struct{}
type ItemRentalPartnerLock struct{}
type CanConfirmItemRental struct{}
type ConfirmItemRental struct{}
type NewRecipeInfo struct{}
type OpenBrowser struct{}
