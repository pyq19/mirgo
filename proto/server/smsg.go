package server

import (
	"github.com/yenkeia/mirgo/common"

	// 使用binary协议，因此匿名引用这个包，底层会自动注册

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

type UserInformation struct {
	ObjectID       uint32
	RealId         uint32
	Name           string
	GuildName      string
	GuildRank      string
	NameColour     uint32
	Class          common.MirClass
	Gender         common.MirGender
	Level          uint16
	Location       common.Point
	Direction      common.MirDirection
	Hair           uint8
	HP             uint16
	MP             uint16
	Experience     int64
	MaxExperience  int64
	LevelEffect    common.LevelEffects
	Inventory      []common.UserItem
	Equipment      []common.UserItem
	QuestInventory []common.UserItem
	Gold           uint32
	Credit         uint32
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

type SetObjectConcentration struct{}
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
type Opendoor struct{}
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
