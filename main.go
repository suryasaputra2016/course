package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// define mux
	mux := http.NewServeMux()

	// define routes
	mux.HandleFunc("/", homeHandler)

	// serving and listening
	fmt.Printf("serving and listening on :8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(fmt.Errorf("listening and serving: %w", err))
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "8080 served")
}
