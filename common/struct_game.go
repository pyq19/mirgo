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

