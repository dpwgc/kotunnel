package base

import (
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
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
	MaxConn    int `yaml:"max-conn" json:"maxConn"`
}

type ClientOptions struct {
	TunnelAddr string `yaml:"tunnel-addr" json:"tunnelAddr"`
	LocalPort  int    `yaml:"local-port" json:"localPort"`
	IdleConn   int    `yaml:"idle-conn" json:"idleConn"`
}

type LogOptions struct {
	Path    string `yaml:"path" json:"path"`
	Size    int    `yaml:"size" json:"size"`
	Age     int    `yaml:"age" json:"age"`
	Backups int    `yaml:"backups" json:"backups"`
}

func GetConfig(args []string) (*ConfigOptions, error) {

	config := &ConfigOptions{}

	// ./main server {secret} {open-port} {tunnel-port} {max-conn}
	// ./main client {secret} {tunnel-addr} {local-port} {idle-conn}
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
			if len(args) == 5 {
				args[5] = "1000"
			}
			open, _ := strconv.Atoi(args[3])
			tunnel, _ := strconv.Atoi(args[4])
			maxC, _ := strconv.Atoi(args[5])
			opts.Servers = []ServerOptions{{
				OpenPort:   open,
				TunnelPort: tunnel,
				MaxConn:    maxC,
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
				IdleConn:   idle,
			}}
		}
		config.App = opts
		return config, nil
	}

	//加载客户端配置
	configBytes, err := os.ReadFile("./config.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
