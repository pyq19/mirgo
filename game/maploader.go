package game

import (
	"fmt"
	"io/ioutil"

	"github.com/yenkeia/mirgo/game/cm"
	"github.com/yenkeia/mirgo/game/util"
)

func LoadMap(filepath string) *Map {
	fileBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	v := DetectMapVersion(fileBytes)

	switch v {
	case 0:
		return GetMapV0(fileBytes)

	case 1:
		return GetMapV1(fileBytes)

	case 3:
		return GetMapV3(fileBytes)

	case 5:
		return GetMapV5(fileBytes)
	default:
		panic(fmt.Sprintf("map version not support! %d", int(v)))
	}

	return nil
}

func DetectMapVersion(input []byte) byte {
	//c# custom map format
	if (input[2] == 0x43) && (input[3] == 0x23) {
		return 100
	}
	//wemade mir3 maps have no title they just start with blank bytes
	if input[0] == 0 {
		return 5
	}
	//shanda mir3 maps start with title: (C) SNDA, MIR3.
	if (input[0] == 0x0F) && (input[5] == 0x53) && (input[14] == 0x33) {
		return 6
	}

	//wemades antihack map (laby maps) title start with: Mir2 AntiHack
	if (input[0] == 0x15) && (input[4] == 0x32) && (input[6] == 0x41) && (input[19] == 0x31) {
		return 4
	}

	//wemades 2010 map format i guess title starts with: Map 2010 Ver 1.0
	if (input[0] == 0x10) && (input[2] == 0x61) && (input[7] == 0x31) && (input[14] == 0x31) {
		return 1
	}

	//shanda's 2012 format and one of shandas(wemades) older formats share same header info, only difference is the filesize
	if (input[4] == 0x0F) && (input[18] == 0x0D) && (input[19] == 0x0A) {
		W := int(input[0] + (input[1] << 8))
		H := int(input[2] + (input[3] << 8))
		if len(input) > (52 + (W * H * 14)) {
			return 3
		}
		return 2
	}

	//3/4 heroes map format (myth/lifcos i guess)
	if (input[0] == 0x0D) && (input[1] == 0x4C) && (input[7] == 0x20) && (input[11] == 0x6D) {
		return 7
	}

	return 0
}

var (
	LowWallCell  = NewCell(cm.CellAttributeLowWall)
	HighWallCell = NewCell(cm.CellAttributeHighWall)
)

func GetMapV0(bytes []byte) *Map {
	offset := 0
	w := util.BytesToUint16(bytes[offset:])
	offset += 2
	h := util.BytesToUint16(bytes[offset:])
	width := int(w)
	height := int(h)

	m := NewMap(width, height, 0)

	var cell *Cell

	offset = 52
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			cell = nil

			if (util.BytesToUint16(bytes[offset:]) & 0x8000) != 0 {
				cell = HighWallCell //Can Fire Over.
			}

			offset += 2
			if (util.BytesToUint16(bytes[offset:]) & 0x8000) != 0 {
				cell = LowWallCell //Can't Fire Over.
			}

			offset += 2
			if (util.BytesToUint16(bytes[offset:]) & 0x8000) != 0 {
				cell = HighWallCell //No Floor Tile.
			}

			if cell == nil {
				cell = NewCell(cm.CellAttributeWalk)
			}

			point := cm.NewPoint(x, y)
			if cell.Attribute == cm.CellAttributeWalk {
				cell.Point = point
				m.SetCell(point, cell)
			}

			offset += 4
			if bytes[offset] > 0 {
				m.AddDoor(bytes[offset], point)
			}

			offset += 3 + 1

			// byte light = fileBytes[offSet++];

			// if (light >= 100 && light <= 119)
			// 	Cells[x, y].FishingAttribute = (sbyte)(light - 100);
		}
	}
	return m
}

func GetMapV1(bytes []byte) *Map {
	offset := 21
	w := util.BytesToUint16(bytes[offset:])
	offset += 2
	xor := util.BytesToUint16(bytes[offset:])
	offset += 2
	h := util.BytesToUint16(bytes[offset:])
	width := int(w ^ xor)
	height := int(h ^ xor)

	m := NewMap(width, height, 1)

	var cell *Cell

	offset = 54
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			cell = nil

			if (util.BytesToUint32(bytes[offset:])^0xAA38AA38)&0x20000000 != 0 {
				cell = HighWallCell
			}

			offset += 6
			if ((util.BytesToUint16(bytes[offset:]) ^ xor) & 0x8000) != 0 {
				cell = LowWallCell
			}

			if cell == nil {
				cell = NewCell(cm.CellAttributeWalk)
			}

			point := cm.NewPoint(x, y)
			if cell.Attribute == cm.CellAttributeWalk {
				cell.Point = point
				m.SetCell(point, cell)
			}

			offset += 2
			if bytes[offset] > 0 {
				m.AddDoor(bytes[offset], point)
			}

			offset += 5

			// byte light = fileBytes[offSet++];
			// if (light >= 100 && light <= 119)
			// 	Cells[x, y].FishingAttribute = (sbyte)(light - 100);
			offset += 1 + 1
		}
	}
	return m
}

func GetMapV3(bytes []byte) *Map {
	offset := 0
	w := util.BytesToUint16(bytes[offset:])
	offset += 2
	h := util.BytesToUint16(bytes[offset:])
	width := int(w)
	height := int(h)

	m := NewMap(width, height, 3)
	var cell *Cell

	offset = 52
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			cell = nil

			if (util.BytesToUint16(bytes[offset:]) & 0x8000) != 0 {
				cell = HighWallCell
			}

			offset += 2
			if (util.BytesToUint16(bytes[offset:]) & 0x8000) != 0 {
				cell = LowWallCell
			}

			offset += 2
			if (util.BytesToUint16(bytes[offset:]) & 0x8000) != 0 {
				cell = HighWallCell
			}

			if cell == nil {
				cell = NewCell(cm.CellAttributeWalk)
			}

			point := cm.NewPoint(x, y)
			if cell.Attribute == cm.CellAttributeWalk {
				cell.Point = point
				m.SetCell(point, cell)
			}

			offset += 2
			if bytes[offset] > 0 {
				m.AddDoor(bytes[offset], point)
			}
			offset += 12

			// byte light = fileBytes[offSet++];

			// if (light >= 100 && light <= 119)
			// 	Cells[x, y].FishingAttribute = (sbyte)(light - 100);

			offset += 17 + 1
		}
	}
	return m
}

func GetMapV5(bytes []byte) *Map {
	offset := 22
	w := util.BytesToUint16(bytes[offset:])
	offset += 2
	h := util.BytesToUint16(bytes[offset:])
	width := int(w)
	height := int(h)

	m := NewMap(width, height, 5)
	var cell *Cell

	offset = 28 + (3 * ((width / 2) + (width % 2)) * (height / 2))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			cell = nil

			if (bytes[offset] & 0x01) != 1 {
				cell = HighWallCell
			} else if (bytes[offset] & 0x02) != 2 {
				cell = LowWallCell
			} else {
				cell = NewCell(cm.CellAttributeWalk)
			}
			offset += 13

			// byte light = fileBytes[offSet++];

			// if (light >= 100 && light <= 119)
			// 	Cells[x, y].FishingAttribute = (sbyte)(light - 100);

			offset += 1

			cell.Point = cm.NewPoint(x, y)
			m.SetCell(cell.Point, cell)
		}
	}
	return m
}
