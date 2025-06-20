package base

import (
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Servers []ServerConfig `yaml:"servers" json:"servers"`
	Clients []ClientConfig `yaml:"clients" json:"clients"`
	Log     LogConfig      `yaml:"log" json:"log"`
}

type ServerConfig struct {
	Name       string `yaml:"name" json:"name"`
	Secret     string `yaml:"secret" json:"secret"`
	OpenPort   int    `yaml:"open-port" json:"openPort"`
	TunnelPort int    `yaml:"tunnel-port" json:"tunnelPort"`
	MaxConn    int    `yaml:"max-conn" json:"maxConn"`
}

type ClientConfig struct {
	Name          string `yaml:"name" json:"name"`
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

	// ./main -server name=test1 secret=123456 openPort=8080 tunnelPort=9090 maxConn=1000
	// ./main -client name=test1 secret=123456 tunnelAddr=0.0.0.0:9090 localPort=7070 idleConn=1 retryInterval=5

	mode, params := ArgsToParams(args[1:])

	// 日志配置读取
	config, err := loadConfig()
	if err != nil {
		Warn("read 'config.yaml' error: " + err.Error())
	}
	if config == nil {
		config = &Config{}
	}

	if mode == "server" {
		config.Servers = append(config.Servers, ServerConfig{
			Name:       params["name"],
			Secret:     params["secret"],
			OpenPort:   ToInt(params["openport"]),
			TunnelPort: ToInt(params["tunnelport"]),
			MaxConn:    ToInt(params["maxconn"]),
		})
	} else if mode == "client" {
		config.Clients = append(config.Clients, ClientConfig{
			Name:          params["name"],
			Secret:        params["secret"],
			TunnelAddr:    params["tunneladdr"],
			LocalPort:     ToInt(params["localport"]),
			IdleConn:      ToInt(params["idleconn"]),
			RetryInterval: ToInt(params["retryinterval"]),
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

func ArgsToParams(args []string) (string, map[string]string) {
	mode := ""
	params := make(map[string]string)
	for _, v := range args {
		if v == "" {
			continue
		}
		if !strings.Contains(v, "-") && !strings.Contains(v, "=") {
			continue
		}
		if strings.HasPrefix(v, "-") {
			mode = strings.ReplaceAll(v, "-", "")
			continue
		}
		if !strings.Contains(v, "=") {
			continue
		}
		arr := strings.Split(v, "=")
		if len(arr) > 1 {
			params[strings.ToLower(arr[0])] = strings.Join(arr[1:], "")
		}
	}
	return mode, params
}

func ToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
