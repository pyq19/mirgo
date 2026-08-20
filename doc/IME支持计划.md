# Linux/macOS 完整 IME 支持计划

> 状态:规划中(未排期)。前置工作已完成:提交字符全平台通路 + Windows IMM32 完整支持(见文末"现状")。
> 本文档基于对 vendored GLFW 3.4 源码的逐行核查,所有技术判断带 `文件:行号` 引用。

---

## 一、现状与能力差距

### 1.1 已实现(批次8)

- **提交字符通路**:所有平台经 GLFW char 回调 → 场景 OnChar → EditBox/聊天输入。
- **Windows IMM32**(`internal/engine/ime_windows.go`):子类化窗口过程,组合窗/候选窗定位到输入框、预编辑串渲染进 EditBox 与聊天行、提交字符走原 WM_CHAR 不重复上屏。
- **引擎抽象已就绪**(`internal/engine/window.go`):`Window.SetTextInputRect(x,y,w,h)` 与 `Window.Preedit()`;平台层钩子 `platformIMEInit/platformIMESetRect/platformIMEShutdown`(`ime_windows.go` / `ime_stub.go` 双实现)。
- **客户端零耦合**:`cmd/client/ime.go` 桥接 + EditBox.paint/聊天行渲染已接入——**后续任何平台层实现都不需要再动客户端代码**。

### 1.2 各平台能力矩阵(当前)

| 平台 | 提交字符(能打进去) | 预编辑内联显示 | 候选窗/组合窗定位 |
|------|:---:|:---:|:---:|
| Windows | ✅ WM_CHAR | ✅ IMM32 | ✅ IMM32 |
| Linux/X11 | ✅ XIM(ibus/fcitx 经 XIM 桥) | ❌ | ❌(输入法浮动窗,位置不受控) |
| Linux/Wayland | ❌ **完全不可用**(见 1.3 更正) | ❌ | ❌ |
| macOS | ✅ NSTextInputClient | ❌ | ❌(候选窗停在视图原点) |

### 1.3 更正:Wayland 无 text-input 支持

此前文档/结论中"Wayland 经 text-input-v3 交付提交字符"的说法**对本项目 vendored 的 GLFW 3.4.0 不成立**。全量核查确认:GLFW 3.4 源码树中不存在任何 `zwp_text_input` 协议代码;Wayland 文本输入仅走 `inputText()` 的 XKB keysym→Unicode 转换(wl_window.c:1188-1203),只能处理组合键/死键,**CJK 输入法完全无法工作**。Wayland 会话下目前连中文都打不出——这是本计划要解决的最大缺口。

---

## 二、技术约束核查(vendored GLFW 3.4.0)

### 2.1 X11 — XIM 风格被锁死

- IC 创建固定请求 `XIMPreeditNothing | XIMStatusNothing`(glfw/src/x11_window.c:1928-1938),即主动放弃预编辑;候选窗位置无任何 XNSpotLocation 设置。
- 初始化时只在 IM 支持的风格里挑 `PreeditNothing|StatusNothing`(glfw/src/x11_init.c:443-478)。
- 事件循环为 GLFW 独占:`XFilterEvent`/`Xutf8LookupString`(x11_window.c:1274)在其内部 processEvent 中,**无法在进程外钩入自己的 XIM 客户端**——任何"绕过 GLFW 自建 XIM"的方案都不可行。

### 2.2 Wayland — 无输入法协议

- 无 zwp_text_input_v3、无 input-method 相关任何代码;`inputText()` 纯 XKB 路径(wl_window.c:1188-1203)。

### 2.3 macOS — 数据在内部,但不暴露

- GLFWContentView 已实现 `NSTextInputClient`(cocoa_window.m:330),内部维护 `markedText`(:350、:654-683),提交文本经 `insertText:`→`_glfwInputChar`(:709-737)——所以提交字符可用。
- 但 `firstRectForCharacterRange:` 返回**视图原点零尺寸矩形**(cocoa_window.m:702-707),系统候选窗因此只能停在窗口左下角;markedText 没有任何公开 API 可取。

### 2.4 结论

三个平台的缺口都在 GLFW 内部,进程外无合法钩点。**唯一彻底方案是 fork GLFW 打补丁**;客户端/引擎侧框架(1.1)已就绪,fork 落地后只需新增平台后端文件。

---

## 三、总体方案

**采用:fork go-gl/glfw + go.mod replace。**

- go-gl/glfw 仓库本身内嵌 GLFW C 源码(本模块缓存中 `glfw@v0.1.0-pre.1/glfw/src/` 即是),fork 一份即同时覆盖 Go 绑定与 C 源码,`go.mod` 加一行 replace 指向自维护 fork(如 `github.com/pyq0109/glfw/v3.4`)。
- 每个平台补丁作为独立 commit,便于上游 GLFW 发布新版时 cherry-pick 迁移;同时跟踪上游 IME 进展(上游若有正式 preedit API,回迁上游版并删除对应补丁)。
- 补丁只加不改语义:所有新行为走新增公开 API,现有 API 行为不变,保证 mapviewer/wilviewer 等工具不受影响。

**不采用的替代方案:**
- 换 SDL2/ebiten:重构整个窗口与渲染接线,成本远超收益。
- 进程外自建 XIM/text-input 客户端:事件循环被 GLFW 独占,不可行(2.1)。
- 等待上游:上游无明确排期,Wayland 用户当前完全打不了中文,不可接受。

---

## 四、分阶段实施

### 阶段 0:X11 ibus D-Bus 候选窗定位实验(免 fork 快赢,可选)

- IBus 的输入上下文 D-Bus 接口提供 `SetCursorLocation(x,y,w,h)`;实验目标:不 fork GLFW,经 D-Bus 把 ibus 浮动候选窗挪到输入框附近。
- 不确定性:XIM 桥(ibus-x11)持有的输入上下文与 D-Bus 新建上下文是否同一个、面板是否跟随,需实测验证;fcitx5 无对应 D-Bus 定位接口,不适用。
- 落点(若可行):新增 `internal/engine/ime_linux_ibus.go`(纯 Go + D-Bus,走现有 platformIMESetRect 钩子);失败则放弃,无沉没成本。
- 工作量:≤0.5 天(实验性质)。

### 阶段 1:Linux/X11 fork 补丁 — XIM 预编辑 + SpotLocation

C 侧(glfw/src/x11_init.c、x11_window.c):
1. 输入风格选择扩展:优先 `XIMPreeditCallbacks | XIMStatusNothing`,IM 不支持时回退现有 `PreeditNothing`(保住存量功能)。
2. 注册预编辑回调(PreeditStart/Draw/Done),经 `Xmb`/`Xutf8` 取预编辑文本与光标位置。
3. 新增公开 API:`glfwSetPreeditCallback(window, cb)`、`glfwSetTextInputRect(window, x, y, w, h)`;后者经 `XSetICValues(XNSpotLocation)` 更新。
4. 事件过滤适配:预编辑期间被 IM 消费的按键不得再产生 key/char 事件(XFilterEvent 语义已覆盖,重点回归测试)。

Go 绑定侧:对应包装两个新 API。
引擎侧:新增 `internal/engine/ime_linux.go`(build tag `linux && !wayland` 或运行时探测)实现 platformIME 三钩子。
客户端:零改动。

- 工作量:2–4 天。风险:fcitx5 XIM 桥对 callbacks 风格的兼容性(重点实测)。

### 阶段 2:Linux/Wayland fork 补丁 — text-input-v3(工作量最大)

C 侧(glfw/src/wl_init.c、wl_window.c、wl_platform.h + 协议生成代码):
1. wayland-scanner 生成 `zwp_text_input_v3` 客户端代码;绑定 seat 的 text_input_manager。
2. 监听器:enter/leave、`preedit_string`、`commit_string`、`delete_surrounding_text`、done。
3. `commit_string`(UTF-8)逐 codepoint 走 `_glfwInputChar` → 现有 char 回调链路;**与现有 XKB 路径共存**(无 IME 时 ASCII 仍走 XKB,需处理键盘焦点 enter/leave 时的启用状态机)。
4. `preedit_string` → 新预编辑 API(同阶段 1 的公开接口)。
5. 每次 commit 前按 `glfwSetTextInputRect` 发送 `set_cursor_rectangle`,使 fcitx5/ibus(Wayland)候选窗定位到输入框。

- 工作量:3–5 天。风险:compositor 差异(Mutter/KDE wlroots 对 text-input-v3 的实现细节)、与 XKB 路径的事件去重、surrounding text 状态机复杂度(首版可只实现最小集:不报告 surrounding text,仅接收 commit/preedit)。

### 阶段 3:macOS fork 补丁 — 暴露 markedText + 修正 firstRect(改动最小)

C/Objective-C 侧(glfw/src/cocoa_window.m):
1. `setMarkedText:` 时记录预编辑文本与光标范围,经新 API `glfwSetPreeditCallback` 上报;`unmarkText`/`insertText:` 时清空。
2. `firstRectForCharacterRange:actualRange:` 改为返回 `glfwSetTextInputRect` 存储的输入框矩形(视图坐标→屏幕坐标换算,注意 Retina 缩放与 Cocoa 左下原点)。

引擎侧:新增 `internal/engine/ime_darwin.go`。
客户端:零改动。

- 工作量:1–2 天。风险:候选窗实际跟随效果取决于系统输入法(拼音/注音)对 firstRect 的消费方式,需真机验证。

---

## 五、实施顺序与验证

**顺序建议:阶段 3(macOS,改动最小先打通管线)→ 阶段 1(X11)→ 阶段 2(Wayland)**;阶段 0 随时可插队实验。若用户主力环境是 Linux,可先做阶段 1。

**验证矩阵**(每阶段落地后):

| 维度 | 用例 |
|------|------|
| 环境 | X11×ibus、X11×fcitx5、Wayland×fcitx5、Wayland×ibus、macOS×拼音/注音、Windows 回归(IMM32 不回退) |
| 场景 | 登录/注册各 EditBox、建角名字、游戏内聊天、ShowInput 对话框、行会公告 MemoBox、调试控制台 |
| 行为 | 候选窗/组合窗出现在输入框附近;预编辑串下划线渲染于光标处;确认候选的 Enter 不误发送聊天;Esc 取消预编辑不触发游戏 Esc;预编辑中热键(B/M/F1…)不误触发;英文直接输入不受影响 |
| 回归 | `go build ./... && go vet ./... && go test ./... -count=1`;mapviewer/wilviewer 输入不受 fork 影响 |

---

## 六、维护与风险

- **fork 漂移**:每次上游 GLFW 发版评估回迁;补丁保持独立 commit + 本文档登记,迁移成本可控。
- **输入法行为差异**:ibus/fcitx5/macOS/微软拼音对 spot location、候选窗样式的消费各不相同,验收以"候选窗大致跟随 + 提交正确"为准,不追求像素级一致。
- **Wayland 生态差异**:个别 compositor 的 text-input-v3 实现不完整时,回退表现为"能提交、候选窗不跟随",属可接受降级。
- **安全**:fork 仓库只含 GLFW 源码与补丁,不引入任何新运行时依赖(D-Bus 方案若采用,依赖 `github.com/godbus/dbus/v5`,纯 Go)。

## 七、关联

- 现状实现:`internal/engine/ime_windows.go`、`internal/engine/window.go`(SetTextInputRect/Preedit)、`cmd/client/ime.go`。
- 路线图:`doc/go版完善路线图.md` §6.5(IME 条目)、§八 backlog。
- 修复背景:`.mimocode/plans/1787177022013-mighty-nebula.md`(ScenePlay 输入系统审查与修复)。
