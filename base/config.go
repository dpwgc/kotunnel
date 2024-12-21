package base

import (
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

const (
	Server = "server"
	Client = "client"
)

type Config struct {
	Servers []ServerConfig `yaml:"servers" json:"servers"`
	Clients []ClientConfig `yaml:"clients" json:"clients"`
	Log     LogConfig      `yaml:"log" json:"log"`
}

type ServerConfig struct {
	Secret     string `yaml:"secret" json:"secret"`
	OpenPort   int    `yaml:"open-port" json:"openPort"`
	TunnelPort int    `yaml:"tunnel-port" json:"tunnelPort"`
	MaxConn    int    `yaml:"max-conn" json:"maxConn"`
}

type ClientConfig struct {
	Secret        string `yaml:"secret" json:"secret"`
	TunnelAddr    string `yaml:"tunnel-addr" json:"tunnelAddr"`
	LocalPort     int    `yaml:"local-port" json:"localPort"`
	IdleConn      int    `yaml:"idle-conn" json:"idleConn"`
	RetryInterval int    `yaml:"retry-interval" json:"retryInterval"`
}

type LogConfig struct {
	ClosePrint bool   `yaml:"close-print" json:"closePrint"`
	Path       string `yaml:"path" json:"path"`
	Size       int    `yaml:"size" json:"size"`
	Age        int    `yaml:"age" json:"age"`
	Backups    int    `yaml:"backups" json:"backups"`
}

func GetConfig(args []string) *Config {

	// ./main -server secret=123456 openPort=8080 tunnelPort=9090 maxConn=1000
	// ./main -client secret=123456 tunnelAddr=0.0.0.0:9090 localPort=7070 idleConn=1 retryInterval=5

	mode := ""
	params := make(map[string]string)
	for _, v := range args[1:] {
		if strings.HasPrefix(v, "-") {
			mode = strings.ReplaceAll(v, "-", "")
			continue
		}
		arr := strings.Split(v, "=")
		if len(arr) > 1 {
			params[strings.ToLower(arr[0])] = strings.Join(arr[1:], "")
		}
	}

	// 日志配置读取
	config, err := loadConfig()
	if err != nil {
		Warn("read 'config.yaml' error: " + err.Error())
	}
	if config == nil {
		config = &Config{}
	}

	if mode == Server {
		config.Servers = append(config.Servers, ServerConfig{
			Secret:     params["secret"],
			OpenPort:   toInt(params["openport"]),
			TunnelPort: toInt(params["tunnelport"]),
			MaxConn:    toInt(params["maxconn"]),
		})
	} else if mode == Client {
		config.Clients = append(config.Clients, ClientConfig{
			Secret:        params["secret"],
			TunnelAddr:    params["tunneladdr"],
			LocalPort:     toInt(params["localport"]),
			IdleConn:      toInt(params["idleconn"]),
			RetryInterval: toInt(params["retryinterval"]),
		})
	}

	return config
}

func loadConfig() (*Config, error) {
	config := &Config{}
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

func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
