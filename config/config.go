package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	ApiURL string `json:"api_url,omitempty"`
}

func InitConfig() error {
	cfg := Config{"https://h2.nealfree.cf:20080/cfg"}
	byts, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Panic(err)
	}
	return os.WriteFile("setting.json", byts, 0666)
}
