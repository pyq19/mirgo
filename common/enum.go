package common

type MirGender uint8

const (
	Male   MirGender = 0
	Female MirGender = 1
)

type MirClass uint8

const (
	Warrior MirClass = iota
	Wizard
	Taoist
	Assassin
	Archer
)
