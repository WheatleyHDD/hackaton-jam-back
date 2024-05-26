package utils

import (
	"database/sql"

	"github.com/danielgtaylor/huma/v2"
)

type JustAccessTokenInput struct {
	Body struct {
		Token string `json:"access_token" example:"82a3682d0d56f40a4d088aee08521663" doc:"Токен пользователя"`
	}
}

type UserEmail struct {
	Email    string
	Username string
	Perms    int
}

func GetUserEmailByToken(token string, db *sql.DB) (*UserEmail, error) {
	row := db.QueryRow("SELECT email, username, perms FROM users JOIN tokens ON tokens.user_email = users.email WHERE token = $1", token)

	userdata := new(UserEmail)
	err := row.Scan(&userdata.Email, &userdata.Username, &userdata.Perms)
	if err != nil {
		return nil, huma.Error403Forbidden("Токен недействительный")
	}
	if userdata.Email == "" {
		return nil, huma.Error403Forbidden("Нет доступа")
	}
	return userdata, nil
}

func GetUserEmailByUsername(username string, db *sql.DB) (*UserEmail, error) {
	row := db.QueryRow("SELECT email, username, perms FROM users WHERE username = $1", username)

	userdata := new(UserEmail)
	err := row.Scan(&userdata.Email, &userdata.Username, &userdata.Perms)
	if err != nil {
		return nil, huma.Error403Forbidden("Токен недействительный")
	}
	if userdata.Email == "" {
		return nil, huma.Error403Forbidden("Нет доступа")
	}
	return userdata, nil
}

func NilToEmpty(val sql.NullString) string {
	return val.String
}
