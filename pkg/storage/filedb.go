package storage

import (
	"bufio"
	"log"
	"os"
)

type FileDB struct {
	file *os.File
}

func NewFileDB(filepath string) (*FileDB, error) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &FileDB{file: file}, nil
}

func (db *FileDB) Records() ([]string, error) {
	var records []string
	_, err := db.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(db.file)
	for scanner.Scan() {
		records = append(records, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func (db *FileDB) Append(value string) error {
	datawriter := bufio.NewWriter(db.file)
	_, err := datawriter.WriteString(value + "\n")
	if err != nil {
		log.Printf("writing file: %s", err)
		return err
	}
	err = datawriter.Flush()
	if err != nil {
		log.Printf("flushing file: %s", err)
		return err
	}
	return nil
}

func (db *FileDB) Contains(value string) bool {
	records, err := db.Records()
	if err != nil {
		log.Println("Error:", err)
		return false
	}

	for _, record := range records {
		if record == value {
			return true
		}
	}

	return false
}
