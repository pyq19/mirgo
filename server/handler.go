package main

import "C"
import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/golog"
	"github.com/yenkeia/mirgo/codec/mircodec"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/client"
	"github.com/yenkeia/mirgo/proto/server"
)

var log = golog.New("server.handler")

func (g *Game) HandleEvent(ev cellnet.Event) {

	var s cellnet.Session
	s = ev.Session()
	switch msg := ev.Message().(type) {

	// 有新的连接
	case *cellnet.SessionAccepted:
		g.SessionAccepted(s, msg)
	// 有连接断开
	case *cellnet.SessionClosed:
		g.SessionClosed(s, msg)
	case *client.ClientVersion:
		g.ClientVersion(s, msg)
	case *client.KeepAlive:
		g.KeepAlive(s, msg)
	case *client.NewAccount:
		g.NewAccount(s, msg)
	case *client.ChangePassword:
		g.ChangePassword(s, msg)
	case *client.Login:
		g.Login(s, msg)
	case *client.NewCharacter:
		g.NewCharacter(s, msg)
	case *client.DeleteCharacter:
		g.DeleteCharacter(s, msg)
	case *client.StartGame:
		g.StartGame(s, msg)
	case *client.LogOut:
		g.LogOut(s, msg)
	case *client.Turn:
		g.Turn(s, msg)
	case *client.Walk:
		g.Walk(s, msg)
	case *client.Run:
		g.Run(s, msg)
	case *client.Chat:
		g.Chat(s, msg)
	case *client.MoveItem:
		g.MoveItem(s, msg)
	case *client.StoreItem:
		g.StoreItem(s, msg)
	case *client.DepositRefineItem:
		g.DepositRefineItem(s, msg)
	case *client.RetrieveRefineItem:
		g.RetrieveRefineItem(s, msg)
	case *client.RefineCancel:
		g.RefineCancel(s, msg)
	case *client.RefineItem:
		g.RefineItem(s, msg)
	case *client.CheckRefine:
		g.CheckRefine(s, msg)
	case *client.ReplaceWedRing:
		g.ReplaceWedRing(s, msg)
	case *client.DepositTradeItem:
		g.DepositTradeItem(s, msg)
	case *client.RetrieveTradeItem:
		g.RetrieveTradeItem(s, msg)
	case *client.TakeBackItem:
		g.TakeBackItem(s, msg)
	case *client.MergeItem:
		g.MergeItem(s, msg)
	case *client.EquipItem:
		g.EquipItem(s, msg)
	case *client.RemoveItem:
		g.RemoveItem(s, msg)
	case *client.RemoveSlotItem:
		g.RemoveSlotItem(s, msg)
	case *client.SplitItem:
		g.SplitItem(s, msg)
	case *client.UseItem:
		g.UseItem(s, msg)
	case *client.DropItem:
		g.DropItem(s, msg)
	case *client.DropGold:
		g.DropGold(s, msg)
	case *client.PickUp:
		g.PickUp(s, msg)
	case *client.Inspect:
		g.Inspect(s, msg)
	case *client.ChangeAMode:
		g.ChangeAMode(s, msg)
	case *client.ChangePMode:
		g.ChangePMode(s, msg)
	case *client.ChangeTrade:
		g.ChangeTrade(s, msg)
	case *client.Attack:
		g.Attack(s, msg)
	case *client.RangeAttack:
		g.RangeAttack(s, msg)
	case *client.Harvest:
		g.Harvest(s, msg)
	case *client.CallNPC:
		g.CallNPC(s, msg)
	case *client.TalkMonsterNPC:
		g.TalkMonsterNPC(s, msg)
	case *client.BuyItem:
		g.BuyItem(s, msg)
	case *client.CraftItem:
		g.CraftItem(s, msg)
	case *client.SellItem:
		g.SellItem(s, msg)
	case *client.RepairItem:
		g.RepairItem(s, msg)
	case *client.BuyItemBack:
		g.BuyItemBack(s, msg)
	case *client.SRepairItem:
		g.SRepairItem(s, msg)
	case *client.MagicKey:
		g.MagicKey(s, msg)
	case *client.Magic:
		g.Magic(s, msg)
	case *client.SwitchGroup:
		g.SwitchGroup(s, msg)
	case *client.AddMember:
		g.AddMember(s, msg)
	case *client.DellMember:
		g.DelMember(s, msg)
	case *client.GroupInvite:
		g.GroupInvite(s, msg)
	case *client.TownRevive:
		g.TownRevive(s, msg)
	case *client.SpellToggle:
		g.SpellToggle(s, msg)
	case *client.ConsignItem:
		g.ConsignItem(s, msg)
	case *client.MarketSearch:
		g.MarketSearch(s, msg)
	case *client.MarketRefresh:
		g.MarketRefresh(s, msg)
	case *client.MarketPage:
		g.MarketPage(s, msg)
	case *client.MarketBuy:
		g.MarketBuy(s, msg)
	case *client.MarketGetBack:
		g.MarketGetBack(s, msg)
	case *client.RequestUserName:
		g.RequestUserName(s, msg)
	case *client.RequestChatItem:
		g.RequestChatItem(s, msg)
	case *client.EditGuildMember:
		g.EditGuildMember(s, msg)
	case *client.EditGuildNotice:
		g.EditGuildNotice(s, msg)
	case *client.GuildInvite:
		g.GuildInvite(s, msg)
	case *client.RequestGuildInfo:
		g.RequestGuildInfo(s, msg)
	case *client.GuildNameReturn:
		g.GuildNameReturn(s, msg)
	case *client.GuildStorageGoldChange:
		g.GuildStorageGoldChange(s, msg)
	case *client.GuildStorageItemChange:
		g.GuildStorageItemChange(s, msg)
	case *client.GuildWarReturn:
		g.GuildWarReturn(s, msg)
	case *client.MarriageRequest:
		g.MarriageRequest(s, msg)
	case *client.MarriageReply:
		g.MarriageReply(s, msg)
	case *client.ChangeMarriage:
		g.ChangeMarriage(s, msg)
	case *client.DivorceRequest:
		g.DivorceRequest(s, msg)
	case *client.DivorceReply:
		g.DivorceReply(s, msg)
	case *client.AddMentor:
		g.AddMentor(s, msg)
	case *client.MentorReply:
		g.MentorReply(s, msg)
	case *client.AllowMentor:
		g.AllowMentor(s, msg)
	case *client.CancelMentor:
		g.CancelMentor(s, msg)
	case *client.TradeRequest:
		g.TradeRequest(s, msg)
	case *client.TradeGold:
		g.TradeGold(s, msg)
	case *client.TradeReply:
		g.TradeReply(s, msg)
	case *client.TradeConfirm:
		g.TradeConfirm(s, msg)
	case *client.TradeCancel:
		g.TradeCancel(s, msg)
	case *client.EquipSlotItem:
		g.EquipSlotItem(s, msg)
	case *client.FishingCast:
		g.FishingCast(s, msg)
	case *client.FishingChangeAutocast:
		g.FishingChangeAutocast(s, msg)
	case *client.AcceptQuest:
		g.AcceptQuest(s, msg)
	case *client.FinishQuest:
		g.FinishQuest(s, msg)
	case *client.AbandonQuest:
		g.AbandonQuest(s, msg)
	case *client.ShareQuest:
		g.ShareQuest(s, msg)
	case *client.AcceptReincarnation:
		g.AcceptReincarnation(s, msg)
	case *client.CancelReincarnation:
		g.CancelReincarnation(s, msg)
	case *client.CombineItem:
		g.CombineItem(s, msg)
	case *client.SetConcentration:
		g.SetConcentration(s, msg)
	case *client.AwakeningNeedMaterials:
		g.AwakeningNeedMaterials(s, msg)
	case *client.AwakeningLockedItem:
		g.AwakeningLockedItem(s, msg)
	case *client.Awakening:
		g.Awakening(s, msg)
	case *client.DisassembleItem:
		g.DisassembleItem(s, msg)
	case *client.DowngradeAwakening:
		g.DowngradeAwakening(s, msg)
	case *client.ResetAddedItem:
		g.ResetAddedItem(s, msg)
	case *client.SendMail:
		g.SendMail(s, msg)
	case *client.ReadMail:
		g.ReadMail(s, msg)
	case *client.CollectParcel:
		g.CollectParcel(s, msg)
	case *client.DeleteMail:
		g.DeleteMail(s, msg)
	case *client.LockMail:
		g.LockMail(s, msg)
	case *client.MailLockedItem:
		g.MailLockedItem(s, msg)
	case *client.MailCost:
		g.MailCost(s, msg)
	case *client.UpdateIntelligentCreature: //IntelligentCreature
		g.UpdateIntelligentCreature(s, msg)
	case *client.IntelligentCreaturePickup: //IntelligentCreature
		g.IntelligentCreaturePickup(s, msg)
	case *client.AddFriend:
		g.AddFriend(s, msg)
	case *client.RemoveFriend:
		g.RemoveFriend(s, msg)
	case *client.RefreshFriends:
		g.RefreshFriends(s, msg)
	case *client.AddMemo:
		g.AddMemo(s, msg)
	case *client.GuildBuffUpdate:
		g.GuildBuffUpdate(s, msg)
	case *client.GameshopBuy:
		g.GameshopBuy(s, msg)
	case *client.NPCConfirmInput:
		g.NPCConfirmInput(s, msg)
	case *client.ReportIssue:
		g.ReportIssue(s, msg)
	case *client.GetRanking:
		g.GetRanking(s, msg)
	case *client.Opendoor:
		g.Opendoor(s, msg)
	case *client.GetRentedItems:
		g.GetRentedItems(s, msg)
	case *client.ItemRentalRequest:
		g.ItemRentalRequest(s, msg)
	case *client.ItemRentalFee:
		g.ItemRentalFee(s, msg)
	case *client.ItemRentalPeriod:
		g.ItemRentalPeriod(s, msg)
	case *client.DepositRentalItem:
		g.DepositRentalItem(s, msg)
	case *client.RetrieveRentalItem:
		g.RetrieveRentalItem(s, msg)
	case *client.CancelItemRental:
		g.CancelItemRental(s, msg)
	case *client.ItemRentalLockFee:
		g.ItemRentalLockFee(s, msg)
	case *client.ItemRentalLockItem:
		g.ItemRentalLockItem(s, msg)
	case *client.ConfirmItemRental:
		g.ConfirmItemRental(s, msg)
	default:
		log.Debugln("default:", msg)
		//MessageQueue.Enqueue(string.Format("Invalid packet received. Index : {0}", p.Index));
	}
}

func (g *Game) SessionAccepted(s cellnet.Session, msg *cellnet.SessionAccepted) {
	connected := server.Connected{}
	s.Send(&connected)
}

func (g *Game) SessionClosed(s cellnet.Session, msg *cellnet.SessionClosed) {

}

func (g *Game) ClientVersion(s cellnet.Session, msg *client.ClientVersion) {
	clientVersion := server.ClientVersion{Result: 1}
	s.Send(&clientVersion)
}

func (g *Game) KeepAlive(s cellnet.Session, msg *client.KeepAlive) {
	keepAlive := server.KeepAlive{Time: 0}
	s.Send(keepAlive)
}

// TODO 保存新账号
func (g *Game) NewAccount(s cellnet.Session, msg *client.NewAccount) {
	log.Debugln(msg.AccountID, msg.Password)
	s.Send(server.NewAccount{8})
}

func (g *Game) ChangePassword(s cellnet.Session, msg *client.NewAccount) {

}

// TODO 登陆
func (g *Game) Login(s cellnet.Session, msg *client.Login) {
	res := new(server.LoginSuccess)

	c1 := new(common.SelectInfo)
	c1.Name = "测试登陆1"
	c1.Index = 1
	c1.Gender = common.MirGenderFemale
	c1.Class = common.MirClassArcher
	res.Characters = append(res.Characters, *c1)

	c2 := new(common.SelectInfo)
	c2.Name = "测试登陆2"
	c2.Index = 2
	c2.Gender = common.MirGenderFemale
	c2.Class = common.MirClassAssassin
	res.Characters = append(res.Characters, *c2)

	s.Send(res)
}

// TODO 创建角色成功
func (g *Game) NewCharacter(s cellnet.Session, msg *client.NewCharacter) {
	log.Debugln(msg.Name, msg.Class, msg.Gender)
	res := new(server.NewCharacterSuccess)
	res.CharInfo.Index = 0
	res.CharInfo.Name = msg.Name
	res.CharInfo.Class = msg.Class
	res.CharInfo.Gender = msg.Gender
	s.Send(res)
}

func (g *Game) DeleteCharacter(s cellnet.Session, msg *client.DeleteCharacter) {

}

// TODO 开始游戏
func (g *Game) StartGame(s cellnet.Session, msg *client.StartGame) {

	// SetConcentration
	sc := new(server.SetConcentration)
	sc.ObjectID = 66432
	sc.Enabled = false
	sc.Interrupted = false
	s.Send(sc)

	// StartGame
	sg := new(server.StartGame)
	sg.Result = 4
	sg.Resolution = 1024
	s.Send(sg)

	// MapInformation
	bytes := []byte{
		1, 48,
		14, 66, 105, 99, 104, 111, 110, 80, 114, 111, 118, 105, 110, 99, 101,
		101, 0,
		135, 0,
		0, 0,
		0,
		0,
		0}
	r := new(server.MapInformation)
	codec := new(mircodec.MirCodec)
	if err := codec.Decode(bytes, r); err != nil {
		panic(err)
	}
	s.Send(r)

	// NewItemInfo
	bytes1 := []byte{146, 2, 0, 0, 13, 40, 72, 80, 41, 68, 114, 117, 103, 83, 109, 97, 108, 108, 13, 0, 0, 31, 3, 0, 0, 0, 1, 0, 0, 142, 1, 0, 0, 20, 0, 0, 0, 40, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 30, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 1, 0}
	item1 := new(server.NewItemInfo)
	codec.Decode(bytes1, item1)
	s.Send(item1)

	bytes2 := []byte{235, 4, 0, 0, 16, 84, 101, 115, 116, 83, 101, 114, 118, 101, 114, 83, 99, 114, 111, 108, 108, 21, 4, 0, 31, 3, 0, 1, 0, 1, 0, 0, 254, 6, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 23, 0, 0, 0, 0, 0, 0, 0, 0, 1, 49, 66, 97, 115, 105, 99, 32, 84, 101, 115, 116, 32, 83, 101, 114, 118, 101, 114, 32, 83, 99, 114, 111, 108, 108, 32, 119, 104, 105, 99, 104, 32, 103, 105, 118, 101, 115, 32, 105, 110, 102, 111, 114, 109, 97, 116, 105, 111, 110, 46}
	item2 := new(server.NewItemInfo)
	codec.Decode(bytes2, item2)
	s.Send(item2)

	bytes3 := []byte{221, 0, 0, 0, 11, 87, 111, 111, 100, 101, 110, 83, 119, 111, 114, 100, 1, 1, 0, 7, 3, 0, 0, 0, 4, 0, 1, 30, 0, 160, 15, 1, 0, 0, 0, 50, 0, 0, 0, 0, 0, 0, 0, 2, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0}
	item3 := new(server.NewItemInfo)
	codec.Decode(bytes3, item3)
	s.Send(item3)

	bytes4 := []byte{61, 1, 0, 0, 12, 66, 97, 115, 101, 68, 114, 101, 115, 115, 40, 77, 41, 2, 1, 0, 31, 1, 0, 1, 0, 5, 0, 1, 60, 0, 136, 19, 1, 0, 0, 0, 120, 0, 0, 0, 2, 2, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 1, 1, 0}
	item4 := new(server.NewItemInfo)
	codec.Decode(bytes4, item4)
	s.Send(item4)

	bytes5 := []byte{210, 2, 0, 0, 6, 67, 97, 110, 100, 108, 101, 12, 0, 0, 31, 3, 0, 0, 0, 1, 38, 0, 130, 0, 64, 31, 1, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 1, 0}
	item5 := new(server.NewItemInfo)
	codec.Decode(bytes5, item5)
	s.Send(item5)

	uiCodec := new(mircodec.MirUserInformationCodec)
	bytes6 := []byte{128, 3, 1, 0, 1, 0, 0, 0, 6, 99, 99, 99, 99, 99, 99, 0, 0, 255, 255, 255, 255, 1, 1, 1, 0, 28, 1, 0, 0, 96, 2, 0, 0, 1, 8, 15, 0, 17, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 1, 46, 0, 0, 0, 1, 3, 0, 0, 0, 0, 0, 0, 0, 146, 2, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 5, 0, 0, 0, 0, 0, 0, 0, 235, 4, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 4, 0, 0, 0, 0, 0, 0, 0, 210, 2, 0, 0, 54, 31, 64, 31, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 14, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 221, 0, 0, 0, 160, 15, 160, 15, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 0, 0, 0, 0, 0, 0, 0, 62, 1, 0, 0, 136, 19, 136, 19, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 40, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	item6 := new(server.UserInformation)
	uiCodec.Decode(bytes6, item6)
	s.Send(item6)
}
func (g *Game) LogOut(s cellnet.Session, msg *client.LogOut) {

}

func (g *Game) Turn(s cellnet.Session, msg *client.Turn) {

}

func (g *Game) Walk(s cellnet.Session, msg *client.Walk) {

}

func (g *Game) Run(s cellnet.Session, msg *client.Run) {

}

func (g *Game) Chat(s cellnet.Session, msg *client.Chat) {

}

func (g *Game) MoveItem(s cellnet.Session, msg *client.MoveItem) {

}

func (g *Game) StoreItem(s cellnet.Session, msg *client.StoreItem) {

}

func (g *Game) DepositRefineItem(s cellnet.Session, msg *client.DepositRefineItem) {

}

func (g *Game) RetrieveRefineItem(s cellnet.Session, msg *client.RetrieveRefineItem) {

}

func (g *Game) RefineCancel(s cellnet.Session, msg *client.RefineCancel) {

}

func (g *Game) RefineItem(s cellnet.Session, msg *client.RefineItem) {

}

func (g *Game) CheckRefine(s cellnet.Session, msg *client.CheckRefine) {

}

func (g *Game) ReplaceWedRing(s cellnet.Session, msg *client.ReplaceWedRing) {

}

func (g *Game) DepositTradeItem(s cellnet.Session, msg *client.DepositTradeItem) {

}

func (g *Game) RetrieveTradeItem(s cellnet.Session, msg *client.RetrieveTradeItem) {

}

func (g *Game) TakeBackItem(s cellnet.Session, msg *client.TakeBackItem) {

}

func (g *Game) MergeItem(s cellnet.Session, msg *client.MergeItem) {

}

func (g *Game) EquipItem(s cellnet.Session, msg *client.EquipItem) {

}

func (g *Game) RemoveItem(s cellnet.Session, msg *client.RemoveItem) {

}

func (g *Game) RemoveSlotItem(s cellnet.Session, msg *client.RemoveSlotItem) {

}

func (g *Game) SplitItem(s cellnet.Session, msg *client.SplitItem) {

}

func (g *Game) UseItem(s cellnet.Session, msg *client.UseItem) {

}

func (g *Game) DropItem(s cellnet.Session, msg *client.DropItem) {

}

func (g *Game) DropGold(s cellnet.Session, msg *client.DropGold) {

}

func (g *Game) PickUp(s cellnet.Session, msg *client.PickUp) {

}

func (g *Game) Inspect(s cellnet.Session, msg *client.Inspect) {

}

func (g *Game) ChangeAMode(s cellnet.Session, msg *client.ChangeAMode) {

}

func (g *Game) ChangePMode(s cellnet.Session, msg *client.ChangePMode) {

}

func (g *Game) ChangeTrade(s cellnet.Session, msg *client.ChangeTrade) {

}

func (g *Game) Attack(s cellnet.Session, msg *client.Attack) {

}

func (g *Game) RangeAttack(s cellnet.Session, msg *client.RangeAttack) {

}

func (g *Game) Harvest(s cellnet.Session, msg *client.Harvest) {

}

func (g *Game) CallNPC(s cellnet.Session, msg *client.CallNPC) {

}

func (g *Game) TalkMonsterNPC(s cellnet.Session, msg *client.TalkMonsterNPC) {

}

func (g *Game) BuyItem(s cellnet.Session, msg *client.BuyItem) {

}

func (g *Game) CraftItem(s cellnet.Session, msg *client.CraftItem) {

}

func (g *Game) SellItem(s cellnet.Session, msg *client.SellItem) {

}

func (g *Game) RepairItem(s cellnet.Session, msg *client.RepairItem) {

}

func (g *Game) BuyItemBack(s cellnet.Session, msg *client.BuyItemBack) {

}

func (g *Game) SRepairItem(s cellnet.Session, msg *client.SRepairItem) {

}

func (g *Game) MagicKey(s cellnet.Session, msg *client.MagicKey) {

}

func (g *Game) Magic(s cellnet.Session, msg *client.Magic) {

}

func (g *Game) SwitchGroup(s cellnet.Session, msg *client.SwitchGroup) {

}

func (g *Game) AddMember(s cellnet.Session, msg *client.AddMember) {

}

func (g *Game) DelMember(s cellnet.Session, msg *client.DellMember) {

}

func (g *Game) GroupInvite(s cellnet.Session, msg *client.GroupInvite) {

}

func (g *Game) TownRevive(s cellnet.Session, msg *client.TownRevive) {

}

func (g *Game) SpellToggle(s cellnet.Session, msg *client.SpellToggle) {

}

func (g *Game) ConsignItem(s cellnet.Session, msg *client.ConsignItem) {

}

func (g *Game) MarketSearch(s cellnet.Session, msg *client.MarketSearch) {

}

func (g *Game) MarketRefresh(s cellnet.Session, msg *client.MarketRefresh) {

}

func (g *Game) MarketPage(s cellnet.Session, msg *client.MarketPage) {

}

func (g *Game) MarketBuy(s cellnet.Session, msg *client.MarketBuy) {

}

func (g *Game) MarketGetBack(s cellnet.Session, msg *client.MarketGetBack) {

}

func (g *Game) RequestUserName(s cellnet.Session, msg *client.RequestUserName) {

}

func (g *Game) RequestChatItem(s cellnet.Session, msg *client.RequestChatItem) {

}

func (g *Game) EditGuildMember(s cellnet.Session, msg *client.EditGuildMember) {

}

func (g *Game) EditGuildNotice(s cellnet.Session, msg *client.EditGuildNotice) {

}

func (g *Game) GuildInvite(s cellnet.Session, msg *client.GuildInvite) {

}

func (g *Game) RequestGuildInfo(s cellnet.Session, msg *client.RequestGuildInfo) {

}

func (g *Game) GuildNameReturn(s cellnet.Session, msg *client.GuildNameReturn) {

}

func (g *Game) GuildStorageGoldChange(s cellnet.Session, msg *client.GuildStorageGoldChange) {

}

func (g *Game) GuildStorageItemChange(s cellnet.Session, msg *client.GuildStorageItemChange) {

}

func (g *Game) GuildWarReturn(s cellnet.Session, msg *client.GuildWarReturn) {

}

func (g *Game) MarriageRequest(s cellnet.Session, msg *client.MarriageRequest) {

}

func (g *Game) MarriageReply(s cellnet.Session, msg *client.MarriageReply) {

}

func (g *Game) ChangeMarriage(s cellnet.Session, msg *client.ChangeMarriage) {

}

func (g *Game) DivorceRequest(s cellnet.Session, msg *client.DivorceRequest) {

}

func (g *Game) DivorceReply(s cellnet.Session, msg *client.DivorceReply) {

}

func (g *Game) AddMentor(s cellnet.Session, msg *client.AddMentor) {

}

func (g *Game) MentorReply(s cellnet.Session, msg *client.MentorReply) {

}

func (g *Game) AllowMentor(s cellnet.Session, msg *client.AllowMentor) {

}

func (g *Game) CancelMentor(s cellnet.Session, msg *client.CancelMentor) {

}

func (g *Game) TradeRequest(s cellnet.Session, msg *client.TradeRequest) {

}

func (g *Game) TradeGold(s cellnet.Session, msg *client.TradeGold) {

}

func (g *Game) TradeReply(s cellnet.Session, msg *client.TradeReply) {

}

func (g *Game) TradeConfirm(s cellnet.Session, msg *client.TradeConfirm) {

}

func (g *Game) TradeCancel(s cellnet.Session, msg *client.TradeCancel) {

}

func (g *Game) EquipSlotItem(s cellnet.Session, msg *client.EquipSlotItem) {

}

func (g *Game) FishingCast(s cellnet.Session, msg *client.FishingCast) {

}

func (g *Game) FishingChangeAutocast(s cellnet.Session, msg *client.FishingChangeAutocast) {

}

func (g *Game) AcceptQuest(s cellnet.Session, msg *client.AcceptQuest) {

}

func (g *Game) FinishQuest(s cellnet.Session, msg *client.FinishQuest) {

}

func (g *Game) AbandonQuest(s cellnet.Session, msg *client.AbandonQuest) {

}

func (g *Game) ShareQuest(s cellnet.Session, msg *client.ShareQuest) {

}

func (g *Game) AcceptReincarnation(s cellnet.Session, msg *client.AcceptReincarnation) {

}

func (g *Game) CancelReincarnation(s cellnet.Session, msg *client.CancelReincarnation) {

}

func (g *Game) CombineItem(s cellnet.Session, msg *client.CombineItem) {

}

func (g *Game) SetConcentration(s cellnet.Session, msg *client.SetConcentration) {

}

func (g *Game) AwakeningNeedMaterials(s cellnet.Session, msg *client.AwakeningNeedMaterials) {

}

func (g *Game) AwakeningLockedItem(s cellnet.Session, msg *client.AwakeningLockedItem) {

}

func (g *Game) Awakening(s cellnet.Session, msg *client.Awakening) {

}

func (g *Game) DisassembleItem(s cellnet.Session, msg *client.DisassembleItem) {

}

func (g *Game) DowngradeAwakening(s cellnet.Session, msg *client.DowngradeAwakening) {

}

func (g *Game) ResetAddedItem(s cellnet.Session, msg *client.ResetAddedItem) {

}

func (g *Game) SendMail(s cellnet.Session, msg *client.SendMail) {

}

func (g *Game) ReadMail(s cellnet.Session, msg *client.ReadMail) {

}

func (g *Game) CollectParcel(s cellnet.Session, msg *client.CollectParcel) {

}

func (g *Game) DeleteMail(s cellnet.Session, msg *client.DeleteMail) {

}

func (g *Game) LockMail(s cellnet.Session, msg *client.LockMail) {

}

func (g *Game) MailLockedItem(s cellnet.Session, msg *client.MailLockedItem) {

}

func (g *Game) MailCost(s cellnet.Session, msg *client.MailCost) {

}

func (g *Game) UpdateIntelligentCreature(s cellnet.Session, msg *client.UpdateIntelligentCreature) {

}

func (g *Game) IntelligentCreaturePickup(s cellnet.Session, msg *client.IntelligentCreaturePickup) {

}

func (g *Game) AddFriend(s cellnet.Session, msg *client.AddFriend) {

}

func (g *Game) RemoveFriend(s cellnet.Session, msg *client.RemoveFriend) {

}

func (g *Game) RefreshFriends(s cellnet.Session, msg *client.RefreshFriends) {

}

func (g *Game) AddMemo(s cellnet.Session, msg *client.AddMemo) {

}

func (g *Game) GuildBuffUpdate(s cellnet.Session, msg *client.GuildBuffUpdate) {

}

func (g *Game) GameshopBuy(s cellnet.Session, msg *client.GameshopBuy) {

}

func (g *Game) NPCConfirmInput(s cellnet.Session, msg *client.NPCConfirmInput) {

}

func (g *Game) ReportIssue(s cellnet.Session, msg *client.ReportIssue) {

}

func (g *Game) GetRanking(s cellnet.Session, msg *client.GetRanking) {

}

func (g *Game) Opendoor(s cellnet.Session, msg *client.Opendoor) {

}

func (g *Game) GetRentedItems(s cellnet.Session, msg *client.GetRentedItems) {

}

func (g *Game) ItemRentalRequest(s cellnet.Session, msg *client.ItemRentalRequest) {

}

func (g *Game) ItemRentalFee(s cellnet.Session, msg *client.ItemRentalFee) {

}

func (g *Game) ItemRentalPeriod(s cellnet.Session, msg *client.ItemRentalPeriod) {

}

func (g *Game) DepositRentalItem(s cellnet.Session, msg *client.DepositRentalItem) {

}

func (g *Game) RetrieveRentalItem(s cellnet.Session, msg *client.RetrieveRentalItem) {

}

func (g *Game) CancelItemRental(s cellnet.Session, msg *client.CancelItemRental) {

}

func (g *Game) ItemRentalLockFee(s cellnet.Session, msg *client.ItemRentalLockFee) {

}

func (g *Game) ItemRentalLockItem(s cellnet.Session, msg *client.ItemRentalLockItem) {

}

func (g *Game) ConfirmItemRental(s cellnet.Session, msg *client.ConfirmItemRental) {

}
