package repo

import (
	"database/sql"
	"fmt"

	"github.com/suryasaputra2016/course/backend/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (ur UserRepo) Create(userPtr *model.User) error {
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

func (ur UserRepo) GetByEmail(email string) (*model.User, error) {
	user := model.User{Email: email}
	queryStr := `
		SELECT  id, password_hash, is_verified, role FROM users
		WHERE email = $1;`
	row := ur.db.QueryRow(queryStr, email)
	err := row.Scan(&user.ID, &user.PasswordHash, &user.IsVerified, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("selecting user by email in repo: %w", err)
	}
	return &user, nil
}

func (ur UserRepo) GetByID(id int) (*model.User, error) {
	user := model.User{ID: id}
	queryStr := `
		SELECT  email, password_hash, is_verified, role 
		FROM users
		WHERE id = $1;`
	row := ur.db.QueryRow(queryStr, id)
	err := row.Scan(&user.Email, &user.PasswordHash, &user.IsVerified, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("selecting user by id in repo: %w", err)
	}
	return &user, nil
}

func (ur UserRepo) UpdatePassword(id int, newPasswordHash string) error {
	queryStr := `
	UPDATE users
	SET password_hash = $1
	WHERE id = $2;`
	_, err := ur.db.Exec(queryStr, newPasswordHash, id)
	if err != nil {
		return fmt.Errorf("updating user password in repo: %w", err)
	}
	return nil
}

func (ur UserRepo) UpdateEmailVerification(id int) error {
	queryStr := `
	UPDATE users
	SET is_verified = TRUE
	WHERE id = $1;`
	_, err := ur.db.Exec(queryStr, id)
	if err != nil {
		return fmt.Errorf("updating user password in repo: %w", err)
	}
	return nil
}
