// 提供基于项目配置创建日志记录器的辅助方法。
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"server_go/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg config.LoggerConfig) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	// zapcore.Level 支持 debug/info/warn/error 等文本解析；这里启动时失败比静默降级更安全，
	// 否则线上可能以错误级别运行很久才被发现。
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, fmt.Errorf("parse logger level %q: %w", cfg.Level, err)
	}

	var zapCfg zap.Config
	switch loggerFormat(cfg) {
	case "json":
		// ProductionConfig 默认使用 JSON encoder。容器/日志平台通常按行采集 JSON，
		// 扩展知识：如果后续接 ELK/Loki/云日志，字段化日志比 fmt 字符串更容易检索和聚合。
		zapCfg = zap.NewProductionConfig()
	default:
		// DevelopmentConfig 偏向本地可读性；这里把时间格式调成 Gin 风格，
		// 方便和 Gin/HTTP 访问日志一起观察。
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeTime = ginStyleTimeEncoder
	}

	// AtomicLevel 允许 zap 在运行期调整日志级别；当前项目未暴露动态接口，
	// 但保留这个结构，后续接管理端点时无需重建 logger。
	zapCfg.Level = zap.NewAtomicLevelAt(level)
	zapCfg.OutputPaths = append([]string{}, cfg.OutputPaths...)
	if cfg.Stdout {
		// zap 约定 "stdout"/"stderr" 是特殊 sink；普通字符串会当作文件路径打开。
		zapCfg.OutputPaths = append(zapCfg.OutputPaths, "stdout")
	}
	if err := ensureLogDirs(zapCfg.OutputPaths); err != nil {
		return nil, err
	}
	if cfg.ClearOnStart {
		if err := clearLogFiles(zapCfg.OutputPaths); err != nil {
			return nil, err
		}
	}

	log, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}
	return log, nil
}

func loggerFormat(cfg config.LoggerConfig) string {
	format := strings.ToLower(strings.TrimSpace(cfg.Format))
	switch format {
	case "json", "console":
		return format
	default:
		return "console"
	}
}

func ginStyleTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 - 15:04:05"))
}

func ensureLogDirs(paths []string) error {
	for _, path := range paths {
		if path == "" || path == "stdout" || path == "stderr" {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("create log dir for %s: %w", path, err)
		}
	}
	return nil
}

func clearLogFiles(paths []string) error {
	for _, path := range paths {
		if path == "" || path == "stdout" || path == "stderr" {
			continue
		}
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("clear log file %s: %w", path, err)
		}
		if err := file.Close(); err != nil {
			return fmt.Errorf("close cleared log file %s: %w", path, err)
		}
	}
	return nil
}
