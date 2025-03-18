package models

import "time"

//Book represents book entity
type Book struct {
	BookID          string    `json:"book_id"`
	AuthorID        string    `json:"author_id"`
	PublisherID     string    `json:"publisher_id"`
	Title           string    `json:"title"`
	PublicationDate time.Time `json:"publication_date"`
	ISBN            string    `json:"isbn"`
	Pages           int       `json:"pages"`
	Genre           string    `json:"genre"`
	Price           float64   `json:"price"`
	Quantity        int       `json:"quantity"`
}
