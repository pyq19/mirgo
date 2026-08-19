//go:build !windows

package engine

import "image"

// gdiFont 非 Windows 平台占位：GDI 后端不可用，TextRenderer 走 x/image 回退。
type gdiFont struct {
	ascent, heightPx int
}

func newGdiFont(family string, px int) (*gdiFont, error) { return nil, nil }

func (g *gdiFont) close() {}

func (g *gdiFont) glyph(ch rune) (img *image.RGBA, w, h, advance, topX, topY int, ok bool) {
	return nil, 0, 0, 0, 0, 0, false
}
