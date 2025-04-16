package model

import "time"

type PasswordReset struct {
	ID             int       `json:"-"`
	UserID         int       `json:"user_id"`
	TokenHash      string    `json:"token_hash"`
	ExpirationTime time.Time `json:"expiration_time"`
}

type PasswordChange struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}
