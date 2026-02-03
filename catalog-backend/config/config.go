package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port         string
	Mode         string
	ReadTimeout  int
	WriteTimeout int
}

type DatabaseConfig struct {
	Path string
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 10)
	viper.SetDefault("server.write_timeout", 10)
	viper.SetDefault("database.path", "./catalog.db")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: could not read config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
