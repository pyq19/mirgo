package game

import (
	"container/list"
	"time"
)

type DelayedType int

const (
	DelayedTypeMagic DelayedType = iota
	DelayedTypeDamage
	DelayedTypeRangeDamage
	DelayedTypeSpawn
	DelayedTypeDie
	DelayedTypeRecall
	DelayedTypeMapMovement
	DelayedTypeMine
	DelayedTypeNPC
	DelayedTypePoison
	DelayedTypeDamageIndicator
)

type DelayedAction struct {
	DelayedType DelayedType
	ActionTime  time.Time
	CB          Callback
}

type Callback func()

type ActionList struct {
	List *list.List
}

func NewActionList() *ActionList {
	ret := &ActionList{}

	ret.List = list.New()

	return ret
}

func (lst *ActionList) PushAction(typ DelayedType, cb Callback) {
	lst.PushDelayAction(typ, 500, cb)
}

func (lst *ActionList) PushDelayAction(typ DelayedType, delay int, cb Callback) {
	lst.List.PushBack(&DelayedAction{
		DelayedType: typ,
		ActionTime:  time.Now().Add(time.Millisecond * time.Duration(delay)),
		CB:          cb,
	})
}

func (lst *ActionList) PushActionSuper(typ DelayedType, cb func(...interface{}), args ...interface{}) {
	lst.PushAction(typ, func() { cb(args...) })
}

func (lst *ActionList) Execute() {
	now := time.Now()
	for it := lst.List.Front(); it != nil; {
		action := it.Value.(*DelayedAction)
		if now.Before(action.ActionTime) {
			it = it.Next()
			continue
		}
		action.CB()
		tmp := it
		it = it.Next()
		lst.List.Remove(tmp)
	}
}
