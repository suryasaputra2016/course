package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/suryasaputra2016/course/frontend/handler"
	"github.com/suryasaputra2016/course/frontend/templates"
)

func main() {
	tmpl, err := template.ParseFS(templates.FS, "home.html")
	if err != nil {
		panic("cannot parse files")
	}
	homeHandler := handler.NewHomeHandler(tmpl)

	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler.ShowHome)

	server := http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	fmt.Println("serving and listening front-end on :8081...")
	log.Fatal(server.ListenAndServe())
}
