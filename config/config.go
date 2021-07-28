package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// --------------------------------------------------------------------
// 类型定义
// --------------------------------------------------------------------
type gateway_proxy struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type gateway_proxy_reverse struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type microserver struct {
	Host      string `json:"host"`
	Http_port string `json:"http_port"`
	Grpc_port string `json:"grpc_port"`
	Name      string `json:"name"`
	ID        string `json:"id"`
}

type consul struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type mongodb struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

type system_config struct {
	Consul  consul  `json:"consul"`
	Mongodb mongodb `json:"mongodb"`

	Gateway_proxy         gateway_proxy         `json:"gateway_proxy"`
	Gateway_proxy_reverse gateway_proxy_reverse `json:"gateway_proxy_reverse"`
	Keystone              microserver           `json:"keystone"`
	Message               microserver           `json:"message"`
	Device                microserver           `json:"device"`
}

// --------------------------------------------------------------------
// 全局变量
// --------------------------------------------------------------------

var __config = &system_config{}

// --------------------------------------------------------------------
// 初始化方法
// --------------------------------------------------------------------

func Init_config() {
	config_file, err := ioutil.ReadFile("/Users/haifengfang/itable/server/cctable/config/config.json")
	if err != nil {
		log.Fatalf("reading config.json is failure, %v", err)
	}

	err = json.Unmarshal(config_file, __config)
	if err != nil {
		log.Fatalf("decoding config.json is failure, %v", err)
	}

}

// --------------------------------------------------------------------
// API：导出方法
// --------------------------------------------------------------------

func Get_config() *system_config {
	return __config
}
