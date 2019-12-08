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

	Resolution int16
}

type StartGameBanned struct{}
type StartGameDelay struct{}
type MapInformation struct{}
type UserInformation struct{}
type UserLocation struct{}
type ObjectPlayer struct{}
type ObjectRemove struct{}
type ObjectTurn struct{}
type ObjectWalk struct{}
type ObjectRun struct{}
type Chat struct{}
type ObjectChat struct{}
type NewItemInfo struct{}
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
}
