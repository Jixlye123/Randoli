package handlers

import (
	"encoding/json"
	"net/http"

	"book-api/models"
	"book-api/storage"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

//Get book by id (retrieving the book by id)

func GetBookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	books, _ := storage.LoadBooks()
	for _, book := range books {
		if book.BookID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)
			return
		}
	}

	http.Error(w, "Book not found", http.StatusNotFound)
}

//Updatebook updates a book by id

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	books, _ := storage.LoadBooks()
	for i, book := range books {
		if book.BookID == id {
			updatedBook.BookID = id
			books[i] = updatedBook
			storage.SaveBooks(books)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedBook)
			return
		}
	}

	http.Error(w, "Book not found", http.StatusNotFound)
}

//Deletes a book by id

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	books, _ := storage.LoadBooks()
	newBooks := []models.Book{}
	for _, book := range books {
		if book.BookID != id {
			newBooks = append(newBooks, book)
		}
	}

	if len(newBooks) == len(books) {
		http.Error(w, "Book not Found", http.StatusNotFound)
		return
	}

	storage.SaveBooks(newBooks)
	w.WriteHeader(http.StatusNoContent)
}
