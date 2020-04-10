package cm

import "github.com/yenkeia/mirgo/game/util"

// 随机方向
func RandomDirection() MirDirection {
	return MirDirection(util.RandomInt(0, MirDirectionCount))
}

func MaxDistance(p1, p2 Point) int {
	return util.MaxInt(util.AbsInt(int(p1.X)-int(p2.X)), util.AbsInt(int(p1.Y)-int(p2.Y)))
}

func NextDirection(d MirDirection) MirDirection {
	switch d {
	case MirDirectionUp:
		return MirDirectionUpRight
	case MirDirectionUpRight:
		return MirDirectionRight
	case MirDirectionRight:
		return MirDirectionDownRight
	case MirDirectionDownRight:
		return MirDirectionDown
	case MirDirectionDown:
		return MirDirectionDownLeft
	case MirDirectionDownLeft:
		return MirDirectionLeft
	case MirDirectionLeft:
		return MirDirectionUpLeft
	case MirDirectionUpLeft:
		return MirDirectionUp
	default:
		return d
	}
}

func PreviousDirection(d MirDirection) MirDirection {

	switch d {
	case MirDirectionUp:
		return MirDirectionUpLeft
	case MirDirectionUpRight:
		return MirDirectionUp
	case MirDirectionRight:
		return MirDirectionUpRight
	case MirDirectionDownRight:
		return MirDirectionRight
	case MirDirectionDown:
		return MirDirectionDownRight
	case MirDirectionDownLeft:
		return MirDirectionDown
	case MirDirectionLeft:
		return MirDirectionDownLeft
	case MirDirectionUpLeft:
		return MirDirectionLeft
	default:
		return d
	}
}

// FacingEachOther 判断两个玩家是否面对面
func FacingEachOther(ad MirDirection, ap Point, bd MirDirection, bp Point) bool {
	if bd > 3 {
		bd -= 4
	} else {
		bd += 4
	}
	if ad != bd {
		return false
	}
	return ap.NextPoint(ad, 1).Equal(bp)
}

func DirectionFromPoint(source, dest Point) MirDirection {
	if source.X < dest.X {
		if source.Y < dest.Y {
			return MirDirectionDownRight
		}
		if source.Y > dest.Y {
			return MirDirectionUpRight
		}
		return MirDirectionRight
	}
	if source.X > dest.X {
		if source.Y < dest.Y {
			return MirDirectionDownLeft
		}
		if source.Y > dest.Y {

			return MirDirectionUpLeft
		}
		return MirDirectionLeft
	}
	if source.Y < dest.Y {
		return MirDirectionDown
	} else {
		return MirDirectionUp
	}
}
