package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"sun-position/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10040"
	}

	// Register routes with /sun-pos prefix - order matters!
	// Static files first
	http.Handle("/sun-pos/static/", http.StripPrefix("/sun-pos/static/", handlers.StaticFileServer()))

	// Then API routes
	http.HandleFunc("/sun-pos/api/sun-position", handlers.SunPositionHandler)

	// Finally, catch-all for the home page (should be last)
	http.HandleFunc("/sun-pos/", handlers.HomeHandler)

	fmt.Printf("Server starting on http://localhost:%s/sun-pos\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
