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
	SWITCH_GROUP
	ADD_MEMBER
	DELL_MEMBER
	GROUP_INVITE
	TOWN_REVIVE
	SPELL_TOGGLE
	CONSIGN_ITEM
	MARKET_SEARCH
	MARKET_REFRESH
	MARKET_PAGE
	MARKET_BUY
	MARKET_GET_BACK
	REQUEST_USER_NAME
	REQUEST_CHAT_ITEM
	EDIT_GUILD_MEMBER
	EDIT_GUILD_NOTICE
	GUILD_INVITE
	GUILD_NAME_RETURN
	REQUEST_GUILD_INFO
	GUILD_STORAGE_GOLD_CHANGE
	GUILD_STORAGE_ITEM_CHANGE
	GUILD_WAR_RETURN
	MARRIAGE_REQUEST
	MARRIAGE_REPLY
	CHANGE_MARRIAGE
	DIVORCE_REQUEST
	DIVORCE_REPLY
	ADD_MENTOR
	MENTOR_REPLY
	ALLOW_MENTOR
	CANCEL_MENTOR
	TRADE_REQUEST
	TRADE_REPLY
	TRADE_GOLD
	TRADE_CONFIRM
	TRADE_CANCEL
	EQUIP_SLOT_ITEM
	FISHING_CAST
	FISHING_CHANGE_AUTOCAST
	ACCEPT_QUEST
	FINISH_QUEST
	ABANDON_QUEST
	SHARE_QUEST
	ACCEPT_REINCARNATION
	CANCEL_REINCARNATION
	COMBINE_ITEM
	SET_CONCENTRATION
	AWAKENING_NEED_MATERIALS
	AWAKENING_LOCKED_ITEM
	AWAKENING
	DISASSEMBLE_ITEM
	DOWNGRADE_AWAKENING
	RESET_ADDED_ITEM
	SEND_MAIL
	READ_MAIL
	COLLECT_PARCEL
	DELETE_MAIL
	LOCK_MAIL
	MAIL_LOCKED_ITEM
	MAIL_COST
	UPDATE_INTELLIGENT_CREATURE
	INTELLIGENT_CREATURE_PICKUP
	ADD_FRIEND
	REMOVE_FRIEND
	REFRESH_FRIENDS
	ADD_MEMO
	GUILD_BUFF_UPDATE
	NPC_CONFIRM_INPUT
	GAMESHOP_BUY
	REPORT_ISSUE
	GET_RANKING
	OPENDOOR
	GET_RENTED_ITEMS
	ITEM_RENTAL_REQUEST
	ITEM_RENTAL_FEE
	ITEM_RENTAL_PERIOD
	DEPOSIT_RENTAL_ITEM
	RETRIEVE_RENTAL_ITEM
	CANCEL_ITEM_RENTAL
	ITEM_RENTAL_LOCK_FEE
	ITEM_RENTAL_LOCK_ITEM
	CONFIRM_ITEM_RENTAL
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
type SwitchGroup struct{}
type AddMember struct{}
type DellMember struct{}
type GroupInvite struct{}
type TownRevive struct{}
type SpellToggle struct{}
type ConsignItem struct{}
type MarketSearch struct{}
type MarketRefresh struct{}
type MarketPage struct{}
type MarketBuy struct{}
type MarketGetBack struct{}
type RequestUserName struct{}
type RequestChatItem struct{}
type EditGuildMember struct{}
type EditGuildNotice struct{}
type GuildInvite struct{}
type GuildNameReturn struct{}
type RequestGuildInfo struct{}
type GuildStorageGoldChange struct{}
type GuildStorageItemChange struct{}
type GuildWarReturn struct{}
type MarriageRequest struct{}
type MarriageReply struct{}
type ChangeMarriage struct{}
type DivorceRequest struct{}
type DivorceReply struct{}
type AddMentor struct{}
type MentorReply struct{}
type AllowMentor struct{}
type CancelMentor struct{}
type TradeRequest struct{}
type TradeReply struct{}
type TradeGold struct{}
type TradeConfirm struct{}
type TradeCancel struct{}
type EquipSlotItem struct{}
type FishingCast struct{}
type FishingChangeAutocast struct{}
type AcceptQuest struct{}
type FinishQuest struct{}
type AbandonQuest struct{}
type ShareQuest struct{}
type AcceptReincarnation struct{}
type CancelReincarnation struct{}
type CombineItem struct{}
type SetConcentration struct{}
type AwakeningNeedMaterials struct{}
type AwakeningLockedItem struct{}
type Awakening struct{}
type DisassembleItem struct{}
type DowngradeAwakening struct{}
type ResetAddedItem struct{}
type SendMail struct{}
type ReadMail struct{}
type CollectParcel struct{}
type DeleteMail struct{}
type LockMail struct{}
type MailLockedItem struct{}
type MailCost struct{}
type UpdateIntelligentCreature struct{}
type IntelligentCreaturePickup struct{}
type AddFriend struct{}
type RemoveFriend struct{}
type RefreshFriends struct{}
type AddMemo struct{}
type GuildBuffUpdate struct{}
type NPCConfirmInput struct{}
type GameshopBuy struct{}
type ReportIssue struct{}
type GetRanking struct{}
type Opendoor struct{}
type GetRentedItems struct{}
type ItemRentalRequest struct{}
type ItemRentalFee struct{}
type ItemRentalPeriod struct{}
type DepositRentalItem struct{}
type RetrieveRentalItem struct{}
type CancelItemRental struct{}
type ItemRentalLockFee struct{}
type ItemRentalLockItem struct{}
type ConfirmItemRental struct{}

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
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SwitchGroup)(nil)).Elem(),
		ID:    SWITCH_GROUP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*AddMember)(nil)).Elem(),
		ID:    ADD_MEMBER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DellMember)(nil)).Elem(),
		ID:    DELL_MEMBER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GroupInvite)(nil)).Elem(),
		ID:    GROUP_INVITE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*TownRevive)(nil)).Elem(),
		ID:    TOWN_REVIVE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SpellToggle)(nil)).Elem(),
		ID:    SPELL_TOGGLE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ConsignItem)(nil)).Elem(),
		ID:    CONSIGN_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MarketSearch)(nil)).Elem(),
		ID:    MARKET_SEARCH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MarketRefresh)(nil)).Elem(),
		ID:    MARKET_REFRESH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MarketPage)(nil)).Elem(),
		ID:    MARKET_PAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MarketBuy)(nil)).Elem(),
		ID:    MARKET_BUY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MarketGetBack)(nil)).Elem(),
		ID:    MARKET_GET_BACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RequestUserName)(nil)).Elem(),
		ID:    REQUEST_USER_NAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RequestChatItem)(nil)).Elem(),
		ID:    REQUEST_CHAT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*EditGuildMember)(nil)).Elem(),
		ID:    EDIT_GUILD_MEMBER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*EditGuildNotice)(nil)).Elem(),
		ID:    EDIT_GUILD_NOTICE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GuildInvite)(nil)).Elem(),
		ID:    GUILD_INVITE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GuildNameReturn)(nil)).Elem(),
		ID:    GUILD_NAME_RETURN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RequestGuildInfo)(nil)).Elem(),
		ID:    REQUEST_GUILD_INFO,
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
		Type:  reflect.TypeOf((*GuildWarReturn)(nil)).Elem(),
		ID:    GUILD_WAR_RETURN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MarriageRequest)(nil)).Elem(),
		ID:    MARRIAGE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MarriageReply)(nil)).Elem(),
		ID:    MARRIAGE_REPLY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ChangeMarriage)(nil)).Elem(),
		ID:    CHANGE_MARRIAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DivorceRequest)(nil)).Elem(),
		ID:    DIVORCE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DivorceReply)(nil)).Elem(),
		ID:    DIVORCE_REPLY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*AddMentor)(nil)).Elem(),
		ID:    ADD_MENTOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MentorReply)(nil)).Elem(),
		ID:    MENTOR_REPLY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*AllowMentor)(nil)).Elem(),
		ID:    ALLOW_MENTOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CancelMentor)(nil)).Elem(),
		ID:    CANCEL_MENTOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*TradeRequest)(nil)).Elem(),
		ID:    TRADE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*TradeReply)(nil)).Elem(),
		ID:    TRADE_REPLY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*TradeGold)(nil)).Elem(),
		ID:    TRADE_GOLD,
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
		Type:  reflect.TypeOf((*EquipSlotItem)(nil)).Elem(),
		ID:    EQUIP_SLOT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*FishingCast)(nil)).Elem(),
		ID:    FISHING_CAST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*FishingChangeAutocast)(nil)).Elem(),
		ID:    FISHING_CHANGE_AUTOCAST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*AcceptQuest)(nil)).Elem(),
		ID:    ACCEPT_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*FinishQuest)(nil)).Elem(),
		ID:    FINISH_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*AbandonQuest)(nil)).Elem(),
		ID:    ABANDON_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ShareQuest)(nil)).Elem(),
		ID:    SHARE_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*AcceptReincarnation)(nil)).Elem(),
		ID:    ACCEPT_REINCARNATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CancelReincarnation)(nil)).Elem(),
		ID:    CANCEL_REINCARNATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CombineItem)(nil)).Elem(),
		ID:    COMBINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SetConcentration)(nil)).Elem(),
		ID:    SET_CONCENTRATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*AwakeningNeedMaterials)(nil)).Elem(),
		ID:    AWAKENING_NEED_MATERIALS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*AwakeningLockedItem)(nil)).Elem(),
		ID:    AWAKENING_LOCKED_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Awakening)(nil)).Elem(),
		ID:    AWAKENING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DisassembleItem)(nil)).Elem(),
		ID:    DISASSEMBLE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DowngradeAwakening)(nil)).Elem(),
		ID:    DOWNGRADE_AWAKENING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ResetAddedItem)(nil)).Elem(),
		ID:    RESET_ADDED_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SendMail)(nil)).Elem(),
		ID:    SEND_MAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ReadMail)(nil)).Elem(),
		ID:    READ_MAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CollectParcel)(nil)).Elem(),
		ID:    COLLECT_PARCEL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DeleteMail)(nil)).Elem(),
		ID:    DELETE_MAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*LockMail)(nil)).Elem(),
		ID:    LOCK_MAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MailLockedItem)(nil)).Elem(),
		ID:    MAIL_LOCKED_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*MailCost)(nil)).Elem(),
		ID:    MAIL_COST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*UpdateIntelligentCreature)(nil)).Elem(),
		ID:    UPDATE_INTELLIGENT_CREATURE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*IntelligentCreaturePickup)(nil)).Elem(),
		ID:    INTELLIGENT_CREATURE_PICKUP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*AddFriend)(nil)).Elem(),
		ID:    ADD_FRIEND,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RemoveFriend)(nil)).Elem(),
		ID:    REMOVE_FRIEND,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RefreshFriends)(nil)).Elem(),
		ID:    REFRESH_FRIENDS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*AddMemo)(nil)).Elem(),
		ID:    ADD_MEMO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GuildBuffUpdate)(nil)).Elem(),
		ID:    GUILD_BUFF_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*NPCConfirmInput)(nil)).Elem(),
		ID:    NPC_CONFIRM_INPUT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GameshopBuy)(nil)).Elem(),
		ID:    GAMESHOP_BUY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ReportIssue)(nil)).Elem(),
		ID:    REPORT_ISSUE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GetRanking)(nil)).Elem(),
		ID:    GET_RANKING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*Opendoor)(nil)).Elem(),
		ID:    OPENDOOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*GetRentedItems)(nil)).Elem(),
		ID:    GET_RENTED_ITEMS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ItemRentalRequest)(nil)).Elem(),
		ID:    ITEM_RENTAL_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ItemRentalFee)(nil)).Elem(),
		ID:    ITEM_RENTAL_FEE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ItemRentalPeriod)(nil)).Elem(),
		ID:    ITEM_RENTAL_PERIOD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*DepositRentalItem)(nil)).Elem(),
		ID:    DEPOSIT_RENTAL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*RetrieveRentalItem)(nil)).Elem(),
		ID:    RETRIEVE_RENTAL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CancelItemRental)(nil)).Elem(),
		ID:    CANCEL_ITEM_RENTAL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ItemRentalLockFee)(nil)).Elem(),
		ID:    ITEM_RENTAL_LOCK_FEE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ItemRentalLockItem)(nil)).Elem(),
		ID:    ITEM_RENTAL_LOCK_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*ConfirmItemRental)(nil)).Elem(),
		ID:    CONFIRM_ITEM_RENTAL,
	})
}
