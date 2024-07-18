package main

import (
	"encoding/json"
	"fmt"
	"kotunnel/base"
	"kotunnel/cli"
	"kotunnel/ser"
	"os"
)

func main() {

	// 配置加载
	base.InitConfig(os.Args)
	marshal, _ := json.Marshal(base.Config().App)
	base.Println(33, 40, "config: "+string(marshal))

	// 日志加载
	base.InitLog()

	base.Println(36, 40, fmt.Sprintf("mode: %s / protocol: %s", base.Config().App.Mode, base.Config().App.Protocol))

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
