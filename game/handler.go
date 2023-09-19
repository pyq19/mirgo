package game

import (
	"time"

	"github.com/davyxu/cellnet"
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
	case *CM_ClientVersion:
		g.ClientVersion(s, msg)
	case *CM_KeepAlive:
		g.KeepAlive(s, msg)
	case *CM_NewAccount:
		g.NewAccount(s, msg)
	case *CM_ChangePassword:
		g.ChangePassword(s, msg)
	case *CM_Login:
		g.Login(s, msg)
	case *CM_NewCharacter:
		g.NewCharacter(s, msg)
	case *CM_DeleteCharacter:
		g.DeleteCharacter(s, msg)
	case *CM_StartGame:
		g.StartGame(s, msg)
	case *CM_LogOut:
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
	case *CM_Turn:
		g.Turn(p, msg)
	case *CM_Walk:
		g.Walk(p, msg)
	case *CM_Run:
		g.Run(p, msg)
	case *CM_Chat:
		g.Chat(p, msg)
	case *CM_MoveItem:
		g.MoveItem(p, msg)
	case *CM_StoreItem:
		g.StoreItem(p, msg)
	case *CM_DepositRefineItem:
		g.DepositRefineItem(p, msg)
	case *CM_RetrieveRefineItem:
		g.RetrieveRefineItem(p, msg)
	case *CM_RefineCancel:
		g.RefineCancel(p, msg)
	case *CM_RefineItem:
		g.RefineItem(p, msg)
	case *CM_CheckRefine:
		g.CheckRefine(p, msg)
	case *CM_ReplaceWedRing:
		g.ReplaceWedRing(p, msg)
	case *CM_DepositTradeItem:
		g.DepositTradeItem(p, msg)
	case *CM_RetrieveTradeItem:
		g.RetrieveTradeItem(p, msg)
	case *CM_TakeBackItem:
		g.TakeBackItem(p, msg)
	case *CM_MergeItem:
		g.MergeItem(p, msg)
	case *CM_EquipItem:
		g.EquipItem(p, msg)
	case *CM_RemoveItem:
		g.RemoveItem(p, msg)
	case *CM_RemoveSlotItem:
		g.RemoveSlotItem(p, msg)
	case *CM_SplitItem:
		g.SplitItem(p, msg)
	case *CM_UseItem:
		g.UseItem(p, msg)
	case *CM_DropItem:
		g.DropItem(p, msg)
	case *CM_DropGold:
		g.DropGold(p, msg)
	case *CM_PickUp:
		g.PickUp(p, msg)
	case *CM_Inspect:
		g.Inspect(p, msg)
	case *CM_ChangeAMode:
		g.ChangeAMode(p, msg)
	case *CM_ChangePMode:
		g.ChangePMode(p, msg)
	case *CM_ChangeTrade:
		g.ChangeTrade(p, msg)
	case *CM_Attack:
		g.Attack(p, msg)
	case *CM_RangeAttack:
		g.RangeAttack(p, msg)
	case *CM_Harvest:
		g.Harvest(p, msg)
	case *CM_CallNPC:
		g.CallNPC(p, msg)
	case *CM_TalkMonsterNPC:
		g.TalkMonsterNPC(p, msg)
	case *CM_BuyItem:
		g.BuyItem(p, msg)
	case *CM_CraftItem:
		g.CraftItem(p, msg)
	case *CM_SellItem:
		g.SellItem(p, msg)
	case *CM_RepairItem:
		g.RepairItem(p, msg)
	case *CM_BuyItemBack:
		g.BuyItemBack(p, msg)
	case *CM_SRepairItem:
		g.SRepairItem(p, msg)
	case *CM_MagicKey:
		g.MagicKey(p, msg)
	case *CM_Magic:
		g.Magic(p, msg)
	case *CM_SwitchGroup:
		g.SwitchGroup(p, msg)
	case *CM_AddMember:
		g.AddMember(p, msg)
	case *CM_DelMember:
		g.DelMember(p, msg)
	case *CM_GroupInvite:
		g.GroupInvite(p, msg)
	case *CM_TownRevive:
		g.TownRevive(p, msg)
	case *CM_SpellToggle:
		g.SpellToggle(p, msg)
	case *CM_ConsignItem:
		g.ConsignItem(p, msg)
	case *CM_MarketSearch:
		g.MarketSearch(p, msg)
	case *CM_MarketRefresh:
		g.MarketRefresh(p, msg)
	case *CM_MarketPage:
		g.MarketPage(p, msg)
	case *CM_MarketBuy:
		g.MarketBuy(p, msg)
	case *CM_MarketGetBack:
		g.MarketGetBack(p, msg)
	case *CM_RequestUserName:
		g.RequestUserName(p, msg)
	case *CM_RequestChatItem:
		g.RequestChatItem(p, msg)
	case *CM_EditGuildMember:
		g.EditGuildMember(p, msg)
	case *CM_EditGuildNotice:
		g.EditGuildNotice(p, msg)
	case *CM_GuildInvite:
		g.GuildInvite(p, msg)
	case *CM_RequestGuildInfo:
		g.RequestGuildInfo(p, msg)
	case *CM_GuildNameReturn:
		g.GuildNameReturn(p, msg)
	case *CM_GuildStorageGoldChange:
		g.GuildStorageGoldChange(p, msg)
	case *CM_GuildStorageItemChange:
		g.GuildStorageItemChange(p, msg)
	case *CM_GuildWarReturn:
		g.GuildWarReturn(p, msg)
	case *CM_MarriageRequest:
		g.MarriageRequest(p, msg)
	case *CM_MarriageReply:
		g.MarriageReply(p, msg)
	case *CM_ChangeMarriage:
		g.ChangeMarriage(p, msg)
	case *CM_DivorceRequest:
		g.DivorceRequest(p, msg)
	case *CM_DivorceReply:
		g.DivorceReply(p, msg)
	case *CM_AddMentor:
		g.AddMentor(p, msg)
	case *CM_MentorReply:
		g.MentorReply(p, msg)
	case *CM_AllowMentor:
		g.AllowMentor(p, msg)
	case *CM_CancelMentor:
		g.CancelMentor(p, msg)
	case *CM_TradeRequest:
		g.TradeRequest(p, msg)
	case *CM_TradeGold:
		g.TradeGold(p, msg)
	case *CM_TradeReply:
		g.TradeReply(p, msg)
	case *CM_TradeConfirm:
		g.TradeConfirm(p, msg)
	case *CM_TradeCancel:
		g.TradeCancel(p, msg)
	case *CM_EquipSlotItem:
		g.EquipSlotItem(p, msg)
	case *CM_FishingCast:
		g.FishingCast(p, msg)
	case *CM_FishingChangeAutocast:
		g.FishingChangeAutocast(p, msg)
	case *CM_AcceptQuest:
		g.AcceptQuest(p, msg)
	case *CM_FinishQuest:
		g.FinishQuest(p, msg)
	case *CM_AbandonQuest:
		g.AbandonQuest(p, msg)
	case *CM_ShareQuest:
		g.ShareQuest(p, msg)
	case *CM_AcceptReincarnation:
		g.AcceptReincarnation(p, msg)
	case *CM_CancelReincarnation:
		g.CancelReincarnation(p, msg)
	case *CM_CombineItem:
		g.CombineItem(p, msg)
	case *CM_SetConcentration:
		g.SetConcentration(p, msg)
	case *CM_AwakeningNeedMaterials:
		g.AwakeningNeedMaterials(p, msg)
	case *CM_AwakeningLockedItem:
		g.AwakeningLockedItem(p, msg)
	case *CM_Awakening:
		g.Awakening(p, msg)
	case *CM_DisassembleItem:
		g.DisassembleItem(p, msg)
	case *CM_DowngradeAwakening:
		g.DowngradeAwakening(p, msg)
	case *CM_ResetAddedItem:
		g.ResetAddedItem(p, msg)
	case *CM_SendMail:
		g.SendMail(p, msg)
	case *CM_ReadMail:
		g.ReadMail(p, msg)
	case *CM_CollectParcel:
		g.CollectParcel(p, msg)
	case *CM_DeleteMail:
		g.DeleteMail(p, msg)
	case *CM_LockMail:
		g.LockMail(p, msg)
	case *CM_MailLockedItem:
		g.MailLockedItem(p, msg)
	case *CM_MailCost:
		g.MailCost(p, msg)
	case *CM_UpdateIntelligentCreature: //IntelligentCreature
		g.UpdateIntelligentCreature(p, msg)
	case *CM_IntelligentCreaturePickup: //IntelligentCreature
		g.IntelligentCreaturePickup(p, msg)
	case *CM_AddFriend:
		g.AddFriend(p, msg)
	case *CM_RemoveFriend:
		g.RemoveFriend(p, msg)
	case *CM_RefreshFriends:
		g.RefreshFriends(p, msg)
	case *CM_AddMemo:
		g.AddMemo(p, msg)
	case *CM_GuildBuffUpdate:
		g.GuildBuffUpdate(p, msg)
	case *CM_GameshopBuy:
		g.GameshopBuy(p, msg)
	case *CM_NPCConfirmInput:
		g.NPCConfirmInput(p, msg)
	case *CM_ReportIssue:
		g.ReportIssue(p, msg)
	case *CM_GetRanking:
		g.GetRanking(p, msg)
	case *CM_Opendoor:
		g.Opendoor(p, msg)
	case *CM_GetRentedItems:
		g.GetRentedItems(p, msg)
	case *CM_ItemRentalRequest:
		g.ItemRentalRequest(p, msg)
	case *CM_ItemRentalFee:
		g.ItemRentalFee(p, msg)
	case *CM_ItemRentalPeriod:
		g.ItemRentalPeriod(p, msg)
	case *CM_DepositRentalItem:
		g.DepositRentalItem(p, msg)
	case *CM_RetrieveRentalItem:
		g.RetrieveRentalItem(p, msg)
	case *CM_CancelItemRental:
		g.CancelItemRental(p, msg)
	case *CM_ItemRentalLockFee:
		g.ItemRentalLockFee(p, msg)
	case *CM_ItemRentalLockItem:
		g.ItemRentalLockItem(p, msg)
	case *CM_ConfirmItemRental:
		g.ConfirmItemRental(p, msg)
	default:
		log.Debugln("default:", msg)
		//MessageQueue.Enqueue(string.Format("Invalid packet received. Index : {0}", p.Index));
	}
}

// SessionAccepted ...
func (g *Game) SessionAccepted(s cellnet.Session, msg *cellnet.SessionAccepted) {
	s.Send(&SM_Connected{})
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
func (g *Game) ClientVersion(s cellnet.Session, msg *CM_ClientVersion) {
	clientVersion := SM_ClientVersion{Result: 1}
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
func (g *Game) KeepAlive(s cellnet.Session, msg *CM_KeepAlive) {
	s.Send(SM_KeepAlive{Time: 0})
}

// NewAccount 保存新账号
func (g *Game) NewAccount(s cellnet.Session, msg *CM_NewAccount) {
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
	ac := new(Account)
	adb.Table("account").Where("username = ?", msg.AccountID).Find(ac)
	if ac.ID == 0 && ac.Username == "" {
		ac.Username = msg.AccountID
		ac.Password = msg.Password
		adb.Table("account").Create(&ac)
		res = 8
	}
	s.Send(&SM_NewAccount{Result: res})
}

// ChangePassword 改密码
func (g *Game) ChangePassword(s cellnet.Session, msg *CM_ChangePassword) {
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
	ac := new(Account)
	adb.Table("account").Where("username = ? AND password = ?", msg.AccountID, msg.CurrentPassword).Find(ac)
	if ac.ID != 0 {
		ac.Password = msg.NewPassword
		adb.Table("account").Model(ac).Updates(Account{Password: msg.NewPassword})
		res = 6
	}
	s.Send(&SM_ChangePassword{Result: res})
}

func (g *Game) getAccountCharacters(AccountID int) []SM_SelectInfo {
	ac := make([]AccountCharacter, 3)
	adb.Table("account_character").Where("account_id = ?", AccountID).Limit(3).Find(&ac)
	ids := make([]int, 3)
	for _, c := range ac {
		ids = append(ids, c.ID)
	}
	cs := make([]Character, 3)
	adb.Table("character").Where("id in (?)", ids).Find(&cs)
	si := make([]SM_SelectInfo, len(cs))
	for i, c := range cs {
		s := new(SM_SelectInfo)
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
func (g *Game) Login(s cellnet.Session, msg *CM_Login) {
	p, ok := g.GetPlayer(s, LOGIN)
	if !ok {
		return
	}
	a := new(Account)
	adb.Table("account").Where("username = ? AND password = ?", msg.AccountID, msg.Password).Find(a)
	if a.ID == 0 {
		s.Send(ServerMessage{}.Login(4))
		return
	}
	p.AccountID = a.ID
	p.GameStage = SELECT
	res := new(SM_LoginSuccess)
	res.Characters = g.getAccountCharacters(p.AccountID)
	s.Send(res)
}

// NewCharacter 创建角色
func (g *Game) NewCharacter(s cellnet.Session, msg *CM_NewCharacter) {
	p, ok := g.GetPlayer(s, SELECT)
	if !ok {
		return
	}
	acs := make([]AccountCharacter, 3)
	adb.Table("account_character").Where("account_id = ?", p.AccountID).Limit(3).Find(&acs)
	if len(acs) >= 3 {
		s.Send(ServerMessage{}.NewCharacter(4))
		return
	}
	s.Send(ServerMessage{}.NewCharacterSuccess(g, p.AccountID, msg.Name, msg.Class, msg.Gender))
}

// DeleteCharacter 删除角色
func (g *Game) DeleteCharacter(s cellnet.Session, msg *CM_DeleteCharacter) {
	_, ok := g.GetPlayer(s, SELECT)
	if !ok {
		return
	}

	c := new(Character)
	adb.Table("character").Where("id = ?", msg.CharacterIndex).Find(c)
	if c.ID == 0 {
		res := new(SM_DeleteCharacter)
		res.Result = 4
		s.Send(res)
		return
	}
	adb.Table("character").Delete(c)
	adb.Table("account_character").Where("character_id = ?", c.ID).Delete(Character{})
	res := new(SM_DeleteCharacterSuccess)
	res.CharacterIndex = msg.CharacterIndex
	s.Send(res)
}

func updatePlayerInfo(g *Game, p *Player, c *Character) {
	p.GameStage = GAME
	p.ID = uint32(c.ID)
	p.Name = c.Name
	p.NameColor = Color{R: 255, G: 255, B: 255}
	p.Direction = c.Direction
	p.CurrentLocation = NewPoint(int(c.CurrentLocationX), int(c.CurrentLocationY))
	p.BindLocation = NewPoint(c.BindLocationX, c.BindLocationY)
	p.BindMapIndex = c.BindMapID

	magics := make([]*UserMagic, 0)
	adb.Table("user_magic").Where("character_id = ?", c.ID).Find(&magics)
	for _, v := range magics {
		v.Info = data.GetMagicInfoByID(v.MagicID)
	}

	p.Inventory = BagLoadFromDB(p, UserItemTypeInventory, 46)
	p.Equipment = BagLoadFromDB(p, UserItemTypeEquipment, 14)
	p.QuestInventory = BagLoadFromDB(p, UserItemTypeQuestInventory, 40)
	p.Storage = BagLoadFromDB(p, UserItemTypeStorage, 80)
	p.Trade = NewBag(p, UserItemTypeTrade, 10)

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
	p.SendItemInfo = make([]*ItemInfo, 0)
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
func (g *Game) StartGame(s cellnet.Session, msg *CM_StartGame) {
	p, ok := g.GetPlayer(s, SELECT)
	if !ok {
		return
	}

	c := new(Character)
	adb.Table("character").Where("id = ?", msg.CharacterIndex).Find(c)
	if c.ID == 0 {
		return
	}
	ac := new(AccountCharacter)
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

func (g *Game) LogOut(s cellnet.Session, msg *CM_LogOut) {
	p, ok := g.GetPlayer(s, GAME)
	if !ok {
		return
	}
	p.GameStage = SELECT
	p.StopGame(StopGameUserReturnedToSelectChar)
	env.Players.Remove(p)
	p.Map.DeleteObject(p)
	s.Send(&SM_LogOutSuccess{Characters: g.getAccountCharacters(p.AccountID)})
	p.SaveData()
}

func (g *Game) Turn(p *Player, msg *CM_Turn) {
	p.Turn(msg.Direction)
}

func (g *Game) Walk(p *Player, msg *CM_Walk) {
	p.Walk(msg.Direction)
}

func (g *Game) Run(p *Player, msg *CM_Run) {
	p.Run(msg.Direction)
}

func (g *Game) Chat(p *Player, msg *CM_Chat) {
	p.Chat(msg.Message)
}

func (g *Game) MoveItem(p *Player, msg *CM_MoveItem) {
	p.MoveItem(msg.Grid, msg.From, msg.To)
}

func (g *Game) StoreItem(p *Player, msg *CM_StoreItem) {
	p.StoreItem(msg.From, msg.To)
}

func (g *Game) DepositRefineItem(p *Player, msg *CM_DepositRefineItem) {
	p.DepositRefineItem(msg.From, msg.To)
}

func (g *Game) RetrieveRefineItem(p *Player, msg *CM_RetrieveRefineItem) {
	p.RetrieveRefineItem(msg.From, msg.To)
}

func (g *Game) RefineCancel(p *Player, msg *CM_RefineCancel) {
	p.RefineCancel()
}

func (g *Game) RefineItem(p *Player, msg *CM_RefineItem) {
	p.RefineItem(msg.UniqueID)
}

func (g *Game) CheckRefine(p *Player, msg *CM_CheckRefine) {
	p.CheckRefine(msg.UniqueID)
}

func (g *Game) ReplaceWedRing(p *Player, msg *CM_ReplaceWedRing) {
	p.ReplaceWeddingRing(msg.UniqueID)
}

func (g *Game) DepositTradeItem(p *Player, msg *CM_DepositTradeItem) {
	p.DepositTradeItem(msg.From, msg.To)
}

func (g *Game) RetrieveTradeItem(p *Player, msg *CM_RetrieveTradeItem) {
	p.RetrieveTradeItem(msg.From, msg.To)
}

func (g *Game) TakeBackItem(p *Player, msg *CM_TakeBackItem) {
	p.TakeBackItem(msg.From, msg.To)
}

func (g *Game) MergeItem(p *Player, msg *CM_MergeItem) {
	p.MergeItem(msg.GridFrom, msg.GridTo, msg.IDFrom, msg.IDTo)
}

func (g *Game) EquipItem(p *Player, msg *CM_EquipItem) {
	p.EquipItem(msg.Grid, msg.UniqueID, msg.To)
}

func (g *Game) RemoveItem(p *Player, msg *CM_RemoveItem) {
	p.RemoveItem(msg.Grid, msg.UniqueID, msg.To)
}

func (g *Game) RemoveSlotItem(p *Player, msg *CM_RemoveSlotItem) {
	p.RemoveSlotItem(msg.Grid, msg.UniqueID, msg.To, msg.GridTo)
}

func (g *Game) SplitItem(p *Player, msg *CM_SplitItem) {
	p.SplitItem(msg.Grid, msg.UniqueID, msg.Count)
}

func (g *Game) UseItem(p *Player, msg *CM_UseItem) {
	p.UseItem(msg.UniqueID)
}

func (g *Game) DropItem(p *Player, msg *CM_DropItem) {
	p.DropItem(msg.UniqueID, msg.Count)
}

func (g *Game) DropGold(p *Player, msg *CM_DropGold) {
	p.DropGold(uint64(msg.Amount))
}

func (g *Game) PickUp(p *Player, msg *CM_PickUp) {
	p.PickUp()
}

func (g *Game) Inspect(p *Player, msg *CM_Inspect) {
	p.Inspect(msg.ObjectID)
}

func (g *Game) ChangeAMode(p *Player, msg *CM_ChangeAMode) {
	p.ChangeAMode(msg.Mode)
}

func (g *Game) ChangePMode(p *Player, msg *CM_ChangePMode) {
	p.ChangePMode(msg.Mode)
}

func (g *Game) ChangeTrade(p *Player, msg *CM_ChangeTrade) {
	p.ChangeTrade(msg.AllowTrade)
}

func (g *Game) Attack(p *Player, msg *CM_Attack) {
	p.Attack(msg.Direction, msg.Spell)
}

func (g *Game) RangeAttack(p *Player, msg *CM_RangeAttack) {
	p.RangeAttack(msg.Direction, msg.TargetLocation, msg.TargetID)
}

func (g *Game) Harvest(p *Player, msg *CM_Harvest) {
	p.Harvest(msg.Direction)
}

func (g *Game) CallNPC(p *Player, msg *CM_CallNPC) {
	p.CallNPC(msg.ObjectID, msg.Key)
}

func (g *Game) TalkMonsterNPC(p *Player, msg *CM_TalkMonsterNPC) {
	p.TalkMonsterNPC(msg.ObjectID)
}

func (g *Game) BuyItem(p *Player, msg *CM_BuyItem) {
	p.BuyItem(msg.ItemIndex, msg.Count, msg.Type)
}

// TODO
func (g *Game) CraftItem(p *Player, msg *CM_CraftItem) {
	p.CraftItem(msg.UniqueID, msg.Count, msg.Slots)
}

func (g *Game) SellItem(p *Player, msg *CM_SellItem) {
	p.SellItem(msg.UniqueID, msg.Count)
}

func (g *Game) RepairItem(p *Player, msg *CM_RepairItem) {
	p.RepairItem(msg.UniqueID, false)
}

func (g *Game) BuyItemBack(p *Player, msg *CM_BuyItemBack) {
	p.BuyItemBack(msg.UniqueID, msg.Count)
}

func (g *Game) SRepairItem(p *Player, msg *CM_SRepairItem) {
	p.RepairItem(msg.UniqueID, true)
}

func (g *Game) MagicKey(p *Player, msg *CM_MagicKey) {
	p.MagicKey(msg.Spell, msg.Key)
}

func (g *Game) Magic(p *Player, msg *CM_Magic) {
	p.Magic(msg.Spell, msg.Direction, msg.TargetID, msg.Location)
}

func (g *Game) SwitchGroup(p *Player, msg *CM_SwitchGroup) {
	p.SwitchGroup(msg.AllowGroup)
}

func (g *Game) AddMember(p *Player, msg *CM_AddMember) {
	p.AddMember(msg.Name)
}

func (g *Game) DelMember(p *Player, msg *CM_DelMember) {
	p.DelMember(msg.Name)
}

func (g *Game) GroupInvite(p *Player, msg *CM_GroupInvite) {
	p.GroupInvite(msg.AcceptInvite)
}

func (g *Game) TownRevive(p *Player, msg *CM_TownRevive) {
	p.TownRevive()
}

func (g *Game) SpellToggle(p *Player, msg *CM_SpellToggle) {
	p.SpellToggle(msg.Spell, msg.CanUse)
}

func (g *Game) ConsignItem(p *Player, msg *CM_ConsignItem) {
	p.ConsignItem(msg.UniqueID, msg.Price)
}

func (g *Game) MarketSearch(p *Player, msg *CM_MarketSearch) {
	p.MarketSearch(msg.Match)
}

func (g *Game) MarketRefresh(p *Player, msg *CM_MarketRefresh) {
	p.MarketRefresh()
}

func (g *Game) MarketPage(p *Player, msg *CM_MarketPage) {
	p.MarketPage(msg.Page)
}

func (g *Game) MarketBuy(p *Player, msg *CM_MarketBuy) {
	p.MarketBuy(msg.AuctionID)
}

func (g *Game) MarketGetBack(p *Player, msg *CM_MarketGetBack) {
	p.MarketGetBack(msg.AuctionID)
}

func (g *Game) RequestUserName(p *Player, msg *CM_RequestUserName) {
	p.RequestUserName(msg.UserID)
}

func (g *Game) RequestChatItem(p *Player, msg *CM_RequestChatItem) {
	p.RequestChatItem(msg.ChatItemID)
}

func (g *Game) EditGuildMember(p *Player, msg *CM_EditGuildMember) {
	// p.EditGuildMember(msg.Name, msg.RankName, msg.RankIndex, msg.ChangeType)
}

func (g *Game) EditGuildNotice(p *Player, msg *CM_EditGuildNotice) {
	// p.EditGuildNotice(msg.Notice)
}

func (g *Game) GuildInvite(p *Player, msg *CM_GuildInvite) {
	// p.GuildInvite(msg.AcceptInvite)
}

func (g *Game) RequestGuildInfo(p *Player, msg *CM_RequestGuildInfo) {
	// p.RequestGuildInfo(msg.Type)
}

func (g *Game) GuildNameReturn(p *Player, msg *CM_GuildNameReturn) {
	// p.GuildNameReturn(msg.Name)
}

func (g *Game) GuildStorageGoldChange(p *Player, msg *CM_GuildStorageGoldChange) {
	p.GuildStorageGoldChange(msg.Type, msg.Amount)
}

func (g *Game) GuildStorageItemChange(p *Player, msg *CM_GuildStorageItemChange) {
	p.GuildStorageItemChange(msg.Type, msg.From, msg.To)
}

func (g *Game) GuildWarReturn(p *Player, msg *CM_GuildWarReturn) {
	p.GuildWarReturn(msg.Name)
}

func (g *Game) MarriageRequest(p *Player, msg *CM_MarriageRequest) {
	p.MarriageRequest()
}

func (g *Game) MarriageReply(p *Player, msg *CM_MarriageReply) {
	p.MarriageReply(msg.AcceptInvite)
}

func (g *Game) ChangeMarriage(p *Player, msg *CM_ChangeMarriage) {
	p.ChangeMarriage()
}

func (g *Game) DivorceRequest(p *Player, msg *CM_DivorceRequest) {
	p.DivorceRequest()
}

func (g *Game) DivorceReply(p *Player, msg *CM_DivorceReply) {
	p.DivorceReply(msg.AcceptInvite)
}

func (g *Game) AddMentor(p *Player, msg *CM_AddMentor) {
	p.AddMentor(msg.Name)
}

func (g *Game) MentorReply(p *Player, msg *CM_MentorReply) {
	p.MentorReply(msg.AcceptInvite)
}

func (g *Game) AllowMentor(p *Player, msg *CM_AllowMentor) {
	p.AllowMentor()
}

func (g *Game) CancelMentor(p *Player, msg *CM_CancelMentor) {
	p.CancelMentor()
}

func (g *Game) TradeRequest(p *Player, msg *CM_TradeRequest) {
	p.TradeRequest()
}

func (g *Game) TradeGold(p *Player, msg *CM_TradeGold) {
	p.TradeGold(msg.Amount)
}

func (g *Game) TradeReply(p *Player, msg *CM_TradeReply) {
	p.TradeReply(msg.AcceptInvite)
}

func (g *Game) TradeConfirm(p *Player, msg *CM_TradeConfirm) {
	p.TradeConfirm(msg.Locked)
}

func (g *Game) TradeCancel(p *Player, msg *CM_TradeCancel) {
	p.TradeCancel()
}

func (g *Game) EquipSlotItem(p *Player, msg *CM_EquipSlotItem) {
	p.EquipSlotItem(msg.Grid, msg.UniqueID, msg.To, msg.GridTo)
}

func (g *Game) FishingCast(p *Player, msg *CM_FishingCast) {
	p.FishingCast(msg.CastOut)
}

func (g *Game) FishingChangeAutocast(p *Player, msg *CM_FishingChangeAutocast) {
	p.FishingChangeAutocast(msg.AutoCast)
}

func (g *Game) AcceptQuest(p *Player, msg *CM_AcceptQuest) {
	p.AcceptQuest(msg.NPCIndex, msg.QuestIndex)
}

func (g *Game) FinishQuest(p *Player, msg *CM_FinishQuest) {
	p.FinishQuest(msg.QuestIndex, msg.SelectedItemIndex)
}

func (g *Game) AbandonQuest(p *Player, msg *CM_AbandonQuest) {
	p.AbandonQuest(msg.QuestIndex)
}

func (g *Game) ShareQuest(p *Player, msg *CM_ShareQuest) {
	p.ShareQuest(msg.QuestIndex)
}

func (g *Game) AcceptReincarnation(p *Player, msg *CM_AcceptReincarnation) {
	p.AcceptReincarnation()
}

func (g *Game) CancelReincarnation(p *Player, msg *CM_CancelReincarnation) {
	p.CancelReincarnation()
}

func (g *Game) CombineItem(p *Player, msg *CM_CombineItem) {
	p.CombineItem(msg.IDFrom, msg.IDTo)
}

func (g *Game) SetConcentration(p *Player, msg *CM_SetConcentration) {
	p.SetConcentration(msg.ObjectID, msg.Enabled, msg.Interrupted)
}

func (g *Game) AwakeningNeedMaterials(p *Player, msg *CM_AwakeningNeedMaterials) {

}

func (g *Game) AwakeningLockedItem(p *Player, msg *CM_AwakeningLockedItem) {

}

func (g *Game) Awakening(p *Player, msg *CM_Awakening) {

}

func (g *Game) DisassembleItem(p *Player, msg *CM_DisassembleItem) {

}

func (g *Game) DowngradeAwakening(p *Player, msg *CM_DowngradeAwakening) {

}

func (g *Game) ResetAddedItem(p *Player, msg *CM_ResetAddedItem) {

}

func (g *Game) SendMail(p *Player, msg *CM_SendMail) {

}

func (g *Game) ReadMail(p *Player, msg *CM_ReadMail) {

}

func (g *Game) CollectParcel(p *Player, msg *CM_CollectParcel) {

}

func (g *Game) DeleteMail(p *Player, msg *CM_DeleteMail) {

}

func (g *Game) LockMail(p *Player, msg *CM_LockMail) {

}

func (g *Game) MailLockedItem(p *Player, msg *CM_MailLockedItem) {

}

func (g *Game) MailCost(p *Player, msg *CM_MailCost) {

}

func (g *Game) UpdateIntelligentCreature(p *Player, msg *CM_UpdateIntelligentCreature) {

}

func (g *Game) IntelligentCreaturePickup(p *Player, msg *CM_IntelligentCreaturePickup) {

}

func (g *Game) AddFriend(p *Player, msg *CM_AddFriend) {

}

func (g *Game) RemoveFriend(p *Player, msg *CM_RemoveFriend) {

}

func (g *Game) RefreshFriends(p *Player, msg *CM_RefreshFriends) {

}

func (g *Game) AddMemo(p *Player, msg *CM_AddMemo) {

}

func (g *Game) GuildBuffUpdate(p *Player, msg *CM_GuildBuffUpdate) {

}

func (g *Game) GameshopBuy(p *Player, msg *CM_GameshopBuy) {

}

func (g *Game) NPCConfirmInput(p *Player, msg *CM_NPCConfirmInput) {

}

func (g *Game) ReportIssue(p *Player, msg *CM_ReportIssue) {

}

func (g *Game) GetRanking(p *Player, msg *CM_GetRanking) {

}

func (g *Game) Opendoor(p *Player, msg *CM_Opendoor) {
	p.OpenDoor(msg.DoorIndex)
}

func (g *Game) GetRentedItems(p *Player, msg *CM_GetRentedItems) {

}

func (g *Game) ItemRentalRequest(p *Player, msg *CM_ItemRentalRequest) {

}

func (g *Game) ItemRentalFee(p *Player, msg *CM_ItemRentalFee) {

}

func (g *Game) ItemRentalPeriod(p *Player, msg *CM_ItemRentalPeriod) {

}

func (g *Game) DepositRentalItem(p *Player, msg *CM_DepositRentalItem) {

}

func (g *Game) RetrieveRentalItem(p *Player, msg *CM_RetrieveRentalItem) {

}

func (g *Game) CancelItemRental(p *Player, msg *CM_CancelItemRental) {

}

func (g *Game) ItemRentalLockFee(p *Player, msg *CM_ItemRentalLockFee) {

}

func (g *Game) ItemRentalLockItem(p *Player, msg *CM_ItemRentalLockItem) {

}

func (g *Game) ConfirmItemRental(p *Player, msg *CM_ConfirmItemRental) {

}
