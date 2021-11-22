package web

import (
	"github.com/brpaz/echozap"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/x1unix/docusearch/internal/config"
	"github.com/x1unix/docusearch/internal/services/search"
	"github.com/x1unix/docusearch/internal/services/store"
	"go.uber.org/zap"
)

// NewService builds application service handler.
func NewService(log *zap.Logger, cfg *config.Config, redisConn redis.Cmdable) *echo.Echo {
	echo.NotFoundHandler = FancyHandleNotFound
	e := echo.New()
	e.Use(echozap.ZapLogger(log))
	e.Use(middleware.Recover())

	searchProvider := search.NewRedisProvider(log.Named("search.redis"), redisConn)
	syncStore := store.NewSyncedDocumentStore(log.Named("store"), store.NewFileDocumentStore(cfg.Storage.UploadsDirectory),
		searchProvider, store.TextIndexConfig{IgnoreCommonWords: cfg.Search.IgnoreCommonWords})
	docHandler := NewDocumentsHandler(log.Named("handler.docs"), syncStore)
	searchHandler := NewSearchHandler(log.Named("handler.search"), searchProvider)

	e.POST("/document/:id", docHandler.UploadDocument)
	e.GET("/document/:id", docHandler.GetDocument)
	e.DELETE("/document/:id", docHandler.DeleteDocument)
	e.GET("/search", searchHandler.SearchWord)
	return e
}
