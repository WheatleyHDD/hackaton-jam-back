package auth

import (
	"context"
	"database/sql"
	"hackaton-jam-back/controllers/auth"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func Route(api huma.API, db *sql.DB) {
	huma.Register(api, huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        "/api/login",
		Summary:     "Вход в аккаунт",
		Tags:        []string{"Авторизация"},
	}, func(ctx context.Context, input *auth.LoginInput) (*auth.LoginResponseOutput, error) {
		return auth.Login(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "register",
		Method:      http.MethodPost,
		Path:        "/api/register",
		Summary:     "Регистрация аккаунта",
		Tags:        []string{"Авторизация"},
	}, func(ctx context.Context, input *auth.RegisterInput) (*auth.LoginResponseOutput, error) {
		return auth.Register(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID:   "logout",
		Method:        http.MethodPost,
		Path:          "/api/logout",
		Summary:       "Выход из аккаунта",
		Tags:          []string{"Авторизация"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *auth.LogoutInput) (*struct{}, error) {
		return auth.Logout(input, db)
	})
}
