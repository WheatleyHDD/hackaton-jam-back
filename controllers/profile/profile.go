package profile

import (
	"database/sql"
	"hackaton-jam-back/controllers/utils"
	"log"

	"github.com/danielgtaylor/huma/v2"
)

type EditProfileInput struct {
	Body struct {
		Token string `json:"access_token" example:"82a3682d0d56f40a4d088aee08521663" doc:"Токен пользователя"`

		Username   string `json:"username,omitempty" example:"ThatMaidGuy" doc:"Никнейм пользователя"`
		Avatar     string `json:"avatar,omitempty" example:"http://example.com/avatar.jpg" doc:"Аватар пользователя"`
		FirstName  string `json:"first_name,omitempty" example:"Иван" doc:"Имя пользователя"`
		LastName   string `json:"last_name,omitempty" example:"Иванов" doc:"Фамилия пользователя"`
		MiddleName string `json:"middle_name,omitempty" example:"Иванович" doc:"Отчество пользователя"`
		About      string `json:"about,omitempty" example:"" doc:"Описание пользователя"`
		WorkPlace  string `json:"work_place,omitempty" example:"IT" doc:"Место работы"`
		WorkTime   string `json:"work_time,omitempty" example:"2 месяца" doc:"Опыт работы"`
		Location   string `json:"location,omitempty" example:"Екатеринбург" doc:"Место жительства"`
	}
}

type ProfileOutput struct {
	Body struct {
		Email       string `json:"email" example:"example@mail.ru" doc:"E-mail пользователя"`
		Username    string `json:"username" example:"ThatMaidGuy" doc:"Никнейм пользователя"`
		Avatar      string `json:"avatar" example:"http://example.com/avatar.jpg" doc:"Аватар пользователя"`
		FirstName   string `json:"first_name" example:"Иван" doc:"Имя пользователя"`
		LastName    string `json:"last_name" example:"Иванов" doc:"Фамилия пользователя"`
		MiddleName  string `json:"middle_name" example:"Иванович" doc:"Отчество пользователя"`
		About       string `json:"about" example:"" doc:"Описание пользователя"`
		WorkPlace   string `json:"work_place" example:"IT" doc:"Место работы"`
		WorkTime    string `json:"work_time" example:"2 месяца" doc:"Опыт работы"`
		Location    string `json:"location" example:"Екатеринбург" doc:"Место жительства"`
		Permissions int    `json:"permissions" example:"0" doc:"Права доступа"`
	}
}

func GetCurrentProfile(input *utils.JustAccessTokenInput, db *sql.DB) (*ProfileOutput, error) {
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil {
		return nil, err
	}

	return GetProfile(user.Username, db)
}

func GetProfile(username string, db *sql.DB) (*ProfileOutput, error) {
	row := db.QueryRow("SELECT email, username, avatar, first_name, last_name, middle_name, about, work_place, work_time, loc, perms FROM users WHERE username = $1", username)

	var avatar sql.NullString
	var middleName sql.NullString
	var about sql.NullString
	var workPlace sql.NullString
	var workTime sql.NullString
	var location sql.NullString

	userdata := new(ProfileOutput)
	err := row.Scan(
		&userdata.Body.Email,
		&userdata.Body.Username,
		&avatar,
		&userdata.Body.FirstName,
		&userdata.Body.LastName,
		&middleName,
		&about,
		&workPlace,
		&workTime,
		&location,
		&userdata.Body.Permissions)
	if err != nil {
		log.Println(err.Error())
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}

	userdata.Body.Avatar = avatar.String
	userdata.Body.MiddleName = middleName.String
	userdata.Body.About = about.String
	userdata.Body.WorkPlace = workPlace.String
	userdata.Body.WorkTime = workTime.String
	userdata.Body.Location = location.String

	return userdata, nil
}

func EditProfile(input *EditProfileInput, db *sql.DB) (*ProfileOutput, error) {
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil {
		return nil, err
	}

	// Никнейм
	if input.Body.Username != "" {
		err := db.QueryRow("UPDATE users SET username = $2 WHERE email = $1", user.Email, input.Body.Username).Scan()
		if err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}
	}

	// Аватар
	if input.Body.Avatar != "" {
		err := db.QueryRow("UPDATE users SET avatar = $2 WHERE email = $1", user.Email, input.Body.Avatar).Scan()
		if err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}
	}

	// Имя
	if input.Body.FirstName != "" {
		err := db.QueryRow("UPDATE users SET first_name = $2 WHERE email = $1", user.Email, input.Body.FirstName).Scan()
		if err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}
	}

	// Фамилия
	if input.Body.LastName != "" {
		err := db.QueryRow("UPDATE users SET last_name = $2 WHERE email = $1", user.Email, input.Body.LastName).Scan()
		if err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}
	}

	// Отчество
	if input.Body.Username != "" {
		err := db.QueryRow("UPDATE users SET username = $2 WHERE email = $1", user.Email, input.Body.Username).Scan()
		if err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}
	}

	// Аватар
	if input.Body.Avatar != "" {
		err := db.QueryRow("UPDATE users SET avatar = $2 WHERE email = $1", user.Email, input.Body.Avatar).Scan()
		if err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}
	}

	// Имя
	if input.Body.FirstName != "" {
		err := db.QueryRow("UPDATE users SET first_name = $2 WHERE email = $1", user.Email, input.Body.FirstName).Scan()
		if err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}
	}

	// Фамилия
	if input.Body.LastName != "" {
		err := db.QueryRow("UPDATE users SET last_name = $2 WHERE email = $1", user.Email, input.Body.LastName).Scan()
		if err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}
	}

	return GetProfile(user.Username, db)
}
