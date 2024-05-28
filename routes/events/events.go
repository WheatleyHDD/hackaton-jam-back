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
		Count int `query:"count" default:"20" example:"20" doc:"Количество мероприятий на страницу"`
		Page  int `query:"page" default:"0" example:"0" doc:"Страница"`
	}) (*events.GetEventsOutput, error) {
		return events.GetAllEvents(input.Count, input.Page, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-event-info",
		Method:      http.MethodGet,
		Path:        "/api/event/{urid}",
		Summary:     "Получить полную информацию мероприятия",
		Tags:        []string{"Мероприятия"},
	}, func(ctx context.Context, input *struct {
		Urid string `path:"urid" doc:"Urid мероприятия"`
	}) (*events.FullEventOutput, error) {
		return events.GetFullEventInfo(input.Urid, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-event",
		Method:      http.MethodPost,
		Path:        "/api/event/create",
		Summary:     "Создать мероприятие (только для организаторов)",
		Tags:        []string{"Мероприятия"},
	}, func(ctx context.Context, input *events.EventCreationInput) (*events.FullEventOutput, error) {
		return events.CreateEvent(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "edit-event",
		Method:      http.MethodPost,
		Path:        "/api/event/{urid}/edit",
		Summary:     "Редактировать мероприятие",
		Tags:        []string{"Мероприятия"},
	}, func(ctx context.Context, input *events.EventEditInput) (*events.FullEventOutput, error) {
		return events.EditEvent(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "edit-event",
		Method:      http.MethodPost,
		Path:        "/api/event/{urid}/join",
		Summary:     "Редактировать мероприятие",
		Tags:        []string{"Мероприятия"},
	}, func(ctx context.Context, input *events.EventEditInput) (*events.FullEventOutput, error) {
		return events.EditEvent(input, db)
	})
}
