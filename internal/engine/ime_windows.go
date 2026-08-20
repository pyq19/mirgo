//go:build windows

package engine

// Windows IME 集成 (IMM32) — 对齐 Delphi 原版 TEdit 的中文输入体验:
// 组合窗/候选窗定位到输入框, 预编辑串回报给游戏渲染。
//
// GLFW 3.4 自身不处理 WM_IME_* (win32_window.c 无 IMM32 调用), 提交
// 字符经 WM_CHAR → GLFW char 回调, 这条路不受影响。这里子类化窗口
// 过程只补三件事: 组合窗定位、候选窗定位、预编辑串提取。
//
// 不依赖 GLFW native API (glfwGetWin32Window 的返回类型是 glfw 包私有
// 的 cgo 类型, 跨包不可用): 用 EnumWindows 按进程 ID 找到自己的窗口。
// 该文件所有代码只在主线程执行 (GLFW 消息循环线程 = engine init 里
// LockOSThread 的线程), 无竞态。

/*
#cgo LDFLAGS: -limm32
#include <windows.h>
#include <imm.h>
#include <stdlib.h>

extern void goIMEPreedit(char* textUTF8, long long cursor);

static WNDPROC ime_prevProc = NULL;
static HWND    ime_hwnd     = NULL;
static RECT    ime_rect     = {0, 0, 0, 0};
static int     ime_hasRect  = 0;
static int     ime_composing = 0;

static void imeApplyPosition(HWND hwnd) {
	if (!ime_hasRect) return;
	HIMC himc = ImmGetContext(hwnd);
	if (!himc) return;
	COMPOSITIONFORM cf;
	ZeroMemory(&cf, sizeof(cf));
	cf.dwStyle = CFS_POINT;
	cf.ptCurrentPos.x = ime_rect.left;
	cf.ptCurrentPos.y = ime_rect.top;
	ImmSetCompositionWindow(himc, &cf);
	CANDIDATEFORM cd;
	ZeroMemory(&cd, sizeof(cd));
	cd.dwIndex = 0;
	cd.dwStyle = CFS_CANDIDATEPOS;
	cd.ptCurrentPos.x = ime_rect.left;
	cd.ptCurrentPos.y = ime_rect.bottom;
	ImmSetCandidateWindow(himc, &cd);
	ImmReleaseContext(hwnd, himc);
}

static void imeReportCompStr(HIMC himc) {
	LONG bytes = ImmGetCompositionStringW(himc, GCS_COMPSTR, NULL, 0);
	if (bytes <= 0) {
		goIMEPreedit(NULL, 0);
		return;
	}
	wchar_t* wbuf = (wchar_t*)malloc(bytes + sizeof(wchar_t));
	if (!wbuf) return;
	ImmGetCompositionStringW(himc, GCS_COMPSTR, wbuf, bytes + sizeof(wchar_t));
	wbuf[bytes / sizeof(wchar_t)] = 0;
	int n = WideCharToMultiByte(CP_UTF8, 0, wbuf, -1, NULL, 0, NULL, NULL);
	char* utf8 = (char*)malloc(n > 0 ? n : 1);
	if (utf8) {
		WideCharToMultiByte(CP_UTF8, 0, wbuf, -1, utf8, n, NULL, NULL);
		LONG cp = ImmGetCompositionStringW(himc, GCS_CURSORPOS, NULL, 0);
		goIMEPreedit(utf8, (long long)(cp & 0xFFFF));
		free(utf8);
	}
	free(wbuf);
}

static LRESULT CALLBACK imeWndProc(HWND hwnd, UINT msg, WPARAM wp, LPARAM lp) {
	switch (msg) {
	case WM_IME_STARTCOMPOSITION:
		ime_composing = 1;
		imeApplyPosition(hwnd);
		return CallWindowProcW(ime_prevProc, hwnd, msg, wp, lp);
	case WM_IME_COMPOSITION: {
		// 先交给原窗口过程 (GLFW): GCS_RESULTSTR 由 DefWindowProc
		// 转成 WM_CHAR → GLFW char 回调, 提交字符不重复上屏。
		LRESULT r = CallWindowProcW(ime_prevProc, hwnd, msg, wp, lp);
		HIMC himc = ImmGetContext(hwnd);
		if (himc) {
			if (lp & GCS_COMPSTR) imeReportCompStr(himc);
			ImmReleaseContext(hwnd, himc);
		}
		if (lp & GCS_RESULTSTR) goIMEPreedit(NULL, 0);
		return r;
	}
	case WM_IME_ENDCOMPOSITION:
		ime_composing = 0;
		goIMEPreedit(NULL, 0);
		return CallWindowProcW(ime_prevProc, hwnd, msg, wp, lp);
	case WM_IME_SETCONTEXT:
		if (wp) imeApplyPosition(hwnd);
		return CallWindowProcW(ime_prevProc, hwnd, msg, wp, lp);
	}
	return CallWindowProcW(ime_prevProc, hwnd, msg, wp, lp);
}

static BOOL CALLBACK imeEnumWnd(HWND hwnd, LPARAM lp) {
	DWORD pid = 0;
	GetWindowThreadProcessId(hwnd, &pid);
	if (pid == (DWORD)GetCurrentProcessId() && IsWindowVisible(hwnd)) {
		ime_hwnd = hwnd;
		return FALSE;
	}
	return TRUE;
}

static int imeInstall(void) {
	EnumWindows(imeEnumWnd, 0);
	if (!ime_hwnd) return 0;
	ime_prevProc = (WNDPROC)SetWindowLongPtrW(ime_hwnd, GWLP_WNDPROC, (LONG_PTR)imeWndProc);
	return ime_prevProc != NULL;
}

static void imeUninstall(void) {
	if (ime_hwnd && ime_prevProc) {
		SetWindowLongPtrW(ime_hwnd, GWLP_WNDPROC, (LONG_PTR)ime_prevProc);
	}
	ime_hwnd = NULL;
	ime_prevProc = NULL;
}

static void imeSetRect(int x, int y, int w, int h) {
	ime_rect.left = x;
	ime_rect.top = y;
	ime_rect.right = x + w;
	ime_rect.bottom = y + h;
	ime_hasRect = (w > 0 && h > 0);
	if (ime_hwnd && ime_composing) imeApplyPosition(ime_hwnd);
}
*/
import "C"

var imeWindow *Window

//export goIMEPreedit
func goIMEPreedit(textUTF8 *C.char, cursor int64) {
	if imeWindow == nil {
		return
	}
	if textUTF8 == nil {
		imeWindow.preedit = ""
		imeWindow.preeditCur = 0
		return
	}
	imeWindow.preedit = C.GoString(textUTF8)
	imeWindow.preeditCur = int(cursor)
}

func platformIMEInit(w *Window) {
	imeWindow = w
	C.imeInstall()
}

func platformIMESetRect(x, y, width, height int) {
	C.imeSetRect(C.int(x), C.int(y), C.int(width), C.int(height))
}

func platformIMEShutdown() {
	C.imeUninstall()
	imeWindow = nil
}
