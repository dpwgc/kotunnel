package main

import (
	"fmt"
	"kotunnel/base"
	"kotunnel/core"
	"os"
	"sync"
)

func main() {

	// 配置加载
	config, err := base.GetConfig(os.Args)
	if err != nil {
		base.Tips(base.Red, fmt.Sprintf("load config error: %s", err.Error()))
		return
	}
	// 日志加载
	base.InitLog(config.Log)

	// 服务端 or 客户端
	if config.Mode == base.Server {
		server(config)
	} else if config.Mode == base.Client {
		client(config)
	} else {
		base.Tips(base.Red, "mode must be 'server' or 'client'")
	}
}

func server(opts *base.ConfigOptions) {

	wait := sync.WaitGroup{}

	var servers []*core.Server

	for _, v := range opts.Servers {
		base.Tips(base.Blue, fmt.Sprintf("server [%v] -> [%v] launch", v.TunnelPort, v.OpenPort))
		servers = append(servers, core.NewServer(v.OpenPort, v.TunnelPort, v.MaxConn, opts.Secret))
	}

	if len(servers) <= 0 {
		base.Tips(base.Red, "no server instances")
	}

	wait.Add(len(servers))
	for _, v := range servers {
		go func(v *core.Server) {
			defer wait.Done()
			oErr, tErr := v.Run()
			if oErr != nil {
				base.Tips(base.Red, fmt.Sprintf("open port [%v] listen error: %s", v.OpenPort(), oErr.Error()))
			}
			if tErr != nil {
				base.Tips(base.Red, fmt.Sprintf("tunnel port [%v] listen error: %s", v.TunnelPort(), tErr.Error()))
			}
		}(v)
	}

	wait.Wait()
}

func client(opts *base.ConfigOptions) {

	wait := sync.WaitGroup{}

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
		base.Tips(base.Red, "no client instances")
	}

	wait.Add(len(clients))
	for _, v := range clients {
		go func(v *core.Client) {
			defer wait.Done()
			v.Run()
		}(v)
	}

	wait.Wait()
}
