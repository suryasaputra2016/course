package repo

import (
	"database/sql"
	"fmt"

	"github.com/suryasaputra2016/course/model"
)

type PasswordResetRepo struct {
	db *sql.DB
}

func NewPasswordResetRepo(db *sql.DB) *PasswordResetRepo {
	return &PasswordResetRepo{
		db: db,
	}
}

func (prr PasswordResetRepo) Create(prPtr *model.PasswordReset) error {
	queryStr := `
	INSERT INTO password_resets (user_id, token_hash, expiration_time)
	VALUES ($1, $2, $3)
	RETURNING id;`
	row := prr.db.QueryRow(queryStr, prPtr.UserID, prPtr.TokenHash, prPtr.ExpirationTime)
	err := row.Scan(&prPtr.ID)
	if err != nil {
		return fmt.Errorf("creating password reset in repo: %w", err)
	}
	return nil
}

func (prr PasswordResetRepo) GetFromTokenHash(tokenHash string) (*model.PasswordReset, error) {
	passReset := model.PasswordReset{
		TokenHash: tokenHash,
	}
	queryStr := `
	SELECT id, user_id, expiration_time
	FROM password_resets
	WHERE token_hash = $1;`
	row := prr.db.QueryRow(queryStr, tokenHash)
	err := row.Scan(&passReset.ID, &passReset.UserID, &passReset.ExpirationTime)
	if err != nil {
		return nil, fmt.Errorf("getting pasword reset from repo: %w", err)
	}

	return &passReset, nil
}
