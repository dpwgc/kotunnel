package main

import (
	"encoding/json"
	"fmt"
	"kotunnel/base"
	"kotunnel/cli"
	"kotunnel/ser"
)

func main() {

	// 配置加载
	base.InitConfig()
	marshal, _ := json.Marshal(base.Config().App)
	base.Println(33, 40, "config: "+string(marshal))

	// 日志加载
	base.InitLog()

	// 服务端 or 客户端
	if base.Config().App.Mode == "server" {
		server()
	} else {
		client()
	}
}

func server() {
	if base.Config().App.Protocol == "udp" {
		// TODO
	} else {
		base.Println(36, 40, "mode: tcp server")
		for _, v := range base.Config().App.Servers {
			bytes, _ := json.Marshal(v)
			base.Println(36, 40, fmt.Sprintf("create tcp server: %s", string(bytes)))
			ser.TCP(v.OpenPort, v.ClientPort)
		}
	}
}

func client() {
	if base.Config().App.Protocol == "udp" {
		// TODO
	} else {
		base.Println(36, 40, "mode: tcp client")
		for _, v := range base.Config().App.Clients {
			bytes, _ := json.Marshal(v)
			base.Println(36, 40, fmt.Sprintf("create tcp client: %s", string(bytes)))
			for i := 0; i < v.TunnelNum-1; i++ {
				go cli.TCP(v.RemoteAddr, v.LocalPort)
			}
			cli.TCP(v.RemoteAddr, v.LocalPort)
		}
	}
}
