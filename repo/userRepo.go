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

func (ur UserRepo) GetUserByEmail(email string) (*model.User, error) {
	var u model.User
	queryStr := `
	SELECT  id, passwordhash, role FROM users
	WHERE email = $1;`
	row := ur.DB.QueryRow(queryStr, email)
	err := row.Scan(&u.ID, &u.PasswordHash, &u.Role)
	if err != nil {
		return nil, fmt.Errorf("selecting user by email in repo: %w", err)
	}
	u.Email = email
	return &u, nil
}
