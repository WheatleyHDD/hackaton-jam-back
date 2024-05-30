package notifications

import (
	"context"
	"database/sql"
	"hackaton-jam-back/controllers/notifications"
	"hackaton-jam-back/controllers/utils"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func Route(api huma.API, db *sql.DB) {
	huma.Register(api, huma.Operation{
		OperationID: "get-notifications",
		Method:      http.MethodGet,
		Path:        "/api/notifications",
		Summary:     "Получить уведомления",
		Tags:        []string{"Уведомления"},
	}, func(ctx context.Context, input *utils.JustAccessTokenInput) (*notifications.NotificationsOutput, error) {
		return notifications.GetNotifications(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "del-notifications",
		Method:      http.MethodDelete,
		Path:        "/api/notifications",
		Summary:     "Удалить все уведомления",
		Tags:        []string{"Уведомления"},
	}, func(ctx context.Context, input *utils.JustAccessTokenInput) (*notifications.NotificationsOutput, error) {
		return notifications.DeleteAllNotifications(input, db)
	})
}
