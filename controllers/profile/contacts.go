package profile

import (
	"database/sql"
	"hackaton-jam-back/controllers/utils"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
)

type AddDelProfileContactInput struct {
	Username string `path:"username" maxLength:"30" example:"ThatMaidGuy" doc:"Никнейм пользователя"`
	Body     struct {
		Token string `json:"access_token" example:"82a3682d0d56f40a4d088aee08521663" doc:"Токен пользователя"`
		Link  string `json:"contact_link" example:"https://vk.com/id0" doc:"Ссылка"`
	}
}

type ProfileContactsOutput struct {
	Body struct {
		Contacts []string `json:"contacts" doc:"Список контактов пользователя"`
	}
}

/// ==============================================
/// ==============================================
/// ===== Список контактов для пользователя ======
/// ==============================================
/// ==============================================

func GetContacts(username string, db *sql.DB) (*ProfileContactsOutput, error) {
	user, err := utils.GetUserEmailByUsername(username, db)
	if err != nil {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}

	rows, err := db.Query("SELECT contact_link FROM contacts WHERE user_email = $1", user.Email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []string

	for rows.Next() {
		var link string
		if err := rows.Scan(&link); err != nil {
			return nil, err
		}
		contacts = append(contacts, link)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	result := new(ProfileContactsOutput)
	result.Body.Contacts = contacts

	return result, nil
}

func AddContact(input *AddDelProfileContactInput, db *sql.DB) (*ProfileContactsOutput, error) {
	// Проверить можем ли менять?
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}
	if input.Username != user.Username && user.Perms != 10 {
		return nil, huma.Error403Forbidden("Нет прав")
	}

	// Проверка ссылки
	_, err = url.ParseRequestURI(input.Body.Link)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity("Неверная ссылка")
	}

	// Получить уже имеющиеся ссылки
	existingContacts, err := GetContacts(user.Username, db)
	if err != nil {
		return nil, huma.Error403Forbidden(err.Error())
	}

	// Проверить на существование
	for _, c := range existingContacts.Body.Contacts {
		if c == input.Body.Link {
			return nil, huma.Error403Forbidden("Такая ссылка уже есть")
		}
	}

	// Добавить контакт
	_, err = db.Query("INSERT INTO contacts (user_email, contact_link) VALUES ($1, $2)", user.Email, input.Body.Link)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	existingContacts.Body.Contacts = append(existingContacts.Body.Contacts, input.Body.Link)

	return existingContacts, nil
}

func DelContact(input *AddDelProfileContactInput, db *sql.DB) (*ProfileContactsOutput, error) {
	// Проверить можем ли менять?
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}
	if input.Username != user.Username && user.Perms != 10 {
		return nil, huma.Error403Forbidden("Нет прав")
	}

	// Удалить ссылку
	_, err = db.Query("DELETE FROM contacts WHERE user_email=$1 AND contact_link=$2", user.Email, input.Body.Link)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity("Похоже ссылки не было")
	}

	return GetContacts(user.Username, db)
}
