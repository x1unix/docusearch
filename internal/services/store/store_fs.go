package store

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileDocumentStore is filesystem document storage.
type FileDocumentStore struct {
	storageDir string
}

// AddDocument implements DocumentStore
func (f FileDocumentStore) AddDocument(_ context.Context, name string, data io.Reader) error {
	// Pre-create directory if not exists
	if err := os.MkdirAll(f.storageDir, os.ModeSticky|os.ModePerm); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Create a file but check if file already exists.
	// os.OpenFile returns fs.ErrExist if file already exists.
	fd, err := os.OpenFile(filepath.Join(f.storageDir, name), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}

	defer fd.Close()
	_, err = io.Copy(fd, data)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// RemoveDocument implements DocumentStore
func (f FileDocumentStore) RemoveDocument(_ context.Context, name string) error {
	// os.Remove returns fs.ErrNotExists if file not exists.
	return os.Remove(filepath.Join(f.storageDir, name))
}

// GetDocument implements DocumentStore
func (f FileDocumentStore) GetDocument(name string) (io.ReadCloser, error) {
	// os.Open returns fs.ErrNotExists if file not exists.
	return os.Open(filepath.Join(f.storageDir, name))
}

func NewFileDocumentStore(storageDir string) *FileDocumentStore {
	return &FileDocumentStore{storageDir: storageDir}
}
