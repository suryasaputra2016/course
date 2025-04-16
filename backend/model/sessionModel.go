package model

type Session struct {
	ID        int    `json:"-"`
	UserID    int    `json:"user_id"`
	TokenHash string `json:"token_hash"`
}
