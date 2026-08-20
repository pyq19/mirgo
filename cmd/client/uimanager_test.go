package main

// 框架分发语义测试（无需 GL——不覆盖 Paint）。
// 验证 DWinCtl 移植：按钮抬起且在内部才触发点击、网格按下==抬起选中、
// 窗口拖拽带边界限制、模态窗口拦截、背景点击+焦点释放、键盘焦点路由。

import "testing"

func newTestManager() *UIManager {
	return NewUIManager(nil, nil, nil)
}

func TestButtonClickOnUpInside(t *testing.T) {
	m := newTestManager()
	btn := NewUIControl("btn", KindButton)
	btn.Left, btn.Top, btn.Width, btn.Height = 100, 100, 40, 30
	clicks := 0
	btn.OnClick = func(c *UIControl, x, y int) { clicks++ }
	btn.Visible = true
	m.Root.AddChild(btn)

	// 内部按下，内部抬起 → 一次点击。
	if !m.RouteMouseDown(110, 110, 0) {
		t.Fatal("mouse down inside button should be consumed")
	}
	if !btn.Downed {
		t.Fatal("button should be Downed after press")
	}
	if m.Capture != btn {
		t.Fatal("button should own the capture after press")
	}
	if !m.RouteMouseUp(120, 115, 0) {
		t.Fatal("mouse up inside button should be consumed")
	}
	if clicks != 1 {
		t.Fatalf("expected 1 click, got %d", clicks)
	}
	if btn.Downed {
		t.Fatal("button should reset Downed after release")
	}

	// 内部按下，拖出范围，外部抬起 → 不触发点击。
	m.RouteMouseDown(110, 110, 0)
	m.RouteMouseMove(300, 300)
	if btn.Downed {
		t.Fatal("button Downed should track cursor out of range while captured")
	}
	m.RouteMouseUp(300, 300, 0)
	if clicks != 1 {
		t.Fatalf("expected still 1 click after drag-out, got %d", clicks)
	}
}

func TestGridSelectDownEqualsUp(t *testing.T) {
	m := newTestManager()
	grid := NewUIControl("grid", KindGrid)
	grid.Left, grid.Top = 33, 43
	grid.ColCount, grid.RowCount = 8, 6
	grid.ColWidth, grid.RowHeight = 36, 32
	var selCol, selRow int
	selects := 0
	grid.OnGridSelect = func(c *UIControl, col, row int) {
		selects++
		selCol, selRow = col, row
	}
	grid.Visible = true
	m.Root.AddChild(grid)

	// 格子 (2,1): x = 33 + 2*36 + 5 = 110, y = 43 + 1*32 + 5 = 80。
	m.RouteMouseDown(110, 80, 0)
	if m.Capture != grid {
		t.Fatal("grid should own capture after cell press")
	}
	// 同一格子抬起 → 选中。
	m.RouteMouseUp(115, 85, 0)
	if selects != 1 || selCol != 2 || selRow != 1 {
		t.Fatalf("expected select (2,1) once, got %d selects col=%d row=%d", selects, selCol, selRow)
	}

	// 按下格子 (0,0)，抬起格子 (3,2) → 不选中。
	m.RouteMouseDown(40, 50, 0)
	m.RouteMouseUp(33+3*36+5, 43+2*32+5, 0)
	if selects != 1 {
		t.Fatalf("drag across cells must not select, got %d selects", selects)
	}
}

func TestWindowDragAndClamp(t *testing.T) {
	m := newTestManager()
	win := NewUIControl("win", KindWindow)
	win.Floating = true
	win.Left, win.Top, win.Width, win.Height = 200, 100, 120, 120
	win.Visible = true
	m.Root.AddChild(win)

	// 按下窗口主体 → 记录抓取点。
	if !m.RouteMouseDown(210, 110, 0) {
		t.Fatal("window should consume press")
	}
	if win.SpotX != 210 || win.SpotY != 110 {
		t.Fatalf("grab spot not recorded: %d,%d", win.SpotX, win.SpotY)
	}
	// 拖拽 +50,+30 → 窗口移动。
	m.RouteMouseMove(260, 140)
	if win.Left != 250 || win.Top != 130 {
		t.Fatalf("window did not follow drag: %d,%d", win.Left, win.Top)
	}
	// 拖到远超右边界限制处。
	m.RouteMouseMove(2000, 140)
	maxLeft := m.Root.Width - 60
	if win.Left > maxLeft {
		t.Fatalf("window Left %d exceeds max %d", win.Left, maxLeft)
	}
}

func TestModalShortCircuit(t *testing.T) {
	m := newTestManager()
	behind := NewUIControl("behind", KindButton)
	behind.Left, behind.Top, behind.Width, behind.Height = 100, 100, 40, 30
	behindClicks := 0
	behind.OnClick = func(c *UIControl, x, y int) { behindClicks++ }
	behind.Visible = true
	m.Root.AddChild(behind)

	modal := NewUIControl("modal", KindWindow)
	modal.Left, modal.Top, modal.Width, modal.Height = 90, 90, 200, 150
	modalDown := 0
	modal.OnMouseDown = func(c *UIControl, button, x, y int) { modalDown++ }
	m.ShowModal(modal)

	// 点击 "behind" 所在位置——模态窗口必须拦截。
	if !m.RouteMouseDown(110, 110, 0) {
		t.Fatal("modal should consume all mouse input")
	}
	m.RouteMouseUp(110, 110, 0)
	if behindClicks != 0 {
		t.Fatalf("control behind modal was clicked (%d)", behindClicks)
	}
	if modalDown != 1 {
		t.Fatalf("modal got %d downs, want 1", modalDown)
	}
	m.CloseModal(modal)
}

// 修复 1.1：模态可见时点在模态框**矩形之外**也必须恒被吞掉
// （Delphi TDWinManager 恒返回 TRUE，DWinCtl.pas:985,1012），
// 否则点击会落到游戏世界导致走路/攻击。
func TestModalConsumesClicksOutsideItsRect(t *testing.T) {
	m := newTestManager()
	modal := NewUIControl("modal", KindWindow)
	modal.Left, modal.Top, modal.Width, modal.Height = 90, 90, 100, 80
	m.ShowModal(modal)

	// (500,500) 远在模态框之外。
	if !m.RouteMouseDown(500, 500, 0) {
		t.Fatal("mouse down outside modal rect must still be consumed")
	}
	if !m.RouteMouseUp(500, 500, 0) {
		t.Fatal("mouse up outside modal rect must still be consumed")
	}
	if !m.RouteMouseMove(500, 500) {
		t.Fatal("mouse move outside modal rect must still be consumed")
	}
	if !m.RouteDblClick(500, 500) {
		t.Fatal("dblclick outside modal rect must still be consumed")
	}
	m.CloseModal(modal)
	// 模态关闭后恢复正常命中：空白处不再被吞。
	if m.RouteMouseDown(500, 500, 0) {
		t.Fatal("after CloseModal, empty area should not consume down")
	}
}

func TestBackgroundClickAndFocusRelease(t *testing.T) {
	m := newTestManager()
	bgClicks := 0
	m.Root.OnBackgroundClick = func(c *UIControl) { bgClicks++ }

	btn := NewUIControl("btn", KindButton)
	btn.Left, btn.Top, btn.Width, btn.Height = 10, 10, 30, 30
	btn.Visible = true
	m.Root.AddChild(btn)

	// 点击按钮 → 获得焦点。
	m.RouteMouseDown(15, 15, 0)
	m.RouteMouseUp(15, 15, 0)
	if m.Focused != btn {
		t.Fatal("pressing EnableFocus button should focus it")
	}
	// 点击空白区域 → 触发背景点击并清除焦点。
	m.RouteMouseDown(500, 500, 0)
	if bgClicks != 1 {
		t.Fatalf("expected 1 background click, got %d", bgClicks)
	}
	if m.Focused != nil {
		t.Fatal("background click should release focus")
	}
	// WantReturn=true → 背景点击消费事件。
	m.Root.OnBackgroundClick = func(c *UIControl) { c.WantReturn = true }
	if !m.RouteMouseDown(500, 500, 0) {
		t.Fatal("background click with WantReturn should be consumed")
	}
}

func TestKeyboardFocusRouting(t *testing.T) {
	m := newTestManager()
	btn := NewUIControl("btn", KindButton)
	btn.Left, btn.Top, btn.Width, btn.Height = 10, 10, 30, 30
	got := ""
	btn.OnChar = func(c *UIControl, ch rune) { got += string(ch) }
	btn.Visible = true
	m.Root.AddChild(btn)

	// 无焦点 → 字符不投递。
	if m.RouteChar('a') {
		t.Fatal("char should not route without focus")
	}
	// 通过点击使按钮获得焦点 → 字符投递成功。
	m.RouteMouseDown(15, 15, 0)
	m.RouteMouseUp(15, 15, 0)
	if !m.RouteChar('a') || !m.RouteChar('b') {
		t.Fatal("focused control should receive chars")
	}
	if got != "ab" {
		t.Fatalf("expected 'ab', got %q", got)
	}
}

func TestRaiseToTopChangesHitOrder(t *testing.T) {
	m := newTestManager()
	w1 := NewUIControl("w1", KindWindow)
	w1.Floating = true
	w1.Left, w1.Top, w1.Width, w1.Height = 100, 100, 100, 100
	w2 := NewUIControl("w2", KindWindow)
	w2.Floating = true
	w2.Left, w2.Top, w2.Width, w2.Height = 150, 150, 100, 100
	who := ""
	w1.OnMouseDown = func(c *UIControl, b, x, y int) { who = "w1" }
	w2.OnMouseDown = func(c *UIControl, b, x, y int) { who = "w2" }
	w1.Visible = true
	w2.Visible = true
	m.Root.AddChild(w1)
	m.Root.AddChild(w2)

	// 重叠区域 (160,160): w2 后添加 → 在最上层 → 先命中。
	m.RouteMouseDown(160, 160, 0)
	if who != "w2" {
		t.Fatalf("expected w2 on top, got %q", who)
	}
	// 点击 w1 非重叠区域 → 提升到 w2 之上。
	m.RouteMouseUp(160, 160, 0)
	m.RouteMouseDown(110, 110, 0)
	m.RouteMouseUp(110, 110, 0)
	if who != "w1" {
		t.Fatalf("expected w1 hit, got %q", who)
	}
	// 此时重叠区域应命中 w1（已提升）。
	m.RouteMouseDown(160, 160, 0)
	if who != "w1" {
		t.Fatalf("expected w1 raised on top, got %q", who)
	}
}

// 焦点控件没有 OnChar/OnKeyDown 时必须放行 (修复 M1/M2):
// Go 单通道下, "无处理器也吞" 会让点过 HUD 后聊天完全打不了字。
func TestFocusWithoutHandlerPassesThrough(t *testing.T) {
	m := newTestManager()
	btn := NewUIControl("btn", KindButton)
	btn.Left, btn.Top, btn.Width, btn.Height = 10, 10, 30, 30
	btn.Visible = true
	m.Root.AddChild(btn)

	// 点击按钮获得焦点 (KindButton 默认 EnableFocus=true)。
	m.RouteMouseDown(15, 15, 0)
	m.RouteMouseUp(15, 15, 0)
	if m.Focused != btn {
		t.Fatal("pressing button should focus it")
	}
	// 按钮没有 OnChar/OnKeyDown → 字符/按键放行给场景。
	if m.RouteChar('a') {
		t.Fatal("char must pass through when focused control has no OnChar")
	}
	if m.RouteKeyDown(257) {
		t.Fatal("key must pass through when focused control has no OnKeyDown")
	}
}

// 焦点控件自身或祖先不可见时释放焦点并放行 (修复 M1/M2,
// 补 Delphi DWinCtl.pas:933-936 的缺失检查)。
func TestInvisibleFocusReleased(t *testing.T) {
	m := newTestManager()
	win := NewUIControl("win", KindWindow)
	win.Left, win.Top, win.Width, win.Height = 100, 100, 200, 100
	win.Visible = true
	m.Root.AddChild(win)
	edit := NewUIControl("edit", KindControl)
	edit.Left, edit.Top, edit.Width, edit.Height = 10, 10, 100, 20
	edit.Visible = true
	got := ""
	edit.OnChar = func(c *UIControl, ch rune) { got += string(ch) }
	win.AddChild(edit)
	m.SetFocus(edit)

	// 正常状态: 字符投递成功。
	if !m.RouteChar('x') {
		t.Fatal("visible focused edit should receive chars")
	}
	// 父窗口隐藏 → 焦点失效 → 放行并释放焦点。
	win.Visible = false
	if m.RouteChar('y') {
		t.Fatal("char must pass through when focused ancestor is hidden")
	}
	if m.Focused != nil {
		t.Fatal("invisible focus should be released")
	}
	if got != "x" {
		t.Fatalf("expected 'x', got %q", got)
	}
}

// 模态可见时键盘恒吞 (修复 M1/M2): 即使模态与焦点都没有处理器,
// 也不能漏键给游戏 (对话框开着时不许放技能/开聊天)。
func TestModalConsumesKeyboardWithoutHandler(t *testing.T) {
	m := newTestManager()
	modal := NewUIControl("modal", KindWindow)
	m.ShowModal(modal) // 焦点落在模态自身, 无 OnChar/OnKeyDown

	if !m.RouteChar('a') {
		t.Fatal("visible modal must consume chars even without handler")
	}
	if !m.RouteKeyDown(257) {
		t.Fatal("visible modal must consume keys even without handler")
	}

	// 模态不可见但引用残留 → 自清理后回到焦点路径
	// (DWinCtl.pas:924-929); 焦点(模态自身)不可见 → 释放放行。
	modal.Visible = false
	if m.RouteChar('b') {
		t.Fatal("stale invisible modal must not consume chars")
	}
	if m.Modal != nil {
		t.Fatal("stale invisible modal should be self-cleaned")
	}
	if m.Focused != nil {
		t.Fatal("invisible focus should be released")
	}
}

// 模态内聚焦子控件正常收字符; CloseModal 释放子树内焦点 (修复 M4)。
func TestModalCharDeliveryAndCloseModalFocusRelease(t *testing.T) {
	m := newTestManager()
	modal := NewUIControl("dlg", KindWindow)
	edit := NewUIControl("edit", KindControl)
	got := ""
	edit.OnChar = func(c *UIControl, ch rune) { got += string(ch) }
	modal.AddChild(edit)
	m.ShowModal(modal)
	m.SetFocus(edit)

	if !m.RouteChar('h') || !m.RouteChar('i') {
		t.Fatal("focused edit inside modal should receive chars")
	}
	if got != "hi" {
		t.Fatalf("expected 'hi', got %q", got)
	}

	m.CloseModal(modal)
	if m.Focused != nil {
		t.Fatal("CloseModal must release focus inside the closed subtree")
	}

	// 焦点在模态子树之外时, CloseModal 不动它。
	m2 := newTestManager()
	outside := NewUIControl("btn", KindButton)
	outside.Left, outside.Top, outside.Width, outside.Height = 10, 10, 30, 30
	outside.Visible = true
	m2.Root.AddChild(outside)
	m2.SetFocus(outside)
	modal2 := NewUIControl("dlg2", KindWindow)
	modal2.EnableFocus = false // ShowModal 不抢焦点
	m2.ShowModal(modal2)
	m2.CloseModal(modal2)
	if m2.Focused != outside {
		t.Fatal("CloseModal must not touch focus outside the closed subtree")
	}
}

func TestControlDesc(t *testing.T) {
	m := newTestManager()
	btn := NewUIControl("DSW", KindButton)
	btn.Left, btn.Top, btn.Width, btn.Height = 25, 141, 34, 31
	btn.Visible = true
	m.Root.AddChild(btn)

	got := controlDesc(btn)
	want := "DBackground>DSW (button) [25,141 34x31] vis=true"
	if got != want {
		t.Fatalf("controlDesc = %q, want %q", got, want)
	}

	btn.Visible = false
	wantHidden := "DBackground>DSW (button) [25,141 34x31] vis=false"
	if got := controlDesc(btn); got != wantHidden {
		t.Fatalf("controlDesc = %q, want %q", got, wantHidden)
	}
}
