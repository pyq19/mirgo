package main

import (
	"fmt"
	"github.com/yenkeia/mirgo/setting"
	"strconv"
	"strings"
	"sync"
	"time"

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
	Gold               uint64
	GuildName          string
	GuildRankName      string
	Class              common.MirClass
	Gender             common.MirGender
	Hair               uint8
	Light              uint8
	Inventory          []common.UserItem // 46
	Equipment          []common.UserItem // 14
	QuestInventory     []common.UserItem // 40
	Trade              []common.UserItem // 10
	Refine             []common.UserItem // 16
	LooksArmour        int
	LooksWings         int
	LooksWeapon        int
	LooksWeaponEffect  int
	SendItemInfo       []common.ItemInfo
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
	MaxExperience      int64
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
	Magics             []common.UserMagic
	ActionList         *sync.Map // map[uint32]DelayedAction
	Health             Health    // 状态恢复
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

func (p *Player) GetID() uint32 {
	return p.ID
}

func (p *Player) Point() common.Point {
	return p.GetPoint()
}

func (p *Player) GetRace() common.ObjectType {
	return common.ObjectTypePlayer
}

func (p *Player) GetPoint() common.Point {
	return p.CurrentLocation
}

func (p *Player) GetCell() *Cell {
	return p.Map.GetCell(p.CurrentLocation)
}

func (p *Player) GetDirection() common.MirDirection {
	return p.CurrentDirection
}

func (p *Player) GetCurrentGrid() *Grid {
	return p.Map.AOI.GetGridByPoint(p.Point())
}

// IsAttackTarget 判断玩家是否是攻击者的攻击对象
func (p *Player) IsAttackTarget(attacker IMapObject) bool {
	return false
}

func (p *Player) IsFriendlyTarget(attacker IMapObject) bool {
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

func (p *Player) NewObjectID() uint32 {
	return p.Map.Env.NewObjectID()
}

func (p *Player) IsDead() bool {
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
	p.Map.Submit(NewTask(func(args ...interface{}) {
		grids := p.Map.AOI.GetSurroundGrids(p.Point())
		for i := range grids {
			areaPlayers := grids[i].GetAllPlayer()
			for i := range areaPlayers {
				if p.GetID() == areaPlayers[i].GetID() {
					continue
				}
				areaPlayers[i].Enqueue(msg)
			}
		}
	}))
}

func (p *Player) Process() {
	finishID := make([]uint32, 0)
	now := time.Now()
	p.ActionList.Range(func(k, v interface{}) bool {
		action := v.(*DelayedAction)
		if action.Finish || now.Before(action.ActionTime) {
			return true
		}
		action.Task.Execute()
		action.Finish = true
		if action.Finish {
			finishID = append(finishID, action.ID)
		}
		return true
	})
	for i := range finishID {
		p.ActionList.Delete(finishID[i])
	}
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

func (p *Player) EnqueueItemInfos() {
	gdb := p.Map.Env.GameDB
	itemInfos := make([]*common.ItemInfo, 0)
	for i := range p.Inventory {
		itemID := int(p.Inventory[i].ItemID)
		if itemID == 0 {
			continue
		}
		itemInfos = append(itemInfos, gdb.GetItemInfoByID(itemID))
	}
	for i := range p.Equipment {
		itemID := int(p.Equipment[i].ItemID)
		if itemID == 0 {
			continue
		}
		itemInfos = append(itemInfos, gdb.GetItemInfoByID(itemID))
	}
	for i := range p.QuestInventory {
		itemID := int(p.QuestInventory[i].ItemID)
		if itemID == 0 {
			continue
		}
		itemInfos = append(itemInfos, gdb.GetItemInfoByID(itemID))
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
	item := p.Map.Env.GameDB.GetItemInfoByID(int(itemID))
	if item == nil {
		return
	}
	p.Enqueue(ServerMessage{}.NewItemInfo(item))
	p.SendItemInfo = append(p.SendItemInfo, *item)
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
	p.MaxExperience = 100
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
	for i := range p.Inventory {
		ui := p.Inventory[i]
		if ui.ID != 0 {
			it := p.Map.Env.GameDB.GetItemInfoByID(int(ui.ItemID))
			p.CurrentBagWeight += int(it.Weight)
		}
	}
}

func (p *Player) RefreshEquipmentStats() {
	gdb := p.Map.Env.GameDB
	for i := range p.Equipment {
		e := gdb.GetItemInfoByID(int(p.Equipment[i].ItemID))
		if e == nil {
			continue
		}
		switch e.Type {
		case common.ItemTypeArmour:
			p.LooksArmour = int(e.Shape)
			p.LooksWings = int(e.Effect)
		case common.ItemTypeWeapon:
			p.LooksWeapon = int(e.Shape)
			p.LooksWeaponEffect = int(e.Effect)
		}
	}
}

func (p *Player) RefreshItemSetStats() {

}

func (p *Player) RefreshMirSetStats() {

}

func (p *Player) RefreshSkills() {

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
	var arr []common.UserItem
	switch mirGridType {
	case common.MirGridTypeInventory:
		arr = p.Inventory
	case common.MirGridTypeEquipment:
		arr = p.Equipment
	default:
		panic("error mirGridType")
	}
	for i := range arr {
		item := arr[i]
		if item.ID == id {
			return i, &item
		}
	}
	return -1, nil
}

// GainItem 为玩家增加物品，增加成功返回 true
func (p *Player) GainItem(ui *common.UserItem) bool {
	item := p.Map.Env.GameDB.GetItemInfoByID(int(ui.ItemID))
	if item == nil {
		return false
	}
	i, j := 6, 46
	if item.Type == common.ItemTypePotion ||
		item.Type == common.ItemTypeScroll ||
		item.Type == common.ItemTypeScript ||
		item.Type == common.ItemTypeAmulet {
		i = 0
		j = 4
	} else if item.Type == common.ItemTypeAmulet {
		i = 4
		j = 6
	} else {
		i = 6
	}
	for i < j {
		if p.Inventory[i].ID != 0 {
			i++
			continue
		}
		p.Inventory[i] = *ui
		break
	}
	p.EnqueueItemInfo(ui.ItemID)
	p.Enqueue(ServerMessage{}.GainedItem(ui))
	p.RefreshBagWeight()
	return true
}

// GainGold 为玩家增加金币
func (p *Player) GainGold(gold uint64) {
	if gold <= 0 {
		return
	}
	p.Gold += gold
	p.Enqueue(ServerMessage{}.GainedGold(gold))
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
	return G_Rand.RandInt(min, max+1)
}

// TODO
func (p *Player) Attacked(attacker IMapObject, damageFinal int, defenceType common.DefenceType, damageWeapon bool) {

}

func (p *Player) GainExp(amount uint32) {
	p.Experience += int64(amount)
	p.Enqueue(ServerMessage{}.GainExperience(amount))
	if p.Experience < p.MaxExperience {
		return
	}
	p.Experience -= p.MaxExperience
	p.Level++
	p.LevelUp()
}

func (p *Player) SetHP(amount uint32) {
	p.HP = uint16(amount)
	msg := ServerMessage{}.HealthChanged(p.HP, p.MP)
	p.Enqueue(msg)
	p.Broadcast(msg)
}

func (p *Player) SetMP(amount uint32) {
	p.MP = uint16(amount)
	msg := ServerMessage{}.HealthChanged(p.HP, p.MP)
	p.Enqueue(msg)
	p.Broadcast(msg)
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

}

func (p *Player) EnqueueAreaObjects(oldGrid, newGrid *Grid) {
	oldAreaGrids := make([]*Grid, 0)
	if oldGrid != nil {
		oldAreaGrids = p.Map.AOI.GetSurroundGridsByGridID(oldGrid.GID)
	}
	newAreaGrids := p.Map.AOI.GetSurroundGridsByGridID(newGrid.GID)
	send := make(map[int]bool)
	for i := range newAreaGrids {
		ng := newAreaGrids[i]
		send[ng.GID] = true
		for j := range oldAreaGrids {
			og := oldAreaGrids[j]
			if ng.GID == og.GID {
				send[ng.GID] = false
			}
		}
	}
	newAreaObjects := make([]IMapObject, 0)
	for i := range newAreaGrids {
		ng := newAreaGrids[i]
		if send[ng.GID] {
			newAreaObjects = append(newAreaObjects, ng.GetAllObjects()...)
		}
	}
	for i := range newAreaObjects {
		if o := newAreaObjects[i]; o.GetID() != p.GetID() {
			p.Enqueue(ServerMessage{}.Object(o))
		}
	}
	drop := make(map[int]bool)
	for i := range oldAreaGrids {
		og := oldAreaGrids[i]
		drop[og.GID] = true
		for j := range newAreaGrids {
			ng := newAreaGrids[j]
			if og.GID == ng.GID {
				drop[og.GID] = false
			}
		}
	}
	oldAreaObjects := make([]IMapObject, 0)
	for i := range oldAreaGrids {
		og := oldAreaGrids[i]
		if drop[og.GID] {
			oldAreaObjects = append(oldAreaObjects, og.GetAllObjects()...)
		}
	}
	for i := range oldAreaObjects {
		if o := oldAreaObjects[i]; o.GetID() != p.GetID() {
			p.Enqueue(ServerMessage{}.ObjectRemove(o))
		}
	}
}

func (p *Player) CompleteAttack(args ...interface{})          {}
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
	p.EnqueueAreaObjects(nil, p.Map.AOI.GetGridByPoint(p.GetPoint()))
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
	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Broadcast(ServerMessage{}.ObjectTurn(p))
}

func (p *Player) Walk(direction common.MirDirection) {
	if !p.CanMove() || !p.CanWalk() {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	n := p.Point().NextPoint(direction, 1)
	ok := p.Map.UpdateObject(p, n)
	if !ok {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	p.CurrentDirection = direction
	p.CurrentLocation = n
	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Broadcast(ServerMessage{}.ObjectWalk(p))
}

func (p *Player) Run(direction common.MirDirection) {
	n1 := p.Point().NextPoint(direction, 1)
	n2 := p.Point().NextPoint(direction, 2)
	if ok := p.Map.UpdateObject(p, n1, n2); !ok {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	p.CurrentDirection = direction
	p.CurrentLocation = n2
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

	curMap := p.Map

	if strings.HasPrefix(message, "@") {
		parts := strings.Split(message[1:], " ")
		switch strings.ToUpper(parts[0]) {
		case "LOGIN":
		case "KILL": // @kill 杀死面前的怪物，@kill name 杀死名字为 name 的玩家
			if len(parts) == 2 {
				o := curMap.Env.GetPlayerByName(parts[1])
				if o == nil {
					p.ReceiveChat(fmt.Sprintf("找不到玩家(%s)", parts[1]), common.ChatTypeSystem)
					return
				}
				o.Die()
				return
			}
			c := curMap.GetNextCell(p.GetCell(), p.GetDirection(), 1)
			if c == nil {
				return
			}
			c.Objects.Range(func(k, v interface{}) bool {
				if o, ok := v.(*Monster); ok {
					o.Die()
				}
				return true
			})
		case "RESTORE":
		case "CHANGEGENDER":
		case "LEVEL":
		case "MAKE": // @make 物品名 数量
			if len(parts) != 3 {
				return
			}
			info := curMap.Env.GameDB.GetItemInfoByName(parts[1])
			if info == nil {
				return
			}
			tmp, err := strconv.Atoi(parts[2])
			if err != nil || tmp > 100 {
				return
			}
			count := uint32(tmp)
			for count > 0 {
				if info.StackSize >= count {
					userItem := curMap.Env.NewUserItem(info)
					userItem.Count = count
					p.GainItem(userItem)
					return
				}
				userItem := curMap.Env.NewUserItem(info)
				userItem.Count = count
				count -= info.StackSize
				p.GainItem(userItem)
			}
			p.ReceiveChat(fmt.Sprintf("%s x %d 创建成功", info.Name, count), common.ChatTypeSystem)
		case "CLEARBUFFS":
		case "CLEARBAG":
		case "SUPERMAN":
		case "GAMEMASTER":
		case "OBSERVER":
		case "ALLOWGUILD":
		case "RECALL":
		case "ENABLEGROUPRECALL":
		case "GROUPRECALL":
		case "RECALLMEMBER":
		case "RECALLLOVER":
		case "TIME":
		case "ROLL":
		case "MAP":
			p.ReceiveChat(fmt.Sprintf("当前地图: %s, ID: %d", curMap.Info.Title, curMap.Info.ID), common.ChatTypeSystem)
		case "MOVE": // @move x y
			if len(parts) != 3 {
				p.ReceiveChat(fmt.Sprintf("移动失败，正确命令格式: @move 123 456"), common.ChatTypeSystem)
				return
			}
			x, err := strconv.Atoi(parts[1])
			if err != nil {
				p.ReceiveChat(fmt.Sprintf("移动失败，正确命令格式: @move 123 456"), common.ChatTypeSystem)
				return
			}
			y, err := strconv.Atoi(parts[2])
			if err != nil {
				p.ReceiveChat(fmt.Sprintf("移动失败，正确命令格式: @move 123 456"), common.ChatTypeSystem)
				return
			}
			p.Teleport(curMap, common.NewPoint(x, y))
		case "MAPMOVE":
		case "GOTO":
		case "MOB": // @mob 怪物名称		在玩家周围生成 1 个怪物
			if len(parts) != 2 {
				p.ReceiveChat(fmt.Sprintf("生成怪物失败，正确命令格式: @mob 怪物名"), common.ChatTypeSystem)
				return
			}
			c := curMap.GetNextCell(p.GetCell(), p.GetDirection(), 1)
			if c == nil || c.HasObject() {
				p.ReceiveChat(fmt.Sprintf("生成怪物失败"), common.ChatTypeSystem)
				return
			}
			mi := curMap.Env.GameDB.GetMonsterInfoByName(parts[1])
			if mi == nil {
				p.ReceiveChat(fmt.Sprintf("生成怪物失败，找不到怪物 %s", parts[1]), common.ChatTypeSystem)
				return
			}
			curMap.AddObject(NewMonster(curMap, c.Point, mi))
		case "RECALLMOB":
		case "RELOADDROPS":
		case "RELOADNPCS":
		case "GIVEGOLD":
		case "GIVEPEARLS":
		case "GIVECREDIT":
		case "GIVESKILL":
		case "FIND":
		case "LEAVEGUILD":
		case "CREATEGUILD":
		case "ALLOWTRADE":
		case "TRIGGER":
		case "RIDE":
		case "SETFLAG":
		case "LISTFLAGS":
		case "CLEARFLAGS":
		case "CLEARMOB":
		case "CHANGECLASS": //@changeclass [Player] [Class]
		case "DIE":
		case "HAIR":
		case "DECO": //TEST CODE
		case "ADJUSTPKPOINT":
		case "ADDINVENTORY":
		case "ADDSTORAGE":
		case "INFO": // @info
			if len(parts) != 1 {
				return
			}
			c := curMap.GetNextCell(p.GetCell(), p.GetDirection(), 1)
			if c == nil {
				return
			}
			c.Objects.Range(func(k, v interface{}) bool {
				o := v.(IMapObject)
				if o.GetRace() == common.ObjectTypeMonster {
					mo := o.(*Monster)
					p.ReceiveChat("--Monster Info--", common.ChatTypeSystem2)
					p.ReceiveChat(fmt.Sprintf("ID: %d, Name: %s", mo.ID, mo.Name), common.ChatTypeSystem2)
					p.ReceiveChat(fmt.Sprintf("Level: %d, Pos: %s", mo.Level, mo.GetPoint()), common.ChatTypeSystem2)
					p.ReceiveChat(fmt.Sprintf("HP: %d, MinDC: %d, MaxDC: %d", mo.HP, mo.MinDC, mo.MaxDC), common.ChatTypeSystem2)
				}
				if o.GetRace() == common.ObjectTypePlayer {
					po := o.(*Player)
					p.ReceiveChat("--Player Info--", common.ChatTypeSystem2)
					p.ReceiveChat(fmt.Sprintf("Name: %s, Level: %d, Pos: %s", po.Name, po.Level, po.GetPoint()), common.ChatTypeSystem2)
				}
				return true
			})
		case "CLEARQUESTS":
		case "SETQUEST":
		case "TOGGLETRANSFORM":
		case "CREATEMAPINSTANCE": //TEST CODE
		case "STARTCONQUEST":
		case "RESETCONQUEST":
		case "GATES":
		case "CHANGEFLAG":
		case "CHANGEFLAGCOLOUR":
		case "REVIVE":
		case "DELETESKILL":
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
	switch mirGridType {
	case common.MirGridTypeInventory:
		l := len(p.Inventory)
		if from > 0 && to > 0 && int(from) < l && int(to) < l {
			array := p.Inventory
			i := array[to]
			array[to] = array[from]
			array[from] = i
			msg.Success = true
		}
	case common.MirGridTypeStorage:
		// TODO
	case common.MirGridTypeTrade:
		// TODO
	case common.MirGridTypeRefine:
		// TODO
	}
	p.Enqueue(msg)
}

func (p *Player) StoreItem(from int32, to int32) {
	msg := &server.StoreItem{
		From:    from,
		To:      to,
		Success: false,
	}
	p.Enqueue(msg)
}

func (p *Player) DepositRefineItem(from int32, to int32) {

}

func (p *Player) RetrieveRefineItem(from int32, to int32) {

}

func (p *Player) RefineCancel() {

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

func (p *Player) TakeBackItem(from int32, to int32) {

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
	if l := len(p.Equipment); to < 0 || int(to) >= l {
		p.Enqueue(msg)
		return
	}
	switch mirGridType {
	case common.MirGridTypeInventory:
		index, item := p.GetUserItemByID(mirGridType, id)
		if item == nil {
			p.Enqueue(msg)
			return
		}
		p.Inventory[index] = p.Equipment[to]
		p.Equipment[to] = *item
	case common.MirGridTypeStorage:
		// TODO
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
	if l := len(p.Inventory); to < 0 || int(to) >= l {
		p.Enqueue(msg)
		return
	}
	switch mirGridType {
	case common.MirGridTypeInventory:
		index, item := p.GetUserItemByID(common.MirGridTypeEquipment, id)
		if item == nil {
			p.Enqueue(msg)
			return
		}
		invItem := p.Inventory[to]
		if invItem.ID == 0 {
			p.Inventory[to], p.Equipment[index] = p.Equipment[index], p.Inventory[to]
			break
		}
		for i := range p.Inventory[6:] {
			tmp := p.Inventory[6:][i]
			if tmp.ID != 0 {
				continue
			}
			p.Inventory[6:][i], p.Equipment[index] = p.Equipment[index], p.Inventory[6:][i]
			break
		}
	case common.MirGridTypeStorage:
		// TODO
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
	ph := &p.Health
	info := p.Map.Env.GameDB.GetItemInfoByID(int(item.ItemID))
	switch info.Type {
	case common.ItemTypePotion:
		switch info.Shape {
		case 0: // NormalPotion 一般药水
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
	case common.ItemTypeScroll:
	case common.ItemTypeBook:
	case common.ItemTypeScript:
	case common.ItemTypeFood:
	case common.ItemTypePets:
	case common.ItemTypeTransform: //Transforms
	default:
		p.Enqueue(msg)
		return
	}
	if item.Count > 1 {
		item.Count--
	} else {
		p.Inventory[index] = common.UserItem{}
	}
	p.RefreshBagWeight()
	msg.Success = true
	p.Enqueue(msg)
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
	obj := p.Map.Env.CreateDropItem(p.Map, userItem, 0)
	if dropMsg, ok := obj.Drop(p.GetPoint(), 1); !ok {
		p.ReceiveChat(dropMsg, common.ChatTypeSystem)
		return
	}
	if count >= userItem.Count {
		p.Inventory[index] = common.UserItem{}
	} else {
		p.Inventory[index].Count -= count
	}
	p.RefreshBagWeight()
	msg.Success = true
	p.Enqueue(msg)
}

func (p *Player) DropGold(gold uint64) {
	if p.Gold < gold {
		return
	}
	obj := p.Map.Env.CreateDropItem(p.Map, nil, gold)
	if dropMsg, ok := obj.Drop(p.GetPoint(), 3); !ok {
		p.ReceiveChat(dropMsg, common.ChatTypeSystem)
		return
	}
	p.Gold -= gold
	p.Enqueue(&server.LoseGold{Gold: uint32(gold)})
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
	c.Objects.Range(func(k, v interface{}) bool {
		if o, ok := v.(*Item); ok {
			if o.UserItem == nil {
				p.GainGold(o.Gold)
				items = append(items, o)
			} else {
				if p.GainItem(o.UserItem) {
					items = append(items, o)
				}
			}
		}
		return true
	})
	for i := range items {
		o := items[i]
		p.Map.DeleteObject(o)
		o.Broadcast(ServerMessage{}.ObjectRemove(o))
	}
}

func (p *Player) Inspect(id uint32) {
	o := p.Map.Env.GetPlayer(id)
	for i := range o.Equipment {
		item := p.Map.Env.GameDB.GetItemInfoByID(int(o.Equipment[i].ItemID))
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
	p.CurrentDirection = direction
	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Broadcast(ServerMessage{}.ObjectAttack(p, common.SpellNone, 0, 0))
	target := p.GetPoint().NextPoint(p.GetDirection(), 1)
	damageBase := p.GetAttackPower(int(p.MinDC), int(p.MaxDC)) // = the original damage from your gear (+ bonus from moonlight and darkbody)
	damageFinal := damageBase                                  // = the damage you're gonna do with skills added
	cell := p.Map.GetCell(target)
	if !cell.CanWalk() {
		return
	}
	cell.Objects.Range(func(k, v interface{}) bool {
		o := v.(IMapObject)
		if !o.IsAttackTarget(p) {
			return true
		}
		switch o.GetRace() {
		case common.ObjectTypePlayer:
			o.(*Player).Attacked(p, damageFinal, common.DefenceTypeAgility, false)
		case common.ObjectTypeMonster:
			o.(*Monster).Attacked(p, damageFinal, common.DefenceTypeAgility, false)
		}
		return true
	})
}

func (p *Player) RangeAttack(direction common.MirDirection, location common.Point, id uint32) {

}

func (p *Player) Harvest(direction common.MirDirection) {

}

func (p *Player) CallNPC(id uint32, key string) {
	npc := p.Map.Env.GetNPC(id)
	if npc == nil {
		return
	}
	say, err := npc.CallScript(p, key)
	if err != nil {
		log.Warnf("NPC 脚本执行失败: %d %s %s\n", id, key, err.Error())
	}

	p.Enqueue(ServerMessage{}.NPCResponse(replaceTemplates(npc, p, say)))

	// ProcessSpecial
	switch strings.ToUpper(key) {
	case "[@BUY]":
		sendBuyKey(p, npc)
	case "[@SELL]":
		p.Enqueue(&server.NPCSell{})
	case "[@BUYSELL]":
		sendBuyKey(p, npc)
		p.Enqueue(&server.NPCSell{})
	default:
		// TODO
	}
}

func sendBuyKey(p *Player, npc *NPC) {

	goods := []*common.UserItem{}

	// TODO: fix..
	// for _, name := range npc.Script.Goods {
	// 	item := p.Map.Env.GameDB.GetItemInfoByName(name)
	// 	if item != nil {
	// 		p.EnqueueItemInfo(item.ID)
	// 		goods = append(goods, p.Map.Env.NewUserItem(item))
	// 	}
	// }

	p.Enqueue(&server.NPCGoods{
		Goods: goods,
		Rate:  1.0,
		Type:  common.PanelTypeBuy,
	})
}

func (p *Player) TalkMonsterNPC(id uint32) {

}

func (p *Player) BuyItem(index uint64, count uint32, panelType common.PanelType) {

}

func (p *Player) CraftItem() {

}

func (p *Player) SellItem(id uint64, count uint32) {

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
	userMagic := p.GetMagic(spell)
	if userMagic == nil {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	info := p.Map.Env.GameDB.GetMagicInfoByID(userMagic.MagicID)
	cost := info.BaseCost + info.LevelCost*userMagic.Level
	if uint16(cost) > p.MP {
		p.Enqueue(ServerMessage{}.UserLocation(p))
		return
	}
	p.CurrentDirection = direction
	p.ChangeMP(-cost)
	target := p.Map.GetObjectInAreaByID(targetID, targetLocation)
	cast, targetID := p.UseMagic(spell, userMagic, target)
	p.Enqueue(ServerMessage{}.UserLocation(p))
	p.Enqueue(ServerMessage{}.Magic(spell, targetID, targetLocation, cast, userMagic.Level))
	p.Broadcast(ServerMessage{}.ObjectMagic(p, spell, targetID, targetLocation, cast, userMagic.Level))
}

func (p *Player) MagicKey(spell common.Spell, key uint8) {

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
