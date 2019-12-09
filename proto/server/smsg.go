package server

import (
	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/codec/mircodec"
	"github.com/yenkeia/mirgo/common"

	// 使用binary协议，因此匿名引用这个包，底层会自动注册
	"reflect"

	_ "github.com/davyxu/cellnet/codec/binary"
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
	/*
	 * 0: Disabled
	 * 1: Bad AccountID
	 * 2: Bad Password
	 * 3: Bad Email
	 * 4: Bad Name
	 * 5: Bad Question
	 * 6: Bad Answer
	 * 7: Account Exists.
	 * 8: Success
	 */
}

type ChangePassword struct {
	Result uint8
	/*
	 * 0: Disabled
	 * 1: Bad AccountID
	 * 2: Bad Current Password
	 * 3: Bad New Password
	 * 4: Account Not Exist
	 * 5: Wrong Password
	 * 6: Success
	 */
}

// TODO
type ChangePasswordBanned struct {
	//public string Reason = string.Empty;
	//public DateTime ExpiryDate;
}

type Login struct {
	Result uint8
	/*
	* 0: Disabled
	* 1: Bad AccountID
	* 2: Bad Password
	* 3: Account Not Exist
	* 4: Wrong Password
	 */
}

// TODO
type LoginBanned struct {
	//public string Reason = string.Empty;
	//public DateTime ExpiryDate
}

type LoginSuccess struct {
	Characters []common.SelectInfo
}

type NewCharacter struct {
	Result uint8
	/*
	 * 0: Disabled.
	 * 1: Bad Character Name
	 * 2: Bad Gender
	 * 3: Bad Class
	 * 4: Max Characters
	 * 5: Character Exists.
	 * */
}

type NewCharacterSuccess struct {
	CharInfo common.SelectInfo
}

type DeleteCharacter struct {
	Result uint8
	/*
	 * 0: Disabled.
	 * 1: Character Not Found
	 * */
}

type DeleteCharacterSuccess struct {
	CharacterIndex int16
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
	Music        uint16
	Lights       common.LightSetting
	Lightning    bool
	MapDarkLight uint8
}

// TODO
type UserInformation struct {
	Bytes []byte
}

// TODO
type UserInformation_tmp struct {
	ObjectID                  uint32
	RealId                    uint32
	Name                      string
	GuildName                 string
	GuildRank                 string
	NameColour                uint32
	Class                     common.MirClass
	Gender                    common.MirGender
	Level                     uint16
	Location                  common.Point
	Direction                 common.MirDirection
	Hair                      uint8
	HP                        uint16
	MP                        uint16
	Experience                uint64
	MaxExperience             uint64
	LevelEffect               uint8 // LevelEffects
	Inventory                 []common.UserItem
	Equipment                 []common.UserItem
	QuestInventory            []common.UserItem
	Gold                      uint32
	Credit                    uint32
	HasExpandedStorage        bool
	ExpandedStorageExpiryTime uint64 // DateTime
	//Magics                    interface{} // []ClientMagic
	//IntelligentCreatures      interface{} // []ClientIntelligentCreature // TODO
	//IntelligentCreatureType   uint8       // IntelligentCreatureType 忽略
	//CreatureSummoned          bool
}

type UserLocation struct{}
type ObjectPlayer struct{}
type ObjectRemove struct{}
type ObjectTurn struct{}
type ObjectWalk struct{}
type ObjectRun struct{}
type Chat struct{}
type ObjectChat struct{}

type NewItemInfo struct {
	Info common.ItemInfo
}

type MoveItem struct{}
type EquipItem struct{}
type MergeItem struct{}
type RemoveItem struct{}
type RemoveSlotItem struct{}
type TakeBackItem struct{}
type StoreItem struct{}
type SplitItem struct{}
type SplitItem1 struct{}
type DepositRefineItem struct{}
type RetrieveRefineItem struct{}
type RefineCancel struct{}
type RefineItem struct{}
type DepositTradeItem struct{}
type RetrieveTradeItem struct{}
type UseItem struct{}
type DropItem struct{}
type PlayerUpdate struct{}
type PlayerInspect struct{}
type LogOutSuccess struct{}
type LogOutFailed struct{}
type TimeOfDay struct{}
type ChangeAMode struct{}
type ChangePMode struct{}
type ObjectItem struct{}
type ObjectGold struct{}
type GainedItem struct{}
type GainedGold struct{}
type LoseGold struct{}
type GainedCredit struct{}
type LoseCredit struct{}
type ObjectMonster struct{}
type ObjectAttack struct{}
type Struck struct{}
type ObjectStruck struct{}
type DamageIndicator struct{}
type DuraChanged struct{}
type HealthChanged struct{}
type DeleteItem struct{}
type Death struct{}
type ObjectDied struct{}
type ColourChanged struct{}
type ObjectColourChanged struct{}
type ObjectGuildNameChanged struct{}
type GainExperience struct{}
type LevelChanged struct{}
type ObjectLeveled struct{}
type ObjectHarvest struct{}
type ObjectHarvested struct{}
type ObjectNpc struct{}
type NPCResponse struct{}
type ObjectHide struct{}
type ObjectShow struct{}
type Poisoned struct{}
type ObjectPoisoned struct{}
type MapChanged struct{}
type ObjectTeleportOut struct{}
type ObjectTeleportIn struct{}
type TeleportIn struct{}
type NPCGoods struct{}
type NPCSell struct{}
type NPCRepair struct{}
type NPCSRepair struct{}
type NPCRefine struct{}
type NPCCheckRefine struct{}
type NPCCollectRefine struct{}
type NPCReplaceWedRing struct{}
type NPCStorage struct{}
type SellItem struct{}
type CraftItem struct{}
type RepairItem struct{}
type ItemRepaired struct{}
type NewMagic struct{}
type RemoveMagic struct{}
type MagicLeveled struct{}
type Magic struct{}
type MagicDelay struct{}
type MagicCast struct{}
type ObjectMagic struct{}
type ObjectEffect struct{}
type RangeAttack struct{}
type Pushed struct{}
type ObjectPushed struct{}
type ObjectName struct{}
type UserStorage struct{}
type SwitchGroup struct{}
type DeleteGroup struct{}
type DeleteMember struct{}
type GroupInvite struct{}
type AddMember struct{}
type Revived struct{}
type ObjectRevived struct{}
type SpellToggle struct{}
type ObjectHealth struct{}
type MapEffect struct{}
type ObjectRangeAttack struct{}
type AddBuff struct{}
type RemoveBuff struct{}
type ObjectHidden struct{}
type RefreshItem struct{}
type ObjectSpell struct{}
type UserDash struct{}
type ObjectDash struct{}
type UserDashFail struct{}
type ObjectDashFail struct{}
type NPCConsign struct{}
type NPCMarket struct{}
type NPCMarketPage struct{}
type ConsignItem struct{}
type MarketFail struct{}
type MarketSuccess struct{}
type ObjectSitDown struct{}
type InTrapRock struct{}
type BaseStatsInfo struct{}
type UserName struct{}
type ChatItemStats struct{}
type GuildNoticeChange struct{}
type GuildMemberChange struct{}
type GuildStatus struct{}
type GuildInvite struct{}
type GuildExpGain struct{}
type GuildNameRequest struct{}
type GuildStorageGoldChange struct{}
type GuildStorageItemChange struct{}
type GuildStorageList struct{}
type GuildRequestWar struct{}
type DefaultNPC struct{}
type NPCUpdate struct{}
type NPCImageUpdate struct{}
type MarriageRequest struct{}
type DivorceRequest struct{}
type MentorRequest struct{}
type TradeRequest struct{}
type TradeAccept struct{}
type TradeGold struct{}
type TradeItem struct{}
type TradeConfirm struct{}
type TradeCancel struct{}
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

// 引用消息时，自动注册消息，这个文件可以由代码生成自动生成
func init() {

	mirCodec := new(mircodec.MirCodec)

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Connected)(nil)).Elem(),
		ID:    CONNECTED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ClientVersion)(nil)).Elem(),
		ID:    CLIENT_VERSION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Disconnect)(nil)).Elem(),
		ID:    DISCONNECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*KeepAlive)(nil)).Elem(),
		ID:    KEEP_ALIVE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NewAccount)(nil)).Elem(),
		ID:    NEW_ACCOUNT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ChangePassword)(nil)).Elem(),
		ID:    CHANGE_PASSWORD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ChangePasswordBanned)(nil)).Elem(),
		ID:    CHANGE_PASSWORD_BANNED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Login)(nil)).Elem(),
		ID:    LOGIN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*LoginBanned)(nil)).Elem(),
		ID:    LOGIN_BANNED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*LoginSuccess)(nil)).Elem(),
		ID:    LOGIN_SUCCESS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NewCharacter)(nil)).Elem(),
		ID:    NEW_CHARACTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NewCharacterSuccess)(nil)).Elem(),
		ID:    NEW_CHARACTER_SUCCESS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DeleteCharacter)(nil)).Elem(),
		ID:    DELETE_CHARACTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DeleteCharacterSuccess)(nil)).Elem(),
		ID:    DELETE_CHARACTER_SUCCESS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*StartGame)(nil)).Elem(),
		ID:    START_GAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*StartGameBanned)(nil)).Elem(),
		ID:    START_GAME_BANNED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*StartGameDelay)(nil)).Elem(),
		ID:    START_GAME_DELAY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MapInformation)(nil)).Elem(),
		ID:    MAP_INFORMATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*UserInformation)(nil)).Elem(),
		ID:    USER_INFORMATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*UserLocation)(nil)).Elem(),
		ID:    USER_LOCATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectPlayer)(nil)).Elem(),
		ID:    OBJECT_PLAYER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectRemove)(nil)).Elem(),
		ID:    OBJECT_REMOVE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectTurn)(nil)).Elem(),
		ID:    OBJECT_TURN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectWalk)(nil)).Elem(),
		ID:    OBJECT_WALK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectRun)(nil)).Elem(),
		ID:    OBJECT_RUN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Chat)(nil)).Elem(),
		ID:    CHAT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectChat)(nil)).Elem(),
		ID:    OBJECT_CHAT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NewItemInfo)(nil)).Elem(),
		ID:    NEW_ITEM_INFO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MoveItem)(nil)).Elem(),
		ID:    MOVE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*EquipItem)(nil)).Elem(),
		ID:    EQUIP_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MergeItem)(nil)).Elem(),
		ID:    MERGE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RemoveItem)(nil)).Elem(),
		ID:    REMOVE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RemoveSlotItem)(nil)).Elem(),
		ID:    REMOVE_SLOT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*TakeBackItem)(nil)).Elem(),
		ID:    TAKE_BACK_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*StoreItem)(nil)).Elem(),
		ID:    STORE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SplitItem)(nil)).Elem(),
		ID:    SPLIT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SplitItem1)(nil)).Elem(),
		ID:    SPLIT_ITEM1,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DepositRefineItem)(nil)).Elem(),
		ID:    DEPOSIT_REFINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RetrieveRefineItem)(nil)).Elem(),
		ID:    RETRIEVE_REFINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RefineCancel)(nil)).Elem(),
		ID:    REFINE_CANCEL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RefineItem)(nil)).Elem(),
		ID:    REFINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DepositTradeItem)(nil)).Elem(),
		ID:    DEPOSIT_TRADE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RetrieveTradeItem)(nil)).Elem(),
		ID:    RETRIEVE_TRADE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*UseItem)(nil)).Elem(),
		ID:    USE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DropItem)(nil)).Elem(),
		ID:    DROP_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*PlayerUpdate)(nil)).Elem(),
		ID:    PLAYER_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*PlayerInspect)(nil)).Elem(),
		ID:    PLAYER_INSPECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*LogOutSuccess)(nil)).Elem(),
		ID:    LOG_OUT_SUCCESS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*LogOutFailed)(nil)).Elem(),
		ID:    LOG_OUT_FAILED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*TimeOfDay)(nil)).Elem(),
		ID:    TIME_OF_DAY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ChangeAMode)(nil)).Elem(),
		ID:    CHANGE_A_MODE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ChangePMode)(nil)).Elem(),
		ID:    CHANGE_P_MODE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectItem)(nil)).Elem(),
		ID:    OBJECT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectGold)(nil)).Elem(),
		ID:    OBJECT_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GainedItem)(nil)).Elem(),
		ID:    GAINED_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GainedGold)(nil)).Elem(),
		ID:    GAINED_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*LoseGold)(nil)).Elem(),
		ID:    LOSE_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GainedCredit)(nil)).Elem(),
		ID:    GAINED_CREDIT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*LoseCredit)(nil)).Elem(),
		ID:    LOSE_CREDIT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectMonster)(nil)).Elem(),
		ID:    OBJECT_MONSTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectAttack)(nil)).Elem(),
		ID:    OBJECT_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Struck)(nil)).Elem(),
		ID:    STRUCK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectStruck)(nil)).Elem(),
		ID:    OBJECT_STRUCK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DamageIndicator)(nil)).Elem(),
		ID:    DAMAGE_INDICATOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DuraChanged)(nil)).Elem(),
		ID:    DURA_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*HealthChanged)(nil)).Elem(),
		ID:    HEALTH_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DeleteItem)(nil)).Elem(),
		ID:    DELETE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Death)(nil)).Elem(),
		ID:    DEATH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectDied)(nil)).Elem(),
		ID:    OBJECT_DIED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ColourChanged)(nil)).Elem(),
		ID:    COLOUR_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectColourChanged)(nil)).Elem(),
		ID:    OBJECT_COLOUR_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectGuildNameChanged)(nil)).Elem(),
		ID:    OBJECT_GUILD_NAME_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GainExperience)(nil)).Elem(),
		ID:    GAIN_EXPERIENCE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*LevelChanged)(nil)).Elem(),
		ID:    LEVEL_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectLeveled)(nil)).Elem(),
		ID:    OBJECT_LEVELED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectHarvest)(nil)).Elem(),
		ID:    OBJECT_HARVEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectHarvested)(nil)).Elem(),
		ID:    OBJECT_HARVESTED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectNpc)(nil)).Elem(),
		ID:    OBJECT_NPC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NPCResponse)(nil)).Elem(),
		ID:    NPC_RESPONSE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectHide)(nil)).Elem(),
		ID:    OBJECT_HIDE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectShow)(nil)).Elem(),
		ID:    OBJECT_SHOW,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Poisoned)(nil)).Elem(),
		ID:    POISONED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectPoisoned)(nil)).Elem(),
		ID:    OBJECT_POISONED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MapChanged)(nil)).Elem(),
		ID:    MAP_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectTeleportOut)(nil)).Elem(),
		ID:    OBJECT_TELEPORT_OUT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectTeleportIn)(nil)).Elem(),
		ID:    OBJECT_TELEPORT_IN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*TeleportIn)(nil)).Elem(),
		ID:    TELEPORT_IN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NPCGoods)(nil)).Elem(),
		ID:    NPC_GOODS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NPCSell)(nil)).Elem(),
		ID:    NPC_SELL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NPCRepair)(nil)).Elem(),
		ID:    NPC_REPAIR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NPCSRepair)(nil)).Elem(),
		ID:    NPC_S_REPAIR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NPCRefine)(nil)).Elem(),
		ID:    NPC_REFINE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NPCCheckRefine)(nil)).Elem(),
		ID:    NPC_CHECK_REFINE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NPCCollectRefine)(nil)).Elem(),
		ID:    NPC_COLLECT_REFINE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NPCReplaceWedRing)(nil)).Elem(),
		ID:    NPC_REPLACE_WED_RING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NPCStorage)(nil)).Elem(),
		ID:    NPC_STORAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SellItem)(nil)).Elem(),
		ID:    SELL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CraftItem)(nil)).Elem(),
		ID:    CRAFT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RepairItem)(nil)).Elem(),
		ID:    REPAIR_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ItemRepaired)(nil)).Elem(),
		ID:    ITEM_REPAIRED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NewMagic)(nil)).Elem(),
		ID:    NEW_MAGIC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RemoveMagic)(nil)).Elem(),
		ID:    REMOVE_MAGIC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MagicLeveled)(nil)).Elem(),
		ID:    MAGIC_LEVELED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Magic)(nil)).Elem(),
		ID:    MAGIC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MagicDelay)(nil)).Elem(),
		ID:    MAGIC_DELAY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MagicCast)(nil)).Elem(),
		ID:    MAGIC_CAST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectMagic)(nil)).Elem(),
		ID:    OBJECT_MAGIC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectEffect)(nil)).Elem(),
		ID:    OBJECT_EFFECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RangeAttack)(nil)).Elem(),
		ID:    RANGE_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Pushed)(nil)).Elem(),
		ID:    PUSHED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectPushed)(nil)).Elem(),
		ID:    OBJECT_PUSHED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectName)(nil)).Elem(),
		ID:    OBJECT_NAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*UserStorage)(nil)).Elem(),
		ID:    USER_STORAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SwitchGroup)(nil)).Elem(),
		ID:    SWITCH_GROUP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DeleteGroup)(nil)).Elem(),
		ID:    DELETE_GROUP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DeleteMember)(nil)).Elem(),
		ID:    DELETE_MEMBER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GroupInvite)(nil)).Elem(),
		ID:    GROUP_INVITE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*AddMember)(nil)).Elem(),
		ID:    ADD_MEMBER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Revived)(nil)).Elem(),
		ID:    REVIVED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectRevived)(nil)).Elem(),
		ID:    OBJECT_REVIVED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SpellToggle)(nil)).Elem(),
		ID:    SPELL_TOGGLE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectHealth)(nil)).Elem(),
		ID:    OBJECT_HEALTH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MapEffect)(nil)).Elem(),
		ID:    MAP_EFFECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectRangeAttack)(nil)).Elem(),
		ID:    OBJECT_RANGE_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*AddBuff)(nil)).Elem(),
		ID:    ADD_BUFF,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RemoveBuff)(nil)).Elem(),
		ID:    REMOVE_BUFF,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectHidden)(nil)).Elem(),
		ID:    OBJECT_HIDDEN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RefreshItem)(nil)).Elem(),
		ID:    REFRESH_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectSpell)(nil)).Elem(),
		ID:    OBJECT_SPELL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*UserDash)(nil)).Elem(),
		ID:    USER_DASH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectDash)(nil)).Elem(),
		ID:    OBJECT_DASH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*UserDashFail)(nil)).Elem(),
		ID:    USER_DASH_FAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectDashFail)(nil)).Elem(),
		ID:    OBJECT_DASH_FAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NPCConsign)(nil)).Elem(),
		ID:    NPC_CONSIGN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NPCMarket)(nil)).Elem(),
		ID:    NPC_MARKET,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NPCMarketPage)(nil)).Elem(),
		ID:    NPC_MARKET_PAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ConsignItem)(nil)).Elem(),
		ID:    CONSIGN_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MarketFail)(nil)).Elem(),
		ID:    MARKET_FAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MarketSuccess)(nil)).Elem(),
		ID:    MARKET_SUCCESS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectSitDown)(nil)).Elem(),
		ID:    OBJECT_SIT_DOWN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*InTrapRock)(nil)).Elem(),
		ID:    IN_TRAP_ROCK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*BaseStatsInfo)(nil)).Elem(),
		ID:    BASE_STATS_INFO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*UserName)(nil)).Elem(),
		ID:    USER_NAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ChatItemStats)(nil)).Elem(),
		ID:    CHAT_ITEM_STATS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GuildNoticeChange)(nil)).Elem(),
		ID:    GUILD_NOTICE_CHANGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GuildMemberChange)(nil)).Elem(),
		ID:    GUILD_MEMBER_CHANGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GuildStatus)(nil)).Elem(),
		ID:    GUILD_STATUS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GuildInvite)(nil)).Elem(),
		ID:    GUILD_INVITE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GuildExpGain)(nil)).Elem(),
		ID:    GUILD_EXP_GAIN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GuildNameRequest)(nil)).Elem(),
		ID:    GUILD_NAME_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GuildStorageGoldChange)(nil)).Elem(),
		ID:    GUILD_STORAGE_GOLD_CHANGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GuildStorageItemChange)(nil)).Elem(),
		ID:    GUILD_STORAGE_ITEM_CHANGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GuildStorageList)(nil)).Elem(),
		ID:    GUILD_STORAGE_LIST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GuildRequestWar)(nil)).Elem(),
		ID:    GUILD_REQUEST_WAR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DefaultNPC)(nil)).Elem(),
		ID:    DEFAULT_NPC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NPCUpdate)(nil)).Elem(),
		ID:    NPC_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NPCImageUpdate)(nil)).Elem(),
		ID:    NPC_IMAGE_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MarriageRequest)(nil)).Elem(),
		ID:    MARRIAGE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DivorceRequest)(nil)).Elem(),
		ID:    DIVORCE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MentorRequest)(nil)).Elem(),
		ID:    MENTOR_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*TradeRequest)(nil)).Elem(),
		ID:    TRADE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*TradeAccept)(nil)).Elem(),
		ID:    TRADE_ACCEPT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*TradeGold)(nil)).Elem(),
		ID:    TRADE_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*TradeItem)(nil)).Elem(),
		ID:    TRADE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*TradeConfirm)(nil)).Elem(),
		ID:    TRADE_CONFIRM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*TradeCancel)(nil)).Elem(),
		ID:    TRADE_CANCEL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MountUpdate)(nil)).Elem(),
		ID:    MOUNT_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*EquipSlotItem)(nil)).Elem(),
		ID:    EQUIP_SLOT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*FishingUpdate)(nil)).Elem(),
		ID:    FISHING_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ChangeQuest)(nil)).Elem(),
		ID:    CHANGE_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CompleteQuest)(nil)).Elem(),
		ID:    COMPLETE_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ShareQuest)(nil)).Elem(),
		ID:    SHARE_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NewQuestInfo)(nil)).Elem(),
		ID:    NEW_QUEST_INFO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GainedQuestItem)(nil)).Elem(),
		ID:    GAINED_QUEST_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DeleteQuestItem)(nil)).Elem(),
		ID:    DELETE_QUEST_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CancelReincarnation)(nil)).Elem(),
		ID:    CANCEL_REINCARNATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RequestReincarnation)(nil)).Elem(),
		ID:    REQUEST_REINCARNATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*UserBackStep)(nil)).Elem(),
		ID:    USER_BACK_STEP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectBackStep)(nil)).Elem(),
		ID:    OBJECT_BACK_STEP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*UserDashAttack)(nil)).Elem(),
		ID:    USER_DASH_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ObjectDashAttack)(nil)).Elem(),
		ID:    OBJECT_DASH_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*UserAttackMove)(nil)).Elem(),
		ID:    USER_ATTACK_MOVE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CombineItem)(nil)).Elem(),
		ID:    COMBINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ItemUpgraded)(nil)).Elem(),
		ID:    ITEM_UPGRADED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SetConcentration)(nil)).Elem(),
		ID:    SET_CONCENTRATION,
	})
}
