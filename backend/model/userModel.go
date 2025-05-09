package model

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	IsVerified   bool   `json:"is_verified"`
	Role         string `json:"role"`
}

type RegisterUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
