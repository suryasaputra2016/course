package repo

import (
	"database/sql"
	"fmt"

	"github.com/suryasaputra2016/course/model"
)

type SessionRepo struct {
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (sr SessionRepo) CreateSession(sPtr *model.Session) error {
	queryStr := `
		INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2)
		RETURNING id;`
	row := sr.db.QueryRow(queryStr, sPtr.UserID, sPtr.TokenHash)
	err := row.Scan(&sPtr.ID)
	if err != nil {
		return fmt.Errorf("creating session in repo: %w", err)
	}
	return nil
}
