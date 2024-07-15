package base

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type ConfigOptions struct {
	Application ApplicationOptions `yaml:"application" json:"application"`
}

type ApplicationOptions struct {
	Protocol string        `yaml:"protocol" json:"protocol"`
	Mode     string        `yaml:"mode" json:"mode"`
	Server   ServerOptions `yaml:"server" json:"server"`
	Client   ClientOptions `yaml:"client" json:"client"`
	Log      LogOptions    `yaml:"log" json:"log"`
}

type ServerOptions struct {
	ListenPort int `yaml:"listen-port" json:"listenPort"`
	ClientPort int `yaml:"client-port" json:"clientPort"`
}

type ClientOptions struct {
	RemoteAddr string `yaml:"remote-addr" json:"remoteAddr"`
	LocalPort  int    `yaml:"local-port" json:"localPort"`
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
		fmt.Printf("\033[1;31;40m%s\033[0m\n", fmt.Sprintf("配置加载失败: %s", err.Error()))
		time.Sleep(3 * time.Second)
		panic(err)
	}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Printf("\033[1;31;40m%s\033[0m\n", fmt.Sprintf("配置加载失败: %s", err.Error()))
		time.Sleep(3 * time.Second)
		panic(err)
	}
	marshal, _ := json.Marshal(config.Application)
	fmt.Printf("\033[1;33;40m%s\033[0m\n", "配置加载完毕: "+string(marshal))
}
