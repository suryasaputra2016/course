package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/suryasaputra2016/backend/course/model"
	"github.com/suryasaputra2016/backend/course/repo"
	"github.com/suryasaputra2016/backend/course/utils"
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

	_, err = uh.ur.GetByEmail(regUser.Email)
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

	err = uh.ur.Create(&newUser)
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

	user, err := uh.ur.GetByEmail(loginUser.Email)
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

	token, err := utils.GenerateToken(32)
	if err != nil {
		log.Printf("generating token: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	tokenHashString := utils.HashToken(token)

	newSession := model.Session{
		UserID:    user.ID,
		TokenHash: tokenHashString,
	}
	err = uh.sr.Create(&newSession)
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

func (uh UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	userIDString := r.PathValue("userid")
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		log.Printf("verifying email, user id not found: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	_, err = uh.ur.GetByID(userID)
	if err != nil {
		log.Printf("user id not found: %s", err)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	err = uh.ur.UpdateEmailVerification(userID)
	if err != nil {
		log.Printf("verifying email: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"message": "email verification success"})
	if err != nil {
		log.Printf("encoding user: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (uh UserHandler) CheckLoginUser(w http.ResponseWriter, r *http.Request) {
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

	tokenHashString := utils.HashToken(tokenMap["token"])

	session, err := uh.sr.GetFromTokenHash(tokenHashString)
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

	tokenHashString := utils.HashToken(tokenMap["token"])

	err = uh.sr.DeleteFromTokenHash(tokenHashString)
	if err != nil {
		log.Printf("deleting session from handler: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]string{"message": "log out sucessful"})
	if err != nil {
		log.Printf("marshaling data to json: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (uh UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
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

	user, err := uh.ur.GetByEmail(email)
	if err != nil {
		log.Printf("email not found: %s", err)
		http.Error(w, "email not found", http.StatusNotFound)
		return
	}

	token, err := utils.GenerateToken(32)
	if err != nil {
		log.Printf("generating token: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	tokenHashString := utils.HashToken(token)

	newPasswordReset := model.PasswordReset{
		UserID:         user.ID,
		TokenHash:      tokenHashString,
		ExpirationTime: time.Now().Local().Add(5 * time.Minute),
	}
	uh.prr.Create(&newPasswordReset)

	// delete previous reset database if any

	err = utils.SendPasswordResetEmail(email, token)
	if err != nil {
		log.Printf("sending email from handler: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]string{"message": "reset email sent"})
	if err != nil {
		log.Printf("marshaling data to json: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (uh UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var passChange model.PasswordChange
	err := json.NewDecoder(r.Body).Decode(&passChange)
	if err != nil {
		log.Printf("decoding token: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	token := passChange.Token
	if token == "" {
		log.Printf("empty token")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	tokenHashString := utils.HashToken(token)

	passResetPtr, err := uh.prr.GetFromTokenHash(tokenHashString)
	if err != nil {
		log.Printf("getting password reset from handler: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// check expiration date
	expired := time.Now().After(passResetPtr.ExpirationTime)
	if expired {
		log.Printf("password reset expired")
		http.Error(w, "password reset link expired", http.StatusInternalServerError)
		return
	}

	// check if new password is the same as the old one

	user, err := uh.ur.GetByID(passResetPtr.UserID)
	if err != nil {
		log.Printf("getting user from handler: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(passChange.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("hashing password: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = uh.ur.UpdatePassword(user.ID, string(passwordHash))
	if err != nil {
		log.Printf("updating user password: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// delete reset database

	response, _ := json.Marshal(map[string]string{"message": "password changed success"})
	w.Write(response)
}
