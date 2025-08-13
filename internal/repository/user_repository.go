package repository

import (
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(username, password string) error {
	_, err := r.db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
	return err
}

func (r *UserRepository) GetUserPassword(username string) (string, error) {
	var password string
	err := r.db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&password)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return password, err
}
