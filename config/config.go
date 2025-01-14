package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	UdpServer *UdpServer `toml:"udpserver"`
}

type UdpServer struct {
	Address           string `toml:"address"`
	HeartBeatInterval int    `toml:"heartBeatInterval"`
}

// LoadConfig loads the configuration from a TOML file
func LoadConfig() (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
