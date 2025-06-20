package main

import (
	"fmt"
	"kotunnel/base"
	"kotunnel/core"
	"os"
)

func main() {
	config := base.GetConfig(os.Args)
	base.InitLog(config.Log)
	runServers(config.Servers)
	runClients(config.Clients)
	<-make(chan byte)
}

func runServers(servers []base.ServerConfig) {
	for _, v := range servers {
		go func(v base.ServerConfig) {
			oErr, tErr := core.NewServer(v.Name, v.OpenPort, v.TunnelPort, v.MaxConn, v.Secret).Run()
			if oErr != nil {
				base.Error(fmt.Sprintf("{server:%s} open port [%v] listen error: %s", v.Name, v.OpenPort, oErr.Error()))
			}
			if tErr != nil {
				base.Error(fmt.Sprintf("{server:%s} tunnel port [%v] listen error: %s", v.Name, v.TunnelPort, tErr.Error()))
			}
		}(v)
	}
}

func runClients(clients []base.ClientConfig) {
	for _, v := range clients {
		if v.IdleConn <= 0 {
			v.IdleConn = 1
		}
		for i := 0; i < v.IdleConn; i++ {
			go func(v base.ClientConfig) {
				core.NewClient(v.Name, v.TunnelAddr, v.LocalPort, v.RetryInterval, v.Secret).Run()
			}(v)
		}
	}
}
