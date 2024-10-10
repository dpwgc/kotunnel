# KoTunnel

### 一个基于Go开发的内网穿透工具

***

## 使用方式

#### 程序入口 `/kotunnel/cmd/main.go`

#### cd 到 cmd 目录下，执行 go build main.go 进行打包

***

### 服务端配置与启动

* 修改 config.yaml 配置，将 mode 设为 server

```yaml
app:
  mode: server
  # 连接密钥（不能为空，客户端和服务端密钥需一致）
  secret: default
  # 可配置多条隧道
  servers:
      # 对外暴露的访问端口（对外暴露访问）
    - open-port: 9090
      # 隧道连接端口（用来与客户端建立连接）
      tunnel-port: 8080
  # 日志存储配置（日志文件保存于./logs目录下，单文件最大1mb，保存7天内的文件，最多存储1000个日志文件）
  log:
    # 日志文件存储路径
    path: ./logs
    # 单个日志文件最多几mb
    size: 1
    # 单个日志的最大保留天数
    age: 7
    # 最多保留多少个日志文件
    backups: 1000
```

* 配置好后，运行 main 程序，启动服务端

```
./main
```

***

### 客户端配置与启动

* 修改 config.yaml 配置，将 mode 设为 client

```yaml
app:
  mode: client
  secret: default
  clients:
      # 服务端隧道连接地址（对应服务端配置里的tunnel-port）
    - tunnel-addr: 0.0.0.0:8080
      # 需要对外暴露的本地端口号
      local-port: 7070
      # 最大空闲连接数（非高并发场景设为1即可）
      idle-conn: 1
  log:
    path: ./logs
    size: 1
    age: 7
    backups: 1000
```

* 配置好后，运行 main 程序，启动客户端

```
./main
```

***

### 另一种启动方式：在运行 main 程序时传入配置参数

* 服务端运行命令

```
模版:
./main server {secret} {open-port} {tunnel-port}

样例:
./main server default 9090 8080
```

* 客户端运行命令

```
模版:
./main client {secret} {tunnel-addr} {local-port} {idle-conn}

样例:
./main client default 0.0.0.0:8080 7070 1

末尾的空闲连接数配置可省略，默认是1:
./main client default 0.0.0.0:8080 7070
```

#### 使用此方式运行程序，config.yaml 配置只需要填写 log 相关配置，多余的客户端与服务端配置将被忽略

```yaml
app:
  log:
    path: ./logs
    size: 1
    age: 7
    backups: 1000
```