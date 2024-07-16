package main

import (
	"encoding/json"
	"fmt"
	"kotunnel/base"
	"kotunnel/cli"
	"kotunnel/ser"
)

func main() {
	fmt.Println("  _  __    _______                          _ \n | |/ /   |__   __|                        | |\n | ' /  ___  | | _   _  _ __   _ __    ___ | |\n |  <  / _ \\ | || | | || '_ \\ | '_ \\  / _ \\| |\n | . \\| (_) || || |_| || | | || | | ||  __/| |\n |_|\\_\\\\___/ |_| \\__,_||_| |_||_| |_| \\___||_|\n                                              ")
	base.InitConfig()
	marshal, _ := json.Marshal(base.Config().App)
	base.Println(33, 40, "config: "+string(marshal))
	base.InitLog()
	if base.Config().App.Mode == "server" {
		if base.Config().App.Protocol == "udp" {
			// TODO
		} else {
			base.Println(36, 40, "mode: tcp server")
			for _, v := range base.Config().App.Servers {
				bytes, _ := json.Marshal(v)
				base.Println(32, 40, fmt.Sprintf("create tcp server: %s", string(bytes)))
				ser.TCP(v.OpenPort, v.ClientPort)
			}
		}
	} else {
		if base.Config().App.Protocol == "udp" {
			// TODO
		} else {
			base.Println(36, 40, "mode: tcp client")
			for _, v := range base.Config().App.Clients {
				bytes, _ := json.Marshal(v)
				base.Println(32, 40, fmt.Sprintf("create tcp client: %s", string(bytes)))
				cli.TCP(v.RemoteAddr, v.LocalPort, v.TunnelNum)
			}
		}
	}
}
