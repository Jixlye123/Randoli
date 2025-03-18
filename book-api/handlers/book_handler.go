package handlers

import (
	"encoding/json"
	"net/http"

	"book-api/models"
	"book-api/storage"

	"github.com/google/uuid"
)

//GetAllBooks to retrieve all the books

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := storage.LoadBooks()
	if err != nil {
		http.Error(w, "Error loading books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var newBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//Generate unique BookId
	newBook.BookID = uuid.New().String()

	books, _ := storage.LoadBooks()
	books = append(books, newBook)

	if err := storage.SaveBooks(books); err != nil {
		http.Error(w, "Error saving books", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)

}
