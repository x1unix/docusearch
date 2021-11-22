package web

import (
	"errors"
	"io"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/x1unix/docusearch/internal/services/store"
	"go.uber.org/zap"
)

type DocumentsHandler struct {
	log            *zap.Logger
	documentsStore store.DocumentStore
}

func NewDocumentsHandler(log *zap.Logger, s store.DocumentStore) *DocumentsHandler {
	return &DocumentsHandler{
		log:            log,
		documentsStore: s,
	}
}

func (h DocumentsHandler) UploadDocument(c echo.Context) error {
	docID := c.Param("id")

	body := c.Request().Body
	defer body.Close()
	if err := h.documentsStore.AddDocument(c.Request().Context(), docID, body); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return echo.NewHTTPError(http.StatusBadRequest, "item already exists")
		}

		h.log.Error("failed to save document", zap.String("id", docID), zap.Error(err))
		return err
	}

	return nil
}

func (h DocumentsHandler) DeleteDocument(c echo.Context) error {
	docID := c.Param("id")
	if err := h.documentsStore.RemoveDocument(c.Request().Context(), docID); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return echo.NewHTTPError(http.StatusNotFound, "document not found")
		}

		h.log.Error("failed to remove document", zap.String("id", docID), zap.Error(err))
		return err
	}
	return nil
}

func (h DocumentsHandler) GetDocument(c echo.Context) error {
	docID := c.Param("id")
	r, err := h.documentsStore.GetDocument(docID)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return echo.NewHTTPError(http.StatusNotFound, "document not found")
		}

		h.log.Error("failed to get document", zap.String("id", docID), zap.Error(err))
		return err
	}

	defer r.Close()
	_, err = io.Copy(c.Response(), r)
	return err
}
