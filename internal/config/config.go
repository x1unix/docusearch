package config

import (
	"fmt"
	"os"

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
		// Address is Redis server address
		Address string `yaml:"address"`
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
