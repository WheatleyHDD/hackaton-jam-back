package utils

import (
	"database/sql"

	"github.com/danielgtaylor/huma/v2"
)

type UserEmail struct {
	Email string
	Perms int
}

func GetUserEmailByToken(token string, db *sql.DB) (*UserEmail, error) {
	row := db.QueryRow("SELECT email, perms FROM users JOIN tokens ON tokens.user_email = users.email WHERE token = $1", token)

	userdata := new(UserEmail)
	err := row.Scan(&userdata.Email, &userdata.Perms)
	if err != nil {
		return nil, huma.Error403Forbidden("Токен недействительный")
	}
	if userdata.Email == "" {
		return nil, huma.Error403Forbidden("Нет доступа")
	}
	return userdata, nil
}
