package main

import (
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
)

type ServerMessage struct{}

func (ServerMessage) SetConcentration() *server.SetConcentration {
	sc := new(server.SetConcentration)
	sc.ObjectID = 66432
	sc.Enabled = false
	sc.Interrupted = false
	return sc
}

// TODO
func (ServerMessage) ObjectPlayer(p *Player) *server.ObjectPlayer {
	return &server.ObjectPlayer{
		ObjectID:         p.ID,
		Name:             p.Name,
		GuildName:        p.GuildName,
		GuildRankName:    p.GuildRankName,
		NameColor:        p.NameColour.ToInt32(),
		Class:            p.Class,
		Gender:           p.Gender,
		Level:            p.Level,
		Location:         p.GetPoint(),
		Direction:        p.GetDirection(),
		Hair:             p.Hair,
		Light:            p.Light,
		Weapon:           0, // TODO
		WeaponEffect:     0,
		Armour:           0,
		Poison:           common.PoisonTypeNone,
		Dead:             false,
		Hidden:           false,
		Effect:           common.SpellEffectNone,
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
		LevelEffects:     common.LevelEffectsNone,
	}
}

func (ServerMessage) ObjectMonster(m *Monster) *server.ObjectMonster {
	return &server.ObjectMonster{
		ObjectID:          m.ID,
		Name:              m.Name,
		NameColor:         common.Color{}.ToInt32(),
		Location:          m.GetPoint(),
		Image:             m.Image,
		Direction:         m.GetDirection(),
		Effect:            uint8(m.Effect),
		AI:                uint8(m.AI),
		Light:             m.Light,
		Dead:              m.IsDead(),
		Skeleton:          m.IsSkeleton(),
		Poison:            m.Poison,
		Hidden:            false, // TODO
		ShockTime:         0,     // TODO
		BindingShotCenter: false, // TODO
		Extra:             false, // TODO
		ExtraByte:         nil,   // TODO
	}
}

func (ServerMessage) MapInformation(info *common.MapInfo) *server.MapInformation {
	mi := new(server.MapInformation)
	mi.FileName = info.Filename
	mi.Title = info.Title
	mi.MiniMap = uint16(info.MineIndex)
	mi.BigMap = uint16(info.BigMap)
	mi.Music = uint16(info.Music)
	mi.Lights = common.LightSetting(info.Light)
	mi.Lightning = true
	mi.MapDarkLight = 0
	return mi
}

func (ServerMessage) StartGame() *server.StartGame {
	sg := new(server.StartGame)
	sg.Result = 4
	sg.Resolution = 1024
	return sg
}

func (ServerMessage) UserInformation(p *Player) *server.UserInformation {
	ui := new(server.UserInformation)
	ui.ObjectID = p.GetID()
	ui.RealID = p.GetID()
	ui.Name = p.Name
	ui.GuildName = p.GuildName
	ui.GuildRank = p.GuildRankName
	ui.NameColor = common.Color{R: 255, G: 255, B: 255, A: 255}.ToInt32()
	ui.Class = p.Class
	ui.Gender = p.Gender
	ui.Level = p.Level
	ui.Location = p.CurrentLocation
	ui.Direction = p.CurrentDirection
	ui.Hair = p.Hair
	ui.HP = p.HP
	ui.MP = p.MP
	ui.Experience = p.Experience
	ui.MaxExperience = 100 // TODO
	ui.LevelEffect = common.LevelEffects(1)
	ui.Gold = uint32(p.Gold)
	ui.Credit = 100 // TODO
	ui.Inventory = p.Inventory
	ui.Equipment = p.Equipment
	ui.QuestInventory = p.QuestInventory
	return ui
}

func (ServerMessage) UserLocation(p *Player) *server.UserLocation {
	return &server.UserLocation{
		Location:  p.Point(),
		Direction: p.CurrentDirection,
	}
}

func (ServerMessage) ObjectTurn(o IMapObject) *server.ObjectTurn {
	return &server.ObjectTurn{
		ObjectID:  o.GetID(),
		Location:  o.GetPoint(),
		Direction: o.GetDirection(),
	}
}

func (ServerMessage) ObjectWalk(o IMapObject) *server.ObjectWalk {
	return &server.ObjectWalk{
		ObjectID:  o.GetID(),
		Location:  o.GetPoint(),
		Direction: o.GetDirection(),
	}
}

func (ServerMessage) ObjectRun(o IMapObject) *server.ObjectRun {
	return &server.ObjectRun{
		ObjectID:  o.GetID(),
		Location:  o.GetPoint(),
		Direction: o.GetDirection(),
	}
}

func (ServerMessage) ObjectRemove(o IMapObject) *server.ObjectRemove {
	return &server.ObjectRemove{ObjectID: o.GetID()}
}

func (ServerMessage) ObjectChat(p *Player, message string, chatType common.ChatType) *server.ObjectChat {
	text := p.Name + ":" + message
	return &server.ObjectChat{
		ObjectID: p.ID,
		Text:     text,
		Type:     chatType,
	}
}

func (ServerMessage) ObjectNPC(n *NPC) *server.ObjectNPC {
	return &server.ObjectNPC{
		ObjectID:  n.ID,
		Name:      n.Name,
		NameColor: common.Color{R: 255, G: 255, B: 255, A: 255}.ToInt32(), // TODO
		Image:     uint16(n.Image),
		Color:     common.Color{R: 255, G: 255, B: 255, A: 255}.ToInt32(), // TODO
		Location:  n.GetPoint(),
		Direction: n.GetDirection(),
		QuestIDs:  nil, // TODO
	}
}

func (ServerMessage) NewItemInfo(item *common.ItemInfo) *server.NewItemInfo {
	if item == nil {
		return nil
	}
	return &server.NewItemInfo{Info: *item}
}

func (ServerMessage) PlayerInspect(p *Player) *server.PlayerInspect {
	return &server.PlayerInspect{
		Name:      p.Name,
		GuildName: p.GuildName,
		GuildRank: p.GuildRankName,
		Equipment: p.Equipment,
		Class:     p.Class,
		Gender:    p.Gender,
		Hair:      p.Hair,
		Level:     p.Level,
		LoverName: "",
	}
}

func (ServerMessage) Login(result int) *server.Login {
	/*
	 * 0: Disabled
	 * 1: Bad AccountID
	 * 2: Bad Password
	 * 3: Account Not Exist
	 * 4: Wrong Password
	 */
	return &server.Login{Result: uint8(result)}
}

func (ServerMessage) NewCharacter(result int) interface{} {
	/*
	 * 0: Disabled.
	 * 1: Bad Character Name
	 * 2: Bad Gender
	 * 3: Bad Class
	 * 4: Max Characters
	 * 5: Character Exists.
	 * */
	return &server.NewCharacter{Result: uint8(result)}
}

func (ServerMessage) NewCharacterSuccess(g *Game, AccountID int, name string, class common.MirClass, gender common.MirGender) *server.NewCharacterSuccess {
	c := new(common.Character)
	c.Name = name
	c.Level = 8
	c.Class = class
	c.Gender = gender
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
	g.DB.Table("character").Where("name = ?", name).Last(c)
	ac := new(common.AccountCharacter)
	ac.AccountID = AccountID
	ac.CharacterID = int(c.ID)
	g.DB.Table("account_character").Create(ac)
	res := new(server.NewCharacterSuccess)
	res.CharInfo.Index = uint32(c.ID)
	res.CharInfo.Name = name
	res.CharInfo.Class = class
	res.CharInfo.Gender = gender
	return res
}

func (ServerMessage) LogOutSuccess(characters []common.SelectInfo) *server.LogOutSuccess {
	return &server.LogOutSuccess{Characters: characters}
}
