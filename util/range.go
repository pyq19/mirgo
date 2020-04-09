package util

func InRange(a, b Point, i int) bool {
	return AbsInt(int(a.X)-int(b.X)) <= i && AbsInt(int(a.Y)-int(b.Y)) <= i
}

func InRangeXY(a Point, x, y, i int) bool {
	return AbsInt(int(a.X)-x) <= i && AbsInt(int(a.Y)-y) <= i
}
