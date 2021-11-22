package store

import (
	"bytes"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileDocumentStore_AddDocument(t *testing.T) {
	cases := map[string]struct {
		name    string
		data    []byte
		wantErr error
		preRun  func(t *testing.T, testdir string) error
	}{
		"creates file": {
			name: "foobar",
			data: []byte{0x50, 0x4B, 0x03, 0x04},
		},
		"check if file exist": {
			name:    "exist",
			wantErr: fs.ErrExist,
			preRun: func(t *testing.T, testdir string) error {
				return ioutil.WriteFile(filepath.Join(testdir, "exist"), []byte{'M', 'Z'}, 0777)
			},
		},
	}

	tmpDir, err := ioutil.TempDir(os.TempDir(), "test-store-fs-*")
	require.NoError(t, err, "failed to create temp dir")
	defer assert.NoError(t, os.RemoveAll(tmpDir), "failed to remove temp dir")

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			if c.preRun != nil {
				require.NoError(t, c.preRun(t, tmpDir), "preRun failed")
			}

			s := NewFileDocumentStore(tmpDir)
			err := s.AddDocument(nil, c.name, bytes.NewBuffer(c.data))
			if c.wantErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, c.wantErr))
				return
			}
			require.NoError(t, err)

			got, err := ioutil.ReadFile(filepath.Join(tmpDir, c.name))
			require.NoError(t, err, "failed to read created file")
			require.Equal(t, c.data, got, "created file and input mismatch")
		})
	}
}

func TestFileDocumentStore_RemoveDocument(t *testing.T) {
	cases := map[string]struct {
		name    string
		wantErr error
		preRun  func(t *testing.T, testdir string) error
	}{
		"removes file": {
			name: "exists",
			preRun: func(t *testing.T, testdir string) error {
				require.NoError(t, os.MkdirAll(testdir, os.ModeSticky|os.ModePerm))
				return ioutil.WriteFile(filepath.Join(testdir, "exists"), []byte{'M', 'Z'}, 0777)
			},
		},
		"check if file exists": {
			name:    "not-exist",
			wantErr: fs.ErrNotExist,
		},
	}

	tmpDir, err := ioutil.TempDir(os.TempDir(), "test-store-fs-*")
	require.NoError(t, err, "failed to create temp dir")
	defer assert.NoError(t, os.RemoveAll(tmpDir), "failed to remove temp dir")

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			if c.preRun != nil {
				require.NoError(t, c.preRun(t, tmpDir), "preRun failed")
			}

			s := NewFileDocumentStore(tmpDir)
			err := s.RemoveDocument(nil, c.name)
			if c.wantErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, c.wantErr))
				return
			}
			require.NoError(t, err)

			_, err = os.Stat(filepath.Join(tmpDir, c.name))
			require.Error(t, err, "file still accessible after delete")
			require.Truef(t, errors.Is(err, fs.ErrNotExist), "invalid error (got: %v)", err)
		})
	}
}

func TestFileDocumentStore_GetDocument(t *testing.T) {
	cases := map[string]struct {
		name    string
		data    []byte
		wantErr error
		preRun  func(t *testing.T, testdir string) error
	}{
		"return error if not exists": {
			name:    "not-exist",
			wantErr: fs.ErrNotExist,
		},
		"get existing file": {
			name: "exist",
			data: []byte{0x50, 0x4B, 0x03, 0x04},
			preRun: func(t *testing.T, testdir string) error {
				require.NoError(t, os.MkdirAll(testdir, os.ModeSticky|os.ModePerm))
				return ioutil.WriteFile(filepath.Join(testdir, "exist"), []byte{0x50, 0x4B, 0x03, 0x04}, 0777)
			},
		},
	}

	tmpDir, err := ioutil.TempDir(os.TempDir(), "test-store-fs-*")
	require.NoError(t, err, "failed to create temp dir")
	defer assert.NoError(t, os.RemoveAll(tmpDir), "failed to remove temp dir")

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			if c.preRun != nil {
				require.NoError(t, c.preRun(t, tmpDir), "preRun failed")
			}

			s := NewFileDocumentStore(tmpDir)
			f, err := s.GetDocument(c.name)
			if c.wantErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, c.wantErr))
				return
			}
			require.NoError(t, err)
			defer f.Close()
			got, err := ioutil.ReadAll(f)
			require.NoError(t, err, "failed to read file from returned reader")
			require.Equal(t, c.data, got, "created file and input mismatch")
		})
	}
}
