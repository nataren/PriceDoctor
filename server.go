package main

import (
	"net/http"
	"log"
)

func main() {
	log.Println("Starting server")
	http.Handle("/", http.FileServer(http.Dir("./public/")))

	log.Println("Listing on port 8080")
	http.ListenAndServe(":8080", nil)
}
