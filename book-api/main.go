package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	//temporary check route
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Book API is running")
	})

	fmt.Println("Starting server on port 8000")
	log.Fatal(http.ListenAndServe((":8000"), r))
}
