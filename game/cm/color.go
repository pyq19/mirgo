package cm

import "github.com/yenkeia/mirgo/game/util"

type Color struct {
	R uint8
	G uint8
	B uint8
}

func NewColor(c uint32) Color {
	return Color{
		R: uint8((c >> 16) & 0xff),
		G: uint8((c >> 8) & 0xff),
		B: uint8((c >> 0) & 0xff),
	}
}

func (c Color) ToInt32() int32 {
	return int32(c.ToUint32())
}

func (c Color) ToUint32() uint32 {
	return util.BytesToUint32([]uint8{c.R, c.G, c.B, 255})
}

var (
	ColorWhite       = NewColor(0xFFFFFFFF)
	ColorDeepSkyBlue = NewColor(0xFF00BFFF)
	ColorDarkOrange  = NewColor(0xFFFF8C00)
	ColorPlum        = NewColor(0xFFDDA0DD)
	ColorCyan        = NewColor(0xFF00FFFF)
	ColorLime        = NewColor(0xFF00FF00)
)
