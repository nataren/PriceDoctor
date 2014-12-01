package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"kittens": [
      {"id": 1, "name": "Bobby", "picture": "http://placekitten.com/g/200/300"},
      {"id": 2, "name": "Wally", "picture": "http://placekitten.com/g/200/400"}
	]}`))
}

func main() {
	log.Println("Starting server")

	r := mux.NewRouter()
	r.HandleFunc("/api/kittens", SearchHandler).Methods("GET")
	http.Handle("/api/", r)

	// Server public assets
	http.Handle("/", http.FileServer(http.Dir("./public/")))

	log.Println("Listing on port 8080")
	http.ListenAndServe(":8080", nil)
}
