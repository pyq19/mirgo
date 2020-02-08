package main

import (
	"io/ioutil"
	"math/rand"
	"sync"

	"github.com/yenkeia/mirgo/common"
)

func GetMapBytes(mapAbsPath string) []byte {
	fileBytes, err := ioutil.ReadFile(mapAbsPath)
	if err != nil {
		panic(err)
	}
	return fileBytes
}

func GetMapV1(bytes []byte) *Map {
	offset := 21
	w := common.BytesToUint16(bytes[offset : offset+2])
	offset += 2
	xor := common.BytesToUint16(bytes[offset : offset+2])
	offset += 2
	h := common.BytesToUint16(bytes[offset : offset+2])
	width := int(w ^ xor)
	height := int(h ^ xor)

	m := NewMap(width, height)

	offset = 54
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			p := common.Point{X: uint32(i), Y: uint32(j)}
			c := new(Cell)
			c.Map = m
			c.Point = p
			c.Objects = new(sync.Map)
			if (common.BytesToUint32(bytes[offset:offset+4])^0xAA38AA38)&0x20000000 != 0 {
				c.Attribute = common.CellAttributeHighWall
			}
			if ((common.BytesToUint16(bytes[offset+6:offset+8]) ^ xor) & 0x8000) != 0 {
				c.Attribute = common.CellAttributeLowWall
			}
			if c.Attribute == common.CellAttributeWalk {
				m.SetCell(p, c)
			}
			offset += 15
		}
	}
	m.Width = width
	m.Height = height
	return m
}

func newAOI(m *Map, width int, height int) (aoi *AOIManager) {
	cntX := width / 20
	cntY := height / 20
	if width < 20 {
		cntX = 1
	}
	if height < 20 {
		cntY = 1
	}
	aoi = NewAOIManager(m, 0, width, cntX, 0, height, cntY)
	return
}

func AbsInt(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// 随机 [low, high]
func RandomInt(low int, high int) int {
	if low == high {
		return low
	}

	return rand.Intn(high-low+1) + low
}

// c# random.next [0, high)
func RandomNext(high int) int {
	return RandomInt(0, high-1)
}

func RandomString(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		b := rand.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// 随机方向
func RandomDirection() common.MirDirection {
	return common.MirDirection(RandomInt(0, common.MirDirectionCount))
}

func NextDirection(d common.MirDirection) common.MirDirection {
	switch d {
	case common.MirDirectionUp:
		return common.MirDirectionUpRight
	case common.MirDirectionUpRight:
		return common.MirDirectionRight
	case common.MirDirectionRight:
		return common.MirDirectionDownRight
	case common.MirDirectionDownRight:
		return common.MirDirectionDown
	case common.MirDirectionDown:
		return common.MirDirectionDownLeft
	case common.MirDirectionDownLeft:
		return common.MirDirectionLeft
	case common.MirDirectionLeft:
		return common.MirDirectionUpLeft
	case common.MirDirectionUpLeft:
		return common.MirDirectionUp
	default:
		return d
	}
}

func PreviousDirection(d common.MirDirection) common.MirDirection {

	switch d {
	case common.MirDirectionUp:
		return common.MirDirectionUpLeft
	case common.MirDirectionUpRight:
		return common.MirDirectionUp
	case common.MirDirectionRight:
		return common.MirDirectionUpRight
	case common.MirDirectionDownRight:
		return common.MirDirectionRight
	case common.MirDirectionDown:
		return common.MirDirectionDownRight
	case common.MirDirectionDownLeft:
		return common.MirDirectionDown
	case common.MirDirectionLeft:
		return common.MirDirectionDownLeft
	case common.MirDirectionUpLeft:
		return common.MirDirectionLeft
	default:
		return d
	}
}

func InRange(a, b common.Point, i int) bool {
	return AbsInt(int(a.X-b.X)) <= i && AbsInt(int(a.Y-b.Y)) <= i
}

func DirectionFromPoint(source, dest common.Point) common.MirDirection {
	if source.X < dest.X {
		if source.Y < dest.Y {
			return common.MirDirectionDownRight
		}
		if source.Y > dest.Y {
			return common.MirDirectionUpRight
		}
		return common.MirDirectionRight
	}
	if source.X > dest.X {
		if source.Y < dest.Y {
			return common.MirDirectionDownLeft
		}
		if source.Y > dest.Y {

			return common.MirDirectionUpLeft
		}
		return common.MirDirectionLeft
	}
	if source.Y < dest.Y {
		return common.MirDirectionDown
	} else {
		return common.MirDirectionUp
	}
}
