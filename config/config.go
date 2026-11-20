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
	// AllowedOrigins 是浏览器跨域白名单，例如 "http://localhost:18090"。
	// 含 "*" 表示回显任意 Origin；留空时 local/test 环境回显任意 Origin，生产环境不放行跨域。
	AllowedOrigins []string `mapstructure:"allowedOrigins"`
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
	Address        string `mapstructure:"address"`
	Pass           string `mapstructure:"pass"`
	DB             int    `mapstructure:"db"`
	DialTimeoutMs  int    `mapstructure:"dialTimeoutMs"`
	ReadTimeoutMs  int    `mapstructure:"readTimeoutMs"`
	WriteTimeoutMs int    `mapstructure:"writeTimeoutMs"`
	PoolSize       int    `mapstructure:"poolSize"`
	MinIdleConns   int    `mapstructure:"minIdleConns"`
	PoolTimeoutMs  int    `mapstructure:"poolTimeoutMs"`
}

type LoggerConfig struct {
	Level        string   `mapstructure:"level"`
	Format       string   `mapstructure:"format"`
	Stdout       bool     `mapstructure:"stdout"`
	OutputPaths  []string `mapstructure:"outputPaths"`
	ClearOnStart bool     `mapstructure:"clearOnStart"`
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

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	if port := strings.TrimSpace(os.Getenv("APP_PORT")); port != "" {
		cfg.Server.Address = ":" + strings.TrimPrefix(port, ":")
	}

	// 跨域白名单支持用环境变量覆盖，便于同一镜像在不同机器/域名下部署时调整，无需改配置文件。
	if origins := strings.TrimSpace(os.Getenv("APP_CORS_ORIGINS")); origins != "" {
		cfg.App.AllowedOrigins = splitOrigins(origins)
	}

	if cfg.Server.Address == "" {
		cfg.Server.Address = ":7001"
	}
	if cfg.Server.GracefulTimeout <= 0 {
		cfg.Server.GracefulTimeout = 15
	}

	return &cfg, nil
}

// splitOrigins 把逗号分隔的来源列表拆成切片，忽略空项，便于写进 docker compose 环境变量。
func splitOrigins(raw string) []string {
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			origins = append(origins, part)
		}
	}
	return origins
}
