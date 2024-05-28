package events

import (
	"database/sql"
	"hackaton-jam-back/controllers/utils"
	"reflect"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type EventCreationInput struct {
	Body struct {
		Token string `json:"access_token" example:"82a3682d0d56f40a4d088aee08521663" doc:"Токен пользователя"`

		Urid                  string    `json:"urid" maxLength:"30" example:"example_events" doc:"Ссылка на мероприятие (Поменять потом нельзя!!!)"`
		Name                  string    `json:"name" example:"Example GameJam" doc:"Название мероприятия"`
		StartTime             time.Time `json:"start_time" doc:"Начало проведения"`
		EndTime               time.Time `json:"end_time" doc:"Конец проведения"`
		Location              string    `json:"location" example:"Свердловская область, г. Екатеринбург" doc:"Место проведения"`
		Icon                  string    `json:"icon" doc:"Превью мероприятия"`
		IsIrl                 bool      `json:"is_irl" doc:"Очное ли мероприятие?"`
		TeamRequirementsType  int       `json:"team_requirements_type" doc:"Тип равенства требования к количеству сокомандников"`
		TeamRequirementsValue int       `json:"team_requirements_value" doc:"Количество сокомандников"`

		Description  string `json:"desc,omitempty" doc:"Описание мероприятия"`
		Prize        string `json:"prize,omitempty" doc:"Призы мероприятия"`
		Partners     string `json:"partners,omitempty" doc:"Партнеры мероприятия"`
		Requirements string `json:"requirements,omitempty" doc:"Необходимые навыки для мероприятия"`
	}
}

type EventEditInput struct {
	Urid string `path:"urid" maxLength:"30" example:"example_events" doc:"Ссылка на мероприятие"`

	Body struct {
		Token string `json:"access_token" example:"82a3682d0d56f40a4d088aee08521663" doc:"Токен пользователя"`

		Name                  string    `json:"name,omitempty" example:"Example GameJam" doc:"Название мероприятия"`
		StartTime             time.Time `json:"start_time,omitempty" doc:"Начало проведения"`
		EndTime               time.Time `json:"end_time,omitempty" doc:"Конец проведения"`
		Location              string    `json:"location,omitempty" example:"Свердловская область, г. Екатеринбург" doc:"Место проведения"`
		Icon                  string    `json:"icon,omitempty" doc:"Превью мероприятия"`
		IsIrl                 bool      `json:"is_irl,omitempty" doc:"Очное ли мероприятие?"`
		TeamRequirementsType  int       `json:"team_requirements_type,omitempty" doc:"Тип равенства требования к количеству сокомандников"`
		TeamRequirementsValue int       `json:"team_requirements_value,omitempty" doc:"Количество сокомандников"`

		Description  string `json:"desc,omitempty" doc:"Описание мероприятия"`
		Prize        string `json:"prize,omitempty" doc:"Призы мероприятия"`
		Partners     string `json:"partners,omitempty" doc:"Партнеры мероприятия"`
		Requirements string `json:"requirements,omitempty" doc:"Необходимые навыки для мероприятия"`
	}
}

type Organizators struct {
	Email    string `json:"email" doc:"E-mail организатора"`
	Username string `json:"username" doc:"Имя пользователя организатора"`
}

type FullEventOutput struct {
	Body struct {
		Urid                  string    `json:"urid" example:"example_events" doc:"Ссылка на мероприятие"`
		Name                  string    `json:"name" example:"Example GameJam" doc:"Название мероприятия"`
		StartTime             time.Time `json:"start_time" doc:"Начало проведения"`
		EndTime               time.Time `json:"end_time" doc:"Конец проведения"`
		Location              string    `json:"location" example:"Свердловская область, г. Екатеринбург" doc:"Место проведения"`
		Description           string    `json:"desc" doc:"Описание мероприятия"`
		Prize                 string    `json:"prize" doc:"Призы мероприятия"`
		Requirements          string    `json:"requirements" doc:"Необходимые навыки для мероприятия"`
		Partners              string    `json:"partners" doc:"Партнеры мероприятия"`
		Icon                  string    `json:"icon" doc:"Превью мероприятия"`
		IsIrl                 bool      `json:"is_irl" doc:"Очное ли мероприятие?"`
		TeamRequirementsType  int       `json:"team_requirements_type" doc:"Тип равенства требования к количеству сокомандников"`
		TeamRequirementsValue int       `json:"team_requirements_value" doc:"Количество сокомандников"`

		Organizators []*Organizators `json:"organisators" doc:"Список организаторов"`
	}
}

func CreateEvent(input *EventCreationInput, db *sql.DB) (*FullEventOutput, error) {
	// Проверить можем ли создать меро?
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}
	if user.Perms < 1 {
		return nil, huma.Error403Forbidden("Нет прав")
	}

	// Запись в базу
	_, err = db.Query("INSERT INTO events ("+
		"urid, name, start_time, end_time, prize, location, desc, requirements,"+
		"partners, icon, is_irl, team_requirements_type, team_requirements_value)"+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",

		input.Body.Urid, input.Body.Name,
		input.Body.StartTime, input.Body.EndTime,
		input.Body.Prize, input.Body.Location,
		input.Body.Description, input.Body.Requirements,
		input.Body.Partners, input.Body.Icon, input.Body.IsIrl,
		input.Body.TeamRequirementsType, input.Body.TeamRequirementsValue,
	)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	_, err = db.Query("INSERT INTO event_orgs (event_uri, organizator_email) VALUES ($1, $2)", input.Body.Urid, user.Email)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return getFullEventInfo(input.Body.Urid, db)
}

func EditEvent(input *EventEditInput, db *sql.DB) (*FullEventOutput, error) {
	// Проверить наша ли меро?
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}

	if user.Perms != 10 {
		if err := db.QueryRow("SELECT * FROM event_orgs WHERE organizator_email = $1 AND event_uri = $2", user.Email, input.Urid).Scan(); err != nil {
			if err == sql.ErrNoRows {
				return nil, huma.Error403Forbidden("Это не твое мероприятие :/")
			}
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
	}

	val := reflect.ValueOf(input.Body)
	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		column_name := t.Field(i)
		fvalue := reflect.Indirect(val).FieldByName(column_name.Name)

		if fvalue.IsZero() {
			continue
		}

		err := db.QueryRow("UPDATE events SET "+strings.Split(column_name.Tag.Get("json"), ",")[0]+" = $2 WHERE urid = $1", user.Email, fvalue.String()).Scan()
		//                                                                               дебилизм ^                                  идиотизм ^
		if err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}
	}

	return getFullEventInfo(input.Urid, db)
}