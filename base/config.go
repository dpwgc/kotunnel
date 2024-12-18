package base

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

const (
	Server = "server"
	Client = "client"
)

type ConfigOptions struct {
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

	// ./main -server secret=123456 openPort=8080 tunnelPort=9090 maxConn=1000
	// ./main -client secret=123456 tunnelAddr=0.0.0.0:9090 localPort=7070 idleConn=1

	var params = make(map[string]string)
	first := true
	for _, v := range args {
		if first {
			first = false
			continue
		}
		if strings.HasPrefix(v, "-") {
			params["mode"] = strings.ReplaceAll(v, "-", "")
			continue
		}
		arr := strings.Split(v, "=")
		if len(arr) > 1 {
			params[strings.ToLower(arr[0])] = strings.Join(arr[1:], "")
		}
	}

	if len(params) > 0 {

		config := &ConfigOptions{
			Mode:   params["mode"],
			Secret: params["secret"],
			Log: LogOptions{
				Path:    "./logs",
				Size:    1,
				Age:     7,
				Backups: 1000,
			},
		}
		if config.Mode == Server {
			open, _ := strconv.Atoi(params["openport"])
			tunnel, _ := strconv.Atoi(params["tunnelport"])
			maxC, _ := strconv.Atoi(params["maxconn"])
			if maxC <= 0 {
				maxC = 1000
			}
			config.Servers = []ServerOptions{{
				OpenPort:   open,
				TunnelPort: tunnel,
				MaxConn:    maxC,
			}}
		} else if config.Mode == Client {
			local, _ := strconv.Atoi(params["localport"])
			idle, _ := strconv.Atoi(params["idleconn"])
			if idle <= 0 {
				idle = 1
			}
			config.Clients = []ClientOptions{{
				TunnelAddr: params["tunneladdr"],
				LocalPort:  local,
				IdleConn:   idle,
			}}
		} else {
			return nil, errors.New("mode must be 'server' or 'client'")
		}
		// 日志配置读取
		file, err := loadConfig()
		if file != nil && err != nil {
			if len(file.Log.Path) > 0 {
				config.Log.Path = file.Log.Path
			}
			if file.Log.Size > 0 {
				config.Log.Size = file.Log.Size
			}
			if file.Log.Age > 0 {
				config.Log.Age = file.Log.Age
			}
			if file.Log.Backups > 0 {
				config.Log.Backups = file.Log.Backups
			}
		}
		return config, nil
	}

	return loadConfig()
}

func loadConfig() (*ConfigOptions, error) {
	config := &ConfigOptions{}
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
