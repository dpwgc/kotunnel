package main

import (
	"encoding/json"
	"fmt"
	"kotunnel/base"
	"kotunnel/cli"
	"kotunnel/ser"
	"os"
	"strconv"
)

func main() {

	// 配置加载
	base.InitConfig()
	marshal, _ := json.Marshal(base.Config().App)
	base.Println(33, 40, "config: "+string(marshal))

	// 日志加载
	base.InitLog()

	args := os.Args
	argc := len(os.Args)

	// main tcp server {secret} {open-port} {tunnel-port}
	// main tcp client {secret} {tunnel-addr} {local-port} {idle-num}
	if argc == 6 {
		open, _ := strconv.Atoi(args[4])
		tunnel, _ := strconv.Atoi(args[5])
		opts := base.AppOptions{
			Protocol: args[1],
			Mode:     args[2],
			Secret:   args[3],
			Servers: []base.ServerOptions{{
				OpenPort:   open,
				TunnelPort: tunnel,
			}},
		}
		server(opts)
	}
	if argc == 7 {
		local, _ := strconv.Atoi(args[5])
		idle, _ := strconv.Atoi(args[6])
		opts := base.AppOptions{
			Protocol: args[1],
			Mode:     args[2],
			Secret:   args[3],
			Clients: []base.ClientOptions{{
				TunnelAddr: args[4],
				LocalPort:  local,
				IdleNum:    idle,
			}},
		}
		client(opts)
	}

	// 服务端 or 客户端
	if base.Config().App.Mode == "server" {
		server(base.Config().App)
	} else {
		client(base.Config().App)
	}
}

func server(opts base.AppOptions) {
	if opts.Protocol == "udp" {
		// TODO
	} else {
		base.Println(36, 40, "mode: tcp server")
		for _, v := range opts.Servers {
			bytes, _ := json.Marshal(v)
			base.Println(36, 40, fmt.Sprintf("create tcp server: %s", string(bytes)))
			ser.TCP(v.OpenPort, v.TunnelPort)
		}
	}
}

func client(opts base.AppOptions) {
	if opts.Protocol == "udp" {
		// TODO
	} else {
		base.Println(36, 40, "mode: tcp client")
		for _, v := range opts.Clients {
			bytes, _ := json.Marshal(v)
			base.Println(36, 40, fmt.Sprintf("create tcp client: %s", string(bytes)))
			for i := 0; i < v.IdleNum-1; i++ {
				go cli.TCP(v.TunnelAddr, v.LocalPort)
			}
			cli.TCP(v.TunnelAddr, v.LocalPort)
		}
	}
}
