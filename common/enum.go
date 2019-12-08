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

type LightSetting uint8

const (
	Normal  LightSetting = 0
	Dawn                 = 1
	Day                  = 2
	Evening              = 3
	Night                = 4
)

type MirDirection uint8

const (
	Up        MirDirection = 0
	UpRight                = 1
	Right                  = 2
	DownRight              = 3
	Down                   = 4
	DownLeft               = 5
	Left                   = 6
	UpLeft                 = 7
)
