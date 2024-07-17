package base

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type ConfigOptions struct {
	App AppOptions `yaml:"app" json:"app"`
}

type AppOptions struct {
	Protocol string          `yaml:"protocol" json:"protocol"`
	Mode     string          `yaml:"mode" json:"mode"`
	Servers  []ServerOptions `yaml:"servers" json:"servers"`
	Clients  []ClientOptions `yaml:"clients" json:"clients"`
	Log      LogOptions      `yaml:"log" json:"log"`
}

type ServerOptions struct {
	OpenPort   int `yaml:"open-port" json:"openPort"`
	ClientPort int `yaml:"client-port" json:"clientPort"`
}

type ClientOptions struct {
	RemoteAddr string `yaml:"remote-addr" json:"remoteAddr"`
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

func InitConfig() {
	//加载客户端配置
	configBytes, err := os.ReadFile("./config.yaml")
	if err != nil {
		Println(31, 40, fmt.Sprintf("read config error: %s", err.Error()))
		time.Sleep(3 * time.Second)
		panic(err)
	}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		Println(31, 40, fmt.Sprintf("parse config error: %s", err.Error()))
		time.Sleep(3 * time.Second)
		panic(err)
	}
}
