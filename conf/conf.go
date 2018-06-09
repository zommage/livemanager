package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type Config struct {
	BaseConf `json:"BaseConf"`
	DbConf   `json:"DbConf"`
	LogConf  `json:"LogConf"`
}

type BaseConf struct {
	HttpPort   string `json:"HttpPort"`   // http port
	HttpsPort  string `json:"HttpsPort"`  // https port
	SslKey     string `json:"SslKey"`     // ssl key
	SslCrt     string `json:"SslCrt"`     // ssl crt
	Env        string `json:"Env"`        // 环境信息
	RsaSertKey string `json:"RsaSertKey"` // rsa 的私钥路径
}

// mysql db config
type DbConf struct {
	DbName        string `json:"DbName"`
	DbHost        string `json:"DbHost"`
	DbPort        string `json:"DbPort"`
	DbUser        string `json:"DbUser"`
	DbPassword    string `json:"DbPassword"`
	DbLogEnable   bool   `json:"DbLogEnable"`
	DbMaxConnect  int    `json:"DbMaxConnect"`
	DbIdleConnect int    `json:"DbIdleConnect"`
}

// Log config
type LogConf struct {
	LogPath  string `json:"LogPath"`
	LogLevel string `json:"LogLevel"`
}

var Conf = new(Config)

func InitConfig(confPath *string) (*Config, error) {
	confBytes, err := ioutil.ReadFile(*confPath)
	if err != nil {
		tmpStr := fmt.Sprintf("read config %v err: %v", *confPath, err)
		return nil, errors.New(tmpStr)
	}

	err = json.Unmarshal(confBytes, Conf)
	if err != nil {
		tmpStr := fmt.Sprintf("fmt config to json err: %v", err)
		return nil, errors.New(tmpStr)
	}

	return Conf, nil
}
