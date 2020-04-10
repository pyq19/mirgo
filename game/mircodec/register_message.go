package mircodec

import (
	"reflect"

	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/game/proto/client"
	"github.com/yenkeia/mirgo/game/proto/server"
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
		Type:  reflect.TypeOf((*client.ClientVersion)(nil)).Elem(),
		ID:    client.CLIENT_VERSION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.Disconnect)(nil)).Elem(),
		ID:    client.DISCONNECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.KeepAlive)(nil)).Elem(),
		ID:    client.KEEP_ALIVE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.NewAccount)(nil)).Elem(),
		ID:    client.NEW_ACCOUNT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ChangePassword)(nil)).Elem(),
		ID:    client.CHANGE_PASSWORD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.Login)(nil)).Elem(),
		ID:    client.LOGIN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.NewCharacter)(nil)).Elem(),
		ID:    client.NEW_CHARACTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.DeleteCharacter)(nil)).Elem(),
		ID:    client.DELETE_CHARACTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.StartGame)(nil)).Elem(),
		ID:    client.START_GAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.LogOut)(nil)).Elem(),
		ID:    client.LOG_OUT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.Turn)(nil)).Elem(),
		ID:    client.TURN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.Walk)(nil)).Elem(),
		ID:    client.WALK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.Run)(nil)).Elem(),
		ID:    client.RUN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.Chat)(nil)).Elem(),
		ID:    client.CHAT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.MoveItem)(nil)).Elem(),
		ID:    client.MOVE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.StoreItem)(nil)).Elem(),
		ID:    client.STORE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.TakeBackItem)(nil)).Elem(),
		ID:    client.TAKE_BACK_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.MergeItem)(nil)).Elem(),
		ID:    client.MERGE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.EquipItem)(nil)).Elem(),
		ID:    client.EQUIP_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.RemoveItem)(nil)).Elem(),
		ID:    client.REMOVE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.RemoveSlotItem)(nil)).Elem(),
		ID:    client.REMOVE_SLOT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.SplitItem)(nil)).Elem(),
		ID:    client.SPLIT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.UseItem)(nil)).Elem(),
		ID:    client.USE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.DropItem)(nil)).Elem(),
		ID:    client.DROP_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.DepositRefineItem)(nil)).Elem(),
		ID:    client.DEPOSIT_REFINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.RetrieveRefineItem)(nil)).Elem(),
		ID:    client.RETRIEVE_REFINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.RefineCancel)(nil)).Elem(),
		ID:    client.REFINE_CANCEL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.RefineItem)(nil)).Elem(),
		ID:    client.REFINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.CheckRefine)(nil)).Elem(),
		ID:    client.CHECK_REFINE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ReplaceWedRing)(nil)).Elem(),
		ID:    client.REPLACE_WED_RING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.DepositTradeItem)(nil)).Elem(),
		ID:    client.DEPOSIT_TRADE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.RetrieveTradeItem)(nil)).Elem(),
		ID:    client.RETRIEVE_TRADE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.DropGold)(nil)).Elem(),
		ID:    client.DROP_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.PickUp)(nil)).Elem(),
		ID:    client.PICK_UP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.Inspect)(nil)).Elem(),
		ID:    client.INSPECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ChangeAMode)(nil)).Elem(),
		ID:    client.CHANGE_A_MODE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ChangePMode)(nil)).Elem(),
		ID:    client.CHANGE_P_MODE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ChangeTrade)(nil)).Elem(),
		ID:    client.CHANGE_TRADE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.Attack)(nil)).Elem(),
		ID:    client.ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.RangeAttack)(nil)).Elem(),
		ID:    client.RANGE_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.Harvest)(nil)).Elem(),
		ID:    client.HARVEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.CallNPC)(nil)).Elem(),
		ID:    client.CALL_NPC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.TalkMonsterNPC)(nil)).Elem(),
		ID:    client.TALK_MONSTER_NPC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.BuyItem)(nil)).Elem(),
		ID:    client.BUY_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.SellItem)(nil)).Elem(),
		ID:    client.SELL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.CraftItem)(nil)).Elem(),
		ID:    client.CRAFT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.RepairItem)(nil)).Elem(),
		ID:    client.REPAIR_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.BuyItemBack)(nil)).Elem(),
		ID:    client.BUY_ITEM_BACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.SRepairItem)(nil)).Elem(),
		ID:    client.S_REPAIR_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.MagicKey)(nil)).Elem(),
		ID:    client.MAGIC_KEY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.Magic)(nil)).Elem(),
		ID:    client.MAGIC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.SwitchGroup)(nil)).Elem(),
		ID:    client.SWITCH_GROUP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.AddMember)(nil)).Elem(),
		ID:    client.ADD_MEMBER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.DelMember)(nil)).Elem(),
		ID:    client.DEL_MEMBER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.GroupInvite)(nil)).Elem(),
		ID:    client.GROUP_INVITE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.TownRevive)(nil)).Elem(),
		ID:    client.TOWN_REVIVE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.SpellToggle)(nil)).Elem(),
		ID:    client.SPELL_TOGGLE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ConsignItem)(nil)).Elem(),
		ID:    client.CONSIGN_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.MarketSearch)(nil)).Elem(),
		ID:    client.MARKET_SEARCH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.MarketRefresh)(nil)).Elem(),
		ID:    client.MARKET_REFRESH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.MarketPage)(nil)).Elem(),
		ID:    client.MARKET_PAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.MarketBuy)(nil)).Elem(),
		ID:    client.MARKET_BUY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.MarketGetBack)(nil)).Elem(),
		ID:    client.MARKET_GET_BACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.RequestUserName)(nil)).Elem(),
		ID:    client.REQUEST_USER_NAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.RequestChatItem)(nil)).Elem(),
		ID:    client.REQUEST_CHAT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.EditGuildMember)(nil)).Elem(),
		ID:    client.EDIT_GUILD_MEMBER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.EditGuildNotice)(nil)).Elem(),
		ID:    client.EDIT_GUILD_NOTICE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.GuildInvite)(nil)).Elem(),
		ID:    client.GUILD_INVITE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.GuildNameReturn)(nil)).Elem(),
		ID:    client.GUILD_NAME_RETURN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.RequestGuildInfo)(nil)).Elem(),
		ID:    client.REQUEST_GUILD_INFO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.GuildStorageGoldChange)(nil)).Elem(),
		ID:    client.GUILD_STORAGE_GOLD_CHANGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.GuildStorageItemChange)(nil)).Elem(),
		ID:    client.GUILD_STORAGE_ITEM_CHANGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.GuildWarReturn)(nil)).Elem(),
		ID:    client.GUILD_WAR_RETURN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.MarriageRequest)(nil)).Elem(),
		ID:    client.MARRIAGE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.MarriageReply)(nil)).Elem(),
		ID:    client.MARRIAGE_REPLY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ChangeMarriage)(nil)).Elem(),
		ID:    client.CHANGE_MARRIAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.DivorceRequest)(nil)).Elem(),
		ID:    client.DIVORCE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.DivorceReply)(nil)).Elem(),
		ID:    client.DIVORCE_REPLY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.AddMentor)(nil)).Elem(),
		ID:    client.ADD_MENTOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.MentorReply)(nil)).Elem(),
		ID:    client.MENTOR_REPLY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.AllowMentor)(nil)).Elem(),
		ID:    client.ALLOW_MENTOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.CancelMentor)(nil)).Elem(),
		ID:    client.CANCEL_MENTOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.TradeRequest)(nil)).Elem(),
		ID:    client.TRADE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.TradeReply)(nil)).Elem(),
		ID:    client.TRADE_REPLY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.TradeGold)(nil)).Elem(),
		ID:    client.TRADE_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.TradeConfirm)(nil)).Elem(),
		ID:    client.TRADE_CONFIRM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.TradeCancel)(nil)).Elem(),
		ID:    client.TRADE_CANCEL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.EquipSlotItem)(nil)).Elem(),
		ID:    client.EQUIP_SLOT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.FishingCast)(nil)).Elem(),
		ID:    client.FISHING_CAST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.FishingChangeAutocast)(nil)).Elem(),
		ID:    client.FISHING_CHANGE_AUTOCAST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.AcceptQuest)(nil)).Elem(),
		ID:    client.ACCEPT_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.FinishQuest)(nil)).Elem(),
		ID:    client.FINISH_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.AbandonQuest)(nil)).Elem(),
		ID:    client.ABANDON_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ShareQuest)(nil)).Elem(),
		ID:    client.SHARE_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.AcceptReincarnation)(nil)).Elem(),
		ID:    client.ACCEPT_REINCARNATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.CancelReincarnation)(nil)).Elem(),
		ID:    client.CANCEL_REINCARNATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.CombineItem)(nil)).Elem(),
		ID:    client.COMBINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.SetConcentration)(nil)).Elem(),
		ID:    client.SET_CONCENTRATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.AwakeningNeedMaterials)(nil)).Elem(),
		ID:    client.AWAKENING_NEED_MATERIALS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.AwakeningLockedItem)(nil)).Elem(),
		ID:    client.AWAKENING_LOCKED_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.Awakening)(nil)).Elem(),
		ID:    client.AWAKENING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.DisassembleItem)(nil)).Elem(),
		ID:    client.DISASSEMBLE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.DowngradeAwakening)(nil)).Elem(),
		ID:    client.DOWNGRADE_AWAKENING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ResetAddedItem)(nil)).Elem(),
		ID:    client.RESET_ADDED_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.SendMail)(nil)).Elem(),
		ID:    client.SEND_MAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ReadMail)(nil)).Elem(),
		ID:    client.READ_MAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.CollectParcel)(nil)).Elem(),
		ID:    client.COLLECT_PARCEL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.DeleteMail)(nil)).Elem(),
		ID:    client.DELETE_MAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.LockMail)(nil)).Elem(),
		ID:    client.LOCK_MAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.MailLockedItem)(nil)).Elem(),
		ID:    client.MAIL_LOCKED_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.MailCost)(nil)).Elem(),
		ID:    client.MAIL_COST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.UpdateIntelligentCreature)(nil)).Elem(),
		ID:    client.UPDATE_INTELLIGENT_CREATURE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.IntelligentCreaturePickup)(nil)).Elem(),
		ID:    client.INTELLIGENT_CREATURE_PICKUP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.AddFriend)(nil)).Elem(),
		ID:    client.ADD_FRIEND,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.RemoveFriend)(nil)).Elem(),
		ID:    client.REMOVE_FRIEND,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.RefreshFriends)(nil)).Elem(),
		ID:    client.REFRESH_FRIENDS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.AddMemo)(nil)).Elem(),
		ID:    client.ADD_MEMO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.GuildBuffUpdate)(nil)).Elem(),
		ID:    client.GUILD_BUFF_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.NPCConfirmInput)(nil)).Elem(),
		ID:    client.NPC_CONFIRM_INPUT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.GameshopBuy)(nil)).Elem(),
		ID:    client.GAMESHOP_BUY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ReportIssue)(nil)).Elem(),
		ID:    client.REPORT_ISSUE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.GetRanking)(nil)).Elem(),
		ID:    client.GET_RANKING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.Opendoor)(nil)).Elem(),
		ID:    client.OPENDOOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.GetRentedItems)(nil)).Elem(),
		ID:    client.GET_RENTED_ITEMS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ItemRentalRequest)(nil)).Elem(),
		ID:    client.ITEM_RENTAL_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ItemRentalFee)(nil)).Elem(),
		ID:    client.ITEM_RENTAL_FEE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ItemRentalPeriod)(nil)).Elem(),
		ID:    client.ITEM_RENTAL_PERIOD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.DepositRentalItem)(nil)).Elem(),
		ID:    client.DEPOSIT_RENTAL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.RetrieveRentalItem)(nil)).Elem(),
		ID:    client.RETRIEVE_RENTAL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.CancelItemRental)(nil)).Elem(),
		ID:    client.CANCEL_ITEM_RENTAL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ItemRentalLockFee)(nil)).Elem(),
		ID:    client.ITEM_RENTAL_LOCK_FEE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ItemRentalLockItem)(nil)).Elem(),
		ID:    client.ITEM_RENTAL_LOCK_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*client.ConfirmItemRental)(nil)).Elem(),
		ID:    client.CONFIRM_ITEM_RENTAL,
	})
}

func initServerMessage() {

	mirCodec := new(MirCodec)

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.Connected)(nil)).Elem(),
		ID:    server.CONNECTED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ClientVersion)(nil)).Elem(),
		ID:    server.CLIENT_VERSION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.Disconnect)(nil)).Elem(),
		ID:    server.DISCONNECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.KeepAlive)(nil)).Elem(),
		ID:    server.KEEP_ALIVE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NewAccount)(nil)).Elem(),
		ID:    server.NEW_ACCOUNT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ChangePassword)(nil)).Elem(),
		ID:    server.CHANGE_PASSWORD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ChangePasswordBanned)(nil)).Elem(),
		ID:    server.CHANGE_PASSWORD_BANNED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.Login)(nil)).Elem(),
		ID:    server.LOGIN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.LoginBanned)(nil)).Elem(),
		ID:    server.LOGIN_BANNED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.LoginSuccess)(nil)).Elem(),
		ID:    server.LOGIN_SUCCESS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NewCharacter)(nil)).Elem(),
		ID:    server.NEW_CHARACTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NewCharacterSuccess)(nil)).Elem(),
		ID:    server.NEW_CHARACTER_SUCCESS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.DeleteCharacter)(nil)).Elem(),
		ID:    server.DELETE_CHARACTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.DeleteCharacterSuccess)(nil)).Elem(),
		ID:    server.DELETE_CHARACTER_SUCCESS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.StartGame)(nil)).Elem(),
		ID:    server.START_GAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.StartGameBanned)(nil)).Elem(),
		ID:    server.START_GAME_BANNED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.StartGameDelay)(nil)).Elem(),
		ID:    server.START_GAME_DELAY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MapInformation)(nil)).Elem(),
		ID:    server.MAP_INFORMATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.UserInformation)(nil)).Elem(),
		ID:    server.USER_INFORMATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.UserLocation)(nil)).Elem(),
		ID:    server.USER_LOCATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectPlayer)(nil)).Elem(),
		ID:    server.OBJECT_PLAYER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectRemove)(nil)).Elem(),
		ID:    server.OBJECT_REMOVE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectTurn)(nil)).Elem(),
		ID:    server.OBJECT_TURN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectWalk)(nil)).Elem(),
		ID:    server.OBJECT_WALK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectRun)(nil)).Elem(),
		ID:    server.OBJECT_RUN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.Chat)(nil)).Elem(),
		ID:    server.CHAT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectChat)(nil)).Elem(),
		ID:    server.OBJECT_CHAT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NewItemInfo)(nil)).Elem(),
		ID:    server.NEW_ITEM_INFO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MoveItem)(nil)).Elem(),
		ID:    server.MOVE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.EquipItem)(nil)).Elem(),
		ID:    server.EQUIP_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MergeItem)(nil)).Elem(),
		ID:    server.MERGE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.RemoveItem)(nil)).Elem(),
		ID:    server.REMOVE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.RemoveSlotItem)(nil)).Elem(),
		ID:    server.REMOVE_SLOT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.TakeBackItem)(nil)).Elem(),
		ID:    server.TAKE_BACK_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.StoreItem)(nil)).Elem(),
		ID:    server.STORE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.SplitItem)(nil)).Elem(),
		ID:    server.SPLIT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.SplitItem1)(nil)).Elem(),
		ID:    server.SPLIT_ITEM1,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.DepositRefineItem)(nil)).Elem(),
		ID:    server.DEPOSIT_REFINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.RetrieveRefineItem)(nil)).Elem(),
		ID:    server.RETRIEVE_REFINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.RefineCancel)(nil)).Elem(),
		ID:    server.REFINE_CANCEL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.RefineItem)(nil)).Elem(),
		ID:    server.REFINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.DepositTradeItem)(nil)).Elem(),
		ID:    server.DEPOSIT_TRADE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.RetrieveTradeItem)(nil)).Elem(),
		ID:    server.RETRIEVE_TRADE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.UseItem)(nil)).Elem(),
		ID:    server.USE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.DropItem)(nil)).Elem(),
		ID:    server.DROP_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.PlayerUpdate)(nil)).Elem(),
		ID:    server.PLAYER_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.PlayerInspect)(nil)).Elem(),
		ID:    server.PLAYER_INSPECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.LogOutSuccess)(nil)).Elem(),
		ID:    server.LOG_OUT_SUCCESS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.LogOutFailed)(nil)).Elem(),
		ID:    server.LOG_OUT_FAILED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.TimeOfDay)(nil)).Elem(),
		ID:    server.TIME_OF_DAY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ChangeAMode)(nil)).Elem(),
		ID:    server.CHANGE_A_MODE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ChangePMode)(nil)).Elem(),
		ID:    server.CHANGE_P_MODE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectItem)(nil)).Elem(),
		ID:    server.OBJECT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectGold)(nil)).Elem(),
		ID:    server.OBJECT_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GainedItem)(nil)).Elem(),
		ID:    server.GAINED_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GainedGold)(nil)).Elem(),
		ID:    server.GAINED_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.LoseGold)(nil)).Elem(),
		ID:    server.LOSE_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GainedCredit)(nil)).Elem(),
		ID:    server.GAINED_CREDIT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.LoseCredit)(nil)).Elem(),
		ID:    server.LOSE_CREDIT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectMonster)(nil)).Elem(),
		ID:    server.OBJECT_MONSTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectAttack)(nil)).Elem(),
		ID:    server.OBJECT_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.Struck)(nil)).Elem(),
		ID:    server.STRUCK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectStruck)(nil)).Elem(),
		ID:    server.OBJECT_STRUCK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.DamageIndicator)(nil)).Elem(),
		ID:    server.DAMAGE_INDICATOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.DuraChanged)(nil)).Elem(),
		ID:    server.DURA_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.HealthChanged)(nil)).Elem(),
		ID:    server.HEALTH_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.DeleteItem)(nil)).Elem(),
		ID:    server.DELETE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.Death)(nil)).Elem(),
		ID:    server.DEATH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectDied)(nil)).Elem(),
		ID:    server.OBJECT_DIED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ColourChanged)(nil)).Elem(),
		ID:    server.COLOUR_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectColourChanged)(nil)).Elem(),
		ID:    server.OBJECT_COLOUR_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectGuildNameChanged)(nil)).Elem(),
		ID:    server.OBJECT_GUILD_NAME_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GainExperience)(nil)).Elem(),
		ID:    server.GAIN_EXPERIENCE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.LevelChanged)(nil)).Elem(),
		ID:    server.LEVEL_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectLeveled)(nil)).Elem(),
		ID:    server.OBJECT_LEVELED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectHarvest)(nil)).Elem(),
		ID:    server.OBJECT_HARVEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectHarvested)(nil)).Elem(),
		ID:    server.OBJECT_HARVESTED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectNPC)(nil)).Elem(),
		ID:    server.OBJECT_NPC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCResponse)(nil)).Elem(),
		ID:    server.NPC_RESPONSE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectHide)(nil)).Elem(),
		ID:    server.OBJECT_HIDE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectShow)(nil)).Elem(),
		ID:    server.OBJECT_SHOW,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.Poisoned)(nil)).Elem(),
		ID:    server.POISONED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectPoisoned)(nil)).Elem(),
		ID:    server.OBJECT_POISONED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MapChanged)(nil)).Elem(),
		ID:    server.MAP_CHANGED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectTeleportOut)(nil)).Elem(),
		ID:    server.OBJECT_TELEPORT_OUT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectTeleportIn)(nil)).Elem(),
		ID:    server.OBJECT_TELEPORT_IN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.TeleportIn)(nil)).Elem(),
		ID:    server.TELEPORT_IN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCGoods)(nil)).Elem(),
		ID:    server.NPC_GOODS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCSell)(nil)).Elem(),
		ID:    server.NPC_SELL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCRepair)(nil)).Elem(),
		ID:    server.NPC_REPAIR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCSRepair)(nil)).Elem(),
		ID:    server.NPC_S_REPAIR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCRefine)(nil)).Elem(),
		ID:    server.NPC_REFINE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCCheckRefine)(nil)).Elem(),
		ID:    server.NPC_CHECK_REFINE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCCollectRefine)(nil)).Elem(),
		ID:    server.NPC_COLLECT_REFINE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCReplaceWedRing)(nil)).Elem(),
		ID:    server.NPC_REPLACE_WED_RING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCStorage)(nil)).Elem(),
		ID:    server.NPC_STORAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.SellItem)(nil)).Elem(),
		ID:    server.SELL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.CraftItem)(nil)).Elem(),
		ID:    server.CRAFT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.RepairItem)(nil)).Elem(),
		ID:    server.REPAIR_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ItemRepaired)(nil)).Elem(),
		ID:    server.ITEM_REPAIRED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NewMagic)(nil)).Elem(),
		ID:    server.NEW_MAGIC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.RemoveMagic)(nil)).Elem(),
		ID:    server.REMOVE_MAGIC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MagicLeveled)(nil)).Elem(),
		ID:    server.MAGIC_LEVELED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.Magic)(nil)).Elem(),
		ID:    server.MAGIC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MagicDelay)(nil)).Elem(),
		ID:    server.MAGIC_DELAY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MagicCast)(nil)).Elem(),
		ID:    server.MAGIC_CAST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectMagic)(nil)).Elem(),
		ID:    server.OBJECT_MAGIC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectEffect)(nil)).Elem(),
		ID:    server.OBJECT_EFFECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.RangeAttack)(nil)).Elem(),
		ID:    server.RANGE_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.Pushed)(nil)).Elem(),
		ID:    server.PUSHED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectPushed)(nil)).Elem(),
		ID:    server.OBJECT_PUSHED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectName)(nil)).Elem(),
		ID:    server.OBJECT_NAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.UserStorage)(nil)).Elem(),
		ID:    server.USER_STORAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.SwitchGroup)(nil)).Elem(),
		ID:    server.SWITCH_GROUP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.DeleteGroup)(nil)).Elem(),
		ID:    server.DELETE_GROUP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.DeleteMember)(nil)).Elem(),
		ID:    server.DELETE_MEMBER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GroupInvite)(nil)).Elem(),
		ID:    server.GROUP_INVITE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.AddMember)(nil)).Elem(),
		ID:    server.ADD_MEMBER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.Revived)(nil)).Elem(),
		ID:    server.REVIVED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectRevived)(nil)).Elem(),
		ID:    server.OBJECT_REVIVED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.SpellToggle)(nil)).Elem(),
		ID:    server.SPELL_TOGGLE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectHealth)(nil)).Elem(),
		ID:    server.OBJECT_HEALTH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MapEffect)(nil)).Elem(),
		ID:    server.MAP_EFFECT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectRangeAttack)(nil)).Elem(),
		ID:    server.OBJECT_RANGE_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.AddBuff)(nil)).Elem(),
		ID:    server.ADD_BUFF,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.RemoveBuff)(nil)).Elem(),
		ID:    server.REMOVE_BUFF,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectHidden)(nil)).Elem(),
		ID:    server.OBJECT_HIDDEN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.RefreshItem)(nil)).Elem(),
		ID:    server.REFRESH_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectSpell)(nil)).Elem(),
		ID:    server.OBJECT_SPELL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.UserDash)(nil)).Elem(),
		ID:    server.USER_DASH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectDash)(nil)).Elem(),
		ID:    server.OBJECT_DASH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.UserDashFail)(nil)).Elem(),
		ID:    server.USER_DASH_FAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectDashFail)(nil)).Elem(),
		ID:    server.OBJECT_DASH_FAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCConsign)(nil)).Elem(),
		ID:    server.NPC_CONSIGN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCMarket)(nil)).Elem(),
		ID:    server.NPC_MARKET,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCMarketPage)(nil)).Elem(),
		ID:    server.NPC_MARKET_PAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ConsignItem)(nil)).Elem(),
		ID:    server.CONSIGN_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MarketFail)(nil)).Elem(),
		ID:    server.MARKET_FAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MarketSuccess)(nil)).Elem(),
		ID:    server.MARKET_SUCCESS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectSitDown)(nil)).Elem(),
		ID:    server.OBJECT_SIT_DOWN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.InTrapRock)(nil)).Elem(),
		ID:    server.IN_TRAP_ROCK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.BaseStatsInfo)(nil)).Elem(),
		ID:    server.BASE_STATS_INFO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.UserName)(nil)).Elem(),
		ID:    server.USER_NAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ChatItemStats)(nil)).Elem(),
		ID:    server.CHAT_ITEM_STATS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GuildNoticeChange)(nil)).Elem(),
		ID:    server.GUILD_NOTICE_CHANGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GuildMemberChange)(nil)).Elem(),
		ID:    server.GUILD_MEMBER_CHANGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GuildStatus)(nil)).Elem(),
		ID:    server.GUILD_STATUS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GuildInvite)(nil)).Elem(),
		ID:    server.GUILD_INVITE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GuildExpGain)(nil)).Elem(),
		ID:    server.GUILD_EXP_GAIN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GuildNameRequest)(nil)).Elem(),
		ID:    server.GUILD_NAME_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GuildStorageGoldChange)(nil)).Elem(),
		ID:    server.GUILD_STORAGE_GOLD_CHANGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GuildStorageItemChange)(nil)).Elem(),
		ID:    server.GUILD_STORAGE_ITEM_CHANGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GuildStorageList)(nil)).Elem(),
		ID:    server.GUILD_STORAGE_LIST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GuildRequestWar)(nil)).Elem(),
		ID:    server.GUILD_REQUEST_WAR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.DefaultNPC)(nil)).Elem(),
		ID:    server.DEFAULT_NPC,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCUpdate)(nil)).Elem(),
		ID:    server.NPC_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCImageUpdate)(nil)).Elem(),
		ID:    server.NPC_IMAGE_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MarriageRequest)(nil)).Elem(),
		ID:    server.MARRIAGE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.DivorceRequest)(nil)).Elem(),
		ID:    server.DIVORCE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MentorRequest)(nil)).Elem(),
		ID:    server.MENTOR_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.TradeRequest)(nil)).Elem(),
		ID:    server.TRADE_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.TradeAccept)(nil)).Elem(),
		ID:    server.TRADE_ACCEPT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.TradeGold)(nil)).Elem(),
		ID:    server.TRADE_GOLD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.TradeItem)(nil)).Elem(),
		ID:    server.TRADE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.TradeConfirm)(nil)).Elem(),
		ID:    server.TRADE_CONFIRM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.TradeCancel)(nil)).Elem(),
		ID:    server.TRADE_CANCEL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MountUpdate)(nil)).Elem(),
		ID:    server.MOUNT_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.EquipSlotItem)(nil)).Elem(),
		ID:    server.EQUIP_SLOT_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.FishingUpdate)(nil)).Elem(),
		ID:    server.FISHING_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ChangeQuest)(nil)).Elem(),
		ID:    server.CHANGE_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.CompleteQuest)(nil)).Elem(),
		ID:    server.COMPLETE_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ShareQuest)(nil)).Elem(),
		ID:    server.SHARE_QUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NewQuestInfo)(nil)).Elem(),
		ID:    server.NEW_QUEST_INFO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GainedQuestItem)(nil)).Elem(),
		ID:    server.GAINED_QUEST_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.DeleteQuestItem)(nil)).Elem(),
		ID:    server.DELETE_QUEST_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.CancelReincarnation)(nil)).Elem(),
		ID:    server.CANCEL_REINCARNATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.RequestReincarnation)(nil)).Elem(),
		ID:    server.REQUEST_REINCARNATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.UserBackStep)(nil)).Elem(),
		ID:    server.USER_BACK_STEP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectBackStep)(nil)).Elem(),
		ID:    server.OBJECT_BACK_STEP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.UserDashAttack)(nil)).Elem(),
		ID:    server.USER_DASH_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectDashAttack)(nil)).Elem(),
		ID:    server.OBJECT_DASH_ATTACK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.UserAttackMove)(nil)).Elem(),
		ID:    server.USER_ATTACK_MOVE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.CombineItem)(nil)).Elem(),
		ID:    server.COMBINE_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ItemUpgraded)(nil)).Elem(),
		ID:    server.ITEM_UPGRADED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.SetConcentration)(nil)).Elem(),
		ID:    server.SET_CONCENTRATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.SetObjectConcentration)(nil)).Elem(),
		ID:    server.SET_OBJECT_CONCENTRATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.SetElemental)(nil)).Elem(),
		ID:    server.SET_ELEMENTAL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.SetObjectElemental)(nil)).Elem(),
		ID:    server.SET_OBJECT_ELEMENTAL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.RemoveDelayedExplosion)(nil)).Elem(),
		ID:    server.REMOVE_DELAYED_EXPLOSION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectDeco)(nil)).Elem(),
		ID:    server.OBJECT_DECO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectSneaking)(nil)).Elem(),
		ID:    server.OBJECT_SNEAKING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ObjectLevelEffects)(nil)).Elem(),
		ID:    server.OBJECT_LEVEL_EFFECTS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.SetBindingShot)(nil)).Elem(),
		ID:    server.SET_BINDING_SHOT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.SendOutputMessage)(nil)).Elem(),
		ID:    server.SEND_OUTPUT_MESSAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCAwakening)(nil)).Elem(),
		ID:    server.NPC_AWAKENING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCDisassemble)(nil)).Elem(),
		ID:    server.NPC_DISASSEMBLE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCDowngrade)(nil)).Elem(),
		ID:    server.NPC_DOWNGRADE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCReset)(nil)).Elem(),
		ID:    server.NPC_RESET,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.AwakeningNeedMaterials)(nil)).Elem(),
		ID:    server.AWAKENING_NEED_MATERIALS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.AwakeningLockedItem)(nil)).Elem(),
		ID:    server.AWAKENING_LOCKED_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.Awakening)(nil)).Elem(),
		ID:    server.AWAKENING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ReceiveMail)(nil)).Elem(),
		ID:    server.RECEIVE_MAIL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MailLockedItem)(nil)).Elem(),
		ID:    server.MAIL_LOCKED_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MailSendRequest)(nil)).Elem(),
		ID:    server.MAIL_SEND_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MailSent)(nil)).Elem(),
		ID:    server.MAIL_SENT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ParcelCollected)(nil)).Elem(),
		ID:    server.PARCEL_COLLECTED,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MailCost)(nil)).Elem(),
		ID:    server.MAIL_COST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ResizeInventory)(nil)).Elem(),
		ID:    server.RESIZE_INVENTORY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ResizeStorage)(nil)).Elem(),
		ID:    server.RESIZE_STORAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NewIntelligentCreature)(nil)).Elem(),
		ID:    server.NEW_INTELLIGENT_CREATURE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.UpdateIntelligentCreatureList)(nil)).Elem(),
		ID:    server.UPDATE_INTELLIGENT_CREATURElIST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.IntelligentCreatureEnableRename)(nil)).Elem(),
		ID:    server.INTELLIGENT_CREATURE_ENABLE_RENAME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.IntelligentCreaturePickup)(nil)).Elem(),
		ID:    server.INTELLIGENT_CREATURE_PICKUP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCPearlGoods)(nil)).Elem(),
		ID:    server.NPC_PEARL_GOODS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.TransformUpdate)(nil)).Elem(),
		ID:    server.TRANSFORM_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.FriendUpdate)(nil)).Elem(),
		ID:    server.FRIEND_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.LoverUpdate)(nil)).Elem(),
		ID:    server.LOVER_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.MentorUpdate)(nil)).Elem(),
		ID:    server.MENTOR_UPDATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GuildBuffList)(nil)).Elem(),
		ID:    server.GUILD_BUFF_LIST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NPCRequestInput)(nil)).Elem(),
		ID:    server.NPC_REQUEST_INPUT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GameShopInfo)(nil)).Elem(),
		ID:    server.GAME_SHOP_INFO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GameShopStock)(nil)).Elem(),
		ID:    server.GAME_SHOP_STOCK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.Rankings)(nil)).Elem(),
		ID:    server.RANKINGS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.Opendoor)(nil)).Elem(),
		ID:    server.OPENDOOR,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.GetRentedItems)(nil)).Elem(),
		ID:    server.GET_RENTED_ITEMS,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ItemRentalRequest)(nil)).Elem(),
		ID:    server.ITEM_RENTAL_REQUEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ItemRentalFee)(nil)).Elem(),
		ID:    server.ITEM_RENTAL_FEE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ItemRentalPeriod)(nil)).Elem(),
		ID:    server.ITEM_RENTAL_PERIOD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.DepositRentalItem)(nil)).Elem(),
		ID:    server.DEPOSIT_RENTAL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.RetrieveRentalItem)(nil)).Elem(),
		ID:    server.RETRIEVE_RENTAL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.UpdateRentalItem)(nil)).Elem(),
		ID:    server.UPDATE_RENTAL_ITEM,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.CancelItemRental)(nil)).Elem(),
		ID:    server.CANCEL_ITEM_RENTAL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ItemRentalLock)(nil)).Elem(),
		ID:    server.ITEM_RENTAL_LOCK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ItemRentalPartnerLock)(nil)).Elem(),
		ID:    server.ITEM_RENTAL_PARTNER_LOCK,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.CanConfirmItemRental)(nil)).Elem(),
		ID:    server.CAN_CONFIRM_ITEM_RENTAL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.ConfirmItemRental)(nil)).Elem(),
		ID:    server.CONFIRM_ITEM_RENTAL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.NewRecipeInfo)(nil)).Elem(),
		ID:    server.NEW_RECIPE_INFO,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: mirCodec,
		Type:  reflect.TypeOf((*server.OpenBrowser)(nil)).Elem(),
		ID:    server.OPEN_BROWSER,
	})

}
