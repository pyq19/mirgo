package game

type ServerMessage struct{}

func (ServerMessage) SetConcentration(p *Player) *SM_SetConcentration {
	sc := new(SM_SetConcentration)
	sc.ObjectID = uint32(p.AccountID)
	sc.Enabled = false
	sc.Interrupted = false
	return sc
}

func (ServerMessage) SetObjectConcentration(p *Player) *SM_SetObjectConcentration {
	return &SM_SetObjectConcentration{
		ObjectID:    uint32(p.AccountID),
		Enabled:     false,
		Interrupted: false,
	}
}

func (ServerMessage) ObjectPlayer(p *Player) (res *SM_ObjectPlayer) {
	return &SM_ObjectPlayer{
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
		Poison:           PoisonTypeNone, // TODO
		Dead:             p.IsDead(),
		Hidden:           p.IsHidden(),
		Effect:           SpellEffectNone, // TODO
		WingEffect:       uint8(p.LooksWings),
		Extra:            false,               // TODO
		MountType:        -1,                  // TODO
		RidingMount:      false,               // TODO
		Fishing:          false,               // TODO
		TransformType:    -1,                  // TODO
		ElementOrbEffect: 0,                   // TODO
		ElementOrbLvl:    0,                   // TODO
		ElementOrbMax:    200,                 // TODO
		Buffs:            make([]BuffType, 0), // TODO
		LevelEffects:     LevelEffectsNone,    // TODO
	}
}

func (ServerMessage) ObjectMonster(m *Monster) *SM_ObjectMonster {
	return &SM_ObjectMonster{
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

func (ServerMessage) ObjectGold(i *Item) *SM_ObjectGold {
	return &SM_ObjectGold{
		ObjectID:  i.GetID(),
		Gold:      uint32(i.Gold),
		LocationX: int32(i.GetPoint().X),
		LocationY: int32(i.GetPoint().Y),
	}
}

func (ServerMessage) ObjectItem(i *Item) *SM_ObjectItem {
	return &SM_ObjectItem{
		ObjectID:  i.GetID(),
		Name:      i.Name,
		NameColor: i.NameColor.ToInt32(),
		LocationX: int32(i.GetPoint().X),
		LocationY: int32(i.GetPoint().Y),
		Image:     i.GetImage(),
		Grade:     ItemGradeNone, // TODO
	}
}

func (ServerMessage) ObjectNPC(n *NPC) *SM_ObjectNPC {
	return &SM_ObjectNPC{
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

func (ServerMessage) MapInformation(info *MapInfo) *SM_MapInformation {
	mi := new(SM_MapInformation)
	mi.FileName = info.Filename
	mi.Title = info.Title
	mi.MiniMap = uint16(info.MiniMap)
	mi.BigMap = uint16(info.BigMap)
	mi.Music = uint16(info.Music)
	mi.Lights = LightSetting(info.Light)
	mi.Lightning = true
	mi.MapDarkLight = 0
	return mi
}

func (ServerMessage) StartGame(result, resolution int) *SM_StartGame {
	/*
	 * 0: Disabled.
	 * 1: Not logged in
	 * 2: Character not found.
	 * 3: Start Game Error
	 * 4: Success
	 * */
	sg := new(SM_StartGame)
	sg.Result = uint8(result)
	sg.Resolution = int32(resolution)
	return sg
}

func (ServerMessage) UserInformation(p *Player) *SM_UserInformation {
	ui := new(SM_UserInformation)
	ui.ObjectID = p.GetID()
	ui.RealID = p.GetID()
	ui.Name = p.Name
	ui.GuildName = p.GuildName
	ui.GuildRank = p.GuildRankName
	ui.NameColor = Color{R: 255, G: 255, B: 255}.ToInt32()
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
	ui.LevelEffect = LevelEffectsNone // TODO
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

func (p *Player) GetClientMagics() []*ClientMagic {
	res := make([]*ClientMagic, 0)
	for i := range p.Magics {
		userMagic := p.Magics[i]
		res = append(res, userMagic.GetClientMagic(userMagic.Info))
	}
	return res
}

func (ServerMessage) UserLocation(p *Player) *SM_UserLocation {
	return &SM_UserLocation{
		Location:  p.Point(),
		Direction: p.Direction,
	}
}

func (ServerMessage) ObjectTurn(o IMapObject) *SM_ObjectTurn {
	return &SM_ObjectTurn{
		ObjectID:  o.GetID(),
		Location:  o.GetPoint(),
		Direction: o.GetDirection(),
	}
}

func (ServerMessage) ObjectWalk(o IMapObject) *SM_ObjectWalk {
	return &SM_ObjectWalk{
		ObjectID:  o.GetID(),
		Location:  o.GetPoint(),
		Direction: o.GetDirection(),
	}
}

func (ServerMessage) ObjectRun(o IMapObject) *SM_ObjectRun {
	return &SM_ObjectRun{
		ObjectID:  o.GetID(),
		Location:  o.GetPoint(),
		Direction: o.GetDirection(),
	}
}

func (ServerMessage) ObjectRemove(o IMapObject) *SM_ObjectRemove {
	return &SM_ObjectRemove{ObjectID: o.GetID()}
}

func (ServerMessage) ObjectChat(p *Player, message string, chatType ChatType) *SM_ObjectChat {
	text := p.Name + ":" + message
	return &SM_ObjectChat{
		ObjectID: p.ID,
		Text:     text,
		Type:     chatType,
	}
}

func (ServerMessage) PlayerInspect(p *Player) *SM_PlayerInspect {
	return &SM_PlayerInspect{
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

func (ServerMessage) Login(result int) *SM_Login {
	/*
	 * 0: Disabled
	 * 1: Bad AccountID
	 * 2: Bad Password
	 * 3: Account Not Exist
	 * 4: Wrong Password
	 */
	return &SM_Login{Result: uint8(result)}
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
	return &SM_NewCharacter{Result: uint8(result)}
}

func (ServerMessage) NewCharacterSuccess(g *Game, AccountID int, name string, class MirClass, gender MirGender) *SM_NewCharacterSuccess {
	c := new(Character)
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

	c.Direction = MirDirectionDown
	c.HP = 15
	c.MP = 17
	c.Experience = 0
	c.AttackMode = AttackModeAll
	c.PetMode = PetModeBoth
	adb.Table("character").Create(c)
	adb.Table("character").Where("name = ?", name).Last(c)
	ac := new(AccountCharacter)
	ac.AccountID = AccountID
	ac.CharacterID = int(c.ID)
	adb.Table("account_character").Create(ac)
	res := new(SM_NewCharacterSuccess)
	res.CharInfo.Index = uint32(c.ID)
	res.CharInfo.Name = name
	res.CharInfo.Class = class
	res.CharInfo.Gender = gender
	return res
}

func (ServerMessage) LogOutSuccess(characters []SM_SelectInfo) *SM_LogOutSuccess {
	return &SM_LogOutSuccess{Characters: characters}
}

func (ServerMessage) NPCResponse(page []string) *SM_NPCResponse {
	return &SM_NPCResponse{Page: page}
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

func (ServerMessage) GainedItem(ui *UserItem) *SM_GainedItem {
	return &SM_GainedItem{Item: ui}
}

func (m ServerMessage) GainedGold(gold uint64) *SM_GainedGold {
	return &SM_GainedGold{Gold: uint32(gold)}
}

func (ServerMessage) PlayerUpdate(p *Player) *SM_PlayerUpdate {
	return &SM_PlayerUpdate{
		ObjectID:     p.GetID(),
		Light:        p.Light,
		Weapon:       int16(p.LooksWeapon),
		WeaponEffect: int16(p.LooksWeaponEffect),
		Armour:       int16(p.LooksArmour),
		WingEffect:   uint8(p.LooksWings),
	}
}

func (ServerMessage) ObjectAttack(obj IMapObject, spell Spell, level int, typ int) *SM_ObjectAttack {
	return &SM_ObjectAttack{
		ObjectID:  obj.GetID(),
		LocationX: int32(obj.GetPoint().X),
		LocationY: int32(obj.GetPoint().Y),
		Direction: obj.GetDirection(),
		Spell:     spell,
		Level:     uint8(level),
		Type:      uint8(typ),
	}
}

func (ServerMessage) DamageIndicator(dmg int32, typ DamageType, id uint32) *SM_DamageIndicator {
	return &SM_DamageIndicator{
		Damage:   dmg,
		Type:     typ,
		ObjectID: id,
	}
}

func (ServerMessage) ObjectStruck(m IMapObject, attackerID uint32) *SM_ObjectStruck {
	return &SM_ObjectStruck{
		ObjectID:   m.GetID(),
		AttackerID: attackerID,
		LocationX:  int32(m.GetPoint().X),
		LocationY:  int32(m.GetPoint().Y),
		Direction:  m.GetDirection(),
	}
}

func (ServerMessage) ObjectHealth(id uint32, percent, expire uint8) *SM_ObjectHealth {
	return &SM_ObjectHealth{
		ObjectID: id,
		Percent:  percent,
		Expire:   expire,
	}
}

func (ServerMessage) ObjectDied(id uint32, direction MirDirection, location Point) *SM_ObjectDied {
	return &SM_ObjectDied{
		ObjectID:  id,
		LocationX: int32(location.X),
		LocationY: int32(location.Y),
		Direction: direction,
		Type:      0,
	}
}

func (ServerMessage) HealthChanged(hp, mp uint16) *SM_HealthChanged {
	return &SM_HealthChanged{
		HP: hp,
		MP: mp,
	}
}

func (ServerMessage) GainExperience(amount uint32) *SM_GainExperience {
	return &SM_GainExperience{Amount: amount}
}

func (ServerMessage) LevelChanged(level uint16, experience, maxExperience int64) *SM_LevelChanged {
	return &SM_LevelChanged{
		Level:         level,
		Experience:    experience,
		MaxExperience: maxExperience,
	}
}

func (ServerMessage) ObjectLeveled(id uint32) *SM_ObjectLeveled {
	return &SM_ObjectLeveled{ObjectID: id}
}
