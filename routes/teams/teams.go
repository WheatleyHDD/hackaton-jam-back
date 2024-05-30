package teams

import (
	"context"
	"database/sql"
	"hackaton-jam-back/controllers/teams"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func Route(api huma.API, db *sql.DB) {
	huma.Register(api, huma.Operation{
		OperationID: "create-team",
		Method:      http.MethodPost,
		Path:        "/api/team/сreate",
		Summary:     "Создать команду",
		Tags:        []string{"Команды"},
	}, func(ctx context.Context, input *teams.TeamCreationInput) (*teams.TeamInfoOutput, error) {
		return teams.CreateTeam(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-team",
		Method:      http.MethodGet,
		Path:        "/api/team/{id}",
		Summary:     "Получить информацию о команде",
		Tags:        []string{"Команды"},
	}, func(ctx context.Context, input *struct {
		Id int64 `path:"id" example:"0" doc:"Идентификатор команды"`
	}) (*teams.TeamInfoOutput, error) {
		return teams.GetTeamInfo(input.Id, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "team-invite-user",
		Method:      http.MethodPost,
		Path:        "/api/team/{id}/invite",
		Summary:     "Пригласить пользователя",
		Tags:        []string{"Команды"},
	}, func(ctx context.Context, input *teams.TeamInviteInput) (*teams.TeamInfoOutput, error) {
		return teams.InviteUser(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "team-invite-accept",
		Method:      http.MethodPut,
		Path:        "/api/invite/accept",
		Summary:     "Принять приглашение",
		Tags:        []string{"Команды"},
	}, func(ctx context.Context, input *teams.TeamInviteAcceptCancelInput) (*teams.TeamInfoOutput, error) {
		return teams.AcceptInvite(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "team-invite-cancel",
		Method:      http.MethodPut,
		Path:        "/api/invite/cancel",
		Summary:     "Отклонить приглашение",
		Tags:        []string{"Команды"},
	}, func(ctx context.Context, input *teams.TeamInviteAcceptCancelInput) (*teams.TeamInviteCancelOutput, error) {
		return teams.CancelInvite(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "team-kick",
		Method:      http.MethodPut,
		Path:        "/api/team/{id}/kick",
		Summary:     "Выгнать пользователя из команды",
		Tags:        []string{"Команды"},
	}, func(ctx context.Context, input *teams.TeamKickInput) (*teams.TeamInfoOutput, error) {
		return teams.KickUser(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "team-change-name",
		Method:      http.MethodPatch,
		Path:        "/api/team/{id}/change-name",
		Summary:     "Изменить название команды",
		Tags:        []string{"Команды"},
	}, func(ctx context.Context, input *teams.TeamChangeNameInput) (*teams.TeamInfoOutput, error) {
		return teams.ChangeTeamName(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "team-change-role",
		Method:      http.MethodPatch,
		Path:        "/api/team/{id}/member-role",
		Summary:     "Изменить роль участника в команде",
		Tags:        []string{"Команды"},
	}, func(ctx context.Context, input *teams.TeamChangeMemberRoleInput) (*teams.TeamInfoOutput, error) {
		return teams.ChangeRole(input, db)
	})
}
