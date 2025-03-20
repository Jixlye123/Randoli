package handlers

import (
	"book-api/models"
	"book-api/storage"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTest() {
	storage.SaveBooks([]models.Book{
		{
			BookID:      "1",
			Title:       "Go in Action",
			AuthorID:    "101",
			PublisherID: "201",
			Description: "abc",
		},
		{
			BookID:      "2",
			Title:       "Docker Deep Dive",
			AuthorID:    "102",
			PublisherID: "202",
			Description: "test02",
		},
	})
}

func TestGetAllBooks(t *testing.T) {
	setupTest()

	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal("Failed to create a request", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllBooks)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var books []models.Book
	if err := json.NewDecoder(rr.Body).Decode(&books); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if len(books) != 2 {
		t.Errorf("Expected 2 books, got %d", len(books))
	}

	if books[0].Title != "Go in Action" {
		t.Errorf("Unexpected book title: got %v, want %v", books[0].Title, "Go in Action")
	}
}

func TestGetAllBooks_Empty(t *testing.T) {
	// Clear book storage
	storage.SaveBooks([]models.Book{})

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllBooks)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var books []models.Book
	json.NewDecoder(rr.Body).Decode(&books)

	if len(books) != 0 {
		t.Errorf("Expected no books, got %d", len(books))
	}
}

func TestGetAllBooksByPagination(t *testing.T) {
	setupTest()

	req, _ := http.NewRequest("GET", "/books?limit=1&offset=1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllBooks)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status OK but got %v", rr.Code)
	}

	var books []models.Book
	json.Unmarshal(rr.Body.Bytes(), &books)

	if len(books) != 1 {
		t.Fatalf("Expected 1 book but got %v", len(books))
	}

	if books[0].Title != "Docker Deep Dive" {
		t.Errorf("Unexpected book title: got %s, want Docker Deep Dive", books[0].Title)
	}
}
