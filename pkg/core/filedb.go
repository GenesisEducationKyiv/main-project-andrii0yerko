package bitcoinrateapp

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
)

type Storage[T any] interface {
	Append(T) error
	GetRecords() ([]T, error)
}

type FileDB struct {
	Filepath string
}

func (db *FileDB) GetRecords() ([]string, error) {
	var records []string
	if _, err := os.Stat(db.Filepath); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		log.Println("file does not exists")
		return records, nil
	}

	file, err := os.Open(db.Filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		records = append(records, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func (db *FileDB) checkExists(value string) bool {
	records, err := db.GetRecords()
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	for _, record := range records {
		if record == value {
			return true
		}
	}

	return false
}

func (db *FileDB) Append(value string) error {
	// TODO: lowercase, trim whitespaces
	if db.checkExists(value) {
		return fmt.Errorf("already exists: %s", value)
	}

	file, err := os.OpenFile(db.Filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close()

	datawriter := bufio.NewWriter(file)
	datawriter.WriteString(value + "\n")
	datawriter.Flush()
	return nil
}
