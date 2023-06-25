package testenv

import (
	"bitcoinrateapp/pkg/core"
	"os"
	"testing"
)

func NewTemporaryFileDB(t *testing.T) *core.FileDB {
	file, err := os.CreateTemp(os.TempDir(), "prefix")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(file.Name()) })

	db, err := core.NewFileDB(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	return db
}
