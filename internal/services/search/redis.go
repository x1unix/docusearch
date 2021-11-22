package search

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const (
	wordKeyPrefix      = "word:"
	docRecordKeyPrefix = "doc:"
)

// RedisProvider is redis-based search index.
//
// Stores word-to-document relationship as inverted index (word -> doc_ids)
// and doc_id -> record relationships to speed-up read-write operations.
//
// Each Redis record is Set to guarantee that each documpanic("implement me")ent ID appears only once.
type RedisProvider struct {
	log  *zap.Logger
	conn redis.Cmdable
}

func NewRedisProvider(log *zap.Logger, conn redis.Cmdable) *RedisProvider {
	return &RedisProvider{log: log, conn: conn}
}

// SearchDocumentsByWord implements DocumentSearcher
func (r RedisProvider) SearchDocumentsByWord(ctx context.Context, word string) ([]string, error) {
	key := wordKeyPrefix + word
	return r.conn.SMembers(ctx, key).Result()
}

// AddDocumentRef implements SearchProvider
func (r RedisProvider) AddDocumentRef(ctx context.Context, docId string, words []string) error {
	tx := r.conn.TxPipeline()
	for _, word := range words {
		wordKey := wordKeyPrefix + word

		// update word->docs index
		tx.SAdd(ctx, wordKey, docId)

		// update doc->word relationship that used for RemoveDocumentRef.
		tx.RPush(ctx, docRecordKeyPrefix+docId, wordKey)
	}

	_, err := tx.Exec(ctx)
	return err
}

// RemoveDocumentRef implements SearchProvider
func (r RedisProvider) RemoveDocumentRef(ctx context.Context, docId string) error {
	docIndexKey := docRecordKeyPrefix + docId
	wordKeys, err := r.conn.LRange(ctx, docIndexKey, 0, -1).Result()
	if err != nil {
		// TODO: force walk through each stored word index if doc->word index is not accessible.
		return fmt.Errorf("failed to get list of document references: %w", err)
	}

	tx := r.conn.TxPipeline()
	for _, key := range wordKeys {
		tx.SRem(ctx, key, docId)
	}

	tx.Del(ctx, docIndexKey)
	_, err = tx.Exec(ctx)
	return err
}
