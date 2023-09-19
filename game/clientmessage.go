package game

import (
	_ "github.com/davyxu/cellnet/codec/binary"
)

const (
	CM_CLIENT_VERSION = 1000 + iota
	CM_DISCONNECT
	CM_KEEP_ALIVE
	CM_NEW_ACCOUNT
	CM_CHANGE_PASSWORD
	CM_LOGIN
	CM_NEW_CHARACTER
	CM_DELETE_CHARACTER
	CM_START_GAME
	CM_LOG_OUT
	CM_TURN
	CM_WALK
	CM_RUN
	CM_CHAT
	CM_MOVE_ITEM
	CM_STORE_ITEM
	CM_TAKE_BACK_ITEM
	CM_MERGE_ITEM
	CM_EQUIP_ITEM
	CM_REMOVE_ITEM
	CM_REMOVE_SLOT_ITEM
	CM_SPLIT_ITEM
	CM_USE_ITEM
	CM_DROP_ITEM
	CM_DEPOSIT_REFINE_ITEM
	CM_RETRIEVE_REFINE_ITEM
	CM_REFINE_CANCEL
	CM_REFINE_ITEM
	CM_CHECK_REFINE
	CM_REPLACE_WED_RING
	CM_DEPOSIT_TRADE_ITEM
	CM_RETRIEVE_TRADE_ITEM
	CM_DROP_GOLD
	CM_PICK_UP
	CM_INSPECT
	CM_CHANGE_A_MODE
	CM_CHANGE_P_MODE
	CM_CHANGE_TRADE
	CM_ATTACK
	CM_RANGE_ATTACK
	CM_HARVEST
	CM_CALL_NPC
	CM_TALK_MONSTER_NPC
	CM_BUY_ITEM
	CM_SELL_ITEM
	CM_CRAFT_ITEM
	CM_REPAIR_ITEM
	CM_BUY_ITEM_BACK
	CM_S_REPAIR_ITEM
	CM_MAGIC_KEY
	CM_MAGIC
	CM_SWITCH_GROUP
	CM_ADD_MEMBER
	CM_DEL_MEMBER
	CM_GROUP_INVITE
	CM_TOWN_REVIVE
	CM_SPELL_TOGGLE
	CM_CONSIGN_ITEM
	CM_MARKET_SEARCH
	CM_MARKET_REFRESH
	CM_MARKET_PAGE
	CM_MARKET_BUY
	CM_MARKET_GET_BACK
	CM_REQUEST_USER_NAME
	CM_REQUEST_CHAT_ITEM
	CM_EDIT_GUILD_MEMBER
	CM_EDIT_GUILD_NOTICE
	CM_GUILD_INVITE
	CM_GUILD_NAME_RETURN
	CM_REQUEST_GUILD_INFO
	CM_GUILD_STORAGE_GOLD_CHANGE
	CM_GUILD_STORAGE_ITEM_CHANGE
	CM_GUILD_WAR_RETURN
	CM_MARRIAGE_REQUEST
	CM_MARRIAGE_REPLY
	CM_CHANGE_MARRIAGE
	CM_DIVORCE_REQUEST
	CM_DIVORCE_REPLY
	CM_ADD_MENTOR
	CM_MENTOR_REPLY
	CM_ALLOW_MENTOR
	CM_CANCEL_MENTOR
	CM_TRADE_REQUEST
	CM_TRADE_REPLY
	CM_TRADE_GOLD
	CM_TRADE_CONFIRM
	CM_TRADE_CANCEL
	CM_EQUIP_SLOT_ITEM
	CM_FISHING_CAST
	CM_FISHING_CHANGE_AUTOCAST
	CM_ACCEPT_QUEST
	CM_FINISH_QUEST
	CM_ABANDON_QUEST
	CM_SHARE_QUEST
	CM_ACCEPT_REINCARNATION
	CM_CANCEL_REINCARNATION
	CM_COMBINE_ITEM
	CM_SET_CONCENTRATION
	CM_AWAKENING_NEED_MATERIALS
	CM_AWAKENING_LOCKED_ITEM
	CM_AWAKENING
	CM_DISASSEMBLE_ITEM
	CM_DOWNGRADE_AWAKENING
	CM_RESET_ADDED_ITEM
	CM_SEND_MAIL
	CM_READ_MAIL
	CM_COLLECT_PARCEL
	CM_DELETE_MAIL
	CM_LOCK_MAIL
	CM_MAIL_LOCKED_ITEM
	CM_MAIL_COST
	CM_UPDATE_INTELLIGENT_CREATURE
	CM_INTELLIGENT_CREATURE_PICKUP
	CM_ADD_FRIEND
	CM_REMOVE_FRIEND
	CM_REFRESH_FRIENDS
	CM_ADD_MEMO
	CM_GUILD_BUFF_UPDATE
	CM_NPC_CONFIRM_INPUT
	CM_GAMESHOP_BUY
	CM_REPORT_ISSUE
	CM_GET_RANKING
	CM_OPENDOOR
	CM_GET_RENTED_ITEMS
	CM_ITEM_RENTAL_REQUEST
	CM_ITEM_RENTAL_FEE
	CM_ITEM_RENTAL_PERIOD
	CM_DEPOSIT_RENTAL_ITEM
	CM_RETRIEVE_RENTAL_ITEM
	CM_CANCEL_ITEM_RENTAL
	CM_ITEM_RENTAL_LOCK_FEE
	CM_ITEM_RENTAL_LOCK_ITEM
	CM_CONFIRM_ITEM_RENTAL
)

type CM_ClientVersion struct {
	VersionHash []uint8
}

type CM_Disconnect struct{}

type CM_KeepAlive struct {
	Time int64
}

type CM_NewAccount struct {
	AccountID      string
	Password       string
	DateTime       int64  // 无用字段 c# 中 DateTime 8 字节
	UserName       string // 无用字段
	SecretQuestion string // 无用字段
	SecretAnswer   string // 无用字段
	EMailAddress   string // 无用字段
}

type CM_ChangePassword struct {
	AccountID       string
	CurrentPassword string
	NewPassword     string
}

type CM_Login struct {
	AccountID string
	Password  string
}

type CM_NewCharacter struct {
	Name   string
	Gender MirGender
	Class  MirClass
}

type CM_DeleteCharacter struct {
	CharacterIndex int32
}

type CM_StartGame struct {
	CharacterIndex int16
}

type CM_LogOut struct{}

type CM_Turn struct {
	Direction MirDirection
}

type CM_Walk struct {
	Direction MirDirection
}

type CM_Run struct {
	Direction MirDirection
}

type CM_Chat struct {
	Message string
}

type CM_MoveItem struct {
	Grid MirGridType
	From int32
	To   int32
}

type CM_StoreItem struct {
	From int32
	To   int32
}

type CM_TakeBackItem struct {
	From int32
	To   int32
}

type CM_MergeItem struct {
	GridFrom MirGridType
	GridTo   MirGridType
	IDFrom   uint64
	IDTo     uint64
}

type CM_EquipItem struct {
	Grid     MirGridType
	UniqueID uint64
	To       int32
}

type CM_RemoveItem struct {
	Grid     MirGridType
	UniqueID uint64
	To       int32
}

type CM_RemoveSlotItem struct {
	Grid     MirGridType
	GridTo   MirGridType
	UniqueID uint64
	To       int32
}

type CM_SplitItem struct {
	Grid     MirGridType
	UniqueID uint64
	Count    uint32
}

type CM_UseItem struct {
	UniqueID uint64
}

type CM_DropItem struct {
	UniqueID uint64
	Count    uint32
}

type CM_DepositRefineItem struct {
	From int32
	To   int32
}

type CM_RetrieveRefineItem struct {
	From int32
	To   int32
}

type CM_RefineCancel struct{}

type CM_RefineItem struct {
	UniqueID uint64
}

type CM_CheckRefine struct {
	UniqueID uint64
}

type CM_ReplaceWedRing struct {
	UniqueID uint64
}

type CM_DepositTradeItem struct {
	From int32
	To   int32
}

type CM_RetrieveTradeItem struct {
	From int32
	To   int32
}

type CM_DropGold struct {
	Amount uint32
}

type CM_PickUp struct{}

type CM_Inspect struct {
	ObjectID uint32
	Ranking  bool
}

type CM_ChangeAMode struct {
	Mode AttackMode
}

type CM_ChangePMode struct {
	Mode PetMode
}

type CM_ChangeTrade struct {
	AllowTrade bool
}

type CM_Attack struct {
	Direction MirDirection
	Spell     Spell
}

type CM_RangeAttack struct {
	Direction      MirDirection
	Location       Point
	TargetID       uint32
	TargetLocation Point
}

type CM_Harvest struct {
	Direction MirDirection
}

type CM_CallNPC struct {
	ObjectID uint32
	Key      string
}

type CM_TalkMonsterNPC struct {
	ObjectID uint32
}

type CM_BuyItem struct {
	ItemIndex uint64
	Count     uint32
	Type      PanelType
}

type CM_SellItem struct {
	UniqueID uint64
	Count    uint32
}

// TODO
type CM_CraftItem struct {
	UniqueID uint64
	Count    uint32
	Slots    []int
}

type CM_RepairItem struct {
	UniqueID uint64
}

type CM_BuyItemBack struct {
	UniqueID uint64
	Count    uint32
}

type CM_SRepairItem struct {
	UniqueID uint64
}

type CM_MagicKey struct {
	Spell Spell
	Key   uint8
}

type CM_Magic struct {
	Spell     Spell
	Direction MirDirection
	TargetID  uint32
	Location  Point
}

type CM_SwitchGroup struct {
	AllowGroup bool
}

type CM_AddMember struct {
	Name string
}

type CM_DelMember struct {
	Name string
}

type CM_GroupInvite struct {
	AcceptInvite bool
}

type CM_TownRevive struct{}

type CM_SpellToggle struct {
	Spell  Spell
	CanUse bool
}

type CM_ConsignItem struct {
	UniqueID uint64
	Price    uint32
}

type CM_MarketSearch struct {
	Match string
}

type CM_MarketRefresh struct{}

type CM_MarketPage struct {
	Page int32
}

type CM_MarketBuy struct {
	AuctionID uint64
}

type CM_MarketGetBack struct {
	AuctionID uint64
}

type CM_RequestUserName struct {
	UserID uint32
}

type CM_RequestChatItem struct {
	ChatItemID uint64
}

type CM_EditGuildMember struct {
	Changetype uint8
	RankIndex  uint8
	Name       string
	RankName   string
}

type CM_EditGuildNotice struct {
	Notice []string
}

type CM_GuildInvite struct {
	AcceptInvite bool
}

type CM_GuildNameReturn struct {
	Name string
}

type CM_RequestGuildInfo struct {
	Type uint8
}

type CM_GuildStorageGoldChange struct {
	Type   uint8
	Amount uint32
}

type CM_GuildStorageItemChange struct {
	Type uint8
	From int32
	To   int32
}

type CM_GuildWarReturn struct {
	Name string
}

type CM_MarriageRequest struct{}

type CM_MarriageReply struct {
	AcceptInvite bool
}

type CM_ChangeMarriage struct{}

type CM_DivorceRequest struct{}

type CM_DivorceReply struct {
	AcceptInvite bool
}

type CM_AddMentor struct {
	Name string
}

type CM_MentorReply struct {
	AcceptInvite bool
}

type CM_AllowMentor struct{}

type CM_CancelMentor struct{}

type CM_TradeRequest struct{}

type CM_TradeReply struct {
	AcceptInvite bool
}

type CM_TradeGold struct {
	Amount uint32
}

type CM_TradeConfirm struct {
	Locked bool
}

type CM_TradeCancel struct{}

type CM_EquipSlotItem struct {
	Grid     MirGridType
	UniqueID uint64
	To       int32
	GridTo   MirGridType
}

type CM_FishingCast struct {
	CastOut bool
}

type CM_FishingChangeAutocast struct {
	AutoCast bool
}

type CM_AcceptQuest struct {
	NPCIndex   uint32
	QuestIndex int32
}

type CM_FinishQuest struct {
	QuestIndex        int32
	SelectedItemIndex int32
}

type CM_AbandonQuest struct {
	QuestIndex int32
}

type CM_ShareQuest struct {
	QuestIndex int32
}

type CM_AcceptReincarnation struct{}

type CM_CancelReincarnation struct{}

type CM_CombineItem struct {
	IDFrom uint64
	IDTo   uint64
}

type CM_SetConcentration struct {
	ObjectID    uint32
	Enabled     bool
	Interrupted bool
}

// TODO
type CM_AwakeningNeedMaterials struct{}

// TODO
type CM_AwakeningLockedItem struct{}

// TODO
type CM_Awakening struct{}

// TODO
type CM_DisassembleItem struct{}

// TODO
type CM_DowngradeAwakening struct{}

// TODO
type CM_ResetAddedItem struct{}

// TODO
type CM_SendMail struct{}

// TODO
type CM_ReadMail struct{}

// TODO
type CM_CollectParcel struct{}

// TODO
type CM_DeleteMail struct{}

// TODO
type CM_LockMail struct{}

// TODO
type CM_MailLockedItem struct{}

// TODO
type CM_MailCost struct{}

// TODO
type CM_UpdateIntelligentCreature struct{}

// TODO
type CM_IntelligentCreaturePickup struct{}

// TODO
type CM_AddFriend struct{}

// TODO
type CM_RemoveFriend struct{}

// TODO
type CM_RefreshFriends struct{}

// TODO
type CM_AddMemo struct{}

// TODO
type CM_GuildBuffUpdate struct{}

// TODO
type CM_NPCConfirmInput struct{}

// TODO
type CM_GameshopBuy struct{}

// TODO
type CM_ReportIssue struct{}

// TODO
type CM_GetRanking struct{}

type CM_Opendoor struct {
	DoorIndex byte
}

// TODO
type CM_GetRentedItems struct{}

// TODO
type CM_ItemRentalRequest struct{}

// TODO
type CM_ItemRentalFee struct{}

// TODO
type CM_ItemRentalPeriod struct{}

// TODO
type CM_DepositRentalItem struct{}

// TODO
type CM_RetrieveRentalItem struct{}

// TODO
type CM_CancelItemRental struct{}

// TODO
type CM_ItemRentalLockFee struct{}

// TODO
type CM_ItemRentalLockItem struct{}

// TODO
type CM_ConfirmItemRental struct{}
