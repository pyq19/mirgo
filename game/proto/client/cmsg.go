package client

import (
	_ "github.com/davyxu/cellnet/codec/binary"
	"github.com/yenkeia/mirgo/game/cm"
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
	DEL_MEMBER
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
	Gender cm.MirGender
	Class  cm.MirClass
}

type DeleteCharacter struct {
	CharacterIndex int32
}

type StartGame struct {
	CharacterIndex int16
}

type LogOut struct{}

type Turn struct {
	Direction cm.MirDirection
}

type Walk struct {
	Direction cm.MirDirection
}

type Run struct {
	Direction cm.MirDirection
}

type Chat struct {
	Message string
}

type MoveItem struct {
	Grid cm.MirGridType
	From int32
	To   int32
}

type StoreItem struct {
	From int32
	To   int32
}

type TakeBackItem struct {
	From int32
	To   int32
}

type MergeItem struct {
	GridFrom cm.MirGridType
	GridTo   cm.MirGridType
	IDFrom   uint64
	IDTo     uint64
}

type EquipItem struct {
	Grid     cm.MirGridType
	UniqueID uint64
	To       int32
}

type RemoveItem struct {
	Grid     cm.MirGridType
	UniqueID uint64
	To       int32
}

type RemoveSlotItem struct {
	Grid     cm.MirGridType
	GridTo   cm.MirGridType
	UniqueID uint64
	To       int32
}

type SplitItem struct {
	Grid     cm.MirGridType
	UniqueID uint64
	Count    uint32
}

type UseItem struct {
	UniqueID uint64
}

type DropItem struct {
	UniqueID uint64
	Count    uint32
}

type DepositRefineItem struct {
	From int32
	To   int32
}

type RetrieveRefineItem struct {
	From int32
	To   int32
}

type RefineCancel struct{}

type RefineItem struct {
	UniqueID uint64
}

type CheckRefine struct {
	UniqueID uint64
}

type ReplaceWedRing struct {
	UniqueID uint64
}

type DepositTradeItem struct {
	From int32
	To   int32
}

type RetrieveTradeItem struct {
	From int32
	To   int32
}

type DropGold struct {
	Amount uint32
}

type PickUp struct{}

type Inspect struct {
	ObjectID uint32
	Ranking  bool
}

type ChangeAMode struct {
	Mode cm.AttackMode
}

type ChangePMode struct {
	Mode cm.PetMode
}

type ChangeTrade struct {
	AllowTrade bool
}

type Attack struct {
	Direction cm.MirDirection
	Spell     cm.Spell
}

type RangeAttack struct {
	Direction      cm.MirDirection
	Location       cm.Point
	TargetID       uint32
	TargetLocation cm.Point
}

type Harvest struct {
	Direction cm.MirDirection
}

type CallNPC struct {
	ObjectID uint32
	Key      string
}

type TalkMonsterNPC struct {
	ObjectID uint32
}

type BuyItem struct {
	ItemIndex uint64
	Count     uint32
	Type      cm.PanelType
}

type SellItem struct {
	UniqueID uint64
	Count    uint32
}

// TODO
type CraftItem struct {
	UniqueID uint64
	Count    uint32
	Slots    []int
}

type RepairItem struct {
	UniqueID uint64
}

type BuyItemBack struct {
	UniqueID uint64
	Count    uint32
}

type SRepairItem struct {
	UniqueID uint64
}

type MagicKey struct {
	Spell cm.Spell
	Key   uint8
}

type Magic struct {
	Spell     cm.Spell
	Direction cm.MirDirection
	TargetID  uint32
	Location  cm.Point
}

type SwitchGroup struct {
	AllowGroup bool
}

type AddMember struct {
	Name string
}

type DelMember struct {
	Name string
}

type GroupInvite struct {
	AcceptInvite bool
}

type TownRevive struct{}

type SpellToggle struct {
	Spell  cm.Spell
	CanUse bool
}

type ConsignItem struct {
	UniqueID uint64
	Price    uint32
}

type MarketSearch struct {
	Match string
}

type MarketRefresh struct{}

type MarketPage struct {
	Page int32
}

type MarketBuy struct {
	AuctionID uint64
}

type MarketGetBack struct {
	AuctionID uint64
}

type RequestUserName struct {
	UserID uint32
}

type RequestChatItem struct {
	ChatItemID uint64
}

type EditGuildMember struct {
	ChangeType uint8
	RankIndex  uint8
	Name       string
	RankName   string
}

type EditGuildNotice struct {
	Notice []string
}

type GuildInvite struct {
	AcceptInvite bool
}

type GuildNameReturn struct {
	Name string
}

type RequestGuildInfo struct {
	Type uint8
}

type GuildStorageGoldChange struct {
	Type   uint8
	Amount uint32
}

type GuildStorageItemChange struct {
	Type uint8
	From int32
	To   int32
}

type GuildWarReturn struct {
	Name string
}

type MarriageRequest struct{}

type MarriageReply struct {
	AcceptInvite bool
}

type ChangeMarriage struct{}

type DivorceRequest struct{}

type DivorceReply struct {
	AcceptInvite bool
}

type AddMentor struct {
	Name string
}

type MentorReply struct {
	AcceptInvite bool
}

type AllowMentor struct{}

type CancelMentor struct{}

type TradeRequest struct{}

type TradeReply struct {
	AcceptInvite bool
}

type TradeGold struct {
	Amount uint32
}

type TradeConfirm struct {
	Locked bool
}

type TradeCancel struct{}

type EquipSlotItem struct {
	Grid     cm.MirGridType
	UniqueID uint64
	To       int32
	GridTo   cm.MirGridType
}

type FishingCast struct {
	CastOut bool
}

type FishingChangeAutocast struct {
	AutoCast bool
}

type AcceptQuest struct {
	NPCIndex   uint32
	QuestIndex int32
}

type FinishQuest struct {
	QuestIndex        int32
	SelectedItemIndex int32
}

type AbandonQuest struct {
	QuestIndex int32
}

type ShareQuest struct {
	QuestIndex int32
}

type AcceptReincarnation struct{}

type CancelReincarnation struct{}

type CombineItem struct {
	IDFrom uint64
	IDTo   uint64
}

type SetConcentration struct {
	ObjectID    uint32
	Enabled     bool
	Interrupted bool
}

// TODO
type AwakeningNeedMaterials struct{}

// TODO
type AwakeningLockedItem struct{}

// TODO
type Awakening struct{}

// TODO
type DisassembleItem struct{}

// TODO
type DowngradeAwakening struct{}

// TODO
type ResetAddedItem struct{}

// TODO
type SendMail struct{}

// TODO
type ReadMail struct{}

// TODO
type CollectParcel struct{}

// TODO
type DeleteMail struct{}

// TODO
type LockMail struct{}

// TODO
type MailLockedItem struct{}

// TODO
type MailCost struct{}

// TODO
type UpdateIntelligentCreature struct{}

// TODO
type IntelligentCreaturePickup struct{}

// TODO
type AddFriend struct{}

// TODO
type RemoveFriend struct{}

// TODO
type RefreshFriends struct{}

// TODO
type AddMemo struct{}

// TODO
type GuildBuffUpdate struct{}

// TODO
type NPCConfirmInput struct{}

// TODO
type GameshopBuy struct{}

// TODO
type ReportIssue struct{}

// TODO
type GetRanking struct{}

type Opendoor struct {
	DoorIndex byte
}

// TODO
type GetRentedItems struct{}

// TODO
type ItemRentalRequest struct{}

// TODO
type ItemRentalFee struct{}

// TODO
type ItemRentalPeriod struct{}

// TODO
type DepositRentalItem struct{}

// TODO
type RetrieveRentalItem struct{}

// TODO
type CancelItemRental struct{}

// TODO
type ItemRentalLockFee struct{}

// TODO
type ItemRentalLockItem struct{}

// TODO
type ConfirmItemRental struct{}
