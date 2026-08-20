//go:build !windows

package engine

// 非 Windows 平台无 IME 平台层集成。已提交字符仍可用:
//   - X11: GLFW 经 XIM 交付 (锁定 XIMPreeditNothing 风格, 无内联预编辑,
//     候选窗是输入法自带浮动窗, 位置不受游戏控制);
//   - Wayland: GLFW 经 text-input-v3 交付提交字符。
//
// 能力边界与后续方向见 doc/go版完善路线图.md。
func platformIMEInit(w *Window) {}

func platformIMESetRect(x, y, width, height int) {}

func platformIMEShutdown() {}
