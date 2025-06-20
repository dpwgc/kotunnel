package main

import (
	"bufio"
	"fmt"
	"kotunnel/base"
	"kotunnel/core"
	"os"
	"strings"
	"sync"
)

var nidLast int64 = 0
var nidLock sync.Mutex

func NID() string {
	nidLock.Lock()
	defer nidLock.Unlock()
	nidLast = nidLast + 1
	return fmt.Sprintf("%v", nidLast)
}

var nodes sync.Map

func main() {
	config := base.GetConfig(os.Args)
	base.InitLog(config.Log)
	runServers(config.Servers)
	runClients(config.Clients)
	// 命令行指令接收
	reader := bufio.NewReader(os.Stdin)
	for {
		in, _ := reader.ReadString('\n')
		// convert CRLF to LF
		in = strings.Replace(in, "\n", "", -1)

		if strings.HasPrefix(in, "ls") {
			nodes.Range(func(k, v interface{}) bool {
				node, ok := v.(Node)
				if ok {
					if node.Server != nil {
						fmt.Printf("> id=%s mode=server name=%s tunnelPort=%v openPort=%v maxConn=%v\n", k.(string), node.ServerConfig.Name, node.ServerConfig.TunnelPort, node.ServerConfig.OpenPort, node.ServerConfig.MaxConn)
					}
					if node.Client != nil {
						fmt.Printf("> id=%s mode=client name=%s localPort=%v tunnelAddr=%s idleConn=%v retryInterval=%v\n", k.(string), node.ClientConfig.Name, node.ClientConfig.LocalPort, node.ClientConfig.TunnelAddr, node.ClientConfig.IdleConn, node.ClientConfig.RetryInterval)
					}
				}
				return true
			})
		}
		if strings.HasPrefix(in, "kill") {
			in = strings.ReplaceAll(in, "kill", "")
			id := strings.ReplaceAll(in, " ", "")
			v, ok := nodes.Load(id)
			if ok {
				node, ok := v.(Node)
				if ok {
					if node.Server != nil {
						node.Server.Close()
						nodes.Delete(id)
					}
					if node.Client != nil {
						node.Client.Close()
						nodes.Delete(id)
					}
					fmt.Println("> success")
				} else {
					fmt.Println("> fail (node not found)")
				}
			}
		}
		if strings.HasPrefix(in, "add") {
			mode, params := base.ArgsToParams(strings.Split(strings.ReplaceAll(in, "add", ""), " "))
			if mode == "server" {
				fmt.Println("> add server")
				runServer(base.ServerConfig{
					Name:       params["name"],
					Secret:     params["secret"],
					OpenPort:   base.ToInt(params["openport"]),
					TunnelPort: base.ToInt(params["tunnelport"]),
					MaxConn:    base.ToInt(params["maxconn"]),
				})
			} else if mode == "client" {
				fmt.Println("> add client")
				runClient(base.ClientConfig{
					Name:          params["name"],
					Secret:        params["secret"],
					TunnelAddr:    params["tunneladdr"],
					LocalPort:     base.ToInt(params["localport"]),
					IdleConn:      base.ToInt(params["idleconn"]),
					RetryInterval: base.ToInt(params["retryinterval"]),
				})
			}
		}
		if in == "shutdown" {
			fmt.Println("> shutdown")
			return
		}
	}
}

func runServers(servers []base.ServerConfig) {
	for _, v := range servers {
		runServer(v)
	}
}

func runServer(v base.ServerConfig) {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				base.Error(fmt.Sprintf("{server:%s} recover error: %s", v.Name, err))
			}
		}()
		server := core.NewServer(v.Name, v.OpenPort, v.TunnelPort, v.MaxConn, v.Secret)
		nid := NID()
		nodes.Store(nid, Node{
			Server:       server,
			ServerConfig: v,
		})
		defer nodes.Delete(nid)
		oErr, tErr := server.Run()
		if oErr != nil {
			base.Error(fmt.Sprintf("{server:%s} open port [%v] listen error: %s", v.Name, v.OpenPort, oErr.Error()))
		}
		if tErr != nil {
			base.Error(fmt.Sprintf("{server:%s} tunnel port [%v] listen error: %s", v.Name, v.TunnelPort, tErr.Error()))
		}
	}()
}

func runClients(clients []base.ClientConfig) {
	for _, v := range clients {
		if v.IdleConn <= 0 {
			v.IdleConn = 1
		}
		for i := 0; i < v.IdleConn; i++ {
			runClient(v)
		}
	}
}

func runClient(v base.ClientConfig) {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				base.Error(fmt.Sprintf("{client:%s} recover error: %s", v.Name, err))
			}
		}()
		client := core.NewClient(v.Name, v.TunnelAddr, v.LocalPort, v.RetryInterval, v.Secret)
		nid := NID()
		nodes.Store(nid, Node{
			Client:       client,
			ClientConfig: v,
		})
		defer nodes.Delete(nid)
		client.Run()
	}()
}

type Node struct {
	ServerConfig base.ServerConfig
	ClientConfig base.ClientConfig
	Server       *core.Server
	Client       *core.Client
}
