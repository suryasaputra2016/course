package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/suryasaputra2016/course/model"
	"github.com/suryasaputra2016/course/repo"
	"github.com/suryasaputra2016/course/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	ur  *repo.UserRepo
	sr  *repo.SessionRepo
	prr *repo.PasswordResetRepo
}

func NewUserHandler(
	ur *repo.UserRepo,
	sr *repo.SessionRepo,
	prr *repo.PasswordResetRepo,
) *UserHandler {
	return &UserHandler{
		ur:  ur,
		sr:  sr,
		prr: prr,
	}
}

func (uh UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var regUser model.RegisterUser
	err := json.NewDecoder(r.Body).Decode(&regUser)
	if err != nil {
		log.Printf("decoding register user: %s", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if regUser.Email == "" || regUser.Password == "" {
		log.Printf("empty email or password")
		http.Error(w, "email or pasword is empty", http.StatusBadRequest)
		return
	}

	_, err = uh.ur.GetUserByEmail(regUser.Email)
	if err == nil {
		log.Printf("email used")
		http.Error(w, "email is already in used", http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(regUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("hashing password: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	newUser := model.User{
		Email:        regUser.Email,
		PasswordHash: string(passwordHash),
		Role:         "user",
	}

	err = uh.ur.CreateUser(&newUser)
	if err != nil {
		log.Printf("creating user in handler: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(newUser)
	if err != nil {
		log.Printf("encoding user: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (uh UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var loginUser model.RegisterUser
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		log.Printf("decoding login user: %s", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if loginUser.Email == "" || loginUser.Password == "" {
		log.Printf("empty email or password: %s", err)
		http.Error(w, "email or pasword is empty", http.StatusBadRequest)
		return
	}

	user, err := uh.ur.GetUserByEmail(loginUser.Email)
	if err != nil {
		log.Printf("email not found: %s", err)
		http.Error(w, "email not found", http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginUser.Password))
	if err != nil {
		log.Printf("password not match: %s", err)
		http.Error(w, "password doesn't match", http.StatusNotFound)
		return
	}

	lengthByte := 32
	tokenByte := make([]byte, lengthByte)
	totalRead, err := rand.Read(tokenByte)
	if err != nil {
		log.Printf("creating random bytes: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if totalRead < lengthByte {
		log.Printf("not enough read bytes: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	token := base64.URLEncoding.EncodeToString(tokenByte)
	tokenHash := sha256.Sum256([]byte(token))
	tokenHashString := base64.URLEncoding.EncodeToString(tokenHash[:])

	newSession := model.Session{
		UserID:    user.ID,
		TokenHash: tokenHashString,
	}
	err = uh.sr.CreateSession(&newSession)
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

func (uh UserHandler) CheckLoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tokenMap map[string]string
	err := json.NewDecoder(r.Body).Decode(&tokenMap)
	if err != nil {
		log.Printf("decoding token map: %s", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if tokenMap["token"] == "" {
		log.Printf("token map empty")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	tokenHash := sha256.Sum256([]byte(tokenMap["token"]))
	tokenHashString := base64.URLEncoding.EncodeToString(tokenHash[:])

	session, err := uh.sr.GetSessionFromTokenHash(tokenHashString)
	if err != nil {
		log.Printf("session hash not found: %s", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(session)
	if err != nil {
		log.Printf("encoding session: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (uh UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tokenMap map[string]string
	err := json.NewDecoder(r.Body).Decode(&tokenMap)
	if err != nil {
		log.Printf("decoding token map: %s", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if tokenMap["token"] == "" {
		log.Printf("token map empty")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	tokenHash := sha256.Sum256([]byte(tokenMap["token"]))
	tokenHashString := base64.URLEncoding.EncodeToString(tokenHash[:])

	err = uh.sr.DeleteSessionFromTokenHash(tokenHashString)
	if err != nil {
		log.Printf("deleting session from handler: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(map[string]string{"message": "log out sucessful"})
	w.Write(response)
}

func (uh UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var emailMap map[string]string
	err := json.NewDecoder(r.Body).Decode(&emailMap)
	if err != nil {
		log.Printf("decoding email map: %s", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	email := emailMap["email"]
	if email == "" {
		log.Printf("email map empty")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = utils.CheckEmailFormat(email)
	if err != nil {
		log.Printf("email not well formatted: %s", err)
		http.Error(w, "bad request", http.StatusNotFound)
		return
	}

	user, err := uh.ur.GetUserByEmail(email)
	if err != nil {
		log.Printf("email not found: %s", err)
		http.Error(w, "email not found", http.StatusNotFound)
		return
	}

	lengthByte := 32
	tokenByte := make([]byte, lengthByte)
	totalRead, err := rand.Read(tokenByte)
	if err != nil {
		log.Printf("creating random bytes: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if totalRead < lengthByte {
		log.Printf("not enough read bytes: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	token := base64.URLEncoding.EncodeToString(tokenByte)
	tokenHash := sha256.Sum256([]byte(token))
	tokenHashString := base64.URLEncoding.EncodeToString(tokenHash[:])

	newPasswordReset := model.PasswordReset{
		UserID:         user.ID,
		TokenHash:      tokenHashString,
		ExpirationTime: time.Now().Local().Add(5 * time.Minute),
	}
	uh.prr.Create(&newPasswordReset)

	err = utils.SendPasswordResetEmail(email, token)
	if err != nil {
		log.Printf("sending email from handler: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(map[string]string{"message": "reset email sent"})
	w.Write(response)
}
