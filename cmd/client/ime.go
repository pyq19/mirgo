package main

import "github.com/pyq0109/mirgo/internal/engine"

// IME 桥接。GLFW 无原生 IME API, 文本输入分两层:
//   - 已提交字符: 所有平台经 char 回调送达 (X11 走 XIM, Wayland 走
//     text-input-v3, Windows 走 WM_CHAR);
//   - 预编辑串与候选窗定位: 仅 Windows 平台层支持 (IMM32,
//     internal/engine/ime_windows.go), 其他平台 imePreedit 恒为空。

// mainWindow 在 main() 里设置, 供桥接函数访问引擎窗口。
var mainWindow *engine.Window

// imePreedit 是当前预编辑 (组合中) 文本, imePreeditCursor 是其光标
// 偏移 (rune)。由主循环每帧从引擎窗口刷新。
var imePreedit string
var imePreeditCursor int

func imeRefresh() {
	if mainWindow == nil {
		return
	}
	imePreedit, imePreeditCursor = mainWindow.Preedit()
}

// notifyTextInputRect 向平台层通告当前输入框位置, 用于候选窗/组合窗
// 定位。输入为逻辑 800×600 坐标, 内部换算为窗口像素。
func notifyTextInputRect(x, y, w, h int) {
	if mainWindow == nil || winW <= 0 || winH <= 0 {
		return
	}
	mainWindow.SetTextInputRect(
		x*winW/ScreenWidth, y*winH/ScreenHeight,
		w*winW/ScreenWidth, h*winH/ScreenHeight,
	)
}

// clearTextInputRect 表示当前没有活动的文本输入。
func clearTextInputRect() {
	if mainWindow != nil {
		mainWindow.SetTextInputRect(0, 0, 0, 0)
	}
}
