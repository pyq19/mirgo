# CodeGraph × opencode 用法

CodeGraph 把整个代码库解析成本地 SQLite 知识图谱（符号/调用边/文件，100% 本地、无 API key），通过 MCP 接入 opencode。一次 `codegraph_explore` 调用返回相关符号的**带行号源码 + 调用路径 + 影响范围**，替代 grep + 逐文件 Read 循环。

## 1. 安装 CLI

自带 Node 运行时，无需预装 Node：

```bash
# Windows (PowerShell)
irm https://raw.githubusercontent.com/colbymchenry/codegraph/main/install.ps1 | iex

# macOS / Linux
curl -fsSL https://raw.githubusercontent.com/colbymchenry/codegraph/main/install.sh | sh

# 已有 Node 也可用 npm
npm i -g @colbymchenry/codegraph
```

装完**开新终端**（PATH 才生效），`codegraph --version` 验证。升级用 `codegraph upgrade`（`--check` 只查不装）。

## 2. 接入 opencode

### 方案 A：installer 自动配置（全局）

```bash
codegraph install            # 交互式选择 target/location
codegraph install -t opencode -l global -y   # 非交互直达
```

opencode 是 installer 官方支持的代理，自动完成两件事：

- 往全局配置 `~/.config/opencode/opencode.jsonc` 写入 MCP 服务器（`codegraph serve --mcp`）
- 往全局 `~/.config/opencode/AGENTS.md` 注入 `CODEGRAPH_START/END` 引导块，让**所有已索引仓库**的会话都带上 CodeGraph 使用指引

只想看会写什么而不落盘：`codegraph install --print-config opencode`。

### 方案 B：项目级 MCP 配置（本仓库做法）

在项目根 `opencode.jsonc` 加：

```jsonc
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "codegraph": {
      "type": "local",
      "command": ["codegraph", "serve", "--mcp"],
      "enabled": true
    }
  }
}
```

### 生效与验证

- 配置改动需**重启 opencode** 生效（MCP 服务器随进程启动）
- 生效后工具列表中出现 `codegraph_codegraph_explore`（opencode 的 MCP 工具带 `<服务器名>_` 前缀；模型视角的工具名是 `codegraph_explore`）
- 可用 `tools` 配置全局或按代理开关 MCP 工具（如 `"codegraph_*": false`），一般不需要

## 3. 索引项目

```bash
cd your-project
codegraph init     # 创建 .codegraph/ 并一次性建好全图
codegraph status   # 查看统计：文件数/符号数/边数
```

默认跳过 node_modules/vendor/dist/build 等依赖与构建目录、`.gitignore` 内容、大于 1MB 的文件。需要定制时用项目根的 `codegraph.json`（改完运行 `codegraph index` 重建生效）：

```json
{
  "exclude": ["static/", "**/vendor/**"],
  "include": ["asset/delphi/"],
  "extensions": { ".dota_lua": "lua", ".tpl": "php", ".h": "cpp" }
}
```

| 字段 | 作用 |
|------|------|
| `exclude` | 排除已提交入库但不想进图的目录（gitignore 风格，相对仓库根；`.gitignore` 挡不住已提交的文件） |
| `include` | 纳入被 `.gitignore` 挡掉的真实源码（覆盖 .gitignore；内置跳过如 node_modules、.git 不可恢复） |
| `extensions` | 非标扩展名映射到受支持语言（值必须是支持的语言 id；拼错只警告跳过，不会破坏索引） |

本仓库实例：`{"include": ["asset/delphi"]}` — 把 gitignore 的 Delphi 参考源码拉进图。另：想把默认排除的目录拉回索引（如确实要索引某个 vendor 依赖），在 `.gitignore` 加否定规则 `!vendor/`。

## 4. 同步机制：无需手动 sync

- 文件监视器（OS 原生事件）监听增删改，debounce 后增量同步，编辑后约 1s 内入图
- MCP 服务器（重）连接时先做 catch-up 对账，吸收离线期间的改动（git pull、其他编辑器的修改）
- 同步窗口内若有查询引用了未同步文件，响应开头会带 `⚠️` 横幅指名该文件 → 代理应改用 Read
- `codegraph status` 随时验证（有积压会显示 Pending sync 段）；仅当监视被禁用（沙箱、`--no-watch`、`CODEGRAPH_NO_DAEMON=1`）才需要手动 `codegraph sync`

## 5. 使用规则

- **编辑前先 explore**：一次拿到要改符号的源码 + 谁调用它 + 改动波及，带着爆炸半径动手
- 传自然语言问题或一组符号/文件名，通常一次调用回答整个问题；不够就再 explore 一次，换更具体的名字
- **流程类问题**（"X 如何走到 Y"）把两端符号写进同一次查询，动态分发（回调、接口→实现）的跳转也能追到
- **返回的源码视为已 Read**：不要重新打开这些文件，可直接照此 Edit
- **信任索引**（来自完整 AST 解析），不要用 grep 复核
- 多仓库/monorepo：给工具调用传 `projectPath`（CLI 用 `-p/--path`），可查询任意已索引项目
- Read/Grep 留给索引不覆盖的内容（配置、文档、二进制）或确认个别细节
- 索引不含正确性判断——编译/测试/lint 的职责不变，CodeGraph 只补结构上下文

## 6. opencode 代理与子代理配合

MCP 的使用指引随 `initialize` 下发，只到达**主代理**；opencode 的子代理也**不会**继承 AGENTS.md。但：

- **build/plan 主代理**与 **general 子代理**（全工具权限）继承 MCP 工具 → 直接用 `codegraph_explore`
- **explore/scout 子代理**是只读白名单、没有 MCP 工具，但 bash 可用 → 走 `codegraph explore "..."` CLI

所以做法是：在 AGENTS.md 里写一节给主代理遵循，并要求主代理派发子代理时把指引随 prompt 转发。可直接粘贴的模板——本仓库 AGENTS.md 的 CodeGraph 一节即以此为准，逐字一致，仅多两句项目专属内容（索引覆盖范围、指向本文档的引用）：

```markdown
## CodeGraph 代码索引

本项目已建 CodeGraph 知识图谱索引（.codegraph/），文件监视器自动同步，写入后约 1s 延迟。

**任何代理查找代码、理解结构、评估改动影响、获取编辑前上下文时，优先用 CodeGraph，替代 grep + 逐文件读取。**

### build/plan 主代理与 general 子代理（有 codegraph MCP 工具）

核心工具 `codegraph_explore`（工具列表里显示为 `codegraph_codegraph_explore`）：

- "X 怎么工作 / X 在哪 / 架构如何"、编辑前取上下文 → 传自然语言问题或一组符号/文件名；一次调用返回相关符号的带行号源码（按文件分组）+ 调用路径 + 影响范围
- 返回的源码视为已 Read：不要重新打开这些文件，可直接照此 Edit
- 流程类问题（"X 如何走到 Y"）把两端符号写进同一次查询
- 信任索引结果（来自完整 AST 解析），不要用 grep 复核
- Read/Grep 只用于索引不覆盖的内容（配置、文档、二进制资源）或确认个别细节
- 响应开头出现 "⚠️ Some files referenced below were edited since the last index sync…" 时，对列出的文件改用 Read；出现 "⚠️ CodeGraph auto-sync is DISABLED…" 时先运行 `codegraph sync`

### explore/scout 子代理（无 MCP 工具）

用 bash 调 `codegraph` CLI（只读）：

- `codegraph explore "问题或符号名"` — 主命令，等价 codegraph_explore
- `codegraph query <符号>` — 按名找位置；`codegraph callers|callees <符号>` — 调用方/被调用方
- `codegraph impact <符号>` — 改动波及；`codegraph node <符号>` — 单符号源码+调用线索

### 派发子代理时

子代理不会继承本 AGENTS.md。主代理派发涉及代码的任务时，prompt 中必须附上：
"查找/理解代码优先用 CodeGraph：general 子代理用 codegraph_explore MCP 工具；explore/scout 子代理用 bash 运行 codegraph explore \"...\"。返回的源码视为已读，不要重复打开文件。"
```

套用到自己的项目时，在开头补一句索引覆盖什么、不覆盖什么（本仓库的写法："覆盖全部 Go 源码与 Delphi 参考源码；vendor/、serverconfig/、文档不在索引内"）。

## 7. CLI 速查（终端 / 无 MCP 工具的子代理）

```bash
codegraph explore "问题或符号名"      # = codegraph_explore，主命令
codegraph query <符号>               # 按名找符号位置（--kind --limit --json）
codegraph node <符号|文件>            # 单符号源码+调用线索 / 按行读文件
codegraph callers <符号>             # 调用方
codegraph callees <符号>             # 被调用方
codegraph impact <符号>              # 改动波及（--depth）
codegraph files [路径]               # 文件树（--format tree|flat|grouped）
codegraph affected [files...]        # 源文件改动影响哪些测试文件（可 --stdin 接 git diff）
codegraph status                     # 索引状态
codegraph sync / index               # 增量同步 / 全量重建
codegraph daemon                     # 管理后台守护进程（停掉卡住的）
codegraph unlock                     # 清理卡住索引的过期锁
```

## 8. 卸载

```bash
codegraph uninstall      # 从所有已配置代理移除配置并卸载 CLI（--keep-cli 只移除配置）
codegraph uninit         # 删除当前项目的 .codegraph/ 索引
```
