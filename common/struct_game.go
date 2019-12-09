package common

type Point struct {
	X uint32
	Y uint32
}

type SelectInfo struct {
	Index      uint32
	Name       string
	Level      uint16
	Class      MirClass
	Gender     MirGender
	LastAccess int64
}

type UserItem struct {
	UniqueID       uint64
	ItemIndex      int32
	CurrentDura    uint16
	MaxDura        uint16
	Count          uint32
	AC             uint8
	MAC            uint8
	DC             uint8
	MC             uint8
	SC             uint8
	Accuracy       uint8
	Agility        uint8
	HP             uint8
	MP             uint8
	AttackSpeed    int8
	Luck           int8
	SoulBoundId    uint32
	Bools          uint8
	Strong         uint8
	MagicResist    uint8
	PoisonResist   uint8
	HealthRecovery uint8
	ManaRecovery   uint8
	PoisonRecovery uint8
	CriticalRate   uint8
	CriticalDamage uint8
	Freezing       uint8
	PoisonAttack   uint8
}
