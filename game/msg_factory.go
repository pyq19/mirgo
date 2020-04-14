package game

import (
	"github.com/yenkeia/mirgo/game/cm"
	"github.com/yenkeia/mirgo/game/proto/server"
)

type ServerMessage struct{}

func (ServerMessage) SetConcentration(p *Player) *server.SetConcentration {
	sc := new(server.SetConcentration)
	sc.ObjectID = uint32(p.AccountID)
	sc.Enabled = false
	sc.Interrupted = false
	return sc
}

func (ServerMessage) SetObjectConcentration(p *Player) *server.SetObjectConcentration {
	return &server.SetObjectConcentration{
		ObjectID:    uint32(p.AccountID),
		Enabled:     false,
		Interrupted: false,
	}
}

func (ServerMessage) ObjectPlayer(p *Player) (res *server.ObjectPlayer) {
	return &server.ObjectPlayer{
		ObjectID:         p.ID,
		Name:             p.Name,
		GuildName:        p.GuildName,
		GuildRankName:    p.GuildRankName,
		NameColor:        p.NameColor.ToInt32(),
		Class:            p.Class,
		Gender:           p.Gender,
		Level:            p.Level,
		Location:         p.GetPoint(),
		Direction:        p.GetDirection(),
		Hair:             p.Hair,
		Light:            p.Light,
		Weapon:           int16(p.LooksWeapon),
		WeaponEffect:     int16(p.LooksWeaponEffect),
		Armour:           int16(p.LooksArmour),
		Poison:           cm.PoisonTypeNone, // TODO
		Dead:             p.IsDead(),
		Hidden:           p.IsHidden(),
		Effect:           cm.SpellEffectNone, // TODO
		WingEffect:       uint8(p.LooksWings),
		Extra:            false,                  // TODO
		MountType:        -1,                     // TODO
		RidingMount:      false,                  // TODO
		Fishing:          false,                  // TODO
		TransformType:    -1,                     // TODO
		ElementOrbEffect: 0,                      // TODO
		ElementOrbLvl:    0,                      // TODO
		ElementOrbMax:    200,                    // TODO
		Buffs:            make([]cm.BuffType, 0), // TODO
		LevelEffects:     cm.LevelEffectsNone,    // TODO
	}
}

func (ServerMessage) ObjectMonster(m *Monster) *server.ObjectMonster {
	return &server.ObjectMonster{
		ObjectID:          m.ID,
		Name:              m.Name,
		NameColor:         m.NameColor.ToInt32(),
		Location:          m.GetPoint(),
		Image:             m.Image,
		Direction:         m.GetDirection(),
		Effect:            uint8(m.Effect),
		AI:                uint8(m.AI),
		Light:             m.Light,
		Dead:              m.IsDead(),
		Skeleton:          m.IsSkeleton(),
		Poison:            m.Poison,
		Hidden:            m.IsHidden(),
		ShockTime:         0,     // TODO
		BindingShotCenter: false, // TODO
		Extra:             false, // TODO
		ExtraByte:         0,     // TODO
	}
}

func (ServerMessage) ObjectGold(i *Item) *server.ObjectGold {
	return &server.ObjectGold{
		ObjectID:  i.GetID(),
		Gold:      uint32(i.Gold),
		LocationX: int32(i.GetPoint().X),
		LocationY: int32(i.GetPoint().Y),
	}
}

func (ServerMessage) ObjectItem(i *Item) *server.ObjectItem {
	return &server.ObjectItem{
		ObjectID:  i.GetID(),
		Name:      i.Name,
		NameColor: i.NameColor.ToInt32(),
		LocationX: int32(i.GetPoint().X),
		LocationY: int32(i.GetPoint().Y),
		Image:     i.GetImage(),
		Grade:     cm.ItemGradeNone, // TODO
	}
}

func (ServerMessage) ObjectNPC(n *NPC) *server.ObjectNPC {
	return &server.ObjectNPC{
		ObjectID:  n.ID,
		Name:      n.Name,
		NameColor: n.NameColor.ToInt32(),
		Image:     uint16(n.Image),
		Color:     0, // TODO
		Location:  n.GetPoint(),
		Direction: n.GetDirection(),
		QuestIDs:  []int32{}, // TODO
	}
}

func (ServerMessage) MapInformation(info *cm.MapInfo) *server.MapInformation {
	mi := new(server.MapInformation)
	mi.FileName = info.Filename
	mi.Title = info.Title
	mi.MiniMap = uint16(info.MiniMap)
	mi.BigMap = uint16(info.BigMap)
	mi.Music = uint16(info.Music)
	mi.Lights = cm.LightSetting(info.Light)
	mi.Lightning = true
	mi.MapDarkLight = 0
	return mi
}

func (ServerMessage) StartGame(result, resolution int) *server.StartGame {
	/*
	 * 0: Disabled.
	 * 1: Not logged in
	 * 2: Character not found.
	 * 3: Start Game Error
	 * 4: Success
	 * */
	sg := new(server.StartGame)
	sg.Result = uint8(result)
	sg.Resolution = int32(resolution)
	return sg
}

func (ServerMessage) UserInformation(p *Player) *server.UserInformation {
	ui := new(server.UserInformation)
	ui.ObjectID = p.GetID()
	ui.RealID = p.GetID()
	ui.Name = p.Name
	ui.GuildName = p.GuildName
	ui.GuildRank = p.GuildRankName
	ui.NameColor = cm.Color{R: 255, G: 255, B: 255}.ToInt32()
	ui.Class = p.Class
	ui.Gender = p.Gender
	ui.Level = p.Level
	ui.Location = p.CurrentLocation
	ui.Direction = p.Direction
	ui.Hair = p.Hair
	ui.HP = p.HP
	ui.MP = p.MP
	ui.Experience = p.Experience
	ui.MaxExperience = int64(data.ExpList[p.Level-1])
	ui.LevelEffect = cm.LevelEffectsNone // TODO
	ui.Gold = uint32(p.Gold)
	ui.Credit = 100 // TODO
	ui.Inventory = p.Inventory.Items
	ui.Equipment = p.Equipment.Items
	ui.QuestInventory = p.QuestInventory.Items
	ui.HasExpandedStorage = false    // TODO
	ui.ExpandedStorageExpiryTime = 0 // TODO
	ui.ClientMagics = p.GetClientMagics()
	return ui
}

func (p *Player) GetClientMagics() []*cm.ClientMagic {
	res := make([]*cm.ClientMagic, 0)
	for i := range p.Magics {
		userMagic := p.Magics[i]
		res = append(res, userMagic.GetClientMagic(userMagic.Info))
	}
	return res
}

func (ServerMessage) UserLocation(p *Player) *server.UserLocation {
	return &server.UserLocation{
		Location:  p.Point(),
		Direction: p.Direction,
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

func (ServerMessage) ObjectChat(p *Player, message string, chatType cm.ChatType) *server.ObjectChat {
	text := p.Name + ":" + message
	return &server.ObjectChat{
		ObjectID: p.ID,
		Text:     text,
		Type:     chatType,
	}
}

func (ServerMessage) PlayerInspect(p *Player) *server.PlayerInspect {
	return &server.PlayerInspect{
		Name:      p.Name,
		GuildName: p.GuildName,
		GuildRank: p.GuildRankName,
		Equipment: p.Equipment.Items,
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

func (ServerMessage) NewCharacterSuccess(g *Game, AccountID int, name string, class cm.MirClass, gender cm.MirGender) *server.NewCharacterSuccess {
	c := new(cm.Character)
	c.Name = name
	c.Level = 1
	c.Class = class
	c.Gender = gender
	c.Hair = 1

	startPoint := data.RandomStartPoint()
	c.CurrentMapID = startPoint.MapID
	c.CurrentLocationX = startPoint.LocationX
	c.CurrentLocationY = startPoint.LocationY
	c.BindMapID = startPoint.MapID
	c.BindLocationX = startPoint.LocationX
	c.BindLocationY = startPoint.LocationY

	c.Direction = cm.MirDirectionDown
	c.HP = 15
	c.MP = 17
	c.Experience = 0
	c.AttackMode = cm.AttackModeAll
	c.PetMode = cm.PetModeBoth
	adb.Table("character").Create(c)
	adb.Table("character").Where("name = ?", name).Last(c)
	ac := new(cm.AccountCharacter)
	ac.AccountID = AccountID
	ac.CharacterID = int(c.ID)
	adb.Table("account_character").Create(ac)
	res := new(server.NewCharacterSuccess)
	res.CharInfo.Index = uint32(c.ID)
	res.CharInfo.Name = name
	res.CharInfo.Class = class
	res.CharInfo.Gender = gender
	return res
}

func (ServerMessage) LogOutSuccess(characters []server.SelectInfo) *server.LogOutSuccess {
	return &server.LogOutSuccess{Characters: characters}
}

func (ServerMessage) NPCResponse(page []string) *server.NPCResponse {
	return &server.NPCResponse{Page: page}
}

func (m ServerMessage) Object(obj IMapObject) interface{} {
	switch obj := obj.(type) {
	case *Player:
		return m.ObjectPlayer(obj)
	case *Monster:
		return m.ObjectMonster(obj)
	case *NPC:
		return m.ObjectNPC(obj)
	case *Item:
		if obj.UserItem == nil {
			return m.ObjectGold(obj)
		} else {
			return m.ObjectItem(obj)
		}
	default:
		panic("unknown object")
	}
}

func (ServerMessage) GainedItem(ui *cm.UserItem) *server.GainedItem {
	return &server.GainedItem{Item: ui}
}

func (m ServerMessage) GainedGold(gold uint64) *server.GainedGold {
	return &server.GainedGold{Gold: uint32(gold)}
}

func (ServerMessage) PlayerUpdate(p *Player) *server.PlayerUpdate {
	return &server.PlayerUpdate{
		ObjectID:     p.GetID(),
		Light:        p.Light,
		Weapon:       int16(p.LooksWeapon),
		WeaponEffect: int16(p.LooksWeaponEffect),
		Armour:       int16(p.LooksArmour),
		WingEffect:   uint8(p.LooksWings),
	}
}

func (ServerMessage) ObjectAttack(obj IMapObject, spell cm.Spell, level int, typ int) *server.ObjectAttack {
	return &server.ObjectAttack{
		ObjectID:  obj.GetID(),
		LocationX: int32(obj.GetPoint().X),
		LocationY: int32(obj.GetPoint().Y),
		Direction: obj.GetDirection(),
		Spell:     spell,
		Level:     uint8(level),
		Type:      uint8(typ),
	}
}

func (ServerMessage) DamageIndicator(dmg int32, typ cm.DamageType, id uint32) *server.DamageIndicator {
	return &server.DamageIndicator{
		Damage:   dmg,
		Type:     typ,
		ObjectID: id,
	}
}

func (ServerMessage) ObjectStruck(m IMapObject, attackerID uint32) *server.ObjectStruck {
	return &server.ObjectStruck{
		ObjectID:   m.GetID(),
		AttackerID: attackerID,
		LocationX:  int32(m.GetPoint().X),
		LocationY:  int32(m.GetPoint().Y),
		Direction:  m.GetDirection(),
	}
}

func (ServerMessage) ObjectHealth(id uint32, percent, expire uint8) *server.ObjectHealth {
	return &server.ObjectHealth{
		ObjectID: id,
		Percent:  percent,
		Expire:   expire,
	}
}

func (ServerMessage) ObjectDied(id uint32, direction cm.MirDirection, location cm.Point) *server.ObjectDied {
	return &server.ObjectDied{
		ObjectID:  id,
		LocationX: int32(location.X),
		LocationY: int32(location.Y),
		Direction: direction,
		Type:      0,
	}
}

func (ServerMessage) HealthChanged(hp, mp uint16) *server.HealthChanged {
	return &server.HealthChanged{
		HP: hp,
		MP: mp,
	}
}

func (ServerMessage) GainExperience(amount uint32) *server.GainExperience {
	return &server.GainExperience{Amount: amount}
}

func (ServerMessage) LevelChanged(level uint16, experience, maxExperience int64) *server.LevelChanged {
	return &server.LevelChanged{
		Level:         level,
		Experience:    experience,
		MaxExperience: maxExperience,
	}
}

func (ServerMessage) ObjectLeveled(id uint32) *server.ObjectLeveled {
	return &server.ObjectLeveled{ObjectID: id}
}
