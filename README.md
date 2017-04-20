# script-web

云霁科技推出的脚本管理工具，用来管理企业内部的脚本。

具有以下功能：

- 用户
- 脚本版本
- 目录和标签支持
- 支持脚本审核
- 支持导入现有仓库


## 安装


go get -u github.com/cloudcli/script-web

## 运行

首先到script-web目录下

```bash
$ cd $GOPATH/src/github.com/cloudcli/script-web
```

修改app.ini下的mysql连接，仓库所在目录


```bash
$ make seed
$ make build
```

本机上运行

``` bash
$ ./script-web web start --port=3000
```

或
```bash
$ make run
```


如果需要发布到远程服务器，使用
```bash
$ make pack
```

进行打包后，整体复制到服务器上，然后在服务器上修改app.ini文件

