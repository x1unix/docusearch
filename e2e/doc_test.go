package e2e

import (
	"bytes"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/x1unix/docusearch/pkg/api"
)

func TestUploadDocument(t *testing.T) {
	cleanData(t)
	files := []string{"kafka1.txt", "kafka2.txt", "pangram1.txt"}
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			data := readTestData(t, file)
			fileID := strings.TrimSuffix(file, filepath.Ext(file))

			t.Run("check_not_exists", func(t *testing.T) {
				// file should not exist
				_, err := client.GetDocument(fileID)
				assertResponseError(t, err, api.ErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "document not found",
				})
				assertResponseError(t, client.RemoveDocument(fileID), api.ErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "document not found",
				})
			})

			t.Run("check_upload", func(t *testing.T) {
				// Upload should succeed
				require.NoErrorf(t, client.AddDocument(fileID, bytes.NewReader(data)), "failed to upload %q", file)
				// Re-upload should fail
				assertResponseError(t, client.AddDocument(fileID, bytes.NewReader(data)), api.ErrorResponse{
					StatusCode: http.StatusBadRequest,
					Message:    "item already exists",
				})

				// Uploaded contents should be equal to original
				got, err := client.GetDocument(fileID)
				require.NoErrorf(t, err, "GetDocument - unexpected error")
				require.ElementsMatch(t, data, got, "file content mismatch")
			})

			t.Run("check_delete", func(t *testing.T) {
				// Delete should succeed
				require.NoErrorf(t, client.RemoveDocument(fileID), "document delete failed")
				_, err := client.GetDocument(fileID)
				assertResponseError(t, err, api.ErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "document not found",
				})
			})
		})
	}
}
