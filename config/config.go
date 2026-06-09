// 读取配置文件并提供全局访问接口。
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

	// Viper 的 SetConfigName 只接收不带扩展名的名字，SetConfigType 再指定 yaml。
	// 这样 APP_CONFIG_FILE 既可以写 config，也可以写 config.yaml。
	v.SetConfigName(cfgFile)
	v.SetConfigType("yaml")
	// 查找顺序保留 manifest/config 优先，便于 docker 镜像和本地运行使用同一套默认路径；
	// 追加 "." 是为了支持临时把配置放在工作目录进行调试。
	v.AddConfigPath("manifest/config")
	v.AddConfigPath(".")

	// AutomaticEnv 只让 Viper 感知环境变量；它不会自动把 APP_PORT 映射到 server.address
	// 这类嵌套字段，所以本文件下面对少量运行时变量做了显式覆盖。
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
