# KoTunnel

### 一个基于Go开发的内网穿透工具

***

## 使用方式

#### 程序入口 `/kotunnel/cmd/main.go`

#### cd 到 cmd 目录下，执行 go build main.go 进行打包

***

### 服务端配置

* 修改 config.yaml 配置，在 servers 数组中添加服务端配置
  * open-port: 需要对外暴露访问的端口号，用于接收外部请求（必填）
  * tunnel-port: 隧道端口号，用于与客户端建立隧道连接（必填）
  * max-conn: tunnel-port 的最大连接数，默认1000（非必填）
  * secret: 鉴权密钥，为空则不开启鉴权功能，默认不开启（非必填）
  * name: 服务端名称，仅作为标识使用，无代码逻辑（非必填）

```yaml
open-port: 9090
tunnel-port: 8080
max-conn: 100
secret: 123456
name: test1
```

***

### 客户端配置

* 修改 config.yaml 配置，在 clients 数组中添加服务端配置
  * tunnel-addr: 服务器隧道地址（必填）
  * local-port: 需要映射到外网的本地端口号（必填）
  * idle-conn: 最大空闲隧道连接数，默认1（非必填）
  * retry-interval: 重试间隔时间，单位秒，发生连接异常时默认5秒后重试（非必填）
  * secret: 鉴权密钥，需与服务端配置保持一致（非必填）
  * name: 客户端名称，仅作为标识使用，无代码逻辑（非必填）

```yaml
tunnel-addr: 0.0.0.0:8080
local-port: 7070
idle-conn: 3
retry-interval: 10
secret: 123456
name: test1
```

***

* 配置好后，运行 main 程序（可在一个程序里同时配置服务端和客户端，一起启动）

```
./main
```

### 快捷启动方式：在运行 main 程序时传入客户端/服务端配置（参数名不区分大小写）

* 服务端运行命令

```
./main -server name=test1 secret=123456 openPort=8080 tunnelPort=9090 maxConn=1000
```

* 客户端运行命令

```
./main -client name=test1 secret=123456 tunnelAddr=0.0.0.0:9090 localPort=7070 idleConn=1 retryInterval=5
```