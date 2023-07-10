package testenv

import (
	"bitcoinrateapp/pkg/storage"
	"os"
	"testing"
)

func NewTemporaryFileDB(t *testing.T) (*storage.FileDB, *os.File) {
	file, err := os.CreateTemp(os.TempDir(), "prefix")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(file.Name()) })

	db, err := storage.NewFileDB(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	return db, file
}
