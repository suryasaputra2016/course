package middleware

import (
	"log"
	"net/http"

	"github.com/suryasaputra2016/course/repo"
	"github.com/suryasaputra2016/course/utils"
)

type AuthMid struct {
	SessionRepo *repo.SessionRepo
}

func NewAuthMid(sr *repo.SessionRepo) *AuthMid {
	return &AuthMid{SessionRepo: sr}
}

// AccountOnly check if the user is logged in
func (am AuthMid) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			log.Printf("getting token from cookie: %s", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if token.Value == "" {
			log.Printf("token map empty")
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		tokenHashString := utils.HashToken(token.Value)

		_, err = am.SessionRepo.GetFromTokenHash(tokenHashString)
		if err != nil {
			log.Printf("session hash not found: %s", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
