# client

- [ ] 游戏声音
- [ ] 反作弊
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
- [ ] ScenePlay 场景的 UI 问题: 右下角四个按钮
- [x] 进商店再出商店后, 场景渲染出错
- [ ] 鼠标/键盘输入与 Delphi 版一致, 例如攻击一个怪物后, 取消攻击, 继续跑步
- [ ] 物品背包 UI 按钮调整
- [ ] NPC 购买物品, 弹出的对话框 UI 问题
- [ ] 和 NPC 对话后购买物品功能
- [ ] 输入系统
- [x] 聊天框字体（GDI 光栅化复刻原版宋体点阵；单一字体 9pt≡Font.Height=-12，聊天行统一字号）
- [x] bug: 女性角色创建人物后, 进入游戏场景, 渲染的是男性角色的身体
- [x] bug: 从游戏场景退回到选择角色场景, 默认激活的是第一个角色
- [x] 把物品从背包扔到地上的逻辑有问题, 全部堆叠在一起
- [x] 击杀怪物后掉落的物品位置不对
- [x] 怪物死亡动画问题: 稻草人死后的动画 (原版 THuSuABi 死亡火焰特效层: Mon3.wil 2860-2869 随倒下动画逐帧加法混合叠加, 已按数据化特效框架实现)
- [x] 怪物死亡后尸体消失过快 (SearchViewRange 不再排除死亡对象, 尸体保留 3 分钟至 ghost, ghost 广播 RM_DISAPPEAR; 后进入者经 SM_DEATH/SM_SKELETON 建尸体)
- [x] 怪物死亡后, 鸡死后直接爆鸡肉, 明显不是原版 176 行为. 应该杀鸡后, 要挖肉才得到鸡肉 (动物死亡不散落, 掉落物刷怪时 roll 进尸体背包, Alt 挖肉扣皮量耗尽出货+骨架化)
- [x] 特化怪物渲染全量补齐(数据化方案): 死亡特效族18种/攻击特效族20+种/投射物6种/施法/物体类8种/城门城墙全状态机, 数据表见 cmd/client/moneffect.go (对照 Delphi AxeMon.pas+HerbActor.pas)
- [x] bug: 不同分辨率下的字体模糊问题 (字形按物理像素倍率光栅化+缓存物理度量, DrawText 物理像素走笔/整数对齐, 缩放源改 GetFramebufferSize+回调; engine/text.go)
- [x] bug: 血条会挡住人物/怪物, 以及怪物被攻击后, 血条没有改变. 血条的行为应该完全遵守delphi版实现 (血条改为 DrawScrn.pas:277-302 式顶层覆盖层 pass; 血条机制: 心灵启示技能28完整流程+STATE_OPENHEATH 状态位+到期 SMCloseHealth; SMStruck 改 Delphi 格式携带 HP/MaxHP 实时更新; SMHealthSpellChanged Delphi 格式+ShowHP 广播; 移动消息 body 第二槽 featureEx→charStatus, 马状态改由 SMFeatureChanged 补偿; 受击开启血条=Go 专属 jsonc 配置 game.showHPOnStruck 默认开启/10秒无命中自动关, serverconfig 转换写入默认值)
- [x] 拾取物品时, 左上角绿字提示, 遵守delphi版实现 (拾取改发 SMAddItem+TClientItem body, 客户端 AddSysMsg/DrawScreenTop: 左上角(30,40)绿字黑描边/行距16/3000ms/上限10条; "<物品名> 放到包囊."; 金币改绿字 "%d 金币 金币.")
- [x] 击杀怪物后, 增加经验值的消息要在聊天框中提示, 遵守delphi版实现 (SMWinExp 改 Delphi 格式 Recog=总经验/Param+Tag=增量, 聊天框白字红底"已获得 N 点经验值.."; SMLevelUp Recog=Exp/Param=Level, 升级绿字"恭喜！你升级了!"+紫底祝贺语)
- [ ] 聊天框样式细节微调
- [ ] 检查所有 GM 命令
- [ ] 玩家角色死亡后, 聊天框应该有提示, 目前没有

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

- [x] serverconfig 转换出来的配置文件要带上本项目的 go 服务端额外的一些配置, 而不仅仅是 176 服务端原本的配置 (server.jsonc 已输出与 config.go 对齐的全部 15 节, Go 专属参数按默认值写入)

- [x] serverconfig 转换出来的配置文件要带上中文注释, 发挥 jsonc 可以注释的优势 (jsonc_writer.go 反射写出器 + 各文件中文注释表; 数据库大文件头部字段说明表 + 首条记录注释)

# 综合

- [x] WIL 解码缓存加 LRU 淘汰（当前解码结果永久驻留，内存只增不减）
- [x] 封包拆包逻辑是否和 Delphi 完全一致
- [x] doc 目录文档整理
- [x] 恢复提交 go.sum（当前被 gitignore）
- [ ] 检查 internal\protocol\message.go 与 Delphi 版是否完全一致, 不要有遗漏
- [ ] cmd/client/testdata 这个目录是什么鬼?
