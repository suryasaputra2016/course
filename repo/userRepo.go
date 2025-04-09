package repo

import (
	"database/sql"
	"fmt"

	"github.com/suryasaputra2016/course/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (ur UserRepo) CreateUser(userPtr *model.User) error {
	queryStr := `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id;`
	row := ur.db.QueryRow(queryStr, userPtr.Email, userPtr.PasswordHash, userPtr.Role)
	err := row.Scan(&userPtr.ID)
	if err != nil {
		return fmt.Errorf("creating user in repo: %w", err)
	}
	return nil
}

func (ur UserRepo) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	queryStr := `
		SELECT  id, password_hash, role FROM users
		WHERE email = $1;`
	row := ur.db.QueryRow(queryStr, email)
	err := row.Scan(&user.ID, &user.PasswordHash, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("selecting user by email in repo: %w", err)
	}
	user.Email = email
	return &user, nil
}
