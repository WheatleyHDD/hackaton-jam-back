package events

import (
	"database/sql"
	"hackaton-jam-back/controllers/utils"

	"github.com/danielgtaylor/huma/v2"
)

type EventJoinExitInput struct {
	Urid string `path:"urid" maxLength:"30" example:"example_events" doc:"Ссылка на мероприятие"`
	Body struct {
		Token string `json:"access_token" example:"82a3682d0d56f40a4d088aee08521663" doc:"Токен пользователя"`
	}
}

type EventJoinExitOutput struct {
	Success bool `json:"success" example:"true" doc:"Успех выполнения"`
}

type UserEventsOutput struct {
	Body struct {
		Events []EventType `json:"event" doc:"Лист с событиями"`
	}
}

func JoinEvent(input *EventJoinExitInput, db *sql.DB) (*EventJoinExitOutput, error) {
	if err := db.QueryRow("SELECT urid FROM events WHERE urid = $1", input.Urid).Scan(); err != nil {
		if err == sql.ErrNoRows {
			return nil, huma.Error403Forbidden("Этого события нет XP")
		}
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	// Проверить можем ли присоединится к меро?
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}
	if user.Perms == 0 {
		return nil, huma.Error403Forbidden("Участвовать можно только обычным пользователям")
	}

	if err := db.QueryRow("INSERT INTO event_members (event_uri, member_email) VALUES ($1, $2)", input.Urid, user.Email).Scan(); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return &EventJoinExitOutput{Success: true}, nil
}

func ExitEvent(input *EventJoinExitInput, db *sql.DB) (*EventJoinExitOutput, error) {
	if err := db.QueryRow("SELECT urid FROM events WHERE urid = $1", input.Urid).Scan(); err != nil {
		if err == sql.ErrNoRows {
			return nil, huma.Error403Forbidden("Этого события нет XP")
		}
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	// Проверить можем ли присоединится к меро?
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}
	if user.Perms == 0 {
		return nil, huma.Error403Forbidden("Участвовать можно только обычным пользователям")
	}

	// Удаляем
	if err := db.QueryRow("DELETE FROM event_members WHERE event_uri = $1 AND member_email = $2", input.Urid, user.Email).Scan(); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return &EventJoinExitOutput{Success: true}, nil
}

func GetAllJoinedEvents(email string, db *sql.DB) (*UserEventsOutput, error) {
	rows, err := db.Query(
		"SELECT events.id, events.urid, events.name, events.start_time, events.end_time, events.location, event_members.member_email"+
			"FROM event_members INNER JOIN events ON event_members.event_uri=events.urid"+
			"WHERE event_members.member_email = $1 ORDER BY events.id DESC", email,
	)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity("Проблемки с вызовом SQL")
	}
	defer rows.Close()

	result := new(UserEventsOutput)

	for rows.Next() {
		var id int
		var member string
		var event EventType
		if err := rows.Scan(&id, &event.Urid, &id, &event.Name, &event.StartTime, &event.EndTime, &event.Location, member); err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
		result.Body.Events = append(result.Body.Events, event)
	}
	if err = rows.Err(); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return result, nil
}
