package web

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/x1unix/docusearch/internal/models"
	"github.com/x1unix/docusearch/internal/services/search"
	"go.uber.org/zap"
)

type SearchHandler struct {
	log            *zap.Logger
	searchProvider search.DocumentSearcher
}

func NewSearchHandler(log *zap.Logger, searchProvider search.DocumentSearcher) *SearchHandler {
	return &SearchHandler{log: log, searchProvider: searchProvider}
}

func (h SearchHandler) SearchWord(c echo.Context) error {
	query := strings.TrimSpace(c.QueryParam("q"))
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "empty search query")
	}

	ids, err := h.searchProvider.SearchDocumentsByWord(c.Request().Context(), query)
	if err != nil {
		h.log.Error("failed to get search results", zap.Error(err), zap.String("query", query))
		return err
	}

	return c.JSON(http.StatusOK, models.DocumentIDsResponse{
		IDs: ids,
	})
}
