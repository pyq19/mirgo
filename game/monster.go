package game

import (
	"container/list"
	"fmt"
	"time"

	"github.com/yenkeia/mirgo/game/cm"
	"github.com/yenkeia/mirgo/game/proto/server"
	"github.com/yenkeia/mirgo/game/util"
)

type IBehavior interface {
	Process(dt time.Duration)
}

type BehaviroFactory func(id int, mon *Monster) IBehavior

var behaviorFactory BehaviroFactory

func SetMonsterBehaviorFactory(fac BehaviroFactory) {
	behaviorFactory = fac
}

// Monster ...
type Monster struct {
	MapObject
	Image         cm.Monster
	AI            int
	Behavior      IBehavior
	Effect        int
	Poison        cm.PoisonType
	Light         uint8
	Target        IMapObject
	Level         uint16
	PetLevel      uint16
	Experience    uint16
	HP            uint32
	MaxHP         uint32
	MinAC         uint16
	MaxAC         uint16
	MinMAC        uint16
	MaxMAC        uint16
	MinDC         uint16
	MaxDC         uint16
	MinMC         uint16
	MaxMC         uint16
	MinSC         uint16
	MaxSC         uint16
	Accuracy      uint8
	Agility       uint8
	MoveSpeed     uint16
	AttackSpeed   int32
	ArmourRate    float32
	DamageRate    float32
	ViewRange     int
	Master        *Player
	EXPOwner      *Player
	EXPOwnerTime  time.Time
	ActionList    *ActionList
	ActionTime    time.Time
	AttackTime    time.Time
	DeadTime      time.Time
	MoveTime      time.Time
	PoisonList    *PoisonList
	CurrentPoison cm.PoisonType
}

func (m *Monster) String() string {
	return fmt.Sprintf("Monster: %s, ID: %d, AI: %d, ptr: %p\n", m.Name, m.ID, m.AI, m)
}

// NewMonster ...
func NewMonster(mp *Map, p cm.Point, mi *cm.MonsterInfo) (m *Monster) {
	m = new(Monster)
	m.ID = env.NewObjectID()
	m.Map = mp
	m.Name = mi.Name
	m.NameColor = cm.Color{R: 255, G: 255, B: 255}
	m.Image = cm.Monster(mi.Image)
	m.AI = mi.AI
	m.Effect = mi.Effect
	m.Light = uint8(mi.Light)
	m.Target = nil
	m.Poison = cm.PoisonTypeNone
	m.CurrentLocation = p
	m.Direction = cm.RandomDirection()
	m.Dead = false
	m.Level = uint16(mi.Level)
	m.PetLevel = 0
	m.Experience = uint16(mi.Experience)
	m.HP = uint32(mi.HP)
	m.MaxHP = uint32(mi.HP)
	m.MinAC = uint16(mi.MinAC)
	m.MaxAC = uint16(mi.MaxAC)
	m.MinMAC = uint16(mi.MinMAC)
	m.MaxMAC = uint16(mi.MaxMAC)
	m.MinDC = uint16(mi.MinDC)
	m.MaxDC = uint16(mi.MaxDC)
	m.MinMC = uint16(mi.MinMC)
	m.MaxMC = uint16(mi.MaxMC)
	m.MinSC = uint16(mi.MinSC)
	m.MaxSC = uint16(mi.MaxSC)
	m.Accuracy = uint8(mi.Accuracy)
	m.Agility = uint8(mi.Agility)
	m.MoveSpeed = uint16(mi.MoveSpeed)
	m.AttackSpeed = int32(mi.AttackSpeed)
	m.ArmourRate = 1.0
	m.DamageRate = 1.0
	m.ActionList = NewActionList()
	now := time.Now()
	m.ActionTime = now
	m.MoveTime = now
	m.ViewRange = mi.ViewRange
	m.Behavior = behaviorFactory(m.AI, m)
	m.PoisonList = NewPoisonList()
	m.CurrentPoison = cm.PoisonTypeNone
	return m
}

func (i *Monster) GetMap() *Map {
	return i.Map
}

func (m *Monster) GetID() uint32 {
	return m.ID
}

func (m *Monster) AddPlayerCount(n int) {
	m.PlayerCount += n
	switch m.PlayerCount {
	case 1:
		m.Map.AddActiveObj(m)
	case 0:
		m.Map.DelActiveObj(m)
	}
}

func (m *Monster) GetPlayerCount() int {
	return m.PlayerCount
}

func (m *Monster) GetName() string {
	return m.Name
}

func (m *Monster) GetLevel() int {
	return int(m.Level)
}

func (m *Monster) GetRace() cm.ObjectType {
	return cm.ObjectTypeMonster
}

func (m *Monster) GetPoint() cm.Point {
	return m.CurrentLocation
}

func (m *Monster) GetCell() *Cell {
	return m.Map.GetCell(m.CurrentLocation)
}

func (m *Monster) GetDirection() cm.MirDirection {
	return m.Direction
}

func (m *Monster) GetHP() int {
	return int(m.HP)
}

func (m *Monster) SetHP(amount uint32) {
	if m.HP == amount {
		return
	}
	if amount <= m.MaxHP {
		m.HP = amount
	} else {
		m.HP = m.MaxHP
	}
	if !m.Dead && m.HP == 0 {
		m.Die()
	}
	m.BroadcastHealthChange()
}

func (m *Monster) GetMaxHP() int {
	return int(m.MaxHP)
}

func (m *Monster) GetPercentHealth() uint8 {
	return uint8((float32(m.GetHP()) / float32(m.GetMaxHP()) * 100))
}

func (m *Monster) BroadcastHealthChange() {
	IMapObject_BroadcastHealthChange(m)
}

func (m *Monster) BroadcastInfo() {
	m.Broadcast(ServerMessage{}.ObjectMonster(m))
}

func (m *Monster) Spawned() {
	IMapObject_Spawned(m)
}

func (m *Monster) GetBaseStats() BaseStats {
	return BaseStats{
		MinAC:    m.MinAC,
		MaxAC:    m.MaxAC,
		MinMAC:   m.MinMAC,
		MaxMAC:   m.MaxMAC,
		MinDC:    m.MinDC,
		MaxDC:    m.MaxDC,
		MinMC:    m.MinMC,
		MaxMC:    m.MaxMC,
		MinSC:    m.MinSC,
		MaxSC:    m.MaxSC,
		Accuracy: m.Accuracy,
		Agility:  m.Agility,
	}
}

func (m *Monster) AddBuff(buff *Buff) {}

func (m *Monster) ApplyPoison(p *Poison, caster IMapObject) {

	ignoreDefence := false

	if p.Owner != nil && p.Owner.IsAttackTarget(m) {
		m.Target = p.Owner
	}
	// TODO
	/*
	  if (Master != null && p.Owner != null && p.Owner.Race == ObjectType.Player && p.Owner != Master)
	  {
	      if (Envir.Time > Master.BrownTime && Master.PKPoints < 200)
	          p.Owner.BrownTime = Envir.Time + Settings.Minute;
	  }

	*/

	if !ignoreDefence && (p.PType == cm.PoisonTypeGreen) {
		armour := m.GetDefencePower(int(m.MinMAC), int(m.MaxMAC))

		if p.Value < armour {

			p.PType = cm.PoisonTypeNone
		} else {
			p.Value -= armour
		}
	}

	if p.PType == cm.PoisonTypeNone {
		return
	}
	// TODO
	/*
	  for (int i = 0; i < PoisonList.Count; i++)
	  {
	      if (PoisonList[i].PType != p.PType) continue;
	      if ((PoisonList[i].PType == PoisonType.Green) && (PoisonList[i].Value > p.Value)) return;//cant cast weak poison to cancel out strong poison
	      if ((PoisonList[i].PType != PoisonType.Green) && ((PoisonList[i].Duration - PoisonList[i].Time) > p.Duration)) return;//cant cast 1 second poison to make a 1minute poison go away!
	      if (p.PType == PoisonType.DelayedExplosion) return;
	      if ((PoisonList[i].PType == PoisonType.Frozen) || (PoisonList[i].PType == PoisonType.Slow) || (PoisonList[i].PType == PoisonType.Paralysis) || (PoisonList[i].PType == PoisonType.LRParalysis)) return;//prevents mobs from being perma frozen/slowed
	      PoisonList[i] = p;
	      return;
	  }

	  if (p.PType == PoisonType.DelayedExplosion)
	  {
	      ExplosionInflictedTime = Envir.Time + 4000;
	      Broadcast(new S.ObjectEffect { ObjectID = ObjectID, Effect = SpellEffect.DelayedExplosion });
	  }
	*/
	m.PoisonList.List.PushBack(p)
}

func (m *Monster) Broadcast(msg interface{}) {
	m.Map.BroadcastP(m.CurrentLocation, msg, nil)
}

// Spawn 怪物生成
func (m *Monster) Spawn() {
	msg, ok := m.Map.AddObject(m)
	if !ok {
		log.Warnln(msg)
		return
	}
	// RefreshAll();
	// SetHP(MaxHP);
	m.Broadcast(ServerMessage{}.Object(m))
}

func (m *Monster) BroadcastDamageIndicator(typ cm.DamageType, dmg int) {
	m.Broadcast(ServerMessage{}.DamageIndicator(int32(dmg), typ, m.GetID()))
}

func (m *Monster) IsDead() bool {
	return m.Dead
}

func (m *Monster) IsUndead() bool {
	return false
}

func (m *Monster) IsBlocking() bool {
	return !m.IsDead()
}

func (m *Monster) IsSkeleton() bool {
	return false
}

func (m *Monster) IsHidden() bool {
	return false
}

func (m *Monster) IsAttackTargetMonster(attacker *Monster) bool {
	if attacker == m {
		return false
	}

	if m.AI == 6 || m.AI == 58 {
		return false
	}

	if attacker.AI == 6 {
		if m.AI != 1 && m.AI != 2 && m.AI != 3 { //Not Dear/Hen/Tree/Pets or Red Master
			return true
		}
	} else if attacker.AI == 58 {
		if m.AI != 1 && m.AI != 2 && m.AI != 3 {
			return true
		}
	}
	return false
}

func (m *Monster) IsAttackTargetPlayer(attacker *Player) bool {
	if m.IsDead() {
		return false
	}
	if m.Master == nil {
		return true
	}
	if attacker.AMode == cm.AttackModePeace {
		return false
	}
	if m.Master == attacker {
		return attacker.AMode == cm.AttackModeAll
	}
	if m.Master.GetRace() == cm.ObjectTypePlayer { // TODO && (attacker.InSafeZone || InSafeZone) {
		return false
	}
	/*
		switch attacker.AMode {
			case cm.AttackModeGroup:
				return Master.GroupMembers == null || !Master.GroupMembers.Contains(attacker);
			case cm.AttackModeGuild:
				{
					if (!(Master is PlayerObject)) return false;
					PlayerObject master = (PlayerObject)Master;
					return master.MyGuild == null || master.MyGuild != attacker.MyGuild;
				}
			case cm.AttackModeEnemyGuild:
				{
					if (!(Master is PlayerObject)) return false;
					PlayerObject master = (PlayerObject)Master;
					return (master.MyGuild != null && attacker.MyGuild != null) && master.MyGuild.IsEnemy(attacker.MyGuild);
				}
			case cm.AttackModeRedBrown:
				return Master.PKPoints >= 200 || Envir.Time < Master.BrownTime;
			default:
				return true;
		}
	*/
	return true
}

func (m *Monster) IsAttackTarget(attacker IMapObject) bool {

	switch attacker.(type) {
	case *Monster:
		return m.IsAttackTargetMonster(attacker.(*Monster))
	case *Player:
		return m.IsAttackTargetPlayer(attacker.(*Player))
	}
	return true
}

func (m *Monster) IsFriendlyTarget(attacker IMapObject) bool {
	return false
}

func (m *Monster) CanMove() bool {
	return time.Now().After(m.MoveTime)
}

func (m *Monster) CanAttack() bool {
	now := time.Now()
	if m.IsDead() {
		return false
	}
	return now.After(m.AttackTime)
}

// InAttackRange 是否在怪物攻击范围内
func (m *Monster) InAttackRange() bool {
	if m.Target.GetMap() != m.GetMap() {
		return false
	}
	return !m.Target.GetPoint().Equal(m.CurrentLocation) && cm.InRange(m.CurrentLocation, m.Target.GetPoint(), 1)
}

// Process 怪物定时轮询
func (m *Monster) Process(dt time.Duration) {
	if m.Target != nil &&
		(m.Target.GetMap() != m.GetMap() || !m.Target.IsAttackTarget(m) || !cm.InRange(m.CurrentLocation, m.Target.GetPoint(), DataRange)) {
		m.Target = nil
	}

	now := time.Now()

	if m.IsDead() && m.DeadTime.Before(now) {
		m.Map.DeleteObject(m)
		m.Broadcast(&server.ObjectRemove{ObjectID: m.GetID()})
		return
	}

	if m.EXPOwner != nil && now.After(m.EXPOwnerTime) {
		m.EXPOwner = nil
	}

	m.Behavior.Process(dt)

	m.ProcessBuffs()
	m.ProcessRegan()
	m.ProcessPoison()

	m.ActionList.Execute()
}

// ProcessBuffs 处理怪物增益效果
func (m *Monster) ProcessBuffs() {

}

// ProcessRegan 怪物自身回血
func (m *Monster) ProcessRegan() {

}

// ProcessPoison 处理怪物中毒效果
func (m *Monster) ProcessPoison() {
	if m.IsDead() {
		return
	}
	ptype := cm.PoisonTypeNone
	l := m.PoisonList.List
	var next *list.Element
	for e := l.Front(); e != nil; e = next {
		next = e.Next()
		poison := e.Value.(*Poison)
		if poison.Owner == nil || poison.TickCnt > poison.TickNum {
			l.Remove(e)
			continue
		}
		// log.Debugln("----")
		// log.Debugln(time.Now())
		// log.Debugln(poison.TickTime)
		// log.Debugln("----")
		if time.Now().After(poison.TickTime) {
			poison.TickTime = poison.TickTime.Add(poison.TickSpeed)
			poison.TickCnt++

			if poison.PType == cm.PoisonTypeGreen || poison.PType == cm.PoisonTypeBleeding {

				// TODO
				// if (m.EXPOwner == nil || m.EXPOwner.Dead) {
				// 	EXPOwner = poison.Owner;
				// 	EXPOwnerTime = Envir.Time + EXPOwnerDelay;
				// } else if (m.EXPOwner == poison.Owner) {
				// 	EXPOwnerTime = Envir.Time + EXPOwnerDelay;
				// }

				if poison.PType == cm.PoisonTypeBleeding {
					m.Broadcast(&server.ObjectEffect{ObjectID: m.GetID(), Effect: cm.SpellEffectBleeding, EffectType: 0})
				}

				m.ChangeHP(-poison.Value)
				// if (PoisonStopRegen) {	// 停止回血
				// 	RegenTime = Envir.Time + RegenDelay;
				// }
			}

			// TODO
			// if (poison.PType == PoisonType.DelayedExplosion)
		}

		switch poison.PType {
		case cm.PoisonTypeRed:
			m.ArmourRate -= 0.5
		case cm.PoisonTypeStun:
			m.DamageRate += 0.5
		case cm.PoisonTypeSlow:
			m.MoveSpeed += 100
			m.AttackSpeed += 100
			/*
				if poison.Time >= poison.Duration {
					m.MoveSpeed = Info.MoveSpeed
					m.AttackSpeed = Info.AttackSpeed
				}
			*/
		}
		ptype |= poison.PType
	}
	if ptype == m.CurrentPoison {
		return
	}
	m.CurrentPoison = ptype
	m.Broadcast(&server.ObjectPoisoned{ObjectID: m.GetID(), Poison: ptype})
}

// GetDefencePower 获取防御值
func (m *Monster) GetDefencePower(min, max int) int {
	if min < 0 {
		min = 0
	}
	if min > max {
		max = min
	}
	return util.RandomInt(min, max)
}

// GetAttackPower 获取攻击值
func (m *Monster) GetAttackPower(min, max int) int {
	if min < 0 {
		min = 0
	}
	if min > max {
		max = min
	}
	// TODO luck
	return util.RandomInt(min, max+1)
}

// Die ...
func (m *Monster) Die() {
	if m.IsDead() {
		return
	}

	m.HP = 0
	m.Dead = true
	m.DeadTime = time.Now().Add(5 * time.Second)

	m.Broadcast(ServerMessage{}.ObjectDied(m.GetID(), m.GetDirection(), m.GetPoint()))
	// EXPOwner.WinExp(Experience, Level);

	if m.EXPOwner != nil && m.Master == nil && m.EXPOwner.GetRace() == cm.ObjectTypePlayer {
		m.EXPOwner.WinExp(int(m.Experience), int(m.Level))
		// PlayerObject playerObj = (PlayerObject)EXPOwner;
		// playerObj.CheckGroupQuestKill(Info);
		// m.EXPOwner.CheckGroupQuestKill(Info)
	}

	m.Drop()
}

// ChangeHP 怪物改变血量 amount 可以是负数(扣血)
func (m *Monster) ChangeHP(amount int) {
	if m.IsDead() {
		return
	}
	value := int(m.HP) + amount
	if value == int(m.HP) {
		return
	}
	if value <= 0 {
		m.Die()
		m.HP = 0
	} else {
		m.HP = uint32(value)
	}
	percent := uint8(float32(m.HP) / float32(m.MaxHP) * 100)
	log.Debugf("monster HP: %d, MaxHP: %d, percent: %d\n", m.HP, m.MaxHP, percent)
	m.Broadcast(ServerMessage{}.ObjectHealth(m.GetID(), percent, 5))
}

// Attacked 被攻击
func (m *Monster) Attacked(attacker IMapObject, damage int, defenceType cm.DefenceType, damageWeapon bool) int {
	if m.Target == nil && attacker.IsAttackTarget(m) {
		m.Target = attacker
	}
	armour := 0
	switch defenceType {
	case cm.DefenceTypeACAgility:
		if util.RandomInt(0, int(m.Agility)) > int(attacker.GetBaseStats().Accuracy) {
			m.BroadcastDamageIndicator(cm.DamageTypeMiss, 0)
			return 0
		}
		armour = m.GetDefencePower(int(m.MinAC), int(m.MaxAC))
	case cm.DefenceTypeAC:
		armour = m.GetDefencePower(int(m.MinAC), int(m.MaxAC))
	case cm.DefenceTypeMACAgility:
		if util.RandomInt(0, int(m.Agility)) > int(attacker.GetBaseStats().Accuracy) {
			m.BroadcastDamageIndicator(cm.DamageTypeMiss, 0)
			return 0
		}
		armour = m.GetDefencePower(int(m.MinMAC), int(m.MaxMAC))
	case cm.DefenceTypeMAC:
		armour = m.GetDefencePower(int(m.MinMAC), int(m.MaxMAC))
	case cm.DefenceTypeAgility:
		if util.RandomInt(0, int(m.Agility)) > int(attacker.GetBaseStats().Accuracy) {
			m.BroadcastDamageIndicator(cm.DamageTypeMiss, 0)
			return 0
		}
	}
	armour = int(float32(armour) * m.ArmourRate)
	damage = int(float32(damage) * m.DamageRate)
	value := damage - armour
	log.Debugf("attacker damage: %d, monster armour: %d\n", damage, armour)
	if value <= 0 {
		m.BroadcastDamageIndicator(cm.DamageTypeMiss, 0)
		return 0
	}

	switch attacker := attacker.(type) {
	case *Monster:
		if attacker.AI == 6 || attacker.AI == 58 {
			m.EXPOwner = nil
		} else if attacker.Master != nil {
			if attacker.GetMap() != attacker.Master.GetMap() || !cm.InRange(attacker.CurrentLocation, attacker.Master.CurrentLocation, DataRange) {
				m.EXPOwner = nil
			} else {
				if m.EXPOwner == nil || m.EXPOwner.Dead {
					m.EXPOwner = attacker.Master
				}
				if m.EXPOwner == attacker.Master {
					m.EXPOwnerTime = time.Now().Add(5 * time.Second)
				}
			}
		}
	case *Player:
		if m.EXPOwner == nil || m.EXPOwner.Dead {
			m.EXPOwner = attacker
		}
		if m.EXPOwner == attacker {
			m.EXPOwnerTime = time.Now().Add(5 * time.Second)
		}
	}

	// TODO 还有很多没做
	m.Broadcast(ServerMessage{}.ObjectStruck(m, attacker.GetID()))
	m.BroadcastDamageIndicator(cm.DamageTypeHit, value)
	m.ChangeHP(-value)
	// log.Debugln("monster->", m)
	// log.Debugln("attacker->", attacker.(*Monster))

	return 0
}

// Drop 怪物掉落物品
func (m *Monster) Drop() {
	dropInfos, ok := data.DropInfoMap[m.Name]
	if !ok {
		return
	}
	mapItems := make([]*Item, 0)
	for _, drop := range dropInfos {
		// FIXME 任务还没做 任务相关的物品就不掉落
		if drop.QuestRequired {
			continue
		}
		if util.RandomNext(drop.High) > drop.Low {
			continue
		}
		if drop.ItemName == "Gold" {
			mapItems = append(mapItems, NewGold(m, uint64(drop.Count)))
			continue
		}
		info := data.GetItemInfoByName(drop.ItemName)
		if info == nil {
			continue
		}

		mapItems = append(mapItems, NewItem(m, env.NewUserItem(info)))
	}
	for i := range mapItems {
		if msg, ok := mapItems[i].Drop(m.GetPoint(), 3); !ok {
			log.Warnln(msg)
		}
	}
}

// Walk 移动，成功返回 true
func (m *Monster) Walk(dir cm.MirDirection) bool {
	if !m.CanMove() {
		return false
	}

	dest := m.CurrentLocation.NextPoint(dir, 1)
	destcell := m.Map.GetCell(dest)

	if destcell != nil && destcell.objects != nil {
		for _, o := range destcell.objects {
			if o.IsBlocking() || m.GetRace() == cm.ObjectTypeCreature {
				return false
			}
		}
	} else {
		return false
	}

	m.Map.GetCell(m.CurrentLocation).DeleteObject(m)
	destcell.AddObject(m)

	oldpos := m.CurrentLocation

	m.Direction = dir
	m.CurrentLocation = dest
	m.UpdateInSafeZone()

	m.WalkNotify(oldpos, destcell.Point)

	m.MoveTime = m.MoveTime.Add(time.Duration(int64(m.MoveSpeed)) * time.Millisecond)

	m.Broadcast(&server.ObjectWalk{
		ObjectID:  m.GetID(),
		Direction: dir,
		Location:  dest,
	})

	return true
}

func (m *Monster) WalkNotify(from, to cm.Point) {
	cells := m.Map.CalcDiff(from, to, DataRange)
	for c, isadd := range cells.M {
		if isadd {
			for _, o := range c.objects {
				switch o.(type) {
				case *Player:
					m.AddPlayerCount(1)
					o.(*Player).Enqueue(ServerMessage{}.Object(m))
				}
			}
		} else {
			for _, o := range c.objects {
				switch o.(type) {
				case *Player:
					m.AddPlayerCount(-1)
					o.(*Player).Enqueue(ServerMessage{}.ObjectRemove(m))
				}
			}
		}

	}
}

func (m *Monster) Turn(dir cm.MirDirection) {
	if !m.CanMove() {
		return
	}
	m.Direction = dir

	m.Broadcast(&server.ObjectTurn{
		ObjectID:  m.GetID(),
		Direction: dir,
		Location:  m.CurrentLocation,
	})

	m.UpdateInSafeZone()

	// TODO:
	// InSafeZone = CurrentMap.GetSafeZone(CurrentLocation) != null

	// Cell cell = CurrentMap.GetCell(CurrentLocation);
	// for (int i = 0; i < cell.Objects.Count; i++)
	// {
	//     if (cell.Objects[i].Race != ObjectType.Spell) continue;
	//     SpellObject ob = (SpellObject)cell.Objects[i];

	//     ob.ProcessSpell(this);
	//     //break;
	// }

}

func ObjectBack(m IMapObject) cm.Point {
	return m.GetPoint().NextPoint(m.GetDirection(), 1)
}

func (m *Monster) Attack() {
	if m.Target == nil {
		return
	}
	if !m.Target.IsAttackTarget(m) {
		m.Target = nil
		return
	}
	if p, ok := m.Target.(*Player); ok && p.GameStage != GAME {
		log.Debugf("monster[%s] lose target[%s]\n", m.Name, p.Name)
		m.Target = nil
		return
	}
	log.Debugf("Monster[%s]AI[%d] Attack [%s]\n", m.Name, m.AI, m.Target.GetName())
	m.Direction = cm.DirectionFromPoint(m.CurrentLocation, m.Target.GetPoint())
	m.Broadcast(ServerMessage{}.ObjectAttack(m, cm.SpellNone, 0, 0))
	now := time.Now()
	// ActionTime = Envir.Time + 300;
	m.AttackTime = now.Add(time.Duration(m.AttackSpeed) * time.Millisecond)
	damage := m.GetAttackPower(int(m.MinDC), int(m.MaxDC))
	if damage <= 0 {
		return
	}
	m.Target.Attacked(m, damage, cm.DefenceTypeAgility, false)
}

func (m *Monster) MoveTo(location cm.Point) {
	if m.CurrentLocation.Equal(location) {
		return
	}
	inRange := cm.InRange(location, m.CurrentLocation, 1)
	if inRange {
		cell := m.Map.GetCell(location)
		if cell == nil || !cell.IsValid() {
			return
		}
		for _, o := range cell.objects {
			if o.IsBlocking() {
				return
			}
		}
	}
	dir := cm.DirectionFromPoint(m.CurrentLocation, location)
	if m.Walk(dir) {
		return
	}
	switch util.RandomNext(2) { //No favour
	case 0:
		for i := 0; i < 7; i++ {
			dir = cm.NextDirection(dir)
			if m.Walk(dir) {
				return
			}
		}
	default:
		for i := 0; i < 7; i++ {
			dir = cm.PreviousDirection(dir)
			if m.Walk(dir) {
				return
			}
		}
	}
}

// FindTarget 怪物寻找攻击目标
func (m *Monster) FindTarget() {
	m.Map.RangeObject(m.CurrentLocation, m.ViewRange, func(o IMapObject) bool {

		if o == m {
			return true
		}

		switch o.GetRace() {
		case cm.ObjectTypeMonster:
			if !o.IsAttackTarget(m) {
				return true
			}
			// if (ob.Hidden && (!CoolEye || Level < ob.Level)) continue;
			m.Target = o

		case cm.ObjectTypePlayer:

			if !o.IsAttackTarget(m) { // continue
				return true
			}

			// TODO:
			// if (playerob.GMGameMaster || ob.Hidden && (!CoolEye || Level < ob.Level) || Envir.Time < HallucinationTime) continue;

			m.Target = o

			return false
		}

		return true
	})
}

func (m *Monster) CheckStacked() bool {
	cell := m.Map.GetCell(m.CurrentLocation)
	if cell != nil && cell.objects != nil {
		for _, o := range cell.objects {
			if o == m || o.IsBlocking() {
				continue
			}
			break
		}
	}

	return false
}

// PetRecall 宠物传送回玩家身边
func (m *Monster) PetRecall() {
	log.Debugln("PetRecall", m.GetID())
}

func (m *Monster) CompleteAttack(target IMapObject, damage int, def cm.DefenceType) {
	if target == nil || !target.IsAttackTarget(m) || target.GetMap() != m.GetMap() {
		return
	}

	target.Attacked(m, damage, def, false)
}
