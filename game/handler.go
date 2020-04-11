package game

import (
	"time"

	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/game/cm"
	_ "github.com/yenkeia/mirgo/game/mirtcp"
	"github.com/yenkeia/mirgo/game/proto/client"
	"github.com/yenkeia/mirgo/game/proto/server"
)

type Game struct {
}

func (g *Game) HandleEvent(ev cellnet.Event) {
	// g.Pool.Submit(NewTask(_HandleEvent, g, ev))

	s := ev.Session()

	switch msg := ev.Message().(type) {
	case *cellnet.SessionAccepted: // 有新的连接
		g.SessionAccepted(s, msg)
	case *cellnet.SessionClosed: // 有连接断开
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

	default:
		p, ok := g.GetPlayer(s, GAME)
		if !ok {
			return
		}

		_HandleEvent(p, g, ev, s)
	}
}

func _HandleEvent(p *Player, g *Game, ev cellnet.Event, s cellnet.Session) {

	switch msg := ev.Message().(type) {
	case *client.Turn:
		g.Turn(p, msg)
	case *client.Walk:
		g.Walk(p, msg)
	case *client.Run:
		g.Run(p, msg)
	case *client.Chat:
		g.Chat(p, msg)
	case *client.MoveItem:
		g.MoveItem(p, msg)
	case *client.StoreItem:
		g.StoreItem(p, msg)
	case *client.DepositRefineItem:
		g.DepositRefineItem(p, msg)
	case *client.RetrieveRefineItem:
		g.RetrieveRefineItem(p, msg)
	case *client.RefineCancel:
		g.RefineCancel(p, msg)
	case *client.RefineItem:
		g.RefineItem(p, msg)
	case *client.CheckRefine:
		g.CheckRefine(p, msg)
	case *client.ReplaceWedRing:
		g.ReplaceWedRing(p, msg)
	case *client.DepositTradeItem:
		g.DepositTradeItem(p, msg)
	case *client.RetrieveTradeItem:
		g.RetrieveTradeItem(p, msg)
	case *client.TakeBackItem:
		g.TakeBackItem(p, msg)
	case *client.MergeItem:
		g.MergeItem(p, msg)
	case *client.EquipItem:
		g.EquipItem(p, msg)
	case *client.RemoveItem:
		g.RemoveItem(p, msg)
	case *client.RemoveSlotItem:
		g.RemoveSlotItem(p, msg)
	case *client.SplitItem:
		g.SplitItem(p, msg)
	case *client.UseItem:
		g.UseItem(p, msg)
	case *client.DropItem:
		g.DropItem(p, msg)
	case *client.DropGold:
		g.DropGold(p, msg)
	case *client.PickUp:
		g.PickUp(p, msg)
	case *client.Inspect:
		g.Inspect(p, msg)
	case *client.ChangeAMode:
		g.ChangeAMode(p, msg)
	case *client.ChangePMode:
		g.ChangePMode(p, msg)
	case *client.ChangeTrade:
		g.ChangeTrade(p, msg)
	case *client.Attack:
		g.Attack(p, msg)
	case *client.RangeAttack:
		g.RangeAttack(p, msg)
	case *client.Harvest:
		g.Harvest(p, msg)
	case *client.CallNPC:
		g.CallNPC(p, msg)
	case *client.TalkMonsterNPC:
		g.TalkMonsterNPC(p, msg)
	case *client.BuyItem:
		g.BuyItem(p, msg)
	case *client.CraftItem:
		g.CraftItem(p, msg)
	case *client.SellItem:
		g.SellItem(p, msg)
	case *client.RepairItem:
		g.RepairItem(p, msg)
	case *client.BuyItemBack:
		g.BuyItemBack(p, msg)
	case *client.SRepairItem:
		g.SRepairItem(p, msg)
	case *client.MagicKey:
		g.MagicKey(p, msg)
	case *client.Magic:
		g.Magic(p, msg)
	case *client.SwitchGroup:
		g.SwitchGroup(p, msg)
	case *client.AddMember:
		g.AddMember(p, msg)
	case *client.DelMember:
		g.DelMember(p, msg)
	case *client.GroupInvite:
		g.GroupInvite(p, msg)
	case *client.TownRevive:
		g.TownRevive(p, msg)
	case *client.SpellToggle:
		g.SpellToggle(p, msg)
	case *client.ConsignItem:
		g.ConsignItem(p, msg)
	case *client.MarketSearch:
		g.MarketSearch(p, msg)
	case *client.MarketRefresh:
		g.MarketRefresh(p, msg)
	case *client.MarketPage:
		g.MarketPage(p, msg)
	case *client.MarketBuy:
		g.MarketBuy(p, msg)
	case *client.MarketGetBack:
		g.MarketGetBack(p, msg)
	case *client.RequestUserName:
		g.RequestUserName(p, msg)
	case *client.RequestChatItem:
		g.RequestChatItem(p, msg)
	case *client.EditGuildMember:
		g.EditGuildMember(p, msg)
	case *client.EditGuildNotice:
		g.EditGuildNotice(p, msg)
	case *client.GuildInvite:
		g.GuildInvite(p, msg)
	case *client.RequestGuildInfo:
		g.RequestGuildInfo(p, msg)
	case *client.GuildNameReturn:
		g.GuildNameReturn(p, msg)
	case *client.GuildStorageGoldChange:
		g.GuildStorageGoldChange(p, msg)
	case *client.GuildStorageItemChange:
		g.GuildStorageItemChange(p, msg)
	case *client.GuildWarReturn:
		g.GuildWarReturn(p, msg)
	case *client.MarriageRequest:
		g.MarriageRequest(p, msg)
	case *client.MarriageReply:
		g.MarriageReply(p, msg)
	case *client.ChangeMarriage:
		g.ChangeMarriage(p, msg)
	case *client.DivorceRequest:
		g.DivorceRequest(p, msg)
	case *client.DivorceReply:
		g.DivorceReply(p, msg)
	case *client.AddMentor:
		g.AddMentor(p, msg)
	case *client.MentorReply:
		g.MentorReply(p, msg)
	case *client.AllowMentor:
		g.AllowMentor(p, msg)
	case *client.CancelMentor:
		g.CancelMentor(p, msg)
	case *client.TradeRequest:
		g.TradeRequest(p, msg)
	case *client.TradeGold:
		g.TradeGold(p, msg)
	case *client.TradeReply:
		g.TradeReply(p, msg)
	case *client.TradeConfirm:
		g.TradeConfirm(p, msg)
	case *client.TradeCancel:
		g.TradeCancel(p, msg)
	case *client.EquipSlotItem:
		g.EquipSlotItem(p, msg)
	case *client.FishingCast:
		g.FishingCast(p, msg)
	case *client.FishingChangeAutocast:
		g.FishingChangeAutocast(p, msg)
	case *client.AcceptQuest:
		g.AcceptQuest(p, msg)
	case *client.FinishQuest:
		g.FinishQuest(p, msg)
	case *client.AbandonQuest:
		g.AbandonQuest(p, msg)
	case *client.ShareQuest:
		g.ShareQuest(p, msg)
	case *client.AcceptReincarnation:
		g.AcceptReincarnation(p, msg)
	case *client.CancelReincarnation:
		g.CancelReincarnation(p, msg)
	case *client.CombineItem:
		g.CombineItem(p, msg)
	case *client.SetConcentration:
		g.SetConcentration(p, msg)
	case *client.AwakeningNeedMaterials:
		g.AwakeningNeedMaterials(p, msg)
	case *client.AwakeningLockedItem:
		g.AwakeningLockedItem(p, msg)
	case *client.Awakening:
		g.Awakening(p, msg)
	case *client.DisassembleItem:
		g.DisassembleItem(p, msg)
	case *client.DowngradeAwakening:
		g.DowngradeAwakening(p, msg)
	case *client.ResetAddedItem:
		g.ResetAddedItem(p, msg)
	case *client.SendMail:
		g.SendMail(p, msg)
	case *client.ReadMail:
		g.ReadMail(p, msg)
	case *client.CollectParcel:
		g.CollectParcel(p, msg)
	case *client.DeleteMail:
		g.DeleteMail(p, msg)
	case *client.LockMail:
		g.LockMail(p, msg)
	case *client.MailLockedItem:
		g.MailLockedItem(p, msg)
	case *client.MailCost:
		g.MailCost(p, msg)
	case *client.UpdateIntelligentCreature: //IntelligentCreature
		g.UpdateIntelligentCreature(p, msg)
	case *client.IntelligentCreaturePickup: //IntelligentCreature
		g.IntelligentCreaturePickup(p, msg)
	case *client.AddFriend:
		g.AddFriend(p, msg)
	case *client.RemoveFriend:
		g.RemoveFriend(p, msg)
	case *client.RefreshFriends:
		g.RefreshFriends(p, msg)
	case *client.AddMemo:
		g.AddMemo(p, msg)
	case *client.GuildBuffUpdate:
		g.GuildBuffUpdate(p, msg)
	case *client.GameshopBuy:
		g.GameshopBuy(p, msg)
	case *client.NPCConfirmInput:
		g.NPCConfirmInput(p, msg)
	case *client.ReportIssue:
		g.ReportIssue(p, msg)
	case *client.GetRanking:
		g.GetRanking(p, msg)
	case *client.Opendoor:
		g.Opendoor(p, msg)
	case *client.GetRentedItems:
		g.GetRentedItems(p, msg)
	case *client.ItemRentalRequest:
		g.ItemRentalRequest(p, msg)
	case *client.ItemRentalFee:
		g.ItemRentalFee(p, msg)
	case *client.ItemRentalPeriod:
		g.ItemRentalPeriod(p, msg)
	case *client.DepositRentalItem:
		g.DepositRentalItem(p, msg)
	case *client.RetrieveRentalItem:
		g.RetrieveRentalItem(p, msg)
	case *client.CancelItemRental:
		g.CancelItemRental(p, msg)
	case *client.ItemRentalLockFee:
		g.ItemRentalLockFee(p, msg)
	case *client.ItemRentalLockItem:
		g.ItemRentalLockItem(p, msg)
	case *client.ConfirmItemRental:
		g.ConfirmItemRental(p, msg)
	default:
		log.Debugln("default:", msg)
		//MessageQueue.Enqueue(string.Format("Invalid packet received. Index : {0}", p.Index));
	}
}

// SessionAccepted ...
func (g *Game) SessionAccepted(s cellnet.Session, msg *cellnet.SessionAccepted) {
	s.Send(&server.Connected{})
}

// SessionClosed ...
func (g *Game) SessionClosed(s cellnet.Session, msg *cellnet.SessionClosed) {
	pm := env.SessionIDPlayerMap
	v, ok := pm.Load(s.ID())
	if !ok {
		return
	}
	p := v.(*Player)
	log.Debugln("SessionClosed")
	if p.GameStage == GAME {
		p.GameStage = DISCONNECTED
		p.StopGame(StopGameUserClosedGame)
		env.Players.Remove(p)
		p.Map.DeleteObject(p)
		p.SaveData()
		log.Debugf("删除玩家: %s\n", p.Name)
	}
	pm.Delete(s.ID())
}

// ClientVersion ...
func (g *Game) ClientVersion(s cellnet.Session, msg *client.ClientVersion) {
	clientVersion := server.ClientVersion{Result: 1}
	p := new(Player)
	p.GameStage = LOGIN
	p.Session = &s
	env.SessionIDPlayerMap.Store(s.ID(), p)
	s.Send(&clientVersion)
}

func (g *Game) GetPlayer(s cellnet.Session, gameStage int) (p *Player, ok bool) {
	v, ok := env.SessionIDPlayerMap.Load(s.ID())
	if !ok {
		return nil, false
	}
	p = v.(*Player)
	if p.GameStage != gameStage {
		return nil, false
	}
	return p, true
}

// KeepAlive ...
func (g *Game) KeepAlive(s cellnet.Session, msg *client.KeepAlive) {
	s.Send(server.KeepAlive{Time: 0})
}

// NewAccount 保存新账号
func (g *Game) NewAccount(s cellnet.Session, msg *client.NewAccount) {
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
	_, ok := g.GetPlayer(s, LOGIN)
	if !ok {
		return
	}

	res := uint8(0)
	ac := new(cm.Account)
	adb.Table("account").Where("username = ?", msg.AccountID).Find(ac)
	if ac.ID == 0 && ac.Username == "" {
		ac.Username = msg.AccountID
		ac.Password = msg.Password
		adb.Table("account").Create(&ac)
		res = 8
	}
	s.Send(&server.NewAccount{Result: res})
}

// ChangePassword 改密码
func (g *Game) ChangePassword(s cellnet.Session, msg *client.ChangePassword) {
	/*
	 * 0: Disabled
	 * 1: Bad AccountID
	 * 2: Bad Current Password
	 * 3: Bad New Password
	 * 4: Account Not Exist
	 * 5: Wrong Password
	 * 6: Success
	 */
	_, ok := g.GetPlayer(s, LOGIN)
	if !ok {
		return
	}

	res := uint8(5)
	ac := new(cm.Account)
	adb.Table("account").Where("username = ? AND password = ?", msg.AccountID, msg.CurrentPassword).Find(ac)
	if ac.ID != 0 {
		ac.Password = msg.NewPassword
		adb.Table("account").Model(ac).Updates(cm.Account{Password: msg.NewPassword})
		res = 6
	}
	s.Send(&server.ChangePassword{Result: res})
}

func (g *Game) getAccountCharacters(AccountID int) []server.SelectInfo {
	ac := make([]cm.AccountCharacter, 3)
	adb.Table("account_character").Where("account_id = ?", AccountID).Limit(3).Find(&ac)
	ids := make([]int, 3)
	for _, c := range ac {
		ids = append(ids, c.ID)
	}
	cs := make([]cm.Character, 3)
	adb.Table("character").Where("id in (?)", ids).Find(&cs)
	si := make([]server.SelectInfo, len(cs))
	for i, c := range cs {
		s := new(server.SelectInfo)
		s.Index = uint32(c.ID)
		s.Name = c.Name
		s.Level = c.Level
		s.Class = c.Class
		s.Gender = c.Gender
		s.LastAccess = 0
		si[i] = *s
	}
	return si
}

// Login 登陆
func (g *Game) Login(s cellnet.Session, msg *client.Login) {
	p, ok := g.GetPlayer(s, LOGIN)
	if !ok {
		return
	}
	a := new(cm.Account)
	adb.Table("account").Where("username = ? AND password = ?", msg.AccountID, msg.Password).Find(a)
	if a.ID == 0 {
		s.Send(ServerMessage{}.Login(4))
		return
	}
	p.AccountID = a.ID
	p.GameStage = SELECT
	res := new(server.LoginSuccess)
	res.Characters = g.getAccountCharacters(p.AccountID)
	s.Send(res)
}

// NewCharacter 创建角色
func (g *Game) NewCharacter(s cellnet.Session, msg *client.NewCharacter) {
	p, ok := g.GetPlayer(s, SELECT)
	if !ok {
		return
	}
	acs := make([]cm.AccountCharacter, 3)
	adb.Table("account_character").Where("account_id = ?", p.AccountID).Limit(3).Find(&acs)
	if len(acs) >= 3 {
		s.Send(ServerMessage{}.NewCharacter(4))
		return
	}
	s.Send(ServerMessage{}.NewCharacterSuccess(g, p.AccountID, msg.Name, msg.Class, msg.Gender))
}

// DeleteCharacter 删除角色
func (g *Game) DeleteCharacter(s cellnet.Session, msg *client.DeleteCharacter) {
	_, ok := g.GetPlayer(s, SELECT)
	if !ok {
		return
	}

	c := new(cm.Character)
	adb.Table("character").Where("id = ?", msg.CharacterIndex).Find(c)
	if c.ID == 0 {
		res := new(server.DeleteCharacter)
		res.Result = 4
		s.Send(res)
		return
	}
	adb.Table("character").Delete(c)
	adb.Table("account_character").Where("character_id = ?", c.ID).Delete(cm.Character{})
	res := new(server.DeleteCharacterSuccess)
	res.CharacterIndex = msg.CharacterIndex
	s.Send(res)
}

func updatePlayerInfo(g *Game, p *Player, c *cm.Character) {
	p.GameStage = GAME
	p.ID = uint32(c.ID)
	p.Name = c.Name
	p.NameColor = cm.Color{R: 255, G: 255, B: 255}
	p.Direction = c.Direction
	p.CurrentLocation = cm.NewPoint(int(c.CurrentLocationX), int(c.CurrentLocationY))
	p.BindLocation = cm.NewPoint(c.BindLocationX, c.BindLocationY)
	p.BindMapIndex = c.BindMapID

	magics := make([]*cm.UserMagic, 0)
	adb.Table("user_magic").Where("character_id = ?", c.ID).Find(&magics)
	for _, v := range magics {
		v.Info = data.GetMagicInfoByID(v.MagicID)
	}

	p.Inventory = BagLoadFromDB(p, cm.UserItemTypeInventory, 46)
	p.Equipment = BagLoadFromDB(p, cm.UserItemTypeEquipment, 14)
	p.QuestInventory = BagLoadFromDB(p, cm.UserItemTypeQuestInventory, 40)
	p.Storage = BagLoadFromDB(p, cm.UserItemTypeStorage, 80)
	p.Trade = NewBag(p, cm.UserItemTypeTrade, 10)

	healNextTime := time.Now().Add(10 * time.Second)
	p.Dead = false
	p.HP = c.HP
	p.MP = c.MP
	p.Level = c.Level
	p.Experience = c.Experience
	p.Gold = c.Gold
	p.GuildName = ""     // TODO
	p.GuildRankName = "" // TODO
	p.Class = c.Class
	p.Gender = c.Gender
	p.Hair = c.Hair
	p.SendItemInfo = make([]*cm.ItemInfo, 0)
	p.MaxExperience = int64(data.ExpList[p.Level-1])
	p.Magics = magics
	p.ActionList = NewActionList()
	p.PoisonList = NewPoisonList()
	p.BuffList = NewBuffList()
	p.Health = Health{
		HPPotNextTime: new(time.Time),
		HPPotDuration: 1 * time.Second,
		MPPotNextTime: new(time.Time),
		MPPotDuration: 1 * time.Second,
		HealNextTime:  &healNextTime,
		HealDuration:  10 * time.Second,
	}
	p.Pets = make([]*Monster, 0)
	p.PKPoints = 0
	p.AMode = c.AttackMode
	p.PMode = c.PetMode
	p.CallingNPC = nil
	p.StruckTime = time.Now()
	p.DamageRate = 1.0
	p.ArmourRate = 1.0
	p.AllowGroup = c.AllowGroup
}

// StartGame 开始游戏
func (g *Game) StartGame(s cellnet.Session, msg *client.StartGame) {
	p, ok := g.GetPlayer(s, SELECT)
	if !ok {
		return
	}

	c := new(cm.Character)
	adb.Table("character").Where("id = ?", msg.CharacterIndex).Find(c)
	if c.ID == 0 {
		return
	}
	ac := new(cm.AccountCharacter)
	adb.Table("account_character").Where("account_id = ? and character_id = ?", p.AccountID, c.ID).Find(&ac)
	if ac.ID == 0 {
		s.Send(ServerMessage{}.StartGame(2, 1024))
		return
	}
	s.Send(ServerMessage{}.SetConcentration(p))
	s.Send(ServerMessage{}.StartGame(4, 1024))

	updatePlayerInfo(g, p, c)

	if p.Level == 0 {
		p.Level = 1
		adb.SyncLevel(p)
		for _, v := range data.StartItems {
			p.GainItem(env.NewUserItem(v))
		}
	}

	log.Debugf("player login, AccountID(%d) Name(%s)\n", p.AccountID, p.Name)
	p.Map = env.GetMap(int(c.CurrentMapID))
	env.Players.Add(p)
	p.Map.AddObject(p)
	p.StartGame()
}

func (g *Game) LogOut(s cellnet.Session, msg *client.LogOut) {
	p, ok := g.GetPlayer(s, GAME)
	if !ok {
		return
	}
	p.GameStage = SELECT
	p.StopGame(StopGameUserReturnedToSelectChar)
	env.Players.Remove(p)
	p.Map.DeleteObject(p)
	s.Send(&server.LogOutSuccess{Characters: g.getAccountCharacters(p.AccountID)})
	p.SaveData()
}

func (g *Game) Turn(p *Player, msg *client.Turn) {
	p.Turn(msg.Direction)
}

func (g *Game) Walk(p *Player, msg *client.Walk) {
	p.Walk(msg.Direction)
}

func (g *Game) Run(p *Player, msg *client.Run) {
	p.Run(msg.Direction)
}

func (g *Game) Chat(p *Player, msg *client.Chat) {
	p.Chat(msg.Message)
}

func (g *Game) MoveItem(p *Player, msg *client.MoveItem) {
	p.MoveItem(msg.Grid, msg.From, msg.To)
}

func (g *Game) StoreItem(p *Player, msg *client.StoreItem) {
	p.StoreItem(msg.From, msg.To)
}

func (g *Game) DepositRefineItem(p *Player, msg *client.DepositRefineItem) {
	p.DepositRefineItem(msg.From, msg.To)
}

func (g *Game) RetrieveRefineItem(p *Player, msg *client.RetrieveRefineItem) {
	p.RetrieveRefineItem(msg.From, msg.To)
}

func (g *Game) RefineCancel(p *Player, msg *client.RefineCancel) {
	p.RefineCancel()
}

func (g *Game) RefineItem(p *Player, msg *client.RefineItem) {
	p.RefineItem(msg.UniqueID)
}

func (g *Game) CheckRefine(p *Player, msg *client.CheckRefine) {
	p.CheckRefine(msg.UniqueID)
}

func (g *Game) ReplaceWedRing(p *Player, msg *client.ReplaceWedRing) {
	p.ReplaceWeddingRing(msg.UniqueID)
}

func (g *Game) DepositTradeItem(p *Player, msg *client.DepositTradeItem) {
	p.DepositTradeItem(msg.From, msg.To)
}

func (g *Game) RetrieveTradeItem(p *Player, msg *client.RetrieveTradeItem) {
	p.RetrieveTradeItem(msg.From, msg.To)
}

func (g *Game) TakeBackItem(p *Player, msg *client.TakeBackItem) {
	p.TakeBackItem(msg.From, msg.To)
}

func (g *Game) MergeItem(p *Player, msg *client.MergeItem) {
	p.MergeItem(msg.GridFrom, msg.GridTo, msg.IDFrom, msg.IDTo)
}

func (g *Game) EquipItem(p *Player, msg *client.EquipItem) {
	p.EquipItem(msg.Grid, msg.UniqueID, msg.To)
}

func (g *Game) RemoveItem(p *Player, msg *client.RemoveItem) {
	p.RemoveItem(msg.Grid, msg.UniqueID, msg.To)
}

func (g *Game) RemoveSlotItem(p *Player, msg *client.RemoveSlotItem) {
	p.RemoveSlotItem(msg.Grid, msg.UniqueID, msg.To, msg.GridTo)
}

func (g *Game) SplitItem(p *Player, msg *client.SplitItem) {
	p.SplitItem(msg.Grid, msg.UniqueID, msg.Count)
}

func (g *Game) UseItem(p *Player, msg *client.UseItem) {
	p.UseItem(msg.UniqueID)
}

func (g *Game) DropItem(p *Player, msg *client.DropItem) {
	p.DropItem(msg.UniqueID, msg.Count)
}

func (g *Game) DropGold(p *Player, msg *client.DropGold) {
	p.DropGold(uint64(msg.Amount))
}

func (g *Game) PickUp(p *Player, msg *client.PickUp) {
	p.PickUp()
}

func (g *Game) Inspect(p *Player, msg *client.Inspect) {
	p.Inspect(msg.ObjectID)
}

func (g *Game) ChangeAMode(p *Player, msg *client.ChangeAMode) {
	p.ChangeAMode(msg.Mode)
}

func (g *Game) ChangePMode(p *Player, msg *client.ChangePMode) {
	p.ChangePMode(msg.Mode)
}

func (g *Game) ChangeTrade(p *Player, msg *client.ChangeTrade) {
	p.ChangeTrade(msg.AllowTrade)
}

func (g *Game) Attack(p *Player, msg *client.Attack) {
	p.Attack(msg.Direction, msg.Spell)
}

func (g *Game) RangeAttack(p *Player, msg *client.RangeAttack) {
	p.RangeAttack(msg.Direction, msg.TargetLocation, msg.TargetID)
}

func (g *Game) Harvest(p *Player, msg *client.Harvest) {
	p.Harvest(msg.Direction)
}

func (g *Game) CallNPC(p *Player, msg *client.CallNPC) {
	p.CallNPC(msg.ObjectID, msg.Key)
}

func (g *Game) TalkMonsterNPC(p *Player, msg *client.TalkMonsterNPC) {
	p.TalkMonsterNPC(msg.ObjectID)
}

func (g *Game) BuyItem(p *Player, msg *client.BuyItem) {
	p.BuyItem(msg.ItemIndex, msg.Count, msg.Type)
}

// TODO
func (g *Game) CraftItem(p *Player, msg *client.CraftItem) {
	p.CraftItem(msg.UniqueID, msg.Count, msg.Slots)
}

func (g *Game) SellItem(p *Player, msg *client.SellItem) {
	p.SellItem(msg.UniqueID, msg.Count)
}

func (g *Game) RepairItem(p *Player, msg *client.RepairItem) {
	p.RepairItem(msg.UniqueID, false)
}

func (g *Game) BuyItemBack(p *Player, msg *client.BuyItemBack) {
	p.BuyItemBack(msg.UniqueID, msg.Count)
}

func (g *Game) SRepairItem(p *Player, msg *client.SRepairItem) {
	p.RepairItem(msg.UniqueID, true)
}

func (g *Game) MagicKey(p *Player, msg *client.MagicKey) {
	p.MagicKey(msg.Spell, msg.Key)
}

func (g *Game) Magic(p *Player, msg *client.Magic) {
	p.Magic(msg.Spell, msg.Direction, msg.TargetID, msg.Location)
}

func (g *Game) SwitchGroup(p *Player, msg *client.SwitchGroup) {
	p.SwitchGroup(msg.AllowGroup)
}

func (g *Game) AddMember(p *Player, msg *client.AddMember) {
	p.AddMember(msg.Name)
}

func (g *Game) DelMember(p *Player, msg *client.DelMember) {
	p.DelMember(msg.Name)
}

func (g *Game) GroupInvite(p *Player, msg *client.GroupInvite) {
	p.GroupInvite(msg.AcceptInvite)
}

func (g *Game) TownRevive(p *Player, msg *client.TownRevive) {
	p.TownRevive()
}

func (g *Game) SpellToggle(p *Player, msg *client.SpellToggle) {
	p.SpellToggle(msg.Spell, msg.CanUse)
}

func (g *Game) ConsignItem(p *Player, msg *client.ConsignItem) {
	p.ConsignItem(msg.UniqueID, msg.Price)
}

func (g *Game) MarketSearch(p *Player, msg *client.MarketSearch) {
	p.MarketSearch(msg.Match)
}

func (g *Game) MarketRefresh(p *Player, msg *client.MarketRefresh) {
	p.MarketRefresh()
}

func (g *Game) MarketPage(p *Player, msg *client.MarketPage) {
	p.MarketPage(msg.Page)
}

func (g *Game) MarketBuy(p *Player, msg *client.MarketBuy) {
	p.MarketBuy(msg.AuctionID)
}

func (g *Game) MarketGetBack(p *Player, msg *client.MarketGetBack) {
	p.MarketGetBack(msg.AuctionID)
}

func (g *Game) RequestUserName(p *Player, msg *client.RequestUserName) {
	p.RequestUserName(msg.UserID)
}

func (g *Game) RequestChatItem(p *Player, msg *client.RequestChatItem) {
	p.RequestChatItem(msg.ChatItemID)
}

func (g *Game) EditGuildMember(p *Player, msg *client.EditGuildMember) {
	p.EditGuildMember(msg.Name, msg.RankName, msg.RankIndex, msg.ChangeType)
}

func (g *Game) EditGuildNotice(p *Player, msg *client.EditGuildNotice) {
	p.EditGuildNotice(msg.Notice)
}

func (g *Game) GuildInvite(p *Player, msg *client.GuildInvite) {
	p.GuildInvite(msg.AcceptInvite)
}

func (g *Game) RequestGuildInfo(p *Player, msg *client.RequestGuildInfo) {
	p.RequestGuildInfo(msg.Type)
}

func (g *Game) GuildNameReturn(p *Player, msg *client.GuildNameReturn) {
	p.GuildNameReturn(msg.Name)
}

func (g *Game) GuildStorageGoldChange(p *Player, msg *client.GuildStorageGoldChange) {
	p.GuildStorageGoldChange(msg.Type, msg.Amount)
}

func (g *Game) GuildStorageItemChange(p *Player, msg *client.GuildStorageItemChange) {
	p.GuildStorageItemChange(msg.Type, msg.From, msg.To)
}

func (g *Game) GuildWarReturn(p *Player, msg *client.GuildWarReturn) {
	p.GuildWarReturn(msg.Name)
}

func (g *Game) MarriageRequest(p *Player, msg *client.MarriageRequest) {
	p.MarriageRequest()
}

func (g *Game) MarriageReply(p *Player, msg *client.MarriageReply) {
	p.MarriageReply(msg.AcceptInvite)
}

func (g *Game) ChangeMarriage(p *Player, msg *client.ChangeMarriage) {
	p.ChangeMarriage()
}

func (g *Game) DivorceRequest(p *Player, msg *client.DivorceRequest) {
	p.DivorceRequest()
}

func (g *Game) DivorceReply(p *Player, msg *client.DivorceReply) {
	p.DivorceReply(msg.AcceptInvite)
}

func (g *Game) AddMentor(p *Player, msg *client.AddMentor) {
	p.AddMentor(msg.Name)
}

func (g *Game) MentorReply(p *Player, msg *client.MentorReply) {
	p.MentorReply(msg.AcceptInvite)
}

func (g *Game) AllowMentor(p *Player, msg *client.AllowMentor) {
	p.AllowMentor()
}

func (g *Game) CancelMentor(p *Player, msg *client.CancelMentor) {
	p.CancelMentor()
}

func (g *Game) TradeRequest(p *Player, msg *client.TradeRequest) {
	p.TradeRequest()
}

func (g *Game) TradeGold(p *Player, msg *client.TradeGold) {
	p.TradeGold(msg.Amount)
}

func (g *Game) TradeReply(p *Player, msg *client.TradeReply) {
	p.TradeReply(msg.AcceptInvite)
}

func (g *Game) TradeConfirm(p *Player, msg *client.TradeConfirm) {
	p.TradeConfirm(msg.Locked)
}

func (g *Game) TradeCancel(p *Player, msg *client.TradeCancel) {
	p.TradeCancel()
}

func (g *Game) EquipSlotItem(p *Player, msg *client.EquipSlotItem) {
	p.EquipSlotItem(msg.Grid, msg.UniqueID, msg.To, msg.GridTo)
}

func (g *Game) FishingCast(p *Player, msg *client.FishingCast) {
	p.FishingCast(msg.CastOut)
}

func (g *Game) FishingChangeAutocast(p *Player, msg *client.FishingChangeAutocast) {
	p.FishingChangeAutocast(msg.AutoCast)
}

func (g *Game) AcceptQuest(p *Player, msg *client.AcceptQuest) {
	p.AcceptQuest(msg.NPCIndex, msg.QuestIndex)
}

func (g *Game) FinishQuest(p *Player, msg *client.FinishQuest) {
	p.FinishQuest(msg.QuestIndex, msg.SelectedItemIndex)
}

func (g *Game) AbandonQuest(p *Player, msg *client.AbandonQuest) {
	p.AbandonQuest(msg.QuestIndex)
}

func (g *Game) ShareQuest(p *Player, msg *client.ShareQuest) {
	p.ShareQuest(msg.QuestIndex)
}

func (g *Game) AcceptReincarnation(p *Player, msg *client.AcceptReincarnation) {
	p.AcceptReincarnation()
}

func (g *Game) CancelReincarnation(p *Player, msg *client.CancelReincarnation) {
	p.CancelReincarnation()
}

func (g *Game) CombineItem(p *Player, msg *client.CombineItem) {
	p.CombineItem(msg.IDFrom, msg.IDTo)
}

func (g *Game) SetConcentration(p *Player, msg *client.SetConcentration) {
	p.SetConcentration(msg.ObjectID, msg.Enabled, msg.Interrupted)
}

func (g *Game) AwakeningNeedMaterials(p *Player, msg *client.AwakeningNeedMaterials) {

}

func (g *Game) AwakeningLockedItem(p *Player, msg *client.AwakeningLockedItem) {

}

func (g *Game) Awakening(p *Player, msg *client.Awakening) {

}

func (g *Game) DisassembleItem(p *Player, msg *client.DisassembleItem) {

}

func (g *Game) DowngradeAwakening(p *Player, msg *client.DowngradeAwakening) {

}

func (g *Game) ResetAddedItem(p *Player, msg *client.ResetAddedItem) {

}

func (g *Game) SendMail(p *Player, msg *client.SendMail) {

}

func (g *Game) ReadMail(p *Player, msg *client.ReadMail) {

}

func (g *Game) CollectParcel(p *Player, msg *client.CollectParcel) {

}

func (g *Game) DeleteMail(p *Player, msg *client.DeleteMail) {

}

func (g *Game) LockMail(p *Player, msg *client.LockMail) {

}

func (g *Game) MailLockedItem(p *Player, msg *client.MailLockedItem) {

}

func (g *Game) MailCost(p *Player, msg *client.MailCost) {

}

func (g *Game) UpdateIntelligentCreature(p *Player, msg *client.UpdateIntelligentCreature) {

}

func (g *Game) IntelligentCreaturePickup(p *Player, msg *client.IntelligentCreaturePickup) {

}

func (g *Game) AddFriend(p *Player, msg *client.AddFriend) {

}

func (g *Game) RemoveFriend(p *Player, msg *client.RemoveFriend) {

}

func (g *Game) RefreshFriends(p *Player, msg *client.RefreshFriends) {

}

func (g *Game) AddMemo(p *Player, msg *client.AddMemo) {

}

func (g *Game) GuildBuffUpdate(p *Player, msg *client.GuildBuffUpdate) {

}

func (g *Game) GameshopBuy(p *Player, msg *client.GameshopBuy) {

}

func (g *Game) NPCConfirmInput(p *Player, msg *client.NPCConfirmInput) {

}

func (g *Game) ReportIssue(p *Player, msg *client.ReportIssue) {

}

func (g *Game) GetRanking(p *Player, msg *client.GetRanking) {

}

func (g *Game) Opendoor(p *Player, msg *client.Opendoor) {
	p.OpenDoor(msg.DoorIndex)
}

func (g *Game) GetRentedItems(p *Player, msg *client.GetRentedItems) {

}

func (g *Game) ItemRentalRequest(p *Player, msg *client.ItemRentalRequest) {

}

func (g *Game) ItemRentalFee(p *Player, msg *client.ItemRentalFee) {

}

func (g *Game) ItemRentalPeriod(p *Player, msg *client.ItemRentalPeriod) {

}

func (g *Game) DepositRentalItem(p *Player, msg *client.DepositRentalItem) {

}

func (g *Game) RetrieveRentalItem(p *Player, msg *client.RetrieveRentalItem) {

}

func (g *Game) CancelItemRental(p *Player, msg *client.CancelItemRental) {

}

func (g *Game) ItemRentalLockFee(p *Player, msg *client.ItemRentalLockFee) {

}

func (g *Game) ItemRentalLockItem(p *Player, msg *client.ItemRentalLockItem) {

}

func (g *Game) ConfirmItemRental(p *Player, msg *client.ConfirmItemRental) {

}
