package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-api/models"
	"book-api/storage"

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

	//Validate required firlds
	if newBook.BookID == "" || newBook.AuthorID == "" || newBook.PublisherID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	//Load existing books
	books, _ := storage.LoadBooks()

	//Ensure BookID is unique
	for _, book := range books {
		if book.BookID == newBook.BookID {
			http.Error(w, "Book ID already exists", http.StatusConflict)
			return
		}
	}

	// Add a new book

	books = append(books, newBook)

	// Save the updated books

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

//Serach books handles the seraching for books by title/description (case insensitive)

func SearchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing Search query", http.StatusBadRequest)
		return
	}

	books, err := storage.LoadBooks()
	if err != nil {
		http.Error(w, "Error loading the books", http.StatusInternalServerError)
		return
	}

	//Use goroutines for concurrent search

	resultChan := make(chan models.Book, len(books))

	//Launch the goroutines for the search

	for _, book := range books {
		go func(b models.Book) {
			if matchesQuery(b, query) { //matchesquery undefined error
				resultChan <- b
			} else {
				resultChan <- models.Book{} //Empty result for no match
			}
		}(book)
	}

	var results []models.Book
	for range books {
		book := <-resultChan
		if book.BookID != "" {
			results = append(results, book)
		}
	}

	close(resultChan)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func matchesQuery(book models.Book, query string) bool {
	query = strings.ToLower(query)
	return strings.Contains(strings.ToLower(book.Title), query) ||
		strings.Contains(strings.ToLower(book.Genre), query)
}
