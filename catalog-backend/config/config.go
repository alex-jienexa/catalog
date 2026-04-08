package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Upload   UploadConfig
	JWT      JWTConfig
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
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

type UploadConfig struct {
	Path    string `mapstructure:"path"`
	MaxSize int64  `mapstructure:"max_size"` // в байтах
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Значения по умолчанию
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 10)
	viper.SetDefault("server.write_timeout", 10)
	viper.SetDefault("database.path", "./catalog.db")
	viper.SetDefault("upload.path", "./uploads")
	viper.SetDefault("upload.max_size", 5*1024*1024) // 5 MB
	viper.SetDefault("jwt.secret", "change-me-in-production-please")

	// Пытаемся прочитать файл, если он есть
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Конфиг файл не найден (%s), используем переменные окружения и значения по умолчанию", path)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
