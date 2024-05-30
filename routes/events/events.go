package events

import (
	"context"
	"database/sql"
	"hackaton-jam-back/controllers/events"
	"hackaton-jam-back/controllers/utils"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func Route(api huma.API, db *sql.DB) {
	huma.Register(api, huma.Operation{
		OperationID: "get-last-events",
		Method:      http.MethodGet,
		Path:        "/api/events/last",
		Summary:     "Получить последние события",
		Tags:        []string{"События"},
	}, func(ctx context.Context, input *struct {
		Count int `query:"count" default:"20" example:"20" doc:"Количество событий на страницу"`
		Page  int `query:"page" default:"0" example:"0" doc:"Страница"`
	}) (*events.GetEventsOutput, error) {
		return events.GetAllEvents(input.Count, input.Page, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-event-info",
		Method:      http.MethodGet,
		Path:        "/api/event/{urid}",
		Summary:     "Получить полную информацию события",
		Tags:        []string{"События"},
	}, func(ctx context.Context, input *struct {
		Urid string `path:"urid" doc:"Urid события"`
	}) (*events.FullEventOutput, error) {
		return events.GetFullEventInfo(input.Urid, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-event",
		Method:      http.MethodPost,
		Path:        "/api/event/create",
		Summary:     "Создать событие (только для организаторов)",
		Tags:        []string{"События"},
	}, func(ctx context.Context, input *events.EventCreationInput) (*events.FullEventOutput, error) {
		return events.CreateEvent(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "edit-event",
		Method:      http.MethodPatch,
		Path:        "/api/event/{urid}",
		Summary:     "Редактировать событие",
		Tags:        []string{"События"},
	}, func(ctx context.Context, input *events.EventEditInput) (*events.FullEventOutput, error) {
		return events.EditEvent(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "del-event",
		Method:      http.MethodDelete,
		Path:        "/api/event/{urid}",
		Summary:     "Удалить событие",
		Tags:        []string{"События"},
	}, func(ctx context.Context, input *events.EventDeleteInput) (*events.DeleteEventOutput, error) {
		return events.DeleteEvent(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "join-event",
		Method:      http.MethodPost,
		Path:        "/api/event/{urid}/join",
		Summary:     "Присоединиться к событию",
		Tags:        []string{"События и пользователи"},
	}, func(ctx context.Context, input *events.EventJoinExitInput) (*events.EventJoinExitOutput, error) {
		return events.JoinEvent(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "exit-event",
		Method:      http.MethodPost,
		Path:        "/api/event/{urid}/exit",
		Summary:     "Выйти из события",
		Tags:        []string{"События и пользователи"},
	}, func(ctx context.Context, input *events.EventJoinExitInput) (*events.EventJoinExitOutput, error) {
		return events.ExitEvent(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-user-events",
		Method:      http.MethodGet,
		Path:        "/api/user-events/{email}",
		Summary:     "События пользователя",
		Tags:        []string{"События и пользователи"},
	}, func(ctx context.Context, input *struct {
		Email string `path:"email" example:"thatmaidguy@ya.ru" doc:"E-mail пользователя"`
	}) (*events.UserEventsOutput, error) {
		return events.GetAllJoinedEvents(input.Email, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-users",
		Method:      http.MethodPost,
		Path:        "/api/event/{email}/search-member",
		Summary:     "Поиск людей по навыкам",
		Description: "Ищет пользователей по критериям илли выводит всех участвующих людей",
		Tags:        []string{"События и пользователи"},
	}, func(ctx context.Context, input *events.EventSearchUsers) (*events.EventSearchUsersOutput, error) {
		return events.GetAllEventMembers(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-curr-user-events",
		Method:      http.MethodGet,
		Path:        "/api/user-events",
		Summary:     "События текущего пользователя",
		Tags:        []string{"События и пользователи"},
	}, func(ctx context.Context, input *utils.JustAccessTokenInput) (*events.UserEventsOutput, error) {
		// Проверить можем ли присоединится к меро?
		user, err := utils.GetUserEmailByToken(input.Body.Token, db)
		if err != nil {
			return nil, huma.Error403Forbidden("Пользователь не найден")
		}

		return events.GetAllJoinedEvents(user.Email, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "tag-add-event",
		Method:      http.MethodPut,
		Path:        "/api/event/{urid}/tags",
		Summary:     "Добавить тэги",
		Tags:        []string{"Теги событий"},
	}, func(ctx context.Context, input *events.EventTagAddDelInput) (*events.FullEventOutput, error) {
		return events.AddEventTags(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "tag-del-event",
		Method:      http.MethodDelete,
		Path:        "/api/event/{urid}/tags",
		Summary:     "Удалить тэги",
		Tags:        []string{"Теги событий"},
	}, func(ctx context.Context, input *events.EventTagAddDelInput) (*events.FullEventOutput, error) {
		return events.DelEventTags(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "partners-add-event",
		Method:      http.MethodPut,
		Path:        "/api/event/{urid}/partners",
		Summary:     "Добавить партнеров",
		Tags:        []string{"Партнеры событий"},
	}, func(ctx context.Context, input *events.EventPartnersAddDelInput) (*events.FullEventOutput, error) {
		return events.AddEventPartners(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "partners-del-event",
		Method:      http.MethodDelete,
		Path:        "/api/event/{urid}/partners",
		Summary:     "Удалить партнеров",
		Tags:        []string{"Партнеры событий"},
	}, func(ctx context.Context, input *events.EventPartnersAddDelInput) (*events.FullEventOutput, error) {
		return events.DelEventPartners(input, db)
	})
}
