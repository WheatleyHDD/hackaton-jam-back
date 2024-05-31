package events

import (
	"database/sql"
	"hackaton-jam-back/controllers/utils"
	"strconv"

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

type EventSearchUsers struct {
	Urid string `path:"urid" maxLength:"30" example:"example_events" doc:"Ссылка на мероприятие"`
	Body struct {
		SkillsToSearch []string `json:"skills_to_search" doc:"Токен пользователя"`
	}
}

type EventSearchUsersOutput struct {
	Body struct {
		Users []*utils.UserShortInfo `json:"users" doc:"Таблица подходящих пользователей по критериям (если нет - выводит всех)"`
	}
}

type UserEventsOutput struct {
	Body struct {
		Events []EventType `json:"event" doc:"Лист с событиями"`
	}
}

func JoinEvent(input *EventJoinExitInput, db *sql.DB) (*EventJoinExitOutput, error) {
	if err := isEventExists(input.Urid, db); err != nil {
		return nil, err
	}

	// Проверить можем ли присоединится к меро?
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}
	if user.Perms != 0 {
		return nil, huma.Error403Forbidden("Участвовать можно только обычным пользователям")
	}

	db.QueryRow("INSERT INTO event_members (event_uri, member_email) VALUES ($1, $2)", input.Urid, user.Email).Scan()

	return &EventJoinExitOutput{Success: true}, nil
}

func ExitEvent(input *EventJoinExitInput, db *sql.DB) (*EventJoinExitOutput, error) {
	if err := isEventExists(input.Urid, db); err != nil {
		return nil, err
	}

	// Проверить можем ли присоединится к меро?
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}
	if user.Perms != 0 {
		return nil, huma.Error403Forbidden("Участвовать можно только обычным пользователям")
	}

	// Удаляем
	db.QueryRow("DELETE FROM event_members WHERE event_uri = $1 AND member_email = $2", input.Urid, user.Email).Scan()

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
		if err := rows.Scan(&id, &event.Urid, &id, &event.Name, &event.StartTime, &event.EndTime, &event.Location, &member); err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
		result.Body.Events = append(result.Body.Events, event)
	}
	if err = rows.Err(); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return result, nil
}

func GetAllEventMembers(input *EventSearchUsers, db *sql.DB) (*EventSearchUsersOutput, error) {
	// Проверяем событие на наличие
	if err := isEventExists(input.Urid, db); err != nil {
		return nil, err
	}

	var rows *sql.Rows

	if len(input.Body.SkillsToSearch) > 0 {
		queryPiece := "("

		var args []any
		args = append(args, input.Urid)

		for i, v := range input.Body.SkillsToSearch {
			if i != 0 {
				queryPiece += " AND "
			}
			queryPiece += "skills.skill=$" + strconv.Itoa(i+2)
			args = append(args, v)
		}
		queryPiece += ")"

		rows, err := db.Query(
			"SELECT event_members.member_email, COUNT(*) as count "+
				"FROM event_members INNER JOIN events ON event_members.event_uri=events.urid "+
				"JOIN skills ON event_members.member_email=skills.user_email "+
				"WHERE event_members.event_uri=$1 AND ("+queryPiece+") "+
				"GROUP BY event_members.member_email", args...,
		)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}

		result := new(EventSearchUsersOutput)

		for rows.Next() {
			var member_email string
			var count int
			if err := rows.Scan(&member_email, &count); err != nil {
				return nil, huma.Error422UnprocessableEntity(err.Error())
			}
			if count == len(input.Body.SkillsToSearch) {
				continue
			}

			user, err := utils.GetUserShortInfo(member_email, db)
			if err != nil {
				return nil, err
			}

			result.Body.Users = append(result.Body.Users, user)
		}
		if err := rows.Err(); err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}

		return result, nil

	}

	rows, err := db.Query(
		"SELECT event_members.member_email FROM event_members "+
			"WHERE event_members.event_uri=$1", input.Urid,
	)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	result := new(EventSearchUsersOutput)

	for rows.Next() {
		var member_email string
		if err := rows.Scan(&member_email); err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}

		user, err := utils.GetUserShortInfo(member_email, db)
		if err != nil {
			return nil, err
		}

		result.Body.Users = append(result.Body.Users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return result, nil
}
