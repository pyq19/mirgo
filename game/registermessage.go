package game

import (
	"reflect"

	"github.com/davyxu/cellnet"
)

func init() {
	initClientMessage()
	initServerMessage()
}

// 引用消息时，自动注册消息，这个文件可以由代码生成自动生成
func initClientMessage() {
	mirCodec := new(MirCodec)
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ClientVersion)(nil)).Elem(),
		ID:    CM_CLIENT_VERSION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_Disconnect)(nil)).Elem(),
		ID:    CM_DISCONNECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_KeepAlive)(nil)).Elem(),
		ID:    CM_KEEP_ALIVE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_NewAccount)(nil)).Elem(),
		ID:    CM_NEW_ACCOUNT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ChangePassword)(nil)).Elem(),
		ID:    CM_CHANGE_PASSWORD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_Login)(nil)).Elem(),
		ID:    CM_LOGIN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_NewCharacter)(nil)).Elem(),
		ID:    CM_NEW_CHARACTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_DeleteCharacter)(nil)).Elem(),
		ID:    CM_DELETE_CHARACTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_StartGame)(nil)).Elem(),
		ID:    CM_START_GAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_LogOut)(nil)).Elem(),
		ID:    CM_LOG_OUT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_Turn)(nil)).Elem(),
		ID:    CM_TURN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_Walk)(nil)).Elem(),
		ID:    CM_WALK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_Run)(nil)).Elem(),
		ID:    CM_RUN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_Chat)(nil)).Elem(),
		ID:    CM_CHAT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_MoveItem)(nil)).Elem(),
		ID:    CM_MOVE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_StoreItem)(nil)).Elem(),
		ID:    CM_STORE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_TakeBackItem)(nil)).Elem(),
		ID:    CM_TAKE_BACK_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_MergeItem)(nil)).Elem(),
		ID:    CM_MERGE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_EquipItem)(nil)).Elem(),
		ID:    CM_EQUIP_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_RemoveItem)(nil)).Elem(),
		ID:    CM_REMOVE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_RemoveSlotItem)(nil)).Elem(),
		ID:    CM_REMOVE_SLOT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_SplitItem)(nil)).Elem(),
		ID:    CM_SPLIT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_UseItem)(nil)).Elem(),
		ID:    CM_USE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_DropItem)(nil)).Elem(),
		ID:    CM_DROP_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_DepositRefineItem)(nil)).Elem(),
		ID:    CM_DEPOSIT_REFINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_RetrieveRefineItem)(nil)).Elem(),
		ID:    CM_RETRIEVE_REFINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_RefineCancel)(nil)).Elem(),
		ID:    CM_REFINE_CANCEL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_RefineItem)(nil)).Elem(),
		ID:    CM_REFINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_CheckRefine)(nil)).Elem(),
		ID:    CM_CHECK_REFINE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ReplaceWedRing)(nil)).Elem(),
		ID:    CM_REPLACE_WED_RING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_DepositTradeItem)(nil)).Elem(),
		ID:    CM_DEPOSIT_TRADE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_RetrieveTradeItem)(nil)).Elem(),
		ID:    CM_RETRIEVE_TRADE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_DropGold)(nil)).Elem(),
		ID:    CM_DROP_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_PickUp)(nil)).Elem(),
		ID:    CM_PICK_UP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_Inspect)(nil)).Elem(),
		ID:    CM_INSPECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ChangeAMode)(nil)).Elem(),
		ID:    CM_CHANGE_A_MODE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ChangePMode)(nil)).Elem(),
		ID:    CM_CHANGE_P_MODE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ChangeTrade)(nil)).Elem(),
		ID:    CM_CHANGE_TRADE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_Attack)(nil)).Elem(),
		ID:    CM_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_RangeAttack)(nil)).Elem(),
		ID:    CM_RANGE_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_Harvest)(nil)).Elem(),
		ID:    CM_HARVEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_CallNPC)(nil)).Elem(),
		ID:    CM_CALL_NPC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_TalkMonsterNPC)(nil)).Elem(),
		ID:    CM_TALK_MONSTER_NPC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_BuyItem)(nil)).Elem(),
		ID:    CM_BUY_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_SellItem)(nil)).Elem(),
		ID:    CM_SELL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_CraftItem)(nil)).Elem(),
		ID:    CM_CRAFT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_RepairItem)(nil)).Elem(),
		ID:    CM_REPAIR_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_BuyItemBack)(nil)).Elem(),
		ID:    CM_BUY_ITEM_BACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_SRepairItem)(nil)).Elem(),
		ID:    CM_S_REPAIR_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_MagicKey)(nil)).Elem(),
		ID:    CM_MAGIC_KEY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_Magic)(nil)).Elem(),
		ID:    CM_MAGIC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_SwitchGroup)(nil)).Elem(),
		ID:    CM_SWITCH_GROUP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_AddMember)(nil)).Elem(),
		ID:    CM_ADD_MEMBER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_DelMember)(nil)).Elem(),
		ID:    CM_DEL_MEMBER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_GroupInvite)(nil)).Elem(),
		ID:    CM_GROUP_INVITE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_TownRevive)(nil)).Elem(),
		ID:    CM_TOWN_REVIVE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_SpellToggle)(nil)).Elem(),
		ID:    CM_SPELL_TOGGLE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ConsignItem)(nil)).Elem(),
		ID:    CM_CONSIGN_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_MarketSearch)(nil)).Elem(),
		ID:    CM_MARKET_SEARCH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_MarketRefresh)(nil)).Elem(),
		ID:    CM_MARKET_REFRESH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_MarketPage)(nil)).Elem(),
		ID:    CM_MARKET_PAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_MarketBuy)(nil)).Elem(),
		ID:    CM_MARKET_BUY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_MarketGetBack)(nil)).Elem(),
		ID:    CM_MARKET_GET_BACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_RequestUserName)(nil)).Elem(),
		ID:    CM_REQUEST_USER_NAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_RequestChatItem)(nil)).Elem(),
		ID:    CM_REQUEST_CHAT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_EditGuildMember)(nil)).Elem(),
		ID:    CM_EDIT_GUILD_MEMBER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_EditGuildNotice)(nil)).Elem(),
		ID:    CM_EDIT_GUILD_NOTICE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_GuildInvite)(nil)).Elem(),
		ID:    CM_GUILD_INVITE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_GuildNameReturn)(nil)).Elem(),
		ID:    CM_GUILD_NAME_RETURN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_RequestGuildInfo)(nil)).Elem(),
		ID:    CM_REQUEST_GUILD_INFO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_GuildStorageGoldChange)(nil)).Elem(),
		ID:    CM_GUILD_STORAGE_GOLD_CHANGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_GuildStorageItemChange)(nil)).Elem(),
		ID:    CM_GUILD_STORAGE_ITEM_CHANGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_GuildWarReturn)(nil)).Elem(),
		ID:    CM_GUILD_WAR_RETURN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_MarriageRequest)(nil)).Elem(),
		ID:    CM_MARRIAGE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_MarriageReply)(nil)).Elem(),
		ID:    CM_MARRIAGE_REPLY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ChangeMarriage)(nil)).Elem(),
		ID:    CM_CHANGE_MARRIAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_DivorceRequest)(nil)).Elem(),
		ID:    CM_DIVORCE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_DivorceReply)(nil)).Elem(),
		ID:    CM_DIVORCE_REPLY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_AddMentor)(nil)).Elem(),
		ID:    CM_ADD_MENTOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_MentorReply)(nil)).Elem(),
		ID:    CM_MENTOR_REPLY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_AllowMentor)(nil)).Elem(),
		ID:    CM_ALLOW_MENTOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_CancelMentor)(nil)).Elem(),
		ID:    CM_CANCEL_MENTOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_TradeRequest)(nil)).Elem(),
		ID:    CM_TRADE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_TradeReply)(nil)).Elem(),
		ID:    CM_TRADE_REPLY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_TradeGold)(nil)).Elem(),
		ID:    CM_TRADE_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_TradeConfirm)(nil)).Elem(),
		ID:    CM_TRADE_CONFIRM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_TradeCancel)(nil)).Elem(),
		ID:    CM_TRADE_CANCEL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_EquipSlotItem)(nil)).Elem(),
		ID:    CM_EQUIP_SLOT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_FishingCast)(nil)).Elem(),
		ID:    CM_FISHING_CAST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_FishingChangeAutocast)(nil)).Elem(),
		ID:    CM_FISHING_CHANGE_AUTOCAST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_AcceptQuest)(nil)).Elem(),
		ID:    CM_ACCEPT_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_FinishQuest)(nil)).Elem(),
		ID:    CM_FINISH_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_AbandonQuest)(nil)).Elem(),
		ID:    CM_ABANDON_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ShareQuest)(nil)).Elem(),
		ID:    CM_SHARE_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_AcceptReincarnation)(nil)).Elem(),
		ID:    CM_ACCEPT_REINCARNATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_CancelReincarnation)(nil)).Elem(),
		ID:    CM_CANCEL_REINCARNATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_CombineItem)(nil)).Elem(),
		ID:    CM_COMBINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_SetConcentration)(nil)).Elem(),
		ID:    CM_SET_CONCENTRATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_AwakeningNeedMaterials)(nil)).Elem(),
		ID:    CM_AWAKENING_NEED_MATERIALS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_AwakeningLockedItem)(nil)).Elem(),
		ID:    CM_AWAKENING_LOCKED_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_Awakening)(nil)).Elem(),
		ID:    CM_AWAKENING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_DisassembleItem)(nil)).Elem(),
		ID:    CM_DISASSEMBLE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_DowngradeAwakening)(nil)).Elem(),
		ID:    CM_DOWNGRADE_AWAKENING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ResetAddedItem)(nil)).Elem(),
		ID:    CM_RESET_ADDED_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_SendMail)(nil)).Elem(),
		ID:    CM_SEND_MAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ReadMail)(nil)).Elem(),
		ID:    CM_READ_MAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_CollectParcel)(nil)).Elem(),
		ID:    CM_COLLECT_PARCEL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_DeleteMail)(nil)).Elem(),
		ID:    CM_DELETE_MAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_LockMail)(nil)).Elem(),
		ID:    CM_LOCK_MAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_MailLockedItem)(nil)).Elem(),
		ID:    CM_MAIL_LOCKED_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_MailCost)(nil)).Elem(),
		ID:    CM_MAIL_COST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_UpdateIntelligentCreature)(nil)).Elem(),
		ID:    CM_UPDATE_INTELLIGENT_CREATURE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_IntelligentCreaturePickup)(nil)).Elem(),
		ID:    CM_INTELLIGENT_CREATURE_PICKUP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_AddFriend)(nil)).Elem(),
		ID:    CM_ADD_FRIEND,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_RemoveFriend)(nil)).Elem(),
		ID:    CM_REMOVE_FRIEND,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_RefreshFriends)(nil)).Elem(),
		ID:    CM_REFRESH_FRIENDS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_AddMemo)(nil)).Elem(),
		ID:    CM_ADD_MEMO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_GuildBuffUpdate)(nil)).Elem(),
		ID:    CM_GUILD_BUFF_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_NPCConfirmInput)(nil)).Elem(),
		ID:    CM_NPC_CONFIRM_INPUT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_GameshopBuy)(nil)).Elem(),
		ID:    CM_GAMESHOP_BUY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ReportIssue)(nil)).Elem(),
		ID:    CM_REPORT_ISSUE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_GetRanking)(nil)).Elem(),
		ID:    CM_GET_RANKING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_Opendoor)(nil)).Elem(),
		ID:    CM_OPENDOOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_GetRentedItems)(nil)).Elem(),
		ID:    CM_GET_RENTED_ITEMS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ItemRentalRequest)(nil)).Elem(),
		ID:    CM_ITEM_RENTAL_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ItemRentalFee)(nil)).Elem(),
		ID:    CM_ITEM_RENTAL_FEE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ItemRentalPeriod)(nil)).Elem(),
		ID:    CM_ITEM_RENTAL_PERIOD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_DepositRentalItem)(nil)).Elem(),
		ID:    CM_DEPOSIT_RENTAL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_RetrieveRentalItem)(nil)).Elem(),
		ID:    CM_RETRIEVE_RENTAL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_CancelItemRental)(nil)).Elem(),
		ID:    CM_CANCEL_ITEM_RENTAL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ItemRentalLockFee)(nil)).Elem(),
		ID:    CM_ITEM_RENTAL_LOCK_FEE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ItemRentalLockItem)(nil)).Elem(),
		ID:    CM_ITEM_RENTAL_LOCK_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*CM_ConfirmItemRental)(nil)).Elem(),
		ID:    CM_CONFIRM_ITEM_RENTAL,
	})
}

func initServerMessage() {

	mirCodec := new(MirCodec)

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_Connected)(nil)).Elem(),
		ID:    SM_CONNECTED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ClientVersion)(nil)).Elem(),
		ID:    SM_CLIENT_VERSION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_Disconnect)(nil)).Elem(),
		ID:    SM_DISCONNECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_KeepAlive)(nil)).Elem(),
		ID:    SM_KEEP_ALIVE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NewAccount)(nil)).Elem(),
		ID:    SM_NEW_ACCOUNT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ChangePassword)(nil)).Elem(),
		ID:    SM_CHANGE_PASSWORD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ChangePasswordBanned)(nil)).Elem(),
		ID:    SM_CHANGE_PASSWORD_BANNED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_Login)(nil)).Elem(),
		ID:    SM_LOGIN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_LoginBanned)(nil)).Elem(),
		ID:    SM_LOGIN_BANNED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_LoginSuccess)(nil)).Elem(),
		ID:    SM_LOGIN_SUCCESS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NewCharacter)(nil)).Elem(),
		ID:    SM_NEW_CHARACTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NewCharacterSuccess)(nil)).Elem(),
		ID:    SM_NEW_CHARACTER_SUCCESS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_DeleteCharacter)(nil)).Elem(),
		ID:    SM_DELETE_CHARACTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_DeleteCharacterSuccess)(nil)).Elem(),
		ID:    SM_DELETE_CHARACTER_SUCCESS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_StartGame)(nil)).Elem(),
		ID:    SM_START_GAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_StartGameBanned)(nil)).Elem(),
		ID:    SM_START_GAME_BANNED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_StartGameDelay)(nil)).Elem(),
		ID:    SM_START_GAME_DELAY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MapInformation)(nil)).Elem(),
		ID:    SM_MAP_INFORMATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_UserInformation)(nil)).Elem(),
		ID:    SM_USER_INFORMATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_UserLocation)(nil)).Elem(),
		ID:    SM_USER_LOCATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectPlayer)(nil)).Elem(),
		ID:    SM_OBJECT_PLAYER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectRemove)(nil)).Elem(),
		ID:    SM_OBJECT_REMOVE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectTurn)(nil)).Elem(),
		ID:    SM_OBJECT_TURN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectWalk)(nil)).Elem(),
		ID:    SM_OBJECT_WALK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectRun)(nil)).Elem(),
		ID:    SM_OBJECT_RUN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_Chat)(nil)).Elem(),
		ID:    SM_CHAT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectChat)(nil)).Elem(),
		ID:    SM_OBJECT_CHAT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NewItemInfo)(nil)).Elem(),
		ID:    SM_NEW_ITEM_INFO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MoveItem)(nil)).Elem(),
		ID:    SM_MOVE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_EquipItem)(nil)).Elem(),
		ID:    SM_EQUIP_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MergeItem)(nil)).Elem(),
		ID:    SM_MERGE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_RemoveItem)(nil)).Elem(),
		ID:    SM_REMOVE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_RemoveSlotItem)(nil)).Elem(),
		ID:    SM_REMOVE_SLOT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_TakeBackItem)(nil)).Elem(),
		ID:    SM_TAKE_BACK_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_StoreItem)(nil)).Elem(),
		ID:    SM_STORE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_SplitItem)(nil)).Elem(),
		ID:    SM_SPLIT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_SplitItem1)(nil)).Elem(),
		ID:    SM_SPLIT_ITEM1,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_DepositRefineItem)(nil)).Elem(),
		ID:    SM_DEPOSIT_REFINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_RetrieveRefineItem)(nil)).Elem(),
		ID:    SM_RETRIEVE_REFINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_RefineCancel)(nil)).Elem(),
		ID:    SM_REFINE_CANCEL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_RefineItem)(nil)).Elem(),
		ID:    SM_REFINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_DepositTradeItem)(nil)).Elem(),
		ID:    SM_DEPOSIT_TRADE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_RetrieveTradeItem)(nil)).Elem(),
		ID:    SM_RETRIEVE_TRADE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_UseItem)(nil)).Elem(),
		ID:    SM_USE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_DropItem)(nil)).Elem(),
		ID:    SM_DROP_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_PlayerUpdate)(nil)).Elem(),
		ID:    SM_PLAYER_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_PlayerInspect)(nil)).Elem(),
		ID:    SM_PLAYER_INSPECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_LogOutSuccess)(nil)).Elem(),
		ID:    SM_LOG_OUT_SUCCESS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_LogOutFailed)(nil)).Elem(),
		ID:    SM_LOG_OUT_FAILED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_TimeOfDay)(nil)).Elem(),
		ID:    SM_TIME_OF_DAY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ChangeAMode)(nil)).Elem(),
		ID:    SM_CHANGE_A_MODE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ChangePMode)(nil)).Elem(),
		ID:    SM_CHANGE_P_MODE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectItem)(nil)).Elem(),
		ID:    SM_OBJECT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectGold)(nil)).Elem(),
		ID:    SM_OBJECT_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GainedItem)(nil)).Elem(),
		ID:    SM_GAINED_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GainedGold)(nil)).Elem(),
		ID:    SM_GAINED_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_LoseGold)(nil)).Elem(),
		ID:    SM_LOSE_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GainedCredit)(nil)).Elem(),
		ID:    SM_GAINED_CREDIT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_LoseCredit)(nil)).Elem(),
		ID:    SM_LOSE_CREDIT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectMonster)(nil)).Elem(),
		ID:    SM_OBJECT_MONSTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectAttack)(nil)).Elem(),
		ID:    SM_OBJECT_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_Struck)(nil)).Elem(),
		ID:    SM_STRUCK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectStruck)(nil)).Elem(),
		ID:    SM_OBJECT_STRUCK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_DamageIndicator)(nil)).Elem(),
		ID:    SM_DAMAGE_INDICATOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_DuraChanged)(nil)).Elem(),
		ID:    SM_DURA_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_HealthChanged)(nil)).Elem(),
		ID:    SM_HEALTH_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_DeleteItem)(nil)).Elem(),
		ID:    SM_DELETE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_Death)(nil)).Elem(),
		ID:    SM_DEATH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectDied)(nil)).Elem(),
		ID:    SM_OBJECT_DIED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ColourChanged)(nil)).Elem(),
		ID:    SM_COLOUR_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectColourChanged)(nil)).Elem(),
		ID:    SM_OBJECT_COLOUR_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectGuildNameChanged)(nil)).Elem(),
		ID:    SM_OBJECT_GUILD_NAME_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GainExperience)(nil)).Elem(),
		ID:    SM_GAIN_EXPERIENCE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_LevelChanged)(nil)).Elem(),
		ID:    SM_LEVEL_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectLeveled)(nil)).Elem(),
		ID:    SM_OBJECT_LEVELED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectHarvest)(nil)).Elem(),
		ID:    SM_OBJECT_HARVEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectHarvested)(nil)).Elem(),
		ID:    SM_OBJECT_HARVESTED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectNPC)(nil)).Elem(),
		ID:    SM_OBJECT_NPC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCResponse)(nil)).Elem(),
		ID:    SM_NPC_RESPONSE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectHide)(nil)).Elem(),
		ID:    SM_OBJECT_HIDE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectShow)(nil)).Elem(),
		ID:    SM_OBJECT_SHOW,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_Poisoned)(nil)).Elem(),
		ID:    SM_POISONED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectPoisoned)(nil)).Elem(),
		ID:    SM_OBJECT_POISONED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MapChanged)(nil)).Elem(),
		ID:    SM_MAP_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectTeleportOut)(nil)).Elem(),
		ID:    SM_OBJECT_TELEPORT_OUT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectTeleportIn)(nil)).Elem(),
		ID:    SM_OBJECT_TELEPORT_IN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_TeleportIn)(nil)).Elem(),
		ID:    SM_TELEPORT_IN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCGoods)(nil)).Elem(),
		ID:    SM_NPC_GOODS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCSell)(nil)).Elem(),
		ID:    SM_NPC_SELL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCRepair)(nil)).Elem(),
		ID:    SM_NPC_REPAIR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCSRepair)(nil)).Elem(),
		ID:    SM_NPC_S_REPAIR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCRefine)(nil)).Elem(),
		ID:    SM_NPC_REFINE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCCheckRefine)(nil)).Elem(),
		ID:    SM_NPC_CHECK_REFINE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCCollectRefine)(nil)).Elem(),
		ID:    SM_NPC_COLLECT_REFINE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCReplaceWedRing)(nil)).Elem(),
		ID:    SM_NPC_REPLACE_WED_RING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCStorage)(nil)).Elem(),
		ID:    SM_NPC_STORAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_SellItem)(nil)).Elem(),
		ID:    SM_SELL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_CraftItem)(nil)).Elem(),
		ID:    SM_CRAFT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_RepairItem)(nil)).Elem(),
		ID:    SM_REPAIR_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ItemRepaired)(nil)).Elem(),
		ID:    SM_ITEM_REPAIRED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NewMagic)(nil)).Elem(),
		ID:    SM_NEW_MAGIC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_RemoveMagic)(nil)).Elem(),
		ID:    SM_REMOVE_MAGIC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MagicLeveled)(nil)).Elem(),
		ID:    SM_MAGIC_LEVELED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_Magic)(nil)).Elem(),
		ID:    SM_MAGIC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MagicDelay)(nil)).Elem(),
		ID:    SM_MAGIC_DELAY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MagicCast)(nil)).Elem(),
		ID:    SM_MAGIC_CAST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectMagic)(nil)).Elem(),
		ID:    SM_OBJECT_MAGIC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectEffect)(nil)).Elem(),
		ID:    SM_OBJECT_EFFECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_RangeAttack)(nil)).Elem(),
		ID:    SM_RANGE_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_Pushed)(nil)).Elem(),
		ID:    SM_PUSHED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectPushed)(nil)).Elem(),
		ID:    SM_OBJECT_PUSHED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectName)(nil)).Elem(),
		ID:    SM_OBJECT_NAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_UserStorage)(nil)).Elem(),
		ID:    SM_USER_STORAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_SwitchGroup)(nil)).Elem(),
		ID:    SM_SWITCH_GROUP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_DeleteGroup)(nil)).Elem(),
		ID:    SM_DELETE_GROUP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_DeleteMember)(nil)).Elem(),
		ID:    SM_DELETE_MEMBER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GroupInvite)(nil)).Elem(),
		ID:    SM_GROUP_INVITE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_AddMember)(nil)).Elem(),
		ID:    SM_ADD_MEMBER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_Revived)(nil)).Elem(),
		ID:    SM_REVIVED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectRevived)(nil)).Elem(),
		ID:    SM_OBJECT_REVIVED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_SpellToggle)(nil)).Elem(),
		ID:    SM_SPELL_TOGGLE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectHealth)(nil)).Elem(),
		ID:    SM_OBJECT_HEALTH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MapEffect)(nil)).Elem(),
		ID:    SM_MAP_EFFECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectRangeAttack)(nil)).Elem(),
		ID:    SM_OBJECT_RANGE_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_AddBuff)(nil)).Elem(),
		ID:    SM_ADD_BUFF,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_RemoveBuff)(nil)).Elem(),
		ID:    SM_REMOVE_BUFF,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectHidden)(nil)).Elem(),
		ID:    SM_OBJECT_HIDDEN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_RefreshItem)(nil)).Elem(),
		ID:    SM_REFRESH_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectSpell)(nil)).Elem(),
		ID:    SM_OBJECT_SPELL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_UserDash)(nil)).Elem(),
		ID:    SM_USER_DASH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectDash)(nil)).Elem(),
		ID:    SM_OBJECT_DASH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_UserDashFail)(nil)).Elem(),
		ID:    SM_USER_DASH_FAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectDashFail)(nil)).Elem(),
		ID:    SM_OBJECT_DASH_FAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCConsign)(nil)).Elem(),
		ID:    SM_NPC_CONSIGN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCMarket)(nil)).Elem(),
		ID:    SM_NPC_MARKET,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCMarketPage)(nil)).Elem(),
		ID:    SM_NPC_MARKET_PAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ConsignItem)(nil)).Elem(),
		ID:    SM_CONSIGN_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MarketFail)(nil)).Elem(),
		ID:    SM_MARKET_FAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MarketSuccess)(nil)).Elem(),
		ID:    SM_MARKET_SUCCESS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectSitDown)(nil)).Elem(),
		ID:    SM_OBJECT_SIT_DOWN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_InTrapRock)(nil)).Elem(),
		ID:    SM_IN_TRAP_ROCK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_BaseStatsInfo)(nil)).Elem(),
		ID:    SM_BASE_STATS_INFO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_UserName)(nil)).Elem(),
		ID:    SM_USER_NAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ChatItemStats)(nil)).Elem(),
		ID:    SM_CHAT_ITEM_STATS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GuildNoticeChange)(nil)).Elem(),
		ID:    SM_GUILD_NOTICE_CHANGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GuildMemberChange)(nil)).Elem(),
		ID:    SM_GUILD_MEMBER_CHANGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GuildStatus)(nil)).Elem(),
		ID:    SM_GUILD_STATUS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GuildInvite)(nil)).Elem(),
		ID:    SM_GUILD_INVITE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GuildExpGain)(nil)).Elem(),
		ID:    SM_GUILD_EXP_GAIN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GuildNameRequest)(nil)).Elem(),
		ID:    SM_GUILD_NAME_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GuildStorageGoldChange)(nil)).Elem(),
		ID:    SM_GUILD_STORAGE_GOLD_CHANGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GuildStorageItemChange)(nil)).Elem(),
		ID:    SM_GUILD_STORAGE_ITEM_CHANGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GuildStorageList)(nil)).Elem(),
		ID:    SM_GUILD_STORAGE_LIST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GuildRequestWar)(nil)).Elem(),
		ID:    SM_GUILD_REQUEST_WAR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_DefaultNPC)(nil)).Elem(),
		ID:    SM_DEFAULT_NPC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCUpdate)(nil)).Elem(),
		ID:    SM_NPC_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCImageUpdate)(nil)).Elem(),
		ID:    SM_NPC_IMAGE_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MarriageRequest)(nil)).Elem(),
		ID:    SM_MARRIAGE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_DivorceRequest)(nil)).Elem(),
		ID:    SM_DIVORCE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MentorRequest)(nil)).Elem(),
		ID:    SM_MENTOR_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_TradeRequest)(nil)).Elem(),
		ID:    SM_TRADE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_TradeAccept)(nil)).Elem(),
		ID:    SM_TRADE_ACCEPT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_TradeGold)(nil)).Elem(),
		ID:    SM_TRADE_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_TradeItem)(nil)).Elem(),
		ID:    SM_TRADE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_TradeConfirm)(nil)).Elem(),
		ID:    SM_TRADE_CONFIRM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_TradeCancel)(nil)).Elem(),
		ID:    SM_TRADE_CANCEL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MountUpdate)(nil)).Elem(),
		ID:    SM_MOUNT_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_EquipSlotItem)(nil)).Elem(),
		ID:    SM_EQUIP_SLOT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_FishingUpdate)(nil)).Elem(),
		ID:    SM_FISHING_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ChangeQuest)(nil)).Elem(),
		ID:    SM_CHANGE_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_CompleteQuest)(nil)).Elem(),
		ID:    SM_COMPLETE_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ShareQuest)(nil)).Elem(),
		ID:    SM_SHARE_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NewQuestInfo)(nil)).Elem(),
		ID:    SM_NEW_QUEST_INFO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GainedQuestItem)(nil)).Elem(),
		ID:    SM_GAINED_QUEST_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_DeleteQuestItem)(nil)).Elem(),
		ID:    SM_DELETE_QUEST_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_CancelReincarnation)(nil)).Elem(),
		ID:    SM_CANCEL_REINCARNATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_RequestReincarnation)(nil)).Elem(),
		ID:    SM_REQUEST_REINCARNATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_UserBackStep)(nil)).Elem(),
		ID:    SM_USER_BACK_STEP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectBackStep)(nil)).Elem(),
		ID:    SM_OBJECT_BACK_STEP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_UserDashAttack)(nil)).Elem(),
		ID:    SM_USER_DASH_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectDashAttack)(nil)).Elem(),
		ID:    SM_OBJECT_DASH_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_UserAttackMove)(nil)).Elem(),
		ID:    SM_USER_ATTACK_MOVE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_CombineItem)(nil)).Elem(),
		ID:    SM_COMBINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ItemUpgraded)(nil)).Elem(),
		ID:    SM_ITEM_UPGRADED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_SetConcentration)(nil)).Elem(),
		ID:    SM_SET_CONCENTRATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_SetObjectConcentration)(nil)).Elem(),
		ID:    SM_SET_OBJECT_CONCENTRATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_SetElemental)(nil)).Elem(),
		ID:    SM_SET_ELEMENTAL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_SetObjectElemental)(nil)).Elem(),
		ID:    SM_SET_OBJECT_ELEMENTAL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_RemoveDelayedExplosion)(nil)).Elem(),
		ID:    SM_REMOVE_DELAYED_EXPLOSION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectDeco)(nil)).Elem(),
		ID:    SM_OBJECT_DECO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectSneaking)(nil)).Elem(),
		ID:    SM_OBJECT_SNEAKING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ObjectLevelEffects)(nil)).Elem(),
		ID:    SM_OBJECT_LEVEL_EFFECTS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_SetBindingShot)(nil)).Elem(),
		ID:    SM_SET_BINDING_SHOT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_SendOutputMessage)(nil)).Elem(),
		ID:    SM_SEND_OUTPUT_MESSAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCAwakening)(nil)).Elem(),
		ID:    SM_NPC_AWAKENING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCDisassemble)(nil)).Elem(),
		ID:    SM_NPC_DISASSEMBLE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCDowngrade)(nil)).Elem(),
		ID:    SM_NPC_DOWNGRADE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCReset)(nil)).Elem(),
		ID:    SM_NPC_RESET,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_AwakeningNeedMaterials)(nil)).Elem(),
		ID:    SM_AWAKENING_NEED_MATERIALS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_AwakeningLockedItem)(nil)).Elem(),
		ID:    SM_AWAKENING_LOCKED_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_Awakening)(nil)).Elem(),
		ID:    SM_AWAKENING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ReceiveMail)(nil)).Elem(),
		ID:    SM_RECEIVE_MAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MailLockedItem)(nil)).Elem(),
		ID:    SM_MAIL_LOCKED_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MailSendRequest)(nil)).Elem(),
		ID:    SM_MAIL_SEND_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MailSent)(nil)).Elem(),
		ID:    SM_MAIL_SENT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ParcelCollected)(nil)).Elem(),
		ID:    SM_PARCEL_COLLECTED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MailCost)(nil)).Elem(),
		ID:    SM_MAIL_COST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ResizeInventory)(nil)).Elem(),
		ID:    SM_RESIZE_INVENTORY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ResizeStorage)(nil)).Elem(),
		ID:    SM_RESIZE_STORAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NewIntelligentCreature)(nil)).Elem(),
		ID:    SM_NEW_INTELLIGENT_CREATURE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_UpdateIntelligentCreatureList)(nil)).Elem(),
		ID:    SM_UPDATE_INTELLIGENT_CREATURElIST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_IntelligentCreatureEnableRename)(nil)).Elem(),
		ID:    SM_INTELLIGENT_CREATURE_ENABLE_RENAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_IntelligentCreaturePickup)(nil)).Elem(),
		ID:    SM_INTELLIGENT_CREATURE_PICKUP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCPearlGoods)(nil)).Elem(),
		ID:    SM_NPC_PEARL_GOODS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_TransformUpdate)(nil)).Elem(),
		ID:    SM_TRANSFORM_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_FriendUpdate)(nil)).Elem(),
		ID:    SM_FRIEND_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_LoverUpdate)(nil)).Elem(),
		ID:    SM_LOVER_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_MentorUpdate)(nil)).Elem(),
		ID:    SM_MENTOR_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GuildBuffList)(nil)).Elem(),
		ID:    SM_GUILD_BUFF_LIST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NPCRequestInput)(nil)).Elem(),
		ID:    SM_NPC_REQUEST_INPUT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GameShopInfo)(nil)).Elem(),
		ID:    SM_GAME_SHOP_INFO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GameShopStock)(nil)).Elem(),
		ID:    SM_GAME_SHOP_STOCK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_Rankings)(nil)).Elem(),
		ID:    SM_RANKINGS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_Opendoor)(nil)).Elem(),
		ID:    SM_OPENDOOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_GetRentedItems)(nil)).Elem(),
		ID:    SM_GET_RENTED_ITEMS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ItemRentalRequest)(nil)).Elem(),
		ID:    SM_ITEM_RENTAL_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ItemRentalFee)(nil)).Elem(),
		ID:    SM_ITEM_RENTAL_FEE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ItemRentalPeriod)(nil)).Elem(),
		ID:    SM_ITEM_RENTAL_PERIOD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_DepositRentalItem)(nil)).Elem(),
		ID:    SM_DEPOSIT_RENTAL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_RetrieveRentalItem)(nil)).Elem(),
		ID:    SM_RETRIEVE_RENTAL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_UpdateRentalItem)(nil)).Elem(),
		ID:    SM_UPDATE_RENTAL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_CancelItemRental)(nil)).Elem(),
		ID:    SM_CANCEL_ITEM_RENTAL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ItemRentalLock)(nil)).Elem(),
		ID:    SM_ITEM_RENTAL_LOCK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ItemRentalPartnerLock)(nil)).Elem(),
		ID:    SM_ITEM_RENTAL_PARTNER_LOCK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_CanConfirmItemRental)(nil)).Elem(),
		ID:    SM_CAN_CONFIRM_ITEM_RENTAL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_ConfirmItemRental)(nil)).Elem(),
		ID:    SM_CONFIRM_ITEM_RENTAL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_NewRecipeInfo)(nil)).Elem(),
		ID:    SM_NEW_RECIPE_INFO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*SM_OpenBrowser)(nil)).Elem(),
		ID:    SM_OPEN_BROWSER,
	})

}
