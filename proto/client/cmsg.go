package client

import (
	"github.com/yenkeia/mirgo/common"

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

type Chat struct {
	Message string
}

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
