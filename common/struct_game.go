package common

type SelectInfo struct {
	Index      uint32
	Name       string
	Level      uint16
	Class      MirClass
	Gender     MirGender
	LastAccess int64
}
