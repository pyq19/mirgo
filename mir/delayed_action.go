package main

import "time"

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
	ID          uint32
	DelayedType DelayedType
	ActionTime  time.Time
	Finish      bool
	Task        *Task
}

func NewDelayedAction(id uint32, typ DelayedType, task *Task) *DelayedAction {
	return &DelayedAction{
		ID:          id,
		DelayedType: typ,
		ActionTime:  time.Now().Add(time.Millisecond * 500),
		Finish:      false,
		Task:        task,
	}
}
