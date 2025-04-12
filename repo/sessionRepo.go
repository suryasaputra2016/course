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

func (sr SessionRepo) Create(sPtr *model.Session) error {
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

func (sr SessionRepo) GetFromTokenHash(tokenHash string) (*model.Session, error) {
	var session model.Session
	queryStr := `
		SELECT id, user_id
		FROM sessions
		WHERE token_hash = $1;`
	row := sr.db.QueryRow(queryStr, tokenHash)
	err := row.Scan(&session.ID, &session.UserID)
	if err != nil {
		return nil, fmt.Errorf("selecting session: %w", err)
	}
	session.TokenHash = tokenHash
	return &session, nil
}

func (sr SessionRepo) DeleteFromTokenHash(tokenHash string) error {
	queryStr := `
		DELETE FROM sessions
			WHERE token_hash = $1`
	res, err := sr.db.Exec(queryStr, tokenHash)
	if err != nil {
		return fmt.Errorf("deleting session: %w", err)
	}
	deletedRow, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking deleted row: %w", err)
	}
	if deletedRow == 0 {
		return fmt.Errorf("zero deleted row")
	}
	return nil
}
