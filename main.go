package main

import (
	"log"
	"net/http"
	"quote-api/handler"
)

func main() {
	http.HandleFunc("/quote", handler.GetQuote)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
