package events

import (
	"database/sql"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type GetEventsOutput struct {
	Body struct {
		Count  int
		Events []EventType
	}
}

type EventType struct {
	Urid      string    `json:"urid" example:"example_events" doc:"Ссылка на мероприятие"`
	Name      string    `json:"name" example:"Example GameJam" doc:"Название мероприятия"`
	StartTime time.Time `json:"start_time" doc:"Начало проведения"`
	EndTime   time.Time `json:"end_time" doc:"Конец проведения"`
	Location  string    `json:"location" example:"Свердловская область, г. Екатеринбург" doc:"Место проведения"`
}

func GetEventCount(db *sql.DB) (int, error) {
	row := db.QueryRow("SELECT COUNT(*) AS total_records FROM event")
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
