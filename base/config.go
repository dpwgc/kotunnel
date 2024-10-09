package base

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"time"
)

type ConfigOptions struct {
	App AppOptions `yaml:"app" json:"app"`
}

type AppOptions struct {
	Mode    string          `yaml:"mode" json:"mode"`
	Secret  string          `yaml:"secret" json:"secret"`
	Servers []ServerOptions `yaml:"servers" json:"servers"`
	Clients []ClientOptions `yaml:"clients" json:"clients"`
	Log     LogOptions      `yaml:"log" json:"log"`
}

type ServerOptions struct {
	OpenPort   int `yaml:"open-port" json:"openPort"`
	TunnelPort int `yaml:"tunnel-port" json:"tunnelPort"`
}

type ClientOptions struct {
	TunnelAddr string `yaml:"tunnel-addr" json:"tunnelAddr"`
	LocalPort  int    `yaml:"local-port" json:"localPort"`
	IdleNum    int    `yaml:"idle-num" json:"idleNum"`
}

type LogOptions struct {
	Path    string `yaml:"path" json:"path"`
	Size    int    `yaml:"size" json:"size"`
	Age     int    `yaml:"age" json:"age"`
	Backups int    `yaml:"backups" json:"backups"`
}

var config ConfigOptions

func Config() ConfigOptions {
	return config
}

func InitConfig(args []string) {

	// ./main server {secret} {open-port} {tunnel-port}
	// ./main client {secret} {tunnel-addr} {local-port} {idle-num}
	if len(args) >= 5 {
		opts := AppOptions{
			Mode:   args[1],
			Secret: args[2],
			Log: LogOptions{
				Path:    "./logs",
				Size:    1,
				Age:     7,
				Backups: 1000,
			},
		}
		if opts.Mode == "server" {
			open, _ := strconv.Atoi(args[3])
			tunnel, _ := strconv.Atoi(args[4])
			opts.Servers = []ServerOptions{{
				OpenPort:   open,
				TunnelPort: tunnel,
			}}
		} else {
			if len(args) == 5 {
				args[5] = "1"
			}
			local, _ := strconv.Atoi(args[4])
			idle, _ := strconv.Atoi(args[5])
			opts.Clients = []ClientOptions{{
				TunnelAddr: args[3],
				LocalPort:  local,
				IdleNum:    idle,
			}}
		}
		config.App = opts
		return
	}

	//加载客户端配置
	configBytes, err := os.ReadFile("./config.yaml")
	if err != nil {
		Println(31, 40, fmt.Sprintf("read config error: %s", err.Error()))
		time.Sleep(5 * time.Second)
		panic(err)
	}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		Println(31, 40, fmt.Sprintf("parse config error: %s", err.Error()))
		time.Sleep(5 * time.Second)
		panic(err)
	}
}
