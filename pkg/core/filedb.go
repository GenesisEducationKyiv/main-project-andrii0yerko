package core

import (
	"bufio"
	"errors"
	"log"
	"os"
)

var ErrIsDuplicate = errors.New("is duplicate")

// Simple file storage
// each record is stored in a new line
// stores unique values only
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

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func (db *FileDB) checkExists(value string) bool {
	records, err := db.GetRecords()
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

// add new value to the file
// if the same line already exists returns ErrIsDuplicate
func (db *FileDB) Append(value string) error {
	if db.checkExists(value) {
		return ErrIsDuplicate
	}

	file, err := os.OpenFile(db.Filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
		return err
	}
	defer file.Close()

	datawriter := bufio.NewWriter(file)
	_, err = datawriter.WriteString(value + "\n")
	if err != nil {
		log.Printf("failed reading file: %s", err)
		return err
	}
	datawriter.Flush()
	return nil
}
