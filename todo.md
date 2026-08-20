# client

- [ ] 游戏声音
- [x] 攻击特效 (攻杀剑术) 的动画播放问题
- [x] WIL 解析问题, DnItems.wil 解析错误
- [x] 人物男女角色渲染问题
- [x] NPC 不能点
- [x] 支持更大更多的分辨率 ALT + 回车 切换, 或者 res 1/2/3/4 切换
- [x] 小地图显示错误, 怪物红点和玩家白点不停闪烁 (不是 bug：Delphi 原版每 300ms 翻转 m_boViewBlink（PlayScn.pas:801-819），Go 是忠实移植。决定改为常亮)
- [x] 被攻击后的动画朝向问题
- [x] 走路卡脚问题 (跑步会卡 1 秒左右?)
- [x] 从 切换场景：PlayGame → SelectChr, 再从选角场景无法开始游戏
- [ ] SceneSelectChr 场景的 UI 微调, 新建角色的名字输入框
- [ ] ScenePlay 场景的 UI 问题: NPC 对话框
- [x] ScenePlay 场景的 UI 问题: 聊天框 (根因: UI 焦点吞事件——点过 HUD 后键盘全被无处理器控件吞掉; 已修复路由语义/模态焦点泄漏/状态残留, 见 .mimocode/plans/1787177022013-mighty-nebula.md)
- [ ] ScenePlay 场景的 UI 问题: 右下角四个按钮
- [x] 进商店再出商店后, 场景渲染出错
- [ ] 鼠标/键盘输入与 Delphi 版一致, 例如攻击一个怪物后, 取消攻击, 继续跑步
- [ ] 物品背包 UI 按钮调整
- [ ] NPC 购买物品, 弹出的对话框 UI 问题
- [ ] 和 NPC 对话后购买物品功能
- [x] 输入系统 (键盘路由语义对齐 Delphi: release/repeat 语义、模态恒吞、焦点释放; Windows IMM32 候选窗/组合窗定位+预编辑显示; Linux 提交字符经 XIM/text-input, 预编辑受 GLFW 限制, 见路线图 6.5)
- [x] 聊天框字体（GDI 光栅化复刻原版宋体点阵；单一字体 9pt≡Font.Height=-12，聊天行统一字号）
- [x] bug: 女性角色创建人物后, 进入游戏场景, 渲染的是男性角色的身体
- [x] bug: 从游戏场景退回到选择角色场景, 默认激活的是第一个角色
- [x] 把物品从背包扔到地上的逻辑有问题, 全部堆叠在一起
- [x] 击杀怪物后掉落的物品位置不对
- [ ] 怪物死亡动画问题: 稻草人死后的动画
- [ ] 怪物死亡后尸体消失过快
- [ ] 怪物死亡后, 鸡死后直接爆鸡肉, 明显不是原版 176 行为. 应该杀鸡后, 要挖肉才得到鸡肉

# server

- [x] 所有服务端的设置都通过读取 serverconfig 目录下的 jsonc 配置文件来设置
- [x] config.go 使用 jsonc 解析器 (支持注释), 而不是 json 解析器
- [x] 登录注册逻辑是否正确
- [x] 玩家账号密码检查, 玩家名字检查, 是否合规

# wilviewer

- [x] 左侧 tree 文件列表高亮显示当前打开的是哪个 wil 文件
- [x] 右下角 Preview 窗口支持鼠标滚轮上下滚能缩放图片大小, 并且按住滚轮可以左右拖动图片, 方便观察
- [ ] 三个列表可以随意拖动大小高宽
- [x] 动画播放模式是否应该去掉?
- [x] wilviewer 目录扁平化

# mapviewer

- [x] mapviewer 目录扁平化

# serverconfig

- [ ] serverconfig 转换出来的配置文件要带上本项目的 go 服务端额外的一些配置, 而不仅仅是 176 服务端原本的配置
- [ ] serverconfig 转换出来的配置文件要带上中文注释

# 综合

- [x] WIL 解码缓存加 LRU 淘汰（当前解码结果永久驻留，内存只增不减）
- [x] 封包拆包逻辑是否和 Delphi 完全一致
- [x] doc 目录文档整理
- [x] 恢复提交 go.sum（当前被 gitignore）
- [ ] 检查 internal\protocol\message.go 与 Delphi 版是否完全一致, 不要有遗漏
- [ ] cmd/client/testdata 这个目录是什么鬼?
