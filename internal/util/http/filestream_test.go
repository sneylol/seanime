package httputil

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setTempDir(t *testing.T) string {
	dir := t.TempDir()
	t.Setenv("TMPDIR", dir)
	t.Setenv("TMP", dir)
	t.Setenv("TEMP", dir)
	return dir
}

func TestFileStreamCloseRemovesFile(t *testing.T) {
	tempDir := setTempDir(t)
	logger := zerolog.Nop()

	fs, err := NewFileStream(context.Background(), &logger, 1024)
	require.NoError(t, err)

	name := fs.file.Name()
	assert.Equal(t, filepath.Join(tempDir, "seanime", "filestream"), filepath.Dir(name))
	assert.FileExists(t, name)

	require.NoError(t, fs.Close())
	assert.NoFileExists(t, name)
}

func TestClearFileStreamCache(t *testing.T) {
	tempDir := setTempDir(t)
	logger := zerolog.Nop()

	// Simulate a stream that was never closed
	fs, err := NewFileStream(context.Background(), &logger, 1024)
	require.NoError(t, err)
	name := fs.file.Name()
	require.NoError(t, fs.file.Close())

	unrelated := filepath.Join(tempDir, "unrelated.tmp")
	require.NoError(t, os.WriteFile(unrelated, []byte("x"), 0644))

	ClearFileStreamCache(&logger)

	assert.NoFileExists(t, name)
	assert.FileExists(t, unrelated)
}

func TestClearFileStreamCacheMissingDir(t *testing.T) {
	setTempDir(t)
	logger := zerolog.Nop()

	ClearFileStreamCache(&logger)
}
