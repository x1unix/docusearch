package store

import (
	"context"
	"io"
)

// DocumentStore is abstract document storage.
type DocumentStore interface {
	// AddDocument stores a new document.
	//
	// Should return fs.ErrExist if item already exists.
	AddDocument(ctx context.Context, name string, data io.Reader) error

	// RemoveDocument removes document from storage.
	//
	// Should return fs.ErrNotExist if item doesn't exist.
	RemoveDocument(ctx context.Context, name string) error

	// GetDocument returns document reader by name.
	//
	// Should return fs.ErrNotExist if item doesn't exist.
	GetDocument(name string) (io.ReadCloser, error)
}
