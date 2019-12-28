package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/golog"
	_ "github.com/yenkeia/mirgo/codec/mircodec"
	"github.com/yenkeia/mirgo/common"
	_ "github.com/yenkeia/mirgo/proc/mirtcp"
	"github.com/yenkeia/mirgo/proto/client"
	"github.com/yenkeia/mirgo/proto/server"
)

var log = golog.New("server.handler")

func (g *Game) HandleEvent(ev cellnet.Event) {
	g.Pool.Submit(NewTask(_HandleEvent, g, ev))
}

func _HandleEvent(args ...interface{}) {
	var (
		g  *Game
		ev cellnet.Event
		s  cellnet.Session
	)
	g = args[0].(*Game)
	ev = args[1].(cellnet.Event)
	s = ev.Session()
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
}

// SessionAccepted ...
func (g *Game) SessionAccepted(s cellnet.Session, msg *cellnet.SessionAccepted) {
	connected := server.Connected{}
	s.Send(&connected)
}

// SessionClosed ...
func (g *Game) SessionClosed(s cellnet.Session, msg *cellnet.SessionClosed) {
	pm := g.Env.SessionIDPlayerMap
	pm.Delete(s.ID())
}

// ClientVersion ...
func (g *Game) ClientVersion(s cellnet.Session, msg *client.ClientVersion) {
	clientVersion := server.ClientVersion{Result: 1}
	p := new(Player)
	p.GameStage = LOGIN
	p.Session = &s
	g.Env.SessionIDPlayerMap.Store(s.ID(), p)
	s.Send(&clientVersion)
}

func (g *Game) GetPlayer(s cellnet.Session, gameStage int) (p *Player, ok bool) {
	v, ok := g.Env.SessionIDPlayerMap.Load(s.ID())
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
	keepAlive := server.KeepAlive{Time: 0}
	s.Send(keepAlive)
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
	ac := new(common.Account)
	g.DB.Table("account").Where("username = ?", msg.AccountID).Find(ac)
	if ac.ID == 0 && ac.Username == "" {
		ac.Username = msg.AccountID
		ac.Password = msg.Password
		g.DB.Table("account").Create(&ac)
		res = 8
	}
	s.Send(server.NewAccount{Result: res})
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
	ac := new(common.Account)
	g.DB.Table("account").Where("username = ? AND password = ?", msg.AccountID, msg.CurrentPassword).Find(ac)
	if ac.ID != 0 {
		ac.Password = msg.NewPassword
		g.DB.Table("account").Model(ac).Updates(common.Account{Password: msg.NewPassword})
		res = 6
	}
	s.Send(server.ChangePassword{Result: res})
}

// Login 登陆
func (g *Game) Login(s cellnet.Session, msg *client.Login) {
	/*
	 * 0: Disabled
	 * 1: Bad AccountID
	 * 2: Bad Password
	 * 3: Account Not Exist
	 * 4: Wrong Password
	 */
	p, ok := g.GetPlayer(s, LOGIN)
	if !ok {
		return
	}

	a := new(common.Account)
	g.DB.Table("account").Where("username = ? AND password = ?", msg.AccountID, msg.Password).Find(a)
	if a.ID == 0 {
		s.Send(server.Login{Result: uint8(4)})
		return
	}

	p.AccountID = a.ID
	p.GameStage = SELECT

	ac := make([]common.AccountCharacter, 3)
	g.DB.Table("account_character").Where("account_id = ?", a.ID).Limit(3).Find(&ac)
	ids := make([]int, 3)
	for _, c := range ac {
		ids = append(ids, c.ID)
	}
	cs := make([]common.Character, 3)
	//db.Where("name in (?)", []string{"jinzhu", "jinzhu 2"}).Find(&users)
	g.DB.Table("character").Where("id in (?)", ids).Find(&cs)
	si := make([]common.SelectInfo, len(cs))
	for i, c := range cs {
		s := new(common.SelectInfo)
		s.Index = uint32(c.ID)
		s.Name = c.Name
		s.Level = c.Level
		s.Class = c.Class
		s.Gender = c.Gender
		s.LastAccess = 0
		si[i] = *s
	}
	res := new(server.LoginSuccess)
	res.Characters = si
	s.Send(res)
}

// NewCharacter 创建角色
func (g *Game) NewCharacter(s cellnet.Session, msg *client.NewCharacter) {
	p, ok := g.GetPlayer(s, SELECT)
	if !ok {
		return
	}

	acs := make([]common.AccountCharacter, 3)
	g.DB.Table("account_character").Where("account_id = ?", p.AccountID).Limit(3).Find(&acs)
	if len(acs) >= 3 {
		n := new(server.NewCharacter)
		/*
		 * 0: Disabled.
		 * 1: Bad Character Name
		 * 2: Bad Gender
		 * 3: Bad Class
		 * 4: Max Characters
		 * 5: Character Exists.
		 * */
		n.Result = uint8(4)
		s.Send(n)
		return
	}
	c := new(common.Character)
	c.Name = msg.Name
	c.Level = 8
	c.Class = msg.Class
	c.Gender = msg.Gender
	c.Hair = 1
	c.CurrentMapID = 1
	c.CurrentLocationX = 284
	c.CurrentLocationY = 608
	c.Direction = common.MirDirectionDown
	c.HP = 15
	c.MP = 17
	c.Experience = 0
	c.AttackMode = common.AttackModeAll
	c.PetMode = common.PetModeBoth
	g.DB.Table("character").Create(c)
	g.DB.Table("character").Where("name = ?", msg.Name).Last(c)
	ac := new(common.AccountCharacter)
	ac.AccountID = p.AccountID
	ac.CharacterID = int(c.ID)
	g.DB.Table("account_character").Create(ac)
	// log.Debugln(msg.Name, msg.Class, msg.Gender)
	// user item
	res := new(server.NewCharacterSuccess)
	res.CharInfo.Index = uint32(c.ID)
	res.CharInfo.Name = msg.Name
	res.CharInfo.Class = msg.Class
	res.CharInfo.Gender = msg.Gender
	s.Send(res)
}

// DeleteCharacter 删除角色
func (g *Game) DeleteCharacter(s cellnet.Session, msg *client.DeleteCharacter) {
	_, ok := g.GetPlayer(s, SELECT)
	if !ok {
		return
	}

	c := new(common.Character)
	g.DB.Table("character").Where("id = ?", msg.CharacterIndex).Find(c)
	if c.ID == 0 {
		res := new(server.DeleteCharacter)
		res.Result = 4
		s.Send(res)
		return
	}
	g.DB.Table("character").Delete(c)
	g.DB.Table("account_character").Where("character_id = ?", c.ID).Delete(common.Character{})
	res := new(server.DeleteCharacterSuccess)
	res.CharacterIndex = msg.CharacterIndex
	s.Send(res)
}

// StartGame 开始游戏
func (g *Game) StartGame(s cellnet.Session, msg *client.StartGame) {
	p, ok := g.GetPlayer(s, SELECT)
	if !ok {
		return
	}

	c := new(common.Character)
	g.DB.Table("character").Where("id = ?", msg.CharacterIndex).Find(c)
	if c.ID == 0 {
		return
	}
	p.Character = c
	p.GameStage = GAME

	m := g.Env.GetMap(int(p.Character.CurrentMapID))
	m.AddObject(p)
	p.Map = m

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
	mi := new(server.MapInformation)
	pmi := g.Env.GameDB.GetMapInfoByID(int(p.Character.CurrentMapID))
	mi.FileName = pmi.Filename
	mi.Title = pmi.Title
	mi.MiniMap = uint16(pmi.MineIndex)
	mi.BigMap = uint16(pmi.BigMap)
	mi.Music = uint16(pmi.Music)
	mi.Lights = common.LightSetting(pmi.Light)
	mi.Lightning = true
	mi.MapDarkLight = 0
	s.Send(mi)

	ui := new(server.UserInformation)
	ui.ObjectID = 66432 // TODO
	ui.RealID = uint32(p.Character.ID)
	ui.Name = p.Character.Name
	ui.GuildName = ""
	ui.GuildRank = ""
	ui.NameColour = common.Color{R: 255, G: 255, B: 255, A: 255}.ToUint32()
	ui.Class = p.Character.Class
	ui.Gender = p.Character.Gender
	ui.Level = p.Character.Level
	ui.Location = common.Point{X: uint32(p.Character.CurrentLocationX), Y: uint32(p.Character.CurrentLocationY)}
	ui.Direction = p.Character.Direction
	ui.Hair = p.Character.Hair
	ui.HP = p.Character.HP
	ui.MP = p.Character.MP
	ui.Experience = p.Character.Experience
	ui.MaxExperience = 100 // TODO
	ui.LevelEffect = common.LevelEffects(1)
	ui.Gold = 100   // TODO
	ui.Credit = 100 // TODO

	cui := make([]common.CharacterUserItem, 0, 100)
	g.DB.Table("character_user_item").Where("character_id = ?", p.Character.ID).Find(&cui)
	is := make([]int, 0, 46)
	es := make([]int, 0, 14)
	qs := make([]int, 0, 40)
	for _, i := range cui {
		switch common.UserItemType(i.Type) {
		case common.UserItemTypeInventory:
			is = append(is, i.UserItemID)
		case common.UserItemTypeEquipment:
			es = append(es, i.UserItemID)
		case common.UserItemTypeQuestInventory:
			qs = append(qs, i.UserItemID)
		}
	}
	ui.Inventory = make([]common.UserItem, 46)
	ui.Equipment = make([]common.UserItem, 14)
	ui.QuestInventory = make([]common.UserItem, 40)
	uii := make([]common.UserItem, 0, 46)
	uie := make([]common.UserItem, 0, 14)
	uiq := make([]common.UserItem, 0, 40)
	g.DB.Table("user_item").Where("id in (?)", is).Find(&uii)
	g.DB.Table("user_item").Where("id in (?)", es).Find(&uie)
	g.DB.Table("user_item").Where("id in (?)", qs).Find(&uiq)
	for i, v := range uii {
		ui.Inventory[i] = v
		ii := g.Env.GameDB.GetItemInfoByID(int(v.ItemID))
		s.Send(&server.NewItemInfo{Info: *ii})
	}
	for i, v := range uie {
		ui.Equipment[i] = v
		ii := g.Env.GameDB.GetItemInfoByID(int(v.ItemID))
		s.Send(&server.NewItemInfo{Info: *ii})
	}
	for i, v := range uiq {
		ui.QuestInventory[i] = v
		ii := g.Env.GameDB.GetItemInfoByID(int(v.ItemID))
		s.Send(&server.NewItemInfo{Info: *ii})
	}
	s.Send(ui)

	// TODO
	p.NotifySurroundingPlayer(&server.ObjectPlayer{
		ObjectID:         uint32(c.ID),
		Name:             c.Name,
		GuildName:        "",
		GuildRankName:    "",
		NameColour:       common.Color{R: 255, G: 255, B: 255, A: 255}.ToInt32(),
		Class:            c.Class,
		Gender:           c.Gender,
		Level:            c.Level,
		Location:         common.Point{X: uint32(c.CurrentLocationX), Y: uint32(c.CurrentLocationY)},
		Direction:        c.Direction,
		Hair:             c.Hair,
		Light:            0, // TODO
		Weapon:           0,
		WeaponEffect:     0,
		Armour:           0,
		Poison:           0,
		Dead:             false,
		Hidden:           false,
		Effect:           0,
		WingEffect:       0,
		Extra:            false,
		MountType:        0,
		RidingMount:      false,
		Fishing:          false,
		TransformType:    0,
		ElementOrbEffect: 0,
		ElementOrbLvl:    0,
		ElementOrbMax:    0,
		Buffs:            nil,
		LevelEffects:     0,
	})
}

func (g *Game) LogOut(s cellnet.Session, msg *client.LogOut) {

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
	p.DropGold(msg.Amount)
}

func (g *Game) PickUp(p *Player, msg *client.PickUp) {
	p.PickUp()
}

func (g *Game) Inspect(p *Player, msg *client.Inspect) {

}

func (g *Game) ChangeAMode(p *Player, msg *client.ChangeAMode) {

}

func (g *Game) ChangePMode(p *Player, msg *client.ChangePMode) {

}

func (g *Game) ChangeTrade(p *Player, msg *client.ChangeTrade) {

}

func (g *Game) Attack(p *Player, msg *client.Attack) {

}

func (g *Game) RangeAttack(p *Player, msg *client.RangeAttack) {

}

func (g *Game) Harvest(p *Player, msg *client.Harvest) {

}

func (g *Game) CallNPC(p *Player, msg *client.CallNPC) {

}

func (g *Game) TalkMonsterNPC(p *Player, msg *client.TalkMonsterNPC) {

}

func (g *Game) BuyItem(p *Player, msg *client.BuyItem) {

}

func (g *Game) CraftItem(p *Player, msg *client.CraftItem) {

}

func (g *Game) SellItem(p *Player, msg *client.SellItem) {

}

func (g *Game) RepairItem(p *Player, msg *client.RepairItem) {

}

func (g *Game) BuyItemBack(p *Player, msg *client.BuyItemBack) {

}

func (g *Game) SRepairItem(p *Player, msg *client.SRepairItem) {

}

func (g *Game) MagicKey(p *Player, msg *client.MagicKey) {

}

func (g *Game) Magic(p *Player, msg *client.Magic) {

}

func (g *Game) SwitchGroup(p *Player, msg *client.SwitchGroup) {

}

func (g *Game) AddMember(p *Player, msg *client.AddMember) {

}

func (g *Game) DelMember(p *Player, msg *client.DelMember) {

}

func (g *Game) GroupInvite(p *Player, msg *client.GroupInvite) {

}

func (g *Game) TownRevive(p *Player, msg *client.TownRevive) {

}

func (g *Game) SpellToggle(p *Player, msg *client.SpellToggle) {

}

func (g *Game) ConsignItem(p *Player, msg *client.ConsignItem) {

}

func (g *Game) MarketSearch(p *Player, msg *client.MarketSearch) {

}

func (g *Game) MarketRefresh(p *Player, msg *client.MarketRefresh) {

}

func (g *Game) MarketPage(p *Player, msg *client.MarketPage) {

}

func (g *Game) MarketBuy(p *Player, msg *client.MarketBuy) {

}

func (g *Game) MarketGetBack(p *Player, msg *client.MarketGetBack) {

}

func (g *Game) RequestUserName(p *Player, msg *client.RequestUserName) {

}

func (g *Game) RequestChatItem(p *Player, msg *client.RequestChatItem) {

}

func (g *Game) EditGuildMember(p *Player, msg *client.EditGuildMember) {

}

func (g *Game) EditGuildNotice(p *Player, msg *client.EditGuildNotice) {

}

func (g *Game) GuildInvite(p *Player, msg *client.GuildInvite) {

}

func (g *Game) RequestGuildInfo(p *Player, msg *client.RequestGuildInfo) {

}

func (g *Game) GuildNameReturn(p *Player, msg *client.GuildNameReturn) {

}

func (g *Game) GuildStorageGoldChange(p *Player, msg *client.GuildStorageGoldChange) {

}

func (g *Game) GuildStorageItemChange(p *Player, msg *client.GuildStorageItemChange) {

}

func (g *Game) GuildWarReturn(p *Player, msg *client.GuildWarReturn) {

}

func (g *Game) MarriageRequest(p *Player, msg *client.MarriageRequest) {

}

func (g *Game) MarriageReply(p *Player, msg *client.MarriageReply) {

}

func (g *Game) ChangeMarriage(p *Player, msg *client.ChangeMarriage) {

}

func (g *Game) DivorceRequest(p *Player, msg *client.DivorceRequest) {

}

func (g *Game) DivorceReply(p *Player, msg *client.DivorceReply) {

}

func (g *Game) AddMentor(p *Player, msg *client.AddMentor) {

}

func (g *Game) MentorReply(p *Player, msg *client.MentorReply) {

}

func (g *Game) AllowMentor(p *Player, msg *client.AllowMentor) {

}

func (g *Game) CancelMentor(p *Player, msg *client.CancelMentor) {

}

func (g *Game) TradeRequest(p *Player, msg *client.TradeRequest) {

}

func (g *Game) TradeGold(p *Player, msg *client.TradeGold) {

}

func (g *Game) TradeReply(p *Player, msg *client.TradeReply) {

}

func (g *Game) TradeConfirm(p *Player, msg *client.TradeConfirm) {

}

func (g *Game) TradeCancel(p *Player, msg *client.TradeCancel) {

}

func (g *Game) EquipSlotItem(p *Player, msg *client.EquipSlotItem) {

}

func (g *Game) FishingCast(p *Player, msg *client.FishingCast) {

}

func (g *Game) FishingChangeAutocast(p *Player, msg *client.FishingChangeAutocast) {

}

func (g *Game) AcceptQuest(p *Player, msg *client.AcceptQuest) {

}

func (g *Game) FinishQuest(p *Player, msg *client.FinishQuest) {

}

func (g *Game) AbandonQuest(p *Player, msg *client.AbandonQuest) {

}

func (g *Game) ShareQuest(p *Player, msg *client.ShareQuest) {

}

func (g *Game) AcceptReincarnation(p *Player, msg *client.AcceptReincarnation) {

}

func (g *Game) CancelReincarnation(p *Player, msg *client.CancelReincarnation) {

}

func (g *Game) CombineItem(p *Player, msg *client.CombineItem) {

}

func (g *Game) SetConcentration(p *Player, msg *client.SetConcentration) {

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
