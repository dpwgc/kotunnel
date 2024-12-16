package main

import (
	"fmt"
	"kotunnel/base"
	"kotunnel/core"
	"os"
)

const (
	Server = "server"
	Client = "client"
)

func main() {

	appCh := make(chan bool)

	// 配置加载
	config, err := base.GetConfig(os.Args)
	if err != nil {
		base.Tips(base.Red, fmt.Sprintf("load config error: %s", err.Error()), 5)
		return
	}
	// 日志加载
	base.InitLog(config.App.Log)

	// 服务端 or 客户端
	if config.App.Mode == Server {
		server(config.App)
	} else if config.App.Mode == Client {
		client(config.App)
	} else {
		base.Tips(base.Red, "mode must be 'server' or 'client'", 5)
		close(appCh)
		return
	}

	<-appCh
}

func server(opts base.AppOptions) {

	var servers []*core.Server

	for _, v := range opts.Servers {
		base.Tips(base.Blue, fmt.Sprintf("server [%v] -> [%v] launch", v.TunnelPort, v.OpenPort))
		servers = append(servers, core.NewServer(v.OpenPort, v.TunnelPort, v.MaxConn, opts.Secret))
	}

	if len(servers) <= 0 {
		base.Tips(base.Red, "no server instances", 5)
	}

	for _, v := range servers {
		go v.Run()
	}
}

func client(opts base.AppOptions) {

	var clients []*core.Client

	for _, v := range opts.Clients {
		if v.IdleConn <= 0 {
			v.IdleConn = 1
		}
		base.Tips(base.Blue, fmt.Sprintf("client [%v] -> [%s] launch", v.LocalPort, v.TunnelAddr))
		for i := 0; i < v.IdleConn; i++ {
			clients = append(clients, core.NewClient(v.TunnelAddr, v.LocalPort, opts.Secret))
		}
	}

	if len(clients) <= 0 {
		base.Tips(base.Red, "no client instances", 5)
	}

	for _, v := range clients {
		go v.Run()
	}
}
