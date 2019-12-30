package main

import "github.com/yenkeia/mirgo/common"

type MapObject struct {
	ID               uint32
	Name             string
	Map              *Map
	CurrentLocation  *common.Point
	CurrentDirection common.MirDirection
}
