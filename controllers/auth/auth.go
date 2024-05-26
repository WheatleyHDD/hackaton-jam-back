package auth

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"hackaton-jam-back/controllers/utils"
	"log"
	"strconv"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/crypto/bcrypt"
)

// ==========================
// ======= Структуры ========
// ==========================
type LoginInput struct {
	Body struct {
		Email    string `json:"email" example:"thatmaidguy@ya.ru" doc:"E-mail пользователя"`
		Password string `json:"password" example:"qwerty123" doc:"Пароль пользователя"`
	}
}

type RegisterInput struct {
	Body struct {
		Email         string `json:"email" example:"thatmaidguy@ya.ru" doc:"E-mail пользователя"`
		Username      string `json:"username" example:"ThatMaidGuy" doc:"Никнейм пользователя"`
		FirstName     string `json:"first_name" example:"Иван" doc:"Имя пользователя"`
		LastName      string `json:"last_name" example:"Иванов" doc:"Фамилия пользователя"`
		Password      string `json:"password" example:"qwerty123" doc:"Пароль пользователя"`
		IsOrganisator bool   `json:"is_organisator" doc:"Профиль для организатора?"`
	}
}

type LoginResponseOutput struct {
	Body struct {
		Email    string `json:"email" example:"example@mail.ru" doc:"E-mail пользователя"`
		Username string `json:"username" example:"example@mail.ru" doc:"Никнейм пользователя"`
		Token    string `json:"access_token" example:"82a3682d0d56f40a4d088aee08521663" doc:"Токен для доступа"`
	}
}

// ==========================
// ======== Методы ==========
// ==========================

func Register(input *RegisterInput, db *sql.DB) (*LoginResponseOutput, error) {
	// Проверка на существование пользователя
	rows, err := db.Query("SELECT COUNT(*) AS count FROM users WHERE email = $1", input.Body.Email)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	defer rows.Close()
	count, err := checkCount(rows)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	if count > 0 {
		return nil, huma.Error422UnprocessableEntity("Пользователь с таким Email уже существует")
	}

	// Хэширование пароля
	saltedBytes := []byte(input.Body.Password)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return nil, huma.Error422UnprocessableEntity("Ошибка с авторизацией на стороне сервера 553")
	}
	hash := string(hashedBytes[:])

	perm := 0
	if input.Body.IsOrganisator {
		perm = 1
	}

	// Запись в базу
	_, err = db.Query("INSERT INTO users (email, username, first_name, last_name, password, perms) VALUES ($1, $2, $3, $4, $5, $6)",
		input.Body.Email, input.Body.Username, input.Body.FirstName, input.Body.LastName, hash, perm)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	// Сразу входим
	logindata := new(LoginInput)
	logindata.Body.Email = input.Body.Email
	logindata.Body.Password = input.Body.Password
	return Login(logindata, db)
}

func Login(input *LoginInput, db *sql.DB) (*LoginResponseOutput, error) {
	// Проверка на существование пользователя
	rows, err := db.Query("SELECT COUNT(*) AS count FROM users WHERE email = $1", input.Body.Email)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	defer rows.Close()

	count, err := checkCount(rows)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	if count == 0 {
		return nil, huma.Error422UnprocessableEntity("Пользователя с таким Email не существует")
	}

	// Получаем данные и проверяем пароли
	row := db.QueryRow("SELECT email, username, password FROM users WHERE email = $1", input.Body.Email)

	var email string
	var username string
	var hashed_pass string
	err = row.Scan(&email, &username, &hashed_pass)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	incoming := []byte(input.Body.Password)
	existing := []byte(hashed_pass)
	err = bcrypt.CompareHashAndPassword(existing, incoming)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity("Пароль неверный")
	}

	// Создаем access_token
	timestamp := time.Now().Unix()

	hasher := md5.New()
	_, err = hasher.Write(existing)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	_, err = hasher.Write([]byte(strconv.Itoa(int(timestamp))))
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	token := hex.EncodeToString(hasher.Sum(nil))

	_, err = db.Query("INSERT INTO tokens (user_email, token) VALUES ($1, $2)", email, token)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	// Пишем ответ
	resp := &LoginResponseOutput{}

	resp.Body.Email = email
	resp.Body.Username = username
	resp.Body.Token = token

	return resp, nil
}

func Logout(input *utils.JustAccessTokenInput, db *sql.DB) (*struct{}, error) {
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil {
		return nil, err
	}

	_, err = db.Query("DELETE FROM tokens WHERE user_email = $1 AND token = $2", user.Email, input.Body.Token)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return nil, nil
}

func checkCount(rows *sql.Rows) (count int, err error) {
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}
