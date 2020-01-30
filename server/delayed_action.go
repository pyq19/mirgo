package main

import "time"

type DelayedType int

/*
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
*/

type DelayedAction struct {
	ID         uint32
	ActionTime time.Time
	Finish     bool
	Task       *Task
}
