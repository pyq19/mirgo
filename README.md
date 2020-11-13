传奇服务端 go 语言实现
--------------------
参照[C#传奇](https://github.com/Suprcode/mir2)，实现一个可玩的游戏服务端
这份代码就是写着自娱自乐的玩具，功能很不完善，无法开服赚钱

截图
--------------------
![image](./assets/img1.jpg)

环境搭建
--------------------
  * [安装golang环境](http://docscn.studygolang.com/doc/install)
  * [设置拉取依赖代理](https://goproxy.io/zh/)
  * 获取服务端资源 git clone https://gitee.com/yenkeia/mir2ServerRelease.git
  * 获取服务端代码 git clone https://github.com/yenkeia/mirgo.git
  * 因为我很久以前把一个很大且没用的文件加到代码仓库了导致代码仓库很大，如果你觉得拉代码的速度慢可以把 https://github.com 换成 https://gitee.com 后面地址不变，速度会快很多

编译运行
--------------------
  * 进入目录 mirgo/cmd/server
  * 新建 config.toml, 文件内容 DataPath="服务端资源绝对路径(mir2ServerRelease)"
  * go mod vendor 
  * go build 
  * ./server
  
联系方式
--------------------
  * QQ 群: 32309474

客户端
--------------------
  * [客户端代码](https://gitee.com/yenkeia/mir2.git)，感谢 https://github.com/cjlaaa/mir2 的汉化
  * [客户端资源](https://pan.baidu.com/s/1ELI8pO278v9JRyt6lS-A8Q) 提取码: 0nc3
  * 注意，客户端我为了方便开发改了一些代码，如果你不想编译一遍客户端，就用我上面提到的群里面我编译好的 client.7z 解压覆盖到你下载好的客户端资源上，然后改 host 把 mir.impyq.com 指向 127.0.0.1 就行

感谢贡献者
--------------------
  * @qcdong2016 
  * @firma

参考资料
--------------------
- [LOMCN](https://www.lomcn.org/forum/)
- [知乎 - 行为树](https://www.zhihu.com/search?type=content&q=%E8%A1%8C%E4%B8%BA%E6%A0%91)
- [行为树概念与结构](https://zhuanlan.zhihu.com/p/92298402)
- [行为树（Behavior Tree）实践（1）– 基本概念](http://www.aisharing.com/archives/90)
