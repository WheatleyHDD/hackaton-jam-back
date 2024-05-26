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

	huma.Register(api, huma.Operation{
		OperationID: "edit-profile",
		Method:      http.MethodPost,
		Path:        "/api/profile/edit",
		Summary:     "Редактировать профиль",
		Description: "Редактирует профиль текущего пользователя",
		Tags:        []string{"Профили"},
	}, func(ctx context.Context, input *profile.EditProfileInput) (*profile.ProfileOutput, error) {
		return profile.EditProfile(input, db)
	})

	/// ======================================
	/// ======================================
	/// ======== Контакты профилей ===========
	/// ======================================
	/// ======================================

	huma.Register(api, huma.Operation{
		OperationID: "get-profile-contacts",
		Method:      http.MethodPost,
		Path:        "/api/profile/{username}/contacts",
		Summary:     "Вывести контакты пользователя",
		Description: "Выводит контакты пользователя",
		Tags:        []string{"Контакты пользователя"},
	}, func(ctx context.Context, input *struct {
		Username string `path:"username" maxLength:"30" example:"ThatMaidGuy" doc:"Никнейм пользователя"`
	}) (*profile.ProfileContactsOutput, error) {
		return profile.GetContacts(input.Username, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "add-profile-contacts",
		Method:      http.MethodPost,
		Path:        "/api/profile/{username}/contacts/add",
		Summary:     "Добавление контакта для пользователя",
		Description: "Добавляет контакт для пользователя",
		Tags:        []string{"Контакты пользователя"},
	}, func(ctx context.Context, input *profile.AddDelProfileContactInput) (*profile.ProfileContactsOutput, error) {
		return profile.AddContact(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "del-profile-contacts",
		Method:      http.MethodPost,
		Path:        "/api/profile/{username}/contacts/del",
		Summary:     "Удаление контакта для пользователя",
		Description: "Удаляет контакт для пользователя",
		Tags:        []string{"Контакты пользователя"},
	}, func(ctx context.Context, input *profile.AddDelProfileContactInput) (*profile.ProfileContactsOutput, error) {
		return profile.DelContact(input, db)
	})

	/// ======================================
	/// ======================================
	/// ========= Навыки профилей ============
	/// ======================================
	/// ======================================
	huma.Register(api, huma.Operation{
		OperationID: "get-profile-skills",
		Method:      http.MethodPost,
		Path:        "/api/profile/{username}/skills",
		Summary:     "Вывести навыки пользователя",
		Description: "Выводит навыки пользователя",
		Tags:        []string{"Навыки пользователя"},
	}, func(ctx context.Context, input *struct {
		Username string `path:"username" maxLength:"30" example:"ThatMaidGuy" doc:"Никнейм пользователя"`
	}) (*profile.ProfileSkillsOutput, error) {
		return profile.GetSkills(input.Username, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "add-profile-skills",
		Method:      http.MethodPost,
		Path:        "/api/profile/{username}/skills/add",
		Summary:     "Добавление навыка для пользователя",
		Description: "Добавляет навыки для пользователя",
		Tags:        []string{"Навыки пользователя"},
	}, func(ctx context.Context, input *profile.AddDelProfileSkillsInput) (*profile.ProfileSkillsOutput, error) {
		return profile.AddSkill(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "del-profile-skills",
		Method:      http.MethodPost,
		Path:        "/api/profile/{username}/skills/del",
		Summary:     "Удаление навыка для пользователя",
		Description: "Удаляет навык для пользователя",
		Tags:        []string{"Навыки пользователя"},
	}, func(ctx context.Context, input *profile.AddDelProfileSkillsInput) (*profile.ProfileSkillsOutput, error) {
		return profile.DelSkill(input, db)
	})

	huma.Register(api, huma.Operation{
		OperationID: "search-profile-skills",
		Method:      http.MethodPost,
		Path:        "/api/skills/search",
		Summary:     "Поиск навыков пользователя",
		Description: "Поиск навыков пользователя",
		Tags:        []string{"Навыки пользователя"},
	}, func(ctx context.Context, input *profile.SkillsSearchInput) (*profile.SkillsSearchOutput, error) {
		return profile.GetSkillsByName(input, db)
	})
}
