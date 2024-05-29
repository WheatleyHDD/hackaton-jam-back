package events

import (
	"database/sql"
	"hackaton-jam-back/controllers/utils"
	"log"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type GetEventsOutput struct {
	Body struct {
		Count  int         `json:"count" doc:"Количество мероприятий всего"`
		Events []EventType `json:"urid" doc:"Список мероприятий"`
	}
}

type EventType struct {
	Urid      string    `json:"urid" example:"example_events" doc:"Ссылка на мероприятие"`
	Name      string    `json:"name" example:"Example GameJam" doc:"Название мероприятия"`
	StartTime time.Time `json:"start_time" doc:"Начало проведения"`
	EndTime   time.Time `json:"end_time" doc:"Конец проведения"`
	Location  string    `json:"location" example:"Свердловская область, г. Екатеринбург" doc:"Место проведения"`
	Icon      string    `json:"icon" doc:"Превью мероприятия"`
	IsIrl     bool      `json:"is_irl" doc:"Очное ли мероприятие?"`
}

func GetEventCount(db *sql.DB) (int, error) {
	row := db.QueryRow("SELECT COUNT(*) AS total_records FROM events")
	var result int

	err := row.Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func GetAllEvents(count int, page int, db *sql.DB) (*GetEventsOutput, error) {
	rows, err := db.Query("SELECT urid, id, name, start_time, end_time, location FROM events ORDER BY id DESC LIMIT $1 OFFSET $2", count, page*count)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity("Проблемки с вызовом SQL")
	}
	defer rows.Close()

	var events []EventType

	for rows.Next() {
		var id int
		var event EventType
		if err := rows.Scan(&event.Urid, &id, &event.Name, &event.StartTime, &event.EndTime, &event.Location); err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
		events = append(events, event)
	}
	if err = rows.Err(); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	event_count, err := GetEventCount(db)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	result := new(GetEventsOutput)
	result.Body.Events = events
	result.Body.Count = event_count

	return result, nil
}

func GetFullEventInfo(urid string, db *sql.DB) (*FullEventOutput, error) {
	return getFullEventInfo(urid, db)
}

func getFullEventInfo(urid string, db *sql.DB) (*FullEventOutput, error) {
	row := db.QueryRow("SELECT * FROM events WHERE urid = $1", urid)
	event := new(FullEventOutput)
	event.Body.Icon = "https://i.imgur.com/b0zqmkj.jpeg"

	var id int

	var prize sql.NullString
	var location sql.NullString
	var requirements sql.NullString
	var partners sql.NullString
	var icon sql.NullString

	err := row.Scan(
		&event.Body.Urid,
		&id,
		&event.Body.Name,
		&event.Body.StartTime,
		&event.Body.EndTime,
		&prize,
		&location,
		&requirements,
		&partners,
		&icon,
		&event.Body.IsIrl,
		&event.Body.TeamRequirementsType,
		&event.Body.TeamRequirementsValue,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	event.Body.Prize = prize.String
	event.Body.Location = location.String
	event.Body.Requirements = requirements.String
	event.Body.Partners = partners.String
	event.Body.Icon = icon.String

	event.Body.Organizators, err = getEventOrganizators(urid, db)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func getEventOrganizators(urid string, db *sql.DB) ([]*Organizators, error) {
	rows, err := db.Query("SELECT organizator_email FROM events WHERE event_uri = $1", urid)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	defer rows.Close()

	var result []*Organizators

	for rows.Next() {
		org := new(Organizators)
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}

		user, err := utils.GetUserUsernameByEmail(email, db)
		if err != nil {
			return nil, huma.Error403Forbidden("Пользователь не найден")
		}

		org.Email = user.Email
		org.Username = user.Username

		result = append(result, org)
	}
	if err = rows.Err(); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return result, nil
}
