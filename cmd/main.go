package main

import (
	"encoding/json"
	"fmt"
	"kotunnel/base"
	"kotunnel/core"
	"os"
	"time"
)

func main() {

	// 配置加载
	base.InitConfig(os.Args)
	// 日志加载
	base.InitLog()

	if len(base.Config().App.Secret) <= 0 {
		base.Println(31, 40, "secret cannot be empty")
		time.Sleep(5 * time.Second)
		return
	}

	// 服务端 or 客户端
	if base.Config().App.Mode == "server" {
		server(base.Config().App)
	} else if base.Config().App.Mode == "client" {
		client(base.Config().App)
	} else {
		base.Println(31, 40, "mode must be 'server' or 'client'")
		time.Sleep(5 * time.Second)
		return
	}
}

func server(opts base.AppOptions) {

	var servers []*core.Server

	for _, v := range opts.Servers {
		bytes, _ := json.Marshal(v)
		base.Println(36, 40, fmt.Sprintf("server start: %s", string(bytes)))
		servers = append(servers, core.NewServer(v.OpenPort, v.TunnelPort, opts.Secret))
	}

	if len(servers) <= 0 {
		base.Println(31, 40, "no server instances")
		time.Sleep(5 * time.Second)
		return
	}

	for i := 0; i < len(servers)-1; i++ {
		go servers[i].Run()
	}

	servers[len(servers)-1].Run()
}

func client(opts base.AppOptions) {

	var clients []*core.Client

	for _, v := range opts.Clients {
		bytes, _ := json.Marshal(v)
		base.Println(36, 40, fmt.Sprintf("client start: %s", string(bytes)))
		for i := 0; i < v.IdleNum; i++ {
			clients = append(clients, core.NewClient(v.TunnelAddr, v.LocalPort, opts.Secret))
		}
	}

	if len(clients) <= 0 {
		base.Println(31, 40, "no client instances")
		time.Sleep(5 * time.Second)
		return
	}

	for i := 0; i < len(clients)-1; i++ {
		go clients[i].Run()
	}

	clients[len(clients)-1].Run()
}
