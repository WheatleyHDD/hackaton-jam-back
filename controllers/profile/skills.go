package profile

import (
	"database/sql"
	"hackaton-jam-back/controllers/utils"

	"github.com/danielgtaylor/huma/v2"
)

type AddDelProfileSkillsInput struct {
	Username string `path:"username" maxLength:"30" example:"ThatMaidGuy" doc:"Никнейм пользователя"`
	Body     struct {
		Token string `json:"access_token" example:"82a3682d0d56f40a4d088aee08521663" doc:"Токен пользователя"`
		Skill string `json:"contact_link" example:"C#" doc:"Навык"`
	}
}

type SkillsSearchInput struct {
	Body struct {
		Token       string `json:"access_token" example:"82a3682d0d56f40a4d088aee08521663" doc:"Токен пользователя"`
		SearchValue string `json:"search_value" doc:"Используемый критерий поиска навыков"`
	}
}

type ProfileSkillsOutput struct {
	Body struct {
		Skills []string `json:"skills" doc:"Список навыков пользователя"`
	}
}

type SkillsSearchOutput struct {
	Body struct {
		SearchValue string   `json:"search_value" doc:"Используемый критерий поиска"`
		Skills      []string `json:"skills" doc:"Список навыков навыков"`
	}
}

// / ==============================================
// / ==============================================
// / ====== Список навыков для пользователя =======
// / ==============================================
// / ==============================================
func GetSkills(username string, db *sql.DB) (*ProfileSkillsOutput, error) {
	user, err := utils.GetUserEmailByUsername(username, db)
	if err != nil {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}
	return GetSkillsByEmail(user.Email, db)
}

func GetSkillsByEmail(email string, db *sql.DB) (*ProfileSkillsOutput, error) {
	rows, err := db.Query("SELECT skill FROM skills WHERE user_email = $1", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []string

	for rows.Next() {
		var skill string
		if err := rows.Scan(&skill); err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	result := new(ProfileSkillsOutput)
	result.Body.Skills = skills

	return result, nil
}

func AddSkill(input *AddDelProfileSkillsInput, db *sql.DB) (*ProfileSkillsOutput, error) {
	// Проверить можем ли менять?
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}
	if input.Username != user.Username && user.Perms != 10 {
		return nil, huma.Error403Forbidden("Нет прав")
	}

	// Получить уже имеющиеся навыки
	existingSkills, err := GetSkillsByEmail(user.Email, db)
	if err != nil {
		return nil, huma.Error403Forbidden(err.Error())
	}

	// Проверить на существование
	for _, c := range existingSkills.Body.Skills {
		if c == input.Body.Skill {
			return nil, huma.Error403Forbidden("Такой навык уже есть")
		}
	}

	// Добавить контакт
	_, err = db.Query("INSERT INTO skills (user_email, skill) VALUES ($1, $2)", user.Email, input.Body.Skill)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	existingSkills.Body.Skills = append(existingSkills.Body.Skills, input.Body.Skill)

	return existingSkills, nil
}

func DelSkill(input *AddDelProfileSkillsInput, db *sql.DB) (*ProfileSkillsOutput, error) {
	// Проверить можем ли менять?
	user, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}
	if input.Username != user.Username && user.Perms != 10 {
		return nil, huma.Error403Forbidden("Нет прав")
	}

	// Удалить ссылку
	_, err = db.Query("DELETE FROM skills WHERE user_email=$1 AND skill=$2", user.Email, input.Body.Skill)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity("Похоже такого навыка у пользователя не было")
	}

	return GetSkillsByEmail(user.Email, db)
}

func GetSkillsByName(input *SkillsSearchInput, db *sql.DB) (*SkillsSearchOutput, error) {
	// Найти пользователя
	_, err := utils.GetUserEmailByToken(input.Body.Token, db)
	if err != nil {
		return nil, huma.Error403Forbidden("Пользователь не найден")
	}

	// Найти навыки
	rows, err := db.Query("SELECT skill, COUNT(*) AS count FROM skills WHERE skill LIKE $1 || '%' GROUP BY skill ORDER BY count DESC LIMIT 10;", input.Body.SearchValue)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity("Проблемки с вызовом SQL")
	}
	defer rows.Close()

	var skills []string

	for rows.Next() {
		var skill string
		var count int
		if err := rows.Scan(&skill, &count); err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
		skills = append(skills, skill)
	}
	if err = rows.Err(); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	result := new(SkillsSearchOutput)
	result.Body.SearchValue = input.Body.SearchValue
	result.Body.Skills = skills

	return result, nil
}
