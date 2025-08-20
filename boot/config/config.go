package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type WebConfig struct {
	SiteName string  `json:"site_name"`
	DbName   string  `json:"db_name"`
	Port     int     `json:"port"`
	Timeout  int     `json:"timeout"` // 单位：秒
	Env      EnvType `json:"env"`
}

var Cfg *WebConfig

func InitConfig(env EnvType) error {
	if Cfg != nil {
		return nil
	}

	SetDefaultConfig(env)

	file, err := os.ReadFile(fmt.Sprintf("./data/config_%s.json", string(env)))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, &Cfg); err != nil {
		return err
	}

	return nil
}

func SetDefaultConfig(env EnvType) {
	Cfg = &WebConfig{
		SiteName: "jing-sync",
		DbName:   "jing-sync.db",
		Port:     8888,
		Timeout:  60,
		Env:      env,
	}
}
