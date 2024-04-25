package profile

import (
	"context"
	"database/sql"
	"hackaton-jam-back/controllers/profile"
	"hackaton-jam-back/controllers/utils"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func Route(api huma.API, db *sql.DB) {
	huma.Register(api, huma.Operation{
		OperationID: "get-profile-by-username",
		Method:      http.MethodGet,
		Path:        "/api/profile/{username}",
		Summary:     "Получить профиль пользователя",
		Tags:        []string{"Профили"},
	}, func(ctx context.Context, input *struct {
		Username string `path:"username" maxLength:"30" example:"ThatMaidGuy" doc:"Никнейм пользователя"`
	}) (*profile.ProfileOutput, error) {
		return profile.GetProfile(input.Username, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-profile",
		Method:      http.MethodGet,
		Path:        "/api/profile",
		Summary:     "Получить профиль текущего пользователя",
		Tags:        []string{"Профили"},
	}, func(ctx context.Context, input *utils.JustAccessTokenInput) (*profile.ProfileOutput, error) {
		return profile.GetCurrentProfile(input, db)
	})
}
