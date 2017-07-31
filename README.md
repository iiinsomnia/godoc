# godoc

go开发的一个简单易用的文档管理系统，用于API等文档的管理，采用markdown编辑

## 特点

* 采用golang开发，部署简单
* 只有一个执行文件，不依赖任何前端静态资源
* 文档采用 markdown 编辑，程序员的最爱

## 运行

* 去 `release` 下载编译好的执行文件
* 导入 `godoc.sql`
* 配置 `env.ini` 和 `log.xml` 配置文件
* 执行相应操作系统的可执行文件
* 访问：`localhost:8000`

> 如果你是golang开发者，可以下载源码，自行 `go run` 或者 `go build`

### Mac & Linux

```sh
cd /godoc

# mac
./godoc_darwin
# 如需后台运行，则：
./godoc_darwin &

# linux
./godoc_linux_x64
# 如需后台运行，则：
./godoc_linux_x64 &
```

### Windows

双击运行即可

## 说明

* 测试登录用户：admin/admin、demo/123
* 访问端口可以在配置文件中设置，默认：8000
* 新增用户初始密码可以在配置文件中设置，默认：123
* 如需配置域名，建议使用 `nginx` 反向代理
