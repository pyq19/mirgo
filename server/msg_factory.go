package main

import (
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
)

type ServerMessage struct{}

// TODO
func (ServerMessage) ObjectPlayer(p *Player) server.ObjectPlayer {
	return server.ObjectPlayer{
		ObjectID:         p.ID,
		Name:             p.Name,
		GuildName:        p.GuildName,
		GuildRankName:    p.GuildRankName,
		NameColour:       p.NameColour.ToInt32(),
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
	}
}

// TODO
func (ServerMessage) ObjectMonster(m *Monster) server.ObjectMonster {
	return server.ObjectMonster{
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
