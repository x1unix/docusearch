package store

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/x1unix/docusearch/internal/services/search"
	"github.com/x1unix/docusearch/internal/utils/collections"
	"go.uber.org/zap"
)

type TextIndexConfig struct {
	IgnoreCommonWords bool
}

// initBufferSize is initial buffer size for document parse buffer
const initBufferSize = 500 * 1024 // 500KB

// SyncedDocumentStore is facade over document storage implementation
// that keeps search index in sync on file upload/delete.
type SyncedDocumentStore struct {
	log            *zap.Logger
	store          DocumentStore
	searchProvider search.Provider
	filterList     collections.StringsSet
}

func NewSyncedDocumentStore(log *zap.Logger, store DocumentStore, searchProvider search.Provider, cfg TextIndexConfig) *SyncedDocumentStore {
	s := &SyncedDocumentStore{log: log, store: store, searchProvider: searchProvider}
	if cfg.IgnoreCommonWords {
		s.filterList = search.EnglishCommonVerbs
	}

	return s
}

// AddDocument implements DocumentStore
func (s SyncedDocumentStore) AddDocument(ctx context.Context, name string, data io.Reader) error {
	buff := new(bytes.Buffer)
	buff.Grow(initBufferSize)

	teeReader := io.TeeReader(data, buff)
	if err := s.store.AddDocument(ctx, name, teeReader); err != nil {
		return err
	}

	words := search.WordsFromString(buff.String(), s.filterList)
	if err := s.searchProvider.AddDocumentRef(ctx, name, words); err != nil {
		return fmt.Errorf("failed to index document: %w", err)
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
