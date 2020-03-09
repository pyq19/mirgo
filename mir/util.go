package mir

import (
	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/ut"
)

// 随机方向
func RandomDirection() common.MirDirection {
	return common.MirDirection(ut.RandomInt(0, common.MirDirectionCount))
}

func MaxDistance(p1, p2 common.Point) int {
	return ut.MaxInt(ut.AbsInt(int(p1.X)-int(p2.X)), ut.AbsInt(int(p1.Y)-int(p2.Y)))
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
	return ut.AbsInt(int(a.X)-int(b.X)) <= i && ut.AbsInt(int(a.Y)-int(b.Y)) <= i
}

func InRangeXY(a common.Point, x, y, i int) bool {
	return ut.AbsInt(int(a.X)-x) <= i && ut.AbsInt(int(a.Y)-y) <= i
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
