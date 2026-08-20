package engine

import (
	"fmt"
	"image"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-gl/gl/v3.3-core/gl"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// TextScale 是全局文字光栅化缩放倍率。切换分辨率时调用 SetTextScale，
// 所有已创建的 TextRenderer 会立即重新缩放，之后新建的也自动应用。
var TextScale float64 = 1

var allRenderers []*TextRenderer

// SetTextScale 更新全局缩放倍率并立即应用到所有已注册的 TextRenderer。
func SetTextScale(s float64) {
	if s <= 0 {
		s = 1
	}
	TextScale = s
	for _, tr := range allRenderers {
		tr.SetScale(s)
	}
}

// glyphEntry 是缓存的字形纹理。逻辑度量（w/h/advance/bearingX/bearingY）
// 用于布局（MeasureText 等）；物理度量（pw/ph/padv/pbx/pby）是光栅化纹理
// 的真实像素度量，DrawText 按物理像素走笔以保持像素对齐。scale==1 时两者相同。
type glyphEntry struct {
	tex      uint32 // GL 纹理 ID
	w, h     int    // 字形像素尺寸（逻辑）
	advance  int    // 水平步进（逻辑像素）
	bearingX int    // 左侧 bearing（逻辑）
	bearingY int    // 顶部 bearing（从基线到字形顶部，逻辑）
	pw, ph   int    // 纹理物理像素尺寸
	padv     int    // 物理水平步进
	pbx      int    // 物理左侧 bearing
	pby      int    // 物理顶部 bearing（从基线到字形顶部）
}

// TextRenderer 使用 TTF 字体渲染文本，并缓存字形纹理。
// Windows 上优先 GDI 光栅化（gdi 非 nil），复刻 Delphi 原版 TextOut 的
// hinting/点阵清晰度；其余平台回退 x/image 矢量抗锯齿。
type TextRenderer struct {
	gl         *GLState
	face       font.Face // 逻辑尺寸 face（度量/布局）
	renderFace font.Face // 高分辨率 face（光栅化），nil 时 == face
	ascent     int       // 从基线到行顶的逻辑像素数
	ascentPhys int       // 从基线到行顶的物理像素数（与实际光栅 face 一致，scale==1 时 == ascent）
	scale      float64   // 光栅化倍率（1.0 = 不缩放）
	cache      map[rune]*glyphEntry
	cacheMu    sync.RWMutex
	size       float64
	fontData   []byte // 保存以便 WithSize 复用
	gdi        *gdiFont
	gdiRender  *gdiFont // 高分辨率 GDI 光栅，nil 时 == gdi
	gdiFamily  string
	gdiPx      int
}

// ptToPx 点数→96dpi 像素（9pt→12px，对应 Delphi Font.Height=-12）。
func ptToPx(size float64) int { return int(math.Round(size * 96 / 72)) }

// gdiFamilyForPath 字体文件→GDI 族名；不识别（含非 Windows 路径）返回 ""。
func gdiFamilyForPath(path string) string {
	switch strings.ToLower(filepath.Base(path)) {
	case "simsun.ttc":
		return "SimSun"
	case "msyh.ttc":
		return "Microsoft YaHei"
	case "msyhbd.ttc":
		return "Microsoft YaHei Bold"
	case "simhei.ttf":
		return "SimHei"
	case "arial.ttf":
		return "Arial"
	}
	return ""
}

// fontSearchPaths 列出常见系统的中文字体路径。
var fontSearchPaths = []string{
	// Windows：Delphi 原版默认宋体（MShare.pas:216-228 g_sCurFontName），
	// 9pt 行高恰好容纳聊天板 12px 行距（FState:3703）。
	`C:\Windows\Fonts\simsun.ttc`, // 宋体
	`C:\Windows\Fonts\msyh.ttc`,   // 微软雅黑
	`C:\Windows\Fonts\msyhbd.ttc`, // 微软雅黑 Bold
	`C:\Windows\Fonts\simhei.ttf`, // 黑体
	`C:\Windows\Fonts\arial.ttf`,  // Arial（英文兜底）
	// Linux
	"/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc",
	"/usr/share/fonts/noto-cjk/NotoSansCJK-Regular.ttc",
	"/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc",
	"/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",
	"/usr/share/fonts/truetype/wqy/wqy-zenhei.ttc",
	"/usr/share/fonts/truetype/droid/DroidSansFallbackFull.ttf",
	"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
	// macOS
	"/System/Library/Fonts/PingFang.ttc",
	"/System/Library/Fonts/STHeiti Light.ttc",
}

// NewTextRenderer 创建一个 TextRenderer。若 fontPath 为空，则尝试常见的 Windows 字体。
func NewTextRenderer(glState *GLState, fontPath string, size float64) (*TextRenderer, error) {
	if size <= 0 {
		size = 9
	}

	resolvedPath := fontPath
	if resolvedPath == "" {
		for _, p := range fontSearchPaths {
			if _, err := os.Stat(p); err == nil {
				resolvedPath = p
				break
			}
		}
	}
	if resolvedPath == "" {
		return nil, fmt.Errorf("no font found, specify a TTF/TTC path")
	}

	fontData, err := os.ReadFile(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("read font %s: %w", resolvedPath, err)
	}

	// 先尝试单字体解析（TTF），再尝试字体集合解析（TTC）。
	f, parseErr := opentype.Parse(fontData)
	if parseErr != nil {
		col, colErr := opentype.ParseCollection(fontData)
		if colErr != nil {
			return nil, fmt.Errorf("parse font: single=%v, collection=%v", parseErr, colErr)
		}
		if col.NumFonts() == 0 {
			return nil, fmt.Errorf("font collection is empty")
		}
		f, err = col.Font(0)
		if err != nil {
			return nil, fmt.Errorf("get font from collection: %w", err)
		}
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     96,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("create face: %w", err)
	}

	metrics := face.Metrics()
	ascent := metrics.Ascent.Ceil()

	tr := &TextRenderer{
		gl:        glState,
		face:      face,
		ascent:    ascent,
		scale:     1,
		cache:     make(map[rune]*glyphEntry),
		size:      size,
		fontData:  fontData,
		gdiFamily: gdiFamilyForPath(resolvedPath),
		gdiPx:     ptToPx(size),
	}
	tr.gdi, _ = newGdiFont(tr.gdiFamily, tr.gdiPx)
	if tr.gdi != nil {
		tr.ascent = tr.gdi.ascent
	}
	tr.ascentPhys = tr.ascent
	allRenderers = append(allRenderers, tr)
	if TextScale != 1 {
		tr.SetScale(TextScale)
	}
	return tr, nil
}

// WithSize 用相同字体、不同字号创建一个新的 TextRenderer。
func (tr *TextRenderer) WithSize(size float64) (*TextRenderer, error) {
	f, parseErr := opentype.Parse(tr.fontData)
	if parseErr != nil {
		col, colErr := opentype.ParseCollection(tr.fontData)
		if colErr != nil {
			return nil, fmt.Errorf("parse font: single=%v, collection=%v", parseErr, colErr)
		}
		f, colErr = col.Font(0)
		if colErr != nil {
			return nil, fmt.Errorf("get font from collection: %w", colErr)
		}
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     96,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("create face: %w", err)
	}

	child := &TextRenderer{
		gl:        tr.gl,
		face:      face,
		ascent:    face.Metrics().Ascent.Ceil(),
		scale:     1,
		cache:     make(map[rune]*glyphEntry),
		size:      size,
		fontData:  tr.fontData,
		gdiFamily: tr.gdiFamily,
		gdiPx:     ptToPx(size),
	}
	child.gdi, _ = newGdiFont(child.gdiFamily, child.gdiPx)
	if child.gdi != nil {
		child.ascent = child.gdi.ascent
	}
	child.ascentPhys = child.ascent
	allRenderers = append(allRenderers, child)
	if TextScale != 1 {
		child.SetScale(TextScale)
	}
	return child, nil
}

// SetScale 设置光栅化倍率。scale > 1 时以更高 DPI 光栅化字形使文字清晰，
// 但保持逻辑度量不变（布局/定位不受影响）。切换分辨率后调用。
func (tr *TextRenderer) SetScale(s float64) {
	if s <= 0 {
		s = 1
	}
	if s == tr.scale {
		return
	}
	tr.scale = s
	if tr.gdi != nil {
		if tr.gdiRender != nil {
			tr.gdiRender.close()
			tr.gdiRender = nil
		}
		if s != 1 {
			tr.gdiRender, _ = newGdiFont(tr.gdiFamily, int(math.Round(float64(tr.gdiPx)*s)))
		}
		if tr.gdiRender != nil {
			tr.ascentPhys = tr.gdiRender.ascent
		} else {
			// 高分辨率字体创建失败（或 s==1）：回退逻辑度量。
			tr.ascentPhys = tr.ascent
		}
		tr.clearCache()
		return
	}
	if s == 1 {
		tr.renderFace = nil
	} else {
		f, parseErr := opentype.Parse(tr.fontData)
		if parseErr != nil {
			col, _ := opentype.ParseCollection(tr.fontData)
			f, _ = col.Font(0)
		}
		if f != nil {
			face, err := opentype.NewFace(f, &opentype.FaceOptions{
				Size:    tr.size,
				DPI:     96 * s,
				Hinting: font.HintingFull,
			})
			if err == nil {
				tr.renderFace = face
			}
		}
	}
	if tr.renderFace != nil {
		tr.ascentPhys = tr.renderFace.Metrics().Ascent.Ceil()
	} else {
		tr.ascentPhys = tr.ascent
	}
	tr.clearCache()
}

func (tr *TextRenderer) clearCache() {
	tr.cacheMu.Lock()
	for _, g := range tr.cache {
		if g.tex != 0 {
			gl.DeleteTextures(1, &g.tex)
		}
	}
	tr.cache = make(map[rune]*glyphEntry)
	tr.cacheMu.Unlock()
}

// getGlyph 返回某个 rune 缓存的字形，缓存未命中时进行光栅化。
func (tr *TextRenderer) getGlyph(ch rune) *glyphEntry {
	tr.cacheMu.RLock()
	if g, ok := tr.cache[ch]; ok {
		tr.cacheMu.RUnlock()
		return g
	}
	tr.cacheMu.RUnlock()

	// GDI 后端（Windows）：逻辑度量取逻辑字体，纹理与物理度量取高分辨率字体
	// （gdiFont.glyph 的度量即该字体自身物理像素：GetTextExtentPoint32W + 墨迹框）。
	if tr.gdi != nil {
		img, gw, gh, advance, topX, topY, ok := tr.gdi.glyph(ch)
		if !ok {
			spaceAdv, _ := tr.face.GlyphAdvance(' ')
			g := &glyphEntry{advance: spaceAdv.Ceil(), padv: spaceAdv.Ceil()}
			if tr.gdiRender != nil {
				if _, _, _, radv, _, _, rok := tr.gdiRender.glyph(' '); rok {
					g.padv = radv
				}
			}
			return g
		}
		src := img
		// 物理度量默认为逻辑度量复制（scale==1 或高分辨率光栅失败时），
		// 高分辨率成功后改用 gdiRender 的物理度量。
		pw, ph, padv, pbx, pby := gw, gh, advance, topX, tr.ascent-topY
		if tr.gdiRender != nil {
			rimg, rw, rh, radv, rtopX, rtopY, rok := tr.gdiRender.glyph(ch)
			if rok {
				padv = radv
				if rimg != nil {
					src = rimg
					pw, ph, pbx = rw, rh, rtopX
					pby = tr.gdiRender.ascent - rtopY
				}
			}
		}
		var tex uint32
		if src != nil {
			tex = tr.gl.UploadTexture(src)
		}
		g := &glyphEntry{
			tex:      tex,
			w:        gw,
			h:        gh,
			advance:  advance,
			bearingX: topX,
			// DrawText 以 y+ascent-bearingY 定位字形顶：GDI 裁剪框顶
			// 相对行顶为 topY，等价 bearingY = ascent-topY。
			bearingY: tr.ascent - topY,
			pw:       pw,
			ph:       ph,
			padv:     padv,
			pbx:      pbx,
			pby:      pby,
		}
		tr.cacheMu.Lock()
		tr.cache[ch] = g
		tr.cacheMu.Unlock()
		return g
	}

	// 光栅 face：scale > 1 时为高分辨率 renderFace，否则为逻辑 face。
	// 其 fixed26.6 度量即物理像素值。
	rf := tr.renderFace
	if rf == nil {
		rf = tr.face
	}
	// physAdv 取实际光栅 face 的物理 advance（失败时空格兜底）。
	physAdv := func(r rune) int {
		if a, ok := rf.GlyphAdvance(r); ok {
			return a.Ceil()
		}
		if a, ok := rf.GlyphAdvance(' '); ok {
			return a.Ceil()
		}
		return 0
	}

	// 逻辑度量（用于布局/quad 尺寸）。
	advance, ok := tr.face.GlyphAdvance(ch)
	if !ok {
		spaceAdv, _ := tr.face.GlyphAdvance(' ')
		return &glyphEntry{advance: spaceAdv.Ceil(), padv: physAdv(' ')}
	}

	bounds, _, ok := tr.face.GlyphBounds(ch)
	if !ok {
		spaceAdv, _ := tr.face.GlyphAdvance(' ')
		return &glyphEntry{advance: spaceAdv.Ceil(), padv: physAdv(' ')}
	}

	logMinX := bounds.Min.X.Floor()
	logMaxX := bounds.Max.X.Ceil()
	logMinY := bounds.Min.Y.Floor()
	logMaxY := bounds.Max.Y.Ceil()

	gw := logMaxX - logMinX
	gh := logMaxY - logMinY

	if gw <= 0 || gh <= 0 {
		return &glyphEntry{advance: advance.Ceil(), padv: physAdv(ch)}
	}

	// 光栅化：scale > 1 时用 renderFace 生成高分辨率位图。
	// 物理度量默认为逻辑度量复制（rf==tr.face 即 scale==1 时两者本就相同）。
	pw, ph, pbx, pby, padv := gw, gh, logMinX, -logMinY, advance.Ceil()
	var img *image.RGBA
	if rf == tr.face {
		img = image.NewRGBA(image.Rect(0, 0, gw, gh))
		d := &font.Drawer{Dst: img, Src: image.White, Face: rf, Dot: fixed.P(-logMinX, -logMinY)}
		d.DrawString(string(ch))
	} else {
		// 高分辨率光栅化；物理度量取 renderFace。
		rBounds, _, _ := rf.GlyphBounds(ch)
		rMinX := rBounds.Min.X.Floor()
		rMaxX := rBounds.Max.X.Ceil()
		rMinY := rBounds.Min.Y.Floor()
		rMaxY := rBounds.Max.Y.Ceil()
		rw := rMaxX - rMinX
		rh := rMaxY - rMinY
		if rw <= 0 || rh <= 0 {
			return &glyphEntry{advance: advance.Ceil(), padv: physAdv(ch)}
		}
		img = image.NewRGBA(image.Rect(0, 0, rw, rh))
		d := &font.Drawer{Dst: img, Src: image.White, Face: rf, Dot: fixed.P(-rMinX, -rMinY)}
		d.DrawString(string(ch))
		pw, ph, pbx, pby, padv = rw, rh, rMinX, -rMinY, physAdv(ch)
	}

	tex := tr.gl.UploadTexture(img)

	g := &glyphEntry{
		tex:      tex,
		w:        gw,
		h:        gh,
		advance:  advance.Ceil(),
		bearingX: logMinX,
		bearingY: -logMinY,
		pw:       pw,
		ph:       ph,
		padv:     padv,
		pbx:      pbx,
		pby:      pby,
	}

	tr.cacheMu.Lock()
	tr.cache[ch] = g
	tr.cacheMu.Unlock()

	return g
}

// DrawText 用给定的 RGBA 颜色在 (x, y) 处渲染文本。
// (x, y) 是文本基线区域的左上角（逻辑坐标）。
// 走笔按物理像素进行：起点四舍五入到物理像素网格，步进使用字形的物理
// advance，再换算回逻辑坐标提交 quad。这样在 NEAREST 采样下纹理 texel
// 与物理像素 1:1 对齐，任意缩放倍率（含非整数倍）文字都保持锐利。
// 布局度量（MeasureText 等）仍为逻辑值，不受影响。
func (tr *TextRenderer) DrawText(text string, x, y float32, r, g, b, a float32, proj [16]float32) {
	s := tr.scale
	penX := int(math.Round(float64(x) * s))
	penY := int(math.Round(float64(y) * s))
	inv := float32(1 / s)
	for _, ch := range text {
		glyph := tr.getGlyph(ch)
		if glyph.tex != 0 {
			// 位置：penX + pbx，penY + ascentPhys - pby（物理像素），
			// 除以 s 换算回逻辑坐标。
			dx := float32(penX+glyph.pbx) * inv
			dy := float32(penY+(tr.ascentPhys-glyph.pby)) * inv
			tr.gl.DrawQuadTint(glyph.tex, dx, dy, float32(glyph.pw)*inv, float32(glyph.ph)*inv, r, g, b, a, proj)
		}
		// 空格等无纹理字形同样按物理 advance 步进。
		penX += glyph.padv
	}
}

// DrawTextOutline 在主色后面渲染 1px 描边（颜色为 or, og, ob, oa）的文本。
// 对应 Delphi 用于 NPC 对话和提示的 BoldTextOut。
func (tr *TextRenderer) DrawTextOutline(text string, x, y float32, r, g, b, a float32, or, og, ob, oa float32, proj [16]float32) {
	tr.DrawText(text, x-1, y, or, og, ob, oa, proj)
	tr.DrawText(text, x+1, y, or, og, ob, oa, proj)
	tr.DrawText(text, x, y-1, or, og, ob, oa, proj)
	tr.DrawText(text, x, y+1, or, og, ob, oa, proj)
	tr.DrawText(text, x, y, r, g, b, a, proj)
}

// DrawTextBold 通过水平偏移 1px 绘制两次来渲染伪粗体文本。
// 近似 Delphi 的 fsBold 样式。
func (tr *TextRenderer) DrawTextBold(text string, x, y float32, r, g, b, a float32, proj [16]float32) {
	tr.DrawText(text, x, y, r, g, b, a, proj)
	tr.DrawText(text, x+1, y, r, g, b, a, proj)
}

// DrawTextBoldOutline 渲染带 1px 描边的伪粗体文本。
// 对应 Delphi 中带 fsBold 样式的 BoldTextOut（FState.pas:2274-2279）。
func (tr *TextRenderer) DrawTextBoldOutline(text string, x, y float32, r, g, b, a float32, or, og, ob, oa float32, proj [16]float32) {
	tr.DrawText(text, x-1, y, or, og, ob, oa, proj)
	tr.DrawText(text, x+1, y, or, og, ob, oa, proj)
	tr.DrawText(text, x+2, y, or, og, ob, oa, proj)
	tr.DrawText(text, x, y-1, or, og, ob, oa, proj)
	tr.DrawText(text, x+1, y-1, or, og, ob, oa, proj)
	tr.DrawText(text, x, y+1, or, og, ob, oa, proj)
	tr.DrawText(text, x+1, y+1, or, og, ob, oa, proj)
	tr.DrawText(text, x, y, r, g, b, a, proj)
	tr.DrawText(text, x+1, y, r, g, b, a, proj)
}

// MeasureText 返回文本的像素宽度。
func (tr *TextRenderer) MeasureText(text string) int {
	width := 0
	for _, ch := range text {
		glyph := tr.getGlyph(ch)
		width += glyph.advance
	}
	return width
}

// MeasureChar 返回单个 rune 的像素宽度。
func (tr *TextRenderer) MeasureChar(ch rune) int {
	return tr.getGlyph(ch).advance
}

// Ascent 返回 ascent（像素）。
func (tr *TextRenderer) Ascent() int {
	return tr.ascent
}

// LineHeight 返回完整行高（ascent + descent）。
func (tr *TextRenderer) LineHeight() int {
	if tr.gdi != nil {
		return tr.gdi.heightPx
	}
	metrics := tr.face.Metrics()
	return (metrics.Ascent + metrics.Descent).Ceil()
}

// Destroy 释放所有缓存的 GL 纹理并从全局注册表移除。
func (tr *TextRenderer) Destroy() {
	if tr.gdi != nil {
		tr.gdi.close()
		tr.gdi = nil
	}
	if tr.gdiRender != nil {
		tr.gdiRender.close()
		tr.gdiRender = nil
	}
	tr.cacheMu.Lock()
	for _, g := range tr.cache {
		if g.tex != 0 {
			gl.DeleteTextures(1, &g.tex)
		}
	}
	tr.cache = make(map[rune]*glyphEntry)
	tr.cacheMu.Unlock()
	for i, r := range allRenderers {
		if r == tr {
			allRenderers = append(allRenderers[:i], allRenderers[i+1:]...)
			break
		}
	}
}
