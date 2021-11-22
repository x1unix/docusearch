package store

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/x1unix/docusearch/internal/services/search"
	"github.com/x1unix/docusearch/internal/services/store/mocks"
	"go.uber.org/zap/zaptest"
)

//go:generate mockgen, ctrl -destination ./mocks/search.go -package mocks github.com/x1unix/docusearch/internal/services/search Provider
//go:generate mockgen -destination ./mocks/store.go -package mocks github.com/x1unix/docusearch/internal/services/store DocumentStore

func TestSyncedDocumentStore_AddDocument(t *testing.T) {
	cases := map[string]struct {
		name    string
		data    io.Reader
		wantErr string
		cfg     TextIndexConfig

		wantErrFn   func(err error) bool
		newStoreFn  func(t *testing.T, ctrl *gomock.Controller) DocumentStore
		newSearchFn func(t *testing.T, ctrl *gomock.Controller) search.Provider
	}{
		"should update document and words index on save": {
			name: "correct",
			data: strings.NewReader("The quick brown fox jumps over the lazy dog"),
			cfg:  TextIndexConfig{IgnoreCommonWords: true},

			newStoreFn: func(t *testing.T, ctrl *gomock.Controller) DocumentStore {
				store := mocks.NewMockDocumentStore(ctrl)
				store.EXPECT().
					AddDocument(gomock.Any(), "correct", matchReaderContents(t, []byte("The quick brown fox jumps over the lazy dog"))).
					Return(nil)
				return store
			},

			newSearchFn: func(t *testing.T, ctrl *gomock.Controller) search.Provider {
				sp := mocks.NewMockProvider(ctrl)
				expectWords := search.WordsFromString("The quick brown fox jumps over the lazy dog", search.EnglishCommonVerbs)
				sp.EXPECT().AddDocumentRef(gomock.Any(), "correct", stringsContentsMatch(t, expectWords)).Return(nil)
				return sp
			},
		},
		"should respect index settings": {
			name: "correct",
			data: strings.NewReader("The quick brown fox jumps over the lazy dog"),
			cfg:  TextIndexConfig{IgnoreCommonWords: false},

			newStoreFn: func(t *testing.T, ctrl *gomock.Controller) DocumentStore {
				store := mocks.NewMockDocumentStore(ctrl)
				store.EXPECT().
					AddDocument(gomock.Any(), "correct", matchReaderContents(t, []byte("The quick brown fox jumps over the lazy dog"))).
					Return(nil)
				return store
			},

			newSearchFn: func(t *testing.T, ctrl *gomock.Controller) search.Provider {
				sp := mocks.NewMockProvider(ctrl)
				expectWords := search.WordsFromString("The quick brown fox jumps over the lazy dog", nil)
				sp.EXPECT().AddDocumentRef(gomock.Any(), "correct", stringsContentsMatch(t, expectWords)).Return(nil)
				return sp
			},
		},
		"should raise errors from inner storage": {
			name: "bad",
			data: strings.NewReader("foobar"),
			wantErrFn: func(err error) bool {
				return errors.Is(err, fs.ErrExist)
			},
			newStoreFn: func(t *testing.T, ctrl *gomock.Controller) DocumentStore {
				store := mocks.NewMockDocumentStore(ctrl)
				store.EXPECT().
					AddDocument(gomock.Any(), "bad", matchReaderContents(t, []byte("foobar"))).
					Return(fs.ErrExist)
				return store
			},
			newSearchFn: func(t *testing.T, ctrl *gomock.Controller) search.Provider {
				return nil
			},
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			syncStore := NewSyncedDocumentStore(zaptest.NewLogger(t), c.newStoreFn(t, ctrl), c.newSearchFn(t, ctrl), c.cfg)

			err := syncStore.AddDocument(context.TODO(), c.name, c.data)
			if c.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), c.wantErr)
				return
			}

			if c.wantErrFn != nil {
				require.Error(t, err)
				require.True(t, c.wantErrFn(err))
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestSyncedDocumentStore_RemoveDocument(t *testing.T) {
	cases := map[string]struct {
		name    string
		wantErr string

		wantErrFn   func(err error) bool
		newStoreFn  func(t *testing.T, ctrl *gomock.Controller) DocumentStore
		newSearchFn func(t *testing.T, ctrl *gomock.Controller) search.Provider
	}{
		"should remove document from index": {
			name:    "foobar",
			wantErr: "failed to remove document from search index: test",
			newStoreFn: func(t *testing.T, ctrl *gomock.Controller) DocumentStore {
				store := mocks.NewMockDocumentStore(ctrl)
				store.EXPECT().RemoveDocument(gomock.Any(), "foobar").Return(nil)
				return store
			},
			newSearchFn: func(t *testing.T, ctrl *gomock.Controller) search.Provider {
				sp := mocks.NewMockProvider(ctrl)
				sp.EXPECT().RemoveDocumentRef(gomock.Any(), "foobar").Return(errors.New("test"))
				return sp
			},
		},
		"should stop if document wasn't removed properly": {
			name: "foobar",
			wantErrFn: func(err error) bool {
				return errors.Is(err, fs.ErrNotExist)
			},
			newStoreFn: func(t *testing.T, ctrl *gomock.Controller) DocumentStore {
				store := mocks.NewMockDocumentStore(ctrl)
				store.EXPECT().RemoveDocument(gomock.Any(), "foobar").Return(fs.ErrNotExist)
				return store
			},
			newSearchFn: func(t *testing.T, ctrl *gomock.Controller) search.Provider {
				return nil
			},
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			syncedStore := NewSyncedDocumentStore(zaptest.NewLogger(t), c.newStoreFn(t, ctrl),
				c.newSearchFn(t, ctrl), TextIndexConfig{})
			err := syncedStore.RemoveDocument(nil, c.name)
			if c.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), c.wantErr)
				return
			}

			if c.wantErrFn != nil {
				require.Error(t, err)
				require.True(t, c.wantErrFn(err))
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestSyncedDocumentStore_GetDocument(t *testing.T) {
	wantErr := "error from inner call"
	ctrl := gomock.NewController(t)
	storeMock := mocks.NewMockDocumentStore(ctrl)
	storeMock.EXPECT().GetDocument("testdoc").Return(nil, errors.New(wantErr))
	syncStore := NewSyncedDocumentStore(nil, storeMock, nil, TextIndexConfig{})
	_, err := syncStore.GetDocument("testdoc")
	require.EqualError(t, err, wantErr)
}

type readerMatcher struct {
	t    *testing.T
	want []byte
}

func (m readerMatcher) Matches(v interface{}) bool {
	r, ok := v.(io.Reader)
	if !ok {
		return false
	}

	got, _ := ioutil.ReadAll(r)
	return assert.Equal(m.t, m.want, got)
}

func (m readerMatcher) String() string {
	return fmt.Sprintf("io.Reader with contents: %#v", m.want)
}

func matchReaderContents(t *testing.T, want []byte) gomock.Matcher {
	return readerMatcher{t: t, want: want}
}

type stringsContentsMatcher struct {
	t    *testing.T
	want []string
}

func stringsContentsMatch(t *testing.T, want []string) gomock.Matcher {
	return &stringsContentsMatcher{
		t:    t,
		want: want,
	}
}

func (m stringsContentsMatcher) Matches(v interface{}) bool {
	gotVal, ok := v.([]string)
	if !ok {
		return false
	}

	return assert.ElementsMatch(m.t, m.want, gotVal)
}

func (m stringsContentsMatcher) String() string {
	return fmt.Sprintln(m.want)
}
