package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/x1unix/docusearch/internal/config"
	"github.com/x1unix/docusearch/internal/services/search"
	"github.com/x1unix/docusearch/internal/services/store"
	"github.com/x1unix/docusearch/internal/web"
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
	redisConn, err := cfg.RedisClient()
	if err != nil {
		return err
	}

	defer redisConn.Close()
	if err := redisConn.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	echo.NotFoundHandler = web.FancyHandleNotFound
	e := echo.New()
	e.Use(echozap.ZapLogger(log))
	e.Use(middleware.Recover())

	searchProvider := search.NewRedisProvider(log.Named("search.redis"), redisConn)
	syncStore := store.NewSyncedDocumentStore(log.Named("store"), store.NewFileDocumentStore(cfg.Storage.UploadsDirectory),
		searchProvider, store.TextIndexConfig{IgnoreCommonWords: cfg.Search.IgnoreCommonWords})
	docHandler := web.NewDocumentsHandler(log.Named("handler.docs"), syncStore)
	searchHandler := web.NewSearchHandler(log.Named("handler.search"), searchProvider)

	e.POST("/document/:id", docHandler.UploadDocument)
	e.GET("/document/:id", docHandler.GetDocument)
	e.DELETE("/document/:id", docHandler.DeleteDocument)
	e.GET("/search", searchHandler.SearchWord)

	srv := &http.Server{
		Addr:    cfg.HTTP.Listen,
		Handler: e,
	}

	log.Info("starting http server...", zap.String("addr", cfg.HTTP.Listen))
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func fatal(err error) {
	_, _ = fmt.Fprintln(os.Stderr, "fatal error:", err)
	os.Exit(2)
}
