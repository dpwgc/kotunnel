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
	marshal, _ := json.Marshal(base.Config().Application)
	base.Println(33, 40, "配置加载完毕: "+string(marshal))
	base.InitLog()
	if base.Config().Application.Mode == "server" {
		if base.Config().Application.Protocol == "udp" {
			// TODO
		} else {
			base.Println(32, 40, "以TCP服务端模式启动")
			ser.TCP(base.Config().Application.Server.ListenPort, base.Config().Application.Server.ClientPort)
		}
	} else {
		if base.Config().Application.Protocol == "udp" {
			// TODO
		} else {
			base.Println(32, 40, "以TCP客户端模式启动")
			cli.TCP(base.Config().Application.Client.RemoteAddr, base.Config().Application.Client.LocalPort, base.Config().Application.Client.TunnelNum)
		}
	}
}
