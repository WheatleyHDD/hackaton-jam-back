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

type MemberInfo struct {
	Role string               `json:"role" example:"Разработчик" doc:"Роль участника"`
	User *utils.UserShortInfo `json:"user" doc:"Информация по участнику"`
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
		"INSERT INTO teams_members (team_id, member_email, role) VALUES ($1, $2, $3)",
		teamId, user.Email, "Тимлидер").Scan()

	return getTeamInfo(teamId, db)
}

func getTeamInfo(teamId int64, db *sql.DB) (*TeamInfoOutput, error) {
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
