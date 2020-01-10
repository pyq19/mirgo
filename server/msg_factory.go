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
		Location:         p.CurrentLocation,
		Direction:        p.CurrentDirection,
		Hair:             p.Hair,
		Light:            0, // TODO
		Weapon:           0,
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

// TODO
func (ServerMessage) ObjectMonster(m *Monster) *server.ObjectMonster {
	return &server.ObjectMonster{
		ObjectID:          m.ID,
		Name:              m.Name,
		NameColor:         common.Color{}.ToInt32(),
		Image:             0,
		Direction:         0,
		Effect:            0,
		AI:                0,
		Light:             0,
		Dead:              false,
		Skeleton:          false,
		Poison:            0,
		Hidden:            false,
		Extra:             false,
		ExtraByte:         0,
		ShockTime:         0,
		BindingShotCenter: false,
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
		Image:     0,                                                      // TODO
		Color:     common.Color{R: 255, G: 255, B: 255, A: 255}.ToInt32(), // TODO
		Location:  n.GetPoint(),
		Direction: n.CurrentDirection,
		QuestIDs:  nil, // TODO
	}
}
