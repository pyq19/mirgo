package game

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/davyxu/cellnet"
	"github.com/yenkeia/mirgo/game/cm"
	"github.com/yenkeia/mirgo/game/proto/server"
	"github.com/yenkeia/mirgo/game/util"
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
	HP                   uint16
	MP                   uint16
	Level                uint16
	Experience           int64
	MaxExperience        int64
	Gold                 uint64
	GuildName            string
	GuildRankName        string
	Class                cm.MirClass
	Gender               cm.MirGender
	Hair                 uint8
	Light                uint8
	Inventory            *Bag           // 46
	Equipment            *Bag           // 14
	QuestInventory       *Bag           // 40
	Storage              *Bag           // 80
	Trade                *Bag           // 10	交易框的索引是从上到下的，背包是从左到右
	Refine               []*cm.UserItem // 16	TODO 合成？提炼？
	LooksArmour          int
	LooksWings           int
	LooksWeapon          int
	LooksWeaponEffect    int
	SendItemInfo         []*cm.ItemInfo
	CurrentBagWeight     int
	MaxHP                uint16
	MaxMP                uint16
	MinAC                uint16 // 物理防御力
	MaxAC                uint16
	MinMAC               uint16 // 魔法防御力
	MaxMAC               uint16
	MinDC                uint16 // 攻击力
	MaxDC                uint16
	MinMC                uint16 // 魔法力
	MaxMC                uint16
	MinSC                uint16 // 道术力
	MaxSC                uint16
	Accuracy             uint8
	Agility              uint8
	CriticalRate         uint8
	CriticalDamage       uint8
	MaxBagWeight         uint16 //Other Stats;
	MaxWearWeight        uint16
	MaxHandWeight        uint16
	ASpeed               int8
	Luck                 int8
	LifeOnHit            uint8
	HpDrainRate          uint8
	Reflect              uint8 // TODO
	MagicResist          uint8
	PoisonResist         uint8
	HealthRecovery       uint8
	SpellRecovery        uint8
	PoisonRecovery       uint8
	Holy                 uint8
	Freezing             uint8
	PoisonAttack         uint8
	ExpRateOffset        float32
	ItemDropRateOffset   float32
	MineRate             uint8
	GemRate              uint8
	FishRate             uint8
	CraftRate            uint8
	GoldDropRateOffset   float32
	AttackBonus          uint8
	Magics               []*cm.UserMagic
	ActionList           *ActionList
	PoisonList           *PoisonList
	BuffList             *BuffList
	Health               Health // 状态恢复
	Pets                 []*Monster
	PKPoints             int
	AMode                cm.AttackMode
	PMode                cm.PetMode
	CallingNPC           *NPC
	CallingNPCPage       string
	Slaying              bool        // TODO
	FlamingSword         bool        // TODO
	TwinDrakeBlade       bool        // TODO
	BindMapIndex         int         // 绑定的地图 死亡时复活用
	BindLocation         cm.Point    // 绑定的坐标 死亡时复活用
	MagicShield          bool        // TODO 是否有魔法盾
	MagicShieldLv        int         // TODO 魔法盾等级
	ArmourRate           float32     // 防御
	DamageRate           float32     // 伤害
	StruckTime           time.Time   // 被攻击硬直时间
	AllowGroup           bool        // 是否允许组队
	GroupMembers         *PlayerList // 小队成员
	GroupInvitation      *Player     // 组队邀请人
	AllowTrade           bool        // 是否允许交易
	TradePartner         *Player     // 交易对象
	TradeInvitation      *Player     // 发起交易的玩家
	TradeLocked          bool        // 是否确认交易
	TradeGoldAmount      uint32      // 摆到交易框的金币
	MyGuild              *Guild      // 加入的行会
	MyGuildRank          *cm.Rank    // 所在工会的头衔/职位
	PendingGuildInvite   *Guild      // 被邀请加入的行会
	GuildIndex           int         // 加入的行会Index
	EnableGuildInvite    bool        // 允许加入行会邀请
	GuildNoticeChanged   bool        // 是否有行会公告改变
	GuildCanRequestItems bool        // ??? 不确定干什么的
	GuildMembersChanged  bool        // 是否有行会成员改变
	CanCreateGuild       bool        // 是否可以创建行会
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

func (p *Player) Point() cm.Point {
	return p.GetPoint()
}

// GetFrontPoint 获取玩家面前的点
func (p *Player) GetFrontPoint() cm.Point {
	return p.Point().NextPoint(p.Direction, 1)
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

func (p *Player) GetPercentHealth() uint8 {
	return uint8((float32(p.GetHP()) / float32(p.GetMaxHP()) * 100))
}

func (p *Player) BroadcastHealthChange() {
	IMapObject_BroadcastHealthChange(p)
}

func (p *Player) BroadcastInfo() {
	p.Broadcast(ServerMessage{}.ObjectPlayer(p))
}

func (p *Player) Spawned() {
	IMapObject_Spawned(p)
}

func (m *Player) GetPlayerCount() int {
	return m.PlayerCount
}

func (p *Player) GetRace() cm.ObjectType {
	return cm.ObjectTypePlayer
}

func (p *Player) GetPoint() cm.Point {
	return p.CurrentLocation
}
func (p *Player) IsBlocking() bool {
	return !p.IsDead() // && !Observer;
}

func (p *Player) GetCell() *Cell {
	return p.Map.GetCell(p.CurrentLocation)
}

func (p *Player) GetDirection() cm.MirDirection {
	return p.Direction
}

// func (p *Player) GetCurrentGrid() *Grid {
// 	return p.Map.AOI.GetGridByPoint(p.Point())
// }

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
	case cm.ObjectTypePlayer:
	case cm.ObjectTypeMonster:
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
		case cm.AttackModeAll:
			return true
		case cm.AttackModeGroup:
			// return GroupMembers == null || !GroupMembers.Contains(attacker.Master)
		case cm.AttackModeGuild:
			return true
		case cm.AttackModeEnemyGuild:
			return false
		case cm.AttackModePeace:
			return false
		case cm.AttackModeRedBrown:
			return p.PKPoints >= 200 //|| Envir.Time < BrownTime
		}
	}
	return true
}

func (p *Player) IsFriendlyTarget(obj IMapObject) bool {
	switch obj.GetRace() {
	case cm.ObjectTypePlayer:
		ally := obj.(*Player)
		if ally.GetID() == p.GetID() {
			return true
		}
		switch ally.AMode {
		case cm.AttackModeGroup:
			// return GroupMembers != null && GroupMembers.Contains(ally)
		case cm.AttackModeRedBrown:
			return p.PKPoints < 200 // &Envir.Time > BrownTime
		case cm.AttackModeGuild:
			// return MyGuild != null && MyGuild == ally.MyGuild
		case cm.AttackModeEnemyGuild:
			return true
		}
		return true
	case cm.ObjectTypeMonster:
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
func (p *Player) AddBuff(b *Buff) {
	if p.BuffList.Has(func(temp *Buff) bool { return temp.Infinite && b.Type == temp.Type }) {
		return //cant overwrite infinite buff with regular buff
	}

	p.BuffList.AddBuff(b)

	var caster string
	if b.Caster != nil {
		caster = b.Caster.GetName()
	}

	if b.Values == nil {
		b.Values = []int32{}
	}

	msg := &server.AddBuff{
		Type:     b.Type,
		Caster:   caster,
		Expire:   10000, // TODO
		Values:   b.Values,
		Infinite: b.Infinite,
		ObjectID: p.ID,
		Visible:  b.Visible,
	}

	p.Enqueue(msg)
	if b.Visible {
		p.Broadcast(msg)
	}

	p.RefreshStats()
}

func (p *Player) ApplyPoison(poison *Poison, caster IMapObject) {

}

func (p *Player) IsDead() bool {
	return p.Dead
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

func (p *Player) CanUseItem(item *cm.UserItem) bool {
	return true
}

func (p *Player) Enqueue(msg interface{}) {
	if msg == nil {
		log.Errorln("warning: enqueue nil message")
		return
	}
	(*p.Session).Send(msg)
}

func (p *Player) ReceiveChat(text string, ct cm.ChatType) {
	p.Enqueue(&server.Chat{
		Message: text,
		Type:    ct,
	})
}

func (p *Player) BroadcastDamageIndicator(typ cm.DamageType, dmg int) {
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
		}
		p.BuffList.RemoveBuff(buff.Type)
		p.Enqueue(&server.RemoveBuff{Type: buff.Type, ObjectID: p.GetID()})
		if buff.Visible {
			p.Broadcast(&server.RemoveBuff{Type: buff.Type, ObjectID: p.GetID()})
		}
		// switch buff.Type {
		// case cm.BuffTypeHiding:
		// }
		refresh = true
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
	// 玩家当前位置
	adb.SyncPosition(p)

	// 玩家已经获取的经验
	adb.SyncExperience(p)

	// 玩家 HP MP
	adb.SyncHPMP(p)

	// 玩家 level magic

	// AMode PMode
	adb.SyncAModePMode(p)

	// BindMapIndex BindLocation

	// AllowGroup 组队开关
	adb.SyncAllowGroup(p)
}

func (p *Player) EnqueueItemInfos() {
	gdb := data
	itemInfos := make([]*cm.ItemInfo, 0)
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

func (p *Player) CheckItem(u *cm.UserItem) {
	p.EnqueueItemInfo(u.ItemID)
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
	baseStats := settings.BaseStats[p.Class]
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
	case cm.MirClassWarrior:
		p.MaxHP = uint16(14.0 + (float32(p.Level)/baseStats.HpGain+baseStats.HpGainRate+float32(p.Level)/20.0)*float32(p.Level))
		p.MaxMP = uint16(11.0 + (float32(p.Level) * 3.5) + (float32(p.Level) * baseStats.MpGainRate))
	case cm.MirClassWizard:
		p.MaxMP = uint16(13.0 + (float32(p.Level/5.0+2.0) * 2.2 * float32(p.Level)) + (float32(p.Level) * baseStats.MpGainRate))
	case cm.MirClassTaoist:
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
	oldLooksWeapon := p.LooksWeapon
	oldLooksWeaponEffect := p.LooksWeaponEffect
	oldLooksArmour := p.LooksArmour
	// oldMountType = MountType;
	oldLooksWings := p.LooksWings
	oldLight := p.Light

	p.LooksArmour = 0
	p.LooksWeapon = -1
	p.LooksWeaponEffect = 0
	p.LooksWings = 0
	for _, temp := range p.Equipment.Items {
		if temp == nil {
			continue
		}

		RealItem := data.GetRealItem(temp.Info, p.Level, p.Class, data.ItemInfos)

		p.MinAC = util.Uint16(int(p.MinAC) + int(RealItem.MinAC))
		p.MaxAC = util.Uint16(int(p.MaxAC) + int(RealItem.MaxAC) + int(temp.AC))
		p.MinMAC = util.Uint16(int(p.MinMAC) + int(RealItem.MinMAC))
		p.MaxMAC = util.Uint16(int(p.MaxMAC) + int(RealItem.MaxMAC) + int(temp.MAC))
		p.MinDC = util.Uint16(int(p.MinDC) + int(RealItem.MinDC))
		p.MaxDC = util.Uint16(int(p.MaxDC) + int(RealItem.MaxDC) + int(temp.DC))
		p.MinMC = util.Uint16(int(p.MinMC) + int(RealItem.MinMC))
		p.MaxMC = util.Uint16(int(p.MaxMC) + int(RealItem.MaxMC) + int(temp.MC))
		p.MinSC = util.Uint16(int(p.MinSC) + int(RealItem.MinSC))
		p.MaxSC = util.Uint16(int(p.MaxSC) + int(RealItem.MaxSC) + int(temp.SC))
		p.MaxHP = util.Uint16(int(p.MaxHP) + int(RealItem.HP) + int(temp.HP))
		p.MaxMP = util.Uint16(int(p.MaxMP) + int(RealItem.MP) + int(temp.MP))

		p.MaxBagWeight = util.Uint16(int(p.MaxBagWeight) + int(RealItem.BagWeight))
		p.MaxWearWeight = util.Uint16(int(p.MaxWearWeight) + int(RealItem.WearWeight))
		p.MaxHandWeight = util.Uint16(int(p.MaxHandWeight) + int(RealItem.HandWeight))

		p.ASpeed = util.Int8(int(p.ASpeed) + int(temp.AttackSpeed) + int(RealItem.AttackSpeed))
		p.Luck = util.Int8(int(p.Luck) + int(temp.Luck) + int(RealItem.Luck))

		p.Accuracy = util.Uint8(int(p.Accuracy) + int(RealItem.Accuracy) + int(temp.Accuracy))
		p.Agility = util.Uint8(int(p.Agility) + int(RealItem.Agility) + int(temp.Agility))

		// p.HPrate = util.Int8(HPrate + RealItem.HPrate)
		// p.MPrate = util.Int8(MPrate + RealItem.MPrate)
		// p.Acrate = util.Int8(Acrate + RealItem.MaxAcRate)
		// p.Macrate = util.Int8(Macrate + RealItem.MaxMacRate)

		p.MagicResist = util.Uint8(int(p.MagicResist) + int(temp.MagicResist) + int(RealItem.MagicResist))
		p.PoisonResist = util.Uint8(int(p.PoisonResist) + int(temp.PoisonResist) + int(RealItem.PoisonResist))
		p.HealthRecovery = util.Uint8(int(p.HealthRecovery) + int(temp.HealthRecovery) + int(RealItem.HealthRecovery))
		p.SpellRecovery = util.Uint8(int(p.SpellRecovery) + int(temp.ManaRecovery) + int(RealItem.SpellRecovery))
		p.PoisonRecovery = util.Uint8(int(p.PoisonRecovery) + int(temp.PoisonRecovery) + int(RealItem.PoisonRecovery))
		p.CriticalRate = util.Uint8(int(p.CriticalRate) + int(temp.CriticalRate) + int(RealItem.CriticalRate))
		p.CriticalDamage = util.Uint8(int(p.CriticalDamage) + int(temp.CriticalDamage) + int(RealItem.CriticalDamage))
		p.Holy = util.Uint8(int(p.Holy) + int(RealItem.Holy))
		p.Freezing = util.Uint8(int(p.Freezing) + int(temp.Freezing) + int(RealItem.Freezing))
		p.PoisonAttack = util.Uint8(int(p.PoisonAttack) + int(temp.PoisonAttack) + int(RealItem.PoisonAttack))
		p.Reflect = util.Uint8(int(p.Reflect) + int(RealItem.Reflect))
		p.HpDrainRate = util.Uint8(int(p.HpDrainRate) + int(RealItem.HpDrainRate))

		switch RealItem.Type {
		case cm.ItemTypeArmour:
			p.LooksArmour = int(RealItem.Shape)
			p.LooksWings = int(RealItem.Effect)
		case cm.ItemTypeWeapon:
			p.LooksWeapon = int(RealItem.Shape)
			p.LooksWeaponEffect = int(RealItem.Effect)
		}
	}

	/*
		MaxHP = (ushort)Math.Min(ushort.MaxValue, (((double)HPrate / 100) + 1) * MaxHP);
		MaxMP = (ushort)Math.Min(ushort.MaxValue, (((double)MPrate / 100) + 1) * MaxMP);
		MaxAC = (ushort)Math.Min(ushort.MaxValue, (((double)Acrate / 100) + 1) * MaxAC);
		MaxMAC = (ushort)Math.Min(ushort.MaxValue, (((double)Macrate / 100) + 1) * MaxMAC);

		AddTempSkills(skillsToAdd);
		RemoveTempSkills(skillsToRemove);

		if (HasMuscleRing)
		{
			MaxBagWeight = (ushort)(MaxBagWeight * 2);
			MaxWearWeight = Math.Min(ushort.MaxValue, (ushort)(MaxWearWeight * 2));
			MaxHandWeight = Math.Min(ushort.MaxValue, (ushort)(MaxHandWeight * 2));
		}
	*/

	if (oldLooksArmour != p.LooksArmour) || (oldLooksWeapon != p.LooksWeapon) || (oldLooksWeaponEffect != p.LooksWeaponEffect) || (oldLooksWings != p.LooksWings) || (oldLight != p.Light) {
		p.Broadcast(p.GetUpdateInfo())
		// if (oldLooksWeapon == 49 || oldLooksWeapon == 50) && (p.LooksWeapon != 49 && p.LooksWeapon != 50) {
		// 	p.Enqueue(p.GetFishInfo())
		// }
	}

	/*
		if (Old_MountType != MountType)
		{
			RefreshMount(false);
		}
	*/

	log.Debugf("p.Name: %s\np.LooksArmour = %d\np.LooksWeapon = %d\np.LooksWeaponEffect = %d\np.LooksWings = %d\n", p.Name, p.LooksArmour, p.LooksWeapon, p.LooksWeaponEffect, p.LooksWings)
}

func (p *Player) GetUpdateInfo() *server.PlayerUpdate {
	p.UpdateConcentration()
	return &server.PlayerUpdate{
		ObjectID:     p.GetID(),
		Weapon:       int16(p.LooksWeapon),
		WeaponEffect: int16(p.LooksWeaponEffect),
		Armour:       int16(p.LooksArmour),
		Light:        p.Light,
		WingEffect:   uint8(p.LooksWings),
	}
}

// RefreshItemSetStats ???
func (p *Player) RefreshItemSetStats() {

}

// RefreshMirSetStats ???
func (p *Player) RefreshMirSetStats() {

}

// RefreshSkills 技能加的属性
func (p *Player) RefreshSkills() {
	// 这些技能只是用来加属性
	for _, magic := range p.Magics {
		switch magic.Spell {
		case cm.SpellFencing: // 基本剑术
			p.Accuracy = util.Uint8(int(p.Accuracy) + magic.Level*3)
			p.MaxAC = util.Uint16(int(p.MaxAC) + (magic.Level+1)*3)
		case cm.SpellFatalSword: // 刺客的技能 忽略
		case cm.SpellSpiritSword: // 精神力战法
			p.Accuracy = util.Uint8(int(p.Accuracy) + magic.Level)
			p.MaxAC = util.Uint16(int(p.MaxDC) + int(float32(p.MaxSC)*float32(magic.Level+1)*0.1))
		}
	}
}

// RefreshBuffs 刷新玩家身上的 buff
func (p *Player) RefreshBuffs() {

}

// RefreshStatCaps 刷新各种状态
func (p *Player) RefreshStatCaps() {
	/*
		MagicResist = Math.Min(Settings.MaxMagicResist, MagicResist);
		PoisonResist = Math.Min(Settings.MaxPoisonResist, PoisonResist);
		CriticalRate = Math.Min(Settings.MaxCriticalRate, CriticalRate);
		CriticalDamage = Math.Min(Settings.MaxCriticalDamage, CriticalDamage);
		Freezing = Math.Min(Settings.MaxFreezing, Freezing);
		PoisonAttack = Math.Min(Settings.MaxPoisonAttack, PoisonAttack);
		HealthRecovery = Math.Min(Settings.MaxHealthRegen, HealthRecovery);
		PoisonRecovery = Math.Min(Settings.MaxPoisonRecovery, PoisonRecovery);
		SpellRecovery = Math.Min(Settings.MaxManaRegen, SpellRecovery);
		HpDrainRate = Math.Min((byte)100, HpDrainRate);
	*/
}

// RefreshMountStats 刷新装备嵌套的宝石属性
func (p *Player) RefreshMountStats() {

}

// RefreshGuildBuffs 刷新工会 buff
func (p *Player) RefreshGuildBuffs() {

}

// GetUserItemByID 获取物品，返回该物品在容器的索引和是否成功
func (p *Player) GetUserItemByID(mirGridType cm.MirGridType, id uint64) (index int, item *cm.UserItem) {
	var arr []*cm.UserItem
	switch mirGridType {
	case cm.MirGridTypeInventory:
		arr = p.Inventory.Items
	case cm.MirGridTypeEquipment:
		arr = p.Equipment.Items
	case cm.MirGridTypeStorage:
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
func (p *Player) ConsumeItem(userItem *cm.UserItem, count uint32) {

	bag := p.Equipment
	idx, item := p.GetUserItemByID(cm.MirGridTypeEquipment, userItem.ID)
	if idx == -1 || item != userItem {
		bag = p.Inventory
		idx, item = p.GetUserItemByID(cm.MirGridTypeInventory, userItem.ID)
	}

	if idx == -1 || item != userItem {
		// 没该物品
		return
	}

	p.Enqueue(&server.DeleteItem{UniqueID: item.ID, Count: count})

	bag.UseCount(idx, count)
}

// CanGainItem 是否可以增加物品
func (p *Player) CanGainItem(item *cm.UserItem) bool {
	useWeight := true
	if item.Info.Type == cm.ItemTypeAmulet {
		if p.FreeSpace(p.Inventory) > 0 && (p.CurrentBagWeight+int(item.Info.Weight) <= int(p.MaxBagWeight) || !useWeight) {
			return true
		}
		count := item.Count
		for i := 0; i < p.Inventory.Length(); i++ {
			bagItem := p.Inventory.Get(i)
			if bagItem == nil || bagItem.Info != item.Info {
				continue
			}
			if bagItem.Count+count <= bagItem.Info.StackSize {
				return true
			}
			count -= bagItem.Info.StackSize - bagItem.Count
		}
		return false
	}
	if useWeight && p.CurrentBagWeight+int(item.Info.Weight) > int(p.MaxBagWeight) {
		return false
	}
	if p.FreeSpace(p.Inventory) > 0 {
		return true
	}
	if item.Info.StackSize > 1 {
		count := item.Count
		for i := 0; i < p.Inventory.Length(); i++ {
			bagItem := p.Inventory.Get(i)
			if bagItem.Info != item.Info {
				continue
			}
			if bagItem.Count+count <= bagItem.Info.StackSize {
				return true
			}
			count -= bagItem.Info.StackSize - bagItem.Count
		}
	}
	return false
}

// CanGainItems 是否可以增加很多物品
func (p *Player) CanGainItems(bag *Bag) bool {
	items := bag.Items
	itemCount := bag.ItemCount()
	itemWeight := 0
	stackOffset := 0
	if itemCount < 1 {
		return true
	}
	for i := 0; i < len(items); i++ {
		if items[i] == nil {
			continue
		}
		itemWeight += int(items[i].Info.Weight)
		if items[i].Info.StackSize > 1 {
			count := items[i].Count
			for u := 0; u < p.Inventory.Length(); u++ {
				bagItem := p.Inventory.Get(u)
				if bagItem == nil || bagItem.Info.ID != items[i].Info.ID {
					continue
				}
				if bagItem.Count+count > bagItem.Info.StackSize {
					stackOffset++
				}
				break
			}
		}
	}
	if p.CurrentBagWeight+(itemWeight) > int(p.MaxBagWeight) {
		return false
	}
	if p.FreeSpace(p.Inventory) < itemCount+stackOffset {
		return false
	}
	return true
}

// GainItem 为玩家增加物品，增加成功返回 true
func (p *Player) GainItem(item *cm.UserItem) (res bool) {
	item.SoulBoundId = p.GetID()

	if item.Info.StackSize > 1 {
		for i, v := range p.Inventory.Items {
			if v == nil || item.Info != v.Info || v.Count > item.Info.StackSize {
				continue
			}
			if item.Count+v.Count <= item.Info.StackSize {
				p.Inventory.SetCount(i, v.Count+item.Count)
				p.Enqueue(ServerMessage{}.GainedItem(item))
				return true
			}
			p.Inventory.SetCount(i, v.Count+item.Count)
			item.Count -= item.Info.StackSize - v.Count
		}
	}

	i, j := 0, 46
	if item.Info.Type == cm.ItemTypePotion ||
		item.Info.Type == cm.ItemTypeScroll ||
		(item.Info.Type == cm.ItemTypeScript && item.Info.Effect == 1) {
		i = 0
		j = 4
	} else if item.Info.Type == cm.ItemTypeAmulet {
		i = 4
		j = 6
	} else {
		i = 6
		j = 46
	}
	for i < j {
		if p.Inventory.Items[i] != nil {
			i++
			continue
		}
		p.Inventory.Set(i, item)
		// p.Inventory.Items[i] = ui
		p.EnqueueItemInfo(item.ItemID)
		p.Enqueue(ServerMessage{}.GainedItem(item))
		p.RefreshBagWeight()
		return true
	}
	i = 0
	for i < 46 {
		if p.Inventory.Items[i] != nil {
			i++
			continue
		}
		p.Inventory.Set(i, item)
		// p.Inventory.Items[i] = ui
		p.EnqueueItemInfo(item.ItemID)
		p.Enqueue(ServerMessage{}.GainedItem(item))
		p.RefreshBagWeight()
		return true
	}
	p.ReceiveChat("没有合适的格子放置物品", cm.ChatTypeSystem)
	return false
}

// GainItemMail TODO 往玩家邮箱里发送物品
func (p *Player) GainItemMail(item *cm.UserItem, count int) {
	log.Warnln("还没有实现的功能: 往玩家邮箱里发送物品")
}

// CanGainGold 是否可以增加金币
func (p *Player) CanGainGold(gold uint64) bool {
	return gold+p.Gold <= math.MaxUint64
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
	return util.RandomInt(min, max)
}

func (p *Player) Attacked(attacker IMapObject, damage int, typ cm.DefenceType, damageWeapon bool) int {
	armour := 0
	switch attacker := attacker.(type) {
	case *Player:
		// TODO
	case *Monster:
		switch typ {
		case cm.DefenceTypeACAgility:
			if util.RandomNext(int(p.Agility)+1) > int(attacker.Accuracy) {
				p.BroadcastDamageIndicator(cm.DamageTypeMiss, 0)
				return 0
			}
			armour = p.GetDefencePower(p.MinAC, p.MaxAC)
		case cm.DefenceTypeAC:
			armour = p.GetDefencePower(p.MinAC, p.MaxAC)
		case cm.DefenceTypeMACAgility:
			if util.RandomNext(settings.MagicResistWeight) < int(p.MagicResist) {
				p.BroadcastDamageIndicator(cm.DamageTypeMiss, 0)
				return 0
			}
			if util.RandomNext(int(p.Agility)+1) > int(attacker.Accuracy) {
				return 0
			}
			armour = p.GetDefencePower(p.MinMAC, p.MaxMAC)
		case cm.DefenceTypeMAC:
			if util.RandomNext(settings.MagicResistWeight) < int(p.MagicResist) {
				p.BroadcastDamageIndicator(cm.DamageTypeMiss, 0)
				return 0
			}
			armour = p.GetDefencePower(p.MinMAC, p.MaxMAC)
		case cm.DefenceTypeAgility:
			if util.RandomNext(int(p.Agility)+1) > int(attacker.Accuracy) {
				log.Debugln("Player attacked DefenceTypeAgility")
				log.Debugf("p.Agility: %d, attacker.Accuracy: %d\n", p.Agility, attacker.Accuracy)
				log.Debugf("util.RandomNext(int(p.Agility)+1): %d\n", util.RandomNext(int(p.Agility)+1))
				p.BroadcastDamageIndicator(cm.DamageTypeMiss, 0)
				return 0
			}
		}
		if util.RandomNext(100) < int(p.Reflect) { // TODO ???
			if attacker.IsAttackTarget(p) {
				attacker.Attacked(p, damage, typ, false)
				// CurrentMap.Broadcast(new S.ObjectEffect { ObjectID = ObjectID, Effect = SpellEffect.Reflect }, CurrentLocation);
				p.Broadcast(&server.ObjectEffect{ObjectID: p.GetID(), Effect: cm.SpellEffectReflect, EffectType: 0, DelayTime: 0, Time: 0})
			}
			return 0
		}
		// armour = (int)Math.Max(int.MinValue, (Math.Min(int.MaxValue, (decimal)(armour * ArmourRate))));
		// damage= (int)Math.Max(int.MinValue, (Math.Min(int.MaxValue, (decimal)(damage * DamageRate))));
		armour = util.Int(int(float32(armour) * p.ArmourRate))
		damage = util.Int(int(float32(damage) * p.DamageRate))
		if p.MagicShield {
			damage -= damage * (p.MagicShieldLv + 2) / 10
		}
		if armour >= damage {
			log.Debugf("Player %s armour >= damage. armour: %d, damage: %d\n", p.Name, armour, damage)
			p.BroadcastDamageIndicator(cm.DamageTypeMiss, 0)
			return 0
		}
		// TODO
		/*
			if (p.MagicShield){
				MagicShieldTime -= (damage - armour) * 60;
				AddBuff(new Buff { Type = BuffType.MagicShield, Caster = this, ExpireTime = MagicShieldTime, Values = new int[] { MagicShieldLv } });
			}
			for (int i = PoisonList.Count - 1; i >= 0; i--){
				if (PoisonList[i].PType != PoisonType.LRParalysis){
					continue;
				}
				PoisonList.RemoveAt(i);
				OperateTime = 0;
			}

			LastHitter = attacker.Master ?? attacker;
			LastHitTime = Envir.Time + 10000;
			RegenTime = Envir.Time + RegenDelay;
			LogTime = Envir.Time + Globals.LogDelay;

			DamageDura();
			ActiveBlizzard = false;
			ActiveReincarnation = false;

			CounterAttackCast(GetMagic(Spell.CounterAttack), LastHitter);
		*/

		now := time.Now()
		if now.After(p.StruckTime) {
			// p.Enqueue(new S.Struck { AttackerID = attacker.ObjectID });
			// p.Broadcast(new S.ObjectStruck { ObjectID = ObjectID, AttackerID = attacker.ObjectID, Direction = Direction, Location = CurrentLocation });
			p.Enqueue(&server.Struck{AttackerID: attacker.GetID()})
			p.Broadcast(&server.ObjectStruck{ObjectID: p.GetID(), AttackerID: attacker.GetID(), Direction: p.GetDirection(), LocationX: int32(p.CurrentLocation.X), LocationY: int32(p.CurrentLocation.Y)})
			p.StruckTime = now.Add(500 * time.Millisecond)
		}
		log.Debugf("Player %s attacked, armour: %d, damage: %d. armour - damage: %d\n", p.Name, armour, damage, armour-damage)
		log.Debugf("Player %s HP: %d, MP: %d, MaxHP: %d, MaxMP: %d\n", p.Name, p.HP, p.MP, p.MaxHP, p.MaxMP)
		p.BroadcastDamageIndicator(cm.DamageTypeHit, armour-damage)
		p.ChangeHP(armour - damage)
		log.Debugf("Player %s ChangeHP: %d\n", p.Name, armour-damage)
		log.Debugf("Player %s HP: %d, MP: %d, MaxHP: %d, MaxMP: %d\n", p.Name, p.HP, p.MP, p.MaxHP, p.MaxMP)
		return damage - armour
	}

	return 0
}

func (p *Player) GetDefencePower(min, max uint16) int {
	if min < 0 {
		min = 0
	}
	if min > max {
		max = min
	}
	return util.RandomNext2(int(min), int(max+1))
}

// GainExp 为玩家增加经验
func (p *Player) GainExp(amount uint32) {
	p.Experience += int64(amount)
	p.Enqueue(ServerMessage{}.GainExperience(amount))
	log.Debugf("Player: %s GainExp amount: %d, p.Experience: %d, p.MaxExperience: %d\n", p.Name, amount, p.Experience, p.MaxExperience)
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
	if amount >= uint32(p.MaxHP) {
		amount = uint32(p.MaxHP)
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
	if amount == 0 || p.IsDead() {
		return
	}
	hp := int(p.HP) + amount
	if hp <= 0 {
		hp = 0
	}
	p.SetHP(uint32(hp))
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
	// TODO 复活戒指、凶手武器诅咒、死亡掉落
	p.HP = 0
	p.Dead = true
	// LogTime = Envir.Time;
	// BrownTime = Envir.Time;
	p.Enqueue(&server.Death{Direction: p.GetDirection(), LocationX: int32(p.CurrentLocation.X), LocationY: int32(p.CurrentLocation.Y)})
	p.Broadcast(&server.ObjectDied{ObjectID: p.GetID(), Direction: p.GetDirection(), LocationX: int32(p.CurrentLocation.X), LocationY: int32(p.CurrentLocation.Y), Type: 0})
	/*
		for (int i = 0; i < Buffs.Count; i++)
		{
			if (Buffs[i].Type == BuffType.Curse)
			{
				Buffs.RemoveAt(i);
				break;
			}
		}
		PoisonList.Clear();
		InTrapRock = false;
	*/
	p.CallDefaultNPC(DefaultNPCTypeDie)
}

func (p *Player) Teleport(m *Map, pt cm.Point) bool {
	oldMap := p.Map

	{ // MapObject Teleport
		if m == nil || !m.ValidPoint(pt) {
			log.Warnln("Teleport: map not valid", m == nil)
			return false
		}
		oldMap.DeleteObject(p)
		p.Broadcast(&server.ObjectTeleportOut{ObjectID: p.GetID(), Type: 0})
		p.Broadcast(&server.ObjectRemove{ObjectID: p.GetID()})

		p.Map = m
		p.CurrentLocation = pt

		// InTrapRock = false;
		m.AddObject(p)
		p.BroadcastInfo()

		p.Broadcast(&server.ObjectTeleportIn{ObjectID: p.GetID(), Type: 0})

		p.BroadcastHealthChange()
	}

	p.Enqueue(&server.MapChanged{
		FileName:     m.Info.Filename,
		Title:        m.Info.Title,
		MiniMap:      uint16(m.Info.MiniMap),
		BigMap:       uint16(m.Info.BigMap),
		Lights:       cm.LightSetting(m.Info.Light),
		Location:     p.CurrentLocation,
		Direction:    p.Direction,
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

	return true
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
	defence := args[2].(cm.DefenceType)
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
		case cm.SpellFencing, cm.SpellSpiritSword:
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
	p.ReceiveChat("[欢迎进入游戏，如有任何建议、疑问欢迎交流。联系QQ群：32309474]", cm.ChatTypeHint)
	p.EnqueueItemInfos()
	p.RefreshStats()
	p.EnqueueQuestInfo()
	p.Enqueue(ServerMessage{}.MapInformation(p.Map.Info))
	p.Enqueue(ServerMessage{}.UserInformation(p))
	p.Enqueue(&server.TimeOfDay{Lights: env.Lights})
	// p.EnqueueAreaObjects(nil, p.Map.AOI.GetGridByPoint(p.GetPoint()))
	p.Enqueue(&server.ChangeAMode{Mode: p.AMode})
	p.Enqueue(&server.ChangePMode{Mode: p.PMode})
	p.Enqueue(&server.SwitchGroup{AllowGroup: p.AllowGroup})
	p.EnqueueAreaObjects(nil, p.GetCell())
	p.Enqueue(ServerMessage{}.NPCResponse([]string{}))
	p.Broadcast(ServerMessage{}.ObjectPlayer(p))
}

func (p *Player) StopGame(reason int) {
	p.Broadcast(ServerMessage{}.ObjectRemove(p))
}

func (p *Player) Turn(direction cm.MirDirection) {
	if !p.CanMove() {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	p.Direction = direction
	p.UpdateInSafeZone()

	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Broadcast(ServerMessage{}.ObjectTurn(p))
}

func (p *Player) Walk(direction cm.MirDirection) {
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
	p.Direction = direction
	p.CurrentLocation = n
	p.UpdateInSafeZone()

	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Broadcast(ServerMessage{}.ObjectWalk(p))
}

func (p *Player) Run(direction cm.MirDirection) {
	steps := 2

	var loc cm.Point
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
	p.Direction = direction
	p.CurrentLocation = loc
	p.UpdateInSafeZone()

	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Broadcast(ServerMessage{}.ObjectRun(p))
}

func (p *Player) Chat(message string) {
	// private message 私聊
	if strings.HasPrefix(message, "/") {
		return
	}
	// group 小队
	if strings.HasPrefix(message, "!!") {
		return
	}
	// guild 行会
	if strings.HasPrefix(message, "!~") {
		return
	}
	// mentor 师徒
	if strings.HasPrefix(message, "!#") {
		return
	}
	// shout 喊话
	if strings.HasPrefix(message, "!") {
		return
	}
	// relationship message 夫妻
	if strings.HasPrefix(message, ":)") {
		return
	}
	// GM 喊话
	if strings.HasPrefix(message, "@!") {
		return
	}
	// command 命令
	if strings.HasPrefix(message, "@") {
		msg, err := cmd.Exec(message[1:], p)
		if err != nil {
			p.ReceiveChat(fmt.Sprintf("执行失败(%s)", err), cm.ChatTypeSystem)
		}
		if msg != nil && msg.(string) != "" {
			p.ReceiveChat(msg.(string), cm.ChatTypeSystem)
		}
		return
	}
	msg := ServerMessage{}.ObjectChat(p, message, cm.ChatTypeNormal)
	p.Enqueue(msg)
	p.Broadcast(msg)
}

func (p *Player) MoveItem(mirGridType cm.MirGridType, from int32, to int32) {
	msg := &server.MoveItem{
		Grid:    mirGridType,
		From:    from,
		To:      to,
		Success: false,
	}

	var err error

	switch mirGridType {
	case cm.MirGridTypeInventory:
		err = p.Inventory.Move(int(from), int(to))
	case cm.MirGridTypeStorage:
		err = p.Storage.Move(int(from), int(to))
	case cm.MirGridTypeTrade:
		err = p.Trade.Move(int(from), int(to))
		p.TradeItem()
	case cm.MirGridTypeRefine:
		// TODO
	}

	if err != nil {
		p.ReceiveChat(err.Error(), cm.ChatTypeSystem)
	} else {
		msg.Success = true
	}

	p.Enqueue(msg)
}

func (p *Player) TakeBackItem(from int32, to int32) {
	msg := &server.TakeBackItem{From: from, To: to, Success: false}

	if p.CallingNPC == nil || !util.StringEqualFold(p.CallingNPCPage, StorageKey) || !cm.InRange(p.CurrentLocation, p.CallingNPC.GetPoint(), DataRange) {
		p.Enqueue(msg)
		return
	}

	if int(from) > len(p.Storage.Items) || int(to) > len(p.Inventory.Items) {
		p.Enqueue(msg)
		return
	}

	// item := p.Inventory.Get(int(from))
	// if item.Info.Weight+p.CurrentBagWeight > MaxBagWeight {
	// 	p.ReceiveChat("Too heavy to get back.", cm.ChatTypeSystem)
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

	if p.CallingNPC == nil || !util.StringEqualFold(p.CallingNPCPage, StorageKey) || !cm.InRange(p.CurrentLocation, p.CallingNPC.GetPoint(), DataRange) {
		p.Enqueue(msg)
		return
	}

	if int(from) > len(p.Inventory.Items) || int(to) > len(p.Storage.Items) {
		p.Enqueue(msg)
		return
	}

	item := p.Inventory.Get(int(from))
	if util.HasFlagUint16(item.Info.Bind, cm.BindModeDontStore) {
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

func (p *Player) MergeItem(gridFrom cm.MirGridType, gridTo cm.MirGridType, fromID uint64, toID uint64) {
	msg := &server.MergeItem{
		GridFrom: gridFrom,
		GridTo:   gridTo,
		IDFrom:   fromID,
		IDTo:     toID,
		Success:  false,
	}
	var arrayFrom []*cm.UserItem
	var bagFrom *Bag
	switch gridFrom {
	case cm.MirGridTypeInventory:
		bagFrom = p.Inventory
	case cm.MirGridTypeStorage:
		bagFrom = p.Storage
	case cm.MirGridTypeEquipment:
		bagFrom = p.Equipment
	// case cm.MirGridTypeFishing:
	default:
		p.Enqueue(msg)
		return
	}
	arrayFrom = bagFrom.Items

	var arrayTo []*cm.UserItem
	var bagTo *Bag
	switch gridTo {
	case cm.MirGridTypeInventory:
		bagTo = p.Inventory
	case cm.MirGridTypeStorage:
		bagTo = p.Storage
	case cm.MirGridTypeEquipment:
		bagTo = p.Equipment
	// case cm.MirGridTypeFishing:
	default:
		p.Enqueue(msg)
		return
	}
	arrayTo = bagTo.Items

	var tempFrom *cm.UserItem
	var indexFrom int
	index := -1
	for i := 0; i < len(arrayFrom); i++ {
		if arrayFrom[i] == nil || arrayFrom[i].ID != fromID {
			continue
		}
		index = i
		tempFrom = arrayFrom[i]
		indexFrom = i
		break
	}
	if tempFrom == nil || tempFrom.Info.StackSize == 1 || index == -1 {
		p.Enqueue(msg)
		return
	}

	var tempTo *cm.UserItem
	var indexTo int
	for i := 0; i < len(arrayTo); i++ {
		if arrayTo[i] == nil || arrayTo[i].ID != toID {
			continue
		}
		tempTo = arrayTo[i]
		indexTo = i
		break
	}
	if tempTo == nil || tempTo.Info != tempFrom.Info || tempTo.Count == tempTo.Info.StackSize {
		p.Enqueue(msg)
		return
	}
	if tempTo.Info.Type != cm.ItemTypeAmulet && (gridFrom == cm.MirGridTypeEquipment || gridTo == cm.MirGridTypeEquipment) {
		p.Enqueue(msg)
		return
	}
	if tempTo.Info.Type != cm.ItemTypeBait && (gridFrom == cm.MirGridTypeFishing || gridTo == cm.MirGridTypeFishing) {
		p.Enqueue(msg)
		return
	}
	if tempFrom.Count <= tempTo.Info.StackSize-tempTo.Count {
		tempTo.Count += tempFrom.Count
		bagTo.SetCount(indexTo, tempTo.Count)
		bagFrom.SetCount(indexFrom, 0)
		arrayFrom[index] = nil
	} else {
		tempFrom.Count -= tempTo.Info.StackSize - tempTo.Count
		tempTo.Count = tempTo.Info.StackSize
		bagTo.SetCount(indexTo, tempTo.Count)
		bagFrom.SetCount(indexFrom, tempFrom.Count)
	}
	msg.Success = true
	p.Enqueue(msg)
	p.RefreshStats()
}

func (p *Player) EquipItem(mirGridType cm.MirGridType, id uint64, to int32) {
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
	case cm.MirGridTypeInventory:
		err = p.Inventory.MoveTo(index, int(to), p.Equipment)
	case cm.MirGridTypeStorage:
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

func (p *Player) RemoveItem(mirGridType cm.MirGridType, id uint64, to int32) {
	msg := &server.RemoveItem{
		Grid:     mirGridType,
		UniqueID: id,
		To:       to,
		Success:  false,
	}

	index, item := p.GetUserItemByID(cm.MirGridTypeEquipment, id)
	if item == nil {
		p.Enqueue(msg)
		return
	}

	switch mirGridType {
	case cm.MirGridTypeInventory:
		p.Equipment.MoveTo(index, int(msg.To), p.Inventory)
	case cm.MirGridTypeStorage:

		if !util.StringEqualFold(p.CallingNPCPage, StorageKey) {
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

func (p *Player) RemoveSlotItem(grid cm.MirGridType, id uint64, to int32, to2 cm.MirGridType) {

}

func (p *Player) SplitItem(grid cm.MirGridType, id uint64, count uint32) {
	msg := &server.SplitItem1{
		Grid:     grid,
		UniqueID: id,
		Count:    count,
		Success:  false,
	}
	var bag *Bag
	switch grid {
	case cm.MirGridTypeInventory:
		bag = p.Inventory
	case cm.MirGridTypeStorage:
		bag = p.Storage
	default:
		p.Enqueue(msg)
		return
	}
	index, item := p.GetUserItemByID(grid, id)
	if item == nil || count >= item.Count || p.FreeSpace(bag) == 0 {
		p.Enqueue(msg)
		return
	}
	newItem := env.NewUserItem(item.Info)
	newItem.SoulBoundId = p.GetID()
	newItem.Count = count

	item.Count -= count
	bag.SetCount(index, item.Count)

	msg.Success = true
	p.Enqueue(msg)
	p.Enqueue(&server.SplitItem{Item: newItem, Grid: grid})

	temp := newItem
	array := bag.Items
	if grid == cm.MirGridTypeInventory && (temp.Info.Type == cm.ItemTypePotion || temp.Info.Type == cm.ItemTypeScroll || temp.Info.Type == cm.ItemTypeAmulet || (temp.Info.Type == cm.ItemTypeScript && temp.Info.Effect == 1)) {
		if temp.Info.Type == cm.ItemTypePotion || temp.Info.Type == cm.ItemTypeScroll || (temp.Info.Type == cm.ItemTypeScript && temp.Info.Effect == 1) {
			for i := 0; i < 4; i++ {
				if array[i] != nil {
					continue
				}
				array[i] = temp
				p.Inventory.Set(i, temp)
				p.RefreshBagWeight()
				return
			}
		} else if temp.Info.Type == cm.ItemTypeAmulet {
			for i := 4; i < 6; i++ {
				if array[i] != nil {
					continue
				}
				array[i] = temp
				p.Inventory.Set(i, temp)
				p.RefreshBagWeight()
				return
			}
		}
	}
	for i := 6; i < len(array); i++ {
		if array[i] != nil {
			continue
		}
		array[i] = temp
		p.Inventory.Set(i, temp)
		p.RefreshBagWeight()
		return
	}
	for i := 0; i < 6; i++ {
		if array[i] != nil {
			continue
		}
		array[i] = temp
		p.Inventory.Set(i, temp)
		p.RefreshBagWeight()
		return
	}
}

// FreeSpace Bag 剩余空格数量
func (p *Player) FreeSpace(bag *Bag) int {
	count := 0
	for i := 0; i < len(bag.Items); i++ {
		if bag.Items[i] == nil {
			count++
		}
	}
	return count
}

func (p *Player) TeleportRandom(attempts int, distance uint16, m *Map) bool {
	if m == nil {
		m = p.Map
	}

	for i := 0; i < attempts; i++ {
		loc := cm.NewPoint(util.RandomNext(m.Width), util.RandomNext(m.Height))
		if p.Teleport(m, loc) {
			return true
		}
	}

	return false
}

// Scrolls are consumable items that have various uses.

// Common Name			Shape	Used Stats	Description
// Dungeon Escape		0					Teleports player to a random position on their last saved map.
// Town Teleport		1					Teleports player to their last saved safezone.
// Random Teleport		2					Randomly teleports player to a new coordinate on the same map.
// Benediction Oil		3					Chance to luck the players equipped weapon.
// Repair Oil			4					Repairs equipped weapons durability by 5, whilst reducing its maximum durability.
// WarGod Oil			5					Repairs equipped weapons durability to maximum.
// Resurrection Scroll	6					Revives player with 100% of their HP & MP.
// Credit Scroll		7		Price		Gives player x amount of game shop credits.
// Map Shout Scroll		8					Allows a single special shout across the players current map.
// Server Shout Scroll	9					Allows a single special shout across the server.
// Guild Skill Scroll	10		Effect		Adds a skill to the players guild. Only usable by guild leaders. Skill chosen by effect number.

func (p *Player) UseItemScroll(item *cm.UserItem) bool {
	switch item.Info.Shape {
	case 0: //DE
		temp := env.GetMap(p.BindMapIndex)
		for i := 0; i < 20; i++ {
			x := int(p.BindLocation.X) + util.RandomInt(-100, 100)
			y := int(p.BindLocation.Y) + util.RandomInt(-100, 100)
			loc := cm.NewPoint(x, y)
			if p.Teleport(temp, loc) {
				return true
			}
		}
	case 1: //TT
		if !p.Teleport(env.GetMap(p.BindMapIndex), p.BindLocation) {
			return false
		}
	case 2: //RT
		if !p.TeleportRandom(200, item.Info.Durability, p.Map) {
			return true
		}
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

// Potions are consumable items that will heal or buff the player.

// Common Name		Shape	Used Stats				Description
// Normal Potion	0		HP/MP					Gradually heals player.
// Sun Potion		1		HP/MP					Instantly heals player.
// Mystery Water	2		None					Allows player to unequip a cursed item (officially mystery items only).

// Buff Potion		3		DC/MC/SC/ASpeed/HP/MP/MaxAC/MaxMAC/Durability
//													Gives player the relative buff. Length of buff depends on potions durability. 1 dura = 1 minute.

// Exp Potion		4		Luck/Durability			Increases players percent of exp gain through the luck stat. Length of buff depends on potions durability. 1 dura = 1 minute.

func (p *Player) UserItemPotion(item *cm.UserItem) bool {
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
		log.Debugln("MysteryWater")
	case 3: // Buff
		expireTime := int(info.Durability)
		if info.MaxDC+item.DC > 0 {
			p.AddBuff(NewBuff(cm.BuffTypeImpact, p, expireTime, []int32{int32(info.MaxDC + item.DC)}))
		}
		if info.MaxMC+item.MC > 0 {
			p.AddBuff(NewBuff(cm.BuffTypeMagic, p, expireTime, []int32{int32(info.MaxMC + item.MC)}))
		}
		if info.MaxSC+item.SC > 0 {
			p.AddBuff(NewBuff(cm.BuffTypeTaoist, p, expireTime, []int32{int32(info.MaxSC + item.SC)}))
		}
		if (info.AttackSpeed + item.AttackSpeed) > 0 {
			p.AddBuff(NewBuff(cm.BuffTypeStorm, p, expireTime, []int32{int32(info.AttackSpeed + item.AttackSpeed)}))
		}
		if (info.HP + uint16(item.HP)) > 0 {
			p.AddBuff(NewBuff(cm.BuffTypeHealthAid, p, expireTime, []int32{int32(info.HP + uint16(item.HP))}))
		}
		if (info.MP + uint16(item.MP)) > 0 {
			p.AddBuff(NewBuff(cm.BuffTypeManaAid, p, expireTime, []int32{int32(info.MP + uint16(item.MP))}))
		}
		if (info.MaxAC + item.AC) > 0 {
			p.AddBuff(NewBuff(cm.BuffTypeDefence, p, expireTime, []int32{int32(info.MaxAC + item.AC)}))
		}
		if (info.MaxMAC + item.MAC) > 0 {
			p.AddBuff(NewBuff(cm.BuffTypeMagicDefence, p, expireTime, []int32{int32(info.MaxMAC + item.MAC)}))
		}
	case 4: // Exp 经验
		expireTime := int(info.Durability)
		p.AddBuff(NewBuff(cm.BuffTypeExp, p, expireTime, []int32{int32(info.Luck + item.Luck)}))
	}
	return true
}

func (p *Player) UseItem(id uint64) {
	msg := &server.UseItem{UniqueID: id, Success: false}
	if p.IsDead() {
		p.Enqueue(msg)
		return
	}
	index, item := p.GetUserItemByID(cm.MirGridTypeInventory, id)

	if item == nil || item.ID == 0 || !p.CanUseItem(item) {
		p.Enqueue(msg)
		return
	}
	info := item.Info

	switch info.Type {
	case cm.ItemTypePotion:
		msg.Success = p.UserItemPotion(item)
	case cm.ItemTypeScroll:
		msg.Success = p.UseItemScroll(item)
	case cm.ItemTypeBook:
		msg.Success = p.GiveSkill(cm.Spell(info.Shape), 1)

	case cm.ItemTypeScript:
		p.CallDefaultNPC(DefaultNPCTypeUseItem, info.Shape)
		msg.Success = true
	case cm.ItemTypeFood:
	case cm.ItemTypePets:
	case cm.ItemTypeTransform: //Transforms
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

func (p *Player) GiveSkill(spell cm.Spell, level int) bool {

	info := data.GetMagicInfoBySpell(spell)

	if info != nil {
		for _, v := range p.Magics {
			if v.Spell == spell {
				p.ReceiveChat("你已经学习该技能", cm.ChatTypeSystem)
				return true
			}
		}
		magic := &cm.UserMagic{Info: info, Level: level, CharacterID: int(p.ID), MagicID: info.ID, Spell: spell}
		adb.AddSkill(p, magic)
		p.Magics = append(p.Magics, magic)
		p.Enqueue(&server.NewMagic{Magic: magic.GetClientMagic(magic.Info)})
		p.RefreshStats()
		return true
	}

	return false
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
	index, userItem := p.GetUserItemByID(cm.MirGridTypeInventory, id)
	if userItem == nil || userItem.ID == 0 {
		p.Enqueue(msg)
		return
	}
	obj := NewItem(p, userItem)
	if dropMsg, ok := obj.Drop(p.GetPoint(), 1); !ok {
		p.ReceiveChat(dropMsg, cm.ChatTypeSystem)
		p.Enqueue(msg)
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
		p.ReceiveChat(dropMsg, cm.ChatTypeSystem)
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
	o := env.Players.GetPlayerByID(id)
	if o == nil {
		p.ReceiveChat("获取不到玩家数据", cm.ChatTypeSystem)
		log.Warnln("Player Inspect id error", id)
		return
	}
	for _, i := range o.Equipment.Items {
		i := i
		if i != nil {
			item := data.GetItemInfoByID(int(i.ItemID))
			if item != nil {
				p.EnqueueItemInfo(item.ID)
			}
		}
	}
	p.Enqueue(ServerMessage{}.PlayerInspect(o))
}

func (p *Player) ChangeAMode(mode cm.AttackMode) {
	p.AMode = mode
	p.Enqueue(&server.ChangeAMode{Mode: p.AMode})
}

func (p *Player) ChangePMode(mode cm.PetMode) {
	p.PMode = mode
	p.Enqueue(&server.ChangePMode{Mode: p.PMode})
}

func (p *Player) ChangeTrade(trade bool) {

}

func (p *Player) Attack(direction cm.MirDirection, spell cm.Spell) {
	if !p.CanAttack() {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	level := 0
	switch spell {
	case cm.SpellSlaying:
		if !p.Slaying {
			spell = cm.SpellNone
		} else {
			magic := p.GetMagic(cm.SpellSlaying)
			level = magic.Level
		}
		p.Slaying = false
	case cm.SpellDoubleSlash:
		magic := p.GetMagic(spell)
		if magic == nil || magic.Info.BaseCost+(magic.Level*magic.Info.LevelCost) > int(p.MP) {
			spell = cm.SpellNone
			break
		}
		level = magic.Level
		p.ChangeMP(-(magic.Info.BaseCost + magic.Level*magic.Info.LevelCost))
	case cm.SpellThrusting, cm.SpellFlamingSword:
		magic := p.GetMagic(spell)
		if (magic == nil) || (!p.FlamingSword && (spell == cm.SpellFlamingSword)) {
			spell = cm.SpellNone
			break
		}
		level = magic.Level
	case cm.SpellHalfMoon, cm.SpellCrossHalfMoon:
		magic := p.GetMagic(spell)
		if magic == nil || magic.Info.BaseCost+(magic.Level*magic.Info.LevelCost) > int(p.MP) {
			spell = cm.SpellNone
			break
		}
		level = magic.Level
		p.ChangeMP(-(magic.Info.BaseCost + magic.Level*magic.Info.LevelCost))
	case cm.SpellTwinDrakeBlade:
		magic := p.GetMagic(spell)
		if !p.TwinDrakeBlade || magic == nil || magic.Info.BaseCost+magic.Level*magic.Info.LevelCost > int(p.MP) {
			spell = cm.SpellNone
			break
		}
		level = magic.Level
		p.ChangeMP(-(magic.Info.BaseCost + magic.Level*magic.Info.LevelCost))
	default:
		spell = cm.SpellNone
	}
	if !p.Slaying {
		magic := p.GetMagic(cm.SpellSlaying)
		if magic != nil && util.RandomNext(12) <= magic.Level {
			p.Slaying = true
			p.Enqueue(&server.SpellToggle{Spell: cm.SpellSlaying, CanUse: p.Slaying})
		}
	}
	_ = level // TODO
	p.Direction = direction
	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Broadcast(ServerMessage{}.ObjectAttack(p, spell, 0, 0))
	target := p.GetPoint().NextPoint(p.GetDirection(), 1)
	damageBase := p.GetAttackPower(int(p.MinDC), int(p.MaxDC)) // = the original damage from your gear (+ bonus from moonlight and darkbody)
	cell := p.Map.GetCell(target)
	if !cell.CanWalk() {
		return
	}
	for _, o := range cell.objects {
		if o.GetRace() != cm.ObjectTypePlayer && o.GetRace() != cm.ObjectTypeMonster {
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
		defence := cm.DefenceTypeACAgility
		damageFinal := damageBase
		switch spell {
		case cm.SpellSlaying: // 攻杀剑术
			magic := p.GetMagic(cm.SpellSlaying)
			damageFinal = magic.GetDamage(damageBase)
			p.LevelMagic(magic)
		// case cm.SpellDoubleSlash:
		case cm.SpellThrusting: // 刺杀剑术
			magic := p.GetMagic(cm.SpellThrusting)
			p.LevelMagic(magic)
		case cm.SpellHalfMoon: // 半月弯刀
			magic := p.GetMagic(cm.SpellHalfMoon)
			p.LevelMagic(magic)
		case cm.SpellCrossHalfMoon: // 圆月弯刀
			magic := p.GetMagic(cm.SpellCrossHalfMoon)
			p.LevelMagic(magic)
		case cm.SpellTwinDrakeBlade: // 双龙斩
			magic := p.GetMagic(cm.SpellTwinDrakeBlade)
			damageFinal = magic.GetDamage(damageBase)
			p.TwinDrakeBlade = false
			//   action = new DelayedAction(DelayedType.Damage, Envir.Time + 400, ob, damageFinal, DefenceType.Agility, false);
			p.ActionList.PushDelayAction(DelayedTypeDamage, 400, func() { p.CompleteAttack(o, damageFinal, cm.DefenceTypeAgility, false) })
			p.LevelMagic(magic)
			// TODO
			//   if ((((ob.Race != ObjectType.Player) || Settings.PvpCanResistPoison) && (Envir.Random.Next(Settings.PoisonAttackWeight) >= ob.PoisonResist)) && (ob.Level < Level + 10 && Envir.Random.Next(ob.Race == ObjectType.Player ? 40 : 20) <= magic.Level + 1))
			//   {
			//       ob.ApplyPoison(new Poison { PType = PoisonType.Stun, Duration = ob.Race == ObjectType.Player ? 2 : 2 + magic.Level, TickSpeed = 1000 }, this);
			//       ob.Broadcast(new S.ObjectEffect { ObjectID = ob.ObjectID, Effect = SpellEffect.TwinDrakeBlade });
			//   }
		case cm.SpellFlamingSword: // 烈火剑法
			magic := p.GetMagic(cm.SpellFlamingSword)
			damageFinal = magic.GetDamage(damageBase)
			p.FlamingSword = false
			defence = cm.DefenceTypeAC
			p.LevelMagic(magic)
		}
		p.ActionList.PushDelayAction(DelayedTypeDamage, 300, func() { p.CompleteAttack(o, damageFinal, defence, true) })
	}
}

func (p *Player) RangeAttack(direction cm.MirDirection, location cm.Point, id uint32) {

}

func (p *Player) Harvest(direction cm.MirDirection) {

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
	log.Debugln("call npc", npc.Name, key)

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
	case RepairKey:
		p.Enqueue(&server.NPCRepair{Rate: p.CallingNPC.PriceRate(p, false)})
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
		p.Enqueue(&server.NPCGoods{Goods: goods, Rate: 1.0, Type: cm.PanelTypeBuy})
		return
	}
}

func (p *Player) TalkMonsterNPC(id uint32) {

}

func (p *Player) BuyItem(index uint64, count uint32, panelType cm.PanelType) {
	if p.IsDead() {
		return
	}
	if !util.StringEqualFold(p.CallingNPCPage, BuySellKey, BuyKey, BuyBackKey, BuyUsedKey, PearlBuyKey) {
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

	if !util.StringEqualFold(p.CallingNPCPage, BuySellKey, SellKey) {
		p.Enqueue(msg)
		return
	}

	var index = -1
	var temp *cm.UserItem
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

	if util.HasFlagUint16(temp.Info.Bind, cm.BindModeDontSell) {
		p.Enqueue(msg)
		return
	}
	// if (temp.RentalInformation != null && temp.RentalInformation.BindingFlags.HasFlag(BindMode.DontSell))
	// {
	// 	Enqueue(p);
	// 	return;
	// }
	log.Debugf("SellItem Info.Type: %d\n", temp.Info.Type)
	log.Debugf("CallingNPC.Script.Types: %s\n", p.CallingNPC.Script.Types)
	if !p.CallingNPC.HasType(temp.Info.Type) {
		p.ReceiveChat("不能在这里卖这类商品", cm.ChatTypeSystem)
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

func (p *Player) RepairItem(uniqueID uint64, special bool) {
	p.Enqueue(&server.RepairItem{UniqueID: uniqueID})
	if p.IsDead() {
		return
	}
	if p.CallingNPC == nil ||
		(!util.StringEqualFold(p.CallingNPCPage, RepairKey) && !special) ||
		(!util.StringEqualFold(p.CallingNPCPage, SRepairKey) && special) {
		return
	}
	ob := p.CallingNPC
	// 找到要修理物品temp和物品在背包里的索引index
	var temp *cm.UserItem
	index := -1
	for i := 0; i < p.Inventory.Length(); i++ {
		temp = p.Inventory.Get(i)
		if temp == nil || temp.ID != uniqueID {
			continue
		}
		index = i
		break
	}
	if temp == nil || index == -1 {
		return
	}
	/* FIXME
	if ((temp.Info.Bind.HasFlag(BindMode.DontRepair)) || (temp.Info.Bind.HasFlag(BindMode.NoSRepair) && special)){
		ReceiveChat("你不能修理这个物品。", ChatType.System);
		return;
	}
	*/
	if !ob.HasType(temp.Info.Type) {
		p.ReceiveChat("你不能在这里修理这个物品。", cm.ChatTypeSystem)
		return
	}
	cost := uint32(float32(temp.RepairPrice()) * ob.PriceRate(p, false))
	// baseCost := uint32(float32(temp.RepairPrice()) * ob.PriceRate(p, true))
	if uint64(cost) > p.Gold {
		return
	}
	p.Gold -= uint64(cost)
	p.Enqueue(&server.LoseGold{Gold: cost})
	/* FIXME
	if (ob.Conq != null) {
		ob.Conq.GoldStorage += (cost - baseCost)
	}
	*/
	if !special {
		// temp.MaxDura = (ushort)Math.Max(0, temp.MaxDura - (temp.MaxDura - temp.CurrentDura) / 30)
		temp.MaxDura = util.Uint16(int(temp.MaxDura - (temp.MaxDura-temp.CurrentDura)/30))
	}
	temp.CurrentDura = temp.MaxDura
	temp.DuraChanged = false
	p.Enqueue(&server.ItemRepaired{UniqueID: uniqueID, MaxDura: temp.MaxDura, CurrentDura: temp.CurrentDura})
}

func (p *Player) BuyItemBack(id uint64, count uint32) {

}

func (p *Player) Magic(spell cm.Spell, direction cm.MirDirection, targetID uint32, targetLocation cm.Point) {
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
	p.Direction = direction
	p.ChangeMP(-cost)
	target := p.Map.GetObjectInAreaByID(targetID, targetLocation)

	_, ok := configsMap[spell]
	if !ok {
		p.ReceiveChat("技能还没实现。。。", cm.ChatTypeSystem)
		return
	}

	ctx := &MagicContext{
		Spell:       spell,
		Magic:       magic,
		Target:      target,
		Player:      p,
		TargetPoint: targetLocation,
	}
	err, targetID := startMagic(ctx)
	cast := true
	if err != nil {
		cast = false
		p.ReceiveChat(err.Error(), cm.ChatTypeSystem)
	}

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

func (p *Player) MagicKey(spell cm.Spell, key uint8) {
	for _, m := range p.Magics {
		if m.Spell == spell {
			m.Key = int(key)
			adb.SyncMagicKey(p, spell, key)
			break
		}
	}
}

// SwitchGroup 组队开关切换(是否允许组队)
func (p *Player) SwitchGroup(allow bool) {
	p.Enqueue(&server.SwitchGroup{AllowGroup: allow})
	if p.AllowGroup == allow {
		return
	}
	p.AllowGroup = allow
	if p.AllowGroup || p.GroupMembers == nil || p.GroupMembers.Count() == 0 {
		return
	}
	p.RemoveGroupBuff()
	p.GroupMembers.Remove(p)
	p.Enqueue(&server.DeleteGroup{})
	if p.GroupMembers.Count() > 1 {
		for i := 0; i < p.GroupMembers.Count(); i++ {
			p.GroupMembers.Get(i).Enqueue(&server.DeleteMember{Name: p.Name})
		}
	} else {
		p.GroupMembers.Get(0).Enqueue(&server.DeleteGroup{})
		p.GroupMembers.Get(0).GroupMembers = nil
	}
	p.GroupMembers = nil
	log.Debugf("Player %s SwitchGroup. p.GroupMembers: %s\n", p.Name, p.GroupMembers)
}

// AddMember 添加别的玩家到自己小队
func (p *Player) AddMember(name string) {
	if p.GroupMembers != nil && p.GroupMembers.Get(0) != nil {
		p.ReceiveChat("你不是队长。", cm.ChatTypeSystem)
		return
	}
	if p.GroupMembers != nil && p.GroupMembers.Count() >= MaxGroup {
		p.ReceiveChat("你的队伍人数已满。", cm.ChatTypeSystem)
		return
	}
	player := env.Players.GetPlayerByName(name)
	if player == nil {
		p.ReceiveChat(name+"无法找到。", cm.ChatTypeSystem)
		return
	}
	if player.ID == p.ID {
		p.ReceiveChat("你无法添加你自己。", cm.ChatTypeSystem)
		return
	}
	if !player.AllowGroup {
		p.ReceiveChat(name+"不允许组队。", cm.ChatTypeSystem)
		return
	}
	if player.GroupMembers != nil {
		p.ReceiveChat(name+"已经在另一个队伍。", cm.ChatTypeSystem)
		return
	}
	if player.GroupInvitation != nil {
		p.ReceiveChat(name+"已经收到了另一个玩家的邀请。", cm.ChatTypeSystem)
		return
	}
	p.SwitchGroup(true)
	player.Enqueue(&server.GroupInvite{Name: p.Name})
	player.GroupInvitation = p
	log.Debugf("Player %s AddMember. p.GroupMembers: %s\n", p.Name, p.GroupMembers)
}

// DelMember 删除小队里的玩家
func (p *Player) DelMember(name string) {
	if p.GroupMembers == nil {
		p.ReceiveChat("你不在队伍中。", cm.ChatTypeSystem)
		return
	}
	if p.GroupMembers.Get(0) != p {
		p.ReceiveChat("你不是队长。", cm.ChatTypeSystem)
		return
	}
	var player *Player
	for i := 0; i < p.GroupMembers.Count(); i++ {
		if p.GroupMembers.Get(i).Name != name {
			continue
		}
		player = p.GroupMembers.Get(i)
		break
	}
	if player == nil {
		p.ReceiveChat(name+"不在你的队伍中。", cm.ChatTypeSystem)
		return
	}
	player.RemoveGroupBuff()
	p.GroupMembers.Remove(player)
	player.Enqueue(&server.DeleteGroup{})

	if p.GroupMembers.Count() > 1 {
		packet := &server.DeleteMember{Name: player.Name}
		for i := 0; i < p.GroupMembers.Count(); i++ {
			p.GroupMembers.Get(i).Enqueue(packet)
		}
	} else {
		p.GroupMembers.Get(0).Enqueue(&server.DeleteGroup{})
		p.GroupMembers.Get(0).GroupMembers = nil
	}
	player.GroupMembers = nil
	log.Debugf("Player %s DelMember. p.GroupMembers: %s\n", p.Name, p.GroupMembers)
}

// GroupInvite 玩家是否同意组队
func (p *Player) GroupInvite(accept bool) {
	if p.GroupInvitation == nil {
		p.ReceiveChat("你没有收到邀请或邀请已过期。", cm.ChatTypeSystem)
		return
	}
	if !accept {
		p.GroupInvitation.ReceiveChat(p.Name+"拒绝了你的组队邀请。", cm.ChatTypeSystem)
		p.GroupInvitation = nil
		return
	}
	if p.GroupMembers != nil {
		p.ReceiveChat(fmt.Sprintf("你无法再加入%s的队伍。", p.GroupInvitation.Name), cm.ChatTypeSystem)
		p.GroupInvitation = nil
		return
	}
	if p.GroupInvitation.GroupMembers != nil && p.GroupInvitation.GroupMembers.Get(0) != p.GroupInvitation {
		p.ReceiveChat(p.GroupInvitation.Name+"不再是队长。", cm.ChatTypeSystem)
		p.GroupInvitation = nil
		return
	}
	if p.GroupInvitation.GroupMembers != nil && p.GroupInvitation.GroupMembers.Count() >= MaxGroup {
		p.ReceiveChat(p.GroupInvitation.Name+"的队伍人数已满。", cm.ChatTypeSystem)
		p.GroupInvitation = nil
		return
	}
	if !p.GroupInvitation.AllowGroup {
		p.ReceiveChat(p.GroupInvitation.Name+"不允许组队。", cm.ChatTypeSystem)
		p.GroupInvitation = nil
		return
	}
	if p.GroupInvitation == nil {
		p.ReceiveChat(p.GroupInvitation.Name+"不在线。", cm.ChatTypeSystem)
		p.GroupInvitation = nil
		return
	}
	if p.GroupInvitation.GroupMembers == nil {
		p.GroupInvitation.GroupMembers = NewPlayerList()
		p.GroupInvitation.GroupMembers.Add(p.GroupInvitation)
		p.GroupInvitation.Enqueue(&server.AddMember{Name: p.GroupInvitation.Name})
	}
	packet := &server.AddMember{Name: p.Name}
	p.GroupMembers = p.GroupInvitation.GroupMembers
	p.GroupInvitation = nil
	for i := 0; i < p.GroupMembers.Count(); i++ {
		member := p.GroupMembers.Get(i)
		member.Enqueue(packet)
		p.Enqueue(&server.AddMember{Name: member.Name})
		if p.Map.Info.ID != member.Map.Info.ID || !cm.InRange(p.CurrentLocation, member.CurrentLocation, DataRange) {
			continue
		}

		// byte time = Math.Min(byte.MaxValue, (byte)Math.Max(5, (RevTime - Envir.Time) / 1000));
		// TODO
		time := uint8(10)

		member.Enqueue(&server.ObjectHealth{ObjectID: p.GetID(), Percent: p.GetPercentHealth(), Expire: time})
		p.Enqueue(&server.ObjectHealth{ObjectID: member.GetID(), Percent: member.GetPercentHealth(), Expire: time})
		for j := 0; j < len(member.Pets); j++ {
			pet := member.Pets[j]
			p.Enqueue(&server.ObjectHealth{ObjectID: pet.GetID(), Percent: pet.GetPercentHealth(), Expire: time})
		}
	}

	p.GroupMembers.Add(p)

	// TODO
	/*
		//Adding Buff on for marriage
		if (GroupMembers != null)
		for (int i = 0; i < GroupMembers.Count; i++)
		{
			PlayerObject player = GroupMembers[i];
				if (Info.Married == player.Info.Index)
				{
					AddBuff(new Buff { Type = BuffType.RelationshipEXP, Caster = player, ExpireTime = Envir.Time * 1000, Infinite = true, Values = new int[] { Settings.LoverEXPBonus } });
					player.AddBuff(new Buff { Type = BuffType.RelationshipEXP, Caster = this, ExpireTime = Envir.Time * 1000, Infinite = true, Values = new int[] { Settings.LoverEXPBonus } });
				}
				if (Info.Mentor == player.Info.Index)
				{
					if (Info.isMentor)
					{
						player.AddBuff(new Buff { Type = BuffType.Mentee, Caster = player, ExpireTime = Envir.Time * 1000, Infinite = true, Values = new int[] { Settings.MentorExpBoost } });
						AddBuff(new Buff { Type = BuffType.Mentor, Caster = this, ExpireTime = Envir.Time * 1000, Infinite = true, Values = new int[] { Settings.MentorDamageBoost } });
					}
					else
					{
						AddBuff(new Buff { Type = BuffType.Mentee, Caster = player, ExpireTime = Envir.Time * 1000, Infinite = true, Values = new int[] { Settings.MentorExpBoost } });
						player.AddBuff(new Buff { Type = BuffType.Mentor, Caster = this, ExpireTime = Envir.Time * 1000, Infinite = true, Values = new int[] { Settings.MentorDamageBoost } });
					}
				}
		}
	*/
	for j := 0; j < len(p.Pets); j++ {
		p.Pets[j].BroadcastHealthChange()
	}
	p.Enqueue(packet)
	log.Debugf("Player %s GroupInvite. p.GroupMembers: %s\n", p.Name, p.GroupMembers)
}

// TODO RemoveGroupBuff 删除组队 buff
func (p *Player) RemoveGroupBuff() {

}

func (p *Player) TownRevive() {
	if !p.IsDead() {
		return
	}
	temp := env.GetMap(p.BindMapIndex)
	bindLocation := p.BindLocation
	// TODO 送到红名村
	// if (Info.PKPoints >= 200)
	// {
	// 	temp = Envir.GetMapByNameAndInstance(Settings.PKTownMapName, 1);
	// 	bindLocation = new Point(Settings.PKTownPositionX, Settings.PKTownPositionY);
	// 	if (temp == null)
	// 	{
	// 		temp = Envir.GetMap(BindMapIndex);
	// 		bindLocation = BindLocation;
	// 	}
	// }
	// if (temp == null || !temp.ValidPoint(bindLocation)) return;
	p.Dead = false
	p.SetHP(uint32(p.MaxHP))
	p.SetMP(uint32(p.MaxMP))
	p.RefreshStats()
	p.Broadcast(&server.ObjectRemove{ObjectID: p.GetID()})
	p.Map.DeleteObject(p)
	p.Map = temp
	p.CurrentLocation = bindLocation
	p.Map.AddObject(p)
	p.Enqueue(&server.MapChanged{
		FileName:     p.Map.Info.Filename,
		Title:        p.Map.Info.Title,
		MiniMap:      uint16(p.Map.Info.MiniMap),
		BigMap:       uint16(p.Map.Info.BigMap),
		Lights:       cm.LightSetting(p.Map.Info.Light),
		Location:     p.CurrentLocation,
		Direction:    p.Direction,
		MapDarkLight: uint8(p.Map.Info.MapDarkLight),
		Music:        uint16(p.Map.Info.Music),
	})
	p.EnqueueAreaObjects(nil, p.GetCell())
	p.Broadcast(ServerMessage{}.ObjectPlayer(p))
	p.Enqueue(&server.Revived{})
	p.Broadcast(&server.ObjectRevived{ObjectID: p.GetID(), Effect: true})
	// InSafeZone = true;
	// Fishing = false;
	// Enqueue(GetFishInfo());
}

func (p *Player) SpellToggle(spell cm.Spell, use bool) {

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
	player := env.Players.GetPlayerByID(id)
	if player != nil {
		p.Enqueue(&server.UserName{ID: player.ID, Name: player.Name})
	}
}

func (p *Player) RequestChatItem(id uint64) {

}

func (p *Player) CheckMovement(pos cm.Point) bool {

	// TODO: 优化效率
	for _, v := range data.MovementInfos {
		if v.SourceMap == p.Map.Info.ID {
			if p.CurrentLocation.EqualXY(v.SourceX, v.SourceY) {
				m := env.GetMap(v.DestinationMap)
				if m == nil {
					log.Infoln("no map id=", v.DestinationMap)
				}
				p.Teleport(m, cm.NewPoint(v.DestinationX, v.DestinationY))
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

func (p *Player) MarriageRequest() {

}

func (p *Player) MarriageReply(acceptInvite bool) {

}

func (p *Player) ChangeMarriage() {

}

func (p *Player) DivorceRequest() {

}

func (p *Player) DivorceReply(acceptInvite bool) {

}

func (p *Player) AddMentor(name string) {

}

func (p *Player) MentorReply(acceptInvite bool) {

}

func (p *Player) AllowMentor() {

}

func (p *Player) CancelMentor() {

}

// DepositTradeItem 交易时，把物品从背包放入交易栏
func (p *Player) DepositTradeItem(f int32, t int32) {
	from := int(f)
	to := int(t)
	msg := &server.DepositTradeItem{From: f, To: t, Success: false}
	if from < 0 || from >= p.Inventory.Length() {
		p.Enqueue(msg)
		return
	}
	if to < 0 || to >= p.Trade.Length() {
		p.Enqueue(msg)
		return
	}
	temp := p.Inventory.Get(from)
	if temp == nil {
		p.Enqueue(msg)
		return
	}
	/*	TODO 物品绑定相关
		if (temp.Info.Bind.HasFlag(BindMode.DontTrade))
		{
			Enqueue(p);
			return;
		}
		if (temp.RentalInformation != null && temp.RentalInformation.BindingFlags.HasFlag(BindMode.DontTrade))
		{
			Enqueue(p);
			return;
		}
	*/
	if p.Trade.Get(to) == nil {
		if err := p.Inventory.MoveTo(from, to, p.Trade); err == nil {
			p.RefreshBagWeight()
			p.TradeItem()
			// 记录交易信息 Report.ItemMoved("DepositTradeItem", temp, MirGridType.Inventory, MirGridType.Trade, from, to)
			msg.Success = true
		}
	}
	p.Enqueue(msg)
}

// RetrieveTradeItem 交易时，把物品从交易栏放回背包
func (p *Player) RetrieveTradeItem(f int32, t int32) {
	from := int(f)
	to := int(t)
	msg := &server.RetrieveTradeItem{From: f, To: t, Success: false}
	if from < 0 || from >= p.Trade.Length() {
		p.Enqueue(msg)
		return
	}
	if to < 0 || to >= p.Inventory.Length() {
		p.Enqueue(msg)
		return
	}
	temp := p.Trade.Get(from)
	if temp == nil {
		p.Enqueue(msg)
		return
	}
	if uint16(int(temp.Info.Weight)+p.CurrentBagWeight) > p.MaxBagWeight {
		p.ReceiveChat("负重不足,无法取回。", cm.ChatTypeSystem)
		p.Enqueue(msg)
		return
	}
	if p.Inventory.Get(to) == nil {
		if err := p.Trade.MoveTo(from, to, p.Inventory); err == nil {
			p.RefreshBagWeight()
			p.TradeItem()
			// 记录交易信息 Report.ItemMoved("RetrieveTradeItem", temp, MirGridType.Trade, MirGridType.Inventory, from, to)
			msg.Success = true
		}
	}
	p.Enqueue(msg)
}

// TradeRequest 向面前的别的玩家发出交易请求
func (p *Player) TradeRequest() {
	if p.TradePartner != nil {
		p.ReceiveChat("你已经在交易了。", cm.ChatTypeSystem)
		return
	}
	target := p.GetPoint().NextPoint(p.Direction, 1)
	cell := p.Map.GetCell(target)
	var player *Player
	if cell == nil || !cell.HasObject() {
		return
	}
	for _, ob := range cell.objects {
		if ob.GetRace() != cm.ObjectTypePlayer {
			continue
		}
		player = ob.(*Player)
	}
	if player == nil || p.Direction != player.Direction.NegativeDirection() {
		p.ReceiveChat("交易时你必须和对方面对面。", cm.ChatTypeSystem)
		return
	}
	if player.ID == p.ID {
		p.ReceiveChat("你不能和自己交易。", cm.ChatTypeSystem)
		return
	}
	if player.IsDead() || p.IsDead() {
		p.ReceiveChat("无法和死人交易。", cm.ChatTypeSystem)
		return
	}
	if player.TradeInvitation != nil {
		p.ReceiveChat(fmt.Sprintf("玩家 %s 已经收到了另一个交易邀请。", player.Name), cm.ChatTypeSystem)
		return
	}
	if !player.AllowTrade {
		p.ReceiveChat(fmt.Sprintf("玩家 %s 不允许交易。", player.Name), cm.ChatTypeSystem)
		return
	}
	if !cm.InRange(player.CurrentLocation, p.CurrentLocation, DataRange) || player.Map.Info.ID != p.Map.Info.ID {
		p.ReceiveChat(fmt.Sprintf("玩家 %s 不在交易范围内。", player.Name), cm.ChatTypeSystem)
		return
	}
	if player.TradePartner != nil {
		p.ReceiveChat(fmt.Sprintf("玩家 %s 已经在交易了。", player.Name), cm.ChatTypeSystem)
		return
	}
	player.TradeInvitation = p
	player.Enqueue(&server.TradeRequest{Name: p.Name})
}

// TradeGold 交易时候摆到交易框的金币，如果交易取消则返还金币
func (p *Player) TradeGold(amount uint32) {
	p.TradeUnlock()
	if p.TradePartner == nil {
		return
	}
	if p.Gold < uint64(amount) {
		return
	}
	p.TradeGoldAmount += amount
	p.Gold -= uint64(amount)
	p.Enqueue(&server.LoseGold{Gold: amount})
	p.TradePartner.Enqueue(&server.TradeGold{Amount: p.TradeGoldAmount})
}

// TradeReply 被交易玩家是否同意交易请求（TradeRequest）
func (p *Player) TradeReply(accept bool) {
	if p.TradeInvitation == nil {
		return
	}
	if !accept {
		p.TradeInvitation.ReceiveChat(fmt.Sprintf("玩家 %s 拒绝交易。", p.Name), cm.ChatTypeSystem)
		p.TradeInvitation = nil
		return
	}
	if p.TradePartner != nil {
		p.ReceiveChat("你已经在交易了。", cm.ChatTypeSystem)
		p.TradeInvitation = nil
		return
	}
	if p.TradeInvitation.TradePartner != nil {
		p.ReceiveChat(fmt.Sprintf("玩家 %s 已经在交易了。", p.TradeInvitation.Name), cm.ChatTypeSystem)
		p.TradeInvitation = nil
		return
	}
	p.TradePartner = p.TradeInvitation
	p.TradeInvitation.TradePartner = p
	p.TradeInvitation = nil
	p.Enqueue(&server.TradeAccept{Name: p.TradePartner.Name})
	p.TradePartner.Enqueue(&server.TradeAccept{Name: p.Name})
}

// TradeItem 发送给交易对象所要交易的物品
func (p *Player) TradeItem() {
	p.TradeUnlock()
	if p.TradePartner == nil {
		return
	}
	for i := 0; i < p.Trade.Length(); i++ {
		u := p.Trade.Get(i)
		if u == nil {
			continue
		}
		p.TradePartner.CheckItem(u)
	}
	p.TradePartner.Enqueue(&server.TradeItem{TradeItems: p.Trade.Items})
}

// TradeUnlock 取消确认交易
func (p *Player) TradeUnlock() {
	p.TradeLocked = false
	if p.TradePartner != nil {
		p.TradePartner.TradeLocked = false
	}
}

// TradeConfirm 摆好物品到交易栏后，确认交易
func (p *Player) TradeConfirm(confirm bool) {
	if !confirm {
		p.TradeLocked = false
		return
	}
	if p.TradePartner == nil {
		p.TradeCancel()
		return
	}
	if !cm.InRange(p.TradePartner.CurrentLocation, p.CurrentLocation, DataRange) || p.TradePartner.Map.Info.ID != p.Map.Info.ID ||
		!cm.FacingEachOther(p.Direction, p.CurrentLocation, p.TradePartner.Direction, p.TradePartner.CurrentLocation) {
		p.TradeCancel()
		return
	}
	p.TradeLocked = true
	if p.TradeLocked && !p.TradePartner.TradeLocked {
		p.TradePartner.ReceiveChat(fmt.Sprintf("玩家 %s 正在等待你确认交易。", p.Name), cm.ChatTypeSystem)
	}
	if !p.TradeLocked || !p.TradePartner.TradeLocked {
		return
	}

	tradePair := [2]*Player{p.TradePartner, p}
	canTrade := true
	var u *cm.UserItem

	//check if both people can accept the others items
	for x := 0; x < 2; x++ { // p -> x
		o := 0
		if x == 0 {
			o = 1
		}
		if !tradePair[o].CanGainItems(tradePair[x].Trade) {
			canTrade = false
			tradePair[x].ReceiveChat("对方物品已满。", cm.ChatTypeSystem)
			tradePair[x].Enqueue(&server.TradeCancel{Unlock: true})

			tradePair[o].ReceiveChat("物品已满。", cm.ChatTypeSystem)
			tradePair[o].Enqueue(&server.TradeCancel{Unlock: true})
			return
		}
		if !tradePair[o].CanGainGold(uint64(tradePair[x].TradeGoldAmount)) {
			canTrade = false
			tradePair[x].ReceiveChat("对方金币已满。", cm.ChatTypeSystem)
			tradePair[x].Enqueue(&server.TradeCancel{Unlock: true})

			tradePair[o].ReceiveChat("金币已满。", cm.ChatTypeSystem)
			tradePair[o].Enqueue(&server.TradeCancel{Unlock: true})
			return
		}
	}
	if !canTrade {
		return
	}
	// 交换物品
	for x := 0; x < 2; x++ { // p -> x
		o := 0
		if x == 0 {
			o = 1
		}
		log.Debugf("交换物品 (%s x %d), (%s o %d)\n", tradePair[x].Name, x, tradePair[o].Name, o)
		for i := 0; i < tradePair[x].Trade.Length(); i++ {
			u = tradePair[x].Trade.Get(i)
			if u == nil {
				continue
			}
			tradePair[o].GainItem(u.Clone(env.NewObjectID()))
			tradePair[x].Trade.Set(i, nil)
			// 记录交易信息 Report.ItemMoved("TradeConfirm", u, MirGridType.Trade, MirGridType.Inventory, i, -99, string.Format("Trade from {0} to {1}", TradePair[p].Name, TradePair[o].Name));
			log.Debugf("(%s o %d) 交出 %s 给 (%s x %d)\n", tradePair[o].Name, o, u.Info.Name, tradePair[x].Name, x)
		}
		if tradePair[x].TradeGoldAmount > 0 {
			// Report.GoldChanged("TradeConfirm", TradePair[p].TradeGoldAmount, true, string.Format("Trade from {0} to {1}", TradePair[p].Name, TradePair[o].Name));
			tradePair[o].GainGold(uint64(tradePair[x].TradeGoldAmount))
			tradePair[x].TradeGoldAmount = 0
		}
		tradePair[x].ReceiveChat("交易成功。", cm.ChatTypeSystem)
		tradePair[x].Enqueue(&server.TradeConfirm{})
		tradePair[x].TradeLocked = false
		tradePair[x].TradePartner = nil
	}
}

// TradeCancel 交易双方任何一方取消交易
func (p *Player) TradeCancel() {
	p.TradeUnlock()
	if p.TradePartner == nil {
		return
	}
	tradePair := [2]*Player{p.TradePartner, p}
	for k := 0; k < 2; k++ { // p -> k
		if tradePair[k] == nil {
			continue
		}
		for t := 0; t < tradePair[k].Trade.Length(); t++ {
			temp := tradePair[k].Trade.Get(t)
			if temp == nil {
				continue
			}
			if p.FreeSpace(tradePair[k].Inventory) < 1 {
				tradePair[k].GainItemMail(temp, 1)
				// Report.ItemMailed("TradeCancel", temp, temp.Count, 1);
				tradePair[k].Enqueue(&server.DeleteItem{UniqueID: temp.ID, Count: temp.Count})
				tradePair[k].Trade.Set(t, nil)
				continue
			}
			for i := 0; i < tradePair[k].Inventory.Length(); i++ {
				if tradePair[k].Inventory.Get(i) != nil {
					continue
				}
				// 把物品放回玩家背包
				if tradePair[k].CanGainItem(temp) {
					tradePair[k].RetrieveTradeItem(int32(t), int32(i))
				} else {
					// FIXME 背包放不下的情况
					// 这里原来的 C# 是发送物品到玩家邮箱
					// 但是 mirgo 还没有实现邮箱
					// 临时解决方法:
					// 1. 弄个遗失的物品管理 NPC
					// 2. 存到玩家仓库（玩家仓库 Storage 也有装不下的问题）
					// 3. 简单粗暴 直接丢地上
					tradePair[k].GainItemMail(temp, 1)
					// Report.ItemMailed("TradeCancel", temp, temp.Count, 1);
					tradePair[k].Enqueue(&server.DeleteItem{UniqueID: temp.ID, Count: temp.Count})
				}
				// tradePair[k].Trade.Set(t, nil)
				break
			}
		}
		// 返还放入交易栏的金币
		if tradePair[k].TradeGoldAmount > 0 {
			// Report.GoldChanged("TradeCancel", TradePair[p].TradeGoldAmount, false);
			tradePair[k].GainGold(uint64(tradePair[k].TradeGoldAmount))
			tradePair[k].TradeGoldAmount = 0
		}
		tradePair[k].TradeLocked = false
		tradePair[k].TradePartner = nil
		tradePair[k].Enqueue(&server.TradeCancel{Unlock: false})
	}
}

func (p *Player) EquipSlotItem(grid cm.MirGridType, uniqueID uint64, to int32, gridTo cm.MirGridType) {

}

func (p *Player) FishingCast(castOut bool) {

}

func (p *Player) FishingChangeAutocast(autoCast bool) {

}

func (p *Player) AcceptQuest(npcIndex uint32, questIndex int32) {

}

func (p *Player) FinishQuest(questIndex int32, selectedItemIndex int32) {

}

func (p *Player) AbandonQuest(questIndex int32) {

}

func (p *Player) ShareQuest(questIndex int32) {

}

func (p *Player) AcceptReincarnation() {

}

func (p *Player) CancelReincarnation() {

}

func (p *Player) CombineItem(idFrom uint64, idTo uint64) {

}

func (p *Player) SetConcentration(objectID uint32, enabled bool, interrupted bool) {

}

func (p *Player) CreateGuild(name string) bool {
	if (p.MyGuild != nil) || (p.GuildIndex != -1) {
		return false
	}
	if env.GetGuild(name) != nil {
		return false
	}
	if p.Level < uint16(settings.Guild_RequiredLevel) {
		p.ReceiveChat(fmt.Sprintf("你的等级不够,无法创建行会,需要: %d", settings.Guild_RequiredLevel), cm.ChatTypeSystem)
		return false
	}

	//check if we have the required items 检查是否有创建行会所需要的物品
	// for (int i = 0; i < Settings.Guild_CreationCostList.Count; i++)
	//take the required items 从创建人拿走创建行会所需的物品
	// for (int i = 0; i < Settings.Guild_CreationCostList.Count; i++)

	p.RefreshStats()

	//make the guild
	guild := NewGuild(p, name)
	guild.GuildIndex = int(env.NewObjectID()) // FIXME 溢出
	env.GuildList.Add(guild)
	p.GuildIndex = guild.GuildIndex
	p.MyGuild = guild

	p.MyGuildRank = guild.FindRank(name)
	p.GuildMembersChanged = true
	p.GuildNoticeChanged = true
	p.GuildCanRequestItems = true

	//tell us we now have a guild
	p.BroadcastInfo()
	p.MyGuild.SendGuildStatus(p)
	return true
}

func (p *Player) EditGuildMember(name string, rankName string, rankIndex uint8, changeType uint8) {
	if (p.MyGuild == nil) || (p.MyGuildRank == nil) {
		p.ReceiveChat("你不在行会里。", cm.ChatTypeSystem)
		return
	}
	switch changeType {
	case 0: //add member
		if util.HasFlagUint8(uint8(p.MyGuildRank.Options), cm.RankOptionsCanRecruit) {
			p.ReceiveChat("你没有招募新成员的权限。", cm.ChatTypeSystem)
			return
		}
		if name == "" {
			return
		}
		player := env.GetPlayer(name)
		if player == nil {
			p.ReceiveChat(fmt.Sprintf("%s 不在线。", name), cm.ChatTypeSystem)
			return
		}
		if (player.MyGuild != nil) || (player.MyGuildRank != nil) || (player.GuildIndex != -1) {
			p.ReceiveChat(fmt.Sprintf("%s 已经在一个行会里了。", name), cm.ChatTypeSystem)
			return
		}
		if !player.EnableGuildInvite {
			p.ReceiveChat(fmt.Sprintf("%s 不接受行会邀请。", name), cm.ChatTypeSystem)
			return
		}
		if player.PendingGuildInvite != nil {
			p.ReceiveChat(fmt.Sprintf("%s 已经接到了一个行会的邀请。", name), cm.ChatTypeSystem)
			return
		}

		if p.MyGuild.IsAtWar() {
			p.ReceiveChat("在行会战争期间不能招募成员。", cm.ChatTypeSystem)
			return
		}

		player.Enqueue(&server.GuildInvite{Name: p.MyGuild.Name})
		player.PendingGuildInvite = p.MyGuild
	case 1: //delete member
		if !util.HasFlagUint8(uint8(p.MyGuildRank.Options), cm.RankOptionsCanKick) {
			p.ReceiveChat("你没有移除成员的权限。", cm.ChatTypeSystem)
			return
		}
		if name == "" {
			return
		}

		if !p.MyGuild.DeleteMember(p, name) {
			return
		}
	case 2: //promote member (and it'll auto create a new rank at bottom if the index > total ranks!)
		if !util.HasFlagUint8(uint8(p.MyGuildRank.Options), uint8(cm.RankOptionsCanChangeRank)) {
			p.ReceiveChat("你没有更改其他成员行会职位的权限。", cm.ChatTypeSystem)
			return
		}
		if name == "" {
			return
		}
		p.MyGuild.ChangeRank(p, name, rankIndex, rankName)
	case 3: //change rank name
		if !util.HasFlagUint8(uint8(p.MyGuildRank.Options), uint8(cm.RankOptionsCanChangeRank)) {
			p.ReceiveChat("你没有更改其他成员行会职位的权限。", cm.ChatTypeSystem)
			return
		}
		if (rankName == "") || len(rankName) < 3 {
			p.ReceiveChat("行会职位名称太短。", cm.ChatTypeSystem)
			return
		}
		// FIXME 判断行会名格式是否合法
		// if RankName.Contains("\\") || RankName.Length > 20 {
		// 	return
		// }
		if !p.MyGuild.ChangeRankName(p, rankName, rankIndex) {
			return
		}
	case 4: //new rank
		if !util.HasFlagUint8(uint8(p.MyGuildRank.Options), uint8(cm.RankOptionsCanChangeRank)) {
			p.ReceiveChat("你没有更改行会职位的权限。", cm.ChatTypeSystem)
			return
		}
		if len(p.MyGuild.Ranks) > 254 {
			p.ReceiveChat("没有更多的行会职位位置可用。", cm.ChatTypeSystem)
			return
		}
		p.MyGuild.NewRank(p)
	case 5: //change rank setting
		if !util.HasFlagUint8(uint8(p.MyGuildRank.Options), uint8(cm.RankOptionsCanChangeRank)) {
			p.ReceiveChat("你没有更改行会职位的权限。", cm.ChatTypeSystem)
			return
		}
		/* FIXME
		int temp;
		if (!int.TryParse(RankName, out temp)){
			return;
		}
		*/
		temp := 0
		p.MyGuild.ChangeRankOption(p, rankIndex, temp, name)
	}
}

func (p *Player) EditGuildNotice(notice []string) {
	if (p.MyGuild == nil) || (p.MyGuildRank == nil) {
		p.ReceiveChat("你不在一个行会里。", cm.ChatTypeSystem)
		return
	}
	if !util.HasFlagUint8(uint8(p.MyGuildRank.Options), uint8(cm.RankOptionsCanChangeNotice)) {
		p.ReceiveChat("你没有更改行会公告的权限。", cm.ChatTypeSystem)
		return
	}
	if len(notice) > 200 {
		p.ReceiveChat("行会公告不能超过200行。", cm.ChatTypeSystem)
		return
	}
	p.MyGuild.NewNotice(notice)
}

func (p *Player) GuildInvite(accept bool) {
	if p.PendingGuildInvite == nil {
		p.ReceiveChat("你没有收到行会邀请或邀请已过期。", cm.ChatTypeSystem)
		return
	}
	if !accept {
		return
	}
	if !p.PendingGuildInvite.HasRoom() {
		p.ReceiveChat(fmt.Sprintf("%s 已满。", p.PendingGuildInvite.Name), cm.ChatTypeSystem)
		return
	}
	p.PendingGuildInvite.NewMember(p)
	p.GuildIndex = p.PendingGuildInvite.GuildIndex
	p.MyGuild = p.PendingGuildInvite
	p.MyGuildRank = p.PendingGuildInvite.FindRank(p.Name)
	p.GuildMembersChanged = true
	p.GuildNoticeChanged = true
	//tell us we now have a guild
	p.BroadcastInfo()
	p.MyGuild.SendGuildStatus(p)
	p.PendingGuildInvite = nil
	p.EnableGuildInvite = false
	p.GuildCanRequestItems = true
	//refresh guildbuffs
	p.RefreshStats()
	// FIXME 工会 BUFF
	// if (MyGuild.BuffList.Count > 0) {
	// 	p.Enqueue(&server.GuildBuffList() { ActiveBuffs : p.MyGuild.BuffList});
	// }
}

func (p *Player) RequestGuildInfo(typ uint8) {
	if p.MyGuild == nil {
		return
	}
	if p.MyGuildRank == nil {
		return
	}
	switch typ {
	case 0: //notice
		if p.GuildNoticeChanged {
			p.Enqueue(&server.GuildNoticeChange{Notice: p.MyGuild.Notice})
		}
		p.GuildNoticeChanged = false
	case 1: //memberlist
		if p.GuildMembersChanged {
			p.Enqueue(&server.GuildMemberChange{Status: 255, Ranks: p.MyGuild.Ranks})
		}
	}
}

func (p *Player) GuildNameReturn(name string) {
	if name == "" {
		p.CanCreateGuild = false
	}
	if !p.CanCreateGuild {
		return
	}
	if (len(name) < 3) || (len(name) > 20) {
		p.ReceiveChat("行会名字过长。", cm.ChatTypeSystem)
		p.CanCreateGuild = false
		return
	}
	if strings.Contains(name, "\\") {
		p.CanCreateGuild = false
		return
	}
	if p.MyGuild != nil {
		p.ReceiveChat("你已经是行会的一员了。", cm.ChatTypeSystem)
		p.CanCreateGuild = false
		return
	}
	guild := env.GetGuild(name)
	if guild != nil {
		p.ReceiveChat(fmt.Sprintf("行会 %s 已存在。", name), cm.ChatTypeSystem)
		p.CanCreateGuild = false
		return
	}
	p.CreateGuild(name)
	p.CanCreateGuild = false
}

func (p *Player) GuildStorageGoldChange(tpy uint8, amount uint32) {

}

func (p *Player) GuildStorageItemChange(tpy uint8, from int32, to int32) {

}

func (p *Player) GuildWarReturn(name string) {

}

func (p *Player) AtWar(attacker *Player) {

}

func (p *Player) GuildBuffUpdate(typ byte, id int) {

}
