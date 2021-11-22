package config

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

// DefaultFileName is default config file name.
const DefaultFileName = "config.yaml"

// Config is application configuration
type Config struct {
	// Production toggles production mode
	Production bool `yaml:"production"`

	HTTP struct {
		// Listen is HTTP listen address
		Listen string `yaml:"listen"`
	} `yaml:"http"`

	Redis struct {
		// URL is Redis server connection string.
		URL string `yaml:"url"`
	}

	Storage struct {
		// UploadsDirectory is uploaded files storage directory
		UploadsDirectory string `yaml:"uploads_dir"`
	}

	Log struct {
		Level zapcore.Level `yaml:"level"`
	}
}

// Logger returns a new logger
func (cfg Config) Logger() (*zap.Logger, error) {
	var logCfg zap.Config
	if cfg.Production {
		logCfg = zap.NewProductionConfig()
	} else {
		logCfg = zap.NewDevelopmentConfig()
	}

	logCfg.Level = zap.NewAtomicLevelAt(cfg.Log.Level)
	return logCfg.Build()
}

// RedisClient returns a new redis client from string
func (cfg Config) RedisClient() (*redis.Client, error) {
	connCfg, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid redis connection URL: %w", err)
	}

	return redis.NewClient(connCfg), nil
}

// FromFile loads configuration from file.
func FromFile(fileName string) (*Config, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	cfg := new(Config)
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %q: %w", fileName, err)
	}

	return cfg, nil
}
