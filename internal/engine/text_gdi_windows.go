//go:build windows

package engine

import (
	"image"
	"image/color"
	"syscall"
	"unsafe"
)

// GDI 字形光栅化后端：x/image sfnt 不支持 TrueType hinting 与宋体内置
// 点阵（vendor sfnt.go:630），小字号矢量抗锯齿发糊；Delphi 原版走 GDI
// TextOut（FState.pas:3692-3706），此处用 Win32 GDI 复刻同等清晰度。

var gdi32dll = syscall.NewLazyDLL("gdi32.dll")

var (
	pCreateCompatibleDC = gdi32dll.NewProc("CreateCompatibleDC")
	pSelectObject       = gdi32dll.NewProc("SelectObject")
	pCreateFontIndirect = gdi32dll.NewProc("CreateFontIndirectW")
	pCreateDIBSection   = gdi32dll.NewProc("CreateDIBSection")
	pTextOut            = gdi32dll.NewProc("TextOutW")
	pGetTextExtent      = gdi32dll.NewProc("GetTextExtentPoint32W")
	pGetTextMetrics     = gdi32dll.NewProc("GetTextMetricsW")
	pSetTextColor       = gdi32dll.NewProc("SetTextColor")
	pSetBkColor         = gdi32dll.NewProc("SetBkColor")
	pSetBkMode          = gdi32dll.NewProc("SetBkMode")
	pDeleteDC           = gdi32dll.NewProc("DeleteDC")
	pDeleteObject       = gdi32dll.NewProc("DeleteObject")
)

type gdiLogfontW struct {
	LfHeight, LfWidth, LfEscapement, LfOrientation, LfWeight     int32
	LfItalic, LfUnderline, LfStrikeOut, LfCharSet                byte
	LfOutPrecision, LfClipPrecision, LfQuality, LfPitchAndFamily byte
	LfFaceName                                                   [32]uint16
}

type gdiTextMetricW struct {
	TmHeight, TmAscent, TmDescent, TmInternalLeading, TmExternalLeading int32
	TmAveCharWidth, TmMaxCharWidth, TmWeight, TmOverhang                int32
	TmDigitizedAspectX, TmDigitizedAspectY                              int32
	TmFirstChar, TmLastChar, TmDefaultChar, TmBreakChar                 uint16
	TmItalic, TmUnderlined, TmStruckOut                                 byte
	TmPitchAndFamily, TmCharSet                                         byte
}

type gdiBitmapInfoHeader struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

type gdiSize struct{ CX, CY int32 }

// gdiFont 持有内存 DC + HFONT + 32bpp DIB，逐字形 TextOut 后读回灰度覆盖。
type gdiFont struct {
	dc, hfont, dib   uintptr
	oldFont, oldDib  uintptr
	bits             []byte
	w, h             int
	ascent, heightPx int
	maxCharW         int
}

func newGdiFont(family string, px int) (*gdiFont, error) {
	if family == "" || px <= 0 {
		return nil, nil
	}
	dc, _, _ := pCreateCompatibleDC.Call(0)
	if dc == 0 {
		return nil, nil
	}
	name, _ := syscall.UTF16FromString(family)
	var lf gdiLogfontW
	lf.LfHeight = -int32(px)
	lf.LfCharSet = 1 // DEFAULT_CHARSET
	lf.LfQuality = 4 // ANTIALIASED_QUALITY：灰度抗锯齿（非 ClearType RGB）
	copy(lf.LfFaceName[:], name)
	hf, _, _ := pCreateFontIndirect.Call(uintptr(unsafe.Pointer(&lf)))
	if hf == 0 {
		pDeleteDC.Call(dc)
		return nil, nil
	}
	oldFont, _, _ := pSelectObject.Call(dc, hf)

	var tm gdiTextMetricW
	pGetTextMetrics.Call(dc, uintptr(unsafe.Pointer(&tm)))

	g := &gdiFont{
		dc: dc, hfont: hf, oldFont: oldFont,
		w:        int(tm.TmMaxCharWidth) + 4,
		h:        int(tm.TmHeight) + 4,
		ascent:   int(tm.TmAscent),
		heightPx: int(tm.TmHeight),
		maxCharW: int(tm.TmMaxCharWidth),
	}
	var bmi gdiBitmapInfoHeader
	bmi.BiSize = 40
	bmi.BiWidth = int32(g.w)
	bmi.BiHeight = -int32(g.h) // top-down
	bmi.BiPlanes = 1
	bmi.BiBitCount = 32
	var bitsPtr unsafe.Pointer
	dib, _, _ := pCreateDIBSection.Call(dc, uintptr(unsafe.Pointer(&bmi)), 0,
		uintptr(unsafe.Pointer(&bitsPtr)), 0, 0)
	if dib == 0 || bitsPtr == nil {
		g.close()
		return nil, nil
	}
	g.dib = dib
	g.oldDib, _, _ = pSelectObject.Call(dc, dib)
	g.bits = unsafe.Slice((*byte)(bitsPtr), g.w*g.h*4)
	pSetTextColor.Call(dc, 0xFFFFFF)
	pSetBkColor.Call(dc, 0)
	pSetBkMode.Call(dc, 1) // TRANSPARENT
	return g, nil
}

func (g *gdiFont) close() {
	if g.dc != 0 {
		if g.oldDib != 0 {
			pSelectObject.Call(g.dc, g.oldDib)
		}
		if g.oldFont != 0 {
			pSelectObject.Call(g.dc, g.oldFont)
		}
		pDeleteDC.Call(g.dc)
	}
	if g.dib != 0 {
		pDeleteObject.Call(g.dib)
	}
	if g.hfont != 0 {
		pDeleteObject.Call(g.hfont)
	}
	g.dc, g.dib, g.hfont = 0, 0, 0
}

// glyph 光栅化单个 rune：返回白字+alpha 覆盖图、裁剪框 w/h、
// advance（GetTextExtentPoint32W）、裁剪框左上角 (topX, topY)。
// 空白字符返回 w=0 且 ok=true。
func (g *gdiFont) glyph(ch rune) (img *image.RGBA, w, h, advance, topX, topY int, ok bool) {
	if ch > 0xFFFF || g.dc == 0 {
		return nil, 0, 0, 0, 0, 0, false
	}
	for i := range g.bits {
		g.bits[i] = 0
	}
	s := []uint16{uint16(ch)}
	var sz gdiSize
	pGetTextExtent.Call(g.dc, uintptr(unsafe.Pointer(&s[0])), 1, uintptr(unsafe.Pointer(&sz)))
	advance = int(sz.CX)
	r, _, _ := pTextOut.Call(g.dc, 0, 0, uintptr(unsafe.Pointer(&s[0])), 1)
	if r == 0 {
		return nil, 0, 0, advance, 0, 0, true
	}
	minX, minY, maxX, maxY := g.w, g.h, -1, -1
	for y := 0; y < g.h; y++ {
		row := g.bits[y*g.w*4 : (y+1)*g.w*4]
		for x := 0; x < g.w; x++ {
			if row[x*4+2] != 0 { // R 通道即灰度覆盖
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}
	if maxX < 0 {
		return nil, 0, 0, advance, 0, 0, true
	}
	w, h = maxX-minX+1, maxY-minY+1
	img = image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		src := g.bits[(minY+y)*g.w*4:]
		for x := 0; x < w; x++ {
			v := src[(minX+x)*4+2]
			img.SetRGBA(x, y, color.RGBA{R: 255, G: 255, B: 255, A: v})
		}
	}
	return img, w, h, advance, minX, minY, true
}
