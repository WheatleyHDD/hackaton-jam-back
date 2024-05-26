package events

import (
	"context"
	"database/sql"
	"hackaton-jam-back/controllers/events"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func Route(api huma.API, db *sql.DB) {
	huma.Register(api, huma.Operation{
		OperationID: "get-last-events",
		Method:      http.MethodGet,
		Path:        "/api/events/last",
		Summary:     "Получить последние мероприятия",
		Tags:        []string{"Мероприятия"},
	}, func(ctx context.Context, input *struct {
		Count int `path:"count" example:"20" doc:"Количество мероприятий на страницу"`
		Page  int `path:"page" example:"0" doc:"Страница"`
	}) (*events.GetEventsOutput, error) {
		return events.GetAllEvents(input.Count, input.Page, db)
	})
}
