package repo

import (
	"database/sql"
	"fmt"

	"github.com/suryasaputra2016/course/model"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (ur UserRepo) CreateUser(uPtr *model.User) error {
	queryStr := `
	INSERT INTO users (email, passwordHash, role)
	VALUES ($1, $2, $3)
	RETURNING id;`
	row := ur.DB.QueryRow(queryStr, uPtr.Email, uPtr.PasswordHash, uPtr.Role)
	err := row.Scan(&uPtr.ID)
	if err != nil {
		return fmt.Errorf("creating user in repo: %w", err)
	}
	return nil
}
