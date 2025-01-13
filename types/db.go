package types

import (
	"database/sql"
	"errors"
	"log/slog"
)

type User struct {
	ID       string
	Email    string
	Password string
}

func (u *User) GetUserByEmail(db *sql.DB, email string) error {
	err := db.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).Scan(&u.ID, &u.Email, &u.Password)
	slog.Debug("GetUserByEmail - User found", "user", u)
	if err == sql.ErrNoRows {
		return errors.New("user not found")
	}
	return nil
}
