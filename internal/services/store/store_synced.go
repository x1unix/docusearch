package store

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/x1unix/docusearch/internal/services/search"
	"go.uber.org/zap"
)

// initBufferSize is initial buffer size for document parse buffer
const initBufferSize = 500 * 1024 // 500KB

// SyncedDocumentStore is facade over document storage implementation
// that keeps search index in sync on file upload/delete.
type SyncedDocumentStore struct {
	log            *zap.Logger
	store          DocumentStore
	searchProvider search.Provider
}

// AddDocument implements DocumentStore
func (s SyncedDocumentStore) AddDocument(ctx context.Context, name string, data io.Reader) error {
	buff := new(bytes.Buffer)
	buff.Grow(initBufferSize)

	teeReader := io.TeeReader(data, buff)
	if err := s.store.AddDocument(ctx, name, teeReader); err != nil {
		return err
	}

	return nil
}

// RemoveDocument implements DocumentStore
func (s SyncedDocumentStore) RemoveDocument(ctx context.Context, name string) error {
	if err := s.store.RemoveDocument(ctx, name); err != nil {
		return err
	}

	if err := s.searchProvider.RemoveDocumentRef(ctx, name); err != nil {
		return fmt.Errorf("failed to remove document from search index: %w", err)
	}

	return nil
}

// GetDocument implements DocumentStore
func (s SyncedDocumentStore) GetDocument(name string) (io.ReadCloser, error) {
	return s.store.GetDocument(name)
}
