package mir

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/yenkeia/mirgo/setting"
	"github.com/yenkeia/mirgo/ut"

	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/proto/server"
)

const (
	LOGIN = iota
	SELECT
	GAME
	DISCONNECTED
)

// Player ...
type Player struct {
	AccountID int
	GameStage int
	Session   *cellnet.Session
	MapObject
	HP                 uint16
	MP                 uint16
	Level              uint16
	Experience         int64
	MaxExperience      int64
	Gold               uint64
	GuildName          string
	GuildRankName      string
	Class              common.MirClass
	Gender             common.MirGender
	Hair               uint8
	Light              uint8
	Inventory          *Bag               // 46
	Equipment          *Bag               // 14
	QuestInventory     *Bag               // 40
	Storage            *Bag               // 80
	Trade              []*common.UserItem // 10
	Refine             []*common.UserItem // 16
	LooksArmour        int
	LooksWings         int
	LooksWeapon        int
	LooksWeaponEffect  int
	SendItemInfo       []*common.ItemInfo
	CurrentBagWeight   int
	MaxHP              uint16
	MaxMP              uint16
	MinAC              uint16 // 物理防御力
	MaxAC              uint16
	MinMAC             uint16 // 魔法防御力
	MaxMAC             uint16
	MinDC              uint16 // 攻击力
	MaxDC              uint16
	MinMC              uint16 // 魔法力
	MaxMC              uint16
	MinSC              uint16 // 道术力
	MaxSC              uint16
	Accuracy           uint8
	Agility            uint8
	CriticalRate       uint8
	CriticalDamage     uint8
	MaxBagWeight       uint16 //Other Stats;
	MaxWearWeight      uint16
	MaxHandWeight      uint16
	ASpeed             int8
	Luck               int8
	LifeOnHit          uint8
	HpDrainRate        uint8
	Reflect            uint8
	MagicResist        uint8
	PoisonResist       uint8
	HealthRecovery     uint8
	SpellRecovery      uint8
	PoisonRecovery     uint8
	Holy               uint8
	Freezing           uint8
	PoisonAttack       uint8
	ExpRateOffset      float32
	ItemDropRateOffset float32
	MineRate           uint8
	GemRate            uint8
	FishRate           uint8
	CraftRate          uint8
	GoldDropRateOffset float32
	AttackBonus        uint8
	Magics             []*common.UserMagic
	ActionList         *ActionList
	PoisonList         *PoisonList
	BuffList           *BuffList
	Health             Health // 状态恢复
	Pets               []IMapObject
	PKPoints           int
	AMode              common.AttackMode
	PMode              common.PetMode
	CallingNPC         *NPC
	CallingNPCPage     string
	Slaying            bool // TODO
	FlamingSword       bool // TODO
	TwinDrakeBlade     bool // TODO
}

type Health struct {
	// 生命药水回复
	HPPotValue    int           // 回复总值
	HPPotPerValue int           // 一次回复多少
	HPPotNextTime *time.Time    // 下次生效时间
	HPPotDuration time.Duration // 两次生效时间间隔
	HPPotTickNum  int           // 总共跳几次
	HPPotTickTime int           // 当前第几跳
	// 魔法药水回复
	MPPotValue    int
	MPPotPerValue int
	MPPotNextTime *time.Time
	MPPotDuration time.Duration
	MPPotTickNum  int
	MPPotTickTime int
	// 角色生命/魔法回复
	HealNextTime *time.Time
	HealDuration time.Duration
}

func (p *Player) GetMap() *Map {
	return p.Map
}

func (p *Player) GetID() uint32 {
	return p.ID
}

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) GetLevel() int {
	return int(p.Level)
}

func (p *Player) GetHP() int {
	return int(p.HP)
}

func (p *Player) GetMaxHP() int {
	return int(p.MaxHP)
}

func (p *Player) Point() common.Point {
	return p.GetPoint()
}

// GetFrontPoint 获取玩家面前的点
func (p *Player) GetFrontPoint() common.Point {
	return p.Point().NextPoint(p.CurrentDirection, 1)
}

func (m *Player) AddPlayerCount(n int) {
	m.PlayerCount += n
	switch m.PlayerCount {
	case 1:
		m.Map.AddActiveObj(m)
	case 0:
		m.Map.DelActiveObj(m)
	}
}

func (p *Player) BroadcastHealthChange() {
	IMapObject_BroadcastHealthChange(p)
}

func (p *Player) BroadcastInfo() {
	p.Broadcast(p.GetInfo())
}

func (p *Player) Spawned() {
	IMapObject_Spawned(p)
}

func (m *Player) GetPlayerCount() int {
	return m.PlayerCount
}

func (p *Player) GetRace() common.ObjectType {
	return common.ObjectTypePlayer
}

func (p *Player) GetPoint() common.Point {
	return p.CurrentLocation
}
func (p *Player) IsBlocking() bool {
	return !p.IsDead() // && !Observer;
}

func (p *Player) GetCell() *Cell {
	return p.Map.GetCell(p.CurrentLocation)
}

func (p *Player) GetDirection() common.MirDirection {
	return p.CurrentDirection
}

// func (p *Player) GetCurrentGrid() *Grid {
// 	return p.Map.AOI.GetGridByPoint(p.Point())
// }

func (p *Player) AttackMode() common.AttackMode {
	return common.AttackModeAll
}

// IsAttackTarget 判断玩家是否是攻击者的攻击对象
func (p *Player) IsAttackTarget(attacker IMapObject) bool {
	// return false
	if attacker == nil {
		return false
	}
	if p.IsDead() {
		return false
	}
	switch attacker.GetRace() {
	case common.ObjectTypePlayer:
	case common.ObjectTypeMonster:
		monster := attacker.(*Monster)
		monsterInfo := data.GetMonsterInfoByName(monster.Name)
		if monsterInfo.AI == 6 || monsterInfo.AI == 58 {
			return p.PKPoints >= 200
		}
		if monster.Master == nil {
			break
		}
		if monster.Master.GetID() == p.GetID() {
			return false
		}
		switch monster.Master.AMode {
		case common.AttackModeAll:
			return true
		case common.AttackModeGroup:
			// return GroupMembers == null || !GroupMembers.Contains(attacker.Master)
		case common.AttackModeGuild:
			return true
		case common.AttackModeEnemyGuild:
			return false
		case common.AttackModePeace:
			return false
		case common.AttackModeRedBrown:
			return p.PKPoints >= 200 //|| Envir.Time < BrownTime
		}
	}
	return true
}

func (p *Player) IsFriendlyTarget(obj IMapObject) bool {
	switch obj.GetRace() {
	case common.ObjectTypePlayer:
		ally := obj.(*Player)
		if ally.GetID() == p.GetID() {
			return true
		}
		switch ally.AMode {
		case common.AttackModeGroup:
			// return GroupMembers != null && GroupMembers.Contains(ally)
		case common.AttackModeRedBrown:
			return p.PKPoints < 200 // &Envir.Time > BrownTime
		case common.AttackModeGuild:
			// return MyGuild != null && MyGuild == ally.MyGuild
		case common.AttackModeEnemyGuild:
			return true
		}
		return true
	case common.ObjectTypeMonster:
		ally := obj.(*Monster)
		if ally.Master == nil {
			return false
		}
		// switch (ally.Master.Race)
		// {
		// 	case ObjectType.Player:
		// 		if (!ally.Master.IsFriendlyTarget(this)) return false;
		// 		break;
		// 	case ObjectType.Monster:
		// 		return false;
		// }
		if !ally.Master.IsFriendlyTarget(ally) {
			return false
		}
		return true
	}
	return true
}

func (p *Player) GetInfo() interface{} {
	res := &server.ObjectPlayer{
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
		Poison:           common.PoisonTypeNone, // TODO
		Dead:             p.IsDead(),
		Hidden:           p.IsHidden(),
		Effect:           common.SpellEffectNone, // TODO
		WingEffect:       uint8(p.LooksWings),
		Extra:            false,                      // TODO
		MountType:        0,                          // TODO
		RidingMount:      false,                      // TODO
		Fishing:          false,                      // TODO
		TransformType:    0,                          // TODO
		ElementOrbEffect: 0,                          // TODO
		ElementOrbLvl:    0,                          // TODO
		ElementOrbMax:    0,                          // TODO
		Buffs:            make([]common.BuffType, 0), // TODO
		LevelEffects:     common.LevelEffectsNone,    // TODO
	}
	return res
}

func (p *Player) GetBaseStats() BaseStats {
	return BaseStats{
		MinAC:    p.MinAC,
		MaxAC:    p.MaxAC,
		MinMAC:   p.MinMAC,
		MaxMAC:   p.MaxMAC,
		MinDC:    p.MinDC,
		MaxDC:    p.MaxDC,
		MinMC:    p.MinMC,
		MaxMC:    p.MaxMC,
		MinSC:    p.MinSC,
		MaxSC:    p.MaxSC,
		Accuracy: p.Accuracy,
		Agility:  p.Agility,
	}
}

// AddBuff ...
func (p *Player) AddBuff(buff *Buff) {
	p.BuffList.List.PushBack(buff)
}

func (p *Player) ApplyPoison(poison *Poison, caster IMapObject) {

}

func (p *Player) NewObjectID() uint32 {
	return env.NewObjectID()
}

func (p *Player) IsDead() bool {
	return false
}

func (p *Player) IsUndead() bool {
	return false
}

func (p *Player) IsHidden() bool {
	return false
}

func (p *Player) CanMove() bool {
	return true
}

func (p *Player) CanWalk() bool {
	return true
}

func (p *Player) CanRun() bool {
	return true
}

func (p *Player) CanAttack() bool {
	return true
}

func (p *Player) CanRegen() bool {
	return true
}

func (p *Player) CanCast() bool {
	return true
}

func (p *Player) CanUseItem(item *common.UserItem) bool {
	return true
}

func (p *Player) Enqueue(msg interface{}) {
	if msg == nil {
		log.Errorln("warning: enqueue nil message")
		return
	}
	(*p.Session).Send(msg)
}

func (p *Player) ReceiveChat(text string, ct common.ChatType) {
	p.Enqueue(&server.Chat{
		Message: text,
		Type:    ct,
	})
}

func (p *Player) BroadcastDamageIndicator(typ common.DamageType, dmg int) {
	msg := ServerMessage{}.DamageIndicator(int32(dmg), typ, p.GetID())
	p.Enqueue(msg)
	p.Broadcast(msg)
}

func (p *Player) Broadcast(msg interface{}) {
	p.Map.BroadcastP(p.CurrentLocation, msg, p)
}

func (p *Player) Process(dt time.Duration) {
	now := time.Now()

	p.ActionList.Execute()

	p.ProcessBuffs()
	p.ProcessPoison()

	ch := &p.Health
	if ch.HPPotValue != 0 && ch.HPPotNextTime.Before(now) {
		p.ChangeHP(ch.HPPotPerValue)
		ch.HPPotTickTime += 1
		if ch.HPPotTickTime >= ch.HPPotTickNum {
			ch.HPPotValue = 0
		} else {
			*ch.HPPotNextTime = now.Add(ch.HPPotDuration)
		}
	}
	if ch.MPPotValue != 0 && ch.MPPotNextTime.Before(now) {
		p.ChangeMP(ch.MPPotPerValue)
		ch.MPPotTickTime += 1
		if ch.MPPotTickTime >= ch.MPPotTickNum {
			ch.MPPotValue = 0
		} else {
			*ch.MPPotNextTime = now.Add(ch.MPPotDuration)
		}
	}
	if ch.HealNextTime.Before(now) {
		*ch.HealNextTime = now.Add(ch.HealDuration)
		p.ChangeHP(int(float32(p.MaxHP)*0.03) + 1)
		p.ChangeMP(int(float32(p.MaxMP)*0.03) + 1)
	}
}

// ProcessBuffs 增益效果
func (p *Player) ProcessBuffs() {
	refresh := false
	now := time.Now()
	for e := p.BuffList.List.Front(); e != nil; e = e.Next() {
		buff := e.Value.(*Buff)
		if now.Before(buff.ExpireTime) || buff.Infinite || buff.Paused {
			continue
			// TODO
		}
	}
	/*
			if (Concentrating && !ConcentrateInterrupted && (ConcentrateInterruptTime != 0))
		   	{
		        //check for reenable
		        if (ConcentrateInterruptTime <= SMain.Envir.Time)
		        {
		            ConcentrateInterruptTime = 0;
		            UpdateConcentration();//Update & send to client
		        }
		   	}
	*/
	if refresh {
		p.RefreshStats()
	}
}

// ProcessPoison 中毒效果
func (p *Player) ProcessPoison() {

}

// SaveData 保存玩家数据
func (p *Player) SaveData() {
	adb.SyncPosition(p)
}

func (p *Player) EnqueueItemInfos() {
	gdb := data
	itemInfos := make([]*common.ItemInfo, 0)
	for _, v := range p.Inventory.Items {
		if v != nil {
			itemID := int(v.ItemID)
			itemInfos = append(itemInfos, gdb.GetItemInfoByID(itemID))
		}
	}
	for _, v := range p.Equipment.Items {
		if v != nil {
			itemID := int(v.ItemID)
			itemInfos = append(itemInfos, gdb.GetItemInfoByID(itemID))
		}
	}
	for _, v := range p.QuestInventory.Items {
		if v != nil {
			itemID := int(v.ItemID)
			itemInfos = append(itemInfos, gdb.GetItemInfoByID(itemID))
		}
	}
	for i := range itemInfos {
		p.EnqueueItemInfo(itemInfos[i].ID)
	}
}

func (p *Player) EnqueueItemInfo(itemID int32) {
	for m := range p.SendItemInfo {
		s := p.SendItemInfo[m]
		if s.ID == itemID {
			return
		}
	}
	item := data.GetItemInfoByID(int(itemID))
	if item == nil {
		return
	}
	p.Enqueue(&server.NewItemInfo{Info: item})
	p.SendItemInfo = append(p.SendItemInfo, item)
}

func (p *Player) EnqueueQuestInfo() {

}

func (p *Player) RefreshStats() {
	p.RefreshLevelStats()
	p.RefreshBagWeight()
	p.RefreshEquipmentStats()
	p.RefreshItemSetStats()
	p.RefreshMirSetStats()
	p.RefreshSkills()
	p.RefreshBuffs()
	p.RefreshStatCaps()
	p.RefreshMountStats()
	p.RefreshGuildBuffs()
}

func (p *Player) RefreshLevelStats() {
	baseStats := setting.BaseStats[p.Class]
	p.Accuracy = uint8(baseStats.StartAccuracy)
	p.Agility = uint8(baseStats.StartAgility)
	p.CriticalRate = uint8(baseStats.StartCriticalRate)
	p.CriticalDamage = uint8(baseStats.StartCriticalDamage)
	if int(p.Level) < len(data.ExpList) {
		p.MaxExperience = int64(data.ExpList[p.Level-1])
	} else {
		p.MaxExperience = 0
	}
	p.MaxHP = uint16(14 + (float32(p.Level)/baseStats.HpGain+baseStats.HpGainRate)*float32(p.Level))
	p.MinAC = 0
	if baseStats.MinAc > 0 {
		p.MinAC = uint16(int(p.Level) / baseStats.MinAc)
	}
	p.MaxAC = 0
	if baseStats.MaxAc > 0 {
		p.MaxAC = uint16(int(p.Level) / baseStats.MaxAc)
	}
	p.MinMAC = 0
	if baseStats.MinMac > 0 {
		p.MinMAC = uint16(int(p.Level) / baseStats.MinMac)
	}
	p.MaxMAC = 0
	if baseStats.MaxMac > 0 {
		p.MaxMAC = uint16(int(p.Level) / baseStats.MaxMac)
	}
	p.MinDC = 0
	if baseStats.MinDc > 0 {
		p.MinDC = uint16(int(p.Level) / baseStats.MinDc)
	}
	p.MaxDC = 0
	if baseStats.MaxDc > 0 {
		p.MaxDC = uint16(int(p.Level) / baseStats.MaxDc)
	}
	p.MinMC = 0
	if baseStats.MinMc > 0 {
		p.MinMC = uint16(int(p.Level) / baseStats.MinMc)
	}
	p.MaxMC = 0
	if baseStats.MaxMc > 0 {
		p.MaxMC = uint16(int(p.Level) / baseStats.MaxMc)
	}
	p.MinSC = 0
	if baseStats.MinSc > 0 {
		p.MinSC = uint16(int(p.Level) / baseStats.MinSc)
	}
	p.MaxSC = 0
	if baseStats.MaxSc > 0 {
		p.MaxSC = uint16(int(p.Level) / baseStats.MaxSc)
	}
	p.CriticalRate = 0
	if baseStats.CritialRateGain > 0 {
		p.CriticalRate = uint8(float32(p.CriticalRate) + (float32(p.Level) / baseStats.CritialRateGain))
	}
	p.CriticalDamage = 0
	if baseStats.CriticalDamageGain > 0 {
		p.CriticalDamage = uint8(float32(p.CriticalDamage) + (float32(p.Level) / baseStats.CriticalDamageGain))
	}
	p.MaxBagWeight = uint16(50.0 + float32(p.Level)/baseStats.BagWeightGain*float32(p.Level))
	p.MaxWearWeight = uint16(15.0 + float32(p.Level)/baseStats.WearWeightGain*float32(p.Level))
	p.MaxHandWeight = uint16(12.0 + float32(p.Level)/baseStats.HandWeightGain*float32(p.Level))
	switch p.Class {
	case common.MirClassWarrior:
		p.MaxHP = uint16(14.0 + (float32(p.Level)/baseStats.HpGain+baseStats.HpGainRate+float32(p.Level)/20.0)*float32(p.Level))
		p.MaxMP = uint16(11.0 + (float32(p.Level) * 3.5) + (float32(p.Level) * baseStats.MpGainRate))
	case common.MirClassWizard:
		p.MaxMP = uint16(13.0 + (float32(p.Level/5.0+2.0) * 2.2 * float32(p.Level)) + (float32(p.Level) * baseStats.MpGainRate))
	case common.MirClassTaoist:
		p.MaxMP = uint16((13 + float32(p.Level)/8.0*2.2*float32(p.Level)) + (float32(p.Level) * baseStats.MpGainRate))
	}
}

func (p *Player) RefreshBagWeight() {
	p.CurrentBagWeight = 0
	for _, ui := range p.Inventory.Items {
		if ui != nil {
			it := data.GetItemInfoByID(int(ui.ItemID))
			p.CurrentBagWeight += int(it.Weight)
		}
	}
}

func (p *Player) RefreshEquipmentStats() {
	for _, temp := range p.Equipment.Items {
		if temp == nil {
			continue
		}

		RealItem := temp.Info

		p.MinAC = ut.Uint16(int(p.MinAC) + int(RealItem.MinAC))
		p.MaxAC = ut.Uint16(int(p.MaxAC) + int(RealItem.MaxAC) + int(temp.AC))
		p.MinMAC = ut.Uint16(int(p.MinMAC) + int(RealItem.MinMAC))
		p.MaxMAC = ut.Uint16(int(p.MaxMAC) + int(RealItem.MaxMAC) + int(temp.MAC))
		p.MinDC = ut.Uint16(int(p.MinDC) + int(RealItem.MinDC))
		p.MaxDC = ut.Uint16(int(p.MaxDC) + int(RealItem.MaxDC) + int(temp.DC))
		p.MinMC = ut.Uint16(int(p.MinMC) + int(RealItem.MinMC))
		p.MaxMC = ut.Uint16(int(p.MaxMC) + int(RealItem.MaxMC) + int(temp.MC))
		p.MinSC = ut.Uint16(int(p.MinSC) + int(RealItem.MinSC))
		p.MaxSC = ut.Uint16(int(p.MaxSC) + int(RealItem.MaxSC) + int(temp.SC))
		p.MaxHP = ut.Uint16(int(p.MaxHP) + int(RealItem.HP) + int(temp.HP))
		p.MaxMP = ut.Uint16(int(p.MaxMP) + int(RealItem.MP) + int(temp.MP))

		p.MaxBagWeight = ut.Uint16(int(p.MaxBagWeight) + int(RealItem.BagWeight))
		p.MaxWearWeight = ut.Uint16(int(p.MaxWearWeight) + int(RealItem.WearWeight))
		p.MaxHandWeight = ut.Uint16(int(p.MaxHandWeight) + int(RealItem.HandWeight))

		p.ASpeed = ut.Int8(int(p.ASpeed) + int(temp.AttackSpeed) + int(RealItem.AttackSpeed))
		p.Luck = ut.Int8(int(p.Luck) + int(temp.Luck) + int(RealItem.Luck))

		p.Accuracy = ut.Uint8(int(p.Accuracy) + int(RealItem.Accuracy) + int(temp.Accuracy))
		p.Agility = ut.Uint8(int(p.Agility) + int(RealItem.Agility) + int(temp.Agility))

		// p.HPrate = ut.Int8(HPrate + RealItem.HPrate)
		// p.MPrate = ut.Int8(MPrate + RealItem.MPrate)
		// p.Acrate = ut.Int8(Acrate + RealItem.MaxAcRate)
		// p.Macrate = ut.Int8(Macrate + RealItem.MaxMacRate)

		p.MagicResist = ut.Uint8(int(p.MagicResist) + int(temp.MagicResist) + int(RealItem.MagicResist))
		p.PoisonResist = ut.Uint8(int(p.PoisonResist) + int(temp.PoisonResist) + int(RealItem.PoisonResist))
		p.HealthRecovery = ut.Uint8(int(p.HealthRecovery) + int(temp.HealthRecovery) + int(RealItem.HealthRecovery))
		p.SpellRecovery = ut.Uint8(int(p.SpellRecovery) + int(temp.ManaRecovery) + int(RealItem.SpellRecovery))
		p.PoisonRecovery = ut.Uint8(int(p.PoisonRecovery) + int(temp.PoisonRecovery) + int(RealItem.PoisonRecovery))
		p.CriticalRate = ut.Uint8(int(p.CriticalRate) + int(temp.CriticalRate) + int(RealItem.CriticalRate))
		p.CriticalDamage = ut.Uint8(int(p.CriticalDamage) + int(temp.CriticalDamage) + int(RealItem.CriticalDamage))
		p.Holy = ut.Uint8(int(p.Holy) + int(RealItem.Holy))
		p.Freezing = ut.Uint8(int(p.Freezing) + int(temp.Freezing) + int(RealItem.Freezing))
		p.PoisonAttack = ut.Uint8(int(p.PoisonAttack) + int(temp.PoisonAttack) + int(RealItem.PoisonAttack))
		p.Reflect = ut.Uint8(int(p.Reflect) + int(RealItem.Reflect))
		p.HpDrainRate = ut.Uint8(int(p.HpDrainRate) + int(RealItem.HpDrainRate))

		switch RealItem.Type {
		case common.ItemTypeArmour:
			p.LooksArmour = int(RealItem.Shape)
			p.LooksWings = int(RealItem.Effect)
		case common.ItemTypeWeapon:
			p.LooksWeapon = int(RealItem.Shape)
			p.LooksWeaponEffect = int(RealItem.Effect)
		}
	}
}

func (p *Player) RefreshItemSetStats() {

}

func (p *Player) RefreshMirSetStats() {

}

func (p *Player) RefreshSkills() {
	// 这些技能只是用来加属性
	for _, magic := range p.Magics {
		switch magic.Spell {
		case common.SpellFencing: // 基本剑术
			// Accuracy = (byte)Math.Min(byte.MaxValue, Accuracy + magic.Level * 3);
			// MaxAC = (ushort)Math.Min(ushort.MaxValue, MaxAC + (magic.Level + 1) * 3);
		case common.SpellFatalSword: // 刺客的技能 忽略
			// Accuracy = (byte)Math.Min(byte.MaxValue, Accuracy + magic.Level);
		case common.SpellSpiritSword: // 精神力战法
			// Accuracy = (byte)Math.Min(byte.MaxValue, Accuracy + magic.Level);
			// MaxDC = (ushort)Math.Min(ushort.MaxValue, MaxDC + MaxSC * (magic.Level + 1) * 0.1F);
		}
	}
}

func (p *Player) RefreshBuffs() {

}

func (p *Player) RefreshStatCaps() {

}

func (p *Player) RefreshMountStats() {

}

func (p *Player) RefreshGuildBuffs() {

}

// GetUserItemByID 获取物品，返回该物品在容器的索引和是否成功
func (p *Player) GetUserItemByID(mirGridType common.MirGridType, id uint64) (index int, item *common.UserItem) {
	var arr []*common.UserItem
	switch mirGridType {
	case common.MirGridTypeInventory:
		arr = p.Inventory.Items
	case common.MirGridTypeEquipment:
		arr = p.Equipment.Items
	case common.MirGridTypeStorage:
		arr = p.Storage.Items
	default:
		panic("error mirGridType")
	}
	for i, v := range arr {
		if v != nil && v.ID == id {
			return i, v
		}
	}
	return -1, nil
}

// ConsumeItem 减少物品数量
func (p *Player) ConsumeItem(userItem *common.UserItem, count int) {
	userItem.Count -= uint32(count)
}

// GainItem 为玩家增加物品，增加成功返回 true
func (p *Player) GainItem(item *common.UserItem) bool {
	item.SoulBoundId = p.GetID()

	info := item.Info

	if info.StackSize > 1 {
		for i, v := range p.Inventory.Items {
			if v == nil || info != v.Info || v.Count > info.StackSize {
				continue
			}
			if item.Count+v.Count <= info.StackSize {
				p.Inventory.SetCount(i, v.Count+item.Count)
				return true
			}

			p.Inventory.SetCount(i, v.Count+item.Count)
			item.Count -= info.StackSize - v.Count
		}
	}

	i, j := 6, 46
	if info.Type == common.ItemTypePotion ||
		info.Type == common.ItemTypeScroll ||
		info.Type == common.ItemTypeScript ||
		info.Type == common.ItemTypeAmulet {
		i = 0
		j = 4
	} else if info.Type == common.ItemTypeAmulet {
		i = 4
		j = 6
	} else {
		i = 6
	}
	for i < j {
		if p.Inventory.Items[i] != nil {
			i++
			continue
		}
		p.Inventory.Set(i, item)
		// p.Inventory.Items[i] = ui
		break
	}

	p.EnqueueItemInfo(item.ItemID)
	p.Enqueue(ServerMessage{}.GainedItem(item))
	p.RefreshBagWeight()
	return true
}

// GainGold 为玩家增加金币
func (p *Player) GainGold(gold uint64) {
	if gold <= 0 {
		return
	}
	p.Gold += gold
	adb.SyncGold(p)
	p.Enqueue(ServerMessage{}.GainedGold(gold))
}

func (p *Player) TakeGold(gold uint64) {
	if uint64(gold) > p.Gold {
		log.Warnf("gold error take=%d,has=%d", gold, p.Gold)
		p.Gold = 0
	} else {
		p.Gold -= uint64(gold)
	}
	adb.SyncGold(p)
	p.Enqueue(&server.LoseGold{Gold: uint32(gold)})
}

func (p *Player) UpdateConcentration() {
	p.Enqueue(ServerMessage{}.SetConcentration(p))
	p.Broadcast(ServerMessage{}.SetObjectConcentration(p))
}

func (p *Player) GetAttackPower(min, max int) int {
	if min < 0 {
		min = 0
	}
	if max < min {
		max = min
	}
	// TODO luck
	return ut.RandomInt(min, max)
}

// TODO
func (p *Player) Attacked(attacker IMapObject, damageFinal int, defenceType common.DefenceType, damageWeapon bool) int {
	return 0
}

// GainExp 为玩家增加经验
func (p *Player) GainExp(amount uint32) {
	p.Experience += int64(amount)
	p.Enqueue(ServerMessage{}.GainExperience(amount))
	if p.Experience < p.MaxExperience {
		return
	}

	// 连续升级
	var exp = p.Experience
	for exp >= p.MaxExperience {
		p.Level++
		exp -= p.MaxExperience
		p.RefreshStats()
	}
	adb.SyncLevel(p)

	p.Experience = exp
	p.LevelUp()
}

// WinExp 玩家获取经验
func (p *Player) WinExp(amount, targetLevel int) {
	var expPoint int
	level := int(p.Level)

	if level < targetLevel+10 { //|| !Settings.ExpMobLevelDifference
		expPoint = amount
	} else {
		expPoint = amount - int(math.Round(math.Max(float64(amount)/15.0, 1.0)*float64(level-(targetLevel+10))))
	}
	if expPoint <= 0 {
		expPoint = 1
	}
	// if (GroupMembers != null)
	p.GainExp(uint32(expPoint))
}

func (p *Player) SetHP(amount uint32) {
	if p.HP == uint16(amount) {
		return
	}

	p.HP = uint16(amount)

	if !p.IsDead() && p.HP == 0 {
		p.Die()
	}

	msg := ServerMessage{}.HealthChanged(p.HP, p.MP)
	p.Enqueue(msg)
	p.BroadcastHealthChange()
}

func (p *Player) SetMP(amount uint32) {
	if p.MP == uint16(amount) {
		return
	}

	p.MP = uint16(amount)
	msg := ServerMessage{}.HealthChanged(p.HP, p.MP)
	p.Enqueue(msg)
	p.BroadcastHealthChange()
}

func (p *Player) ChangeHP(amount int) {
	if amount == 0 || p.IsDead() || p.HP >= p.MaxHP {
		return
	}
	p.SetHP(uint32(int(p.HP) + amount))
}

func (p *Player) ChangeMP(amount int) {
	if amount == 0 || p.IsDead() || p.MP >= p.MaxMP {
		return
	}
	p.SetMP(uint32(int(p.MP) + amount))
}

func (p *Player) LevelUp() {
	p.RefreshStats()
	p.SetHP(uint32(p.MaxHP))
	p.SetMP(uint32(p.MaxMP))
	p.Enqueue(ServerMessage{}.LevelChanged(p.Level, p.Experience, p.MaxExperience))
	p.Broadcast(ServerMessage{}.ObjectLeveled(p.GetID()))
}

func (p *Player) Die() {

}

func (p *Player) Teleport(m *Map, pt common.Point) {
	oldMap := p.Map

	{ // MapObject Teleport
		if m == nil || !m.ValidPoint(pt) {
			log.Warnln("Teleport: map not valid", m == nil)
			return
		}
		oldMap.DeleteObject(p)
		p.Broadcast(&server.ObjectTeleportOut{ObjectID: p.GetID(), Type: 0})
		p.Broadcast(&server.ObjectRemove{ObjectID: p.GetID()})

		p.Map = m
		p.CurrentLocation = pt

		// InTrapRock = false;
		m.AddObject(p)
		p.Broadcast(p.GetInfo())

		p.Broadcast(&server.ObjectTeleportIn{ObjectID: p.GetID(), Type: 0})

		p.BroadcastHealthChange()
	}

	p.Enqueue(&server.MapChanged{
		FileName:     m.Info.Filename,
		Title:        m.Info.Title,
		MiniMap:      uint16(m.Info.MiniMap),
		BigMap:       uint16(m.Info.BigMap),
		Lights:       common.LightSetting(m.Info.Light),
		Location:     p.CurrentLocation,
		Direction:    p.CurrentDirection,
		MapDarkLight: uint8(m.Info.MapDarkLight),
		Music:        uint16(m.Info.Music),
	})

	p.EnqueueAreaObjects(nil, p.GetCell())

	p.Enqueue(&server.ObjectTeleportIn{ObjectID: p.GetID(), Type: 0})
	/* TODO
	//Cancel actions
	if (TradePartner != null)
	TradeCancel();

	if (ItemRentalPartner != null)
		CancelItemRental();

	if (RidingMount) RefreshMount();
	if (ActiveBlizzard) ActiveBlizzard = false;

	GetObjectsPassive();

	SafeZoneInfo szi = CurrentMap.GetSafeZone(CurrentLocation);

	if (szi != null)
	{
		BindLocation = szi.Location;
		BindMapIndex = CurrentMapIndex;
		InSafeZone = true;
	}
	else
		InSafeZone = false;

	CheckConquest();

	Fishing = false;
	Enqueue(GetFishInfo());

	if (mapChanged)
	{
		CallDefaultNPC(DefaultNPCType.MapEnter, CurrentMap.Info.FileName);

		if (Info.Married != 0)
		{
			CharacterInfo Lover = Envir.GetCharacterInfo(Info.Married);
			PlayerObject player = Envir.GetPlayer(Lover.Name);

			if (player != null) player.GetRelationship(false);
		}
	}

	if (CheckStacked())
	{
		StackingTime = Envir.Time + 1000;
		Stacking = true;
	}

	Report.MapChange("Teleported", oldMap.Info, CurrentMap.Info);
	*/
}

func (p *Player) EnqueueAreaObjects(oldCell, newCell *Cell) {
	if oldCell == nil {
		p.Map.RangeObject(p.CurrentLocation, DataRange, func(o IMapObject) bool {
			if o != p {
				o.AddPlayerCount(1)
				p.Enqueue(ServerMessage{}.Object(o))
			}
			return true
		})
		return
	}

	cells := p.Map.CalcDiff(oldCell.Point, newCell.Point, DataRange)
	for c, isadd := range cells.M {
		if isadd {
			for _, o := range c.objects {
				o.AddPlayerCount(1)
				p.Enqueue(ServerMessage{}.Object(o))
			}
		} else {
			for _, o := range c.objects {
				o.AddPlayerCount(-1)
				p.Enqueue(ServerMessage{}.ObjectRemove(o))
			}
		}

	}
}

func (p *Player) CompleteAttack(args ...interface{}) {
	target := args[0].(IMapObject)
	damage := args[1].(int)
	defence := args[2].(common.DefenceType)
	damageWeapon := args[3].(bool)

	if target == nil || !target.IsAttackTarget(p) { // || target.CurrentMap != CurrentMap || target.Node == nil) {
		return
	}

	if target.Attacked(p, damage, defence, damageWeapon) <= 0 {
		return
	}

	//Level Fencing / SpiritSword
	for _, magic := range p.Magics {
		switch magic.Spell {
		case common.SpellFencing, common.SpellSpiritSword:
			p.LevelMagic(magic)
			break
		}
	}
}

func (p *Player) CompleteMapMovement(args ...interface{})     {}
func (p *Player) CompleteMine(args ...interface{})            {}
func (p *Player) CompleteNPC(args ...interface{})             {}
func (p *Player) CompletePoison(args ...interface{})          {}
func (p *Player) CompleteDamageIndicator(args ...interface{}) {}

func (p *Player) StartGame() {
	p.ReceiveChat("这是一个以学习为目的传奇服务端", common.ChatTypeSystem)
	p.ReceiveChat("如有任何建议、疑问欢迎交流", common.ChatTypeSystem)
	p.ReceiveChat("源码地址 https://github.com/yenkeia/mirgo", common.ChatTypeSystem)
	p.EnqueueItemInfos()
	p.RefreshStats()
	p.EnqueueQuestInfo()
	p.Enqueue(ServerMessage{}.MapInformation(p.Map.Info))
	p.Enqueue(ServerMessage{}.UserInformation(p))
	p.Enqueue(ServerMessage{}.TimeOfDay(common.LightSettingDay))
	// p.EnqueueAreaObjects(nil, p.Map.AOI.GetGridByPoint(p.GetPoint()))
	p.EnqueueAreaObjects(nil, p.GetCell())
	p.Enqueue(ServerMessage{}.NPCResponse([]string{}))
	p.Broadcast(ServerMessage{}.ObjectPlayer(p))
}

func (p *Player) StopGame(reason int) {
	p.Broadcast(ServerMessage{}.ObjectRemove(p))
}

func (p *Player) Turn(direction common.MirDirection) {
	if !p.CanMove() {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	p.CurrentDirection = direction
	p.UpdateInSafeZone()

	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Broadcast(ServerMessage{}.ObjectTurn(p))
}

func (p *Player) Walk(direction common.MirDirection) {
	if !p.CanMove() || !p.CanWalk() {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	n := p.Point().NextPoint(direction, 1)

	if p.CheckMovement(n) {
		return
	}

	ok := p.Map.UpdateObject(p, n)
	if !ok {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	p.CurrentDirection = direction
	p.CurrentLocation = n
	p.UpdateInSafeZone()

	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Broadcast(ServerMessage{}.ObjectWalk(p))
}

func (p *Player) Run(direction common.MirDirection) {
	steps := 2

	var loc common.Point
	for i := 1; i <= steps; i++ {
		loc = p.CurrentLocation.NextPoint(direction, uint32(i))
		if !p.Map.ValidPoint(loc) {
			p.Enqueue(ServerMessage{}.UserLocation(p))
			return
		}
		if !p.Map.CheckDoorOpen(loc) {
			p.Enqueue(ServerMessage{}.UserLocation(p))
			return
		}

		cell := p.Map.GetCell(loc)
		if cell.objects != nil {
			for _, o := range cell.objects {
				switch o.(type) {
				case *NPC:
					// if (!NPC.Visible || !NPC.VisibleLog[Info.Index]) continue
				default:
					if !o.IsBlocking() {
						continue
					}
				}
				p.Enqueue(ServerMessage{}.UserLocation(p))
				return
			}
		}

		if p.CheckMovement(loc) {
			return
		}
	}

	if ok := p.Map.UpdateObject(p, loc); !ok {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	p.CurrentDirection = direction
	p.CurrentLocation = loc
	p.UpdateInSafeZone()

	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Broadcast(ServerMessage{}.ObjectRun(p))
}

func (p *Player) Chat(message string) {
	// private message
	if strings.HasPrefix(message, "/") {
		return
	}
	// group
	if strings.HasPrefix(message, "!!") {
		return
	}

	if strings.HasPrefix(message, "@") {
		msg, err := cmd.Exec(message[1:], p)
		if err != nil {
			p.ReceiveChat(fmt.Sprintf("执行失败(%s)", err), common.ChatTypeSystem)
		}
		if msg != nil && msg.(string) != "" {
			p.ReceiveChat(msg.(string), common.ChatTypeSystem)
		}
		return
	}
	msg := ServerMessage{}.ObjectChat(p, message, common.ChatTypeNormal)
	p.Enqueue(msg)
	p.Broadcast(msg)
}

func (p *Player) MoveItem(mirGridType common.MirGridType, from int32, to int32) {
	msg := &server.MoveItem{
		Grid:    mirGridType,
		From:    from,
		To:      to,
		Success: false,
	}

	var err error

	switch mirGridType {
	case common.MirGridTypeInventory:
		err = p.Inventory.Move(int(from), int(to))
	case common.MirGridTypeStorage:
		err = p.Storage.Move(int(from), int(to))
	case common.MirGridTypeTrade:
		// TODO
	case common.MirGridTypeRefine:
		// TODO
	}

	if err != nil {
		log.Errorln(err)
	} else {
		msg.Success = true
	}

	p.Enqueue(msg)
}

func (p *Player) TakeBackItem(from int32, to int32) {
	msg := &server.TakeBackItem{From: from, To: to, Success: false}

	if p.CallingNPC == nil || !ut.StringEqualFold(p.CallingNPCPage, StorageKey) || !InRange(p.CurrentLocation, p.CallingNPC.GetPoint(), DataRange) {
		p.Enqueue(msg)
		return
	}

	if int(from) > len(p.Storage.Items) || int(to) > len(p.Inventory.Items) {
		p.Enqueue(msg)
		return
	}

	// item := p.Inventory.Get(int(from))
	// if item.Info.Weight+p.CurrentBagWeight > MaxBagWeight {
	// 	p.ReceiveChat("Too heavy to get back.", common.ChatTypeSystem)
	// 	p.Enqueue(p)
	// }
	err := p.Storage.MoveTo(int(from), int(to), p.Inventory)
	if err != nil {
		log.Infoln(err)
		p.Enqueue(msg)
		return
	}

	msg.Success = true
	p.Enqueue(msg)
}

func (p *Player) TakeItem(itemname string, n int) {
	info := data.GetItemInfoByName(itemname)
	if info == nil {
		return
	}

	for i, item := range p.Inventory.Items {
		if item == nil {
			continue
		}
		if item.Info != info {
			continue
		}
		if n > int(item.Count) {
			p.Enqueue(&server.DeleteItem{UniqueID: item.ID, Count: item.Count})
			p.Inventory.Set(i, nil)
			n -= int(item.Count)
			continue
		}

		p.Enqueue(&server.DeleteItem{UniqueID: item.ID, Count: uint32(n)})
		if n == int(item.Count) {
			p.Inventory.Set(i, nil)
		} else {
			p.Inventory.UseCount(i, uint32(n))
		}
		break
	}

	p.RefreshStats()
}

func (p *Player) StoreItem(from int32, to int32) {
	msg := &server.StoreItem{
		From:    from,
		To:      to,
		Success: false,
	}

	if p.CallingNPC == nil || !ut.StringEqualFold(p.CallingNPCPage, StorageKey) || !InRange(p.CurrentLocation, p.CallingNPC.GetPoint(), DataRange) {
		p.Enqueue(msg)
		return
	}

	if int(from) > len(p.Inventory.Items) || int(to) > len(p.Storage.Items) {
		p.Enqueue(msg)
		return
	}

	item := p.Inventory.Get(int(from))
	if ut.HasFlagUint16(item.Info.Bind, common.BindModeDontStore) {
		p.Enqueue(msg)
		return
	}

	// if (temp.RentalInformation != null && temp.RentalInformation.BindingFlags.HasFlag(BindMode.DontStore))
	//         {
	//             Enqueue(p);
	//             return;
	//         }

	err := p.Inventory.MoveTo(int(from), int(to), p.Storage)
	if err != nil {
		log.Infoln(err)
		p.Enqueue(msg)
		return
	}

	msg.Success = true
	p.Enqueue(msg)
}

func (p *Player) DepositRefineItem(from int32, to int32) {

}

func (p *Player) RetrieveRefineItem(from int32, to int32) {

}

func (p *Player) RefineCancel() {
	p.CallingNPC = nil
}

func (p *Player) RefineItem(id uint64) {

}

func (p *Player) CheckRefine(id uint64) {

}

func (p *Player) ReplaceWeddingRing(id uint64) {

}

func (p *Player) DepositTradeItem(from int32, to int32) {

}

func (p *Player) RetrieveTradeItem(from int32, to int32) {

}

func (p *Player) MergeItem(from common.MirGridType, to common.MirGridType, from2 uint64, to2 uint64) {

}

func (p *Player) EquipItem(mirGridType common.MirGridType, id uint64, to int32) {
	var msg = &server.EquipItem{
		Grid:     mirGridType,
		UniqueID: id,
		To:       to,
		Success:  false,
	}

	index, item := p.GetUserItemByID(mirGridType, id)
	if item == nil {
		p.Enqueue(msg)
		return
	}

	var err error

	switch mirGridType {
	case common.MirGridTypeInventory:
		err = p.Inventory.MoveTo(index, int(to), p.Equipment)
	case common.MirGridTypeStorage:
		err = p.Inventory.MoveTo(index, int(to), p.Storage)
	}

	if err != nil {
		p.Enqueue(msg)
		return
	}

	msg.Success = true
	p.RefreshStats()
	p.Enqueue(msg)
	p.UpdateConcentration()
	p.Broadcast(ServerMessage{}.PlayerUpdate(p))
}

func (p *Player) RemoveItem(mirGridType common.MirGridType, id uint64, to int32) {
	msg := &server.RemoveItem{
		Grid:     mirGridType,
		UniqueID: id,
		To:       to,
		Success:  false,
	}

	index, item := p.GetUserItemByID(common.MirGridTypeEquipment, id)
	if item == nil {
		p.Enqueue(msg)
		return
	}

	switch mirGridType {
	case common.MirGridTypeInventory:
		p.Equipment.MoveTo(index, int(msg.To), p.Inventory)
	case common.MirGridTypeStorage:

		if !ut.StringEqualFold(p.CallingNPCPage, StorageKey) {
			p.Enqueue(msg)
			return
		}

		p.Equipment.MoveTo(index, int(msg.To), p.Storage)
	}
	msg.Success = true
	p.RefreshStats()
	p.Enqueue(msg)
	p.UpdateConcentration()
	p.Broadcast(ServerMessage{}.PlayerUpdate(p))
}

func (p *Player) RemoveSlotItem(grid common.MirGridType, id uint64, to int32, to2 common.MirGridType) {

}

func (p *Player) SplitItem(grid common.MirGridType, id uint64, count uint32) {
	msg := &server.SplitItem1{
		Grid:     grid,
		UniqueID: id,
		Count:    count,
		Success:  false,
	}
	var array []*common.UserItem
	switch grid {
	case common.MirGridTypeInventory:
		_, userItem := p.GetUserItemByID(common.MirGridTypeInventory, id)
		if userItem == nil {
			p.Enqueue(msg)
			return
		}
		userItem.Count -= count
		itemInfo := data.GetItemInfoByID(int(userItem.ItemID))
		newUserItem := env.NewUserItem(itemInfo)
		newUserItem.Count = count
		msg.Success = true
		p.Enqueue(msg)
		p.Enqueue(&server.SplitItem{Item: newUserItem, Grid: grid})
		a, b := 0, 6
		if itemInfo.Type == common.ItemTypePotion || itemInfo.Type == common.ItemTypeScroll { // 药水 卷轴
			a = 0
			b = 4
		} else if itemInfo.Type == common.ItemTypeAmulet { // 护身符
			a = 4
			b = 6
		} else {
			a = 6
			b = len(array)
		}
		for i := a; i < b; i++ {
			if array[i] != nil {
				continue
			}
			array[i] = newUserItem
			p.RefreshBagWeight()
			return
		}
	case common.MirGridTypeStorage:
		// TODO
		p.Enqueue(msg)
	default:
		p.Enqueue(msg)
	}
}

func (p *Player) UseItemScroll(item *common.UserItem) bool {
	switch item.Info.Shape {
	case 0: //DE
		// if (!p.TeleportEscape(20)) {
		// 	return false
		// }
	case 1: //TT
		// if (!p.Teleport(env.GetMap(BindMapIndex), BindLocation)) {
		// 	return false
		// }
	case 2: //RT
		// if (!TeleportRandom(200, item.Info.Durability)) {
		// 	Enqueue(p);
		// }
	case 3: //BenedictionOil
		// if (!TryLuckWeapon()) {
		// 	Enqueue(p);
		// }
		/*
			case 4: //RepairOil
				temp = Info.Equipment[(int)EquipmentSlot.Weapon];
				if (temp == null || temp.MaxDura == temp.CurrentDura) {
					Enqueue(p);
					return;
				}
				if (temp.Info.Bind.HasFlag(BindMode.DontRepair)) {
					Enqueue(p);
					return;
				}
				temp.MaxDura = (ushort)Math.Max(0, temp.MaxDura - Math.Min(5000, temp.MaxDura - temp.CurrentDura) / 30);

				temp.CurrentDura = (ushort)Math.Min(temp.MaxDura, temp.CurrentDura + 5000);
				temp.DuraChanged = false;

				ReceiveChat("Your weapon has been partially repaired", ChatType.Hint);
				Enqueue(new S.ItemRepaired { UniqueID = temp.UniqueID, MaxDura = temp.MaxDura, CurrentDura = temp.CurrentDura });
			case 5: //WarGodOil
				temp = Info.Equipment[(int)EquipmentSlot.Weapon];
				if (temp == null || temp.MaxDura == temp.CurrentDura) {
					Enqueue(p);
					return;
				}
				if (temp.Info.Bind.HasFlag(BindMode.DontRepair) || (temp.Info.Bind.HasFlag(BindMode.NoSRepair))) {
					Enqueue(p);
					return;
				}
				temp.CurrentDura = temp.MaxDura;
				temp.DuraChanged = false;

				ReceiveChat("Your weapon has been completely repaired", ChatType.Hint);
				Enqueue(new S.ItemRepaired { UniqueID = temp.UniqueID, MaxDura = temp.MaxDura, CurrentDura = temp.CurrentDura });
			case 6: //ResurrectionScroll
				if (CurrentMap.Info.NoReincarnation) {
					ReceiveChat(string.Format("Cannot use on this map"), ChatType.System);
					Enqueue(p);
					return;
				}
				if (Dead) {
					MP = MaxMP;
					Revive(MaxHealth, true);
				}
			case 7: //CreditScroll
				if (item.Info.Price > 0)
				{
					GainCredit(item.Info.Price);
					ReceiveChat(String.Format("{0} Credits have been added to your Account", item.Info.Price), ChatType.Hint);
				}
			case 8: //MapShoutScroll
				HasMapShout = true;
				ReceiveChat("You have been given one free shout across your current map", ChatType.Hint);
			case 9://ServerShoutScroll
				HasServerShout = true;
				ReceiveChat("You have been given one free shout across the server", ChatType.Hint);
			case 10://GuildSkillScroll
				MyGuild.NewBuff(item.Info.Effect, false);
			case 11://HomeTeleport
				if (MyGuild != null && MyGuild.Conquest != null && !MyGuild.Conquest.WarIsOn && MyGuild.Conquest.PalaceMap != null && !TeleportRandom(200, 0, MyGuild.Conquest.PalaceMap)) {
					Enqueue(p);
					return;
				}
			case 12://LotteryTicket
				if (Envir.Random.Next(item.Info.Effect * 32) == 1){ // 1st prize : 1,000,000
					ReceiveChat("You won 1st Prize! Received 1,000,000 gold", ChatType.Hint);
					GainGold(1000000);
				} else if (Envir.Random.Next(item.Info.Effect * 16) == 1) { // 2nd prize : 200,000
					ReceiveChat("You won 2nd Prize! Received 200,000 gold", ChatType.Hint);
					GainGold(200000);
				} else if (Envir.Random.Next(item.Info.Effect * 8) == 1)  {// 3rd prize : 100,000
					ReceiveChat("You won 3rd Prize! Received 100,000 gold", ChatType.Hint);
					GainGold(100000);
				} else if (Envir.Random.Next(item.Info.Effect * 4) == 1) {// 4th prize : 10,000
					ReceiveChat("You won 4th Prize! Received 10,000 gold", ChatType.Hint);
					GainGold(10000);
				} else if (Envir.Random.Next(item.Info.Effect * 2) == 1) { // 5th prize : 1,000
					ReceiveChat("You won 5th Prize! Received 1,000 gold", ChatType.Hint);
					GainGold(1000);
				} else if (Envir.Random.Next(item.Info.Effect) == 1)  {// 6th prize 500
					ReceiveChat("You won 6th Prize! Received 500 gold", ChatType.Hint);
					GainGold(500);
				} else {
					ReceiveChat("You haven't won anything.", ChatType.Hint);
				}
		*/
	}

	return true
}

func (p *Player) UserItemPotion(item *common.UserItem) bool {
	info := item.Info
	switch info.Shape {
	case 0: // NormalPotion 一般药水
		ph := &p.Health
		if info.HP > 0 {
			ph.HPPotValue = int(info.HP)                         // 回复总值
			ph.HPPotPerValue = int(info.HP / 3)                  // 一次回复多少
			*ph.HPPotNextTime = time.Now().Add(ph.HPPotDuration) // 下次生效时间
			ph.HPPotTickNum = 3                                  // 总共跳几次
			ph.HPPotTickTime = 0                                 // 当前第几跳
		}
		if info.MP > 0 {
			ph.MPPotValue = int(info.MP)
			ph.MPPotPerValue = int(info.MP / 3)
			*ph.MPPotNextTime = time.Now().Add(ph.MPPotDuration)
			ph.MPPotTickNum = 3
			ph.MPPotTickTime = 0
		}
	case 1: // SunPotion 太阳水
		p.ChangeHP(int(info.HP))
		p.ChangeMP(int(info.MP))
	case 2: // TODO MysteryWater
	case 3: // TODO Buff
	case 4: // TODO Exp 经验
	}

	return true
}

func (p *Player) UseItem(id uint64) {
	msg := &server.UseItem{UniqueID: id, Success: false}
	if p.IsDead() {
		p.Enqueue(msg)
		return
	}
	index, item := p.GetUserItemByID(common.MirGridTypeInventory, id)

	if item == nil || item.ID == 0 || !p.CanUseItem(item) {
		p.Enqueue(msg)
		return
	}
	info := item.Info

	switch info.Type {
	case common.ItemTypePotion:
		msg.Success = p.UserItemPotion(item)
	case common.ItemTypeScroll:
		msg.Success = p.UseItemScroll(item)
	case common.ItemTypeBook:
		magic := &common.UserMagic{}
		magic.Spell = common.Spell(info.Shape)
		magic.Info = data.GetMagicInfoBySpell(magic.Spell)

		if magic.Info != nil {
			p.Magics = append(p.Magics, magic)
			p.Enqueue(magic.Info)
			p.RefreshStats()
			msg.Success = true
		}

	case common.ItemTypeScript:
		p.CallDefaultNPC(DefaultNPCTypeUseItem, info.Shape)
		msg.Success = true
	case common.ItemTypeFood:
	case common.ItemTypePets:
	case common.ItemTypeTransform: //Transforms
	}

	if msg.Success {
		if item.Count > 1 {
			p.Inventory.UseCount(index, 1)
		} else {
			p.Inventory.Set(index, nil)
		}

		p.RefreshBagWeight()
	}

	p.Enqueue(msg)
}

func (p *Player) CallDefaultNPC(calltype DefaultNPCType, args ...interface{}) {
	var key string

	switch calltype {
	case DefaultNPCTypeUseItem:
		key = fmt.Sprintf("UseItem(%v)", args[0])
	}

	key = fmt.Sprintf("[@_%s]", key)

	p.ActionList.PushAction(DelayedTypeNPC, func() {
		p.CallNPC1(env.DefaultNPC, key)
	})

	p.Enqueue(&server.NPCUpdate{NPCID: env.DefaultNPC.GetID()})
}

func (p *Player) DropItem(id uint64, count uint32) {
	msg := &server.DropItem{
		UniqueID: id,
		Count:    count,
		Success:  false,
	}
	index, userItem := p.GetUserItemByID(common.MirGridTypeInventory, id)
	if userItem == nil || userItem.ID == 0 {
		p.Enqueue(msg)
		return
	}
	obj := NewItem(p, userItem)
	if dropMsg, ok := obj.Drop(p.GetPoint(), 1); !ok {
		p.ReceiveChat(dropMsg, common.ChatTypeSystem)
		return
	}
	if count >= userItem.Count {
		p.Inventory.Set(index, nil)
		// p.Inventory[index] = nil
	} else {
		p.Inventory.UseCount(index, count)
		// p.Inventory[index].Count -= count
	}
	p.RefreshBagWeight()
	msg.Success = true
	p.Enqueue(msg)
}

func (p *Player) DropGold(gold uint64) {
	if p.Gold < gold {
		return
	}
	obj := NewGold(p, gold)
	if dropMsg, ok := obj.Drop(p.GetPoint(), 3); !ok {
		p.ReceiveChat(dropMsg, common.ChatTypeSystem)
		return
	}

	p.TakeGold(gold)
}

func (p *Player) PickUp() {
	if p.IsDead() {
		return
	}
	c := p.GetCell()
	if c == nil {
		return
	}
	items := make([]*Item, 0)
	for _, o := range c.objects {
		if item, ok := o.(*Item); ok {
			if item.UserItem == nil {
				p.GainGold(item.Gold)
				items = append(items, item)
			} else {
				if p.GainItem(item.UserItem) {
					items = append(items, item)
				}
			}
		}
	}
	for i := range items {
		o := items[i]
		p.Map.DeleteObject(o)
		o.Broadcast(ServerMessage{}.ObjectRemove(o))
	}
}

func (p *Player) Inspect(id uint32) {
	o := env.GetPlayer(id)
	for i := range o.Equipment.Items {
		item := data.GetItemInfoByID(int(o.Equipment.Items[i].ItemID))
		if item != nil {
			p.EnqueueItemInfo(item.ID)
		}
	}
	p.Enqueue(ServerMessage{}.PlayerInspect(o))
}

func (p *Player) ChangeAMode(mode common.AttackMode) {

}

func (p *Player) ChangePMode(mode common.AttackMode) {

}

func (p *Player) ChangeTrade(trade bool) {

}

func (p *Player) Attack(direction common.MirDirection, spell common.Spell) {
	if !p.CanAttack() {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	level := 0
	switch spell {
	case common.SpellSlaying:
		if !p.Slaying {
			spell = common.SpellNone
		} else {
			magic := p.GetMagic(common.SpellSlaying)
			level = magic.Level
		}
		p.Slaying = false
	case common.SpellDoubleSlash:
		magic := p.GetMagic(spell)
		if magic == nil || magic.Info.BaseCost+(magic.Level*magic.Info.LevelCost) > int(p.MP) {
			spell = common.SpellNone
			break
		}
		level = magic.Level
		p.ChangeMP(-(magic.Info.BaseCost + magic.Level*magic.Info.LevelCost))
	case common.SpellThrusting, common.SpellFlamingSword:
		magic := p.GetMagic(spell)
		if (magic == nil) || (!p.FlamingSword && (spell == common.SpellFlamingSword)) {
			spell = common.SpellNone
			break
		}
		level = magic.Level
	case common.SpellHalfMoon, common.SpellCrossHalfMoon:
		magic := p.GetMagic(spell)
		if magic == nil || magic.Info.BaseCost+(magic.Level*magic.Info.LevelCost) > int(p.MP) {
			spell = common.SpellNone
			break
		}
		level = magic.Level
		p.ChangeMP(-(magic.Info.BaseCost + magic.Level*magic.Info.LevelCost))
	case common.SpellTwinDrakeBlade:
		magic := p.GetMagic(spell)
		if !p.TwinDrakeBlade || magic == nil || magic.Info.BaseCost+magic.Level*magic.Info.LevelCost > int(p.MP) {
			spell = common.SpellNone
			break
		}
		level = magic.Level
		p.ChangeMP(-(magic.Info.BaseCost + magic.Level*magic.Info.LevelCost))
	default:
		spell = common.SpellNone
	}
	if !p.Slaying {
		magic := p.GetMagic(common.SpellSlaying)
		if magic != nil && ut.RandomNext(12) <= magic.Level {
			p.Slaying = true
			p.Enqueue(&server.SpellToggle{Spell: common.SpellSlaying, CanUse: p.Slaying})
		}
	}
	_ = level // TODO
	p.CurrentDirection = direction
	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Broadcast(ServerMessage{}.ObjectAttack(p, spell, 0, 0))
	target := p.GetPoint().NextPoint(p.GetDirection(), 1)
	damageBase := p.GetAttackPower(int(p.MinDC), int(p.MaxDC)) // = the original damage from your gear (+ bonus from moonlight and darkbody)
	cell := p.Map.GetCell(target)
	if !cell.CanWalk() {
		return
	}
	for _, o := range cell.objects {
		if o.GetRace() != common.ObjectTypePlayer && o.GetRace() != common.ObjectTypeMonster {
			continue
		}
		if !o.IsAttackTarget(p) {
			continue
		}
		// if (ob.Undead)
		// {
		// 	damageBase = Math.Min(int.MaxValue, damageBase + Holy);
		// 	damageFinal = damageBase;//incase we're not using skills
		// }
		// #region FatalSword	// TODO
		// #region MPEater		// TODO
		// #region Hemorrhage	// TODO
		defence := common.DefenceTypeACAgility
		damageFinal := damageBase
		switch spell {
		case common.SpellSlaying: // 攻杀剑术
			magic := p.GetMagic(common.SpellSlaying)
			damageFinal = magic.GetDamage(damageBase)
			p.LevelMagic(magic)
		// case common.SpellDoubleSlash:
		case common.SpellThrusting: // 刺杀剑术
			magic := p.GetMagic(common.SpellThrusting)
			p.LevelMagic(magic)
		case common.SpellHalfMoon: // 半月弯刀
			magic := p.GetMagic(common.SpellHalfMoon)
			p.LevelMagic(magic)
		case common.SpellCrossHalfMoon: // 圆月弯刀
			magic := p.GetMagic(common.SpellCrossHalfMoon)
			p.LevelMagic(magic)
		case common.SpellTwinDrakeBlade: // 双龙斩
			magic := p.GetMagic(common.SpellTwinDrakeBlade)
			damageFinal = magic.GetDamage(damageBase)
			p.TwinDrakeBlade = false
			//   action = new DelayedAction(DelayedType.Damage, Envir.Time + 400, ob, damageFinal, DefenceType.Agility, false);
			p.ActionList.PushDelayAction(DelayedTypeDamage, 400, func() { p.CompleteAttack(o, damageFinal, common.DefenceTypeAgility, false) })
			p.LevelMagic(magic)
			// TODO
			//   if ((((ob.Race != ObjectType.Player) || Settings.PvpCanResistPoison) && (Envir.Random.Next(Settings.PoisonAttackWeight) >= ob.PoisonResist)) && (ob.Level < Level + 10 && Envir.Random.Next(ob.Race == ObjectType.Player ? 40 : 20) <= magic.Level + 1))
			//   {
			//       ob.ApplyPoison(new Poison { PType = PoisonType.Stun, Duration = ob.Race == ObjectType.Player ? 2 : 2 + magic.Level, TickSpeed = 1000 }, this);
			//       ob.Broadcast(new S.ObjectEffect { ObjectID = ob.ObjectID, Effect = SpellEffect.TwinDrakeBlade });
			//   }
		case common.SpellFlamingSword: // 烈火剑法
			magic := p.GetMagic(common.SpellFlamingSword)
			damageFinal = magic.GetDamage(damageBase)
			p.FlamingSword = false
			defence = common.DefenceTypeAC
			p.LevelMagic(magic)
		}
		p.ActionList.PushDelayAction(DelayedTypeDamage, 300, func() { p.CompleteAttack(o, damageFinal, defence, true) })
	}
}

func (p *Player) RangeAttack(direction common.MirDirection, location common.Point, id uint32) {

}

func (p *Player) Harvest(direction common.MirDirection) {

}

func (p *Player) CallNPC(id uint32, key string) {

	var npc *NPC

	if id == env.DefaultNPC.GetID() {
		npc = env.DefaultNPC
	} else {
		npc = p.Map.GetNPC(id)
	}

	if npc == nil {
		log.Warnf("NPC 不存在: %d %s\n", id, key)
		return
	}
	p.CallNPC1(npc, key)
}

func (p *Player) CallNPC1(npc *NPC, key string) {

	say, err := npc.CallScript(p, key)
	if err != nil {
		log.Warnf("NPC 脚本执行失败: %d %s %s\n", npc.GetID(), key, err.Error())
	}

	p.CallingNPC = npc
	p.CallingNPCPage = key

	p.Enqueue(ServerMessage{}.NPCResponse(replaceTemplates(npc, p, say)))

	// ProcessSpecial
	key = strings.ToUpper(key)
	fmt.Println("call npc", npc.Name, key)

	switch key {
	case BuyKey:
		sendNpcGoods(p, npc)
	case SellKey:
		p.Enqueue(&server.NPCSell{})
	case BuySellKey:
		sendNpcGoods(p, npc)
		p.Enqueue(&server.NPCSell{})
	case StorageKey:
		sendStorage(p, npc)
		p.Enqueue(&server.NPCStorage{})
	case BuyBackKey:
		sendBuyBackGoods(p, npc, true)

	default:
		// TODO
	}
}

func sendBuyBackGoods(p *Player, npc *NPC, syncItem bool) {
	goods := npc.GetPlayerBuyBack(p)

	if syncItem {
		for _, item := range goods {
			p.EnqueueItemInfo(item.ItemID)
		}
	}

	p.Enqueue(&server.NPCGoods{Goods: goods, Rate: 1})
}

func sendStorage(p *Player, npc *NPC) {
	// if (Connection.StorageSent) return;
	// Connection.StorageSent = true;

	for _, item := range p.Storage.Items {
		if item != nil {
			p.EnqueueItemInfo(item.ItemID)
		}
	}

	p.Enqueue(&server.UserStorage{Storage: p.Storage.Items})
}

func sendNpcGoods(p *Player, npc *NPC) {

	goods := npc.Goods

	for _, item := range npc.Goods {
		p.EnqueueItemInfo(item.ItemID)
	}

	if len(goods) != 0 {
		p.Enqueue(&server.NPCGoods{Goods: goods, Rate: 1.0, Type: common.PanelTypeBuy})
		return
	}
}

func (p *Player) TalkMonsterNPC(id uint32) {

}

func (p *Player) BuyItem(index uint64, count uint32, panelType common.PanelType) {
	if p.IsDead() {
		return
	}
	if !ut.StringEqualFold(p.CallingNPCPage, BuySellKey, BuyKey, BuyBackKey, BuyUsedKey, PearlBuyKey) {
		return
	}

	npc := p.CallingNPC
	if npc == nil {
		return
	}

	npc.Buy(p, index, count)
}

func (p *Player) CraftItem(index uint64, count uint32, slots []int) {
	if p.IsDead() {
		return
	}
	if p.CallingNPCPage == "" {
		return
	}

	p.CallingNPC.Craft(p, index, count, slots)
}

func (p *Player) SellItem(id uint64, count uint32) {
	msg := &server.SellItem{UniqueID: id, Count: count}
	if p.IsDead() || count == 0 {
		p.Enqueue(msg)
		return
	}

	if !ut.StringEqualFold(p.CallingNPCPage, BuySellKey, SellKey) {
		p.Enqueue(msg)
		return
	}

	var index = -1
	var temp *common.UserItem
	for i, v := range p.Inventory.Items {
		if v == nil || v.ID != id {
			continue
		}

		temp = v
		index = i
		break
	}

	if temp == nil || index == -1 || count > temp.Count {
		p.Enqueue(msg)
		return
	}

	if ut.HasFlagUint16(temp.Info.Bind, common.BindModeDontSell) {
		p.Enqueue(msg)
		return
	}
	// if (temp.RentalInformation != null && temp.RentalInformation.BindingFlags.HasFlag(BindMode.DontSell))
	// {
	// 	Enqueue(p);
	// 	return;
	// }

	if !p.CallingNPC.HasType(temp.Info.Type) {
		p.ReceiveChat("You cannot sell this item here.", common.ChatTypeSystem)
		p.Enqueue(msg)
		return
	}

	if temp.Info.StackSize > 1 && count != temp.Count {
		item := env.NewUserItem(temp.Info)
		item.Count = count
		if item.Price()/2+p.Gold > uint64(math.MaxUint64) {
			p.Enqueue(msg)
			return
		}

		temp.Count -= count
		temp = item
	} else {
		p.Inventory.Set(index, nil)
		// p.Inventory[index] = nil
	}

	p.CallingNPC.Sell(p, temp)
	msg.Success = true
	p.Enqueue(msg)
	p.GainGold(temp.Price() / 2)
	p.RefreshBagWeight()
}

func (p *Player) RepairItem(id uint64) {

}

func (p *Player) BuyItemBack(id uint64, count uint32) {

}

func (p *Player) SRepairItem(id uint64) {

}

func (p *Player) Magic(spell common.Spell, direction common.MirDirection, targetID uint32, targetLocation common.Point) {
	if !p.CanCast() {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	magic := p.GetMagic(spell)
	if magic == nil {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	info := data.GetMagicInfoByID(magic.MagicID)
	cost := info.BaseCost + info.LevelCost*magic.Level
	if uint16(cost) > p.MP {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	p.CurrentDirection = direction
	p.ChangeMP(-cost)
	target := p.Map.GetObjectInAreaByID(targetID, targetLocation)
	cast, targetID := p.UseMagic(spell, magic, target)
	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Enqueue(&server.Magic{
		Spell:    spell,
		TargetID: targetID,
		TargetX:  int32(targetLocation.X),
		TargetY:  int32(targetLocation.Y),
		Cast:     cast,
		Level:    uint8(magic.Level),
	})
	p.Broadcast(&server.ObjectMagic{
		ObjectID:      p.GetID(),
		LocationX:     int32(p.GetPoint().X),
		LocationY:     int32(p.GetPoint().Y),
		Direction:     p.GetDirection(),
		Spell:         spell,
		TargetID:      targetID,
		TargetX:       int32(targetLocation.X),
		TargetY:       int32(targetLocation.Y),
		Cast:          cast,
		Level:         uint8(magic.Level),
		SelfBroadcast: false,
	})
}

func (p *Player) MagicKey(spell common.Spell, key uint8) {
	clientMagics := p.GetClientMagics()
	for _, cm := range clientMagics {
		// log.Debugln(cm)
		if cm.Spell == spell {
			cm.Key = key
			adb.SyncMagicKey(p, spell, key)
			// log.Debugln("found: ", cm.Spell)
			return
		}
	}
}

func (p *Player) SwitchGroup(group bool) {

}

func (p *Player) AddMember(name string) {

}

func (p *Player) DelMember(name string) {

}

func (p *Player) GroupInvite(invite bool) {

}

func (p *Player) TownRevive() {

}

func (p *Player) SpellToggle(spell common.Spell, use bool) {

}

func (p *Player) ConsignItem(id uint64, price uint32) {

}

func (p *Player) MarketSearch(match string) {

}

func (p *Player) MarketRefresh() {

}

func (p *Player) MarketPage(page int32) {

}

func (p *Player) MarketBuy(id uint64) {

}

func (p *Player) MarketGetBack(id uint64) {

}

func (p *Player) RequestUserName(id uint32) {

}

func (p *Player) RequestChatItem(id uint64) {

}

func (p *Player) EditGuildMember(name string, name2 string, index uint8, changeType uint8) {

}

func (p *Player) CheckMovement(pos common.Point) bool {

	// TODO: 优化效率
	for _, v := range data.MovementInfos {
		if v.SourceMap == p.Map.Info.ID {
			if p.CurrentLocation.EqualXY(v.SourceX, v.SourceY) {
				m := env.GetMap(v.DestinationMap)
				if m == nil {
					log.Infoln("no map id=", v.DestinationMap)
				}
				p.Teleport(m, common.NewPoint(v.DestinationX, v.DestinationY))
				return true
			}
		}
	}

	return false
}

func (p *Player) OpenDoor(doorIndex byte) {
	if p.Map.OpenDoor(doorIndex) {
		p.Enqueue(&server.Opendoor{DoorIndex: doorIndex})
		p.Broadcast(&server.Opendoor{DoorIndex: doorIndex})
	}
}
