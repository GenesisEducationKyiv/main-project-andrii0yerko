package core_test

import (
	"bitcoinrateapp/pkg/core"
	"errors"
	"os"
	"testing"
)

func NewFileDB(t *testing.T) *core.FileDB {
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

func TestAddNewToFileDB(t *testing.T) {
	value := "test@email.org"

	db := NewFileDB(t)
	err := db.Append(value)
	if err != nil {
		t.Error(err)
	}

	records, err := db.Records()
	if err != nil {
		t.Error(err)
	}

	if len(records) != 1 {
		t.Errorf("expected 1 record, got %d", len(records))
	}
}

func TestAddDuplicateToFileDB(t *testing.T) {
	value := "test@email.org"

	db := NewFileDB(t)
	err := db.Append(value)
	if err != nil {
		t.Error(err)
	}
	err = db.Append(value)
	if !errors.Is(err, core.ErrIsDuplicate) {
		t.Error(err)
	}
}
