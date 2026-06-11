package config

import (
	"fmt"
	"strings"
	"github.com/kavya/content-engine/internal/logger"
	"github.com/spf13/viper"
)

type Config struct{
	Server ServerConfig `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct{
	Port int `mapstructure:"port"`
	Environment string `mapstructure:"environment"`
}

type DatabaseConfig struct{
	URL string `mapstructure:"url"`
}

func LoadConfig() (*Config, error){
	viper.SetDefault("server.port", 9090)
	viper.SetDefault("server.environment", "test")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err:= viper.ReadInConfig(); err!=nil{
		if _, ok := err.(viper.ConfigFileNotFoundError);
		!ok{
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		logger.WarningLog.Println("Config file not found, relying on environment variables/defaults")
	}
	var config Config

	if err := viper.Unmarshal(&config); err!=nil{
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &config, nil
}