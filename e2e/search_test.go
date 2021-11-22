package e2e

import (
	"bytes"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/x1unix/docusearch/internal/services/search"
	"github.com/x1unix/docusearch/pkg/api"
)

func TestSearch(t *testing.T) {
	cleanData(t)
	t.Run("invalid search query", func(t *testing.T) {
		_, err := client.SearchByWord("")
		assertResponseError(t, err, api.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "empty search query",
		})
	})

	t.Run("search after upload", func(t *testing.T) {
		files := []string{"kafka1.txt", "kafka2.txt", "pangram1.txt"}
		for _, file := range files {
			data := readTestData(t, file)
			fileID := strings.TrimSuffix(file, filepath.Ext(file))
			require.NoErrorf(t, client.AddDocument(fileID, bytes.NewReader(data)), "failed to upload %q", file)
		}

		expect := map[string][]string{
			"morning":   {"kafka1", "kafka2"},
			"GREGOR":    {"kafka1", "kafka2"},
			"Pitifully": {"kafka1"},
			"Waltz":     {"pangram1"},
		}

		for word, expectMatches := range expect {
			t.Run("search/"+word, func(t *testing.T) {
				gotIds, err := client.SearchByWord(word)
				require.NoError(t, err)
				require.ElementsMatch(t, expectMatches, gotIds)
			})
		}
	})

	t.Run("article or verb should not indexed", func(t *testing.T) {
		items := search.EnglishCommonVerbs.ToArray()
		for _, w := range items {
			t.Run(w, func(t *testing.T) {
				gotIds, err := client.SearchByWord(w)
				require.NoError(t, err)
				require.Empty(t, gotIds)
			})
		}
	})

	t.Run("search after delete", func(t *testing.T) {
		require.NoError(t, client.RemoveDocument("kafka2"))
		gotIds, err := client.SearchByWord("GREGOR")
		require.NoError(t, err)
		require.Equal(t, []string{"kafka1"}, gotIds)

		require.NoError(t, client.RemoveDocument("kafka1"))
		gotIds, err = client.SearchByWord("GREGOR")
		require.NoError(t, err)
		require.Empty(t, gotIds)
	})
}
