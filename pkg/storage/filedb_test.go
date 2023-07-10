package storage_test

import (
	"bitcoinrateapp/pkg/testenv"
	"testing"
)

func TestAddNewToFileDB(t *testing.T) {
	value := "test@email.org"

	db, _ := testenv.NewTemporaryFileDB(t)
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
