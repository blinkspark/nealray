package config

import (
	"encoding/json"
	"log"
	"os"
)

func GenDefaultConfig() {
	conf := Config{
		Servers: []ServerConfig{
			{
				Addr: ":80",
			},
		},
	}

	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		log.Panic(err)
	}
	os.WriteFile("config.json", data, 0666)
}

func LoadConfig(fname string) (*Config, error) {
	var conf Config
	data, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
