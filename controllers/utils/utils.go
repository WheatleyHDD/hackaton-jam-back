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

type UserShortInfo struct {
	Email      string   `json:"email" example:"example@mail.ru" doc:"E-mail пользователя"`
	Username   string   `json:"username" example:"ThatMaidGuy" doc:"Никнейм пользователя"`
	Avatar     string   `json:"avatar" example:"http://example.com/avatar.jpg" doc:"Аватар пользователя"`
	FirstName  string   `json:"first_name" example:"Иван" doc:"Имя пользователя"`
	LastName   string   `json:"last_name" example:"Иванов" doc:"Фамилия пользователя"`
	MiddleName string   `json:"middle_name" example:"Иванович" doc:"Отчество пользователя"`
	Location   string   `json:"location" example:"г. Екатеринбург" doc:"Место жительства пользователя"`
	Skills     []string `json:"skills" doc:"Навыки пользователя"`
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
		return nil, huma.Error403Forbidden("Имя пользователя недействительное")
	}
	if userdata.Email == "" {
		return nil, huma.Error403Forbidden("Нет доступа")
	}
	return userdata, nil
}

func GetUserUsernameByEmail(email string, db *sql.DB) (*UserEmail, error) {
	row := db.QueryRow("SELECT email, username, perms FROM users WHERE email = $1", email)

	userdata := new(UserEmail)
	err := row.Scan(&userdata.Email, &userdata.Username, &userdata.Perms)
	if err != nil {
		return nil, huma.Error403Forbidden("E-mail недействительный")
	}
	if userdata.Email == "" {
		return nil, huma.Error403Forbidden("Нет доступа")
	}
	return userdata, nil
}

func NilToEmpty(val sql.NullString) string {
	return val.String
}

func GetUserShortInfo(email string, db *sql.DB) (*UserShortInfo, error) {
	result := new(UserShortInfo)

	var middle sql.NullString
	var loca sql.NullString
	if err := db.QueryRow(
		"SELECT email, username, avatar, first_name, last_name, middle_name, loc FROM users WHERE email = $1",
		email).Scan(
		&result.Email,
		&result.Username,
		&result.Avatar,
		&result.FirstName,
		&result.LastName,
		&middle,
		&loca,
	); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	result.MiddleName = middle.String
	result.Location = loca.String

	rows, err := db.Query("SELECT skill FROM skills WHERE user_email = $1", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []string

	for rows.Next() {
		var skill string
		if err := rows.Scan(&skill); err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	result.Skills = skills

	return result, nil
}
