# KoTunnel

### 一个基于Go开发的内网穿透工具

***

## 使用方式

#### 程序入口 `/kotunnel/cmd/main.go`

#### cd 到 cmd 目录下，执行 go build main.go 进行打包

***

### 服务端配置与启动

* 修改 config.yaml 配置，将 mode 设为 server

* 配置好后，运行 main 程序，启动服务端

```
./main
```

***

### 客户端配置与启动

* 修改 config.yaml 配置，将 mode 设为 client

* 配置好后，运行 main 程序，启动客户端

```
./main
```

***

### 另一种启动方式：在运行 main 程序时传入配置参数（参数名不区分大小写）

* 服务端运行命令

```
./main -server secret=123456 openPort=8080 tunnelPort=9090 maxConn=1000
```

* 客户端运行命令

```
./main -client secret=123456 tunnelAddr=0.0.0.0:9090 localPort=7070 idleConn=1
```

#### 使用此方式运行程序，config.yaml 配置只需要填写 log 相关配置，多余的客户端与服务端配置将被忽略

```yaml
log:
  path: ./logs
  size: 1
  age: 7
  backups: 1000
```