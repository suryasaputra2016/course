package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/suryasaputra2016/course/backend/config"
	"github.com/suryasaputra2016/course/backend/handler"
	"github.com/suryasaputra2016/course/backend/middleware"
	"github.com/suryasaputra2016/course/backend/repo"
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
	prr := repo.NewPasswordResetRepo(db)
	uh := handler.NewUserHandler(ur, sr, prr)
	nfh := handler.NewNotFoundHandler()

	// define routes
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", uh.RegisterUser)
	mux.HandleFunc("POST /login", uh.LoginUser)
	mux.HandleFunc("PUT /verifyemail/{userid}", uh.VerifyEmail)
	mux.HandleFunc("PUT /updatepassword", uh.UpdatePassword)
	mux.HandleFunc("GET /checklogin", uh.CheckLoginUser)

	accountMux := http.NewServeMux()
	accountMux.HandleFunc("DELETE /logout", uh.LogoutUser)
	accountMux.HandleFunc("POST /resetpassword", uh.ResetPassword)

	publicMux := http.NewServeMux()
	publicMux.HandleFunc("/", nfh.Home)

	auth := middleware.NewAuthMid(sr)
	mux.Handle("/dashboard/", http.StripPrefix("/dashboard", auth.Authorize(accountMux)))
	mux.Handle("/", http.StripPrefix("", publicMux))

	// serving and listening
	server := http.Server{
		Addr:    ":8080",
		Handler: middleware.SetJSONHeader(mux),
	}
	fmt.Println("serving and listening back-end on :8080...")
	log.Fatal(server.ListenAndServe())
}
