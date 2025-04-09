package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/suryasaputra2016/course/model"
	"github.com/suryasaputra2016/course/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	ur *repo.UserRepo
}

func NewUserHandler(ur *repo.UserRepo) *UserHandler {
	return &UserHandler{ur: ur}
}

func (uh UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ru model.RegisterUser
	err := json.NewDecoder(r.Body).Decode(&ru)
	if err != nil {
		log.Printf("decoding register user: %s", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	_, err = uh.ur.GetUserIDByEmail(ru.Email)
	if err == nil {
		log.Printf("email used")
		http.Error(w, "email is already in used", http.StatusBadRequest)
		return
	}

	if ru.Email == "" || ru.Password == "" {
		log.Printf("empty email or password")
		http.Error(w, "email or pasword is empty", http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(ru.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("hashing password: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	u := model.User{
		Email:        ru.Email,
		PasswordHash: string(passwordHash),
		Role:         "user",
	}

	err = uh.ur.CreateUser(&u)
	if err != nil {
		log.Printf("creating user in handler: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		log.Printf("encoding user: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
