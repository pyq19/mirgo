package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pyq0109/mirgo/internal/engine"
	"github.com/pyq0109/mirgo/internal/log"
	"github.com/pyq0109/mirgo/internal/wil"
)

// UIManager 在控件树之间路由输入和绘制. 移植自 TDWinManager
// (DWinCtl.pas:878-1122): 模态窗口短路一切, 其次是鼠标捕获,
// 再其次是以 Root (全屏 DBackground 控件) 为根的树.
// 绘制顺序 = 树序, 模态窗口最后 (置顶).
type UIManager struct {
	gl        *engine.GLState
	resources *engine.ResourceManager
	text      *engine.TextRenderer

	Root    *UIControl
	Modal   *UIControl
	Capture *UIControl
	Focused *UIControl

	ShowBounds bool // 调试: 绘制所有可见控件的包围盒

	ShowImgs bool // 调试: 绘制可见控件绑定的 WIL 实际图像 (独立于常态渲染路径)

	ShowInspect bool // 调试: 悬停拾取 — 高亮光标下控件并显示一行属性
	// 调试: 包围盒名字模式 0=仅框 1=框+仅光标悬停控件及其祖先的名字 2=框+全部名字(带底条+去重)。
	// 与 ShowBounds 解耦: ShowBounds 只决定是否画框, BoundsNames 决定名字密度。
	BoundsNames    int
	hoverX, hoverY int // 调试: 最近一帧光标位置 (供 BoundsNames=1 取悬停控件)
}

// gActiveUI 是当前活动场景的 UIManager 引用。
// 各场景在 Open() 中设置、Close() 中置 nil，
// 使全局注册的 ui bounds 调试命令可跨场景工作。
var gActiveUI *UIManager

func NewUIManager(gl *engine.GLState, resources *engine.ResourceManager, text *engine.TextRenderer) *UIManager {
	root := NewUIControl("DBackground", KindControl)
	root.Background = true
	root.Visible = true
	root.Left, root.Top = 0, 0
	root.Width, root.Height = ScreenWidth, ScreenHeight
	return &UIManager{
		gl:        gl,
		resources: resources,
		text:      text,
		Root:      root,
	}
}

func (m *UIManager) SetText(t *engine.TextRenderer) { m.text = t }

// ---------------------------------------------------------------------------
// 焦点 / 捕获 (DWinCtl.pas:211-231)
// ---------------------------------------------------------------------------

func (m *UIManager) SetFocus(c *UIControl)   { m.Focused = c }
func (m *UIManager) ReleaseFocus()           { m.Focused = nil }
func (m *UIManager) SetCapture(c *UIControl) { m.Capture = c }
func (m *UIManager) ReleaseCapture()         { m.Capture = nil }

// canFocusMsg: 只有当没有捕获生效, 或本控件 (或其父控件) 持有捕获时,
// 控件才能接收鼠标消息 (DWinCtl.pas:456-462).
func (m *UIManager) canFocusMsg(c *UIControl) bool {
	return m.Capture == nil || m.Capture == c || m.Capture == c.Parent
}

// ShowWindow 显示窗口, 置顶浮动窗口并获取焦点
// (TDWindow.Show, DWinCtl.pas:859-867).
func (m *UIManager) ShowWindow(c *UIControl) {
	c.Show()
	if c.EnableFocus {
		m.SetFocus(c)
	}
}

// ShowModal 将窗口作为模态显示: 所有输入都路由给它, 且绘制在最上层
// (TDWindow.ShowModal, DWinCtl.pas:869-875; 管理器路由 :1008-1015).
// 非阻塞 — 结果通过回调返回.
func (m *UIManager) ShowModal(c *UIControl) {
	c.Visible = true
	m.Modal = c
	if c.EnableFocus {
		m.SetFocus(c)
	}
}

// CloseModal 关闭模态窗口。焦点/捕获若落在被关闭的子树内必须一并释放:
// Delphi 在模态结束时 RestoreHideControls 并按需重聚焦 EdChat
// (FState.pas:712-720, :2119-2123); Go 若不释放, 焦点会永久留在已隐藏
// 的输入控件上, 此后所有键盘输入被吞 (且 Enter 会重复触发旧对话框回调)。
// 必须在此处集中做, 不能依赖调用方: 链式对话框 (ShowConfirm 回调里再
// ShowInput) 依赖"先放焦 → 新模态再抢焦"的顺序。
func (m *UIManager) CloseModal(c *UIControl) {
	if m.Modal == c {
		m.Modal = nil
	}
	if m.Focused != nil && withinSubtree(m.Focused, c) {
		m.ReleaseFocus()
	}
	if m.Capture != nil && withinSubtree(m.Capture, c) {
		m.ReleaseCapture()
	}
	c.Hide()
}

// ---------------------------------------------------------------------------
// 输入路由
// ---------------------------------------------------------------------------

// RouteMouseDown 在 UI 消费了事件时返回 true (游戏世界必须忽略该事件
// — 对应 ClMain.pas 的 `if g_DWinMan.MouseDown then exit`).
func (m *UIManager) RouteMouseDown(absX, absY int, button int) bool {
	if m.Modal != nil && m.Modal.Visible {
		// 模态可见时恒吞掉鼠标事件 (TDWinManager.MouseDown
		// 恒返回 TRUE, DWinCtl.pas:1008-1015), 点在模态框外也不许落到世界.
		log.Logf(log.LevelDebug, "UI", "mouse-down routed -> modal %s pos=(%d,%d)", m.Modal.Name, absX, absY)
		m.dispatchMouseDown(m.Modal, button, m.Modal.ParentSpaceX(absX), m.Modal.ParentSpaceY(absY))
		return true
	}
	if m.Capture != nil {
		log.Logf(log.LevelDebug, "UI", "mouse-down routed -> capture %s pos=(%d,%d)", m.Capture.Name, absX, absY)
		return m.dispatchMouseDown(m.Capture, button, m.Capture.ParentSpaceX(absX), m.Capture.ParentSpaceY(absY))
	}
	return m.dispatchMouseDown(m.Root, button, absX-m.Root.Left, absY-m.Root.Top)
}

func (m *UIManager) RouteMouseUp(absX, absY int, button int) bool {
	if m.Modal != nil && m.Modal.Visible {
		log.Logf(log.LevelDebug, "UI", "mouse-up routed -> modal %s pos=(%d,%d)", m.Modal.Name, absX, absY)
		m.dispatchMouseUp(m.Modal, button, m.Modal.ParentSpaceX(absX), m.Modal.ParentSpaceY(absY))
		return true
	}
	if m.Capture != nil {
		log.Logf(log.LevelDebug, "UI", "mouse-up routed -> capture %s pos=(%d,%d)", m.Capture.Name, absX, absY)
		return m.dispatchMouseUp(m.Capture, button, m.Capture.ParentSpaceX(absX), m.Capture.ParentSpaceY(absY))
	}
	return m.dispatchMouseUp(m.Root, button, absX-m.Root.Left, absY-m.Root.Top)
}

func (m *UIManager) RouteMouseMove(absX, absY int) bool {
	if m.Modal != nil && m.Modal.Visible {
		// 恒吞 (DWinCtl.pas:976-1001).
		m.dispatchMouseMove(m.Modal, m.Modal.ParentSpaceX(absX), m.Modal.ParentSpaceY(absY))
		return true
	}
	if m.Capture != nil {
		return m.dispatchMouseMove(m.Capture, m.Capture.ParentSpaceX(absX), m.Capture.ParentSpaceY(absY))
	}
	return m.dispatchMouseMove(m.Root, absX-m.Root.Left, absY-m.Root.Top)
}

// RouteDblClick 路由合成双击 (捕获优先语义,
// TDControl.DblClick DWinCtl.pas:553-578).
func (m *UIManager) RouteDblClick(absX, absY int) bool {
	if m.Modal != nil && m.Modal.Visible {
		m.dispatchDblClick(m.Modal, m.Modal.ParentSpaceX(absX), m.Modal.ParentSpaceY(absY))
		return true
	}
	if m.Capture != nil {
		return m.dispatchDblClick(m.Capture, m.Capture.ParentSpaceX(absX), m.Capture.ParentSpaceY(absY))
	}
	return m.dispatchDblClick(m.Root, absX-m.Root.Left, absY-m.Root.Top)
}

// RouteChar / RouteKeyDown: 优先模态子树, 否则发给焦点控件
// (TDWinManager.KeyPress/KeyDown, DWinCtl.pas:916-974).
//
// 与 Delphi 的两处有意差异 (根因: Delphi 是双通道 — KeyPreview 吞掉的
// 只是窗体预览, 字符仍由 VCL 送给 OS 焦点控件 EdChat; Go 是单通道,
// 这里吞掉即彻底丢弃):
//  1. 模态可见时恒吞 — 对话框打开期间绝不漏键给游戏 (不误放技能/吃药).
//  2. 焦点控件没有对应处理器时放行 — Delphi 的"无处理器也吞"只影响
//     窗体预览逻辑, 字符照样进 EdChat; Go 单通道下必须放宽, 否则
//     点过任何 HUD 后聊天完全打不了字.
func (m *UIManager) RouteChar(ch rune) bool {
	if m.Modal != nil {
		if m.Modal.Visible {
			m.dispatchChar(m.Modal, ch)
			return true
		}
		m.Modal = nil // 自清理 (DWinCtl.pas:924-929)
	}
	f := m.Focused
	if f == nil {
		return false
	}
	if !f.treeVisible() {
		m.ReleaseFocus() // DWinCtl.pas:933-936
		return false
	}
	if f.OnChar == nil {
		return false
	}
	return m.dispatchChar(f, ch)
}

func (m *UIManager) RouteKeyDown(key int) bool {
	if m.Modal != nil {
		if m.Modal.Visible {
			m.dispatchKeyDown(m.Modal, key)
			return true
		}
		m.Modal = nil
	}
	f := m.Focused
	if f == nil {
		return false
	}
	if !f.treeVisible() {
		m.ReleaseFocus()
		return false
	}
	if f.OnKeyDown == nil {
		return false
	}
	return m.dispatchKeyDown(f, key)
}

// ---------------------------------------------------------------------------
// 递归分发 (TDControl 各方法, DWinCtl.pas:464-605,665-852)
// ---------------------------------------------------------------------------

func (m *UIManager) dispatchMouseDown(c *UIControl, button, x, y int) bool {
	// 子控件从后往前 (最上层最后绘制 = 最先检测).
	// TDGrid 不递归 (其 MouseDown 覆写跳过了 inherited).
	if c.Kind != KindGrid {
		for i := len(c.Children) - 1; i >= 0; i-- {
			ch := c.Children[i]
			if ch.Visible && m.dispatchMouseDown(ch, button, x-c.Left, y-c.Top) {
				return true
			}
		}
	}

	switch c.Kind {
	case KindGrid:
		// TDGrid.MouseDown (:721-736): 左键且命中单元格 → 捕获.
		if button == 0 {
			if col, row, ok := c.ColRowAt(x, y); ok {
				log.Logf(log.LevelDebug, "UI", "grid mouse-down %s cell=(%d,%d)", c.DebugPath(), col, row)
				c.SelectCol, c.SelectRow = col, row
				m.SetCapture(c)
				return true
			}
		}
		return false

	default:
		if c.Background {
			// TDControl.MouseDown 的 Background 分支 (:504-511).
			if c.OnBackgroundClick != nil {
				c.WantReturn = false
				c.OnBackgroundClick(c)
				consumed := c.WantReturn
				m.ReleaseFocus()
				return consumed
			}
			m.ReleaseFocus()
			return false
		}
		if m.canFocusMsg(c) && (c.InRange(x, y) || m.Capture == c) {
			log.Logf(log.LevelDebug, "UI", "mouse-down %s(%s) pos=(%d,%d) btn=%d", c.DebugPath(), c.Kind, x, y, button)
			if c.OnMouseDown != nil {
				c.OnMouseDown(c, button, x, y)
			}
			if c.EnableFocus {
				m.SetFocus(c)
			}
			if c.Kind == KindButton || c.Kind == KindWindow {
				// TDButton.MouseDown (:665-675): 按下视觉 + 捕获.
				// TDWindow 继承自 TDButton, 所以窗口也会捕获 —
				// 这正是浮动拖动得以工作的原因.
				c.Downed = true
				m.SetCapture(c)
			}
			if c.Kind == KindWindow && c.Floating {
				// TDWindow.MouseDown (:841-852): 置顶 + 记录抓取点.
				c.RaiseToTop()
				c.SpotX, c.SpotY = x, y
			}
			return true
		}
		return false
	}
}

func (m *UIManager) dispatchMouseUp(c *UIControl, button, x, y int) bool {
	if c.Kind != KindGrid {
		for i := len(c.Children) - 1; i >= 0; i-- {
			ch := c.Children[i]
			if ch.Visible && m.dispatchMouseUp(ch, button, x-c.Left, y-c.Top) {
				return true
			}
		}
	}

	switch c.Kind {
	case KindGrid:
		// TDGrid.MouseUp (:752-769): 仅当抬起格 == 按下格时选中.
		consumed := false
		if button == 0 {
			if col, row, ok := c.ColRowAt(x, y); ok {
				if c.SelectCol == col && c.SelectRow == row {
					log.Logf(log.LevelDebug, "UI", "grid select %s cell=(%d,%d)", c.DebugPath(), col, row)
					c.Col, c.Row = col, row
					if c.OnGridSelect != nil {
						c.OnGridSelect(c, col, row)
					}
				}
				consumed = true
			}
			m.ReleaseCapture()
		}
		return consumed

	default:
		if m.Capture != nil {
			if m.Capture == c {
				if c.Kind == KindButton || c.Kind == KindWindow {
					// TDButton.MouseUp (:677-695): 仅当抬起仍在按钮内时
					// 才触发点击 (TDWindow 继承此行为 —
					// TDWindow.MouseUp 只是调用 inherited).
					m.ReleaseCapture()
					if !c.Background && c.InRange(x, y) {
						log.Logf(log.LevelDebug, "UI", "click %s(%s) pos=(%d,%d)", c.DebugPath(), c.Kind, x, y)
						if c.OnMouseUp != nil {
							c.OnMouseUp(c, button, x, y)
						}
						if c.ClickSound >= 0 {
							gSound.PlaySound(c.ClickSound)
						}
						if c.OnClick != nil {
							c.OnClick(c, x, y)
						}
					}
					c.Downed = false
					return true
				}
				if c.OnMouseUp != nil {
					c.OnMouseUp(c, button, x, y)
				}
				return true
			}
			return false // 捕获属于其他控件
		}
		if c.Background {
			return false
		}
		if c.InRange(x, y) {
			if c.OnMouseUp != nil {
				c.OnMouseUp(c, button, x, y)
			}
			// Delphi 中抬起之后会跟一个 WM_LBUTTONCLK; 带 OnClick 的
			// 控件 (如 DStateWin 魔法页命中检测) 在此收到它.
			if c.OnClick != nil {
				log.Logf(log.LevelDebug, "UI", "click %s(%s) pos=(%d,%d)", c.DebugPath(), c.Kind, x, y)
				if c.ClickSound >= 0 {
					gSound.PlaySound(c.ClickSound)
				}
				c.OnClick(c, x, y)
			}
			if c.Kind == KindButton || c.Kind == KindWindow {
				c.Downed = false
			}
			return true
		}
		return false
	}
}

func (m *UIManager) dispatchMouseMove(c *UIControl, x, y int) bool {
	if c.Kind != KindGrid {
		for i := len(c.Children) - 1; i >= 0; i-- {
			ch := c.Children[i]
			if ch.Visible && m.dispatchMouseMove(ch, x-c.Left, y-c.Top) {
				return true
			}
		}
	}

	switch c.Kind {
	case KindGrid:
		// TDGrid.MouseMove (:738-750): 悬停单元格通知.
		if col, row, ok := c.ColRowAt(x, y); ok {
			if c.OnGridMouseMove != nil {
				c.OnGridMouseMove(c, col, row)
			}
			return true
		}
		return false

	default:
		if m.Capture != nil {
			if m.Capture == c {
				if c.Kind == KindButton {
					// TDButton 按下视觉跟随光标 (:654-663).
					c.Downed = c.InRange(x, y)
				}
				if c.Kind == KindWindow && c.Floating {
					// TDWindow.MouseMove 带边界钳制的拖动 (:820-839).
					if c.SpotX != x || c.SpotY != y {
						al := c.Left + (x - c.SpotX)
						at := c.Top + (y - c.SpotY)
						if al+c.Width < WinLeft {
							al = WinLeft - c.Width
						}
						if al > m.Root.Width-60 {
							al = m.Root.Width - 60
						}
						if at+c.Height < WinTop {
							at = WinTop - c.Height
						}
						if at+c.Height > m.Root.Height-30 {
							at = m.Root.Height - 30 - c.Height
						}
						log.Logf(log.LevelDebug, "UI", "drag %s to=(%d,%d)", c.Name, al, at)
						c.Left, c.Top = al, at
						c.SpotX, c.SpotY = x, y
					}
				}
				if c.OnMouseMove != nil {
					c.OnMouseMove(c, x, y)
				}
				return true
			}
			return false
		}
		if c.Background {
			return false
		}
		if c.InRange(x, y) {
			if c.OnMouseMove != nil {
				c.OnMouseMove(c, x, y)
			}
			return true
		}
		return false
	}
}

func (m *UIManager) dispatchDblClick(c *UIControl, x, y int) bool {
	// 捕获优先于子控件 (DWinCtl.pas:558-565).
	if m.Capture != nil {
		if m.Capture == c && c.OnDblClick != nil {
			log.Logf(log.LevelDebug, "UI", "double-click %s(%s) pos=(%d,%d)", c.DebugPath(), c.Kind, x, y)
			c.OnDblClick(c, x, y)
			return true
		}
		return false
	}
	for i := len(c.Children) - 1; i >= 0; i-- {
		ch := c.Children[i]
		if ch.Visible && m.dispatchDblClick(ch, x-c.Left, y-c.Top) {
			return true
		}
	}
	if c.Background {
		return false
	}
	if c.InRange(x, y) && c.OnDblClick != nil {
		log.Logf(log.LevelDebug, "UI", "double-click %s(%s) pos=(%d,%d)", c.DebugPath(), c.Kind, x, y)
		c.OnDblClick(c, x, y)
		return true
	}
	return false
}

func (m *UIManager) dispatchChar(c *UIControl, ch rune) bool {
	if c.Background {
		return false
	}
	for i := len(c.Children) - 1; i >= 0; i-- {
		child := c.Children[i]
		if child.Visible && m.dispatchChar(child, ch) {
			return true
		}
	}
	if m.Focused == c {
		if c.OnChar != nil {
			c.OnChar(c, ch)
		}
		return true
	}
	return false
}

func (m *UIManager) dispatchKeyDown(c *UIControl, key int) bool {
	if c.Background {
		return false
	}
	for i := len(c.Children) - 1; i >= 0; i-- {
		child := c.Children[i]
		if child.Visible && m.dispatchKeyDown(child, key) {
			return true
		}
	}
	if m.Focused == c {
		if c.OnKeyDown != nil {
			c.OnKeyDown(c, key)
		}
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// 绘制 (TDControl.DirectPaint, DWinCtl.pas:623-639; TDGrid :783-796;
// TDWinManager :1108-1122)
// ---------------------------------------------------------------------------

func (m *UIManager) Paint(proj [16]float32) {
	m.paintControl(m.Root, proj)
	if m.Modal != nil && m.Modal.Visible {
		m.paintControl(m.Modal, proj)
	}
}

func (m *UIManager) paintControl(c *UIControl, proj [16]float32) {
	if !c.Visible {
		return
	}

	if c.Kind == KindGrid {
		// 网格由属主自绘; 它不绘制自身图像, 也没有参与绘制的子控件
		// (TDGrid.DirectPaint 不调用 inherited).
		if c.OnGridPaint != nil {
			ax, ay := c.AbsX(), c.AbsY()
			for row := 0; row < c.RowCount; row++ {
				for col := 0; col < c.ColCount; col++ {
					sel := c.SelectCol == col && c.SelectRow == row
					c.OnGridPaint(c, col, row, ax+col*c.ColWidth, ay+row*c.RowHeight, c.ColWidth, c.RowHeight, sel, proj)
				}
			}
		}
		return
	}

	if c.OnDirectPaint != nil {
		c.OnDirectPaint(c, proj)
	} else if c.WLib != nil {
		m.BlitImage(c.WLib, c.FaceIndex, c.AbsX(), c.AbsY(), proj)
	}

	for _, ch := range c.Children {
		if ch.Visible {
			m.paintControl(ch, proj)
		}
	}
}

// BlitImage 在绝对屏幕坐标处绘制 WIL 图像 (控件的默认绘制,
// DWinCtl.pas:631-635). 不可用时返回 false.
func (m *UIManager) BlitImage(f *wil.File, idx, x, y int, proj [16]float32) bool {
	if f == nil || m.resources == nil {
		return false
	}
	img := f.GetImage(idx)
	if img == nil || img.RGBA == nil {
		return false
	}
	tex := m.resources.GetTexture(f, idx)
	if tex == 0 {
		return false
	}
	m.gl.DrawQuad(tex, float32(x), float32(y), float32(img.Width), float32(img.Height), proj)
	return true
}

// ---------------------------------------------------------------------------
// 调试检查 (DebugTree 供布局快照测试; debugCollectHits 供悬停检测)
// ---------------------------------------------------------------------------

// DebugTree 返回缩进的控件树文本, maxDepth 限制展开深度。
func (m *UIManager) DebugTree(maxDepth int) string {
	var sb strings.Builder
	m.debugTreeWalk(&sb, m.Root, 0, maxDepth)
	if m.Modal != nil {
		sb.WriteString("(modal)\n")
		m.debugTreeWalk(&sb, m.Modal, 0, maxDepth)
	}
	return strings.TrimRight(sb.String(), "\n")
}

func (m *UIManager) debugTreeWalk(sb *strings.Builder, c *UIControl, depth, maxDepth int) {
	if depth > maxDepth {
		return
	}
	w, h := c.effectiveSize()
	vis := "vis"
	if !c.Visible {
		vis = "HID"
	}
	fmt.Fprintf(sb, "%s%s (%s) [%d,%d %dx%d] %s",
		strings.Repeat("  ", depth), c.Name, c.Kind, c.AbsX(), c.AbsY(), w, h, vis)
	if len(c.Children) > 0 {
		fmt.Fprintf(sb, " %dch", len(c.Children))
	}
	sb.WriteString("\n")
	for _, ch := range c.Children {
		m.debugTreeWalk(sb, ch, depth+1, maxDepth)
	}
}

func debugCtlName(c *UIControl) string {
	if c == nil {
		return "none"
	}
	return c.Name
}

// debugCollectHits 以树序收集所有 InRange 命中 (x,y 处于 c 的父空间)。
func (m *UIManager) debugCollectHits(c *UIControl, x, y int, hits *[]*UIControl) {
	if !c.Visible {
		return
	}
	if c.InRange(x, y) {
		*hits = append(*hits, c)
	}
	if c.Kind != KindGrid {
		for _, ch := range c.Children {
			m.debugCollectHits(ch, x-c.Left, y-c.Top, hits)
		}
	}
}

// HoveredControl 返回鼠标坐标处最顶层的可交互控件 (Button/Grid/EnableFocus)。
func (m *UIManager) HoveredControl(x, y int) *UIControl {
	var hits []*UIControl
	m.debugCollectHits(m.Root, x-m.Root.Left, y-m.Root.Top, &hits)
	for i := len(hits) - 1; i >= 0; i-- {
		c := hits[i]
		if c.Background || c.Kind == KindWindow {
			continue
		}
		if c.Kind == KindButton || c.Kind == KindGrid || c.EnableFocus {
			return c
		}
	}
	return nil
}

// HoverProvidesHint 报告鼠标下的最顶层非背景控件是否自带提示逻辑
//（OnMouseMove，如底栏按钮悬停文案、商店列表行、背包格）。场景层据此
// 决定移动到"消费移动但不产提示"的区域（底板背景等）时清除残留提示。
func (m *UIManager) HoverProvidesHint(x, y int) bool {
	var hits []*UIControl
	m.debugCollectHits(m.Root, x-m.Root.Left, y-m.Root.Top, &hits)
	for i := len(hits) - 1; i >= 0; i-- {
		c := hits[i]
		if c.Background {
			continue
		}
		return c.OnMouseMove != nil
	}
	return false
}

// SetHoverPos 记录最近一帧光标位置, 供包围盒"仅光标下显示名字"模式 (BoundsNames=1) 取悬停控件。
func (m *UIManager) SetHoverPos(x, y int) { m.hoverX, m.hoverY = x, y }

// RenderDebugBounds 绘制可见控件包围盒 (按类型着色: button=绿 window=蓝 grid=黄 control=灰)。
// 名字密度由 BoundsNames 控制: 0=仅框; 1=仅光标悬停控件及其祖先; 2=全部名字(带底条+去重)。
// 网格(KindGrid)子格在模式2不画名字, 避免背包等密集区标签成团。
func (m *UIManager) RenderDebugBounds(proj [16]float32) {
	ctx := &boundsCtx{m: m, mode: m.BoundsNames}
	if m.BoundsNames == 1 {
		if h := m.HoveredControl(m.hoverX, m.hoverY); h != nil {
			ctx.hot = map[*UIControl]bool{}
			for p := h; p != nil; p = p.Parent {
				ctx.hot[p] = true
			}
		}
	}
	ctx.walk(m.Root, proj)
	if m.Modal != nil && m.Modal.Visible {
		ctx.walk(m.Modal, proj)
	}
}

// RenderDebugImgs 在每个可见控件位置绘制其绑定的 WIL 实际图像, 独立于
// 常态渲染路径——常态被抑制 (如底栏按钮仅按下时画面) 或被 OnDirectPaint
// 覆盖的控件也能看到其 FaceIndex 图素。
func (m *UIManager) RenderDebugImgs(proj [16]float32) {
	m.debugImgsWalk(m.Root, proj)
	if m.Modal != nil && m.Modal.Visible {
		m.debugImgsWalk(m.Modal, proj)
	}
}

func (m *UIManager) debugImgsWalk(c *UIControl, proj [16]float32) {
	if !c.Visible {
		return
	}
	if c.WLib != nil {
		m.BlitImage(c.WLib, c.FaceIndex, c.AbsX(), c.AbsY(), proj)
	}
	for _, ch := range c.Children {
		m.debugImgsWalk(ch, proj)
	}
}

// RenderDebug 调试叠加层统一入口, 按开关依次绘制:
// WIL 实际图像 -> 包围盒 -> 悬停拾取。各场景渲染末尾调用一次。
func (m *UIManager) RenderDebug(proj [16]float32) {
	if m.ShowImgs {
		m.RenderDebugImgs(proj)
	}
	if m.ShowBounds {
		m.RenderDebugBounds(proj)
	}
	if m.ShowInspect {
		m.RenderDebugInspect(proj)
	}
}

// RenderDebugInspect 悬停拾取: 白框高亮光标下控件, 旁边一行属性
// (路径/类型/矩形/可见性/WIL 来源), 复用包围盒标签的避让逻辑。
func (m *UIManager) RenderDebugInspect(proj [16]float32) {
	h := m.HoveredControl(m.hoverX, m.hoverY)
	if h == nil {
		return
	}
	w, hgt := h.effectiveSize()
	x, y := float32(h.AbsX()), float32(h.AbsY())
	drawWireRect(m.gl, x, y, float32(w), float32(hgt), 1, 1, 1, 0.9, proj)
	ctx := &boundsCtx{m: m, mode: 2}
	ctx.drawLabel(controlDesc(h), x, y, float32(w), float32(hgt), 1, 1, 1, proj)
}

// controlDesc 生成 inspector 的一行属性描述。
func controlDesc(c *UIControl) string {
	w, h := c.effectiveSize()
	desc := fmt.Sprintf("%s (%s) [%d,%d %dx%d] vis=%v",
		c.DebugPath(), c.Kind, c.AbsX(), c.AbsY(), w, h, c.Visible)
	if c.WLib != nil {
		desc += fmt.Sprintf(" wil=%s:%d", filepath.Base(c.WLib.Path()), c.FaceIndex)
	}
	return desc
}

type boundsCtx struct {
	m      *UIManager
	mode   int
	hot    map[*UIControl]bool
	placed [][4]float32 // 已放置的标签矩形 (x,y,w,h), 用于去重
}

func (ctx *boundsCtx) walk(c *UIControl, proj [16]float32) {
	if !c.Visible {
		return
	}
	if !c.Background {
		ctx.drawOne(c, proj)
	}
	for _, ch := range c.Children {
		ctx.walk(ch, proj)
	}
}

// drawOne 绘制单个控件的类型着色包围盒, 名字密度由 mode 决定。
func (ctx *boundsCtx) drawOne(c *UIControl, proj [16]float32) {
	w, h := c.effectiveSize()
	if w <= 0 || h <= 0 {
		return
	}
	x, y := float32(c.AbsX()), float32(c.AbsY())

	var r, g, b float32
	switch c.Kind {
	case KindButton:
		r, g, b = 0, 1, 0
	case KindWindow:
		r, g, b = 0.3, 0.6, 1
	case KindGrid:
		r, g, b = 1, 0.9, 0
	default:
		r, g, b = 0.7, 0.7, 0.7
	}
	drawWireRect(ctx.m.gl, x, y, float32(w), float32(h), r, g, b, 0.9, proj)

	parentIsGrid := c.Parent != nil && c.Parent.Kind == KindGrid
	drawName := false
	switch ctx.mode {
	case 2:
		drawName = !parentIsGrid
	case 1:
		drawName = ctx.hot[c]
	}
	if drawName {
		ctx.drawLabel(c.Name, x, y, float32(w), float32(h), r, g, b, proj)
	}
}

// drawLabel 在控件附近画带半透明底条的名字, 并贪心下移避让已放置标签, 保证可读不重叠。
func (ctx *boundsCtx) drawLabel(name string, x, y, w, h, r, g, b float32, proj [16]float32) {
	text := ctx.m.text
	if text == nil || name == "" {
		return
	}
	lineH := float32(text.LineHeight())
	if lineH <= 0 {
		lineH = 14
	}
	const padX = 3.0
	chipW := float32(text.MeasureText(name)) + padX*2
	chipH := lineH + 3
	screenW := float32(ScreenWidth)
	screenH := float32(ScreenHeight)
	cx := x
	if cx+chipW > screenW {
		cx = screenW - chipW
	}
	if cx < 0 {
		cx = 0
	}
	// 候选纵坐标: 框上方 -> 框下方 -> 框内顶 -> 依次下移堆叠
	cands := []float32{y - chipH, y + h, y}
	for k := 1; k <= 6; k++ {
		cands = append(cands, y-chipH+float32(k)*chipH)
	}
	chosenY := y - chipH
	if chosenY < 0 {
		chosenY = 0
	}
	for _, cy := range cands {
		if cy < 0 || cy+chipH > screenH {
			continue
		}
		if !ctx.overlaps(cx, cy, chipW, chipH) {
			chosenY = cy
			break
		}
	}
	ctx.placed = append(ctx.placed, [4]float32{cx, chosenY, chipW, chipH})
	ctx.m.gl.DrawQuadColor(cx, chosenY, chipW, chipH, 0, 0, 0, 0.75, proj)
	text.DrawText(name, cx+padX, chosenY+1, r, g, b, 1, proj)
}

func (ctx *boundsCtx) overlaps(x, y, w, h float32) bool {
	for _, p := range ctx.placed {
		if x+w <= p[0] || p[0]+p[2] <= x || y+h <= p[1] || p[1]+p[3] <= y {
			continue
		}
		return true
	}
	return false
}
