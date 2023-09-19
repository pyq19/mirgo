# 停止更新

没时间更新了

# 传奇服务端 go 语言实现

参照 [Suprcode/mir2](https://github.com/Suprcode/mir2), 用 go 语言**翻译**了该项目的**服务端**部分

**这份代码就是写着自娱自乐的玩具, 功能很不完善, 无法开服赚钱**

# 截图

![image](./assets/img1.jpg)

# 环境搭建

## 前置条件

1. [安装 go 语言](http://docscn.studygolang.com/doc/install)
2. (可选)[设置拉取依赖代理](https://goproxy.io/zh/)
3. 获取服务端资源
   ```bash
   git clone https://gitee.com/pyq19/mir2ServerRelease.git
   ```
4. 获取服务端代码
   ```bash
   git clone https://github.com/pyq19/mirgo.git
   ```

## 编译 (unix)

1. 新建配置文件 `config.toml`
   在当前目录下, 新建 `config.toml`, 文件内容 `DataPath="服务端资源绝对路径(mir2ServerRelease)"`
2. `go mod vendor`
3. `go build -o server ./cmd/server`
4. `./server`

## 编译 (windows)

TODO

# 联系方式

- QQ 群: 32309474

# 客户端

- [客户端代码](https://gitee.com/pyq19/mir2.git)，感谢 https://github.com/cjlaaa/mir2 的汉化
- [客户端资源](https://pan.baidu.com/s/1ELI8pO278v9JRyt6lS-A8Q) 提取码: 0nc3
- 注意，客户端我为了方便开发改了一些代码，如果你不想编译一遍客户端，就用我上面提到的群里面我编译好的 client.7z 解压覆盖到你下载好的客户端资源上，然后改 host 把 mir.impyq.com 指向 127.0.0.1 就行

# 感谢贡献者

- @qcdong2016
- @firma
