#### 传奇服务端 go 语言实现
参照 C# Mir2 (https://github.com/Suprcode/mir2)

#### 开发进度
- [x] 注册/登陆/创建角色/进游戏
- [x] 地图/怪物/NPC加载
- [x] 角色移动
- [ ] 玩家退出时保存角色信息
- [ ] 定时保存游戏数据
- [x] 玩家背包/物品掉落/拾取
- [x] 玩家属性(升级/基础属性/装备属性计算)
- [ ] 玩家/怪物状态(Buff/Poison)
- [ ] 技能
- [ ] NPC 交互
- [ ] 怪物 AI
- [ ] 任务

#### BUG
- 别的玩家装备显示不正确

#### 计划
- [ ] CPU 使用率优化
- [ ] 客户端汉化（打算直接用 https://github.com/cjlaaa/mir2)
- [ ] WEB 管理后台
- [ ] 数据库换成 MySQL
- [ ] 刺客/弓箭手

#### 用到的开源库/工具
- [Cellnet](https://github.com/davyxu/cellnet)
- [GORM](https://github.com/jinzhu/gorm)

#### 编译步骤
环境设置
```bash
go env -w GOPROXY=https://goproxy.cn,direct
export GO111MODULE=off
```
新建项目文件夹，设置为 gopath
```bash
mkdir ~/mir
export GOPATH=~/mir
```
获取项目
```bash
mkdir -p $GOPATH/src/github.com/yenkeia
cd $GOPATH/src/github.com/yenkeia
git clone https://github.com/yenkeia/mirgo.git
```
获取依赖
```bash
go get -u -v github.com/mattn/go-sqlite3
go get -u -v github.com/jinzhu/gorm
go get -u -v github.com/davyxu/cellnet
go get -u -v github.com/davyxu/golog
go get -u -v github.com/davyxu/goobjfmt
go get -u -v github.com/davyxu/protoplus
```
解压服务端资源到 mirgo/dotnettools/database 目录下
```bash
cd mirgo/dotnettools
mkdir database
cd database
unrar x ../Daneo1989_Server.rar -pLOMCN -y
```
运行
```bash
cd $GOPATH/src/github.com/yenkeia/mirgo
go run server/main.go
```

#### 联系方式
欢迎 issue 交流
QQ 群: 32309474

#### 客户端
代码: https://github.com/yenkeia/mir2
链接: https://pan.baidu.com/s/1ELI8pO278v9JRyt6lS-A8Q
提取码: 0nc3

#### 截图
![image](https://github.com/yenkeia/mirgo/blob/master/img/img1.png)

#### 参考资料
- [Zinx应用-MMO游戏案例-(2)AOI兴趣点算法](https://www.jianshu.com/p/e5b5db9fa6fe)
- [Zinx应用-MMO游戏案例-(8)移动位置与AOI广播](https://www.jianshu.com/p/8c8fafdace14)
