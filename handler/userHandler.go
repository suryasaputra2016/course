package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/suryasaputra2016/course/model"
	"github.com/suryasaputra2016/course/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	ur *repo.UserRepo
	sr *repo.SessionRepo
}

func NewUserHandler(ur *repo.UserRepo, sr *repo.SessionRepo) *UserHandler {
	return &UserHandler{
		ur: ur,
		sr: sr,
	}
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

	if ru.Email == "" || ru.Password == "" {
		log.Printf("empty email or password")
		http.Error(w, "email or pasword is empty", http.StatusBadRequest)
		return
	}

	_, err = uh.ur.GetUserByEmail(ru.Email)
	if err == nil {
		log.Printf("email used")
		http.Error(w, "email is already in used", http.StatusBadRequest)
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

func (uh UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var lu model.RegisterUser
	err := json.NewDecoder(r.Body).Decode(&lu)
	if err != nil {
		log.Printf("decoding login user: %s", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if lu.Email == "" || lu.Password == "" {
		log.Printf("empty email or password: %s", err)
		http.Error(w, "email or pasword is empty", http.StatusBadRequest)
		return
	}

	u, err := uh.ur.GetUserByEmail(lu.Email)
	if err != nil {
		log.Printf("email not found: %s", err)
		http.Error(w, "email not found", http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(lu.Password))
	if err != nil {
		log.Printf("password not match: %s", err)
		http.Error(w, "password doesn't match", http.StatusNotFound)
		return
	}

	length := 32
	b := make([]byte, length)
	nRead, err := rand.Read(b)
	if err != nil {
		log.Printf("creating random bytes: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if nRead < length {
		log.Printf("not enough read bytes: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	token := base64.URLEncoding.EncodeToString(b)
	tokenHash := sha256.Sum256([]byte(token))
	tokenHashString := base64.URLEncoding.EncodeToString(tokenHash[:])

	ns := model.Session{
		UserID:    u.ID,
		TokenHash: tokenHashString,
	}
	err = uh.sr.CreateSession(&ns)
	if err != nil {
		log.Printf("creating session: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"token": token})
	if err != nil {
		log.Printf("encoding user: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
