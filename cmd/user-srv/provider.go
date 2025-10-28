package main

import (
	"flag"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/spf13/viper"

	"github.com/kyson/e-shop-native/internal/user-srv/conf"
)

func ProvideServerConfig(c *conf.Bootstrap) *conf.Server {
	return c.Server
}

func ProvideDataConfig(c *conf.Bootstrap) *conf.Data {
	return c.Data
}

func ProvideAuthConfig(c *conf.Bootstrap) *conf.Auth {
	return c.Auth
}

func ProvideLogConfig(c *conf.Bootstrap) *conf.Log {
	return c.Log
}

func LoadConfig() (*conf.Bootstrap, error) {
	flag.Parse()

	// viper
	v := viper.New()
	// 设置配置文件
	v.SetConfigFile(flagconf)
	v.SetConfigType("yaml")

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 将配置 unmarshal 到 conf.Bootstrap
	var bc conf.Bootstrap
	if err := v.Unmarshal(&bc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &bc, nil
}

func NewLogger(c *conf.Log) (*zap.Logger, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(c.Level)); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)

	// 根据格式配置 Encoder
	switch strings.ToLower(c.Format) {
	case "json":
		cfg.Encoding = "json"
		cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	case "console", "text", "plain":
		cfg.Encoding = "console"
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	case "logfmt":
		// zap 原生不支持 logfmt，需要自定义或使用第三方库
		cfg.Encoding = "console"
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	default:
		return nil, fmt.Errorf("unsupported log format: %s", c.Format)
	}

	// // 可选：配置输出路径
	// if len(c.OutputPaths) > 0 {
	// 	cfg.OutputPaths = c.OutputPaths
	// }

	// // 可选：配置错误输出路径
	// if len(c.ErrorOutputPaths) > 0 {
	// 	cfg.ErrorOutputPaths = c.ErrorOutputPaths
	// }

	return cfg.Build()
}
