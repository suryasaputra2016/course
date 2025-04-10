package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/suryasaputra2016/course/config"
	"github.com/suryasaputra2016/course/handler"
	"github.com/suryasaputra2016/course/repo"
)

func main() {
	// set up postgres database
	db, err := config.ConnectPostgres()
	if err != nil {
		log.Fatal(fmt.Errorf("connecting database from main: %w", err))
	}
	defer config.ClosePostgres(db)
	fmt.Println("postgres database connected.")

	// migration
	err = config.PrepareTables(db)
	if err != nil {
		log.Fatal(fmt.Errorf("preparing table from main: %w", err))
	}

	// repos and handlers
	ur := repo.NewUserRepo(db)
	sr := repo.NewSessionRepo(db)
	uh := handler.NewUserHandler(ur, sr)

	// define mux
	mux := http.NewServeMux()

	// define routes
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("POST /register", uh.RegisterUser)
	mux.HandleFunc("POST /login", uh.LoginUser)
	mux.HandleFunc("DELETE /logout", uh.LogoutUser)
	mux.HandleFunc("GET /checklogin", uh.CheckLoginUser)

	// serving and listening
	fmt.Printf("serving and listening on :8080...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(fmt.Errorf("listening and serving: %w", err))
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "8080 served")
}
