package storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"book-api/models"
)

var mu sync.Mutex

const dataFile = "books.json"

func LoadBooks() ([]models.Book, error) {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile(dataFile, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var books []models.Book
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&books); err != nil {
		if errors.Is(err, io.EOF) {
			return []models.Book{}, nil
		}
		return nil, err
	}

	return books, nil
}

func SaveBooks(books []models.Book) error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(books)
}
