package views

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// GreetingOutput represents the greeting operation response.
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Привет мир!" doc:"Приветственное сообщение"`
	}
}

func (a *ApiRoutes) Example() {
	// Register GET /api/hello/{name} handler.
	huma.Register(a.Api, huma.Operation{
		OperationID: "get-greeting",
		Method:      http.MethodGet,
		Path:        "/api/hello/{name}",
		Summary:     "Получить \"Привет\"",
		Description: "Получить \"Привет\" кому-то или чему-то",
		Tags:        []string{"Пример", "Привет"},
	}, func(ctx context.Context, input *struct {
		Name string `path:"name" maxLength:"30" example:"мир" doc:"Имя для привета"`
	}) (*GreetingOutput, error) {
		resp := &GreetingOutput{}
		resp.Body.Message = fmt.Sprintf("Привет, %s!", input.Name)
		return resp, nil
	})
}
