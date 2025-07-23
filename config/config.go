package config

import (
	"encoding/json"
	"os"
)

type WebConfig struct {
	SiteName string `json:"site_name"`
	DbName   string `json:"db_name"`
	Port     int    `json:"port"`
	Timeout  int    `json:"timeout"` // 单位：秒
}

var Cfg *WebConfig

func InitConfig() error {
	if Cfg != nil {
		return nil
	}

	SetDefaultConfig()

	file, err := os.ReadFile("./data/config.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, &Cfg); err != nil {
		return err
	}

	return nil
}

func SetDefaultConfig() {
	Cfg = &WebConfig{
		SiteName: "jing-sync",
		DbName:   "jing-sync.db",
		Port:     8888,
		Timeout:  60,
	}
}
