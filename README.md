# KoTunnel

### 一个基于Go开发的内网穿透工具

***

## 使用方式

* cd 到 cmd 或 gui 目录
* 执行 go build main.go

***

### 服务端配置与启动

#### `程序入口（cmd目录下）` /kotunnel/cmd/main.go

* 修改config.yaml配置（mode设为server）

```yaml
app:
  mode: server
  # 网络协议（当前只支持tcp）
  protocol: tcp
  # 连接密钥（客户端和服务端密钥需一致）
  secret: default
  # 可配置多条隧道
  servers:
      # 对外暴露的访问端口（对外暴露访问）
    - open-port: 9090
      # 隧道连接端口（用来与客户端建立连接）
      tunnel-port: 8080
  # 日志存储配置（日志文件保存于./logs目录下，单文件最大1mb，保存7天内的文件，最多存储1000个日志文件）
  log:
    path: ./logs
    size: 1
    age: 7
    backups: 1000
```

* 配置好后，运行main程序，启动服务端

***

### 客户端（桌面端模式）配置与启动

#### `程序入口（gui目录下）` /kotunnel/gui/main.go

* 修改config.yaml配置（桌面端运行时只需要填写日志存储配置即可）

```yaml
app:
  log:
    path: ./logs
    size: 1
    age: 7
    backups: 1000
```

* 配置好后，运行main程序，启动桌面端

![GUI](gui1.png "")

***

### 客户端（终端模式）配置与启动

#### `程序入口（cmd目录下）` /kotunnel/cmd/main.go

* 修改config.yaml配置（mode设为client）

```yaml
app:
  mode: client
  protocol: tcp
  secret: default
  clients:
      # 服务端隧道连接地址（对应服务端配置里的tunnel-port）
    - tunnel-addr: 0.0.0.0:8080
      # 需要对外暴露的本地端口号
      local-port: 7070
      # 最大空闲连接数（非高并发场景设为1即可）
      idle-num: 5
  log:
    path: ./logs
    size: 1
    age: 7
    backups: 1000
```

* 配置好后，运行main程序，启动客户端