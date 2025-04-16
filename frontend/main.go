package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)

	server := http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var data map[string]string
	json.Unmarshal(body, &data)
	json.NewEncoder(w).Encode(data)

}
