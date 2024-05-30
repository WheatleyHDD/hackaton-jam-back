package notifications

import (
	"database/sql"
	"hackaton-jam-back/controllers/utils"

	"github.com/danielgtaylor/huma/v2"
)

type Notify struct {
	NotifyType int                  `json:"notify_type" example:"0" doc:"Тип уведомления (0 - приглашение в команду, 1 - отклонение приглашения, 2 - принятие приглашения, 3 - при кике с команды)"`
	From       *utils.UserShortInfo `json:"from" doc:"От кого уведомление"`
	TeamId     int64                `json:"team_id" doc:"Айдишник команды, чтобы принять приглашение"`
	EventUri   string               `json:"event_urid" doc:"Ссылка на мероприятие"`
}

type NotificationsOutput struct {
	Body struct {
		Notifications []*Notify `json:"notifications" doc:"Лист "`
	}
}

func getNotifys(email string, db *sql.DB) (*NotificationsOutput, error) {
	rows, err := db.Query("SELECT team_id, type, from, event_uri FROM notifications WHERE user = $1", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := new(NotificationsOutput)

	for rows.Next() {
		notify := new(Notify)
		var e string
		if err := rows.Scan(&notify.TeamId, &notify.NotifyType, &e, &notify.EventUri); err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}

		notify.From, err = utils.GetUserShortInfo(e, db)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}

		result.Body.Notifications = append(result.Body.Notifications, notify)
	}
	if err = rows.Err(); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return result, nil
}

func GetNotifications(input *utils.JustAccessTokenInput, db *sql.DB) (*NotificationsOutput, error) {
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}

	return getNotifys(user.Email, db)
}

func DeleteAllNotifications(input *utils.JustAccessTokenInput, db *sql.DB) (*NotificationsOutput, error) {
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}

	rows, err := db.Query("DELETE FROM notifications WHERE user = $1 AND type <> 0", user.Email)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	defer rows.Close()

	return getNotifys(user.Email, db)
}
