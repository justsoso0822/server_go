// Package config provides functionality to load and manage application configuration.
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig              `mapstructure:"server"`
	App      AppConfig                 `mapstructure:"app"`
	Database map[string]DatabaseConfig `mapstructure:"database"`
	Redis    map[string]RedisConfig    `mapstructure:"redis"`
	Logger   LoggerConfig              `mapstructure:"logger"`
}

type ServerConfig struct {
	Address         string `mapstructure:"address"`
	GracefulTimeout int    `mapstructure:"gracefulTimeout"`
}

type AppConfig struct {
	Keys []string `mapstructure:"keys"`
}

type DatabaseConfig struct {
	Link          string `mapstructure:"link"`
	Debug         bool   `mapstructure:"debug"`
	Cache         bool   `mapstructure:"cache"`
	MaxIdle       int    `mapstructure:"maxIdle"`
	MaxOpen       int    `mapstructure:"maxOpen"`
	MaxLifetime   int    `mapstructure:"maxLifetime"`
	SlowThreshold int    `mapstructure:"slowThreshold"`
}

type RedisConfig struct {
	Address string `mapstructure:"address"`
	Pass    string `mapstructure:"pass"`
	DB      int    `mapstructure:"db"`
}

type LoggerConfig struct {
	Level               string   `mapstructure:"level"`
	Format              string   `mapstructure:"format"`
	Stdout              bool     `mapstructure:"stdout"`
	StdoutColorDisabled bool     `mapstructure:"stdoutColorDisabled"`
	OutputPaths         []string `mapstructure:"outputPaths"`
	ClearOnStart        bool     `mapstructure:"clearOnStart"`
}

func Load() (*Config, error) {
	v := viper.New()

	cfgFile := os.Getenv("APP_CONFIG_FILE")
	if cfgFile == "" {
		cfgFile = "config.yaml"
	}
	cfgFile = strings.TrimSuffix(cfgFile, ".yaml")

	v.SetConfigName(cfgFile)
	v.SetConfigType("yaml")
	v.AddConfigPath("manifest/config")
	v.AddConfigPath(".")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if port := strings.TrimSpace(os.Getenv("APP_PORT")); port != "" {
		cfg.Server.Address = ":" + strings.TrimPrefix(port, ":")
	}

	if cfg.Server.Address == "" {
		cfg.Server.Address = ":7001"
	}
	if cfg.Server.GracefulTimeout <= 0 {
		cfg.Server.GracefulTimeout = 15
	}

	return &cfg, nil
}
