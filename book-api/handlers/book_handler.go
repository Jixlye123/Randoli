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

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := storage.LoadBooks()
	if err != nil {
		http.Error(w, "Error loading books", http.StatusInternalServerError)
		return
	}

	// OPTIMIZE - Pagination Query Params
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

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

	if newBook.BookID == "" || newBook.AuthorID == "" || newBook.PublisherID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	books, _ := storage.LoadBooks()

	for _, book := range books {
		if book.BookID == newBook.BookID {
			http.Error(w, "Book ID already exists", http.StatusConflict)
			return
		}
	}
	books = append(books, newBook)

	if err := storage.SaveBooks(books); err != nil {
		http.Error(w, "Error saving books", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)

}

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

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]

	books, err := storage.LoadBooks()
	if err != nil {
		http.Error(w, "Error loading books", http.StatusInternalServerError)
		return
	}

	var updatedBooks []models.Book
	bookFound := false

	for _, book := range books {
		if book.BookID == bookID {
			bookFound = true
			continue
		}
		updatedBooks = append(updatedBooks, book)
	}

	if !bookFound {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	if err := storage.SaveBooks(updatedBooks); err != nil {
		http.Error(w, "Error saving books", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

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

	var wg sync.WaitGroup
	results := make(chan models.Book, len(books))

	for _, book := range books {
		wg.Add(1)
		go func(b models.Book) {
			defer wg.Done()

			//CHECKING IF THE BOOK MATCHES THE QUERY
			fmt.Printf("Checking book: %s\n", b.Title)

			if matchesQuery(b, query) {
				fmt.Printf("Match found for book: %s\n", b.Title)
				results <- b
			} else {
				results <- models.Book{}
			}
		}(book)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

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
	//DEBUG 2
	fmt.Printf("Checking Book: %s | Title Match: %v | Description Match: %v\n", book.Title, titleMatch, descriptionMatch)

	return titleMatch || descriptionMatch
}
