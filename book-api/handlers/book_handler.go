package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

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

	// OPTIMIZE - Parse limit and offset from query params
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	//Implementation of pagination

	start := offset
	end := offset + limit

	if start >= len(books) {
		start = len(books)
	}

	if end > len(books) {
		end = len(books)
	}

	paginatedBooks := books[start:end]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paginatedBooks)
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
	bookID := vars["id"]

	// Load books from storage
	books, err := storage.LoadBooks()
	if err != nil {
		http.Error(w, "Error loading books", http.StatusInternalServerError)
		return
	}

	// Search for the book and delete if found
	var updatedBooks []models.Book
	bookFound := false

	for _, book := range books {
		if book.BookID == bookID {
			bookFound = true
			continue // Skip the book to delete it
		}
		updatedBooks = append(updatedBooks, book)
	}

	if !bookFound {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Save updated book list
	if err := storage.SaveBooks(updatedBooks); err != nil {
		http.Error(w, "Error saving books", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

//Serach books handles the seraching for books by title/description (case insensitive)

func SearchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing Search query", http.StatusBadRequest)
		return
	}
	//Loading the books from storage
	books, err := storage.LoadBooks()
	if err != nil {
		http.Error(w, "Error loading the books", http.StatusInternalServerError)
		return
	}

	// Fixing the seearch error
	fmt.Printf("Searching for: %s\n", query)
	fmt.Printf("Total books in database: %d\n", len(books))

	//Using goroutines for concurrent search
	var wg sync.WaitGroup
	results := make(chan models.Book, len(books))

	//Launch the goroutines for the search

	for _, book := range books {
		wg.Add(1)
		go func(b models.Book) {
			defer wg.Done()

			//CHECKING IF THE BOOK MATCHES THE QUERY
			fmt.Printf("Checking book: %s\n", b.Title)

			if matchesQuery(b, query) { //matchesquery undefined error
				fmt.Printf("Match found for book: %s\n", b.Title)
				results <- b
			} else {
				results <- models.Book{} //Empty result for no match
			}
		}(book)
	}

	//CLOSING CHANNELS
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collecting search results
	var foundBooks []models.Book
	for b := range results {
		foundBooks = append(foundBooks, b)
	}

	fmt.Printf("Total matches: %d\n", len(foundBooks))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foundBooks)
}

func matchesQuery(book models.Book, query string) bool {
	query = strings.ToLower(query)
	titleMatch := strings.Contains(strings.ToLower(book.Title), query)
	descriptionMatch := strings.Contains(strings.ToLower(book.Description), query)

	fmt.Printf("Checking Book: %s | Title Match: %v | Description Match: %v\n", book.Title, titleMatch, descriptionMatch)

	return titleMatch || descriptionMatch
}
