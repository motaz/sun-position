package main

import (
	"log"
	"net/http"
)

func main() {
	// Define the route for the home page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})

	// Serve static files (CSS, JS, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", StaticFileServer()))

	// Start the server
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}