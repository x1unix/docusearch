package services

import "context"

// DocumentSearcher is abstract document search implementation.
type DocumentSearcher interface {
	// SearchDocumentsByWord returns list of document IDs
	// that contain specified word.
	SearchDocumentsByWord(ctx context.Context, word string) ([]string, error)
}

// SearchProvider is abstract document search provider.
type SearchProvider interface {
	DocumentSearcher

	// AddDocumentRef adds references of specified words to document in search index.
	AddDocumentRef(ctx context.Context, docId string, words []string) error

	// RemoveDocumentRef removes any references to document from index.
	RemoveDocumentRef(ctx context.Context, docId string) error
}
