package game

import (
	"fmt"

	// 使用binary协议，因此匿名引用这个包，底层会自动注册

	_ "github.com/davyxu/cellnet/codec/binary"
)

const (
	SM_CONNECTED = 2000 + iota
	SM_CLIENT_VERSION
	SM_DISCONNECT
	SM_KEEP_ALIVE
	SM_NEW_ACCOUNT
	SM_CHANGE_PASSWORD
	SM_CHANGE_PASSWORD_BANNED
	SM_LOGIN
	SM_LOGIN_BANNED
	SM_LOGIN_SUCCESS
	SM_NEW_CHARACTER
	SM_NEW_CHARACTER_SUCCESS
	SM_DELETE_CHARACTER
	SM_DELETE_CHARACTER_SUCCESS
	SM_START_GAME
	SM_START_GAME_BANNED
	SM_START_GAME_DELAY
	SM_MAP_INFORMATION
	SM_USER_INFORMATION
	SM_USER_LOCATION
	SM_OBJECT_PLAYER
	SM_OBJECT_REMOVE
	SM_OBJECT_TURN
	SM_OBJECT_WALK
	SM_OBJECT_RUN
	SM_CHAT
	SM_OBJECT_CHAT
	SM_NEW_ITEM_INFO
	SM_MOVE_ITEM
	SM_EQUIP_ITEM
	SM_MERGE_ITEM
	SM_REMOVE_ITEM
	SM_REMOVE_SLOT_ITEM
	SM_TAKE_BACK_ITEM
	SM_STORE_ITEM
	SM_SPLIT_ITEM
	SM_SPLIT_ITEM1
	SM_DEPOSIT_REFINE_ITEM
	SM_RETRIEVE_REFINE_ITEM
	SM_REFINE_CANCEL
	SM_REFINE_ITEM
	SM_DEPOSIT_TRADE_ITEM
	SM_RETRIEVE_TRADE_ITEM
	SM_USE_ITEM
	SM_DROP_ITEM
	SM_PLAYER_UPDATE
	SM_PLAYER_INSPECT
	SM_LOG_OUT_SUCCESS
	SM_LOG_OUT_FAILED
	SM_TIME_OF_DAY
	SM_CHANGE_A_MODE
	SM_CHANGE_P_MODE
	SM_OBJECT_ITEM
	SM_OBJECT_GOLD
	SM_GAINED_ITEM
	SM_GAINED_GOLD
	SM_LOSE_GOLD
	SM_GAINED_CREDIT
	SM_LOSE_CREDIT
	SM_OBJECT_MONSTER
	SM_OBJECT_ATTACK
	SM_STRUCK
	SM_OBJECT_STRUCK
	SM_DAMAGE_INDICATOR
	SM_DURA_CHANGED
	SM_HEALTH_CHANGED
	SM_DELETE_ITEM
	SM_DEATH
	SM_OBJECT_DIED
	SM_COLOUR_CHANGED
	SM_OBJECT_COLOUR_CHANGED
	SM_OBJECT_GUILD_NAME_CHANGED
	SM_GAIN_EXPERIENCE
	SM_LEVEL_CHANGED
	SM_OBJECT_LEVELED
	SM_OBJECT_HARVEST
	SM_OBJECT_HARVESTED
	SM_OBJECT_NPC
	SM_NPC_RESPONSE
	SM_OBJECT_HIDE
	SM_OBJECT_SHOW
	SM_POISONED
	SM_OBJECT_POISONED
	SM_MAP_CHANGED
	SM_OBJECT_TELEPORT_OUT
	SM_OBJECT_TELEPORT_IN
	SM_TELEPORT_IN
	SM_NPC_GOODS
	SM_NPC_SELL
	SM_NPC_REPAIR
	SM_NPC_S_REPAIR
	SM_NPC_REFINE
	SM_NPC_CHECK_REFINE
	SM_NPC_COLLECT_REFINE
	SM_NPC_REPLACE_WED_RING
	SM_NPC_STORAGE
	SM_SELL_ITEM
	SM_CRAFT_ITEM
	SM_REPAIR_ITEM
	SM_ITEM_REPAIRED
	SM_NEW_MAGIC
	SM_REMOVE_MAGIC
	SM_MAGIC_LEVELED
	SM_MAGIC
	SM_MAGIC_DELAY
	SM_MAGIC_CAST
	SM_OBJECT_MAGIC
	SM_OBJECT_EFFECT
	SM_RANGE_ATTACK
	SM_PUSHED
	SM_OBJECT_PUSHED
	SM_OBJECT_NAME
	SM_USER_STORAGE
	SM_SWITCH_GROUP
	SM_DELETE_GROUP
	SM_DELETE_MEMBER
	SM_GROUP_INVITE
	SM_ADD_MEMBER
	SM_REVIVED
	SM_OBJECT_REVIVED
	SM_SPELL_TOGGLE
	SM_OBJECT_HEALTH
	SM_MAP_EFFECT
	SM_OBJECT_RANGE_ATTACK
	SM_ADD_BUFF
	SM_REMOVE_BUFF
	SM_OBJECT_HIDDEN
	SM_REFRESH_ITEM
	SM_OBJECT_SPELL
	SM_USER_DASH
	SM_OBJECT_DASH
	SM_USER_DASH_FAIL
	SM_OBJECT_DASH_FAIL
	SM_NPC_CONSIGN
	SM_NPC_MARKET
	SM_NPC_MARKET_PAGE
	SM_CONSIGN_ITEM
	SM_MARKET_FAIL
	SM_MARKET_SUCCESS
	SM_OBJECT_SIT_DOWN
	SM_IN_TRAP_ROCK
	SM_BASE_STATS_INFO
	SM_USER_NAME
	SM_CHAT_ITEM_STATS
	SM_GUILD_NOTICE_CHANGE
	SM_GUILD_MEMBER_CHANGE
	SM_GUILD_STATUS
	SM_GUILD_INVITE
	SM_GUILD_EXP_GAIN
	SM_GUILD_NAME_REQUEST
	SM_GUILD_STORAGE_GOLD_CHANGE
	SM_GUILD_STORAGE_ITEM_CHANGE
	SM_GUILD_STORAGE_LIST
	SM_GUILD_REQUEST_WAR
	SM_DEFAULT_NPC
	SM_NPC_UPDATE
	SM_NPC_IMAGE_UPDATE
	SM_MARRIAGE_REQUEST
	SM_DIVORCE_REQUEST
	SM_MENTOR_REQUEST
	SM_TRADE_REQUEST
	SM_TRADE_ACCEPT
	SM_TRADE_GOLD
	SM_TRADE_ITEM
	SM_TRADE_CONFIRM
	SM_TRADE_CANCEL
	SM_MOUNT_UPDATE
	SM_EQUIP_SLOT_ITEM
	SM_FISHING_UPDATE
	SM_CHANGE_QUEST
	SM_COMPLETE_QUEST
	SM_SHARE_QUEST
	SM_NEW_QUEST_INFO
	SM_GAINED_QUEST_ITEM
	SM_DELETE_QUEST_ITEM
	SM_CANCEL_REINCARNATION
	SM_REQUEST_REINCARNATION
	SM_USER_BACK_STEP
	SM_OBJECT_BACK_STEP
	SM_USER_DASH_ATTACK
	SM_OBJECT_DASH_ATTACK
	SM_USER_ATTACK_MOVE
	SM_COMBINE_ITEM
	SM_ITEM_UPGRADED
	SM_SET_CONCENTRATION
	SM_SET_OBJECT_CONCENTRATION
	SM_SET_ELEMENTAL
	SM_SET_OBJECT_ELEMENTAL
	SM_REMOVE_DELAYED_EXPLOSION
	SM_OBJECT_DECO
	SM_OBJECT_SNEAKING
	SM_OBJECT_LEVEL_EFFECTS
	SM_SET_BINDING_SHOT
	SM_SEND_OUTPUT_MESSAGE
	SM_NPC_AWAKENING
	SM_NPC_DISASSEMBLE
	SM_NPC_DOWNGRADE
	SM_NPC_RESET
	SM_AWAKENING_NEED_MATERIALS
	SM_AWAKENING_LOCKED_ITEM
	SM_AWAKENING
	SM_RECEIVE_MAIL
	SM_MAIL_LOCKED_ITEM
	SM_MAIL_SEND_REQUEST
	SM_MAIL_SENT
	SM_PARCEL_COLLECTED
	SM_MAIL_COST
	SM_RESIZE_INVENTORY
	SM_RESIZE_STORAGE
	SM_NEW_INTELLIGENT_CREATURE
	SM_UPDATE_INTELLIGENT_CREATURElIST
	SM_INTELLIGENT_CREATURE_ENABLE_RENAME
	SM_INTELLIGENT_CREATURE_PICKUP
	SM_NPC_PEARL_GOODS
	SM_TRANSFORM_UPDATE
	SM_FRIEND_UPDATE
	SM_LOVER_UPDATE
	SM_MENTOR_UPDATE
	SM_GUILD_BUFF_LIST
	SM_NPC_REQUEST_INPUT
	SM_GAME_SHOP_INFO
	SM_GAME_SHOP_STOCK
	SM_RANKINGS
	SM_OPENDOOR
	SM_GET_RENTED_ITEMS
	SM_ITEM_RENTAL_REQUEST
	SM_ITEM_RENTAL_FEE
	SM_ITEM_RENTAL_PERIOD
	SM_DEPOSIT_RENTAL_ITEM
	SM_RETRIEVE_RENTAL_ITEM
	SM_UPDATE_RENTAL_ITEM
	SM_CANCEL_ITEM_RENTAL
	SM_ITEM_RENTAL_LOCK
	SM_ITEM_RENTAL_PARTNER_LOCK
	SM_CAN_CONFIRM_ITEM_RENTAL
	SM_CONFIRM_ITEM_RENTAL
	SM_NEW_RECIPE_INFO
	SM_OPEN_BROWSER
)

type SM_Connected struct{}

type SM_ClientVersion struct {
	Result uint8
}

type SM_Disconnect struct {
	Reason uint8
	/*
	 * 0: Server Closing.
	 * 1: Another User.
	 * 2: Packet Error.
	 * 3: Server Crashed.
	 */
}

type SM_KeepAlive struct {
	Time int64
}

type SM_NewAccount struct {
	Result uint8
}

type SM_ChangePassword struct {
	Result uint8
}

// TODO
type SM_ChangePasswordBanned struct {
	//public string Reason = string.Empty;
	//public DateTime ExpiryDate;
}

type SM_Login struct {
	Result uint8
}

// TODO
type SM_LoginBanned struct {
	//public string Reason = string.Empty;
	//public DateTime ExpiryDate
}

type SM_SelectInfo struct {
	Index      uint32
	Name       string
	Level      uint16
	Class      MirClass
	Gender     MirGender
	LastAccess int64
}

type SM_LoginSuccess struct {
	Characters []SM_SelectInfo
}

type SM_NewCharacter struct {
	Result uint8
}

type SM_NewCharacterSuccess struct {
	CharInfo SM_SelectInfo
}

type SM_DeleteCharacter struct {
	Result uint8
	/*
	 * 0: Disabled.
	 * 1: Character Not Found
	 * */
}

type SM_DeleteCharacterSuccess struct {
	CharacterIndex int32
}

type SM_StartGame struct {
	Result uint8
	/*
	 * 0: Disabled.
	 * 1: Not logged in
	 * 2: Character not found.
	 * 3: Start Game Error
	 * */

	Resolution int32
}

type SM_StartGameBanned struct {
	Reason     string
	ExpiryDate int64 // DateTime
}

type SM_StartGameDelay struct {
	Milliseconds int64
}

type SM_MapInformation struct {
	FileName     string
	Title        string
	MiniMap      uint16
	BigMap       uint16
	Lights       LightSetting
	Lightning    bool
	MapDarkLight uint8
	Music        uint16
}

type SM_UserInformation struct {
	ObjectID                  uint32
	RealID                    uint32
	Name                      string
	GuildName                 string
	GuildRank                 string
	NameColor                 int32
	Class                     MirClass
	Gender                    MirGender
	Level                     uint16
	Location                  Point
	Direction                 MirDirection
	Hair                      uint8
	HP                        uint16
	MP                        uint16
	Experience                int64
	MaxExperience             int64
	LevelEffect               LevelEffects
	Inventory                 []*UserItem
	Equipment                 []*UserItem
	QuestInventory            []*UserItem
	Gold                      uint32
	Credit                    uint32
	HasExpandedStorage        bool
	ExpandedStorageExpiryTime int64
	ClientMagics              []*ClientMagic
}

func (msg *SM_UserInformation) String() string {
	return fmt.Sprintf("UserInformation: %s(%d)\n", msg.Name, msg.ObjectID)
}

type SM_UserLocation struct {
	Location  Point
	Direction MirDirection
}

type SM_ObjectPlayer struct {
	ObjectID         uint32
	Name             string
	GuildName        string
	GuildRankName    string
	NameColor        int32 // = Color.FromArgb(reader.ReadInt32());
	Class            MirClass
	Gender           MirGender
	Level            uint16
	Location         Point
	Direction        MirDirection
	Hair             uint8
	Light            uint8
	Weapon           int16
	WeaponEffect     int16
	Armour           int16
	Poison           PoisonType // = (PoisonType)reader.ReadUInt16()
	Dead             bool
	Hidden           bool
	Effect           SpellEffect // = (SpellEffect)reader.ReadByte()
	WingEffect       uint8
	Extra            bool
	MountType        int16
	RidingMount      bool
	Fishing          bool
	TransformType    int16
	ElementOrbEffect uint32
	ElementOrbLvl    uint32
	ElementOrbMax    uint32
	Buffs            []BuffType
	LevelEffects     LevelEffects
}

func (msg *SM_ObjectPlayer) String() string {
	return fmt.Sprintf("ObjectPlayer Name: %s(%d), Location: %s, Class: %d, Gender: %d, Location: %s, Direction: %d, Hair: %d, Light: %d, Weapon: %d, WeaponEffect: %d, Armour: %d", msg.Name, msg.ObjectID, msg.Location, msg.Class, msg.Gender, msg.Location, msg.Direction, msg.Hair, msg.Light, msg.Weapon, msg.WeaponEffect, msg.Armour)
}

type SM_ObjectRemove struct {
	ObjectID uint32
}

type SM_ObjectTurn struct {
	ObjectID  uint32
	Location  Point
	Direction MirDirection
}

type SM_ObjectWalk struct {
	ObjectID  uint32
	Location  Point
	Direction MirDirection
}

type SM_ObjectRun struct {
	ObjectID  uint32
	Location  Point
	Direction MirDirection
}

type SM_Chat struct {
	Message string
	Type    ChatType
}

type SM_ObjectChat struct {
	ObjectID uint32
	Text     string
	Type     ChatType
}

type SM_NewItemInfo struct {
	Info *ItemInfo
}

type SM_MoveItem struct {
	Grid    MirGridType
	From    int32
	To      int32
	Success bool
}

type SM_EquipItem struct {
	Grid     MirGridType
	UniqueID uint64
	To       int32
	Success  bool
}

type SM_MergeItem struct {
	GridFrom MirGridType
	GridTo   MirGridType
	IDFrom   uint64
	IDTo     uint64
	Success  bool
}

type SM_RemoveItem struct {
	Grid     MirGridType
	UniqueID uint64
	To       int32
	Success  bool
}

type SM_RemoveSlotItem struct {
	Grid     MirGridType
	GridTo   MirGridType
	UniqueID uint64
	To       int32
	Success  bool
}

type SM_TakeBackItem struct {
	From    int32
	To      int32
	Success bool
}

type SM_StoreItem struct {
	From    int32
	To      int32
	Success bool
}

type SM_SplitItem struct {
	Item *UserItem
	Grid MirGridType
}

type SM_SplitItem1 struct {
	Grid     MirGridType
	UniqueID uint64
	Count    uint32
	Success  bool
}

type SM_DepositRefineItem struct {
	From    int32
	To      int32
	Success bool
}

type SM_RetrieveRefineItem struct {
	From    int32
	To      int32
	Success bool
}

type SM_RefineCancel struct {
	Unlock bool
}

type SM_RefineItem struct {
	UniqueID uint64
}

type SM_DepositTradeItem struct {
	From    int32
	To      int32
	Success bool
}

type SM_RetrieveTradeItem struct {
	From    int32
	To      int32
	Success bool
}

type SM_UseItem struct {
	UniqueID uint64
	Success  bool
}

type SM_DropItem struct {
	UniqueID uint64
	Count    uint32
	Success  bool
}

type SM_PlayerUpdate struct {
	ObjectID     uint32
	Light        uint8
	Weapon       int16
	WeaponEffect int16
	Armour       int16
	WingEffect   uint8
}

type SM_PlayerInspect struct {
	Name      string
	GuildName string
	GuildRank string
	Equipment []*UserItem
	Class     MirClass
	Gender    MirGender
	Hair      uint8
	Level     uint16
	LoverName string
}

type SM_LogOutSuccess struct {
	Characters []SM_SelectInfo
}

// TODO
type SM_LogOutFailed struct {
}

type SM_TimeOfDay struct {
	Lights LightSetting
}

type SM_ChangeAMode struct {
	Mode AttackMode
}

type SM_ChangePMode struct {
	Mode PetMode
}

type SM_ObjectItem struct {
	ObjectID  uint32
	Name      string
	NameColor int32
	LocationX int32
	LocationY int32
	Image     uint16
	Grade     ItemGrade
}

type SM_ObjectGold struct {
	ObjectID  uint32
	Gold      uint32
	LocationX int32
	LocationY int32
}

type SM_GainedItem struct {
	Item *UserItem
}

type SM_GainedGold struct {
	Gold uint32
}

type SM_LoseGold struct {
	Gold uint32
}

type SM_GainedCredit struct {
	Credit uint32
}

type SM_LoseCredit struct {
	Credit uint32
}

type SM_ObjectMonster struct {
	ObjectID          uint32
	Name              string
	NameColor         int32
	Location          Point
	Image             MonsterType
	Direction         MirDirection
	Effect            uint8
	AI                uint8
	Light             uint8
	Dead              bool
	Skeleton          bool
	Poison            PoisonType
	Hidden            bool
	ShockTime         int64
	BindingShotCenter bool
	Extra             bool
	ExtraByte         uint8
}

func (msg *SM_ObjectMonster) String() string {
	return fmt.Sprintf("ObjectMonster Name: %s, ObjectID: %d, Location: %s\n", msg.Name, msg.ObjectID, msg.Location)
}

type SM_ObjectAttack struct {
	ObjectID  uint32
	LocationX int32
	LocationY int32
	Direction MirDirection
	Spell     Spell
	Level     uint8
	Type      uint8
}

type SM_Struck struct {
	AttackerID uint32
}

type SM_ObjectStruck struct {
	ObjectID   uint32
	AttackerID uint32
	LocationX  int32
	LocationY  int32
	Direction  MirDirection
}

type SM_DamageIndicator struct {
	Damage   int32
	Type     DamageType
	ObjectID uint32
}

type SM_DuraChanged struct {
	UniqueID    uint64
	CurrentDura uint16
}

type SM_HealthChanged struct {
	HP uint16
	MP uint16
}

type SM_DeleteItem struct {
	UniqueID uint64
	Count    uint32
}

type SM_Death struct {
	LocationX int32
	LocationY int32
	Direction MirDirection
}

type SM_ObjectDied struct {
	ObjectID  uint32
	LocationX int32
	LocationY int32
	Direction MirDirection
	Type      uint8
}

type SM_ColourChanged struct {
	NameColor Color
}

type SM_ObjectColourChanged struct {
	ObjectID  uint32
	NameColor Color
}

type SM_ObjectGuildNameChanged struct {
	ObjectID  uint32
	GuildName string
}

type SM_GainExperience struct {
	Amount uint32
}

type SM_LevelChanged struct {
	Level         uint16
	Experience    int64
	MaxExperience int64
}

type SM_ObjectLeveled struct {
	ObjectID uint32
}

type SM_ObjectHarvest struct {
	ObjectID  uint32
	Location  Point
	Direction MirDirection
}

type SM_ObjectHarvested struct {
	ObjectID  uint32
	Location  Point
	Direction MirDirection
}

type SM_ObjectNPC struct {
	ObjectID  uint32
	Name      string
	NameColor int32
	Image     uint16
	Color     int32
	Location  Point
	Direction MirDirection
	QuestIDs  []int32
}

func (msg *SM_ObjectNPC) String() string {
	return fmt.Sprintf("\nObjectNPC: ID(%d) Name(%s) NameColor(%d) Image(%d) Color(%d) Location(%s) Direction(%d) QuestIDs(%v)\n",
		msg.ObjectID, msg.Name, msg.NameColor, msg.Image, msg.Color, msg.Location, msg.Direction, msg.QuestIDs)
}

type SM_NPCResponse struct {
	Page []string
}

type SM_ObjectHide struct {
	ObjectID uint32
}

type SM_ObjectShow struct {
	ObjectID uint32
}

type SM_Poisoned struct {
	Poison PoisonType
}

type SM_ObjectPoisoned struct {
	ObjectID uint32
	Poison   PoisonType
}

type SM_MapChanged struct {
	FileName     string
	Title        string
	MiniMap      uint16
	BigMap       uint16
	Lights       LightSetting
	Location     Point
	Direction    MirDirection
	MapDarkLight uint8
	Music        uint16
}

type SM_ObjectTeleportOut struct {
	ObjectID uint32
	Type     uint8
}

type SM_ObjectTeleportIn struct {
	ObjectID uint32
	Type     uint8
}

type SM_TeleportIn struct{}

type SM_NPCGoods struct {
	Goods []*UserItem
	Rate  float32
	Type  PanelType
}

func (n *SM_NPCGoods) String() string {
	return fmt.Sprintf("Goods: %s, Rate: %f, Type: %d", n.Goods, n.Rate, n.Type)
}

type SM_NPCSell struct{}

type SM_NPCRepair struct {
	Rate float32
}

type SM_NPCSRepair struct {
	Rate float32
}

type SM_NPCRefine struct {
	Rate     float32
	Refining bool
}

type SM_NPCCheckRefine struct{}

type SM_NPCCollectRefine struct {
	Success bool
}

type SM_NPCReplaceWedRing struct {
	Rate float32
}

type SM_NPCStorage struct{}

type SM_SellItem struct {
	UniqueID uint64
	Count    uint32
	Success  bool
}

type SM_CraftItem struct {
	Success bool
}

type SM_RepairItem struct {
	UniqueID uint64
}

type SM_ItemRepaired struct {
	UniqueID    uint64
	MaxDura     uint16
	CurrentDura uint16
}

type SM_NewMagic struct {
	Magic *ClientMagic
}

type SM_RemoveMagic struct {
	PlaceID int32
}

type SM_MagicLeveled struct {
	Spell      Spell
	Level      uint8
	Experience uint16
}

type SM_Magic struct {
	Spell    Spell
	TargetID uint32
	TargetX  int32
	TargetY  int32
	Cast     bool
	Level    uint8
}

type SM_MagicDelay struct {
	Spell Spell
	Delay int64
}

type SM_MagicCast struct {
	Spell Spell
}

type SM_ObjectMagic struct {
	ObjectID      uint32
	LocationX     int32
	LocationY     int32
	Direction     MirDirection
	Spell         Spell
	TargetID      uint32
	TargetX       int32
	TargetY       int32
	Cast          bool
	Level         uint8
	SelfBroadcast bool
}

type SM_ObjectEffect struct {
	ObjectID   uint32
	Effect     SpellEffect
	EffectType uint32
	DelayTime  uint32
	Time       uint32
}

type SM_RangeAttack struct {
	TargetID uint32
	Target   Point
	Spell    Spell
}

type SM_Pushed struct {
	Location  Point
	Direction MirDirection
}

type SM_ObjectPushed struct {
	ObjectID  uint32
	Location  Point
	Direction MirDirection
}

type SM_ObjectName struct {
	ObjectID uint32
	Name     string
}

type SM_UserStorage struct {
	Storage []*UserItem `codec:"emptyflag"`
}

type SM_SwitchGroup struct {
	AllowGroup bool
}

type SM_DeleteGroup struct{}

type SM_DeleteMember struct {
	Name string
}

type SM_GroupInvite struct {
	Name string
}

type SM_AddMember struct {
	Name string
}

type SM_Revived struct{}

type SM_ObjectRevived struct {
	ObjectID uint32
	Effect   bool
}

type SM_SpellToggle struct {
	Spell  Spell
	CanUse bool
}

type SM_ObjectHealth struct {
	ObjectID uint32
	Percent  uint8
	Expire   uint8
}

type SM_MapEffect struct {
	Location Point
	Effect   SpellEffect
	Value    uint8
}

type SM_ObjectRangeAttack struct {
	ObjectID  uint32
	Location  Point
	Direction MirDirection
	TargetID  uint32
	Target    Point
	Type      uint8
	Spell     Spell
}

type SM_AddBuff struct {
	Type     BuffType
	Caster   string
	ObjectID uint32
	Visible  bool
	Expire   int64
	Values   []int32
	Infinite bool
}

type SM_RemoveBuff struct {
	Type     BuffType
	ObjectID uint32
}

type SM_ObjectHidden struct {
	ObjectID uint32
	Hidden   bool
}

type SM_RefreshItem struct {
	Item UserItem
}

type SM_ObjectSpell struct {
	ObjectID  uint32
	Location  Point
	Spell     Spell
	Direction MirDirection
	Param     bool
}

type SM_UserDash struct {
	Location  Point
	Direction MirDirection
}

type SM_ObjectDash struct {
	ObjectID  uint32
	Location  Point
	Direction MirDirection
}

type SM_UserDashFail struct {
	Location  Point
	Direction MirDirection
}

type SM_ObjectDashFail struct {
	ObjectID  uint32
	Location  Point
	Direction MirDirection
}

// TODO
type SM_NPCConsign struct{}

// TODO
type SM_NPCMarket struct{}

// TODO
type SM_NPCMarketPage struct{}

type SM_ConsignItem struct{}

type SM_MarketFail struct{}

type SM_MarketSuccess struct{}

type SM_ObjectSitDown struct{}

type SM_InTrapRock struct{}

type SM_BaseStatsInfo struct{}

type SM_UserName struct {
	ID   uint32
	Name string
}

type SM_ChatItemStats struct{}

type SM_GuildNoticeChange struct {
	Update int32
	Notice []string
}

type SM_GuildMemberChange struct {
	Name      string
	Status    uint8
	RankIndex uint8
	Ranks     []*Rank
}

type SM_GuildStatus struct{}

type SM_GuildInvite struct {
	Name string
}

type SM_GuildExpGain struct{}

type SM_GuildNameRequest struct{}

type SM_GuildStorageGoldChange struct{}

type SM_GuildStorageItemChange struct{}

type SM_GuildStorageList struct{}

type SM_GuildRequestWar struct{}

type SM_DefaultNPC struct{}

type SM_NPCUpdate struct {
	NPCID uint32
}

type SM_NPCImageUpdate struct{}

type SM_MarriageRequest struct{}

type SM_DivorceRequest struct{}

type SM_MentorRequest struct{}

type SM_TradeRequest struct {
	Name string
}

type SM_TradeAccept struct {
	Name string
}

type SM_TradeGold struct {
	Amount uint32
}

type SM_TradeItem struct {
	TradeItems []*UserItem
}

type SM_TradeConfirm struct{}

type SM_TradeCancel struct {
	Unlock bool
}

type SM_MountUpdate struct{}

type SM_EquipSlotItem struct{}

type SM_FishingUpdate struct{}

type SM_ChangeQuest struct{}

type SM_CompleteQuest struct{}

type SM_ShareQuest struct{}

type SM_NewQuestInfo struct{}

type SM_GainedQuestItem struct{}

type SM_DeleteQuestItem struct{}

type SM_CancelReincarnation struct{}

type SM_RequestReincarnation struct{}

type SM_UserBackStep struct{}

type SM_ObjectBackStep struct{}

type SM_UserDashAttack struct{}

type SM_ObjectDashAttack struct{}

type SM_UserAttackMove struct{}

type SM_CombineItem struct{}

type SM_ItemUpgraded struct{}

type SM_SetConcentration struct {
	ObjectID    uint32
	Enabled     bool
	Interrupted bool
}

type SM_SetObjectConcentration struct {
	ObjectID    uint32
	Enabled     bool
	Interrupted bool
}

type SM_SetElemental struct{}
type SM_SetObjectElemental struct{}
type SM_RemoveDelayedExplosion struct{}
type SM_ObjectDeco struct{}
type SM_ObjectSneaking struct{}
type SM_ObjectLevelEffects struct{}
type SM_SetBindingShot struct{}
type SM_SendOutputMessage struct{}
type SM_NPCAwakening struct{}
type SM_NPCDisassemble struct{}
type SM_NPCDowngrade struct{}
type SM_NPCReset struct{}
type SM_AwakeningNeedMaterials struct{}
type SM_AwakeningLockedItem struct{}
type SM_Awakening struct{}
type SM_ReceiveMail struct{}
type SM_MailLockedItem struct{}
type SM_MailSendRequest struct{}
type SM_MailSent struct{}
type SM_ParcelCollected struct{}
type SM_MailCost struct{}
type SM_ResizeInventory struct{}
type SM_ResizeStorage struct{}
type SM_NewIntelligentCreature struct{}
type SM_UpdateIntelligentCreatureList struct{}
type SM_IntelligentCreatureEnableRename struct{}
type SM_IntelligentCreaturePickup struct{}
type SM_NPCPearlGoods struct{}
type SM_TransformUpdate struct{}
type SM_FriendUpdate struct{}
type SM_LoverUpdate struct{}
type SM_MentorUpdate struct{}
type SM_GuildBuffList struct{}
type SM_NPCRequestInput struct{}
type SM_GameShopInfo struct{}
type SM_GameShopStock struct{}
type SM_Rankings struct{}
type SM_Opendoor struct {
	DoorIndex byte
	Close     bool
}
type SM_GetRentedItems struct{}
type SM_ItemRentalRequest struct{}
type SM_ItemRentalFee struct{}
type SM_ItemRentalPeriod struct{}
type SM_DepositRentalItem struct{}
type SM_RetrieveRentalItem struct{}
type SM_UpdateRentalItem struct{}
type SM_CancelItemRental struct{}
type SM_ItemRentalLock struct{}
type SM_ItemRentalPartnerLock struct{}
type SM_CanConfirmItemRental struct{}
type SM_ConfirmItemRental struct{}
type SM_NewRecipeInfo struct{}
type SM_OpenBrowser struct{}
