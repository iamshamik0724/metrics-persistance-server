package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	UdpServer      *UdpServer      `toml:"udpserver"`
	DatabaseConfig *DatabaseConfig `toml:"database"`
	Server         *Server         `toml:"server"`
}

type UdpServer struct {
	Address           string `toml:"address"`
	HeartBeatInterval int    `toml:"heartBeatInterval"`
}

type DatabaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

type Server struct {
	Port string `toml:"port"`
}

// LoadConfig loads the configuration from a TOML file
func LoadConfig() (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
