#### 传奇服务端 go 语言实现
参照 C# Mir2 (https://github.com/Suprcode/mir2)

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
编译
```bash
cd $GOPATH/src/github.com/yenkeia/mirgo/server
go build -o smain
```
运行
```bash
./smain
```

#### 客户端
代码: https://github.com/yenkeia/mir2
链接: https://pan.baidu.com/s/1ELI8pO278v9JRyt6lS-A8Q
提取码: 0nc3

#### 截图
![image](https://github.com/yenkeia/mirgo/blob/master/img/img1.png)

#### 参考资料
- [Zinx应用-MMO游戏案例-(2)AOI兴趣点算法](https://www.jianshu.com/p/e5b5db9fa6fe)
- [Zinx应用-MMO游戏案例-(8)移动位置与AOI广播](https://www.jianshu.com/p/8c8fafdace14)
