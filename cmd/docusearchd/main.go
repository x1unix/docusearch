package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/x1unix/docusearch/internal/config"
	"go.uber.org/zap"
)

func main() {
	var cfgFile string
	flag.StringVar(&cfgFile, "config", config.DefaultFileName, "Config file name")
	flag.Parse()

	cfg, err := config.FromFile(cfgFile)
	if err != nil {
		fatal(err)
	}

	log, err := cfg.Logger()
	if err != nil {
		fatal(err)
	}
	defer log.Sync()
	if err := start(log, cfg); err != nil {
		log.Fatal("failed to start service", zap.Error(err))
	}
}

func start(log *zap.Logger, cfg *config.Config) error {
	return nil
}

func fatal(err error) {
	_, _ = fmt.Fprintln(os.Stderr, "fatal error:", err)
	os.Exit(2)
}
