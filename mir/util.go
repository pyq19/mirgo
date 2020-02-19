package mir

import (
	"os"
	"path"
	"path/filepath"

	"github.com/yenkeia/mirgo/common"
	"github.com/yenkeia/mirgo/ut"
)

// 随机方向
func RandomDirection() common.MirDirection {
	return common.MirDirection(ut.RandomInt(0, common.MirDirectionCount))
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

func GetFiles(dir string, allow []string) []string {

	allowMap := map[string]bool{}
	if allow != nil {
		for _, v := range allow {
			allowMap[v] = true
		}
	}

	ret := []string{}
	filepath.Walk(dir, func(fpath string, f os.FileInfo, err error) error {
		if f == nil || f.IsDir() {
			return nil
		}

		ext := path.Ext(fpath)
		if allowMap[ext] {
			ret = append(ret, filepath.ToSlash(fpath))
		}

		return nil
	})

	return ret
}
