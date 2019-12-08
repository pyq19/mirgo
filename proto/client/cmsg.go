package client

import (
	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/codec/mircodec"
	"github.com/yenkeia/mirgo/common"

	// 使用binary协议，因此匿名引用这个包，底层会自动注册
	"reflect"

	_ "github.com/davyxu/cellnet/codec/binary"
)

const (
	CLIENT_VERSION = 1000 + iota
	DISCONNECT
	KEEP_ALIVE
	NEW_ACCOUNT
	CHANGE_PASSWORD
	LOGIN
	NEW_CHARACTER
	DELETE_CHARACTER
	START_GAME
	LOG_OUT
	TURN
	WALK
	RUN
	CHAT
	MOVE_ITEM
	STORE_ITEM
	TAKE_BACK_ITEM
	MERGE_ITEM
	EQUIP_ITEM
	REMOVE_ITEM
	REMOVE_SLOT_ITEM
	SPLIT_ITEM
	USE_ITEM
	DROP_ITEM
	DEPOSIT_REFINE_ITEM
	RETRIEVE_REFINE_ITEM
	REFINE_CANCEL
	REFINE_ITEM
	CHECK_REFINE
	REPLACE_WED_RING
	DEPOSIT_TRADE_ITEM
	RETRIEVE_TRADE_ITEM
	DROP_GOLD
	PICK_UP
	INSPECT
	CHANGE_A_MODE
	CHANGE_P_MODE
	CHANGE_TRADE
	ATTACK
	RANGE_ATTACK
	HARVEST
	CALL_NPC
	TALK_MONSTER_NPC
	BUY_ITEM
	SELL_ITEM
	CRAFT_ITEM
	REPAIR_ITEM
	BUY_ITEM_BACK
	S_REPAIR_ITEM
	MAGIC_KEY
	MAGIC
)

type ClientVersion struct {
	VersionHash []uint8
}

type Disconnect struct{}

type KeepAlive struct {
	Time int64
}

type NewAccount struct {
	AccountID      string
	Password       string
	DateTime       int64  // 无用字段 c# 中 DateTime 8 字节
	UserName       string // 无用字段
	SecretQuestion string // 无用字段
	SecretAnswer   string // 无用字段
	EMailAddress   string // 无用字段
}

type ChangePassword struct {
	AccountID       string
	CurrentPassword string
	NewPassword     string
}

type Login struct {
	AccountID string
	Password  string
}

type NewCharacter struct {
	Name   string
	Gender common.MirGender
	Class  common.MirClass
}

type DeleteCharacter struct {
	CharacterIndex int16
}

type StartGame struct {
	CharacterIndex int16
}

type LogOut struct{}
type Turn struct{}
type Walk struct{}
type Run struct{}
type Chat struct{}
type MoveItem struct{}
type StoreItem struct{}
type TakeBackItem struct{}
type MergeItem struct{}
type EquipItem struct{}
type RemoveItem struct{}
type RemoveSlotItem struct{}
type SplitItem struct{}
type UseItem struct{}
type DropItem struct{}
type DepositRefineItem struct{}
type RetrieveRefineItem struct{}
type RefineCancel struct{}
type RefineItem struct{}
type CheckRefine struct{}
type ReplaceWedRing struct{}
type DepositTradeItem struct{}
type RetrieveTradeItem struct{}
type DropGold struct{}
type PickUp struct{}
type Inspect struct{}
type ChangeAMode struct{}
type ChangePMode struct{}
type ChangeTrade struct{}
type Attack struct{}
type RangeAttack struct{}
type Harvest struct{}
type CallNPC struct{}
type TalkMonsterNPC struct{}
type BuyItem struct{}
type SellItem struct{}
type CraftItem struct{}
type RepairItem struct{}
type BuyItemBack struct{}
type SRepairItem struct{}
type MagicKey struct{}
type Magic struct{}

// 引用消息时，自动注册消息，这个文件可以由代码生成自动生成
func init() {

	mirCodec := new(mircodec.MirCodec)

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
		Type:  reflect.TypeOf((*Login)(nil)).Elem(),
		ID:    LOGIN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NewCharacter)(nil)).Elem(),
		ID:    NEW_CHARACTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DeleteCharacter)(nil)).Elem(),
		ID:    DELETE_CHARACTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*StartGame)(nil)).Elem(),
		ID:    START_GAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*LogOut)(nil)).Elem(),
		ID:    LOG_OUT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Turn)(nil)).Elem(),
		ID:    TURN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Walk)(nil)).Elem(),
		ID:    WALK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Run)(nil)).Elem(),
		ID:    RUN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Chat)(nil)).Elem(),
		ID:    CHAT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MoveItem)(nil)).Elem(),
		ID:    MOVE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*StoreItem)(nil)).Elem(),
		ID:    STORE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*TakeBackItem)(nil)).Elem(),
		ID:    TAKE_BACK_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MergeItem)(nil)).Elem(),
		ID:    MERGE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*EquipItem)(nil)).Elem(),
		ID:    EQUIP_ITEM,
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
		Type:  reflect.TypeOf((*SplitItem)(nil)).Elem(),
		ID:    SPLIT_ITEM,
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
		Type:  reflect.TypeOf((*CheckRefine)(nil)).Elem(),
		ID:    CHECK_REFINE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ReplaceWedRing)(nil)).Elem(),
		ID:    REPLACE_WED_RING,
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
		Type:  reflect.TypeOf((*DropGold)(nil)).Elem(),
		ID:    DROP_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*PickUp)(nil)).Elem(),
		ID:    PICK_UP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Inspect)(nil)).Elem(),
		ID:    INSPECT,
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
		Type:  reflect.TypeOf((*ChangeTrade)(nil)).Elem(),
		ID:    CHANGE_TRADE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Attack)(nil)).Elem(),
		ID:    ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RangeAttack)(nil)).Elem(),
		ID:    RANGE_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Harvest)(nil)).Elem(),
		ID:    HARVEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CallNPC)(nil)).Elem(),
		ID:    CALL_NPC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*TalkMonsterNPC)(nil)).Elem(),
		ID:    TALK_MONSTER_NPC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*BuyItem)(nil)).Elem(),
		ID:    BUY_ITEM,
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
		Type:  reflect.TypeOf((*BuyItemBack)(nil)).Elem(),
		ID:    BUY_ITEM_BACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SRepairItem)(nil)).Elem(),
		ID:    S_REPAIR_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MagicKey)(nil)).Elem(),
		ID:    MAGIC_KEY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Magic)(nil)).Elem(),
		ID:    MAGIC,
	})
}
