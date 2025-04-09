package model

type Session struct {
	ID        int    `json:"-"`
	UserID    int    `json:"-"`
	TokenHash string `json:"token_hash"`
}
