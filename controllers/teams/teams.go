package teams

import (
	"database/sql"
	"hackaton-jam-back/controllers/utils"

	"github.com/danielgtaylor/huma/v2"
)

type TeamCreationInput struct {
	Body struct {
		Token string `json:"access_token" example:"82a3682d0d56f40a4d088aee08521663" doc:"Токен пользователя"`

		Urid string `json:"urid" maxLength:"30" example:"example_events" doc:"Ссылка на событие"`
		Name string `json:"name" example:"Супер-команда" doc:"Название команды"`
	}
}

type TeamInviteInput struct {
	Id   int64 `path:"id" example:"0" doc:"Идентификатор команды"`
	Body struct {
		Token string `json:"access_token" example:"82a3682d0d56f40a4d088aee08521663" doc:"Токен пользователя"`

		Invitee string `json:"invitee" example:"thatmaidguy@ya.ru" doc:"E-mail приглашаемого"`
		Role    string `json:"role" example:"Аналитик" doc:"Роль приглашаемого"`
	}
}

type TeamKickInput struct {
	Id   int64 `path:"id" example:"0" doc:"Идентификатор команды"`
	Body struct {
		Token string `json:"access_token" example:"82a3682d0d56f40a4d088aee08521663" doc:"Токен пользователя"`

		Email string `json:"email" example:"thatmaidguy@ya.ru" doc:"E-mail выгоняемого"`
	}
}

type TeamChangeNameInput struct {
	Id   int64 `path:"id" example:"0" doc:"Идентификатор команды"`
	Body struct {
		Token string `json:"access_token" example:"82a3682d0d56f40a4d088aee08521663" doc:"Токен пользователя"`

		NewName string `json:"new_name" example:"Новое название" doc:"Новое название для команды"`
	}
}

type TeamChangeMemberRoleInput struct {
	Id   int64 `path:"id" example:"0" doc:"Идентификатор команды"`
	Body struct {
		Token string `json:"access_token" example:"82a3682d0d56f40a4d088aee08521663" doc:"Токен пользователя"`

		Member  string `json:"member" example:"thatmaidguy@ya.ru" doc:"Участник, роль которому нужно изменить"`
		NewRole string `json:"role" example:"Аналитик" doc:"Новая роль"`
	}
}

type TeamInviteAcceptCancelInput struct {
	Body struct {
		Token  string `json:"access_token" example:"82a3682d0d56f40a4d088aee08521663" doc:"Токен пользователя"`
		TeamId int64  `json:"team_id" example:"0" doc:"Идентификатор команды"`
		From   string `json:"from" example:"thatmaidguy@ya.ru" doc:"От кого приглашение"`
	}
}

type MemberInfo struct {
	Role    string               `json:"role" example:"Разработчик" doc:"Роль участника"`
	Pending bool                 `json:"pending" doc:"В ожидании ответа"`
	User    *utils.UserShortInfo `json:"user" doc:"Информация по участнику"`
}

type TeamInviteCancelOutput struct {
	Body struct {
		Success bool `json:"success" doc:"Успешно выполнено!"`
	}
}

type TeamInfoOutput struct {
	Body struct {
		Id         int64         `json:"id" example:"2" doc:"Идентификатор команды"`
		Name       string        `json:"name" example:"Супер-команда" doc:"Название команды"`
		Urid       string        `json:"urid" example:"example_events" doc:"Ссылка на событие, привязанного к команде"`
		Teamleader string        `json:"teamleader" example:"thatmaidguy@ya.ru" doc:"Тимлид (участник, который может собирать людей)"`
		Members    []*MemberInfo `json:"members" doc:"Список участников"`
	}
}

func CreateTeam(input *TeamCreationInput, db *sql.DB) (*TeamInfoOutput, error) {
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil || user.Perms != 0 {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}

	// Проверяем записан ли пользователь на событие
	var member_email string
	if err := db.QueryRow("SELECT member_email FROM event_members WHERE event_uri = $1 AND member_email = $2", input.Body.Urid, user.Email).Scan(&member_email); err != nil {
		if err == sql.ErrNoRows {
			return nil, huma.Error403Forbidden("Ты не записан на это событие")
		}
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	// Создаем команду
	var teamId int64
	if err := db.QueryRow(
		"INSERT INTO teams (event_uri, name, teamleader) VALUES ($1, $2, $3);"+
			"SELECT id FROM teams ORDER BY id DESC LIMIT 1;",
		input.Body.Urid, input.Body.Name, user.Email).Scan(&teamId); err != nil {

		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	db.QueryRow(
		"INSERT INTO teams_members (team_id, member_email, role, pending) VALUES ($1, $2, $3, false)",
		teamId, user.Email, "Тимлидер").Scan()

	return GetTeamInfo(teamId, db)
}

func GetTeamInfo(teamId int64, db *sql.DB) (*TeamInfoOutput, error) {
	info := new(TeamInfoOutput)

	if err := db.QueryRow(
		"SELECT * FROM teams WHERE id = $1", teamId).Scan(
		&info.Body.Id,
		&info.Body.Urid,
		&info.Body.Name,
		&info.Body.Teamleader,
	); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	rows, err := db.Query("SELECT member_email, role FROM teams_members WHERE team_id = $1", teamId)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		memberInfo := new(MemberInfo)
		if err := rows.Scan(&memberInfo.User.Email, &memberInfo.Role); err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}

		memberInfo.User, err = utils.GetUserShortInfo(memberInfo.User.Email, db)
		if err != nil {
			return nil, err
		}

		info.Body.Members = append(info.Body.Members, memberInfo)
	}
	if err = rows.Err(); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return info, nil
}

func InviteUser(input *TeamInviteInput, db *sql.DB) (*TeamInfoOutput, error) {
	// Проверяем, что пользователь тимлид
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil || user.Perms != 0 {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}

	var teamleader string
	if err := db.QueryRow("SELECT teamleader FROM teams WHERE id = $1", input.Id).Scan(&teamleader); err != nil {
		if err == sql.ErrNoRows {
			return nil, huma.Error403Forbidden("Этой команды нет")
		}
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	if teamleader != user.Email {
		return nil, huma.Error403Forbidden("Вы не тимлид команды")
	}

	// Получаем чуть больше инфы о команде
	result, err := GetTeamInfo(input.Id, db)
	if err != nil {
		return nil, err
	}

	// Кидаем приглашение
	db.QueryRow(
		"INSERT INTO notifications (user, team_id, type, from, event_uri) VALUES ($1, $2, 0, $3, $4)",
		input.Body.Invitee, input.Id, user.Email, result.Body.Urid).Scan()

	// Записываем "карандашом"
	db.QueryRow(
		"INSERT INTO teams_members (team_id, member_email, role) VALUES ($1, $2, $3)",
		input.Id, input.Body.Invitee, input.Body.Role).Scan()

	return GetTeamInfo(input.Id, db)
}

func AcceptInvite(input *TeamInviteAcceptCancelInput, db *sql.DB) (*TeamInfoOutput, error) {
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil || user.Perms != 0 {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}

	// Получаем чуть больше инфы о команде
	result, err := GetTeamInfo(input.Body.TeamId, db)
	if err != nil {
		return nil, err
	}

	// Удаляем приглашение
	db.QueryRow("DELETE FROM notifications WHERE user=$1 AND team_id=$2 AND type=0", user.Email, result.Body.Id).Scan()

	// Записываем "ручкой"
	db.QueryRow("UPDATE teams_members SET pending = false WHERE member_email=$1 AND team_id=$2", user.Email, result.Body.Id).Scan()

	// Кидаем уведомление о принятии (2)
	db.QueryRow(
		"INSERT INTO notifications (user, team_id, type, from, event_uri) VALUES ($1, $2, 2, $3, $4)",
		input.Body.From, result.Body.Id, user.Email, result.Body.Urid).Scan()

	return GetTeamInfo(input.Body.TeamId, db)
}

func CancelInvite(input *TeamInviteAcceptCancelInput, db *sql.DB) (*TeamInviteCancelOutput, error) {
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil || user.Perms != 0 {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}

	// Получаем чуть больше инфы о команде
	teamInfo, err := GetTeamInfo(input.Body.TeamId, db)
	if err != nil {
		return nil, err
	}

	// Удаляем приглашение
	db.QueryRow("DELETE FROM notifications WHERE user=$1 AND team_id=$2 AND type=0", user.Email, teamInfo.Body.Id).Scan()

	// Зачеркиваем "карандаш"
	db.QueryRow("DELETE FROM teams_members WHERE member_email=$1 AND team_id=$2", user.Email, teamInfo.Body.Id).Scan()

	// Кидаем уведомление о отклонении (1)
	db.QueryRow(
		"INSERT INTO notifications (user, team_id, type, from, event_uri) VALUES ($1, $2, 1, $3, $4)",
		input.Body.From, teamInfo.Body.Id, user.Email, teamInfo.Body.Urid).Scan()

	result := &TeamInviteCancelOutput{Body: struct {
		Success bool "json:\"success\" doc:\"Успешно выполнено!\""
	}{Success: true}}

	return result, nil
}

func KickUser(input *TeamKickInput, db *sql.DB) (*TeamInfoOutput, error) {
	// Проверяем, что пользователь тимлид
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil || user.Perms != 0 {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}

	var teamleader string
	if err := db.QueryRow("SELECT teamleader FROM teams WHERE id = $1", input.Id).Scan(&teamleader); err != nil {
		if err == sql.ErrNoRows {
			return nil, huma.Error403Forbidden("Этой команды нет")
		}
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	if teamleader != user.Email {
		return nil, huma.Error403Forbidden("Вы не тимлид команды")
	}

	// Получаем чуть больше инфы о команде
	result, err := GetTeamInfo(input.Id, db)
	if err != nil {
		return nil, err
	}

	// Кидаем уведомление о сливе олуха
	db.QueryRow(
		"INSERT INTO notifications (user, team_id, type, from, event_uri) VALUES ($1, $2, 3, $3, $4)",
		input.Body.Email, input.Id, user.Email, result.Body.Urid).Scan()

	// Вычеркиваем
	db.QueryRow("DELETE FROM teams_members WHERE member_email=$1 AND team_id=$2", user.Email, result.Body.Id).Scan()

	return GetTeamInfo(input.Id, db)

}

func ChangeTeamName(input *TeamChangeNameInput, db *sql.DB) (*TeamInfoOutput, error) {
	if input.Body.NewName == "" {
		return nil, huma.Error422UnprocessableEntity("Имя команды не должно быть пустым")
	}

	// Проверяем, что пользователь тимлид
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil || user.Perms != 0 {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}

	var teamleader string
	if err := db.QueryRow("SELECT teamleader FROM teams WHERE id = $1", input.Id).Scan(&teamleader); err != nil {
		if err == sql.ErrNoRows {
			return nil, huma.Error403Forbidden("Этой команды нет")
		}
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	if teamleader != user.Email {
		return nil, huma.Error403Forbidden("Вы не тимлид команды")
	}

	// Меняем название
	db.QueryRow(
		"UPDATE teams SET name=$1 false WHERE id=$2",
		input.Body.NewName, input.Id).Scan()
	// Мне лень делать проверку ((

	return GetTeamInfo(input.Id, db)
}

func ChangeRole(input *TeamChangeMemberRoleInput, db *sql.DB) (*TeamInfoOutput, error) {
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil || user.Perms != 0 {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}

	var teamleader string
	if err := db.QueryRow("SELECT teamleader FROM teams WHERE id = $1", input.Id).Scan(&teamleader); err != nil {
		if err == sql.ErrNoRows {
			return nil, huma.Error403Forbidden("Этой команды нет")
		}
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	if teamleader != user.Email {
		return nil, huma.Error403Forbidden("Вы не тимлид команды")
	}

	// Получаем чуть больше инфы о команде
	result, err := GetTeamInfo(input.Id, db)
	if err != nil {
		return nil, err
	}

	// Обновляем роль
	db.QueryRow("UPDATE teams_members SET role = $3 WHERE member_email=$1 AND team_id=$2", input.Body.Member, result.Body.Id, input.Body.NewRole).Scan()

	return GetTeamInfo(input.Id, db)
}
